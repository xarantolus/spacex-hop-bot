module github.com/xarantolus/spacex-hop-bot

go 1.16

require (
	github.com/PuerkitoBio/goquery v1.7.0
	github.com/andybalholm/cascadia v1.2.0 // indirect
	github.com/bcampbell/fuzzytime v0.0.0-20191010161914-05ea0010feac
	github.com/cenkalti/backoff v2.2.1+incompatible // indirect
	github.com/dghubble/go-twitter v0.0.0-20210609183100-2fdbf421508e
	github.com/dghubble/oauth1 v0.7.0
	github.com/docker/go-units v0.4.0
	github.com/esimov/pigo v1.4.4
	github.com/google/go-querystring v1.1.0 // indirect
	github.com/icholy/replace v0.4.0
	github.com/tdewolff/parse/v2 v2.5.18 // indirect
	github.com/xarantolus/jsonextract v1.5.3
	golang.org/x/image v0.0.0-20210628002857-a66eb6448b8d
	golang.org/x/net v0.0.0-20210614182718-04defd469f4e // indirect
	golang.org/x/text v0.3.6
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b
	mvdan.cc/xurls/v2 v2.2.0
)

replace github.com/dghubble/go-twitter => ./bot/go-twitter
