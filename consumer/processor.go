package consumer

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/xarantolus/spacex-hop-bot/match"
	"github.com/xarantolus/spacex-hop-bot/util"
)

// Processor handles tweets by looking at them and deciding whether to retweet.
type Processor struct {
	debug bool
	test  bool

	client TwitterClient

	matcher *match.StarshipMatcher

	selfUser *twitter.User

	// map[URL]last Retweet time
	seenLinks map[string]time.Time

	seenTweets map[int64]bool

	spacePeopleListID      int64
	spacePeopleListMembers map[int64]bool
}

const (
	articlesFilename = "articles.json"
	// How long after we've seen a link will we allow it to be retweeted again?
	seenLinkDelay = 12 * time.Hour
)

// NewProcessor returns a new processor with the given options
func NewProcessor(debug bool, inTest bool, client TwitterClient, selfUser *twitter.User, matcher *match.StarshipMatcher, spacePeopleListID int64) *Processor {
	p := &Processor{
		debug:   debug,
		test:    inTest,
		matcher: matcher,

		client:   client,
		selfUser: selfUser,

		seenLinks: make(map[string]time.Time),

		spacePeopleListID: spacePeopleListID,

		seenTweets:             make(map[int64]bool),
		spacePeopleListMembers: make(map[int64]bool),
	}

	if !p.test {
		util.LogError(util.LoadJSON(articlesFilename, &p.seenLinks), "loading links")

		p.cleanup(true)
	}

	return p
}

// Tweet processes the given tweet and checks whether it should be retweeted.
// Tweets that have already been seen are ignored. It is not safe for concurrent use.
func (p *Processor) Tweet(tweet match.TweetWrapper) {
	// So now we got a tweet. There are three categories that interest us:
	// 1. Elon Musk drops insider info about starship, e.g. as a reply.
	//    We do not care about his other tweets, so we check if any tweet
	//    in the reply chain matches stuff about starship
	// 2. We find a retweet
	// 3. We find a quoted tweet
	// 4. We find a tweet that is about starship

	if (p.seenTweets[tweet.ID] || tweet.Retweeted) && !p.debug {
		return
	}

	// Skip our own tweets
	if tweet.User != nil && tweet.User.ID == p.selfUser.ID {
		return
	}

	switch {
	case isElonTweet(tweet):
		// When elon drops starship info, we want to retweet it.
		// We basically detect if the thread/tweet is about starship and
		// retweet everything that is appropriate
		p.thread(&tweet.Tweet)
	case tweet.RetweetedStatus != nil:
		p.Tweet(match.TweetWrapper{TweetSource: tweet.TweetSource, Tweet: *tweet.RetweetedStatus})
	case tweet.QuotedStatus != nil:
		// If someone quotes a tweet, we check some things.

		// If we have a Starship-Tweet quoting a tweet that does not contain antikeywords,
		// we assume that the quoted tweet also contains relevant information

		if p.seenTweets[tweet.QuotedStatusID] {
			break
		}

		// If the quoted tweet already is about starship, we maybe only look at that one
		quotedWrap := match.TweetWrapper{TweetSource: tweet.TweetSource, Tweet: *tweet.QuotedStatus}
		if p.isStarshipTweet(quotedWrap) {
			// If it's from the *same* user, then we just assume they added additional info.
			// We only retweet if it's media though
			if sameUser(&tweet.Tweet, tweet.QuotedStatus) {
				if hasMedia(tweet.QuotedStatus) {
					p.retweet(tweet.QuotedStatus, "quoted media", tweet.TweetSource)
				}
			} else {
				p.Tweet(quotedWrap)
			}
		}

		// The quoting tweet should be about starship AND have media
		if !(p.isStarshipTweet(tweet) && hasMedia(&tweet.Tweet)) {
			break
		}

		// Make sure the quoted user is not ignored
		if p.matcher.IsOrMentionsIgnoredAccount(tweet.QuotedStatus) {
			break
		}

		// Now we have a tweet about starship, that we haven't seen/retweeted before,
		// that quotes another tweet
		p.retweet(&tweet.Tweet, "quoted", tweet.TweetSource)

		p.seenTweets[tweet.QuotedStatusID] = true
	case tweet.InReplyToStatusID != 0:
		parentTweet, err := p.client.LoadStatus(tweet.InReplyToStatusID)
		if err != nil {
			// Most errors happen because we're not allowed to see protected accounts' tweets.
			// We don't log these errors
			if !strings.Contains(err.Error(), "179 Sorry, you are not authorized to see this status.") &&
				!strings.Contains(err.Error(), "144 No status found with that ID") {
				util.LogError(err, "loading parent of "+util.TweetURL(&tweet.Tweet))
			}
			break
		}

		// If
		// - we retweeted the parent tweet (so must be starship related)
		//      OR the current tweet is Starship-related
		// - the current tweet doesn't contain antiKeywords
		// - it's from the same user who started the thread (looking through all tweets above parent tweet)
		// then we want to go to the retweeting part below
		isStarshipTweet := p.matcher.StarshipTweet(tweet)
		hasAntiKeywords := match.ContainsStarshipAntiKeyword(tweet.Text())

		if !(((parentTweet.Retweeted && hasMedia(&tweet.Tweet)) ||
			isStarshipTweet) &&
			!hasAntiKeywords &&
			!isReactionGIF(&tweet.Tweet) &&
			sameUser(parentTweet, &tweet.Tweet) &&
			!p.isReply(parentTweet)) {
			break
		}

		fallthrough
	case p.isStarshipTweet(tweet):
		// If the tweet itself is about starship, we retweet it
		// We already filtered out replies, which is important because we don't want to
		// retweet every question someone posts under an elon post, only those that
		// elon responded to.
		// Then we also filter out all tweets that tag elon musk, e.g. there could be someone
		// just tweeting something like "Do you think xyz... @elonmusk"

		// Filter out non-english tweets (except for location stream)
		if tweet.TweetSource != match.TweetSourceLocationStream && tweet.Lang != "" && tweet.Lang != "en" && tweet.Lang != "und" {
			log.Printf("[Processor] Ignoring %s because of language %s (%s)", util.TweetURL(&tweet.Tweet), tweet.Lang, tweet.TweetSource)
			break
		}

		if p.shouldIgnoreLink(&tweet.Tweet) {
			log.Printf("[Processor] Ignoring %s because we have seen this link recently (%s)", util.TweetURL(&tweet.Tweet), tweet.TweetSource)
			break
		}

		if tweet.Tweet.PossiblySensitive {
			log.Printf("[Processor] Ignoring %s because it is possibly sensitive (%s)", util.TweetURL(&tweet.Tweet), tweet.TweetSource)
			break
		}

		// Depending on the tweet source, we require media
		if tweet.TweetSource == match.TweetSourceLocationStream {
			switch {
			case hasMedia(&tweet.Tweet):
				// If it's from the location stream, matches etc. and has media
				p.retweet(&tweet.Tweet, "normal + location media", tweet.TweetSource)
			case match.IsPadAnnouncement(tweet.Text()):
				// If we have a pad announcement - those are usually tweets without media
				p.retweet(&tweet.Tweet, "location + pad announcement", tweet.TweetSource)
			default:
				log.Printf("[Processor] Ignoring %s because it's from the location stream and has no media", util.TweetURL(&tweet.Tweet))
			}
		} else {
			// If a tweet contains *only hashtags*, we only retweet it if it has media
			if isTagsOnly(tweet.Text()) {
				if hasMedia(&tweet.Tweet) {
					p.retweet(&tweet.Tweet, "normal matcher, only tags, but media", tweet.TweetSource)
				}
			} else {
				p.retweet(&tweet.Tweet, "normal matcher", tweet.TweetSource)
			}
		}
	}

	p.seenTweets[tweet.ID] = true
}

