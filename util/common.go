package util

import (
	"log"

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
