package consumer

import (
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

	seenTweets      map[int64]bool
	retweetedTweets map[int64]bool

	spacePeopleListID      int64
	spacePeopleListMembers map[int64]bool

	startTime time.Time
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
		retweetedTweets:        make(map[int64]bool),
		spacePeopleListMembers: make(map[int64]bool),

		startTime: time.Now(),
	}

	if !p.test {
		util.LogError(util.LoadJSON(articlesFilename, &p.seenLinks), "loading links")

		p.cleanup(true)
	}

	return p
}

func (p *Processor) Stats() map[string]interface{} {
	return map[string]interface{}{
		"tweets_seen_count":      len(p.seenTweets),
		"tweets_retweeted_count": len(p.retweetedTweets),
		"user":                   p.selfUser,
		"seen_links":             p.seenLinks,
		"start_time":             p.startTime,
		"uptime":                 time.Since(p.startTime).String(),
	}
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

	if (p.seenTweets[tweet.ID] || tweet.Retweeted) && !(p.debug || tweet.EnableLogging) {
		tweet.Log("already saw this tweet")
		return
	}

	// Skip our own tweets
	if tweet.User != nil && tweet.User.ID == p.selfUser.ID {
		tweet.Log("is our own tweet")
		return
	}

	// Some tweets are truncated, this means that twitter did not send the full text of the tweet.
	// Not 100% sure why this happens, but it happens
	if tweet.Truncated {
		t, err := p.client.LoadStatus(tweet.ID)
		if err != nil || t == nil {
			util.LogError(err, "loading trucated status with id %d", tweet.ID)
			return
		}

		tweet = tweet.Wrap(t)
	}

	switch {
	case isElonTweet(tweet):
		// When elon drops starship info, we want to retweet it.
		// We basically detect if the thread/tweet is about starship and
		// retweet everything that is appropriate
		tweet.Log("is elon tweet")
		p.thread(&tweet.Tweet)
	case isSpaceXTweet(tweet):
		tweet.Log("is SpaceX tweet")
		if tweet.QuotedStatus != nil {
			p.Tweet(tweet.Wrap(tweet.QuotedStatus))
		}
		if tweet.RetweetedStatus != nil {
			p.Tweet(tweet.Wrap(tweet.RetweetedStatus))
		}
		if p.isStarshipTweet(tweet) {
			p.retweet(&tweet.Tweet, "SpaceX tweet", tweet.TweetSource)
		}
	case tweet.RetweetedStatus != nil:
		tweet.Log("is retweet")
		p.Tweet(tweet.Wrap(tweet.RetweetedStatus))
	case tweet.QuotedStatus != nil:
		tweet.Log("is quoting")
		// If someone quotes a tweet, we check some things.

		// If we have a Starship-Tweet quoting a tweet that does not contain antikeywords,
		// we assume that the quoted tweet also contains relevant information

		if p.seenTweets[tweet.QuotedStatusID] && !match.IsImportantAcount(tweet.User) {
			tweet.Log("already saw this quoted tweet")
			break
		}

		// If the quoted tweet already is about starship, we maybe only look at that one
		quotedWrap := tweet.Wrap(tweet.QuotedStatus)
		if p.isStarshipTweet(quotedWrap) {
			tweet.Log("quoted is starship tweet")
			// If it's from the *same* user, then we just assume they added additional info.
			// We only retweet if it's media though
			if sameUser(&tweet.Tweet, tweet.QuotedStatus) {
				tweet.Log("quoted is starship tweet with same user")
				if hasMedia(tweet.QuotedStatus) {
					tweet.Log("quoted is starship tweet with same user and has media")
					p.retweet(tweet.QuotedStatus, "quoted media", tweet.TweetSource)
				}
				if hasMedia(&tweet.Tweet) && p.isStarshipTweet(tweet) {
					tweet.Log("quoting starship tweet has media")
					p.retweet(&tweet.Tweet, "quoting starship tweet with media", tweet.TweetSource)
				}
			} else {
				tweet.Log("quoted is starship tweet with different user")
				p.Tweet(quotedWrap)
			}
		}

		// The quoting tweet should be about starship AND have media
		if !(p.isStarshipTweet(tweet) && (hasMedia(&tweet.Tweet) || match.IsImportantAcount(tweet.User))) {
			tweet.Log("quoting tweet is not starship tweet with media")
			break
		}

		// Make sure the quoted user is not ignored
		if p.matcher.IsOrMentionsIgnoredAccount(tweet.QuotedStatus) {
			tweet.Log("quoting tweet user ignored")
			break
		}

		// Now we have a tweet about starship, that we haven't seen/retweeted before,
		// that quotes another tweet
		p.retweet(&tweet.Tweet, "quoted", tweet.TweetSource)

		p.seenTweets[tweet.QuotedStatusID] = true
	case tweet.InReplyToStatusID != 0:
		tweet.Log("tweet is reply")

		parentTweet, err := p.client.LoadStatus(tweet.InReplyToStatusID)
		if err != nil {
			// Most errors happen because we're not allowed to see protected accounts' tweets.
			// We don't log these errors
			if !strings.Contains(err.Error(), "179 Sorry, you are not authorized to see this status.") &&
				!strings.Contains(err.Error(), "144 No status found with that ID") {
				util.LogError(err, "loading parent of %s", util.TweetURL(&tweet.Tweet))
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
		hasMedia := hasMedia(&tweet.Tweet)
		tweet.Log("reply is isStarshipTweet=%v, hasAntiKeywords=%v", isStarshipTweet, hasAntiKeywords)

		if !(((parentTweet.Retweeted && hasMedia) ||
			isStarshipTweet) &&
			!hasAntiKeywords &&
			!(isQuestion(&tweet.Tweet) && !hasMedia) &&
			!isReactionGIF(&tweet.Tweet) &&
			sameUser(parentTweet, &tweet.Tweet) &&
			!p.isReply(parentTweet)) {
			tweet.Log("reply does not match complex criteria")
			break
		}

		tweet.Log("reply did match complex criteria")
		fallthrough
	case p.isStarshipTweet(tweet):
		// If the tweet itself is about starship, we retweet it
		// We already filtered out replies, which is important because we don't want to
		// retweet every question someone posts under an elon post, only those that
		// elon responded to.
		// Then we also filter out all tweets that tag elon musk, e.g. there could be someone
		// just tweeting something like "Do you think xyz... @elonmusk"

		tweet.Log("tweet is starship tweet")

		// Filter out non-english tweets (except for location stream)
		if tweet.TweetSource != match.TweetSourceLocationStream && tweet.Lang != "" && tweet.Lang != "en" && tweet.Lang != "und" {
			tweet.Log("ignored because tweet is starship tweet with language %s", tweet.Lang)
			break
		}

		if p.shouldIgnoreLink(tweet) {
			tweet.Log("ignored because of link")
			break
		}

		if tweet.Tweet.PossiblySensitive {
			tweet.Log("ignored because it's possibly sensitive")
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
			case linksToLiveStream(&tweet.Tweet):
				p.retweet(&tweet.Tweet, "location + live stream", tweet.TweetSource)
			default:
				tweet.Log("location tweet ignored because it doesn't have media and is no pad announcement")
				if !p.test {
					log.Printf("[Processor] Ignoring %s because it's from the location stream and has no media", util.TweetURL(&tweet.Tweet))
				}
			}
		} else {
			switch {
			case isTagsOnly(tweet.Text()):
				// If a tweet contains *only hashtags*, we only retweet it if it has media
				tweet.Log("tweet only has tags")
				if hasMedia(&tweet.Tweet) {
					p.retweet(&tweet.Tweet, "normal matcher, only tags, but media", tweet.TweetSource)
				}
			case match.IsAtSpaceXSite(&tweet.Tweet) && linksToLiveStream(&tweet.Tweet):
				p.retweet(&tweet.Tweet, "live stream at spacex site", tweet.TweetSource)
			default:
				p.retweet(&tweet.Tweet, "normal matcher", tweet.TweetSource)
			}
		}
	}

	p.seenTweets[tweet.ID] = true

	if !tweet.Retweeted {
		p.saveNonRetweetedTweet(&tweet.Tweet)
	}
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
			util.LogError(err, "retweeting %s", util.TweetURL(tweet))
		}
		return
	}

	p.retweetedTweets[tweet.ID] = true

	if !p.test {
		// save tweet so we can reproduce why it was matched
		p.saveRetweetedTweet(tweet)

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
		util.LogError(err, "fetching tweet reply with id %d in thread", tweet.InReplyToStatusID)

		// If we have a matching tweet thread
		if err == nil && parent != nil && !match.ContainsStarshipAntiKeyword(parent.Text()) && p.thread(parent) {
			p.seenTweets[parent.ID] = true
			p.retweet(parent, "thread: matched parent", match.TweetSourceUnknown)
			didRetweet = true
		}
	}

	// A quoted tweet. Let's see if there's anything interesting
	if tweet.QuotedStatusID != 0 && tweet.QuotedStatus != nil {
		if p.thread(tweet.QuotedStatus) {
			p.retweet(tweet, "thread: quoted", match.TweetSourceUnknown)
			return true
		}
	}

	realTweet := tweet
	if tweet.RetweetedStatus != nil {
		realTweet = tweet.RetweetedStatus
	}

	// Now actually match the tweet
	if didRetweet ||
		p.matcher.StarshipTweet(match.TweetWrapper{TweetSource: match.TweetSourceUnknown, Tweet: *realTweet}) ||
		(match.ElonReplyIsStarshipRelated(tweet.Text()) && !isElonTweet(match.Wrap(tweet))) {
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
	util.LogError(err, "adding %s to list", tweet.User.ScreenName)
}
