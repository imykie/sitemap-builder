package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
	link "sitemap-builder/cmd"
	"strings"
)

func main() {
	urlFlag := flag.String("url", "https://github.com", "The URL you want to build sitemap for")
	flag.Parse()

	pages := getPage(*urlFlag)
	for _, page := range pages {
		fmt.Println(page)
	}
}

func getPage(urlFlag string) []string {
	resp, err := http.Get(urlFlag)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	reqUrl := resp.Request.URL
	baseUrl := &url.URL{
		Scheme: reqUrl.Scheme,
		Host:   reqUrl.Host,
	}

	base := baseUrl.String()
	return filter(hrefs(resp.Body, base), withPrefix(base))
}

func hrefs(body io.Reader, base string) []string {
	links, err := link.Parse(body)
	if err != nil {
		panic(err)
	}

	var hrefs []string
	for _, l := range links {
		switch {
		case strings.HasPrefix(l.Href, "/"):
			hrefs = append(hrefs, base+l.Href)
		case strings.HasPrefix(l.Href, "http"):
			hrefs = append(hrefs, l.Href)
		}
	}

	return hrefs
}

func filter(links []string, keepFn func (string) bool) []string{
	var result []string
	for _, link := range links {
		if keepFn(link) {
			result = append(result, link)
		}
	}
	return result

}

func withPrefix(prefix string) func (string) bool {
	return func (link string) bool {
		return strings.HasPrefix(link, prefix)
	}
}
