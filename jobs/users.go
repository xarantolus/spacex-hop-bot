package jobs

import (
	"math/rand"
	"sort"
	"time"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/xarantolus/spacex-hop-bot/util"
)

// CheckUserTimeline requests the given user profile every few minutes or so
func CheckUserTimeline(client *twitter.Client, name string, tweetChan chan<- twitter.Tweet) {
	defer panic("user (" + name + ") follower stopped processing even though it shouldn't")

	var (
		// lastSeenID is the ID of the last tweet we saw
		lastSeenID int64
	)

	for {
		tweets, _, err := client.Timelines.UserTimeline(&twitter.UserTimelineParams{
			ScreenName:     name,
			TweetMode:      "extended",
			ExcludeReplies: twitter.Bool(false),
			SinceID:        lastSeenID,
		})

		if err != nil {
			util.LogError(err, "user "+name)
			goto sleep
		}

		// Sort tweets so the first tweet we process is the oldest one
		sort.Slice(tweets, func(i, j int) bool {
			di, _ := tweets[i].CreatedAtTime()
			dj, _ := tweets[j].CreatedAtTime()

			return dj.After(di)
		})

		for _, tweet := range tweets {
			lastSeenID = tweet.ID

			// OK, process this tweet
			tweetChan <- tweet
		}

	sleep:
		// Add a random delay
		time.Sleep(2*time.Minute + time.Duration(rand.Intn(500))*time.Second)
	}
}
