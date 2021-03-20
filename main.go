package main

import (
	"errors"
	"flag"
	"log"
	"math/rand"
	"time"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/xarantolus/spacex-hop-bot/bot"
	"github.com/xarantolus/spacex-hop-bot/config"
	"github.com/xarantolus/spacex-hop-bot/match"
	"github.com/xarantolus/spacex-hop-bot/scrapers"
)

var flagConfigFile = flag.String("cfg", "deploy/config.yaml", "Config file path")

func main() {
	flag.Parse()

	// Some stuff depends on randomness
	rand.Seed(time.Now().UnixNano())

	cfg, err := config.Parse(*flagConfigFile)
	if err != nil {
		panic("parsing configuration file: " + err.Error())
	}

	client, user, err := bot.Login(cfg)
	if err != nil {
		panic("logging in to twitter: " + err.Error())
	}
	log.Printf("[Twitter] Logged in @%s\n", user.ScreenName)

	// run youtube scraper this in the background
	go checkYouTubeLive(client, user)

}

// checkYouTubeLive checks SpaceX's youtube live stream every 1-2 minutes and tweets if there is a starship launch
func checkYouTubeLive(client *twitter.Client, user *twitter.User) {
	defer panic("for some reason, the youtube live checker stopped running even though it never should")

	const spaceXLiveURL = "https://www.youtube.com/spacex/live"

	var (
		lastTweetedURL string
	)

	for {
		liveVideo, err := scrapers.YouTubeLive(spaceXLiveURL)
		if err == nil {
			if liveVideo.VideoID != "" && (match.Starship(liveVideo.Title) || match.Starship(liveVideo.ShortDescription)) {
				// Get the video URL
				liveURL := liveVideo.URL()

				if liveURL != lastTweetedURL {

					// OK, we can tweet this

					t, _, err := client.Statuses.Update("SpaceX #Starship stream is now live\n"+liveURL, nil)
					if err == nil {
						log.Printf("Tweeted https://twitter.com/%s/status/%s\n", user.ScreenName, t.IDStr)

						// make sure we don't tweet this again
						lastTweetedURL = liveURL
					} else {
						log.Println("Error while tweeting livestream update:", err.Error())
					}
				}
			}
		} else {
			if !errors.Is(err, scrapers.ErrNotLive) {
				log.Println("Unexpected error while scraping YouTube live:", err.Error())
			}
		}

		// Wait up to two minutes, then check again
		time.Sleep(time.Minute + time.Duration(rand.Intn(60))*time.Second)
	}
}
