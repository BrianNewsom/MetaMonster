package metadata

import (
	"encoding/json"
	"fmt"
	"net/url"
)

type Metadata struct {
	/* TitleText often includes */
	HTMLTitle     string  `json:"html_title"`
	Title         string  `json:"title"`
	Author        string  `json:"author"`
	Description   string  `json:"description"`
	Image         string  `json:"image"`
	PublishedDate string  `json:"published_date"`
	URL           url.URL `json:"url"`
}

func (m Metadata) ToJson() []byte {
	j, err := json.Marshal(m)

	if err != nil {
		panic(err)
	}

	return j
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
