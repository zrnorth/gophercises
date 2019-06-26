package main

import (
	"encoding/xml"
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"

	link "github.com/zrnorth/gopher/gopher-4"
)

func get(urlStr string) []string {
	resp, err := http.Get(urlStr)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	reqURL := resp.Request.URL // In case redirected to https or something, use the response req, not the urlFlag.
	baseURL := &url.URL{
		Scheme: reqURL.Scheme,
		Host:   reqURL.Host,
	}
	base := baseURL.String()

	return filter(hrefs(resp.Body, base), withPrefix(base))
}

func hrefs(r io.Reader, base string) []string {
	links, _ := link.Parse(r)
	/*
		Rules for parsed links:
		a. /some-path										<-- need to add domain
		b. https://xxx.xxx/some-path		<-- good as-is
		c. #fragment, mailto:xxx, etc  	<-- ignore
	*/
	var ret []string

	for _, l := range links {
		switch {
		// rule A
		case strings.HasPrefix(l.Href, "/"):
			ret = append(ret, base+l.Href)
		// rule B
		case strings.HasPrefix(l.Href, "http"):
			ret = append(ret, l.Href)
		// rule C
		default: // don't use this link
		}
	}
	return ret
}

func filter(links []string, keepFn func(string) bool) []string {
	var ret []string
	for _, link := range links {
		if keepFn(link) {
			ret = append(ret, link)
		}
	}
	return ret
}

func withPrefix(pfx string) func(string) bool {
	return func(link string) bool {
		return strings.HasPrefix(link, pfx)
	}
}

type empty struct{}

func bfs(urlStr string, maxDepth int) []string {
	seen := make(map[string]empty)
	// People do this (map string -> empty struct) in golang instead of mapping to an int or bool
	// Saves a small bit of memory or something.

	var q map[string]empty     // normal queue for our bfs
	nextQ := map[string]empty{ // keep track of all next-layer links here
		urlStr: empty{},
	}

	for i := 0; i <= maxDepth; i++ {
		q, nextQ = nextQ, make(map[string]empty) // when we finish a layer, move nextQ -> q
		for url := range q {
			if _, exists := seen[url]; exists { // if we've seen this key before, skip it
				continue
			}
			seen[url] = empty{}             // mark this url as seen
			for _, link := range get(url) { // put all links from this url in the nextQ
				nextQ[link] = empty{}
			}
		}
	}

	ret := make([]string, 0, len(seen))
	for url := range seen {
		ret = append(ret, url)
	}
	return ret
}

const xmlns = "http://www.sitemaps.org/schemas/sitemap/0.9"

type loc struct {
	Value string `xml:"loc"`
}

type urlset struct {
	URLs  []loc  `xml:"url"`
	Xmlns string `xml:"xmlns,attr"`
}

func main() {
	urlFlag := flag.String("url", "http://www.github.com", "the url that you want to build a sitemap for")
	maxDepth := flag.Int("depth", 3, "the maximum number of links deep to traverse")
	flag.Parse()

	pages := bfs(*urlFlag, *maxDepth)

	toXML := urlset{
		Xmlns: xmlns,
	}
	for _, page := range pages {
		toXML.URLs = append(toXML.URLs, loc{page})
	}
	fmt.Print(xml.Header)
	enc := xml.NewEncoder(os.Stdout)
	enc.Indent("", "  ")
	if err := enc.Encode(toXML); err != nil {
		panic(err)
	}
}
