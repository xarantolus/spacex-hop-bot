module github.com/xarantolus/spacex-hop-bot

go 1.16

require (
	github.com/PuerkitoBio/goquery v1.6.1
	github.com/andybalholm/cascadia v1.2.0 // indirect
	github.com/bcampbell/fuzzytime v0.0.0-20191010161914-05ea0010feac
	github.com/cenkalti/backoff v2.2.1+incompatible // indirect
	github.com/dghubble/go-twitter v0.0.0-20201011215211-4b180d0cc78d
	github.com/dghubble/oauth1 v0.7.0
	github.com/docker/go-units v0.4.0
	github.com/google/go-querystring v1.1.0 // indirect
	github.com/xarantolus/jsonextract v1.5.1
	golang.org/x/net v0.0.0-20210330075724-22f4162a9025 // indirect
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b
)

replace github.com/dghubble/go-twitter => ./bot/go-twitter
