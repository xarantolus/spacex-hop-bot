package jobs

import (
	"errors"
	"fmt"
	"log"
	"math/rand"
	"strings"
	"time"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/docker/go-units"
	"github.com/xarantolus/spacex-hop-bot/match"
	"github.com/xarantolus/spacex-hop-bot/scrapers"
	"github.com/xarantolus/spacex-hop-bot/util"
)

// CheckYouTubeLive checks SpaceX's youtube live stream every 1-2 minutes and tweets if there is a starship launch stream
func CheckYouTubeLive(client *twitter.Client, user *twitter.User, linkChan <-chan string) {
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
		if match.StarshipText(liveVideo.Title, nil) || match.StarshipText(liveVideo.ShortDescription, nil) {
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

			// See if we can get the starship name, but we tweet without it anyway
			var shipName = scrapers.ShipNameRegex.FindString(strings.ToUpper(liveVideo.Title))

			// Depending on what flies, we tweet different text
			var tweetText string

			switch {
			// Upcoming video
			case strings.HasPrefix(shipName, "BN") && liveVideo.IsUpcoming:
				if haveStartTime {
					tweetText = fmt.Sprintf("Upcoming SpaceX #Starship Booster #SuperHeavy #%s stream posted to YouTube, likely starting in %s\n#WenHop\n%s", shipName, strings.ToLower(units.HumanDuration(d)), liveURL)
				} else {
					tweetText = fmt.Sprintf("Upcoming SpaceX #Starship Booster #SuperHeavy #%s stream posted to YouTube, likely starting soon\n#WenHop\n%s", shipName, liveURL)
				}
			case strings.HasPrefix(shipName, "SN") && liveVideo.IsUpcoming:
				if haveStartTime {
					tweetText = fmt.Sprintf("Upcoming SpaceX #Starship #%s stream posted to YouTube, likely starting in %s\n#WenHop\n%s", shipName, strings.ToLower(units.HumanDuration(d)), liveURL)
				} else {
					tweetText = fmt.Sprintf("Upcoming SpaceX #Starship #%s stream posted to YouTube, likely starting soon\n#WenHop\n%s", shipName, liveURL)
				}

				// If it's not upcoming, it's likely live
			case strings.HasPrefix(shipName, "BN"):
				tweetText = fmt.Sprintf("It's hoppening! SpaceX #Starship Booster #SuperHeavy #%s stream is live\n%s", shipName, liveURL)
			case strings.HasPrefix(shipName, "SN"):
				tweetText = fmt.Sprintf("It's hoppening! SpaceX #Starship #%s stream is live\n%s", shipName, liveURL)

				// If we don't have a SN/BN prefix, we ignore that and tweet anyways
			case liveVideo.IsUpcoming:
				if haveStartTime {
					tweetText = fmt.Sprintf("Upcoming SpaceX #Starship stream was posted to YouTube, likely starting in %s\n#WenHop\n%s", liveURL, strings.ToLower(units.HumanDuration(d)))
				} else {
					tweetText = fmt.Sprintf("Upcoming SpaceX #Starship stream was posted to YouTube, likely starting soon\n#WenHop\n%s", liveURL)
				}
			case liveVideo.IsLive:
				tweetText = fmt.Sprintf("It's hoppening! SpaceX #Starship stream is live\n%s", liveURL)
			default:
				log.Printf("Got stream with title %q and link %s (isUpcoming=%v, isLive=%v), but cannot generate a nice tweet text\n", liveVideo.Title, liveURL, liveVideo.IsUpcoming, liveVideo.IsLive)
				goto sleep
			}

			// Now tweet the text we generated
			tweet, _, err := client.Statuses.Update(tweetText, nil)
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
