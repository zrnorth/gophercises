package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	link "github.com/zrnorth/gopher/gopher-4"
)

func get(urlStr string) []string {
	fmt.Println("GETting from " + urlStr)

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

func main() {
	urlFlag := flag.String("url", "http://www.github.com", "the url that you want to build a sitemap for")
	flag.Parse()

	pages := get(*urlFlag)

	for _, page := range pages {
		fmt.Println(page)
	}
	// 3. filter out any links that have a different domain
	// 4. BFS to next layer of links
	// 5. print the formatted xml
}
