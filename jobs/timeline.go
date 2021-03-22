package jobs

import (
	"math/rand"
	"sort"
	"time"

	"github.com/dghubble/go-twitter/twitter"
)

// CheckHomeTimeline requests the user home timeline about every minute and puts all new tweets in tweetChan.
// it also includes replies which would normally not be shown in the timeline.
// TL;DR: it stalks all users the account follows, even their replies
func CheckHomeTimeline(client *twitter.Client, tweetChan chan<- twitter.Tweet) {
	defer panic("home timeline follower stopped processing even though it shouldn't")

	var (
		// lastSeenID is the ID of the last tweet we saw
		lastSeenID int64

		// The first batch of tweets we receive should not acted upon
		isFirstRequest bool = true
	)

	for {
		// https://developer.twitter.com/en/docs/twitter-api/v1/tweets/timelines/api-reference/get-statuses-home_timeline
		tweets, _, err := client.Timelines.HomeTimeline(&twitter.HomeTimelineParams{
			ExcludeReplies:  twitter.Bool(false), // We want to get everything, including replies to tweets
			TrimUser:        twitter.Bool(false), // We care about the user
			IncludeEntities: twitter.Bool(false), // We also don't care about who was mentioned etc.
			SinceID:         lastSeenID,          // everything since our last request
			Count:           200,                 // Maximum number of tweets we can get at once
			TweetMode:       "extended",          // We have to use tweet.FullText instead of .Text
		})

		if err != nil {
			logError(err, "home timeline")
		} else {

			// Sort tweets so the first tweet we process is the oldest one
			sort.Slice(tweets, func(i, j int) bool {
				di, _ := tweets[i].CreatedAtTime()
				dj, _ := tweets[j].CreatedAtTime()

				return dj.After(di)
			})

			for _, tweet := range tweets {
				lastSeenID = tweet.ID

				// We only look at tweets that appeared after the bot started
				if isFirstRequest {
					continue
				}

				// OK, process this tweet
				tweetChan <- tweet
			}

			if isFirstRequest {
				isFirstRequest = false
			}
		}

		// I guess one request every minute is ok
		time.Sleep(time.Minute + time.Duration(rand.Intn(45))*time.Second)
	}
}
