package jobs

import (
	"reflect"
	"testing"
	"time"

	"github.com/xarantolus/spacex-hop-bot/scrapers"
)

func Test_extractKeywords(t *testing.T) {
	type args struct {
		title       string
		description string
	}
	tests := []struct {
		args         args
		wantKeywords []string
	}{
		{args{
			title: "Starship | S20 & B4 | Orbital Flight Test", description: `As early as Wednesday, May 5, the SpaceX team will attempt an orbital flight test of Starship serial number 20 (S20) – our fifth high-altitude flight test of a Starship prototype from Starbase in Texas. S20 has vehicle improvements across structures, avionics and software, and the engines that will allow more speed and efficiency throughout production and flight: specifically, a new enhanced avionics suite, updated propellant architecture in the aft skirt, and a new Raptor engine design and configuration. 

Similar to previous high-altitude flight tests of Starship, S20 will be powered through ascent by three Raptor engines, each shutting down in sequence prior to the vehicle reaching apogee – approximately 10 km in altitude. S20 will perform a propellant transition to the internal header tanks, which hold landing propellant, before reorienting itself for reentry and a controlled aerodynamic descent.
  
The Starship prototype will descend under active aerodynamic control, accomplished by independent movement of two forward and two aft flaps on the vehicle. All four flaps are actuated by an onboard flight computer to control Starship’s attitude during flight and enable precise landing at the intended location. S20’s Raptor engines will then reignite as the vehicle attempts a landing flip maneuver immediately before touching down on the landing pad adjacent to the launch mount. 

A controlled aerodynamic descent with body flaps and vertical landing capability, combined with in-space refilling, are critical to landing Starship at destinations across the solar system where prepared surfaces or runways do not exist, and returning to Earth. This capability will enable a fully reusable transportation system designed to carry both crew and cargo on long-duration, interplanetary flights and help humanity return to the Moon, and travel to Mars and beyond.  

SuperHeavy Booster number 4 (B4)

Given the dynamic schedule of development testing, stay tuned to our social media channels for updates as we move toward SpaceX’s fifth high-altitude flight test of Starship!`,
		}, []string{"Starship", "SuperHeavy", "S20", "B4"}},
	}
	for _, tt := range tests {
		t.Run(t.Name(), func(t *testing.T) {
			if gotKeywords := extractKeywords(tt.args.title, tt.args.description); !reflect.DeepEqual(gotKeywords, tt.wantKeywords) {
				t.Errorf("extractKeywords() = %v, want %v", gotKeywords, tt.wantKeywords)
			}
		})
	}
}

func Test_matcherMatchesStreamTitle(t *testing.T) {
	// Make sure these would trigger the youtube live stream link tweet
	var titles = []string{
		"Starship Orbital Test flight",
		"S20 & B4 Test flight",
		"Booster 4 Hop",
		"Starship 20 Suborbital Test flight",
		"Starship | SN15 | High-Altitude Flight Test",
		"First Private Passenger on Lunar Starship Mission",
		"Starship Update",
		"Starship Orbital Flight Test",
		"Starship | SN5 | 150m Flight Test",
		"Starship | SN6 | 150m Flight Test",
		"Starship | SN8 | High-Altitude Flight Test",
		"Starship | SN9 | High-Altitude Flight Test",
		"Starship | SN10 | High-Altitude Flight Test",
		"Starship | SN11 | High-Altitude Flight Test",
	}

	for _, title := range titles {
		t.Run(t.Name(), func(t *testing.T) {
			matched := isStarshipStream(&scrapers.LiveVideo{
				Title: title,
			})
			if !matched {
				t.Errorf("expected video title %q to match, but didn't", title)
			}
		})
	}
}

func TestStreamTitles(t *testing.T) {
	tests := []struct {
		text string
		want bool
	}{
		{"Starlink Mission", false},
		{"Starship | SN11 | High-Altitude Flight Test", true},
		{"Starship | SN10 | High-Altitude Flight Recap", true},
		{"Starship | SN9 | High-Altitude Flight Test", true},
		{"Starship | SN8 | High-Altitude Flight Test", true},
		{"Starship SN20 & BN3: Orbital Flight Test", true},
		{"Starship | Starlink Mission", true},
	}
	for _, tt := range tests {
		t.Run(t.Name(), func(t *testing.T) {
			matched := isStarshipStream(&scrapers.LiveVideo{
				Title: tt.text,
			})
			if matched != tt.want {
				t.Errorf("expected video title %q match result in %v, but got %v", tt.text, tt.want, matched)
			}
		})
	}
}

