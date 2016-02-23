package fetcher

import (
	"strings"
)

func (t MetaTagAttribute) matchesOneOf(key string, vals []string) bool {
	for _, val := range vals {
		if t.matches(key, val) {
			return true
		}
	}

	return false
}

func (t MetaTagAttribute) matches(key string, val string) bool {
	if strings.ToLower(t.Key) == strings.ToLower(key) && strings.ToLower(t.Val) == strings.ToLower(val) {
		return true
	}

	return false
}

func descriptionMatcher(t MetaTagAttribute) bool {
	return t.matchesOneOf("name", []string{"description", "twitter:description", "og:description"})
}

func titleMatcher(t MetaTagAttribute) bool {
	return t.matches("property", "og:title")
}

func authorMatcher(t MetaTagAttribute) bool {
	return t.matchesOneOf("name", []string{"author", "sailthru.author"}) ||
		t.matches("property", "article:author")
}

func tagMatcher(t MetaTagAttribute) bool {
	return t.matches("property", "article:tag")
}

func imageMatcher(t MetaTagAttribute) bool {
	return t.matches("property", "og:image")
}

func publishedMatcher(t MetaTagAttribute) bool {
	return t.matchesOneOf("property", []string{"article:published_time", "article:published"}) ||
		t.matches("name", "sailthru.date") ||
		t.matches("itemprop", "datepublished")
}

func urlMatcher(t MetaTagAttribute) bool {
	return t.matches("property", "og:url")
}

func getContent(attrs MetaTagAttributes) string {
	for _, a := range attrs {
		if strings.ToLower(a.Key) == "content" {
			return strings.ToLower(a.Val)
		}
	}
	return ""
}
