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
	tweetChan chan<- match.TweetWrapper
}

func httpErrWrapper(f func(w http.ResponseWriter, r *http.Request) error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := f(w, r)
		if err != nil {
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

func (h *httpServer) submitTweet(w http.ResponseWriter, r *http.Request) (err error) {
	type incomingJSON struct {
		TweetURL string `json:"url"`
	}

	body := new(incomingJSON)
	err = json.NewDecoder(io.LimitReader(r.Body, 1<<20)).Decode(body)
	if err != nil {
		return
	}

	u, err := url.ParseRequestURI(body.TweetURL)
	if err != nil {
		return
	}

	if !strings.HasSuffix(u.Host, "twitter.com") {
		return fmt.Errorf("URL host must be *.twitter.com, but was %s", u.Host)
	}

	pathSplit := strings.Split(u.Path, "/")
	if len(pathSplit) < 4 {
		return fmt.Errorf("URL does not point to a tweet")
	}

	parsedID, err := strconv.ParseInt(pathSplit[3], 10, 64)
	if err != nil {
		return fmt.Errorf("URL does not contain a tweet ID: %w", err)
	}

	status, err := h.twitter.LoadStatus(parsedID)
	if err != nil {
		return fmt.Errorf("loading tweet: %w", err)
	}

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

	return
}

func RunWebServer(c config.Config, t consumer.TwitterClient, tweetChan chan<- match.TweetWrapper) {
	defer panic("web server stopped running, but it should never do that")

	server := &httpServer{
		twitter:   t,
		tweetChan: tweetChan,
	}

	http.HandleFunc("/api/v1/tweet/submit", httpErrWrapper(server.submitTweet))

	port := strconv.Itoa(int(c.Server.Port))
	log.Printf("[HTTP] Server listening on port %s", port)

	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		panic("running web server: " + err.Error())
	}
}
