module github.com/xarantolus/spacex-hop-bot

go 1.16

require (
	github.com/PuerkitoBio/goquery v1.8.0
	github.com/bcampbell/fuzzytime v0.0.0-20191010161914-05ea0010feac
	github.com/cenkalti/backoff v2.2.1+incompatible // indirect
	github.com/dghubble/go-twitter v0.0.0-20220816163853-8a0df96f1e6d
	github.com/dghubble/oauth1 v0.7.1
	github.com/dghubble/sling v1.4.0 // indirect
	github.com/docker/go-units v0.5.0
	github.com/icholy/replace v0.6.0
	github.com/tdewolff/parse/v2 v2.6.4 // indirect
	github.com/xarantolus/jsonextract v1.5.3
	golang.org/x/net v0.0.0-20221019024206-cb67ada4b0ad // indirect
	golang.org/x/text v0.4.0
	gopkg.in/check.v1 v1.0.0-20190902080502-41f04d3bba15 // indirect
	gopkg.in/yaml.v3 v3.0.1
	mvdan.cc/xurls/v2 v2.4.0
)

replace github.com/dghubble/go-twitter => ./bot/go-twitter
