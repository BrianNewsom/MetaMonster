package fetcher

import (
	"fmt"
	"golang.org/x/net/html"
	"io"
	"net/http"
	"net/url"

	"github.com/briannewsom/metamonster/models/metadata"
	"github.com/briannewsom/metamonster/util"
)

type MetaTagAttribute html.Attribute
type MetaTagAttributes []html.Attribute

func GetInfoForUrl(u string) (*metadata.Metadata, error) {
	m := metadata.Metadata{}

	var client http.Client
	util.BuildHttpClient(true, true, 10, &client)
	req, _ := http.NewRequest("GET", u, nil)
	resp, err := client.Do(req)

	if err != nil {
		return nil, fmt.Errorf("GET request failed")
	}

	ParseData(resp.Body, &m)

	/* Fill in any data that we can assume if we don't have metadata */
	if m.URL.String() == "" {
		givenU, _ := url.Parse(u)
		m.URL = *givenU
	}

	return &m, nil
}

func ParseData(b io.Reader, m *metadata.Metadata) error {
	d, err := html.Parse(b)

	if err != nil {
		return fmt.Errorf("Failed to parse HTML - %s", err)
	}

	var f func(*html.Node)

	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "title" {
			title := n.FirstChild.Data
			m.HTMLTitle = title
		}
		if n.Type == html.ElementNode && n.Data == "meta" {
			attrs := MetaTagAttributes(n.Attr)
			for _, a := range attrs {
				t := MetaTagAttribute(a)
				if descriptionMatcher(t) {
					m.Description = getContent(attrs)
				} else if titleMatcher(t) {
					m.Title = getContent(attrs)
				} else if authorMatcher(t) {
					m.Author = getContent(attrs)
				} else if tagMatcher(t) {
					m.Tags = append(m.Tags, getContent(attrs))
				} else if imageMatcher(t) {
					m.Image = getContent(attrs)
				} else if publishedMatcher(t) {
					m.PublishedDate = getContent(attrs)
				} else if urlMatcher(t) {
					u, _ := url.Parse(getContent(attrs))
					m.URL = *u
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(d)
	return nil
}
