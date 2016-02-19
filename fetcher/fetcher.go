package fetcher

import (
	"golang.org/x/net/html"
	"io"
	"net/http"
	"net/url"

	"github.com/briannewsom/metamonster/models/metadata"
	"github.com/briannewsom/metamonster/util"
)

func GetInfoForUrl(u string) *metadata.Metadata {
	m := metadata.Metadata{}

	var client http.Client
	util.BuildHttpClient(true, true, 10, &client)
	req, _ := http.NewRequest("GET", u, nil)
	resp, _ := client.Do(req)

	ParseData(resp.Body, &m)

	return &m
}

func ParseData(b io.Reader, m *metadata.Metadata) {
	d, _ := html.Parse(b)
	var f func(*html.Node)

	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "title" {
			title := n.FirstChild.Data
			m.HTMLTitle = title
		}
		if n.Type == html.ElementNode && n.Data == "meta" {
			for _, a := range n.Attr {
				if descriptionMatcher(a) {
					m.Description = getContent(n.Attr)
				} else if titleMatcher(a) {
					m.Title = getContent(n.Attr)
				} else if authorMatcher(a) {
					m.Author = getContent(n.Attr)
				} else if imageMatcher(a) {
					m.Image = getContent(n.Attr)
				} else if publishedMatcher(a) {
					m.PublishedDate = getContent(n.Attr)
				} else if urlMatcher(a) {
					u, _ := url.Parse(getContent(n.Attr))
					m.URL = *u
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(d)
}

func descriptionMatcher(a html.Attribute) bool {
	if a.Key == "name" && (a.Val == "description" || a.Val == "twitter:description") {
		return true
	}
	return false
}

func titleMatcher(a html.Attribute) bool {
	if a.Key == "property" && a.Val == "og:title" {
		return true
	}
	return false
}

func authorMatcher(a html.Attribute) bool {
	if (a.Key == "name" && (a.Val == "author" || a.Val == "sailthru.author")) || (a.Key == "property" && a.Val == "article:author") {
		return true
	}
	return false
}

func imageMatcher(a html.Attribute) bool {
	if a.Key == "property" && a.Val == "og:image" {
		return true
	}
	return false
}

func publishedMatcher(a html.Attribute) bool {
	if (a.Key == "property" && (a.Val == "article:published_time" || a.Val == "article:published")) || (a.Key == "name" && a.Val == "sailthru.date") {
		return true
	}
	return false
}

func urlMatcher(a html.Attribute) bool {
	if a.Key == "property" && a.Val == "og:url" {
		return true
	}
	return false
}

func getContent(attr []html.Attribute) string {
	for _, a := range attr {
		if a.Key == "content" {
			return a.Val
		}
	}
	return ""
}
