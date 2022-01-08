package util

import (
	"time"

	"github.com/bcampbell/fuzzytime"
)

// This is the time zone for Texas
// For other operations SpaceX also uses EDT, but doesn't seem to be the case with Starship stuff
var NorthAmericaTZ = time.FixedZone("CDT", -5*60*60)

// ExtractDate extracts human-readable dates from text
func ExtractDate(text string, now time.Time) (date time.Time, ok bool) {
	d, _, err := fuzzytime.ExtractDate(text)
	if err != nil || d.Empty() {
		return
	}

	// Now merge with the current date
	fuzzyNow := fuzzytime.NewDate(now.Year(), int(now.Month()), now.Day())
	fuzzyNow.Merge(&d)

	return time.Date(fuzzyNow.Year(), time.Month(fuzzyNow.Month()), fuzzyNow.Day(), 0, 0, 0, 0, NorthAmericaTZ), true
}
