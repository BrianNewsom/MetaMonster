package main

import (
	"crypto/tls"
	"errors"
	"fmt"
	"golang.org/x/net/html"
	"golang.org/x/net/publicsuffix"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
)

type Metadata struct {
	TitleText     string
	Title         string
	Author        string
	Description   string
	Image         string
	PublishedDate string
	URL           url.URL
}

func main() {
	u := "http://www.theguardian.com/technology/2016/feb/17/apple-fbi-encryption-san-bernardino-russia-china"

	m := GetInfoForUrl(u)

	printMetadata(*m)
}

func GetInfoForUrl(u string) *Metadata {
	m := Metadata{}

	var client http.Client
	buildHttpClient(true, true, 10, &client)
	req, _ := http.NewRequest("GET", u, nil)
	resp, _ := client.Do(req)

	parseData(resp.Body, &m)

	return &m
}

func printMetadata(m Metadata) {
	fmt.Printf("TitleText: %s\n", m.TitleText)
	fmt.Printf("Title: %s\n", m.Title)
	fmt.Printf("Author: %s\n", m.Author)
	fmt.Printf("Description: %s\n", m.Description)
	fmt.Printf("Image: %s\n", m.Image)
	fmt.Printf("Published Date: %s\n", m.PublishedDate)
	fmt.Printf("URL: %s\n", m.URL.String())
}

func parseData(b io.Reader, m *Metadata) {
	d, _ := html.Parse(b)
	var f func(*html.Node)

	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "title" {
			title := n.FirstChild.Data
			m.Title = title
		}
		if n.Type == html.ElementNode && n.Data == "meta" {
			match := map[string]bool{}
			for _, a := range n.Attr {
				// fmt.Printf("%s - %s\n", m.Key, m.Val)
				if a.Key == "name" && (a.Val == "description" || a.Val == "twitter:description") {
					match["description"] = true
				}
				if a.Key == "content" && match["description"] {
					m.Description = a.Val
					match["description"] = false
				}
				if a.Key == "property" && a.Val == "og:title" {
					match["og:title"] = true
				}
				if a.Key == "content" && match["og:title"] {
					m.TitleText = a.Val
					match["og:title"] = false
				}
				if (a.Key == "name" && a.Val == "author") || (a.Key == "property" && a.Val == "article:author") {
					match["author"] = true
				}
				if a.Key == "content" && match["author"] {
					m.Author = a.Val
					match["author"] = false
				}
				if a.Key == "property" && a.Val == "og:image" {
					match["image"] = true
				}
				if a.Key == "content" && match["image"] {
					m.Image = a.Val
					match["image"] = false
				}
				if a.Key == "property" && (a.Val == "article:published_time" || a.Val == "article:published") {
					match["published"] = true
				}
				if a.Key == "content" && match["published"] {
					m.PublishedDate = a.Val
					match["published"] = false
				}
				if a.Key == "property" && a.Val == "og:url" {
					match["url"] = true
				}
				if a.Key == "content" && match["url"] {
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

func buildHttpClient(insecureSkipVerify bool, cookieJar bool, maxRedirects int, client *http.Client) {
	// If we're having ssl issues, enable this to ignore the cert
	client.Transport = &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: insecureSkipVerify,
		},
	}

	// Set maximum redirects
	client.CheckRedirect = func() func(req *http.Request, via []*http.Request) error {
		redirects := 0
		return func(req *http.Request, via []*http.Request) error {
			if redirects > maxRedirects {
				return errors.New("stopped after 30 redirects")
			}
			redirects++
			return nil
		}
	}()

	// To store and pass cookies - some sites require this.
	if cookieJar {
		options := cookiejar.Options{
			PublicSuffixList: publicsuffix.List,
		}

		jar, err := cookiejar.New(&options)
		if err != nil {
			panic(err)
		}

		client.Jar = jar
	}
}
