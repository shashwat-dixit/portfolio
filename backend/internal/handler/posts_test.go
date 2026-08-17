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

func TestAgentBodyContentType(t *testing.T) {
	plain := httptest.NewRequest(http.MethodGet, "/api/posts/slug?format=md", nil)
	if got := agentBodyContentType(plain); got != "text/plain; charset=utf-8" {
		t.Fatalf("format=md Content-Type = %q, want text/plain", got)
	}

	md := httptest.NewRequest(http.MethodGet, "/api/posts/slug", nil)
	md.Header.Set("Accept", "text/markdown")
	if got := agentBodyContentType(md); got != "text/markdown; charset=utf-8" {
		t.Fatalf("Accept markdown Content-Type = %q, want text/markdown", got)
	}
}
