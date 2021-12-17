package match

import (
	"log"
	"strings"
	"time"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/xarantolus/spacex-hop-bot/util"
)

// The faceRatio of a tweet is the number of faces in all images (or video thumbnails) divided by the number of images in the tweet
const maxFaceRatio = .75

var faceDetector = NewFaceDetector()

// StarshipTweet returns whether the given tweet mentions starship. It also includes custom matchers for certain users
func StarshipTweet(tweet TweetWrapper) bool {
	// Ignore OLD tweets
	if d, err := tweet.CreatedAtTime(); err == nil && time.Since(d) > 24*time.Hour {
		return false
	}

	text := tweet.Text()

	// We do not care about tweets that are timestamped with a text more than 24 hours ago
	// e.g. if someone posts a photo and then writes "took this on March 15, 2002"
	if d, ok := util.ExtractDate(text); ok && time.Since(d) > 48*time.Hour {
		return false
	}

	// We ignore certain (e.g. satire, artist) accounts
	if tweet.User != nil {
		if _, important := veryImportantAccounts[strings.ToLower(tweet.User.Name)]; !important && IsOrMentionsIgnoredAccount(&tweet.Tweet) {
			return false
		}
	}

	// Now check if the text of the tweet matches what we're looking for.
	text = strings.ToLower(text)

	// Depending on the user, we use different antiKeywords
	antiKeywords := antiStarshipKeywords
	if tweet.User != nil {
		ak, ok := userAntikeywordsOverwrite[strings.ToLower(tweet.User.ScreenName)]
		if ok {
			antiKeywords = ak
		}
	}

	var containsBadWords = containsAntikeyword(antiKeywords, text)

	// If the tweet is tagged with Starbase as location, we just retweet it
	if !containsBadWords && IsAtSpaceXSite(&tweet.Tweet) {
		fr := faceDetector.FaceRatio(&tweet.Tweet)
		log.Printf("[FaceRatio] %s: %f\n", util.TweetURL(&tweet.Tweet), fr)
		return fr <= maxFaceRatio
	}

	// If the tweet mentions raptor without images, we still retweet it.
	// This is mostly for tweets from SpaceX McGregor
	if !containsBadWords && strings.Contains(text, "raptor") && IsAtSpaceXSite(&tweet.Tweet) {
		fr := faceDetector.FaceRatio(&tweet.Tweet)
		log.Printf("[FaceRatio] %s: %f\n", util.TweetURL(&tweet.Tweet), fr)
		return fr <= maxFaceRatio
	}

	// Now check if it mentions too many people
	if strings.Count(text, "@") > 5 {
		return false
	}

	// ignore b4 when lowercase, as it's an abbreviation of "before"
	if strings.Contains(tweet.Text(), "b4") {
		log.Println("Ignored b4 tweet", util.TweetURL(&tweet.Tweet))
		return false
	}

	// Check if the text matches
	if StarshipText(text, antiKeywords) {
		fr := faceDetector.FaceRatio(&tweet.Tweet)
		log.Printf("[FaceRatio] %s: %f\n", util.TweetURL(&tweet.Tweet), fr)
		return fr <= maxFaceRatio
	}

	// Now check if we have a matcher for this specific user.
	// These users usually post high-quality information
	if tweet.User != nil {
		m, ok := specificUserMatchers[strings.ToLower(tweet.User.ScreenName)]
		if ok {
			return m.MatchString(text)
		}

		// There are some accounts that always post high-quality pictures and videos.
		// For them we retweet *everything* that has media
		if hqMediaAccounts[strings.ToLower(tweet.User.ScreenName)] {
			return hasMedia(&tweet.Tweet)
		}

		// If the user mentions a raptor engine keyword (however not all from raptorKeywords)
		if ok && startsWithAny(text, "raptor", "rb", "rc", "rvac") {
			return true
		}
	}

	return false
}

func hasMedia(tweet *twitter.Tweet) bool {
	return tweet.Entities != nil && len(tweet.Entities.Media) > 0 || tweet.ExtendedEntities != nil && len(tweet.ExtendedEntities.Media) > 0
}