package consumer

import (
	"fmt"
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

	parent *ttest
	quoted *ttest

	id int64
}

type TestTwitterClient struct {
	retweetedTweetIDs map[int64]bool

	tweets map[int64]*twitter.Tweet
}

func (r *TestTwitterClient) LoadStatus(tweetID int64) (*twitter.Tweet, error) {
	t, ok := r.tweets[tweetID]
	if ok {
		t.Retweeted = r.HasRetweeted(t.ID)
		return t, nil
	}

	return nil, fmt.Errorf("could not load status with id %d", tweetID)
}

func (n *TestTwitterClient) AddListMember(listID int64, userID int64) (err error) {
	return nil
}

func (r *TestTwitterClient) Retweet(tweet *twitter.Tweet) error {
	r.retweetedTweetIDs[tweet.ID] = true
	return nil
}

func (r *TestTwitterClient) HasRetweeted(tweetID int64) bool {
	return r.retweetedTweetIDs[tweetID]
}

func testStarshipRetweets(t *testing.T, tweets []ttest) {

	var processor = func() (p *Processor, t *TestTwitterClient) {
		t = &TestTwitterClient{
			retweetedTweetIDs: make(map[int64]bool),
			tweets:            make(map[int64]*twitter.Tweet),
		}

		p = NewProcessor(false, true, t, &twitter.User{ID: 5}, match.NewStarshipMatcherForTests(), 0)
		return
	}

	var tweetId int64 = 50
	var userId int64 = 80
	var tweet func(t *ttest) match.TweetWrapper
	tweet = func(t *ttest) match.TweetWrapper {
		if t.userID == 0 {
			t.userID = userId
			userId++
		}
		t.id = tweetId
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

		if t.parent != nil {
			tw.InReplyToStatusID = t.parent.id
		}

		if t.quoted != nil {
			quotedTweet := tweet(t.quoted)
			t.quoted.id = quotedTweet.ID

			tw.QuotedStatus = &quotedTweet.Tweet
			tw.QuotedStatusID = t.quoted.id
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

			// Populate & already show parent tweets to matcher.
			// That way it already knows/retweets tweets before the one we have here, making
			// it perfect for threads
			var matchParents func(t *ttest)
			matchParents = func(t *ttest) {
				if t == nil {
					return
				}

				matchParents(t.parent)

				var prevTweet = tweet(t)
				t.id = prevTweet.ID
				ret.tweets[prevTweet.ID] = &prevTweet.Tweet

				proc.Tweet(prevTweet)
			}
			matchParents(tt.parent)

			// Now we can generate & test the tweet we are actually interested in
			tweet := tweet(&tt)
			proc.Tweet(tweet)

			if !tt.want && ret.HasRetweeted(tweet.Tweet.ID) {
				t.Errorf("Tweet %q by %q was retweeted, but shouldn't have been", tt.text, tt.acc)
			}
			if tt.want && !ret.HasRetweeted(tweet.Tweet.ID) {
				t.Errorf("Tweet %q by %q was NOT retweeted, but should have been", tt.text, tt.acc)
			}

			if tt.quoted != nil {
				if !tt.quoted.want && ret.HasRetweeted(tt.quoted.id) {
					t.Errorf("Quoted tweet %q by %q was retweeted, but shouldn't have been", tt.quoted.text, tt.quoted.acc)
				}
				if tt.quoted.want && !ret.HasRetweeted(tt.quoted.id) {
					t.Errorf("Quoted tweet %q by %q was NOT retweeted, but should have been", tt.quoted.text, tt.quoted.acc)
				}
			}

			parent := tt.parent
			for parent != nil {
				if !parent.want && ret.HasRetweeted(parent.id) {
					t.Errorf("Parent tweet %q by %q was retweeted, but shouldn't have been", parent.text, parent.acc)
				}
				if parent.want && !ret.HasRetweeted(parent.id) {
					t.Errorf("Parent tweet %q by %q was NOT retweeted, but should have been", parent.text, parent.acc)
				}

				if parent.quoted != nil {
					if !parent.quoted.want && ret.HasRetweeted(parent.quoted.id) {
						t.Errorf("Parent quoted tweet %q by %q was retweeted, but shouldn't have been", parent.quoted.text, parent.quoted.acc)
					}
					if parent.quoted.want && !ret.HasRetweeted(parent.quoted.id) {
						t.Errorf("Parent quoted %q by %q was NOT retweeted, but should have been", parent.quoted.text, parent.quoted.acc)
					}
				}

				parent = parent.parent
			}
		})
	}
}
