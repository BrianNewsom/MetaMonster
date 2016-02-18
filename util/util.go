package util

import (
	"crypto/tls"
	"errors"
	"golang.org/x/net/publicsuffix"
	"net/http"
	"net/http/cookiejar"
)

func BuildHttpClient(insecureSkipVerify bool, cookieJar bool, maxRedirects int, client *http.Client) {
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
