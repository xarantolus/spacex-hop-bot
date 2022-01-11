package match

import (
	"fmt"
	"log"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/xarantolus/spacex-hop-bot/util"
)

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

	EnableLogging bool
}

func (t *TweetWrapper) Log(format string, a ...interface{}) {
	if t.EnableLogging {
		log.Printf("[Processor] %s (%s): %s", util.TweetURL(&t.Tweet), t.TweetSource.String(), fmt.Sprintf(format, a...))
	}
}
