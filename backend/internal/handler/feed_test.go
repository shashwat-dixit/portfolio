package handler

import (
	"encoding/xml"
	"strings"
	"testing"
)

func TestRSSIncludesContentEncodedNamespace(t *testing.T) {
	feed := rssRoot{
		Version:   "2.0",
		AtomNS:    "http://www.w3.org/2005/Atom",
		ContentNS: "http://purl.org/rss/1.0/modules/content/",
		Channel: rssChannel{
			Title: "Blog",
			Items: []rssItem{{
				Title:          "Hello",
				Link:           "https://example.com/blog/hello",
				Description:    "teaser",
				ContentEncoded: "<p>Full <strong>HTML</strong></p>",
				GUID:           "https://example.com/blog/hello",
			}},
		},
	}
	raw, err := xml.Marshal(feed)
	if err != nil {
		t.Fatal(err)
	}
	out := string(raw)
	if !strings.Contains(out, `xmlns:content="http://purl.org/rss/1.0/modules/content/"`) {
		t.Fatalf("missing content namespace: %s", out)
	}
	if !strings.Contains(out, "Full") || !strings.Contains(out, "HTML") {
		t.Fatalf("missing full content: %s", out)
	}
}
