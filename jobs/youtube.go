package jobs

import (
	"errors"
	"fmt"
	"log"
	"math/rand"
	"regexp"
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
	var shipNameRegex = regexp.MustCompile(`(SN\d+)`)

	var (
		lastTweetedURL string
	)

	for {
		liveVideo, err := scrapers.YouTubeLive(spaceXLiveURL)
		if err == nil {
			if liveVideo.VideoID != "" && (match.StarshipText(liveVideo.Title, false) || match.StarshipText(liveVideo.ShortDescription, false)) {
				// Get the video URL
				liveURL := liveVideo.URL()

				if liveURL != lastTweetedURL {

					// See if we can get the starship name, but we tweet without it anyway
					var shipName = shipNameRegex.FindString(liveVideo.Title)
					if shipName != "" {
						shipName = " #" + shipName
					}

					tweetText := fmt.Sprintf("It's hoppening! SpaceX #Starship%s stream is live\n%s", shipName, liveURL)

					// OK, we can tweet this

					t, _, err := client.Statuses.Update(tweetText, nil)
					if err == nil {
						log.Println("[Twitter] Tweeted", util.TweetURL(t))

						// make sure we don't tweet this again
						lastTweetedURL = liveURL
					} else {
						log.Println("[Twitter] Error while tweeting livestream update:", err.Error())
					}
				}
			}
		} else {
			if !errors.Is(err, scrapers.ErrNotLive) {
				log.Println("[YouTube] Unexpected error while scraping YouTube live:", err.Error())
			}
		}

		// Wait up to two minutes, then check again
		time.Sleep(time.Minute + time.Duration(rand.Intn(60))*time.Second)
	}
}
