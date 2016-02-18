package main

import (
	"fmt"
	"net/url"
)

type Metadata struct {
	/* TitleText often includes */
	HTMLTitle     string
	Title         string
	Author        string
	Description   string
	Image         string
	PublishedDate string
	URL           url.URL
}

func PrintMetadata(m Metadata) {
	fmt.Printf("---- Printing Metadata ----\n")
	fmt.Printf("HTML Title: %s\n", m.HTMLTitle)
	fmt.Printf("Title: %s\n", m.Title)
	fmt.Printf("Author: %s\n", m.Author)
	fmt.Printf("Description: %s\n", m.Description)
	fmt.Printf("Image: %s\n", m.Image)
	fmt.Printf("Published Date: %s\n", m.PublishedDate)
	fmt.Printf("URL: %s\n", m.URL.String())
	fmt.Printf("---------------------------\n")
}
