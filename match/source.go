package match

import "github.com/dghubble/go-twitter/twitter"

type TweetSource int

const (
	TweetSourceUnknown TweetSource = iota
	TweetSourceLocationStream
	TweetSourceKnownList
	TweetSourceTimeline
	TweetSourceTrustedUser
)

type TweetWrapper struct {
	TweetSource TweetSource
	twitter.Tweet
}
