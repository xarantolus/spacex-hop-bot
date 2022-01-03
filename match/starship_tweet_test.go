package match

import (
	"testing"
	"time"

	"github.com/dghubble/go-twitter/twitter"
)

type ttest struct {
	acc             string
	userID          int64
	userDescription string

	date string

	text string

	location string

	hasMedia bool

	want bool
}

func testStarshipTweets(t *testing.T, tweets []ttest) {
	t.Helper()

	var matcher = NewStarshipMatcherForTests()

	var tweet = func(t ttest) TweetWrapper {
		var tw = TweetWrapper{
			Tweet: twitter.Tweet{
				CreatedAt: t.date,
				User: &twitter.User{
					ScreenName:  t.acc,
					ID:          t.userID,
					Description: t.userDescription,
				},
				FullText: t.text,
			},
		}

		// Set a recent date, aka now (the bot usually sees very recent tweets)
		if tw.CreatedAt == "" {
			tw.CreatedAt = time.Now().Add(-time.Minute).Format(time.RubyDate)
		}

		if tw.User.ScreenName == "" {
			tw.User = &twitter.User{
				ScreenName:  "default_name",
				Description: t.userDescription,
				ID:          t.userID,
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
				t.Errorf("StarshipTweet(%q by %q) = %v, want %v", tt.text, tt.acc, got, tt.want)
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
				text:     "Raptors roaring!",
				location: SpaceXMcGregorPlaceID,
				want:     true,
			},
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
				acc:  "AustinDeSisto",
				text: "Pad clear in 1 hour announcement at pad!",
				want: true,
			},
			{
				acc:  "starshipgazer",
				text: `Just announced "entire pad clear in 45 minutes" over the loud speakers at launch complex.`,
				want: true,
			},
			{
				text: "This pad will clear all your laundry!",
				want: false,
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

func TestOldTweets(t *testing.T) {
	testStarshipTweets(t,
		[]ttest{
			{
				date: "Wed Dec 11 17:52:17 +0000 2021",
				text: "Starship S20 will light its raptor engines soon",
				want: false,
			},
		},
	)
}

func TestStarshipTweetIgnoredAccount(t *testing.T) {
	testStarshipTweets(t,
		[]ttest{
			// Tweet of a render, but not marked as such. However the description contains that info
			{
				text:            "Starship 20 static fire",
				userDescription: "3D artist",
				want:            false,
			},
			{
				text: "Starship 20 static fire",
				want: true,
			},

			// Make sure we ignore ignored accounts
			{
				userID: testIgnoredUserID,
				acc:    "ignored_user",
				text:   "Starship 20 S.C.A.M (Starship Camera) now here",
				want:   false,
			},
			{
				// Same tweet, but by not ignored user
				text: "Starship 20 S.C.A.M (Starship Camera) now here",
				want: true,
			},
		},
	)
}
