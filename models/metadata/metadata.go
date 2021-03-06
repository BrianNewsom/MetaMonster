package metadata

import (
	"encoding/json"
	"fmt"
	"net/url"
)

type Metadata struct {
	/* TitleText often includes */
	HTMLTitle     string   `json:"html_title"`
	Title         string   `json:"title"`
	Author        string   `json:"author"`
	Description   string   `json:"description"`
	Tags          []string `json:"tags"`
	Image         string   `json:"image"`
	PublishedDate string   `json:"published_date"`
	URL           url.URL  `json:"url"`
}

func (m Metadata) ToJson() ([]byte, error) {
	j, err := json.Marshal(m)

	if err != nil {
		return []byte{0}, err
	}

	return j, nil
}

func PrintMetadata(m Metadata) {
	fmt.Printf("---- Printing Metadata ----\n")
	fmt.Printf("HTML Title: %s\n", m.HTMLTitle)
	fmt.Printf("Title: %s\n", m.Title)
	fmt.Printf("Author: %s\n", m.Author)
	fmt.Printf("Description: %s\n", m.Description)
	fmt.Printf("Tags: %s\n", m.Tags)
	fmt.Printf("Image: %s\n", m.Image)
	fmt.Printf("Published Date: %s\n", m.PublishedDate)
	fmt.Printf("URL: %s\n", m.URL.String())
	fmt.Printf("---------------------------\n")
}
