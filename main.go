package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"strings"
	"time"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/xarantolus/spacex-hop-bot/bot"
	"github.com/xarantolus/spacex-hop-bot/config"
	"github.com/xarantolus/spacex-hop-bot/jobs"
	"github.com/xarantolus/spacex-hop-bot/match"
	"github.com/xarantolus/spacex-hop-bot/util"
)

var (
	flagConfigFile = flag.String("cfg", "config.yaml", "Config file path")
	flagDebug      = flag.Bool("debug", false, "Debug mode disables background jobs")
)

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

	if *flagDebug {
		log.Println("[Info] Running in debug mode, no background jobs are started")
	} else {
		// Run YouTube scraper in the background,
		// it will tweet if it discovers that SpaceX is online with a Starship stream
		go jobs.CheckYouTubeLive(client, selfUser)

		// When the webpage mentions a new date/starship, we tweet about that
		go jobs.StarshipWebsiteChanges(client)

		// Check out the home timeline of the bot user, it will contain all kinds of tweets from all kinds of people
		go jobs.CheckHomeTimeline(client, tweetChan)

		// Get tweets from the general area around boca chica
		go jobs.CheckLocationStream(client, tweetChan)

		// Start watching all lists the bot account follows
		lists, _, err := client.Lists.List(&twitter.ListsListParams{})
		if len(lists) == 100 {
			// See https://developer.twitter.com/en/docs/twitter-api/v1/accounts-and-users/create-manage-lists/api-reference/get-lists-list
			log.Println("[Warning] Lists API call returned 100 lists, which means that it is likely that some lists were not included. See API URL in comment above this line")
		}
		if err != nil {
			panic("initializing bot: couldn't retrieve lists: " + err.Error())
		}

		// Those are also background jobs
		for _, l := range lists {
			go jobs.CheckListTimeline(client, l, tweetChan)
		}

		log.Printf("[Twitter] Started watching %d lists\n", len(lists))
	}

	var (
		// seenTweets maps the tweet id to an boolean. If the tweet was already processed/seen, it is put here
		seenTweets = make(map[int64]bool)
	)

	// Now we just pass all tweets to processTweet
	for tweet := range tweetChan {
		processTweet(client, seenTweets, selfUser, tweet)
	}
}

func processTweet(client *twitter.Client, seenTweets map[int64]bool, selfUser *twitter.User, tweet twitter.Tweet) {
	// So now we got a tweet. There are three categories that interest us:
	// 1. Elon Musk drops insider info about starship, e.g. as a reply.
	//    We do not care about his other tweets, so we check if any tweet
	//    in the reply chain matches stuff about starship
	// 2. We find a retweet
	// 3. We find a quoted tweet
	// 4. We find a tweet that is about starship

	if seenTweets[tweet.ID] || tweet.Retweeted {
		return
	}

	// Skip our own tweets
	if tweet.User != nil && tweet.User.ID == selfUser.ID {
		return
	}

	switch {
	case tweet.User != nil && tweet.User.ScreenName == "elonmusk":
		// When elon drops starship info, we want to retweet it.
		// We basically detect if the thread/tweet is about starship and
		// retweet everything that is appropriate
		processThread(client, &tweet, seenTweets)
	case tweet.QuotedStatus != nil:
		// If someone quotes a tweet, we only check the tweet that was quoted.
		processTweet(client, seenTweets, selfUser, *tweet.QuotedStatus)
	case tweet.RetweetedStatus != nil:
		processTweet(client, seenTweets, selfUser, *tweet.RetweetedStatus)
	case match.StarshipTweet(&tweet) && !isReply(&tweet) && !isQuestion(&tweet):
		// If the tweet itself is about starship, we retweet it
		// We already filtered out replies, which is important because we don't want to
		// retweet every question someone posts under an elon post, only those that
		// elon responded to.
		// Then we also filter out all tweets that tag elon musk, e.g. there could be someone
		// just tweeting something like "Do you think xyz... @elonmusk"
		retweet(client, &tweet)
	}

	seenTweets[tweet.ID] = true
}

// isReply returns if the given tweet is a reply to another user
func isReply(t *twitter.Tweet) bool {
	if t.User == nil || t.InReplyToStatusID == 0 {
		return false
	}

	return t.User.ID != t.InReplyToUserID
}

func isQuestion(tweet *twitter.Tweet) bool {
	return strings.Contains(strings.ToLower(tweet.FullText), "@") && strings.Contains(tweet.FullText, "?")
}

const spacePeopleListID = 1375480259840212997

var spacePeopleMembers = map[int64]bool{}

// retweet retweets the given tweet, but if it fails it doesn't care
func retweet(client *twitter.Client, tweet *twitter.Tweet) {
	if tweet.Retweeted {
		return
	}

	_, _, err := client.Statuses.Retweet(tweet.ID, nil)
	if err != nil {
		util.LogError(err, "retweeting "+util.TweetURL(tweet))
		return
	}

	// Collect interesting people in a list
	if tweet.User != nil && !spacePeopleMembers[tweet.User.ID] {
		spacePeopleMembers[tweet.User.ID] = true

		_, err := client.Lists.MembersCreate(&twitter.ListsMembersCreateParams{
			ListID: spacePeopleListID,
			UserID: tweet.User.ID,
		})
		util.LogError(err, fmt.Sprintf("adding %s to list", tweet.User.ScreenName))
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

		// If we have a matching tweet thread
		if len(t) > 0 && processThread(client, &t[0], seenTweets) {
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
