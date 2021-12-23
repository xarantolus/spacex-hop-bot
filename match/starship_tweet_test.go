package match

import (
	"testing"
	"time"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/xarantolus/spacex-hop-bot/bot"
)

func newEmptyStarshipMatcher() *StarshipMatcher {
	return &StarshipMatcher{
		Ignorer: &Ignorer{
			list:     bot.ListMembers(nil, "test"),
			keywords: nil,
		},
	}
}

type ttest struct {
	acc  string
	text string

	location string

	hasMedia bool

	want bool
}

func testStarshipTweets(t *testing.T, tweets []ttest) {
	var matcher = newEmptyStarshipMatcher()

	var tweet = func(t ttest) TweetWrapper {
		var tw = TweetWrapper{
			Tweet: twitter.Tweet{
				User: &twitter.User{
					ScreenName: t.acc,
				},
				FullText: t.text,
			},
		}

		// Set a recent date, aka now (the bot usually sees very recent tweets)
		tw.CreatedAt = time.Now().Add(-time.Minute).Format(time.RubyDate)

		if tw.User.ScreenName == "" {
			tw.User = &twitter.User{
				ScreenName: "default_name",
			}
		}

		if t.location != "" {
			tw.Place = &twitter.Place{
				ID: t.location,
			}
		}

		// Just add a dummy photo
		if t.hasMedia {
			tw.Entities = &twitter.Entities{
				Media: []twitter.MediaEntity{
					{
						ID: 1024,
					},
				},
			}
		}

		return tw
	}

	for _, tt := range tweets {
		t.Run(t.Name(), func(t *testing.T) {
			if got := matcher.StarshipTweet(tweet(tt)); got != tt.want {
				t.Errorf("StarshipTweet(%q) = %v, want %v", tt.text, got, tt.want)
			}
		})
	}
}

func TestStarshipTweetNASA(t *testing.T) {
	testStarshipTweets(t,
		[]ttest{
			{
				acc:  "NASA_Marshall",
				text: "Starship launch hardware stands tall at @SpaceX while NASA HLS experts, @AstroKomrade, and @AstroVicGlover take a firsthand look. A Starship will land @NASAArtemis astronauts on the Moon during #Artemis III after @NASA_SLS and @NASA_Orion deliver the crew to lunar orbit.",
				want: true,
			},
			{
				acc:  "NASA",
				text: "Starship will land humans on the moon.",
				want: true,
			},
			{
				acc:  "NASA",
				text: "Unrelated orion tweet.",
				want: false,
			},
		},
	)
}

func TestStarshipTweetOld(t *testing.T) {
	testStarshipTweets(t,
		[]ttest{
			{
				// Photo from today should be retweeted
				text: "Photo of Starship taken on " + time.Now().Format("02. January 2006"),
				want: true,
			},
			{
				// Older photos should not be retweeted
				text: "Photo of Starship taken on 24. October 2021",
				want: false,
			},
		},
	)
}

func TestStarshipTweetAntiOverwrite(t *testing.T) {
	testStarshipTweets(t,
		[]ttest{
			{
				text: "Normal user tweeting about Starship and Orion",
				want: false,
			},
			{
				acc:  "NASA",
				text: "NASA account tweeting about Starship and Orion",
				want: true,
			},
		},
	)
}

func TestStarshipTweetPlace(t *testing.T) {
	testStarshipTweets(t,
		[]ttest{
			{
				text:     "What a nice ship",
				location: SpaceXBuildSiteID,
				want:     true,
			},
			{
				text:     "This is worse than SLS",
				location: SpaceXBuildSiteID,
				want:     false,
			},
			{
				hasMedia: true,
				location: SpaceXBuildSiteID,
				want:     true,
			},
			{
				text: "Raptors roaring!",
				want: false,
			},
			{
				text:     "Raptors roaring!",
				location: SpaceXMcGregorPlaceID,
				want:     true,
			},
		},
	)
}

func TestStarshipTweetSpam(t *testing.T) {
	testStarshipTweets(t,
		[]ttest{
			{
				text: "@a @b @a @a @b @a @a @b @a @a @b @a @a @b @a are just as annoying as I am",
				want: false,
			},
			{
				text:     "Doing stuff b4 work",
				hasMedia: true,
				want:     false,
			},
		},
	)
}

func TestStarshipTweetSpecificMatchers(t *testing.T) {
	testStarshipTweets(t,
		[]ttest{
			{
				acc:      "Bocachicagal",
				text:     "I have received an alert notice",
				hasMedia: true,
				want:     true,
			},
			{
				acc:  "spacexboca",
				text: "Closure for testing has begun",
				want: true,
			},
		},
	)
}

func TestStarshipTweetSpecificHQMedia(t *testing.T) {
	testStarshipTweets(t,
		[]ttest{
			{
				acc:      "starshipgazer",
				hasMedia: true,
				want:     true,
			},
			{
				hasMedia: true,
				want:     false,
			},
		},
	)
}
