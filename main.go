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

var flagConfigFile = flag.String("cfg", "config.yaml", "Config file path")

func main() {
	flag.Parse()

	// Some stuff depends on randomness
	rand.Seed(time.Now().UnixNano())

	cfg, err := config.Parse(*flagConfigFile)
	if err != nil {
		panic("parsing configuration file: " + err.Error())
	}

	client, selfUser, err := bot.Login(cfg)
	if err != nil {
		panic("logging in to twitter: " + err.Error())
	}
	log.Printf("[Twitter] Logged in @%s\n", selfUser.ScreenName)

	// contains all tweets the bot should check
	var tweetChan = make(chan twitter.Tweet, 250)

	// Run YouTube scraper in the background,
	// it will tweet if it discovers that SpaceX is online with a Starship stream
	go checkYouTubeLive(client, selfUser)

	lists, _, err := client.Lists.List(&twitter.ListsListParams{})
	if len(lists) == 100 {
		// See https://developer.twitter.com/en/docs/twitter-api/v1/accounts-and-users/create-manage-lists/api-reference/get-lists-list
		log.Println("[Warning] Lists API call returned 100 lists, which means that it is likely that some lists were not included. See API URL in comment above this line")
	}
	if err != nil {
		panic("initializing bot: couldn't retrieve lists: " + err.Error())
	}

	// Start watching all lists we follow
	for _, l := range lists {
		go checkListTimeline(client, l, tweetChan)
	}

	log.Printf("[Twitter] Started watching %d lists\n", len(lists))

	// Check out the home timeline of the bot user, it will contain all kinds of tweets from all kinds of people
	go checkHomeTimeline(client, tweetChan)

	var (
		// seenTweets maps the tweet id to an boolean. If the tweet was already processed/seen, it is put here
		seenTweets = make(map[int64]bool)
	)

	for tweet := range tweetChan {
		if seenTweets[tweet.ID] || tweet.Retweeted {
			continue
		}

		// Skip our own tweets
		if tweet.User != nil && tweet.User.ID == selfUser.ID {
			continue
		}

		// So now we got a tweet. There are three categories that interest us:
		// 1. Elon Musk drops insider info about starship, e.g. as a reply.
		//    We do not care about his other tweets, so we check if any tweet
		//    in the reply chain matches stuff about starship
		// 2. We find a retweet of a tweet that contains a certain keyword, e.g. Starship
		// 3. We find a tweet that is about starship

		switch {
		case tweet.User != nil && tweet.User.ScreenName == "elonmusk":
			// When elon drops starship info, we want to retweet it.
			// We basically detect if the thread/tweet is about starship and
			// retweet everything that is appropriate
			processThread(client, &tweet, seenTweets)
		case tweet.RetweetedStatus != nil && match.StarshipTweet(tweet.RetweetedStatus):
			// If it's a retweet of someone, we check that tweet if it's interesting
			retweet(client, tweet.RetweetedStatus)
		case match.StarshipTweet(&tweet):
			// If the tweet itself is about starship, we retweet it
			retweet(client, &tweet)
		}
	}
}

// checkYouTubeLive checks SpaceX's youtube live stream every 1-2 minutes and tweets if there is a starship launch
func checkYouTubeLive(client *twitter.Client, user *twitter.User) {
	defer panic("for some reason, the youtube live checker stopped running even though it never should")

	log.Println("[YouTube] Watching SpaceX channel for live Starship streams")

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
						log.Println("[Twitter] Tweeted", tweetURL(t))

						// make sure we don't tweet this again
						lastTweetedURL = liveURL
					} else {
						log.Println("[Twitter] Error while tweeting livestream update:", err.Error())
					}
				}
			}
		} else {
			if !errors.Is(err, scrapers.ErrNotLive) {
				log.Println("[YouTube] Unexpected error while scraping YouTube live:", err.Error())
			}
		}

		// Wait up to two minutes, then check again
		time.Sleep(time.Minute + time.Duration(rand.Intn(60))*time.Second)
	}
}

