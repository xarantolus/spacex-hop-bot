package util

import (
	"log"
	"strings"

	"github.com/dghubble/go-twitter/twitter"
)

func LogError(err error, location string) {
	if err != nil {
		log.Printf("[Error (%s)]: %s\n", location, err.Error())
	}
}

// TweetURL returns the URL for this tweet
func TweetURL(tweet *twitter.Tweet) string {
	if tweet.User == nil {
		return "https://twitter.com/i/status/" + tweet.IDStr
	}
	return "https://twitter.com/" + tweet.User.ScreenName + "/status/" + tweet.IDStr
}

func HashTagText(words []string) string {
	var joinedWords []string

	// Replace spaces in words with nothing
	for _, w := range words {
		if len(strings.TrimSpace(w)) == 0 {
			continue
		}
		joinedWords = append(joinedWords, strings.Join(strings.Fields(w), ""))
	}

	if len(joinedWords) == 0 {
		return ""
	}

	return "#" + strings.Join(joinedWords, " #")
}
