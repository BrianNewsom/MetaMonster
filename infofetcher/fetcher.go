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
			for _, a := range n.Attr {
				if a.Key == "name" && (a.Val == "description" || a.Val == "twitter:description") {
					m.Description = getContent(n.Attr)
				} else if a.Key == "property" && a.Val == "og:title" {
					m.Title = getContent(n.Attr)
				} else if (a.Key == "name" && a.Val == "author") || (a.Key == "property" && a.Val == "article:author") {
					m.Author = getContent(n.Attr)
				} else if a.Key == "property" && a.Val == "og:image" {
					m.Image = getContent(n.Attr)
				} else if a.Key == "property" && (a.Val == "article:published_time" || a.Val == "article:published") {
					m.PublishedDate = getContent(n.Attr)
				} else if a.Key == "property" && a.Val == "og:url" {
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

func getContent(attr []html.Attribute) string {
	for _, a := range attr {
		if a.Key == "content" {
			return a.Val
		}
	}
	return ""
}
