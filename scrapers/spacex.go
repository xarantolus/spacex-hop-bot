package scrapers

import (
	"bufio"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/xarantolus/spacex-hop-bot/util"
)

const (
	StarshipURL = "https://www.spacex.com/vehicles/starship/"
)

var ShipNameRegex = regexp.MustCompile(`\b((?:SN|BN|S|B|Booster|Star[Ss]hip)-?\s*\d+)\b`)

type StarshipInfo struct {
	ShipName       string
	NextFlightDate time.Time

	LiveStreamID string
}

func (s *StarshipInfo) Equals(b StarshipInfo) bool {
	return s.ShipName == b.ShipName && s.NextFlightDate.Equal(b.NextFlightDate)
}

var ErrNoInfo = errors.New("no info")

// SpaceXStarship returns info about the starship page (ship name & first mentioned date)
func SpaceXStarship() (s StarshipInfo, err error) {
	resp, err := c.Get(StarshipURL)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	// Extract the HTML body text
	doc, err := goquery.NewDocumentFromReader(bufio.NewReader(resp.Body))
	if err != nil {
		return
	}

	var (
		date     time.Time
		shipName string
	)

	// Let's check the text on the website.
	// We first look for .text-column elements,
	// but in case something is changed on the site, we also look at all divs
	doc.Find(".text-column").Add("div").EachWithBreak(func(i int, s *goquery.Selection) bool {
		content := s.Text()

		// Find the first interesting text
		if shipName == "" {
			shipName = ShipNameRegex.FindString(content)
		}

		if date.IsZero() {
			// Try to extract a date
			etime, ok := util.ExtractDate(content)
			if ok && time.Now().In(util.NorthAmericaTZ).Sub(etime) > 0 {
				date = etime
			}
		}

		// if we have both, we break (by returning true)
		return !date.IsZero() && shipName != ""
	})

	// Find the first element with a data-video attribute. This is where SpaceX puts replays and other videos
	vidButton := doc.Find("[data-video]").First()

	// Now check if we can get the video id of a video that is not a replay
	var liveID string
	if vidButton.Length() > 0 {
		vt := strings.TrimSpace(vidButton.Text())
		if !strings.EqualFold(vt, "REPLAY") && !strings.EqualFold(vt, "REWATCH") && !strings.EqualFold(vt, "RE-WATCH") {
			liveID, _ = vidButton.Attr("data-video")
		}
	}

	if date.IsZero() {
		err = fmt.Errorf("couldn't extract date info: %w", ErrNoInfo)
		return
	}
	if shipName == "" {
		err = fmt.Errorf("couldn't extract ship name: %w", ErrNoInfo)
		return
	}

	return StarshipInfo{
		ShipName:       shipName,
		NextFlightDate: date,
		LiveStreamID:   liveID,
	}, nil
}
