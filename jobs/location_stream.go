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
			Locations: []string{
				// This is a large area around boca chica (aka Starbase).
				// We want to catch many tweets from there and then filter them
				// You can see this area on a map here: https://mapper.acme.com/?ll=26.46074,-97.21252&z=9&t=M&marker0=26.90982%2C-96.59729%2Cunnamed&marker1=25.68237%2C-97.80029%2C1.7%20km%20NE%20of%20Valle%20Hermoso%20MX
				"-97.80029,25.68237,-96.59729,26.90982",

				// McGregor test site (they also test raptor engines there, so maybe someone tweets from there)
				// Map: https://mapper.acme.com/?ll=31.39966,-97.46246&z=12&t=M&marker0=31.39930%2C-97.46250%2C31.399308%20-97.462496&marker1=31.34836%2C-97.51740%2Cunnamed&marker2=31.48314%2C-97.36530%2C6.0%20km%20NE%20of%20McGregor%20TX
				"-97.51740,31.34836,-97.36530,31.48314",

				// Port/Cape Canaveral (the oil rigs that could be used for starship landings are stationed there)
				// Also includes SpaceX's LC-39A, where a new starship orbital launch pad is being constructed
				// Map: https://mapper.acme.com/?ll=28.40952,-80.60944&z=10&t=M&marker0=28.21910%2C-80.79552%2Cunnamed&marker1=28.88617%2C-79.96262%2C79.2%20km%20ExNE%20of%20Merritt%20Island%20FL
				"-80.79552,28.21910,-79.96262,28.88617",
			},
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
