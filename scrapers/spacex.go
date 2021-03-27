package scrapers

import (
	"errors"
	"fmt"
	"regexp"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/xarantolus/spacex-hop-bot/util"
)

const (
	StarshipURL = "https://www.spacex.com/vehicles/starship/"
)

var shipNameRegex = regexp.MustCompile(`((?:SN|BN)\s*\d+)`)

type StarshipInfo struct {
	ShipName       string
	NextFlightDate time.Time
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
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return
	}

	var (
		date     time.Time
		shipName string
	)

	// Let's check the text on the website
	doc.Find("div").EachWithBreak(func(i int, s *goquery.Selection) bool {
		content := s.Text()

		// Find the first interesting text
		if shipName == "" {
			var eShipName = shipNameRegex.FindString(content)
			if eShipName != "" {
				shipName = eShipName
			}
		}

		// Try to extract a date
		etime, ok := util.ExtractDate(content)
		if ok {
			date = etime

			// This is like "break" in a for loop
			return false
		}

		// Continues matching
		return true
	})

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
	}, nil
}
