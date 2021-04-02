package main

import (
	"flag"
	"log"
	"math/rand"
	"time"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/xarantolus/spacex-hop-bot/bot"
	"github.com/xarantolus/spacex-hop-bot/config"
	"github.com/xarantolus/spacex-hop-bot/consumer"
	"github.com/xarantolus/spacex-hop-bot/jobs"
	"github.com/xarantolus/spacex-hop-bot/match"
)

var (
	flagConfigFile = flag.String("cfg", "config.yaml", "Config file path")
	flagDebug      = flag.Bool("debug", false, "Debug mode disables background jobs")
)

func main() {
	flag.Parse()

	var dbg string
	if *flagDebug {
		dbg = " in debug mode"
	}

	log.Printf("[Startup] Bot is starting%s\n", dbg)

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
		var linkChan = make(chan string, 2)

		// Run YouTube scraper in the background,
		// it will tweet if it discovers that SpaceX is online with a Starship stream
		go jobs.CheckYouTubeLive(client, selfUser, linkChan)

		// When the webpage mentions a new date/starship, we tweet about that
		go jobs.StarshipWebsiteChanges(client, linkChan)

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
		var watchedLists int
		for _, l := range lists {
			if l.ID == match.SatireListID {
				continue
			}
			go jobs.CheckListTimeline(client, l, tweetChan)
			watchedLists++
		}

		log.Printf("[Twitter] Started watching %d lists\n", watchedLists)

		match.LoadSatireList(client)
	}

	const spacePeopleListID = 1375480259840212997

	// proc handles tweets by filtering & retweeting the interesting ones
	var proc = consumer.NewProcessor(*flagDebug, client, selfUser, spacePeopleListID)

	// Now we just pass all tweets to processTweet
	for tweet := range tweetChan {
		proc.Tweet(tweet)
	}
}
