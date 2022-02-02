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
	"github.com/xarantolus/spacex-hop-bot/util"
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

	// Some stuff depends on randomness, so here we seed it
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

	// Load all ignored accounts to make sure we don't retweet them
	ignoredUserMatcher := match.LoadIgnoredList(client, cfg.Lists.IgnoredListIDs...)
	iu := ignoredUserMatcher.UserIDs()
	if len(iu) > 5 {
		util.LogError(util.SaveJSON("ignored-users.json", iu), "startup: saving ignored users")
	}

	// Now create a matcher instance that ignores those accounts
	var starshipMatcher = match.NewStarshipMatcher(ignoredUserMatcher)

	var twitterClient consumer.TwitterClient = &consumer.NormalTwitterClient{
		Client: client,
		Debug:  *flagDebug,
	}

	// This is the main channel tweets will be sent on. Basically many jobs *search* for tweets
	// and send them on this channel, then the processor will handle each incoming tweet
	var tweetChan = make(chan match.TweetWrapper, 250)

	if *flagDebug {
		log.Println("[Info] Running in debug mode, no background jobs are started")
	} else {
		// Register all background jobs, most of them send tweets on tweetChan
		err = jobs.Register(client, twitterClient, selfUser, starshipMatcher, tweetChan, cfg.IgnoredListsMapping())
		if err != nil {
			panic("registering jobs: " + err.Error())
		}
	}

	// handler handles tweets by filtering & retweeting the interesting ones
	var handler = consumer.NewProcessor(*flagDebug, false, twitterClient, selfUser, starshipMatcher, cfg.Lists.MainStarshipListID)

	// The web server should always run, regardless of debug mode or not
	go jobs.RunWebServer(cfg, twitterClient, handler, tweetChan)

	// Now we just process every tweet we come across
	for tweet := range tweetChan {
		handler.Tweet(tweet)
	}
}
