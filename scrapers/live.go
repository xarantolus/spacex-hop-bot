package scrapers

import (
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"net/url"
	"time"

	"github.com/xarantolus/jsonextract"
)

var (
	c = http.Client{
		Timeout: 10 * time.Second,
	}
	possibleUserAgents = []string{
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:86.0) Gecko/20100101 Firefox/86.0",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/42.0.2311.135 Safari/537.36 Edge/12.246",
		"Mozilla/5.0 (X11; CrOS x86_64 8172.45.0) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/51.0.2704.64 Safari/537.36",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_11_2) AppleWebKit/601.3.9 (KHTML, like Gecko) Version/9.0.2 Safari/601.3.9",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/74.0.3729.169 Safari/537.36",

		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/74.0.3729.169 Safari/537.36",
		"Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/72.0.3626.121 Safari/537.36",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/74.0.3729.157 Safari/537.36",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/60.0.3112.113 Safari/537.36",

		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/79.0.3945.88 Safari/537.36",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/78.0.3904.108 Safari/537.36",
		"Mozilla/5.0 (Windows NT 5.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/60.0.3112.90 Safari/537.36",
		"Mozilla/5.0 (Windows NT 6.2; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/60.0.3112.90 Safari/537.36",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/79.0.3945.130 Safari/537.36",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/85.0.4183.121 Safari/537.36",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/84.0.4147.105 Safari/537.36",
		"Mozilla/5.0 (Windows NT 6.3; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/60.0.3112.113 Safari/537.36",
		"Mozilla/5.0 (Windows NT 5.1; rv:36.0) Gecko/20100101 Firefox/36.0",
	}
)

// This struct only contains minimal info, there is more but I don't care about other info we can get
type LiveVideo struct {
	VideoID          string `json:"videoId"`
	Title            string `json:"title"`
	IsLive           bool   `json:"isLive"`
	ShortDescription string `json:"shortDescription"`
	IsLiveContent    bool   `json:"isLiveContent"`
}

// URL returns the youtube video URL for this live stream
func (lv *LiveVideo) URL() string {
	var u = &url.URL{
		Scheme: "https",
		Host:   "www.youtube.com",
		Path:   "watch",
	}

	var q = u.Query()
	q.Set("v", lv.VideoID)
	u.RawQuery = q.Encode()

	return u.String()
}

var ErrNotLive = errors.New("not live")

// YouTubeLive extracts a live stream from a channel live url. This kind of URL looks like the following:
//     https://www.youtube.com/channel/UCSUu1lih2RifWkKtDOJdsBA/live
//     https://www.youtube.com/spacex/live
func YouTubeLive(channelLiveURL string) (lv LiveVideo, err error) {
	req, err := http.NewRequest(http.MethodGet, channelLiveURL, nil)
	if err != nil {
		return
	}

	// Set a few headers to look like a browser
	userAgent := possibleUserAgents[rand.Intn(len(possibleUserAgents))]
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	req.Header.Set("Accept-Language", "en-US;q=0.7,en;q=0.3")

	resp, err := c.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	// Basically extract the video info and make sure it's live

	var verify = func() bool {
		return lv.VideoID != "" && (lv.IsLive || lv.IsLiveContent)
	}

	err = jsonextract.Objects(resp.Body, []jsonextract.ObjectOption{
		{
			Keys: []string{"isLive"},
			Callback: jsonextract.Unmarshal(&lv, func() bool {
				return lv.VideoID != "" && (lv.IsLive || lv.IsLiveContent)
			}),
		},
	})
	if err != nil {
		return
	}

	// If it looks ok, we return it
	if verify() {
		return lv, nil
	}

	// If it's not ok, we return an error; but we also wrap ErrNotLive
	// Include user agent information in case youtube doesn't serve JSON for certain browsers
	err = fmt.Errorf("couldn't extract livestream info from %q (useragent %q): %w", channelLiveURL, userAgent, ErrNotLive)

	return
}
