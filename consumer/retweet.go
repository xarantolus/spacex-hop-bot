package consumer

import (
	"log"

	"github.com/dghubble/go-twitter/twitter"
)

type Retweeter interface {
	Retweet(*twitter.Tweet) error
}

type NormalRetweeter struct {
	Client *twitter.Client
	Debug  bool
}

func (r *NormalRetweeter) Retweet(tweet *twitter.Tweet) error {
	if r.Debug {
		return nil
	}

	_, _, err := r.Client.Statuses.Retweet(tweet.ID, nil)

	return err
}

type TestRetweeter struct {
	tweetIDs map[int64]bool
}

func (r *TestRetweeter) Retweet(tweet *twitter.Tweet) error {
	if r.tweetIDs == nil {
		r.tweetIDs = make(map[int64]bool)
	}
	r.tweetIDs[tweet.ID] = true
	return nil
}

func (r *TestRetweeter) Contains(tweet *twitter.Tweet) bool {
	if r.tweetIDs == nil {
		return false
	}
	return r.tweetIDs[tweet.ID]
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
