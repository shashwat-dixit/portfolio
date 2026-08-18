package server

import (
	"strings"
	"testing"

	"github.com/shashwat-dixit/portfolio/tui/internal/api"
)

func TestFormatBlogIndex(t *testing.T) {
	got := formatBlogIndex(nil)
	if !strings.Contains(got, "No published posts") {
		t.Fatalf("empty list = %q", got)
	}

	got = formatBlogIndex([]api.PostSummary{
		{Slug: "hello", Title: "Hello", Description: "A greeting."},
	})
	if !strings.Contains(got, "# Blog") || !strings.Contains(got, "hello") || !strings.Contains(got, "A greeting.") {
		t.Fatalf("index = %q", got)
	}
}
