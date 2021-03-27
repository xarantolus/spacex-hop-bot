package util

import (
	"time"

	"github.com/bcampbell/fuzzytime"
)

// This is the time zone SpaceX often uses
var northAmericaTZ = time.FixedZone("EDT", -4*60*60)

// ExtractDate extracts human-readable dates from text
func ExtractDate(text string) (date time.Time, ok bool) {
	d, _, err := fuzzytime.ExtractDate(text)
	if err != nil || d.Empty() {
		return
	}

	// Now merge with the current date
	now := time.Now()
	fuzzyNow := fuzzytime.NewDate(now.Year(), int(now.Month()), now.Day())
	fuzzyNow.Merge(&d)

	return time.Date(fuzzyNow.Year(), time.Month(fuzzyNow.Month()), fuzzyNow.Day(), 0, 0, 0, 0, northAmericaTZ), true
}
