package main

import (
	"crypto/tls"
	"errors"
	"fmt"
	"golang.org/x/net/html"
	"golang.org/x/net/publicsuffix"
	"io"
	//	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
)

type Metadata struct {
	TitleText   string
	Title       string
	Description string
}

func main() {
	var client http.Client

	buildHttpClient(true, true, 10, &client)

	url := "http://www.theguardian.com/technology/2016/feb/17/apple-fbi-encryption-san-bernardino-russia-china"

	req, _ := http.NewRequest("GET", url, nil)

	resp, _ := client.Do(req)

	// print(*resp)
	m := Metadata{}

	parseData(resp.Body, &m)

	print(m)
}

func print(m Metadata) {
	fmt.Printf("TitleText: %s\n", m.TitleText)
	fmt.Printf("Title: %s\n", m.Title)
	fmt.Printf("Description: %s\n", m.Description)
}

func parseData(b io.ReadCloser, m *Metadata) {
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
				if a.Key == "name" && a.Val == "description" {
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
			}

		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(d)
}

/*
func print(resp http.Response) {
	fmt.Printf("Header\n")
	for key, val := range resp.Header {
		fmt.Printf("%s => %s\n", key, val)
	}

	fmt.Printf("Body\n")

	body, _ := ioutil.ReadAll(resp.Body)

	fmt.Printf("%s\n", body)
}
*/

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
