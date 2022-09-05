package jobs

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/xarantolus/spacex-hop-bot/config"
	"github.com/xarantolus/spacex-hop-bot/consumer"
	"github.com/xarantolus/spacex-hop-bot/match"
)

type httpServer struct {
	twitter   consumer.TwitterClient
	processor *consumer.Processor
	tweetChan chan<- match.TweetWrapper
}

func httpErrWrapper(f func(w http.ResponseWriter, r *http.Request) error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := f(w, r)
		if err != nil {
			log.Printf("[Error] in %s %s: %s", r.Method, r.URL.Path, err.Error())
			if strings.Contains(r.URL.Path, "/api/") {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusInternalServerError)
				_ = json.NewEncoder(w).Encode(map[string]interface{}{
					"error_message": err.Error(),
				})
			} else {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		}
	}
}

func parseTweetID(ustr string) (tweetID int64, err error) {
	u, err := url.ParseRequestURI(ustr)
	if err != nil {
		return
	}

	if !strings.HasSuffix(u.Host, "twitter.com") {
		return 0, fmt.Errorf("URL host must be *.twitter.com, but was %s", u.Host)
	}

	pathSplit := strings.Split(u.Path, "/")
	if len(pathSplit) < 4 {
		return 0, fmt.Errorf("URL does not point to a tweet")
	}

	parsedID, err := strconv.ParseInt(pathSplit[3], 10, 64)
	if err != nil {
		return 0, fmt.Errorf("URL does not contain a tweet ID: %w", err)
	}

	return parsedID, nil
}

func (h *httpServer) submitTweet(w http.ResponseWriter, r *http.Request) (err error) {
	type incomingJSON struct {
		TweetURL string `json:"url"`
	}

	body := new(incomingJSON)
	err = json.NewDecoder(io.LimitReader(r.Body, 1<<20)).Decode(body)
	if err != nil {
		return
	}

	parsedID, err := parseTweetID(body.TweetURL)
	if err != nil {
		return
	}

	status, err := h.twitter.LoadStatus(parsedID)

	// If we cannot load the tweet, it could be that we're blocked.
	if err != nil || (status != nil && status.Retweeted) {
		log.Printf("Unretweeting tweet with id %d\n", parsedID)
		err = h.twitter.UnRetweet(parsedID)
		if err != nil {
			return
		}
		log.Printf("Unretweeted tweet with id %d\n", parsedID)
	} else if status != nil {
		// Put it into the matcher
		select {
		case h.tweetChan <- match.TweetWrapper{
			TweetSource:   match.TweetSourceUnknown,
			Tweet:         *status,
			EnableLogging: true,
		}:
			break
		case <-time.After(1 * time.Second):
			return fmt.Errorf("could not send tweet on tweetChan: timeout")
		}
	} else {
		return fmt.Errorf("could not load tweet with id %d", parsedID)
	}

	return
}

func (s *httpServer) stats(w http.ResponseWriter, r *http.Request) (err error) {
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(s.processor.Stats())
}

func RunWebServer(c config.Config, t consumer.TwitterClient, p *consumer.Processor, tweetChan chan<- match.TweetWrapper) {
	defer panic("web server stopped running, but it should never do that")

	server := &httpServer{
		twitter:   t,
		tweetChan: tweetChan,
		processor: p,
	}

	http.HandleFunc("/api/v1/tweet/submit", httpErrWrapper(server.submitTweet))
	http.HandleFunc("/api/v1/stats", httpErrWrapper(server.stats))

	port := strconv.Itoa(int(c.Server.Port))
	log.Printf("[HTTP] Server listening on port %s", port)

	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		panic("running web server: " + err.Error())
	}
}
