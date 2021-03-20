module github.com/xarantolus/spacex-hop-bot

go 1.16

require (
	github.com/dghubble/go-twitter v0.0.0-20201011215211-4b180d0cc78d
	github.com/dghubble/oauth1 v0.7.0
	github.com/tdewolff/parse/v2 v2.5.14 // indirect
	github.com/xarantolus/jsonextract v1.4.3
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b
)

replace github.com/dghubble/go-twitter => ./bot/go-twitter
