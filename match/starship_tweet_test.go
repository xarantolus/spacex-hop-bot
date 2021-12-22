package match

import (
	"testing"

	"github.com/dghubble/go-twitter/twitter"
)

func tweet(acc, text string) TweetWrapper {
	return TweetWrapper{
		Tweet: twitter.Tweet{
			User: &twitter.User{
				ScreenName: acc,
			},
			FullText: text,
		},
	}
}

func TestStarshipTweet(t *testing.T) {
	tests := []struct {
		tweet TweetWrapper
		want  bool
	}{
		{
			tweet: tweet(
				"NASA_Marshall",
				"Starship launch hardware stands tall at @SpaceX while NASA HLS experts, @AstroKomrade, and @AstroVicGlover take a firsthand look. A Starship will land @NASAArtemis astronauts on the Moon during #Artemis III after @NASA_SLS and @NASA_Orion deliver the crew to lunar orbit.",
			),
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(t.Name(), func(t *testing.T) {
			if got := StarshipTweet(tt.tweet); got != tt.want {
				t.Errorf("StarshipTweet() = %v, want %v", got, tt.want)
			}
		})
	}
}
