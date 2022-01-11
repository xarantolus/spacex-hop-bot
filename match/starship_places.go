package match

import "github.com/dghubble/go-twitter/twitter"

const (
	// TODO: find IDs for "Mesa del Gavilan", Stargate and generally places around/between the site.
	// The data seems to come from foursquare, but the IDs are *not* the same on both services

	// https://twitter.com/places/1380f3b60f972001
	StarbasePlaceID = "1380f3b60f972001"

	// https://twitter.com/places/124cb6de55957000
	SpaceXLaunchSiteID = "124cb6de55957000"

	// https://twitter.com/places/124bed061054f000
	SpaceXBuildSiteID = "124bed061054f000"

	// https://twitter.com/places/07d9f642af482000
	SpaceXMcGregorPlaceID = "07d9f642af482000"

	// https://twitter.com/places/07d9f0b85ac83003
	BocaChicaPlaceID = "07d9f0b85ac83003"

	// Other places around the area:
	// "Isla Blanca Park": https://twitter.com/places/11dca9a728950001
	// "South Padre Island, TX": https://twitter.com/places/1d1f665883989434

	// "Cape Canaveral, FL": https://twitter.com/places/1739d72c18edbb1e
)

func IsAtSpaceXSite(tweet *twitter.Tweet) bool {
	return tweet.Place != nil && (tweet.Place.ID == StarbasePlaceID ||
		tweet.Place.ID == SpaceXLaunchSiteID || tweet.Place.ID == SpaceXBuildSiteID ||
		tweet.Place.ID == SpaceXMcGregorPlaceID || tweet.Place.ID == BocaChicaPlaceID)
}
