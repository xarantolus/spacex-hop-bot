package jobs

import (
	"bytes"
	"errors"
	"log"
	"math/rand"
	"regexp"
	"strings"
	"text/template"
	"time"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/docker/go-units"
	"github.com/xarantolus/spacex-hop-bot/consumer"
	"github.com/xarantolus/spacex-hop-bot/match"
	"github.com/xarantolus/spacex-hop-bot/scrapers"
	"github.com/xarantolus/spacex-hop-bot/util"
)

// CheckYouTubeLive checks SpaceX's youtube live stream every 1-2 minutes and tweets if there is a starship launch stream
func CheckYouTubeLive(client consumer.TwitterClient, user *twitter.User, matcher *match.StarshipMatcher, linkChan <-chan string) {
	defer panic("for some reason, the youtube live checker stopped running even though it never should")

	log.Println("[YouTube] Watching SpaceX channel for live Starship streams")

	const spaceXLiveURL = "https://www.youtube.com/spacex/live"

	var (
		lastTweetedURL      string
		lastTweetedUpcoming bool

		lastLiveStart time.Time
		lastTweetTime time.Time
	)

	var linkOverwrite string

	for {
		if linkOverwrite == "" {
			linkOverwrite = spaceXLiveURL
		}

		liveVideo, err := scrapers.YouTubeLive(linkOverwrite)
		if err != nil && !errors.Is(err, scrapers.ErrNoVideo) {
			log.Println("[YouTube] Unexpected error while scraping YouTube live:", err.Error())
		}

		linkOverwrite = ""

		if liveVideo.VideoID == "" || err != nil {
			goto sleep
		}

		// If we have interesting video info
		if isStarshipStream(matcher, &liveVideo) {
			// Get the video URL
			liveURL := liveVideo.URL()

			liveStartTime, d, haveStartTime := liveVideo.TimeUntil()

			// Check if we already tweeted this before - but also tweet if we didn't tweet within the last 15 minutes
			if liveURL == lastTweetedURL && liveVideo.IsUpcoming == lastTweetedUpcoming && lastLiveStart.Equal(liveStartTime) && time.Since(lastTweetTime) < tweetInterval(d) {
				log.Printf("[YouTube] Already tweeted stream link %s with title %q", liveVideo.URL(), liveVideo.Title)
				goto sleep
			}

			if haveStartTime {
				lastLiveStart = liveStartTime
			}

			tweetText := describeLiveStream(&liveVideo)

			// Now tweet the text we generated
			tweet, err := client.Tweet(tweetText, nil)
			if err != nil {
				log.Println("[Twitter] Error while tweeting livestream update:", err.Error())
				goto sleep
			}

			// make sure we don't tweet this again
			lastTweetedURL = liveURL
			lastTweetedUpcoming = liveVideo.IsUpcoming
			lastTweetTime = time.Now()

			log.Println("[Twitter] Tweeted", util.TweetURL(tweet))
		}

	sleep:
		// Wait up to two minutes, then check again
		select {
		case <-time.After(time.Minute + time.Duration(rand.Intn(60))*time.Second):
		case linkOverwrite = <-linkChan:
		}
	}
}

func isStarshipStream(matcher *match.StarshipMatcher, liveVideo *scrapers.LiveVideo) bool {
	return matcher.StarshipText(liveVideo.Title, nil, false) || matcher.StarshipText(liveVideo.ShortDescription, nil, false)
}

func tweetInterval(streamStartsIn time.Duration) time.Duration {
	switch {
	case streamStartsIn < time.Hour:
		return 15 * time.Minute
	case streamStartsIn < 4*time.Hour:
		return time.Hour
	default:
		return 2 * time.Hour
	}
}

var expectedRegexes = []*regexp.Regexp{
	regexp.MustCompile(`\b(Star[sS]hip)\b`),
	regexp.MustCompile(`\b(Super\s*[Hh]eavy)\b`),
	regexp.MustCompile(`\b(SN?\d+)\b`),
	regexp.MustCompile(`\b(BN?\d+)\b`),
}

func extractKeywords(title string, description string) (keywords []string) {
	extr := title + "\n " + description

	for _, matcher := range expectedRegexes {
		res := matcher.FindAllString(extr, 1)
		if len(res) == 0 {
			continue
		}

		if !containsIgnoreCase(keywords, res[0]) {
			keywords = append(keywords, res[0])
		}
	}

	return
}

func containsIgnoreCase(s []string, e string) bool {
	for _, a := range s {
		if strings.EqualFold(a, e) {
			return true
		}
	}
	return false
}

const streamTweetTemplate = `{{if .IsUpcoming}}SpaceX live stream starts {{if .HaveStartTime}}in {{.TimeUntil | duration}}{{else}}soon{{end}}:{{else}}SpaceX is now live on YouTube:{{end}}

{{.Title}}
{{$keywords := (keywords .Title .ShortDescription)}}{{with $keywords}}
{{hashtags $keywords}}{{end}}

{{.URL}}`

var (
	tmplFuncs = map[string]interface{}{
		"hashtags": util.HashTagText,
		"keywords": extractKeywords,
		"duration": func(d time.Duration) string {
			return strings.ToLower(units.HumanDuration(d))
		},
	}
	streamTweetTmpl = template.Must(template.New("streamTweetTemplate").Funcs(tmplFuncs).Parse(streamTweetTemplate))
)

func describeLiveStream(v *scrapers.LiveVideo) string {
	var b bytes.Buffer

	t, dur, haveStartTime := v.TimeUntil()

	var data = struct {
		HaveStartTime bool
		TimeUntil     time.Duration
		StartTime     time.Time
		*scrapers.LiveVideo
	}{
		HaveStartTime: haveStartTime,
		TimeUntil:     dur,
		StartTime:     t,
		LiveVideo:     v,
	}

	err := streamTweetTmpl.Execute(&b, data)
	if err != nil {
		panic("executing Tweet template: " + err.Error())
	}

	return strings.TrimSpace(b.String())
}
