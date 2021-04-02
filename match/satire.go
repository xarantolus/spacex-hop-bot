package match

import (
	"log"
	"strings"

	"github.com/dghubble/go-twitter/twitter"
)

var (
	// See https://twitter.com/ULASeagull/status/1376913976362217472 and
	// https://twitter.com/i/lists/1357527189370130432 for a list of names
	satireNames = []string{}

	satireKeywords = []string{
		"parody", "joke",
	}
)

// LoadSatireList marks the members of this list as satire accounts
func LoadSatireList(client *twitter.Client, satireListID int64) {
	users, _, err := client.Lists.Members(&twitter.ListsMembersParams{
		ListID: satireListID,
		Count:  1000,
	})
	if err != nil {
		log.Println("[Twitter] Failed loading satire account list:", err.Error())
		return
	}

	for _, u := range users.Users {
		satireNames = append(satireNames, strings.ToLower(u.ScreenName))
	}
}

func isSatireAccount(tweet *twitter.Tweet) bool {
	if tweet.User == nil {
		return false
	}
	// If we know the user, it can't be satire
	_, known1 := specificUserMatchers[tweet.User.ScreenName]
	_, known2 := usersWithNoAntikeywords[tweet.User.ScreenName]
	if known1 || known2 {
		return false
	}

	username := strings.ToLower(tweet.User.ScreenName)

	for _, k := range satireNames {
		if username == k {
			return true
		}
	}

	desc := strings.ToLower(tweet.User.Description)
	for _, k := range satireKeywords {
		if strings.Contains(desc, k) {
			return true
		}
	}

	return false
}
