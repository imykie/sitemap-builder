package main

import (
	"encoding/xml"
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	link "sitemap-builder/cmd"
	"strconv"
	"strings"
)

const xmlns = "http://www.sitemaps.org/schemas/sitemap/0.9"

type loc struct {
	Value string `xml:"loc"`
}

type urlset struct {
	Urls    []loc  `xml:"url"`
	Xmlns   string `xml:"xmlns,attr"`
	Comment string `xml:",comment"`
}

func main() {
	urlFlag := flag.String("url", "https://go.dev", "The URL you want to build Sitemap for")
	maxDepth := flag.Int("depth", 3, "The maximum depth of the Sitemap Builder")
	flag.Parse()

	toXml := urlset{
		Xmlns: xmlns,
	}
	pages := bfs(*urlFlag, *maxDepth)
	for _, page := range pages {
		toXml.Urls = append(toXml.Urls, loc{page})
	}
	toXml.Comment = "Total number of URLs: " + strconv.Itoa(len(pages))

	fmt.Printf(xml.Header)
	enc := xml.NewEncoder(os.Stdout)
	enc.Indent(" ", "  ")
	if err := enc.Encode(toXml); err != nil {
		panic(err)
	}
	fmt.Println("")
}

func bfs(baseUrl string, depth int) []string {
	reached := make(map[string]struct{})
	var queue map[string]struct{}
	nextQueue := map[string]struct{}{
		baseUrl: struct{}{},
	}
	for i := 0; i <= depth; i++ {
		queue, nextQueue = nextQueue, make(map[string]struct{})
		if len(queue) == 0 {
			break
		}
		for urlStr, _ := range queue {
			if _, ok := reached[urlStr]; ok {
				continue
			}
			reached[urlStr] = struct{}{}
			for _, l := range getPage(urlStr) {
				nextQueue[l] = struct{}{}
			}
		}
	}
	result := make([]string, 0, len(reached))
	for urlStr, _ := range reached {
		result = append(result, urlStr)
	}
	return result
}

func getPage(urlFlag string) []string {
	resp, err := http.Get(urlFlag)
	if err != nil {
		return []string{}
	}
	defer resp.Body.Close()

	reqUrl := resp.Request.URL
	baseUrl := &url.URL{
		Scheme: reqUrl.Scheme,
		Host:   reqUrl.Host,
	}
	base := baseUrl.String()

	// NOTE: Check functions must resolve to true for url to be included.
	return filter(hrefs(resp.Body, base), withPrefix(base), withoutSubstring("#"))
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

func filter(links []string, checkFns ...func(string) bool) []string {
	var result []string
	for _, lnk := range links {
		r := true
		for i, _ := range checkFns {
			r = r && checkFns[i](lnk)
		}
		if r {
			result = append(result, lnk)
		}
	}
	return result
}

func withPrefix(prefix string) func(string) bool {
	return func(link string) bool {
		return strings.HasPrefix(link, prefix)
	}
}

func withoutSubstring(substr string) func(string) bool {
	return func(link string) bool {
		return !strings.Contains(link, substr)
	}
}
