package infofetcher

import (
	"golang.org/x/net/html"
	"io"
	"net/http"
	"net/url"

	"github.com/briannewsom/metamonster/util"
)

func GetInfoForUrl(u string) *Metadata {
	m := Metadata{}

	var client http.Client
	util.BuildHttpClient(true, true, 10, &client)
	req, _ := http.NewRequest("GET", u, nil)
	resp, _ := client.Do(req)

	ParseData(resp.Body, &m)

	return &m
}

func ParseData(b io.Reader, m *Metadata) {
	d, _ := html.Parse(b)
	var f func(*html.Node)

	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "title" {
			title := n.FirstChild.Data
			m.HTMLTitle = title
		}
		if n.Type == html.ElementNode && n.Data == "meta" {
			match := map[string]bool{}
			for _, a := range n.Attr {
				if a.Key == "name" && (a.Val == "description" || a.Val == "twitter:description") {
					match["description"] = true
				}
				if a.Key == "content" && match["description"] && a.Val != "" {
					m.Description = a.Val
					match["description"] = false
				}
				if a.Key == "property" && a.Val == "og:title" {
					match["title"] = true
				}
				if a.Key == "content" && match["title"] && a.Val != "" {
					m.Title = a.Val
					match["title"] = false
				}
				if (a.Key == "name" && a.Val == "author") || (a.Key == "property" && a.Val == "article:author") {
					match["author"] = true
				}
				if a.Key == "content" && match["author"] && a.Val != "" {
					m.Author = a.Val
					match["author"] = false
				}
				if a.Key == "property" && a.Val == "og:image" {
					match["image"] = true
				}
				if a.Key == "content" && match["image"] && a.Val != "" {
					m.Image = a.Val
					match["image"] = false
				}
				if a.Key == "property" && (a.Val == "article:published_time" || a.Val == "article:published") {
					match["published"] = true
				}
				if a.Key == "content" && match["published"] && a.Val != "" {
					m.PublishedDate = a.Val
					match["published"] = false
				}
				if a.Key == "property" && a.Val == "og:url" {
					match["url"] = true
				}
				if a.Key == "content" && match["url"] && a.Val != "" {
					u, _ := url.Parse(a.Val)
					m.URL = *u
					match["url"] = false
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(d)
}
