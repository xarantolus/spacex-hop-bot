package util

import (
	"bufio"
	"math/rand"
	"net/http"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/icholy/replace"
	"golang.org/x/text/transform"
	"mvdan.cc/xurls/v2"
)

var (
	urlRegex = xurls.Strict()

	noScriptEnter = replace.String("<noscript>", "")
	noScriptLeave = replace.String("</noscript>", "")

	noScriptReplacer = transform.Chain(noScriptEnter, noScriptLeave)
)

// FindCanonicalURL returns the canonical URL of an article if possible,
// else the original input is returned
func FindCanonicalURL(url string) (out string) {
	out = url

	client := http.Client{
		Timeout: 30 * time.Second,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			// If we were redirected, we now get a different URL
			out = req.URL.String()
			return nil
		},
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return
	}

	req.Header.Set("User-Agent", GetUserAgent())
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	req.Header.Set("Accept-Language", "en-US;q=0.7,en;q=0.3")

	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusBadRequest {
		return
	}

	// t.co doesn't use real redirects. And the <meta http-equiv="refresh" tag is inside a <noscript>
	// tag which the Go html parser doesn't parse (because it's in script mode and thus parses that node as text...)
	// So here is the stupid but working solution: replacing any "<noscript>" "</noscript>" text
	bodyReader := transform.NewReader(bufio.NewReader(resp.Body), noScriptReplacer)

	// Still try to find the tag
	doc, err := goquery.NewDocumentFromReader(bodyReader)
	if err != nil {
		return
	}

	canon := doc.Find("[rel=canonical]").First()
	if canon.Length() != 0 {
		return canon.AttrOr("href", url)
	}

	equiv := doc.Find("[http-equiv]").First()
	if equiv.Length() != 0 {
		u := urlRegex.FindString(equiv.AttrOr("content", ""))
		if u != "" {
			return u
		}
	}

	return
}

var possibleUserAgents = []string{
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:87.0) Gecko/20100101 Firefox/87.0",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:86.0) Gecko/20100101 Firefox/86.0",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/42.0.2311.135 Safari/537.36 Edge/12.246",
	"Mozilla/5.0 (X11; CrOS x86_64 8172.45.0) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/51.0.2704.64 Safari/537.36",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_11_2) AppleWebKit/601.3.9 (KHTML, like Gecko) Version/9.0.2 Safari/601.3.9",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/74.0.3729.169 Safari/537.36",

	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/74.0.3729.169 Safari/537.36",
	"Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/72.0.3626.121 Safari/537.36",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/74.0.3729.157 Safari/537.36",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/60.0.3112.113 Safari/537.36",

	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/79.0.3945.88 Safari/537.36",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/78.0.3904.108 Safari/537.36",
	"Mozilla/5.0 (Windows NT 5.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/60.0.3112.90 Safari/537.36",
	"Mozilla/5.0 (Windows NT 6.2; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/60.0.3112.90 Safari/537.36",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/79.0.3945.130 Safari/537.36",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/85.0.4183.121 Safari/537.36",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/84.0.4147.105 Safari/537.36",
	"Mozilla/5.0 (Windows NT 6.3; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/60.0.3112.113 Safari/537.36",
	"Mozilla/5.0 (Windows NT 5.1; rv:36.0) Gecko/20100101 Firefox/36.0",
}

func GetUserAgent() string {
	return possibleUserAgents[rand.Intn(len(possibleUserAgents))]
}
