package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestWantsMarkdown(t *testing.T) {
	tests := []struct {
		name   string
		url    string
		accept string
		want   bool
	}{
		{name: "browser html", url: "/api/posts/slug", accept: "text/html,application/xhtml+xml,*/*;q=0.8", want: false},
		{name: "json api", url: "/api/posts/slug", accept: "application/json", want: false},
		{name: "star accept", url: "/api/posts/slug", accept: "*/*", want: false},
		{name: "empty accept", url: "/api/posts/slug", accept: "", want: false},
		{name: "format query", url: "/api/posts/slug?format=md", accept: "application/json", want: true},
		{name: "markdown accept", url: "/api/posts/slug", accept: "text/markdown", want: true},
		{name: "x-markdown accept", url: "/api/posts/slug", accept: "text/x-markdown; charset=utf-8", want: true},
		{name: "markdown among html", url: "/api/posts/slug", accept: "text/html, text/markdown;q=0.9", want: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, tt.url, nil)
			if tt.accept != "" {
				req.Header.Set("Accept", tt.accept)
			}
			if got := wantsMarkdown(req); got != tt.want {
				t.Fatalf("wantsMarkdown() = %v, want %v", got, tt.want)
			}
		})
	}
}
