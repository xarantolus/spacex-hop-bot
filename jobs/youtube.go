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
		lastTweetedURL      string
		lastTweetedUpcoming bool
	)

	for {
		liveVideo, err := scrapers.YouTubeLive(spaceXLiveURL)
		if err != nil {
			if !errors.Is(err, scrapers.ErrNoVideo) {
				log.Println("[YouTube] Unexpected error while scraping YouTube live:", err.Error())
			}
		}

		// This combines both the error and any other case where no link can be generated
		if liveVideo.VideoID == "" {
			goto sleep
		}

		// If we have interesting video info
		if match.StarshipText(liveVideo.Title, true) || match.StarshipText(liveVideo.ShortDescription, true) {
			// Get the video URL
			liveURL := liveVideo.URL()

			// Check if we already tweeted this before
			if liveURL == lastTweetedURL && liveVideo.IsUpcoming == lastTweetedUpcoming {
				log.Printf("[YouTube] Already tweeted stream link %s with title %q", liveVideo.URL(), liveVideo.Title)
				goto sleep
			}

			// See if we can get the starship name, but we tweet without it anyway
			var shipName = shipNameRegex.FindString(strings.ToUpper(liveVideo.Title))

			// Depending on what flies, we tweet different text
			var tweetText string
			switch {
			// Upcoming video
			case strings.HasPrefix(shipName, "BN") && liveVideo.IsUpcoming:
				tweetText = fmt.Sprintf("SpaceX #Starship Booster #SuperHeavy #%s stream was posted to YouTube, likely starting soon\n#WenHop\n%s", shipName, liveURL)
			case strings.HasPrefix(shipName, "SN") && liveVideo.IsUpcoming:
				tweetText = fmt.Sprintf("SpaceX #Starship #%s stream was posted to YouTube, likely starting soon\n#WenHop\n%s", shipName, liveURL)

				// If it's not upcoming, it's likely live
			case strings.HasPrefix(shipName, "BN"):
				tweetText = fmt.Sprintf("It's hoppening! SpaceX #Starship Booster #SuperHeavy #%s stream is live\n%s", shipName, liveURL)
			case strings.HasPrefix(shipName, "SN"):
				tweetText = fmt.Sprintf("It's hoppening! SpaceX #Starship #%s stream is live\n%s", shipName, liveURL)

				// If we don't have a SN/BN prefix, we ignore that and tweet anyways
			case liveVideo.IsUpcoming:
				tweetText = fmt.Sprintf("SpaceX #Starship stream was posted to YouTube, likely starting soon\n#WenHop\n%s", liveURL)
			case liveVideo.IsLive:
				tweetText = fmt.Sprintf("It's hoppening! SpaceX #Starship stream is live\n%s", liveURL)
			default:
				log.Printf("Got stream with title %q and link %s (isUpcoming=%v, isLive=%v), but cannot generate a nice tweet text\n", liveVideo.Title, liveURL, liveVideo.IsUpcoming, liveVideo.IsLive)
				goto sleep
			}

			// Now tweet the text we generated
			t, _, err := client.Statuses.Update(tweetText, nil)
			if err != nil {
				log.Println("[Twitter] Error while tweeting livestream update:", err.Error())
				goto sleep
			}

			// make sure we don't tweet this again
			lastTweetedURL = liveURL
			lastTweetedUpcoming = liveVideo.IsUpcoming

			log.Println("[Twitter] Tweeted", util.TweetURL(t))
		}

	sleep:
		// Wait up to two minutes, then check again
		time.Sleep(time.Minute + time.Duration(rand.Intn(60))*time.Second)
	}
}
