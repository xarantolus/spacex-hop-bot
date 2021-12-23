package match

import (
	"strings"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/xarantolus/spacex-hop-bot/bot"
)

var (
	ignoredKeywords = []string{
		"parody", "joke", "blender", "3d", "render", "animat", /* e/ion */
	}
)

type Ignorer struct {
	list *bot.UserList

	keywords []string
}

// LoadIgnoredList marks the members of this list as ignored accounts
func LoadIgnoredList(client *twitter.Client, ignoredListIDs ...int64) *Ignorer {
	var list = bot.ListMembers(client, "ignored", ignoredListIDs...)

	return &Ignorer{
		list:     list,
		keywords: ignoredKeywords,
	}
}

func (i *Ignorer) IsOrMentionsIgnoredAccount(tweet *twitter.Tweet) bool {
	username := strings.ToLower(tweet.User.ScreenName)

	// If we know the user, they can't be ignored
	_, known1 := specificUserMatchers[username]
	_, known2 := userAntikeywordsOverwrite[username]
	if known1 || known2 {
		return false
	}

	// If the list of accounts we ignore contains *anything* related to this account
	// we ignore the tweet
	if i.list.TweetAssociatedWithAny(tweet) {
		return true
	}

	// Now search the user description to see if any negative keywords stand out
	desc := strings.ToLower(tweet.User.Description)
	for _, k := range ignoredKeywords {
		if strings.Contains(desc, k) {
			return true
		}
	}

	return false
}
