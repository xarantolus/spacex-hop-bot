package scrapers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"
	"time"

	"github.com/xarantolus/jsonextract"
	"github.com/xarantolus/spacex-hop-bot/util"
)

var (
	c = http.Client{
		Timeout: 30 * time.Second,
	}
)

// This struct only contains minimal info, there is more but I don't care about other info we can get
type LiveVideo struct {
	VideoID string `json:"videoId"`
	Title   string `json:"title"`
	IsLive  bool   `json:"isLive"`

	ChannelID string `json:"channelId"`

	ShortDescription string `json:"shortDescription"`
	IsUpcoming       bool   `json:"isUpcoming"`

	UpcomingInfo LiveBroadcastDetails
	unixInfo     liveBroadcastUnixInfo
}

type LiveBroadcastDetails struct {
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

func (lv *LiveVideo) TimeUntil() (t time.Time, d time.Duration, ok bool) {
	// Check if we got any time info
	t = lv.UpcomingInfo.StartTimestamp
	if t.IsZero() {
		t = time.Time(lv.unixInfo.ScheduledStartTime)

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
	req.Header.Set("User-Agent", util.GetUserAgent())
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	req.Header.Set("Accept-Language", "en-US;q=0.7,en;q=0.3")

	// This cookie is set after going to that consent page (e.g. when visiting a video in private mode)
	// Sometimes a redirect to that page causes the scraper to not work, so this is an attempt to go around that
	req.Header.Set("Cookie", "CONSENT=YES+cb.20210328-17-p0.en+FX+419; PREF=tz=UTC")

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
			Callback: jsonextract.Unmarshal(&lv.UpcomingInfo, func() bool {
				return !lv.UpcomingInfo.StartTimestamp.IsZero()
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
	} else if err != nil {
		rdetails, derr := httputil.DumpRequest(req, false)
		if derr == nil {
			log.Printf("[YouTube Scraper] Failed with %s\n%s", err.Error(), string(rdetails))
		} else {
			log.Println("[YouTube Scraper] Cannot get request details")
		}
	}

	return
}
