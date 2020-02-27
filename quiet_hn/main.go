package main

import (
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

	"github.com/zjbztianya/gophercises/quiet_hn/hn"
)

// item is the same as the hn.Item, but adds the Host field
type item struct {
	hn.Item
	Host string
}

type templateData struct {
	Stories []item
	Time    time.Duration
}

type result struct {
	ok   bool
	pos  int
	item item
}

var (
	cacheExpire  time.Time
	cacheStories []item
	cacheMutex   sync.Mutex
)

func getStories(ids []int, resultCh chan result) []item {
	var client hn.Client
	for i, id := range ids {
		go func(pos, id int) {
			hnItem, err := client.GetItem(id)
			res := result{pos: pos}
			if err != nil {
				resultCh <- res
				return
			}
			item := parseHNItem(hnItem)
			if isStoryLink(item) {
				res.item = item
				res.ok = true
			}
			resultCh <- res
		}(i, id)
	}

	var results []result
	for i := 0; i < len(ids); i++ {
		result := <-resultCh
		if result.ok {
			results = append(results, result)
		}
	}
	sort.Slice(results, func(i, j int) bool {
		return results[i].pos < results[j].pos
	})

	stories := make([]item, 0, len(results))
	for _, result := range results {
		stories = append(stories, result.item)
	}
	return stories
}

func getTopStories(numStories int) ([]item, error) {
	var client hn.Client
	ids, err := client.TopItems()
	if err != nil {
		return nil, err
	}
	var stories []item
	resultCh := make(chan result, numStories)
	pos := 0
	for len(stories) < numStories {
		n := (numStories - len(stories)) * 5 / 4
		stories = append(stories, getStories(ids[pos:pos+n], resultCh)...)
		pos += n
	}
	return stories[:numStories], nil
}

func getCacheStories(numStories int) ([]item, error) {
	cacheMutex.Lock()
	defer cacheMutex.Unlock()

	if time.Now().Sub(cacheExpire) < 0 {
		return cacheStories, nil
	}
	stories, err := getTopStories(numStories)
	if err != nil {
		return nil, err
	}
	cacheStories = stories
	cacheExpire = time.Now().Add(5 * time.Minute)
	return stories, nil
}

func handler(numStories int, tpl *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		stories, err := getCacheStories(numStories)
		if err != nil {
			http.Error(w, "Failed to load top stories", http.StatusInternalServerError)
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
	}
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

func main() {
	// parse flags
	var port, numStories int
	flag.IntVar(&port, "port", 3000, "the port to start the web server on")
	flag.IntVar(&numStories, "num_stories", 30, "the number of top stories to display")
	flag.Parse()

	tpl := template.Must(template.ParseFiles("quiet_hn/index.gohtml"))

	http.HandleFunc("/", handler(numStories, tpl))

	// Start the server
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}
