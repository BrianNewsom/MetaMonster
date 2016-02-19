package main

import (
	"flag"
	"fmt"

	"github.com/briannewsom/metamonster/fetcher"
	"github.com/briannewsom/metamonster/models/metadata"
)

func main() {
	/* If given a url, use it, otherwise default to defaultUrl */
	var u = flag.String("url", "", "URL from which to retrieve metadata")
	var format = flag.String("format", "plaintext", "Output data format. Options - [json,plaintext]")

	flag.Parse()

	if *u == "" {
		fmt.Printf("No URL provided, please provide a url using -url=url")
	} else {
		m, err := fetcher.GetInfoForUrl(*u)

		if err != nil {
			fmt.Printf("Failed to retrieve metadata - %s", err)
			return
		}

		switch *format {
		case "plaintext":
			metadata.PrintMetadata(*m)
		case "json":
			fmt.Printf("%s", m.ToJson())
		default:
			fmt.Printf("Unrecognized output format %s.  Please try json or plaintext", *format)
		}
	}
}
