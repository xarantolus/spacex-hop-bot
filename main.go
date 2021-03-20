package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"regexp"
	"sort"
	"time"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/xarantolus/spacex-hop-bot/bot"
	"github.com/xarantolus/spacex-hop-bot/config"
	"github.com/xarantolus/spacex-hop-bot/match"
	"github.com/xarantolus/spacex-hop-bot/scrapers"
)

var flagConfigFile = flag.String("cfg", "deploy/config.yaml", "Config file path")

func main() {
	flag.Parse()

	// Some stuff depends on randomness
	rand.Seed(time.Now().UnixNano())

	cfg, err := config.Parse(*flagConfigFile)
	if err != nil {
		panic("parsing configuration file: " + err.Error())
	}

	client, user, err := bot.Login(cfg)
	if err != nil {
		panic("logging in to twitter: " + err.Error())
	}
	log.Printf("[Twitter] Logged in @%s\n", user.ScreenName)

	// run youtube scraper this in the background
	go checkYouTubeLive(client, user)

	checkHomeTimeline(client, user)
}

// checkYouTubeLive checks SpaceX's youtube live stream every 1-2 minutes and tweets if there is a starship launch
func checkYouTubeLive(client *twitter.Client, user *twitter.User) {
	defer panic("for some reason, the youtube live checker stopped running even though it never should")

	const spaceXLiveURL = "https://www.youtube.com/spacex/live"
	var shipNameRegex = regexp.MustCompile(`(SN\d+)`)

	var (
		lastTweetedURL string
	)

	for {
		liveVideo, err := scrapers.YouTubeLive(spaceXLiveURL)
		if err == nil {
			if liveVideo.VideoID != "" && (match.StarshipText(liveVideo.Title) || match.StarshipText(liveVideo.ShortDescription)) {
				// Get the video URL
				liveURL := liveVideo.URL()

				if liveURL != lastTweetedURL {

					// See if we can get the starship name, but we tweet without it anyway
					var shipName = shipNameRegex.FindString(liveVideo.Title)
					if shipName != "" {
						shipName = " #" + shipName
					}

					tweetText := fmt.Sprintf("It's hoppening! SpaceX #Starship%s stream is live\n%s", shipName, liveURL)

					// OK, we can tweet this

					t, _, err := client.Statuses.Update(tweetText, nil)
					if err == nil {
						log.Printf("Tweeted https://twitter.com/%s/status/%s\n", user.ScreenName, t.IDStr)

						// make sure we don't tweet this again
						lastTweetedURL = liveURL
					} else {
						log.Println("Error while tweeting livestream update:", err.Error())
					}
				}
			}
		} else {
			if !errors.Is(err, scrapers.ErrNotLive) {
				log.Println("Unexpected error while scraping YouTube live:", err.Error())
			}
		}

		// Wait up to two minutes, then check again
		time.Sleep(time.Minute + time.Duration(rand.Intn(60))*time.Second)
	}
}

func checkHomeTimeline(client *twitter.Client, user *twitter.User) {
	// defer panic("user stream follower stopped processing even though it shouldn't")

	var (
		seenTweets = make(map[int64]bool)
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
			log.Println("Home timeline:", err.Error())
		} else {

			// Sort tweets so the first tweet we process is the oldest one
			sort.Slice(tweets, func(i, j int) bool {
				di, _ := tweets[i].CreatedAtTime()
				dj, _ := tweets[j].CreatedAtTime()

				return dj.After(di)
			})

			for _, tweet := range tweets {
				lastSeenID = tweet.ID
				if seenTweets[tweet.ID] {
					continue
				}
				seenTweets[tweet.ID] = true

				// We only look at tweets that appeared after the bot started
				if isFirstRequest {
					continue
				}

				switch {
				case tweet.RetweetedStatus != nil && match.StarshipTweet(tweet.RetweetedStatus):
					retweet(client, tweet.RetweetedStatus)
				case match.StarshipTweet(&tweet):
					retweet(client, &tweet)
				case tweet.InReplyToStatusID != 0:
					t, _, _ := client.Statuses.Lookup([]int64{tweet.InReplyToStatusID}, &twitter.StatusLookupParams{
						IncludeEntities: twitter.Bool(false),
						TweetMode:       "extended",
					})

					if len(t) > 0 {
						processThread(client, &t[0])
					}
				}
			}

			if isFirstRequest {
				isFirstRequest = false
			}
		}

		// I guess one request every minute is ok
		time.Sleep(time.Minute)
	}
}

func retweet(client *twitter.Client, tweet *twitter.Tweet) {
	if tweet.Retweeted {
		return
	}

	twurl := tweetURL(tweet)
	fmt.Print("Would retweet ", twurl)
	fmt.Println()

	tweet.Retweeted = true
}

func tweetURL(tweet *twitter.Tweet) string {
	if tweet.User == nil {
		return "https://twitter.com/i/status/" + tweet.IDStr
	}
	return "https://twitter.com/" + tweet.User.ScreenName + "/status/" + tweet.IDStr
}

// processThread process tweet threads and should retweet everything in them.
// This is useful because Elon Musk often replies to people that quote tweeted/asked a questions on his tweets
// See this one for example: https://twitter.com/elonmusk/status/1372826575293583366
// OR: https://twitter.com/elonmusk/status/1372725108909957121
func processThread(client *twitter.Client, tweet *twitter.Tweet) (didRetweet bool) {
	if tweet == nil {
		// Just in case
		return false
	}

	// First process the rest of the thread
	if tweet.InReplyToStatusID != 0 {
		// Ok, there was a reply. Check if we can do something with that
		t, _, _ := client.Statuses.Lookup([]int64{tweet.InReplyToStatusID}, &twitter.StatusLookupParams{
			IncludeEntities: twitter.Bool(false),
			TweetMode:       "extended",
		})

		if len(t) > 0 {
			if !processThread(client, &t[0]) {
				return false
			}

			retweet(client, &t[0])
			didRetweet = true
		}
	}

	// A quoted tweet. Let's see if there's anything interesting
	if tweet.QuotedStatusID != 0 && tweet.QuotedStatus != nil {
		return processThread(client, tweet.QuotedStatus)
	}

	// Now actually match the tweet
	if didRetweet || match.StarshipTweet(tweet) {

		retweet(client, tweet)

		return true
	}

	return didRetweet
}
