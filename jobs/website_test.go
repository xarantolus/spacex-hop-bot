package jobs

import (
	"bytes"
	_ "embed"
	"html/template"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/xarantolus/spacex-hop-bot/scrapers"
	"github.com/xarantolus/spacex-hop-bot/util"
)

type TestTweetingClient struct {
	tweetedTweetText  string
	tweetedTweetReply *int64
}

func (r *TestTweetingClient) UnRetweet(tweetID int64) error {
	panic("UnRetweet not implemented")
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
	if r.HasTweeted() {
		panic("test tweeted more than once")
	}
	r.tweetedTweetText = text
	r.tweetedTweetReply = inReplyToID

	return &twitter.Tweet{
		ID:         13589,
		SimpleText: text,
	}, nil
}

func (r *TestTweetingClient) HasTweetedWithoutAnswering(text string) bool {
	return r.tweetedTweetText == text && r.tweetedTweetReply == nil
}

func (r TestTweetingClient) HasTweeted() bool {
	return r.tweetedTweetText != ""
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
		return time.Date(year, month, day, 0, 0, 0, 0, util.NorthAmericaTZ)
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
			// This one should basically do nothing, there are no changes.
			// The problem we have here is that the fuzzy date matcher assumes that May 5 is in the current year, while the flight was in 2021
			// I have decided to handle this case by just editing the saved JSON file once a year
			pageContent: pageContentWithText("On Wednesday, May 5, Starship serial number 15 (SN15) successfully completed SpaceX’s fifth high-altitude flight test of a Starship prototype from Starbase in Texas."),
			now:         date(2022, time.January, 8),
			wantNewInfo: scrapers.StarshipInfo{
				ShipName: "SN15",
				// Yes, this is wrong, but for the test it's right
				NextFlightDate: date(2022, time.May, 5),
			},

			wantTweet: "The SpaceX #Starship website now mentions May 5 for #SN15\n#WenHop\n" + scrapers.StarshipURL,
		},
		{
			pageContent: pageContentWithText("On Wednesday, February 13, Starship serial number 20 (S20) will do another high-altitude flight test of a Starship prototype from Starbase in Texas."),
			now:         date(2022, time.January, 8),
			wantNewInfo: scrapers.StarshipInfo{
				ShipName:       "S20",
				NextFlightDate: date(2022, time.February, 13),
			},
			wantTweet: "The SpaceX #Starship website now mentions February 13 for #S20\n#WenHop\n" + scrapers.StarshipURL,
		},
		{
			pageContent: pageContentWithText("SpaceX plans an orbital test flight of S20 and B4 for Wednesday, February 16"),
			now:         date(2022, time.January, 8),
			wantErr:     false,

			wantNewInfo: scrapers.StarshipInfo{
				ShipName:       "S20",
				NextFlightDate: date(2022, time.February, 16),
				Orbital:        true,
			},
			wantTweet: "The SpaceX #Starship website now mentions February 16 for an orbital flight of #S20\n#WenHop\n" + scrapers.StarshipURL,
		},
		{
			pageContent: pageContentWithText("SpaceX plans an orbital flight test with S20 and B4 for Wednesday, February 16"),
			now:         date(2022, time.January, 8),
			wantErr:     false,

			wantNewInfo: scrapers.StarshipInfo{
				ShipName:       "S20",
				NextFlightDate: date(2022, time.February, 16),
				Orbital:        true,
			},
			wantTweet: "The SpaceX #Starship website now mentions February 16 for an orbital flight of #S20\n#WenHop\n" + scrapers.StarshipURL,
		},
		{
			pageContent: pageContentWithText("SpaceX plans a test flight of S20 and B4 with an orbital trajectory for Wednesday, February 16"),
			now:         date(2022, time.January, 8),
			wantErr:     false,

			wantNewInfo: scrapers.StarshipInfo{
				ShipName:       "S20",
				NextFlightDate: date(2022, time.February, 16),
				Orbital:        true,
			},
			wantTweet: "The SpaceX #Starship website now mentions February 16 for an orbital flight of #S20\n#WenHop\n" + scrapers.StarshipURL,
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
				rt := twitterClient.tweetedTweetText
				if rt == "" {
					t.Errorf("runWebsiteScrape() should have tweeted %q, but didn't tweet anything", tt.wantTweet)
				} else {
					t.Errorf("runWebsiteScrape() should have tweeted %q, but tweeted %q", tt.wantTweet, rt)

				}
			}
			if tt.wantTweet == "" && twitterClient.HasTweeted() {
				t.Errorf("runWebsiteScrape() should not have tweeted anything, but did (%q)", twitterClient.tweetedTweetText)
			}
		})
	}
}
