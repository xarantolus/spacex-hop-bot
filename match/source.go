package match

import "github.com/dghubble/go-twitter/twitter"

// go get -u golang.org/x/tools/cmd/stringer
//go:generate stringer -type=TweetSource
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
