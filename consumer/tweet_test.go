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

	location string

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
		},
	)
}
