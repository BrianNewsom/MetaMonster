package fetcher

import (
	"errors"
	"golang.org/x/net/html"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/briannewsom/metamonster/models/metadata"
	"github.com/briannewsom/metamonster/util"
)

type MetaTag html.Attribute

func GetInfoForUrl(u string) (*metadata.Metadata, error) {
	m := metadata.Metadata{}

	var client http.Client
	util.BuildHttpClient(true, true, 10, &client)
	req, _ := http.NewRequest("GET", u, nil)
	resp, err := client.Do(req)

	if err != nil {
		return nil, errors.New("GET request failed")
	}

	ParseData(resp.Body, &m)

	/* Fill in any data that we can assume if we don't have metadata */
	if m.URL.String() == "" {
		givenU, _ := url.Parse(u)
		m.URL = *givenU
	}

	return &m, nil
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
				t := MetaTag(a)
				if descriptionMatcher(t) {
					m.Description = getContent(n.Attr)
				} else if titleMatcher(t) {
					m.Title = getContent(n.Attr)
				} else if authorMatcher(t) {
					m.Author = getContent(n.Attr)
				} else if imageMatcher(t) {
					m.Image = getContent(n.Attr)
				} else if publishedMatcher(t) {
					m.PublishedDate = getContent(n.Attr)
				} else if urlMatcher(t) {
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

func (t MetaTag) matchesOneOf(key string, vals []string) bool {
	for _, val := range vals {
		if t.matches(key, val) {
			return true
		}
	}

	return false
}

func (t MetaTag) matches(key string, val string) bool {
	if strings.ToLower(t.Key) == strings.ToLower(key) && strings.ToLower(t.Val) == strings.ToLower(val) {
		return true
	}

	return false
}

func descriptionMatcher(t MetaTag) bool {
	return t.matchesOneOf("name", []string{"description", "twitter:description"})
}

func titleMatcher(t MetaTag) bool {
	return t.matches("property", "og:title")
}

func authorMatcher(t MetaTag) bool {
	return t.matchesOneOf("name", []string{"author", "sailthru.author"}) ||
		t.matches("property", "article:author")
}

func imageMatcher(t MetaTag) bool {
	return t.matches("property", "og:image")
}

func publishedMatcher(t MetaTag) bool {
	return t.matchesOneOf("property", []string{"article:published_time", "article:published"}) ||
		t.matches("name", "sailthru.date") ||
		t.matches("itemprop", "datepublished")
}

func urlMatcher(t MetaTag) bool {
	return t.matches("property", "og:url")
}

func getContent(attr []html.Attribute) string {
	for _, a := range attr {
		if strings.ToLower(a.Key) == "content" {
			return strings.ToLower(a.Val)
		}
	}
	return ""
}
