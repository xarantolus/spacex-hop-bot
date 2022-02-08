package consumer

import (
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/xarantolus/spacex-hop-bot/match"
	"github.com/xarantolus/spacex-hop-bot/scrapers"
	"github.com/xarantolus/spacex-hop-bot/util"
	"mvdan.cc/xurls/v2"
)

var (
	ignoredHosts = map[string]bool{
		"gofundme.com":         true,
		"opensea.io":           true,
		"shop.spreadshirt.com": true,
		"spreadshirt.com":      true,
		"instagram.com":        true,
		"soundcloud.com":       true,
		"blueorigin.com":       true,
		"affinitweet.com":      true,
		"boards.greenhouse.io": true,
		"etsy.com":             true,
		// Most of their articles are paywalled, no additional benefit for retweeting them
		"spaceq.ca": true,
	}

	// This map contains very important URLs that should usually be retweeted
	// map[host]path
	importantURLs = map[string]string{
		// NSF live stream
		"nasaspaceflight.com": "/starbaselive",
		// Road closures
		"cameroncountytx.gov": "/spacex/",

		// HQ media site often linked with image tweets
		"cnunezimages.com": "*",
	}

	highQualityYouTubeStreams = map[string]bool{
		// Do not ignore NASASpaceflight, people often tweet updates with a link to their 24/7 stream
		"UCSUu1lih2RifWkKtDOJdsBA": true,
		// Same for LabPadre
		"UCFwMITSkc1Fms6PoJoh1OUQ": true,
		// SpaceX official channel
		"UCtI0Hodo5o5dUb67FeUjDeA": true,

		// Jessica Kirsh
		"UCpThejfzN2EJiXa2mEwdEUw": true,

		// Starship Gazer
		"UCBVnapKtPTNYl4phaGXxYng": true,
	}

	urlRegex *regexp.Regexp
)

func init() {
	var err error
	urlRegex, err = xurls.StrictMatchingScheme("https|http")
	if err != nil {
		panic("parsing URL regex: " + err.Error())
	}
}

func isImportantURL(uri string) (important bool) {
	parsed, err := url.ParseRequestURI(uri)
	if err != nil {
		return
	}

	host := strings.TrimPrefix(strings.ToLower(parsed.Hostname()), "www.")

	requiredPath, ok := importantURLs[host]
	if !ok {
		return false
	}

	return requiredPath == "*" ||
		strings.Trim(requiredPath, "/") == strings.Trim(parsed.Path, "/")
}

// shouldIgnoreLink returns whether this tweet should be ignored because of a linked article
func (p *Processor) shouldIgnoreLink(tweet match.TweetWrapper) (ignore bool) {
	// Get the text *with* URLs
	var textWithURLs = tweet.TextWithURLs()

	// Find all URLs
	urls := urlRegex.FindAllString(textWithURLs, -1)

	// Now check if any of these URLs is ignored
	for _, u := range urls {
		// Ignore important URLs
		if isImportantURL(u) {
			tweet.Log("URL %q is important", u)
			continue
		}

		var canonical string = u
		if !p.test {
			canonical = util.FindCanonicalURL(u, false)
		}
		if isImportantURL(canonical) {
			tweet.Log("URL %q is important", canonical)
			continue
		}

		parsed, err := url.ParseRequestURI(canonical)
		if err != nil {
			util.LogError(err, "parse canonical for %s", u)
			continue
		}

		// Check if the host is ignored
		host := strings.TrimPrefix(strings.ToLower(parsed.Hostname()), "www.")
		if ignoredHosts[host] {
			tweet.Log("URL %q has ignored host %q", parsed.String(), host)
			return true
		}

		// Don't make requests when we're in a test, it makes no sense to have tests depend on behavior
		// of external websites
		if p.test {
			continue
		}

		if host == "youtube.com" || host == "youtu.be" {
			stream, err := scrapers.YouTubeLive(canonical)
			if err == nil {
				// If we know the channel is good, then we don't ignore their live streams
				if (stream.IsLive || stream.IsUpcoming) && highQualityYouTubeStreams[stream.ChannelID] {
					continue
				}

				// Else, we should of course check if we've seen it before
			}
		}

		// If we retweeted this link in the last 12 hours, we should
		// definitely ignore it
		lastRetweetTime, ok := p.seenLinks[u]
		if ok && time.Since(lastRetweetTime) < seenLinkDelay {
			return true
		}
		lastRetweetTime, ok = p.seenLinks[canonical]
		if ok && time.Since(lastRetweetTime) < seenLinkDelay {
			return true
		}

		// Mark this link as seen, but allow a retweet
		p.seenLinks[u] = time.Now()
		p.seenLinks[canonical] = time.Now()

		p.cleanup(false)

		// Now save it to make sure we still know after a restart
		util.LogError(util.SaveJSON(articlesFilename, p.seenLinks), "saving links")

		return false
	}

	return false
}