// checkHomeTimeline requests the user home timeline about every minute and puts all new tweets in tweetChan.
// it also includes replies which would normally not be shown in the timeline.
// TL;DR: it stalks all users the account follows, even their replies
func checkHomeTimeline(client *twitter.Client, tweetChan chan<- twitter.Tweet) {
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

// checkListTimeline requests the given lists about every minute or so. Any new tweets are put in tweetChan.
func checkListTimeline(client *twitter.Client, list twitter.List, tweetChan chan<- twitter.Tweet) {
	defer panic("list (" + list.Name + ") follower stopped processing even though it shouldn't")

	var (
		// lastSeenID is the ID of the last tweet we saw
		lastSeenID int64

		// The first batch of tweets we receive should not acted upon
		isFirstRequest bool = true
	)
	for {
		// https://developer.twitter.com/en/docs/twitter-api/v1/accounts-and-users/create-manage-lists/api-reference/get-lists-statuses
		tweets, _, err := client.Lists.Statuses(&twitter.ListsStatusesParams{
			ListID: list.ID,

			IncludeRetweets: twitter.Bool(true),
			IncludeEntities: twitter.Bool(false), // We also don't care about who was mentioned etc.
			SinceID:         lastSeenID,          // everything since our last request
			Count:           200,                 // Maximum number of tweets we can get at once
		})

		if err != nil {
			logError(err, "list "+list.FullName)
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

		// Add a random delay
		time.Sleep(time.Minute + time.Duration(rand.Intn(45))*time.Second)
	}
}

// retweet retweets the given tweet, but if it fails it doesn't care
func retweet(client *twitter.Client, tweet *twitter.Tweet) {
	if tweet.Retweeted {
		return
	}

	_, _, err := client.Statuses.Retweet(tweet.ID, nil)
	if err != nil {
		logError(err, "retweet")
		return
	}

	twurl := tweetURL(tweet)
	log.Println("[Twitter] Retweeted", twurl)

	// Setting Retweeted can help processThread to detect that it should stop
	tweet.Retweeted = true
}

func tweetURL(tweet *twitter.Tweet) string {
	if tweet.User == nil {
		return "https://twitter.com/i/status/" + tweet.IDStr
	}
	return "https://twitter.com/" + tweet.User.ScreenName + "/status/" + tweet.IDStr
}

func logError(err error, location string) {
	if err != nil {
		log.Printf("[Error (%s)]: %s\n", location, err.Error())
	}
}

// processThread processes tweet threads retweets everything in them if they are on-topif.
// This is useful because Elon Musk often replies to people that quote tweeted/asked a questions on his tweets
// See this for example: https://twitter.com/elonmusk/status/1372826575293583366
// or here: https://twitter.com/elonmusk/status/1372725108909957121
func processThread(client *twitter.Client, tweet *twitter.Tweet, seenTweets map[int64]bool) (didRetweet bool) {
	if tweet == nil {
		// Just in case
		return false
	}

	// Was that tweet interesting the last time we saw it?
	// If yes, then we should probably retweet the next stuff.
	// If not, we can stop here because it won't get any better
	// (we already checked the last time if it's good)
	if seenTweets[tweet.ID] || tweet.Retweeted {
		return tweet.Retweeted
	}
	seenTweets[tweet.ID] = true

	// First process the rest of the thread
	if tweet.InReplyToStatusID != 0 {
		// Ok, there was a reply. Check if we can do something with that
		t, _, err := client.Statuses.Lookup([]int64{tweet.InReplyToStatusID}, &twitter.StatusLookupParams{
			IncludeEntities: twitter.Bool(false),
			TweetMode:       "extended",
		})
		logError(err, "tweet reply status fetch (processThread)")

		if len(t) > 0 {
			if !processThread(client, &t[0], seenTweets) {
				return false
			}

			seenTweets[t[0].ID] = true
			retweet(client, &t[0])
			didRetweet = true
		}
	}

	// A quoted tweet. Let's see if there's anything interesting
	if tweet.QuotedStatusID != 0 && tweet.QuotedStatus != nil {
		seenTweets[tweet.QuotedStatusID] = true

		return processThread(client, tweet.QuotedStatus, seenTweets)
	}

	// Now actually match the tweet
	if didRetweet || match.StarshipTweet(tweet) {

		retweet(client, tweet)

		return true
	}

	return didRetweet
}
