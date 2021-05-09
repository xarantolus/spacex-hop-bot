package jobs

import (
	"math/rand"
	"sort"
	"time"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/xarantolus/spacex-hop-bot/match"
	"github.com/xarantolus/spacex-hop-bot/util"
)

// CheckListTimeline requests the given lists about every minute or so. Any new tweets are put in tweetChan.
func CheckListTimeline(client *twitter.Client, list twitter.List, tweetChan chan<- match.TweetWrapper) {
	defer panic("list (" + list.Name + ") follower stopped processing even though it shouldn't")

	var (
		// lastSeenID is the ID of the last tweet we saw
		lastSeenID int64

		// The first batch of tweets we receive should not acted upon
		isFirstRequest = true
	)
	for {
		// https://developer.twitter.com/en/docs/twitter-api/v1/accounts-and-users/create-manage-lists/api-reference/get-lists-statuses
		tweets, _, err := client.Lists.Statuses(&twitter.ListsStatusesParams{
			ListID: list.ID,

			IncludeRetweets: twitter.Bool(true),
			IncludeEntities: twitter.Bool(true),
			SinceID:         lastSeenID, // everything since our last request
			Count:           200,        // Maximum number of tweets we can get at once
		})

		if err != nil {
			util.LogError(err, "list "+list.FullName)
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

			// We only look at tweets that appeared after the bot started
			if isFirstRequest {
				continue
			}

			// OK, process this tweet
			tweetChan <- match.TweetWrapper{
				TweetSource: match.TweetSourceKnownList,
				Tweet:       tweet,
			}
		}

		if isFirstRequest {
			isFirstRequest = false
		}

	sleep:
		// Add a random delay
		time.Sleep(time.Minute + time.Duration(rand.Intn(45))*time.Second)
	}
}
