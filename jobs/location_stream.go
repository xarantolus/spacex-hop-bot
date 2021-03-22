package jobs

import "github.com/dghubble/go-twitter/twitter"

// CheckLocationStream checks out tweets from a large area around boca chica
func CheckLocationStream(client *twitter.Client, tweetChan chan<- twitter.Tweet) {
	defer panic("location stream ended even though it never should")

	s, err := client.Streams.Filter(&twitter.StreamFilterParams{
		// This is a large area around boca chica. We want to catch many tweets from there and then filter them
		// You can see this area on a map here: https://mapper.acme.com/?ll=26.00002,-97.07932&z=10&t=M&marker0=25.98750%2C-97.18639%2CSpaceX%20South%20Texas%20launch%20site&marker1=26.39190%2C-96.71811%2C26.3919%20-96.7181&marker2=25.52629%2C-97.43501%2C25.5263%20-97.4350
		Locations:   []string{"-97.4350,25.5263,-96.7181,26.3919"},
		FilterLevel: "none",
		Language:    []string{"en"},
	})
	if err != nil {
		panic("setting up location stream: " + err.Error())
	}

	// Stream all tweets and serve them to the channel
	for m := range s.Messages {
		t, ok := m.(*twitter.Tweet)
		if !ok || t == nil {
			continue
		}

		tweetChan <- *t
	}
}
