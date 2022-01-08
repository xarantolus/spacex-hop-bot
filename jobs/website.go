package jobs

import (
	"fmt"
	"log"
	"math/rand"
	"reflect"
	"time"

	"github.com/xarantolus/spacex-hop-bot/consumer"
	"github.com/xarantolus/spacex-hop-bot/scrapers"
	"github.com/xarantolus/spacex-hop-bot/util"
)

const (
	changesFile = "website.json"
)

// StarshipWebsiteChanges watches the SpaceX starship page and tweets when the date or starship serial number change
func StarshipWebsiteChanges(client consumer.TwitterClient, linkChan chan<- string) {
	defer panic("website watcher stopped even though it never should")

	log.Println("[SpaceX] Watching Starship page for updates")

	var lastChange scrapers.StarshipInfo

	// Load our last state
	err := util.LoadJSON(changesFile, &lastChange)
	util.LogError(err, "loading changes file %q", changesFile)

	if err == nil {
		log.Printf("[Website] Waiting for new info, last was %s (NET %s)\n", lastChange.ShipName, lastChange.NextFlightDate.Format("2006-01-02"))
	} else {
		log.Println("[Website] Waiting for new info")
	}

	for {
		info, err := runWebsiteScrape(client, linkChan, scrapers.StarshipURL, lastChange, time.Now())
		if err != nil {
			util.LogError(err, "scraping SpaceX Starship website")
			goto sleep
		}

		if !reflect.DeepEqual(lastChange, info) {
			util.LogError(util.SaveJSON(changesFile, info), "saving changes file")

			lastChange = info
		}

	sleep:
		// Wait 2-4 minutes until checking again
		time.Sleep(2*time.Minute + time.Duration(rand.Intn(120))*time.Second)
	}
}

func runWebsiteScrape(client consumer.TwitterClient, linkChan chan<- string,
	starshipPageURL string, lastChange scrapers.StarshipInfo, datenow time.Time) (info scrapers.StarshipInfo, err error) {

	info, err = scrapers.SpaceXStarship(starshipPageURL, datenow)

	if info.LiveStreamID != "" && linkChan != nil {
		linkChan <- fmt.Sprintf("https://www.youtube.com/watch?v=%s", info.LiveStreamID)
	}

	if err != nil {
		return
	}

	// If it's the same info again, we don't care
	if lastChange.NextFlightDate.YearDay() >= info.NextFlightDate.YearDay() &&
		lastChange.NextFlightDate.Year() == info.NextFlightDate.Year() &&
		lastChange.Orbital == info.Orbital {
		return
	}

	// OK, now we have an interesting and new change
	var tweetText string
	if info.Orbital {
		tweetText = fmt.Sprintf("The SpaceX #Starship website now mentions %s for an orbital flight of #%s\n#WenHop\n%s",
			info.NextFlightDate.Format("January 2"), info.ShipName, scrapers.StarshipURL)
	} else {
		tweetText = fmt.Sprintf("The SpaceX #Starship website now mentions %s for #%s\n#WenHop\n%s",
			info.NextFlightDate.Format("January 2"), info.ShipName, scrapers.StarshipURL)
	}

	t, err := client.Tweet(tweetText, nil)
	if err != nil {
		err = fmt.Errorf("tweeting about update: %w", err)
		return
	}
	log.Println("[Twitter] Tweeted", util.TweetURL(t))

	return info, nil
}
