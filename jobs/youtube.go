package jobs

import (
	"errors"
	"fmt"
	"log"
	"math/rand"
	"regexp"
	"strings"
	"time"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/xarantolus/spacex-hop-bot/match"
	"github.com/xarantolus/spacex-hop-bot/scrapers"
	"github.com/xarantolus/spacex-hop-bot/util"
)

// CheckYouTubeLive checks SpaceX's youtube live stream every 1-2 minutes and tweets if there is a starship launch stream
func CheckYouTubeLive(client *twitter.Client, user *twitter.User) {
	defer panic("for some reason, the youtube live checker stopped running even though it never should")

	log.Println("[YouTube] Watching SpaceX channel for live Starship streams")

	const spaceXLiveURL = "https://www.youtube.com/spacex/live"

	// This finds strings like SN11, BN2 etc.
	var shipNameRegex = regexp.MustCompile(`((?:SN|BN)\s*\d+)`)

	var (
		lastTweetedURL string
	)

	for {
		liveVideo, err := scrapers.YouTubeLive(spaceXLiveURL)
		if err != nil {
			if !errors.Is(err, scrapers.ErrNotLive) {
				log.Println("[YouTube] Unexpected error while scraping YouTube live:", err.Error())
			}
		}

		// This combines both the error and any other case where no link can be generated
		if liveVideo.VideoID == "" {
			goto sleep
		}

		// If we have interesting video info
		if match.StarshipText(liveVideo.Title, false) || match.StarshipText(liveVideo.ShortDescription, false) {
			// Get the video URL
			liveURL := liveVideo.URL()

			// Check if we already tweeted this before
			if liveURL == lastTweetedURL {
				log.Printf("[YouTube] Not Tweeting stream link %s with title %q", liveVideo.URL(), liveVideo.Title)
				goto sleep
			}

			// Generate Tweet text
			tweetText := fmt.Sprintf("It's hoppening! SpaceX #Starship stream is live\n%s", liveURL)

			// See if we can get the starship name, but we tweet without it anyway
			var shipName = shipNameRegex.FindString(strings.ToUpper(liveVideo.Title))
			if shipName != "" {
				// Booster or Starship?
				if strings.HasPrefix(shipName, "BN") {
					tweetText = fmt.Sprintf("It's hoppening! SpaceX #Starship Booster #SuperHeavy #%s stream is live\n%s", shipName, liveURL)
				} else {
					tweetText = fmt.Sprintf("It's hoppening! SpaceX #Starship #%s stream is live\n%s", shipName, liveURL)
				}
			}

			// Now tweet the text we generated
			t, _, err := client.Statuses.Update(tweetText, nil)
			if err != nil {
				log.Println("[Twitter] Error while tweeting livestream update:", err.Error())
				goto sleep
			}

			// make sure we don't tweet this again
			lastTweetedURL = liveURL

			log.Println("[Twitter] Tweeted", util.TweetURL(t))
		}

	sleep:
		// Wait up to two minutes, then check again
		time.Sleep(time.Minute + time.Duration(rand.Intn(60))*time.Second)
	}
}
