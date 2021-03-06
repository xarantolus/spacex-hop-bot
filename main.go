package main

import (
	"flag"
	"log"
	"math/rand"
	"time"

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

	// Let's parse our configuration file
	cfg, err := config.Parse(*flagConfigFile)
	if err != nil {
		panic("parsing configuration file: " + err.Error())
	}

	// Log in to Twitter
	client, selfUser, err := bot.Login(cfg)
	if err != nil {
		panic("logging in to twitter: " + err.Error())
	}
	log.Printf("[Twitter] Logged in @%s\n", selfUser.ScreenName)

	// Twitter list ids for lists we need
	const (
		spacePeopleListID = 1375480259840212997
	)

	// We have different account lists that define ignored accounts.
	// That way I don't have to mute them.
	var ignoredLists = map[int64]bool{
		// Satire
		1377136100574064647: true,

		// Unrelated to SpaceX
		1411299241050488835: true,

		// Animation, renders etc.
		1396191591686184967: true,

		// Other/Muted
		1410664386943930368: true,
	}

	// Load all ignored accounts to make sure we don't retweet them
	for lid := range ignoredLists {
		match.LoadIgnoredList(client, lid)
	}

	// The bot should check all tweets that are sent on this channel
	var tweetChan = make(chan match.TweetWrapper, 250)

	if *flagDebug {
		log.Println("[Info] Running in debug mode, no background jobs are started")
	} else {
		// Register all background jobs, most of them send tweets on tweetChan
		jobs.Register(client, selfUser, tweetChan, ignoredLists)
	}

	// handler handles tweets by filtering & retweeting the interesting ones
	var handler = consumer.NewProcessor(*flagDebug, client, selfUser, spacePeopleListID)

	// Now we just process every tweet we come across
	for tweet := range tweetChan {
		handler.Tweet(tweet)
	}
}
