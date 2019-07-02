package main

import (
	"errors"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/zrnorth/gopher/gopher-13/hn"
)

func main() {
	// parse flags
	var port, numStories int
	flag.IntVar(&port, "port", 3000, "the port to start the web server on")
	flag.IntVar(&numStories, "num_stories", 30, "the number of top stories to display")
	flag.Parse()

	tpl := template.Must(template.ParseFiles("./index.gohtml"))

	http.HandleFunc("/", handler(numStories, tpl))

	// Start the server
	fmt.Printf("Listening on port %d...\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}

func handler(numStories int, tpl *template.Template) http.HandlerFunc {
	// Create a cache
	cacheDuration := 6 * time.Second
	sc := storyCache{
		numStories: numStories,
		duration:   cacheDuration,
	}

	// Create a ticker that updates the cache
	go func() {
		ticker := time.NewTicker(3 * time.Second)
		for {
			temp := storyCache{
				numStories: numStories,
				duration:   cacheDuration,
			}
			temp.stories()
			sc.mutex.Lock()
			sc.cache = temp.cache
			sc.expiration = temp.expiration
			sc.mutex.Unlock()

			<-ticker.C
		}
	}()

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		stories, err := sc.stories()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		data := templateData{
			Stories: stories,
			Time:    time.Now().Sub(start),
		}
		err = tpl.Execute(w, data)
		if err != nil {
			http.Error(w, "Failed to process the template", http.StatusInternalServerError)
			return
		}
	})
}

type storyCache struct {
	numStories int
	cache      []item
	expiration time.Time
	duration   time.Duration
	mutex      sync.Mutex
}

func (sc *storyCache) stories() ([]item, error) {
	sc.mutex.Lock()
	defer sc.mutex.Unlock()

	if sc.expiration.Sub(time.Now()) > 0 {
		return sc.cache, nil
	}

	stories, err := getTopStories(sc.numStories)
	if err != nil {
		return nil, err
	}
	sc.expiration = time.Now().Add(sc.duration)

	sc.cache = stories
	return sc.cache, nil
}

func getTopStories(numStories int) ([]item, error) {
	var client hn.Client
	ids, err := client.TopItems()
	if err != nil {
		return nil, errors.New("Failed to load top stories")
	}

	// Try to get the number of stories we need. If we come up short bc one or more
	// are text posts or AskHN, try again until we have exactly numStories.
	var stories []item
	at := 0
	for len(stories) < numStories {
		need := (numStories - len(stories)) * 5 / 4 // add a little bit extra to reduce loops
		stories = append(stories, getStories(ids[at:at+need])...)
		at += need
	}
	return stories[:numStories], nil
}

// Creates a bunch of parallel goroutines, each one getting a single item by ID
// Blocks until they all finish and then returns the consolidated stories
func getStories(ids []int) []item {
	type result struct {
		idx  int
		item item
		err  error
	}
	resultCh := make(chan result)

	for i := 0; i < len(ids); i++ {
		go func(idx, id int) {
			var client hn.Client
			hnItem, err := client.GetItem(id)
			if err != nil {
				resultCh <- result{idx: idx, err: err}
			}
			resultCh <- result{idx: idx, item: parseHNItem(hnItem)}
		}(i, ids[i])
	}

	// Consolidate the results here and sort them
	var results []result
	for i := 0; i < len(ids); i++ {
		results = append(results, <-resultCh)
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].idx < results[j].idx
	})

	// Filter for just stories (remove text posts, errors, etc)
	var stories []item
	for _, res := range results {
		if res.err != nil {
			continue
		}
		if isStoryLink(res.item) {
			stories = append(stories, res.item)
		}
	}
	return stories
}

func isStoryLink(item item) bool {
	return item.Type == "story" && item.URL != ""
}

func parseHNItem(hnItem hn.Item) item {
	ret := item{Item: hnItem}
	url, err := url.Parse(ret.URL)
	if err == nil {
		ret.Host = strings.TrimPrefix(url.Hostname(), "www.")
	}
	return ret
}

// item is the same as the hn.Item, but adds the Host field
type item struct {
	hn.Item
	Host string
}

type templateData struct {
	Stories []item
	Time    time.Duration
}
