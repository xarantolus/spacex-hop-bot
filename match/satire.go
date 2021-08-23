package match

import (
	"log"
	"strings"

	"github.com/dghubble/go-twitter/twitter"
)

var (
	// These are loaded from a list
	ignoredNames = []string{}

	ignoredKeywords = []string{
		"parody", "joke", "blender", "3d", "render", "animat", /* e/ion */
	}
)

// LoadIgnoredList marks the members of this list as ignored accounts
func LoadIgnoredList(client *twitter.Client, ignoredListID int64) {
	users, _, err := client.Lists.Members(&twitter.ListsMembersParams{
		ListID: ignoredListID,
		Count:  1000,
	})
	if err != nil {
		log.Println("[Twitter] Failed loading ignored account list:", err.Error())
		return
	}

	for _, u := range users.Users {
		ignoredNames = append(ignoredNames, strings.ToLower(u.ScreenName))
	}
}

func isIgnoredName(username string) bool {

	for _, k := range ignoredNames {
		if username == k {
			return true
		}
	}
	return false
}

func IsOrMentionsIgnoredAccount(tweet *twitter.Tweet) bool {
	if tweet.User == nil {
		return false
	}

	username := strings.ToLower(tweet.User.ScreenName)

	// If we know the user, they can't be ignored
	_, known1 := specificUserMatchers[username]
	_, known2 := userAntikeywordsOverwrite[username]
	if known1 || known2 {
		return false
	}

	if isIgnoredName(username) {
		return true
	}

	desc := strings.ToLower(tweet.User.Description)
	for _, k := range ignoredKeywords {
		if strings.Contains(desc, k) {
			return true
		}
	}

	// If someone *mentions* an ignored user, it's likely not important
	if tweet.Entities != nil {
		for _, um := range tweet.Entities.UserMentions {
			if isIgnoredName(strings.ToLower(um.ScreenName)) {
				return true
			}
		}
	}

	return false
}
