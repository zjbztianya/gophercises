package main

import (
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

func getUrls(baseUrl string, links []link.Link) []string {
	var ret []string
	for _, page := range links {
		switch {
		case strings.HasPrefix(page.Href, "/"):
			ret = append(ret, baseUrl+page.Href)
		case strings.HasPrefix(page.Href, "http"):
			ret = append(ret, page.Href)
		}
	}
	return fitter(ret, withPrefix(baseUrl))
}

func bfs(baseUrl string, links []link.Link) []string {
	//visit := make(map[string]struct{})
	return nil
}

func main() {
	urlFlag := flag.String("url", "https://www.baidu.com/", "url for build map site ")
	flag.Parse()

	resp, err := http.Get(*urlFlag)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	links, _ := link.Parse(resp.Body)
	reqUrl := resp.Request.URL
	baseUrl := url.URL{Scheme: reqUrl.Scheme, Host: reqUrl.Host}
	pages := bfs(baseUrl.String(), links)
	for _, page := range pages {
		fmt.Println(page)
	}
}
