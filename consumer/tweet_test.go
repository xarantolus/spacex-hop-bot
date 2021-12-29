package consumer

import (
	"testing"
	"time"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/xarantolus/spacex-hop-bot/match"
)

type ttest struct {
	acc  string
	text string

	userID int64

	tweetSource match.TweetSource
	location    string

	hasMedia bool

	want bool
}

func testStarshipRetweets(t *testing.T, tweets []ttest) {

	var processor = func() (p *Processor, t *TestRetweeter) {
		t = &TestRetweeter{}

		p = NewProcessor(false, true, nil, &twitter.User{ID: 5}, match.NewStarshipMatcherForTests(), t, 0)
		return
	}

	var tweetId int64
	var tweet = func(t ttest) match.TweetWrapper {
		var tw = match.TweetWrapper{
			Tweet: twitter.Tweet{
				User: &twitter.User{
					ScreenName: t.acc,
					ID:         t.userID,
				},
				FullText: t.text,
				ID:       tweetId,
			},

			TweetSource: t.tweetSource,
		}
		tweetId++

		// Set a recent date, aka now (the bot usually sees very recent tweets)
		tw.CreatedAt = time.Now().Add(-time.Minute).Format(time.RubyDate)

		if tw.User.ScreenName == "" {
			tw.User.ScreenName = "default_name"
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
			proc, ret := processor()

			tweet := tweet(tt)

			proc.Tweet(tweet)

			if !tt.want && ret.Contains(&tweet.Tweet) {
				t.Errorf("Tweet %q by %q was retweeted, but shouldn't have been", tt.text, tt.acc)
			}
			if tt.want && !ret.Contains(&tweet.Tweet) {
				t.Errorf("Tweet %q by %q was NOT retweeted, but should have been", tt.text, tt.acc)
			}
		})
	}
}

func TestBasicTweets(t *testing.T) {
	testStarshipRetweets(t,
		[]ttest{
			{
				// this is the test user ID; we don't want to retweet our own tweets
				userID: 5,
				text:   "S20 standing on the pad",
				want:   false,
			},
			{
				text: "S20 standing on the pad",
				want: true,
			},
			{
				text: "Unrelated",
				want: false,
			},
			{
				text: "Road closure with no information where it is",
				want: false,
			},

			{
				text: "Road closure with no information where it is, but trusted account",
				acc:  "nextspaceflight",
				want: true,
			},

			// If we have a tweet that only contains (hash)tags, it should only retweeted if it has media
			{
				text:     "#Starbase #Starbase #SpaceX #Starship @elonmusk",
				hasMedia: true,
				want:     true,
			},
			{
				text: "#Starbase #Starbase #SpaceX #Starship @elonmusk",
				want: false,
			},
		},
	)
}

func TestLocationTweets(t *testing.T) {
	testStarshipRetweets(t,
		[]ttest{
			// If it explicitly mentions a starship, then no need for location
			{
				text: "Pad announcement over the speakers: clearing pad for S20 static fire",
				want: true,
			},
			// Here we have the same tweet, but one with a good location
			{
				text: "Pad announcement over the speakers: clearing pad for static fire",
				want: false,
			},
			{
				text:     "Pad announcement over the speakers: clearing pad for static fire",
				location: "random place",
				want:     false,
			},
			{
				text:     "Pad announcement over the speakers: clearing pad for static fire",
				location: match.StarbasePlaceID,
				hasMedia: true,
				want:     true,
			},

			{
				text:        "Pad announcement over the speakers: clearing pad for static fire",
				want:        false,
				tweetSource: match.TweetSourceLocationStream,
			},
			{
				text:        "Pad announcement over the speakers: clearing pad for static fire",
				location:    "random place",
				tweetSource: match.TweetSourceLocationStream,
				want:        false,
			},
			{
				text:        "Pad announcement over the speakers: clearing pad for static fire",
				location:    match.StarbasePlaceID,
				tweetSource: match.TweetSourceLocationStream,
				hasMedia:    true,
				want:        true,
			},

			// However, we don't want *any* tweet from starbase etc.
			{
				text:        "Drinking some coffee at the beach",
				location:    match.StarbasePlaceID,
				tweetSource: match.TweetSourceLocationStream,
				want:        false,
			},
		},
	)
}
