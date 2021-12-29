package consumer

import (
	"fmt"
	"log"
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

type TestTwitterClient struct {
	retweetedTweetIDs map[int64]bool

	tweets map[int64]*twitter.Tweet
}

func (r *TestTwitterClient) LoadStatus(tweetID int64) (*twitter.Tweet, error) {
	t, ok := r.tweets[tweetID]
	if ok {
		return t, nil
	}

	return nil, fmt.Errorf("could not load status with id %d", tweetID)
}

func (n *TestTwitterClient) AddListMember(listID int64, userID int64) (err error) {
	return nil
}

func (r *TestTwitterClient) Retweet(tweet *twitter.Tweet) error {
	if r.retweetedTweetIDs == nil {
		r.retweetedTweetIDs = make(map[int64]bool)
	}
	r.retweetedTweetIDs[tweet.ID] = true
	return nil
}

func (r *TestTwitterClient) HasRetweeted(tweet *twitter.Tweet) bool {
	if r.retweetedTweetIDs == nil {
		return false
	}
	return r.retweetedTweetIDs[tweet.ID]
}

type DebugRetweeter struct {
}

func (r *DebugRetweeter) Retweet(tweet *twitter.Tweet) error {
	if tweet.User == nil {
		tweet.User = &twitter.User{}
	}
	log.Printf("Would have retweeted %q by %q, but we are in debug mode", tweet.Text(), tweet.User.Name)

	return nil
}

func testStarshipRetweets(t *testing.T, tweets []ttest) {

	var processor = func() (p *Processor, t *TestTwitterClient) {
		t = &TestTwitterClient{}

		p = NewProcessor(false, true, t, &twitter.User{ID: 5}, match.NewStarshipMatcherForTests(), 0)
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

			if !tt.want && ret.HasRetweeted(&tweet.Tweet) {
				t.Errorf("Tweet %q by %q was retweeted, but shouldn't have been", tt.text, tt.acc)
			}
			if tt.want && !ret.HasRetweeted(&tweet.Tweet) {
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

func TestQuestionTweets(t *testing.T) {
	testStarshipRetweets(t,
		[]ttest{
			// Questions only if we have media or are at the spacex locations
			{
				acc:  "considercosmos",
				text: "Super Heavy is now hooked up to @SpaceX crane...\nWill we see a booster 4 lift soon?",
				want: false,
			},
			{
				acc:      "considercosmos",
				text:     "Super Heavy is now hooked up to @SpaceX crane...\nWill we see a booster 4 lift soon?",
				hasMedia: true,
				want:     true,
			},
			{
				text:     "Super Heavy is now hooked up to @SpaceX crane...\nWill we see a booster 4 lift soon?",
				location: match.SpaceXBuildSiteID,
				want:     true,
			},
		},
	)
}

func TestElonTweets(t *testing.T) {
	testStarshipRetweets(t,
		[]ttest{
			{
				acc:  "elonmusk",
				text: "Wow, working on this problem has soaked up a lot of my time & brain cycles over the past ~7 years! This and Starship engines are currently the two hardest problems.",
				want: true,
			},
		},
	)
}

// TODO: Add parent/child tweet, e.g. https://twitter.com/NASASpaceflight/status/1476249730585968647 should be retweeted, and also test that it is
