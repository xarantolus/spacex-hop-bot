package jobs

import (
	"errors"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/xarantolus/spacex-hop-bot/scrapers"
	"github.com/xarantolus/spacex-hop-bot/util"
)

const (
	changesFile = "website.json"
)

// StarshipWebsiteChanges watches the SpaceX starship page and tweets when the date or starship serial number change
func StarshipWebsiteChanges(client *twitter.Client, linkChan chan<- string) {
	defer panic("website watcher stopped even though it never should")

	log.Println("[SpaceX] Watching Starship page for updates")

	var lastChange scrapers.StarshipInfo

	// Load our last state
	err := util.LoadJSON(changesFile, &lastChange)
	util.LogError(err, "loading changes file "+changesFile)

	if err == nil {
		log.Printf("[Website] Waiting for new info, last was %s (NET %s)\n", lastChange.ShipName, lastChange.NextFlightDate.Format("2006-01-02"))
	} else {
		log.Println("[Website] Waiting for new info")
	}

	for {
		info, err := scrapers.SpaceXStarship()

		if info.LiveStreamID != "" {
			linkChan <- fmt.Sprintf("https://www.youtube.com/watch?v=%s", info.LiveStreamID)
		}

		if err != nil {
			// Log only interesting errors
			if !errors.Is(err, scrapers.ErrNoInfo) {
				util.LogError(err, "scraping SpaceX Starship website")
			}

			goto sleep
		}

		// If it's the same info again, we don't care
		if lastChange.NextFlightDate.YearDay() >= info.NextFlightDate.YearDay() {
			goto sleep
		}

		// OK, so now we got an interesting change
		lastChange = info

		// Save this one
		util.LogError(util.SaveJSON(changesFile, lastChange), "saving changes file")

		// OK, now we have an interesting and new change
		{
			var tweetText = fmt.Sprintf("The SpaceX #Starship website now mentions %s for #%s #WenHop\n%s",
				lastChange.NextFlightDate.Format("January 02"), lastChange.ShipName, scrapers.StarshipURL)

			t, _, err := client.Statuses.Update(tweetText, nil)
			if err != nil {
				util.LogError(err, "tweeting starship update")
			} else {
				log.Println("[Twitter] Tweeted", util.TweetURL(t))
			}
		}
	sleep:
		// Wait 5-10 minutes until checking again
		time.Sleep(5*time.Minute + time.Duration(rand.Intn(300))*time.Second)
	}
}
