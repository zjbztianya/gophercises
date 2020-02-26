package main

import (
	"container/list"
	"flag"
	"fmt"
	"github.com/zjbztianya/gophercises/link"
	"net/http"
	"net/url"
	"strings"
)

func withPrefix(baseUrl string) func(string) bool {
	return func(link string) bool {
		return strings.HasPrefix(link, baseUrl)
	}
}

func fitter(links []string, fitterFn func(string) bool) []string {
	var ret []string
	for _, page := range links {
		if fitterFn(page) {
			ret = append(ret, page)
		}
	}
	return ret
}

func getUrls(baseUrl string, pages []link.Link) []string {
	var ret []string
	for _, page := range pages {
		switch {
		case strings.HasPrefix(page.Href, "/"):
			ret = append(ret, baseUrl+page.Href)
		case strings.HasPrefix(page.Href, "http"):
			ret = append(ret, page.Href)
		}
	}
	return fitter(ret, withPrefix(baseUrl))
}

func httpGet(urlStr string) []string {
	resp, err := http.Get(urlStr)
	if err != nil {
		return nil
	}
	defer resp.Body.Close()

	pages, _ := link.Parse(resp.Body)
	reqUrl := resp.Request.URL
	baseUrl := url.URL{Scheme: reqUrl.Scheme, Host: reqUrl.Host}
	return getUrls(baseUrl.String(), pages)
}

func bfs(baseUrl string, maxDepth int) []string {
	visit := make(map[string]struct{})
	queue := list.New()
	queue.PushBack(baseUrl)
	visit[baseUrl] = struct{}{}
	var ret []string
	for i := 0; i < maxDepth && queue.Len() > 0; i++ {
		urlStr := queue.Front().Value.(string)
		queue.Remove(queue.Front())
		ret = append(ret, urlStr)
		for _, page := range httpGet(urlStr) {
			if _, ok := visit[page]; !ok {
				queue.PushBack(page)
				visit[page] = struct{}{}
			}
		}
	}
	return ret
}

func main() {
	urlFlag := flag.String("url", "https://gophercises.com", "url for build map site ")
	depthFlag := flag.Int("depth", 5, "search depth for link ")
	flag.Parse()
	pages := bfs(*urlFlag, *depthFlag)
	for _, page := range pages {
		fmt.Println(page)
	}
}
