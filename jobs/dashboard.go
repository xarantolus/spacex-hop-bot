package jobs

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/xarantolus/spacex-hop-bot/jobs/review"
	"github.com/xarantolus/spacex-hop-bot/util"
)

func CheckDashboard(twitterClient *twitter.Client) {
	defer panic("for some reason, the gov dashboard checker stopped running even though it never should")

	var client = review.NewReviewClient()

	log.Println("[Review] Start watching Environmental Review dashboard")

	var previousTweetID int64

	for {
		diffs, err := client.ReportProjectDiff(review.StarshipBocaProjectID)
		if err != nil {
			util.LogError(err, "Review dashboard")
			goto sleep
		}

		// If there were any changes to the dashboard, we of course tweet about them.
		// We basically create a first tweet with the first diff; that tweet links to
		// the dashboard and tags @elonmusk
		// The next tweets will then be only the diff description as an answer to the first tweet
		for _, diffText := range diffs {
			var tweetText = generateReviewTweetText(diffText, previousTweetID == 0)

			tweet, _, err := twitterClient.Statuses.Update(tweetText, &twitter.StatusUpdateParams{
				InReplyToStatusID: previousTweetID,
			})
			if err != nil {
				log.Printf("[Review] Error while sending tweet with text %q: %s", tweetText, err.Error())
				continue
			}

			previousTweetID = tweet.ID

			log.Println("[Twitter] Tweeted", util.TweetURL(tweet))
		}

		previousTweetID = 0

	sleep:
		time.Sleep(time.Minute + time.Duration(rand.Intn(90))*time.Second)
	}
}

func generateReviewTweetText(description string, withTags bool) string {
	if withTags {
		return fmt.Sprintf("#Starship review update: %s\n\n@elonmusk\n%s", description, review.StarshipBocaDashboardURL)
	}

	return description
}