// retweet retweets the given tweet, but if it fails it doesn't care
func (p *Processor) retweet(tweet *twitter.Tweet, reason string, source match.TweetSource) {
	// If we have already retweeted a tweet, we don't try to do it again, that just leads to errors
	if tweet.Retweeted || tweet.RetweetedStatus != nil && tweet.RetweetedStatus.Retweeted {
		return
	}

	err := p.client.Retweet(tweet)
	if err != nil {
		// Twitter often doesn't send the info that we have already retweeted a tweet.
		// So here we don't log the error if that's the case
		if !strings.Contains(err.Error(), "327 You have already retweeted this Tweet.") {
			util.LogError(err, "retweeting "+util.TweetURL(tweet))
		}
		return
	}

	if !p.test {
		// save tweet so we can reproduce why it was matched
		p.saveTweet(tweet)

		// Add the user to our space people list
		// We ignore those from the location stream as they might not always tweet about starship
		if source != match.TweetSourceLocationStream {
			p.addSpaceMember(tweet)
		}

		twurl := util.TweetURL(tweet)
		log.Printf("[Twitter] Retweeted %s (%s - %s)", twurl, reason, source.String())
	}

	// Setting Retweeted can help thread to detect that it should stop
	tweet.Retweeted = true
}

// thread processes tweet threads and retweets everything on-topic.
// This is useful because Elon Musk often replies to people that quote tweeted/asked a questions on his tweets
// See this for example: https://twitter.com/elonmusk/status/1372826575293583366
// or here: https://twitter.com/elonmusk/status/1372725108909957121
func (p *Processor) thread(tweet *twitter.Tweet) (didRetweet bool) {
	if tweet == nil {
		// Just in case
		return false
	}

	// Was that tweet interesting the last time we saw it?
	// If yes, then we should probably retweet the next stuff.
	if tweet.Retweeted {
		return true
	}
	p.seenTweets[tweet.ID] = true

	// First process the rest of the thread
	if tweet.InReplyToStatusID != 0 {
		// Ok, there was a reply. Check if we can do something with that
		parent, err := p.client.LoadStatus(tweet.InReplyToStatusID)
		util.LogError(err, "tweet reply status fetch (thread)")

		// If we have a matching tweet thread
		if parent != nil && p.thread(parent) {
			p.seenTweets[parent.ID] = true
			p.retweet(parent, "thread: matched parent", match.TweetSourceUnknown)
			didRetweet = true
		}
	}

	// A quoted tweet. Let's see if there's anything interesting
	if tweet.QuotedStatusID != 0 && tweet.QuotedStatus != nil {
		return p.thread(tweet.QuotedStatus)
	}

	realTweet := tweet
	if tweet.RetweetedStatus != nil {
		realTweet = tweet.RetweetedStatus
	}

	// Now actually match the tweet
	if didRetweet || p.matcher.StarshipTweet(match.TweetWrapper{TweetSource: match.TweetSourceUnknown, Tweet: *realTweet}) {
		p.retweet(tweet, "thread: matched", match.TweetSourceUnknown)
		return true
	}

	return didRetweet
}

// addSpaceMember adds the user of the given tweet to the space people list
func (p *Processor) addSpaceMember(tweet *twitter.Tweet) {
	if tweet.User == nil || p.spacePeopleListMembers[tweet.User.ID] {
		return
	}

	p.spacePeopleListMembers[tweet.User.ID] = true

	err := p.client.AddListMember(p.spacePeopleListID, tweet.User.ID)
	util.LogError(err, fmt.Sprintf("adding %s to list", tweet.User.ScreenName))
}
