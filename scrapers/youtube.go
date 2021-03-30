package scrapers

import (
	"encoding/json"
	"errors"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/xarantolus/jsonextract"
)

var (
	c = http.Client{
		Timeout: 10 * time.Second,
	}
	possibleUserAgents = []string{
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:87.0) Gecko/20100101 Firefox/87.0",
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
	IsUpcoming       bool   `json:"isUpcoming"`

	upcomingInfo liveBroadcastDetails
	unixInfo     liveBroadcastUnixInfo
}

type liveBroadcastDetails struct {
	StartTimestamp time.Time `json:"startTimestamp"`
}

type liveBroadcastUnixInfo struct {
	ScheduledStartTime UnixTime `json:"scheduledStartTime"`
}

type UnixTime time.Time

func (t *UnixTime) UnmarshalJSON(b []byte) (err error) {
	var s string
	err = json.Unmarshal(b, &s)
	if err != nil {
		return
	}

	i, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return
	}

	*t = UnixTime(time.Unix(i, 0))

	return nil
}

func (l *LiveVideo) TimeUntil() (t time.Time, d time.Duration, ok bool) {
	// Check if we got any time info
	t = l.upcomingInfo.StartTimestamp
	if t.IsZero() {
		t = time.Time(l.unixInfo.ScheduledStartTime)

		if t.IsZero() {
			return
		}
	}

	d = t.UTC().Sub(time.Now().UTC())
	ok = d > 0
	return
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

var ErrNoVideo = errors.New("not live")

// YouTubeLive extracts a live stream from a channel live url. This kind of URL looks like the following:
//     https://www.youtube.com/channel/UCSUu1lih2RifWkKtDOJdsBA/live
//     https://www.youtube.com/spacex/live
// It also extract streams that are upcoming
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

	// Basically extract the video info and make sure it's live/upcoming
	// We also extract extra info when the livestream will go live
	err = jsonextract.Objects(resp.Body, []jsonextract.ObjectOption{
		{
			Keys: []string{"videoId"},
			Callback: jsonextract.Unmarshal(&lv, func() bool {
				return lv.VideoID != "" && (lv.IsLive || lv.IsUpcoming)
			}),
			Required: true,
		},

		// There are two ways of getting the Upcoming time of a livestream, so we need to handle both
		{
			Keys: []string{"startTimestamp"},
			Callback: jsonextract.Unmarshal(&lv.upcomingInfo, func() bool {
				return !lv.upcomingInfo.StartTimestamp.IsZero()
			}),
		},
		{
			Keys: []string{"scheduledStartTime"},
			Callback: jsonextract.Unmarshal(&lv.unixInfo, func() bool {
				return !time.Time(lv.unixInfo.ScheduledStartTime).IsZero()
			}),
		},
	})

	if errors.Is(err, jsonextract.ErrCallbackNeverCalled) {
		err = ErrNoVideo
	}

	return
}
