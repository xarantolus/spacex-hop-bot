package main

import (
	"flag"
	"log"
	"math/rand"
	"time"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/xarantolus/spacex-hop-bot/bot"
	"github.com/xarantolus/spacex-hop-bot/config"
	"github.com/xarantolus/spacex-hop-bot/jobs"
	"github.com/xarantolus/spacex-hop-bot/match"
	"github.com/xarantolus/spacex-hop-bot/util"
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
	go jobs.CheckYouTubeLive(client, selfUser)

	{
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
			go jobs.CheckListTimeline(client, l, tweetChan)
		}

		log.Printf("[Twitter] Started watching %d lists\n", len(lists))
	}

	// Check out the home timeline of the bot user, it will contain all kinds of tweets from all kinds of people
	go jobs.CheckHomeTimeline(client, tweetChan)

	// Get tweets from the general area around boca chica
	go jobs.CheckLocationStream(client, tweetChan)

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
		case tweet.RetweetedStatus != nil && match.StarshipTweet(tweet.RetweetedStatus) && !isReply(tweet.RetweetedStatus):
			// If it's a retweet of someone, we check that tweet if it's interesting
			retweet(client, tweet.RetweetedStatus)
		case match.StarshipTweet(&tweet) && !isReply(&tweet):
			// If the tweet itself is about starship, we retweet it
			// We already filtered out replies, which is important because we don't want to
			// retweet every question someone posts under an elon post, only those that
			// elon responded to
			retweet(client, &tweet)
		}
	}
}

// isReply returns if the given tweet is a reply to another user
func isReply(t *twitter.Tweet) bool {
	if t.User == nil || t.InReplyToStatusID == 0 {
		return false
	}

	return t.User.ID != t.InReplyToStatusID
}

// retweet retweets the given tweet, but if it fails it doesn't care
func retweet(client *twitter.Client, tweet *twitter.Tweet) {
	if tweet.Retweeted {
		return
	}

	_, _, err := client.Statuses.Retweet(tweet.ID, nil)
	if err != nil {
		util.LogError(err, "retweet")
		return
	}

	twurl := util.TweetURL(tweet)
	log.Println("[Twitter] Retweeted", twurl)

	// Setting Retweeted can help processThread to detect that it should stop
	tweet.Retweeted = true
}

// processThread processes tweet threads and retweets everything on-topic.
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
		util.LogError(err, "tweet reply status fetch (processThread)")

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
