package match

import (
	"strings"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/xarantolus/spacex-hop-bot/bot"
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
		keywords: ignoredAccountDescriptionKeywords,
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
	for _, k := range i.keywords {
		if _, contains := startsWithAny(desc, k); contains {
			return true
		}
	}

	return false
}

func (i Ignorer) UserIDs() []int64 {
	return i.list.ContainedIDs()
}
