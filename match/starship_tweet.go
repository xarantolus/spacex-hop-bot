package match

import (
	"strings"
	"time"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/xarantolus/spacex-hop-bot/util"
)

// StarshipTweet returns whether the given tweet mentions starship. It also includes custom matchers for certain users
func (m *StarshipMatcher) StarshipTweet(tweet TweetWrapper) bool {
	// Ignore OLD tweets
	if d, err := tweet.CreatedAtTime(); err == nil && time.Since(d) > 24*time.Hour {
		return false
	}

	text := tweet.Text()

	// We do not care about tweets that are timestamped with a text more than 24 hours ago
	// e.g. if someone posts a photo and then writes "took this on March 15, 2002"
	if d, ok := util.ExtractDate(text, time.Now()); ok && time.Since(d) > 48*time.Hour {
		return false
	}

	_, isVeryImportant := veryImportantAccounts[strings.ToLower(tweet.User.ScreenName)]
	// We ignore certain (e.g. satire, artist) accounts
	if tweet.User != nil {
		if !isVeryImportant && m.IsOrMentionsIgnoredAccount(&tweet.Tweet) {
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

	// If the tweet is tagged with Starbase as location, we just retweet it.
	if !containsBadWords && IsAtSpaceXSite(&tweet.Tweet) {
		return true
	}
	// In case of antikeywords being present in a tweet at a starship location, we will retweet the tweet anyways if it has media
	if hasMedia(&tweet.Tweet) && IsAtStarshipLocation(&tweet.Tweet) {
		return true
	}

	// Stop if we have antikeywords. However, if e.g. elon tweets about tesla *and* spacex, it should still go to the specificUserMatcher below
	if containsBadWords && !isVeryImportant {
		return false
	}

	// Now check if it mentions too many people
	if strings.Count(text, "@") > 10 {
		return false
	}

	// ignore b4 when lowercase, as it's an abbreviation of "before"
	if strings.Contains(tweet.Text(), "b4") {
		tweet.FullText = strings.ReplaceAll(tweet.Text(), "b4", "")
		text = strings.ToLower(tweet.FullText)
	}

	// Check if the text matches
	if m.StarshipText(text, antiKeywords, false) {
		return true
	}
	// If the text didn't match, maybe it is matched when we don't remove URLs from it.
	// We do want to be a bit more careful here, because URLs can contain tricky sequences
	// of characters that could trick simple matchers (e.g. t.co/s20_513)
	if m.StarshipText(tweet.TextWithURLs(), antiKeywords, true) {
		return true
	}

	// Now check if we have a matcher for this specific user.
	// These users usually post high-quality information
	if tweet.User != nil {
		regexes, ok := specificUserMatchers[strings.ToLower(tweet.User.ScreenName)]
		if ok {
			// If at least one regex matches, we have a match
			for _, m := range regexes {
				if m.MatchString(text) {
					return true
				}
			}
		}

		// There are some accounts that always post high-quality pictures and videos.
		// For them we retweet *everything* that has media
		if hqMediaAccounts[strings.ToLower(tweet.User.ScreenName)] {
			return hasMedia(&tweet.Tweet)
		}
	}

	if tweet.Place != nil {
		pkw, ok := locationKeywords[tweet.Place.ID]
		if ok {
			if startsWithAny(text, pkw...) {
				return true
			}
		}
	}

	return false
}

func hasMedia(tweet *twitter.Tweet) bool {
	return tweet.Entities != nil && len(tweet.Entities.Media) > 0 || tweet.ExtendedEntities != nil && len(tweet.ExtendedEntities.Media) > 0
}
