package jobs

import (
	"log"
	"time"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/xarantolus/spacex-hop-bot/match"
	"github.com/xarantolus/spacex-hop-bot/util"
)

// CheckLocationStream checks out tweets from a large area around boca chica
func CheckLocationStream(client *twitter.Client, tweetChan chan<- match.TweetWrapper) {
	defer panic("location stream ended even though it never should")

	var backoff int = 1
	for {
		s, err := client.Streams.Filter(&twitter.StreamFilterParams{
			// This is a large area around boca chica. We want to catch many tweets from there and then filter them
			// You can see this area on a map here: https://mapper.acme.com/?ll=26.46074,-97.21252&z=9&t=M&marker0=26.90982%2C-96.59729%2Cunnamed&marker1=25.68237%2C-97.80029%2C1.7%20km%20NE%20of%20Valle%20Hermoso%20MX
			Locations:   []string{"-97.80029,25.68237,-96.59729,26.90982"},
			FilterLevel: "none",
			Language:    []string{"en"},
		})
		if err != nil {
			util.LogError(err, "location stream")
			goto sleep
		}

		log.Println("[Twitter] Connected to location stream")

		// Stream all tweets and serve them to the channel
		for m := range s.Messages {
			backoff = 1
			t, ok := m.(*twitter.Tweet)
			if !ok || t == nil {
				continue
			}

			// If we have truncated text, we try to get the whole tweet
			if t.Truncated {
				t, _, err = client.Statuses.Show(t.ID, &twitter.StatusShowParams{
					TweetMode: "extended",
				})
				if err != nil {
					continue
				}
			}

			tweetChan <- match.TweetWrapper{
				TweetSource: match.TweetSourceLocationStream,
				Tweet:       *t,
			}
		}

		backoff *= 2

		log.Printf("[Twitter] Location stream ended for some reason, trying again in %d seconds", backoff*5)
	sleep:
		time.Sleep(time.Duration(backoff) * 5 * time.Second)
	}
}
