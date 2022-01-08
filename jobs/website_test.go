package jobs

import (
	"bytes"
	"html/template"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	_ "embed"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/xarantolus/spacex-hop-bot/scrapers"
	"github.com/xarantolus/spacex-hop-bot/util"
)

type TestTweetingClient struct {
	tweetedTweets map[string]*int64
}

func (r *TestTweetingClient) LoadStatus(tweetID int64) (*twitter.Tweet, error) {
	panic("LoadStatus not implemented")
}

func (r *TestTweetingClient) AddListMember(listID int64, userID int64) (err error) {
	panic("AddListMember not implemented")
}

func (r *TestTweetingClient) Retweet(tweet *twitter.Tweet) error {
	panic("Retweet not implemented")
}

func (r *TestTweetingClient) Tweet(text string, inReplyToID *int64) (t *twitter.Tweet, err error) {
	r.tweetedTweets[text] = inReplyToID

	return
}

func (r *TestTweetingClient) HasTweetedWithoutAnswering(text string) bool {
	rid, ok := r.tweetedTweets[text]
	if ok && rid == nil {
		return true
	}
	return false
}

func (r TestTweetingClient) HasTweeted() bool {
	return len(r.tweetedTweets) > 0
}

//go:embed testdata/spacex_starship_page.html
var pageTemplateText string

func Test_runWebsiteScrape(t *testing.T) {
	tmpl := template.Must(template.New("").Parse(pageTemplateText))
	var pageContentWithText = func(text string) string {
		var b bytes.Buffer
		err := tmpl.Execute(&b, struct {
			Text string
		}{
			Text: text,
		})
		if err != nil {
			panic("executing template: " + err.Error())
		}

		return b.String()
	}

	var date = func(year int, month time.Month, day int) time.Time {
		return time.Date(year, month, day, 18, 24, 0, 0, util.NorthAmericaTZ)
	}

	var lastChange = scrapers.StarshipInfo{
		ShipName:       "SN15",
		NextFlightDate: date(2021, time.May, 5),
	}

	tests := []struct {
		pageContent string

		now         time.Time
		wantTweet   string
		wantNewInfo scrapers.StarshipInfo
		wantErr     bool
	}{
		{
			// This one should basically do nothing, there are no changes
			pageContent: pageContentWithText("On Wednesday, May 5, Starship serial number 15 (SN15) successfully completed SpaceX’s fifth high-altitude flight test of a Starship prototype from Starbase in Texas."),
			now:         date(2022, time.January, 8),
			wantErr:     true,
		},
		{
			pageContent: pageContentWithText("SpaceX plans an orbital test flight of S20 and B4 for Wednesday, February 16"),
			now:         date(2022, time.January, 8),
			wantErr:     false,

			wantNewInfo: scrapers.StarshipInfo{
				ShipName:       "S20",
				NextFlightDate: date(2022, time.February, 16),
			},
			wantTweet: "The #Starship website now mentions February 16, 2022 for an orbital test flight of S20 and B4",
		},
		{
			pageContent: pageContentWithText("SpaceX plans an orbital flight test with S20 and B4 for Wednesday, February 16"),
			now:         date(2022, time.January, 8),
			wantErr:     false,

			wantNewInfo: scrapers.StarshipInfo{
				ShipName:       "S20",
				NextFlightDate: date(2022, time.February, 16),
			},
			wantTweet: "The #Starship website now mentions February 16, 2022 for an orbital test flight of S20 and B4",
		},
		{
			pageContent: pageContentWithText("SpaceX plans a test flight of S20 and B4 with an orbital trajectory for Wednesday, February 16"),
			now:         date(2022, time.January, 8),
			wantErr:     false,

			wantNewInfo: scrapers.StarshipInfo{
				ShipName:       "S20",
				NextFlightDate: date(2022, time.February, 16),
			},
			wantTweet: "The #Starship website now mentions February 16, 2022 for an orbital test flight of S20 and B4",
		},
	}
	for _, tt := range tests {
		t.Run(t.Name(), func(t *testing.T) {
			twitterClient := &TestTweetingClient{}

			server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
				rw.Header().Set("Content-Type", "text/html")
				rw.Write([]byte(tt.pageContent))
			}))

			if tt.now.IsZero() {
				tt.now = time.Now()
			}

			gotNewInfo, err := runWebsiteScrape(twitterClient, nil, server.URL, lastChange, tt.now)
			if (err != nil) != tt.wantErr {
				t.Errorf("runWebsiteScrape() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotNewInfo, tt.wantNewInfo) {
				t.Errorf("runWebsiteScrape() = %v, want %v", gotNewInfo, tt.wantNewInfo)
			}

			if tt.wantTweet != "" && !twitterClient.HasTweetedWithoutAnswering(tt.wantTweet) {
				t.Errorf("runWebsiteScrape() should have tweeted %q, but didn't", tt.wantTweet)
			}
			if tt.wantTweet == "" && twitterClient.HasTweeted() {
				t.Errorf("runWebsiteScrape() should not have tweeted anything, but did")
			}
		})
	}
}
