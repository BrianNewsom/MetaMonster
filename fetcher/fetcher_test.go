package fetcher

import (
	"errors"
	"io"
	"os"
	"testing"

	"github.com/briannewsom/metamonster/models/metadata"
)

const testDir = "_test_res"

func TestGetInfoForUrl(t *testing.T) {
	badURL := "http://this-url-really-shouldn't-exist.com"
	_, err := GetInfoForUrl(badURL)

	if err == nil {
		t.Errorf("A bad URL should give an error, got nil")
	}

}

func TestParseData(t *testing.T) {
	files := []string{"medium.html", "nytimes.html", "theguardian.html", "blog.html", "dailycamera.html", "latimes.html"}

	for _, f := range files {
		m := metadata.Metadata{}
		r, err := getFileReader(f)

		if err != nil {
			t.Errorf("Opening test file %s returned error - %s", f, err)
		} else {
			ParseData(r, &m)
			err = hasAllData(m)

			if err != nil {
				t.Errorf("File %s returned error - %s", f, err)
			}
		}
	}
}

func TestTagMatcher(t *testing.T) {
	f := "nytimes.html"

	r, err := getFileReader(f)
	m := metadata.Metadata{}

	if err != nil {
		t.Errorf("Opening test file %s returned error - %s", f, err)
	} else {
		ParseData(r, &m)
		err = hasTags(m)

		if err != nil {
			t.Errorf("File %s returned error - %s", f, err)
		}
	}
}

func getFileReader(name string) (io.Reader, error) {
	r, err := os.Open(testDir + "/" + name)

	if err != nil {
		return nil, errors.New("Failed to read file")
	}

	return r, nil
}

func hasTags(m metadata.Metadata) error {
	if len(m.Tags) == 0 {
		return errors.New("No Tags")
	}

	return nil
}

func hasAllData(m metadata.Metadata) error {
	if m.HTMLTitle == "" {
		return errors.New("No HTML Title")
	}

	if m.Title == "" {
		return errors.New("No Title")
	}

	if m.Author == "" {
		return errors.New("No Author")
	}

	if m.Description == "" {
		return errors.New("No Description")
	}

	if m.Image == "" {
		return errors.New("No Image")
	}

	if m.PublishedDate == "" {
		return errors.New("No PublishedDate")
	}

	if m.URL.String() == "" {
		return errors.New("No URL")
	}

	return nil
}

/* Where the test data is retrieved from
testUrls := []string{
	"http://www.nytimes.com/2016/02/17/us/politics/senator-charles-grassley-hearings-supreme-court-nominee.html?hp&action=click&pgtype=Homepage&clickSource=story-heading&module=first-column-region&region=top-news&WT.nav=top-news",
	"http://www.theguardian.com/technology/2016/feb/17/apple-fbi-encryption-san-bernardino-russia-china",
	"https://medium.com/@fjmubeen/ai-no-longer-understand-my-phd-dissertation-and-what-this-means-for-mathematics-education-1d40708f61c#.mtqtdl9gm",
	"http://markmanson.net/not-giving-a-fuck"}
*/