func Test_describeLiveStream(t *testing.T) {
	tests := []struct {
		args scrapers.LiveVideo
		want string
	}{
		{
			scrapers.LiveVideo{
				VideoID: "ykdajlsdkf",
				Title:   "Starship Update",
				ShortDescription: `
SpaceX's Starship and Super Heavy launch vehicle is a fully, rapidly reusable transportation system designed to carry both crew and cargo to Earth orbit, the Moon, Mars, and anywhere else in the solar system. On Saturday, September 28 at our launch facility in Cameron County, Texas, SpaceX Chief Engineer and CEO Elon Musk will provide an update on the design and development of Starship.

You can watch the event live at approximately 8:00 p.m. CDT.
`,
				IsLive: true,
			},
			`SpaceX is now live on YouTube:

Starship Update

#Starship #SuperHeavy

https://www.youtube.com/watch?v=ykdajlsdkf`,
		},
		{
			scrapers.LiveVideo{
				VideoID: "9135813491",
				Title:   "Starship | S20 & B4 | Orbital Flight Test",
				ShortDescription: `As early as Wednesday, May 5, the SpaceX team will attempt an orbital flight test of Starship serial number 20 (S20) – our fifth high-altitude flight test of a Starship prototype from Starbase in Texas. S20 has vehicle improvements across structures, avionics and software, and the engines that will allow more speed and efficiency throughout production and flight: specifically, a new enhanced avionics suite, updated propellant architecture in the aft skirt, and a new Raptor engine design and configuration. 

Similar to previous high-altitude flight tests of Starship, S20 will be powered through ascent by three Raptor engines, each shutting down in sequence prior to the vehicle reaching apogee – approximately 10 km in altitude. S20 will perform a propellant transition to the internal header tanks, which hold landing propellant, before reorienting itself for reentry and a controlled aerodynamic descent.
  
The Starship prototype will descend under active aerodynamic control, accomplished by independent movement of two forward and two aft flaps on the vehicle. All four flaps are actuated by an onboard flight computer to control Starship’s attitude during flight and enable precise landing at the intended location. S20’s Raptor engines will then reignite as the vehicle attempts a landing flip maneuver immediately before touching down on the landing pad adjacent to the launch mount. 

A controlled aerodynamic descent with body flaps and vertical landing capability, combined with in-space refilling, are critical to landing Starship at destinations across the solar system where prepared surfaces or runways do not exist, and returning to Earth. This capability will enable a fully reusable transportation system designed to carry both crew and cargo on long-duration, interplanetary flights and help humanity return to the Moon, and travel to Mars and beyond.  

SuperHeavy Booster number 4 (B4)

Given the dynamic schedule of development testing, stay tuned to our social media channels for updates as we move toward SpaceX’s fifth high-altitude flight test of Starship!`,
				IsLive: true,
			},
			`SpaceX is now live on YouTube:

Starship | S20 & B4 | Orbital Flight Test

#Starship #SuperHeavy #S20 #B4

https://www.youtube.com/watch?v=9135813491`,
		},
		{
			scrapers.LiveVideo{
				VideoID: "9135813491",
				Title:   "Starship | S20 & B4 | Orbital Flight Test",
				ShortDescription: `As early as Wednesday, May 5, the SpaceX team will attempt an orbital flight test of Starship serial number 20 (S20) – our fifth high-altitude flight test of a Starship prototype from Starbase in Texas. S20 has vehicle improvements across structures, avionics and software, and the engines that will allow more speed and efficiency throughout production and flight: specifically, a new enhanced avionics suite, updated propellant architecture in the aft skirt, and a new Raptor engine design and configuration. 

Similar to previous high-altitude flight tests of Starship, S20 will be powered through ascent by three Raptor engines, each shutting down in sequence prior to the vehicle reaching apogee – approximately 10 km in altitude. S20 will perform a propellant transition to the internal header tanks, which hold landing propellant, before reorienting itself for reentry and a controlled aerodynamic descent.
  
The Starship prototype will descend under active aerodynamic control, accomplished by independent movement of two forward and two aft flaps on the vehicle. All four flaps are actuated by an onboard flight computer to control Starship’s attitude during flight and enable precise landing at the intended location. S20’s Raptor engines will then reignite as the vehicle attempts a landing flip maneuver immediately before touching down on the landing pad adjacent to the launch mount. 

A controlled aerodynamic descent with body flaps and vertical landing capability, combined with in-space refilling, are critical to landing Starship at destinations across the solar system where prepared surfaces or runways do not exist, and returning to Earth. This capability will enable a fully reusable transportation system designed to carry both crew and cargo on long-duration, interplanetary flights and help humanity return to the Moon, and travel to Mars and beyond.  

SuperHeavy Booster number 4 (B4)

Given the dynamic schedule of development testing, stay tuned to our social media channels for updates as we move toward SpaceX’s fifth high-altitude flight test of Starship!`,
				IsUpcoming: true,
			},
			`SpaceX live stream starts soon:

Starship | S20 & B4 | Orbital Flight Test

#Starship #SuperHeavy #S20 #B4

https://www.youtube.com/watch?v=9135813491`,
		},
		{
			scrapers.LiveVideo{
				VideoID: "9135813491",
				Title:   "Starship | S20 & B4 | Orbital Flight Test",
				ShortDescription: `As early as Wednesday, May 5, the SpaceX team will attempt an orbital flight test of Starship serial number 20 (S20) – our fifth high-altitude flight test of a Starship prototype from Starbase in Texas. S20 has vehicle improvements across structures, avionics and software, and the engines that will allow more speed and efficiency throughout production and flight: specifically, a new enhanced avionics suite, updated propellant architecture in the aft skirt, and a new Raptor engine design and configuration. 

Similar to previous high-altitude flight tests of Starship, S20 will be powered through ascent by three Raptor engines, each shutting down in sequence prior to the vehicle reaching apogee – approximately 10 km in altitude. S20 will perform a propellant transition to the internal header tanks, which hold landing propellant, before reorienting itself for reentry and a controlled aerodynamic descent.
  
The Starship prototype will descend under active aerodynamic control, accomplished by independent movement of two forward and two aft flaps on the vehicle. All four flaps are actuated by an onboard flight computer to control Starship’s attitude during flight and enable precise landing at the intended location. S20’s Raptor engines will then reignite as the vehicle attempts a landing flip maneuver immediately before touching down on the landing pad adjacent to the launch mount. 

A controlled aerodynamic descent with body flaps and vertical landing capability, combined with in-space refilling, are critical to landing Starship at destinations across the solar system where prepared surfaces or runways do not exist, and returning to Earth. This capability will enable a fully reusable transportation system designed to carry both crew and cargo on long-duration, interplanetary flights and help humanity return to the Moon, and travel to Mars and beyond.  

SuperHeavy Booster number 4 (B4)

Given the dynamic schedule of development testing, stay tuned to our social media channels for updates as we move toward SpaceX’s fifth high-altitude flight test of Starship!`,
				IsUpcoming: true,
				UpcomingInfo: scrapers.LiveBroadcastDetails{
					StartTimestamp: time.Now().Add(1*time.Minute + 10*time.Second),
				},
			},
			`SpaceX live stream starts in about a minute:

Starship | S20 & B4 | Orbital Flight Test

#Starship #SuperHeavy #S20 #B4

https://www.youtube.com/watch?v=9135813491`,
		},
		{
			scrapers.LiveVideo{
				VideoID: "9135813491",
				Title:   "Starship | S20 & B4 | Orbital Flight Test",
				ShortDescription: `As early as Wednesday, May 5, the SpaceX team will attempt an orbital flight test of Starship serial number 20 (S20) – our fifth high-altitude flight test of a Starship prototype from Starbase in Texas. S20 has vehicle improvements across structures, avionics and software, and the engines that will allow more speed and efficiency throughout production and flight: specifically, a new enhanced avionics suite, updated propellant architecture in the aft skirt, and a new Raptor engine design and configuration. 

Similar to previous high-altitude flight tests of Starship, S20 will be powered through ascent by three Raptor engines, each shutting down in sequence prior to the vehicle reaching apogee – approximately 10 km in altitude. S20 will perform a propellant transition to the internal header tanks, which hold landing propellant, before reorienting itself for reentry and a controlled aerodynamic descent.
  
The Starship prototype will descend under active aerodynamic control, accomplished by independent movement of two forward and two aft flaps on the vehicle. All four flaps are actuated by an onboard flight computer to control Starship’s attitude during flight and enable precise landing at the intended location. S20’s Raptor engines will then reignite as the vehicle attempts a landing flip maneuver immediately before touching down on the landing pad adjacent to the launch mount. 

A controlled aerodynamic descent with body flaps and vertical landing capability, combined with in-space refilling, are critical to landing Starship at destinations across the solar system where prepared surfaces or runways do not exist, and returning to Earth. This capability will enable a fully reusable transportation system designed to carry both crew and cargo on long-duration, interplanetary flights and help humanity return to the Moon, and travel to Mars and beyond.  

SuperHeavy Booster number 4 (B4)

Given the dynamic schedule of development testing, stay tuned to our social media channels for updates as we move toward SpaceX’s fifth high-altitude flight test of Starship!`,
				IsUpcoming: true,
				UpcomingInfo: scrapers.LiveBroadcastDetails{
					StartTimestamp: time.Now().Add(3*time.Hour + 10*time.Minute),
				},
			},
			`SpaceX live stream starts in 3 hours:

Starship | S20 & B4 | Orbital Flight Test

#Starship #SuperHeavy #S20 #B4

https://www.youtube.com/watch?v=9135813491`,
		},
		{
			scrapers.LiveVideo{
				VideoID: "9135813491",
				Title:   "Starship | S20 & B4 | Orbital Flight Test",
				ShortDescription: `As early as Wednesday, May 5, the SpaceX team will attempt an orbital flight test of Starship serial number 20 (S20) – our fifth high-altitude flight test of a Starship prototype from Starbase in Texas. S20 has vehicle improvements across structures, avionics and software, and the engines that will allow more speed and efficiency throughout production and flight: specifically, a new enhanced avionics suite, updated propellant architecture in the aft skirt, and a new Raptor engine design and configuration. 

Similar to previous high-altitude flight tests of Starship, S20 will be powered through ascent by three Raptor engines, each shutting down in sequence prior to the vehicle reaching apogee – approximately 10 km in altitude. S20 will perform a propellant transition to the internal header tanks, which hold landing propellant, before reorienting itself for reentry and a controlled aerodynamic descent.
  
The Starship prototype will descend under active aerodynamic control, accomplished by independent movement of two forward and two aft flaps on the vehicle. All four flaps are actuated by an onboard flight computer to control Starship’s attitude during flight and enable precise landing at the intended location. S20’s Raptor engines will then reignite as the vehicle attempts a landing flip maneuver immediately before touching down on the landing pad adjacent to the launch mount. 

A controlled aerodynamic descent with body flaps and vertical landing capability, combined with in-space refilling, are critical to landing Starship at destinations across the solar system where prepared surfaces or runways do not exist, and returning to Earth. This capability will enable a fully reusable transportation system designed to carry both crew and cargo on long-duration, interplanetary flights and help humanity return to the Moon, and travel to Mars and beyond.  

SuperHeavy Booster number 4 (B4)

Given the dynamic schedule of development testing, stay tuned to our social media channels for updates as we move toward SpaceX’s fifth high-altitude flight test of Starship!`,
				IsUpcoming: true,
				UpcomingInfo: scrapers.LiveBroadcastDetails{
					StartTimestamp: time.Now().Add(3*time.Hour + 35*time.Minute),
				},
			},
			`SpaceX live stream starts in 4 hours:

Starship | S20 & B4 | Orbital Flight Test

#Starship #SuperHeavy #S20 #B4

https://www.youtube.com/watch?v=9135813491`,
		},
	}
	for _, tt := range tests {
		t.Run(t.Name(), func(t *testing.T) {
			if got := describeLiveStream(&tt.args); got != tt.want {
				t.Errorf("describeLiveStream() = \n%v, want \n%v", got, tt.want)
			}

			matched := isStarshipStream(&tt.args)
			if !matched {
				t.Errorf("describeLiveStream(): expected video with title %q and description %q to match, but didn't", tt.args.Title, tt.args.ShortDescription)
			}
		})
	}
}
