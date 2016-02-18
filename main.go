package main

import (
	"fmt"
	"os"

	"github.com/briannewsom/metamonster/infofetcher"
)

const defaultUrl = "http://www.theguardian.com/technology/2016/feb/17/apple-fbi-encryption-san-bernardino-russia-china"

func main() {
	/* If given a url, use it, otherwise default to defaultUrl */
	var u string

	if len(os.Args) > 1 {
		u = os.Args[1]
	} else {
		u = defaultUrl
	}

	fmt.Printf("Retrieving metadata for url %s\n", u)

	m := infofetcher.GetInfoForUrl(u)

	infofetcher.PrintMetadata(*m)
}
