package consumer

import (
	"encoding/json"
	"errors"
	"os"
	"strings"
	"time"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/xarantolus/spacex-hop-bot/match"
	"github.com/xarantolus/spacex-hop-bot/scrapers"
	"github.com/xarantolus/spacex-hop-bot/util"
)

func (p *Processor) cleanup(save bool) {
	if p.test {
		return
	}

	var changedLinks = false

	for k, d := range p.seenLinks {
		if time.Since(d) > seenLinkDelay {
			// No point in keeping this info
			delete(p.seenLinks, k)
			changedLinks = true
		}
	}

	if save && changedLinks {
		util.LogError(util.SaveJSON(articlesFilename, p.seenLinks), "saving links after cleanup")
	}
}

func isElonTweet(t match.TweetWrapper) bool {
	return t.User != nil && strings.EqualFold(t.User.ScreenName, "elonmusk")
}
func isSpaceXTweet(t match.TweetWrapper) bool {
	return t.User != nil && strings.EqualFold(t.User.ScreenName, "SpaceX")
}

func sameUser(t1, t2 *twitter.Tweet) bool {
	return t1.User != nil && t2.User != nil &&
		(t1.User.ID == t2.User.ID ||
			strings.EqualFold(t1.User.ScreenName, t2.User.ScreenName))
}

func (p *Processor) isStarshipTweet(t match.TweetWrapper) bool {
	// At first, we of course need to match some keywords
	if !p.matcher.StarshipTweet(t) {
		t.Log("tweet not considered a starship tweet by the matcher")
		return false
	}

	// However, we don't want reaction gifs
	if isReactionGIF(&t.Tweet) {
		t.Log("tweet has a reaction gif")
		return false
	}

	// Replies to other people should be filtered
	if p.isReply(&t.Tweet) {
		t.Log("tweet is a reply to other people")
		return false
	}

	// If it's a question, we ignore it, except if at the launch site OR has media OR is a pad announcement
	if isQuestion(&t.Tweet) &&
		!(match.IsAtSpaceXSite(&t.Tweet) ||
			hasMedia(&t.Tweet) ||
			match.IsPadAnnouncement(t.Text())) {
		t.Log("tweet is a question not (at a spacex site | has media | pad announcement)")
		return false
	}

	// Anything else should be OK
	return true
}

// isReply returns if the given tweet is a reply to another user
func (p *Processor) isReply(t *twitter.Tweet) bool {
	if t.User == nil || t.InReplyToStatusID == 0 {
		return false
	}

	if !(t.User.ID == t.InReplyToUserID || strings.EqualFold(t.User.ScreenName, t.InReplyToScreenName)) {
		return true
	}

	t, err := p.client.LoadStatus(t.InReplyToStatusID)
	if err != nil {
		// If something goes wrong, we just assume it is a reply
		return true
	}

	return p.isReply(t)
}

func linksToLiveStream(tweet *twitter.Tweet) bool {
	if tweet.Entities == nil {
		return false
	}

	for _, u := range tweet.Entities.Urls {
		if !strings.Contains(u.ExpandedURL, "youtu") {
			continue
		}
		liveVid, err := scrapers.YouTubeLive(u.ExpandedURL)
		if errors.Is(err, scrapers.ErrNoVideo) ||
			util.LogError(err, "scraping youtube live at %q", u.URL) {
			continue
		}
		if liveVid.IsLive || liveVid.IsUpcoming {
			return true
		}
	}

	return false
}

func isQuestion(tweet *twitter.Tweet) bool {
	txt := tweet.Text()
	// Make sure we don't match "?" in an URL
	txt = urlRegex.ReplaceAllString(txt, "")
	return strings.Contains(txt, "?")
}

func isReactionGIF(tweet *twitter.Tweet) bool {
	if tweet.ExtendedEntities == nil || len(tweet.ExtendedEntities.Media) != 1 {
		return false
	}

	// Type of a GIF is animated_gif
	return strings.Contains(tweet.ExtendedEntities.Media[0].Type, "gif")
}

func hasMedia(tweet *twitter.Tweet) bool {
	return tweet.Entities != nil && len(tweet.Entities.Media) > 0 || tweet.ExtendedEntities != nil && len(tweet.ExtendedEntities.Media) > 0
}

// saveTweet appends the given tweet to a JSON file for later inspections, especially in case of wrong retweets
func (p *Processor) saveTweet(tweet *twitter.Tweet) {
	f, err := os.OpenFile("retweeted.ndjson", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		util.LogError(err, "open tweet file")
		return
	}
	defer f.Close()

	err = json.NewEncoder(f).Encode(tweet)
	if err != nil {
		util.LogError(err, "encoding tweet json")
	}
}

// isTagsOnly returns if the given text only contains words that start with a tag or hashtag
func isTagsOnly(text string) bool {
	var fields = strings.Fields(text)

	if len(fields) == 0 {
		return false
	}

	for _, f := range fields {
		if len(fields) == 0 || f[0] == '#' || f[0] == '@' {
			continue
		}

		return false
	}

	return true
}
