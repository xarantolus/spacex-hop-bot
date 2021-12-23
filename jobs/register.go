package jobs

import (
	"fmt"
	"log"
	"strings"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/xarantolus/spacex-hop-bot/match"
)

func Register(client *twitter.Client, selfUser *twitter.User, matcher *match.StarshipMatcher, tweetChan chan match.TweetWrapper, skipLists map[int64]bool) (err error) {
	var linkChan = make(chan string, 2)

	// Run YouTube scraper in the background,
	// it will tweet if it discovers that SpaceX is online with a Starship stream
	go CheckYouTubeLive(client, selfUser, matcher, linkChan)

	// When the gov dashboard changes, we want to tweet about it
	go CheckDashboard(client)

	// When the webpage mentions a new date/starship, we tweet about that
	go StarshipWebsiteChanges(client, linkChan)

	// Check out the home timeline of the bot user, it will contain all kinds of tweets from all kinds of people
	go CheckHomeTimeline(client, tweetChan)

	// Get tweets from the general area around boca chica
	go CheckLocationStream(client, tweetChan)

	// Make we get all tweets from certain users, before this we sometimes missed stuff
	go CheckUserTimeline(client, "elonmusk", tweetChan)
	go CheckUserTimeline(client, "SpaceX", tweetChan)

	// Start watching all lists the bot account follows
	lists, _, err := client.Lists.List(&twitter.ListsListParams{})
	if len(lists) == 100 {
		// See https://developer.twitter.com/en/docs/twitter-api/v1/accounts-and-users/create-manage-lists/api-reference/get-lists-list
		log.Println("[Warning] Lists API call returned 100 lists, which means that it is likely that some lists were not included. See API URL in comment above this line")
	}
	if err != nil {
		return fmt.Errorf("initializing bot: couldn't retrieve lists: %s", err.Error())
	}

	// Those are also background jobs
	var watchedLists int
	var listNames []string
	for _, l := range lists {
		// We cannot request tweets from a list created by a protected user
		if skipLists[l.ID] || l.User != nil && l.User.Protected {
			continue
		}
		go CheckListTimeline(client, l, tweetChan)

		listNames = append(listNames, fmt.Sprintf("%s's %q", l.User.ScreenName, l.Name))
		watchedLists++
	}

	log.Printf("[Twitter] Started watching %d lists (%s)\n", watchedLists, strings.Join(listNames, ", "))

	return nil
}
