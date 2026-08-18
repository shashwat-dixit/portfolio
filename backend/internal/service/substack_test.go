package service

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"gitlab.com/shashwat-dixit/portfolio/backend/internal/config"
	"gitlab.com/shashwat-dixit/portfolio/backend/internal/model"
)

func TestSubstackCookieHeader(t *testing.T) {
	got := substackCookieHeader("", "s%3Aabc", "")
	if got != "substack.sid=s:abc; connect.sid=s:abc" {
		t.Fatalf("got %q", got)
	}
	got = substackCookieHeader("substack.sid=raw", "ignored", "")
	if got != "substack.sid=raw" {
		t.Fatalf("full cookies win: %q", got)
	}
}

func TestSubstackPublishCreatesDraft(t *testing.T) {
	var gotCookie, createBody string
	var published bool
	mux := http.NewServeMux()
	srv := httptest.NewServer(mux)
	defer srv.Close()

	mux.HandleFunc("/api/v1/user/profile/self", func(w http.ResponseWriter, r *http.Request) {
		gotCookie = r.Header.Get("Cookie")
		json.NewEncoder(w).Encode(map[string]any{
			"id": 99,
			"publicationUsers": []map[string]any{
				{"id": 7, "publication": map[string]any{"subdomain": "notes", "custom_domain": ""}},
			},
		})
	})
	mux.HandleFunc("/api/v1/drafts", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method", http.StatusMethodNotAllowed)
			return
		}
		b, _ := io.ReadAll(r.Body)
		createBody = string(b)
		json.NewEncoder(w).Encode(map[string]any{"id": 123, "slug": "hello"})
	})
	mux.HandleFunc("/api/v1/drafts/123/publish", func(w http.ResponseWriter, r *http.Request) {
		published = true
		json.NewEncoder(w).Encode(map[string]any{"id": 123, "slug": "hello", "canonical_url": "https://notes.substack.com/p/hello"})
	})

	client := NewSubstackClient(&config.Config{
		SiteURL:                "https://shashwatdixit.com",
		SubstackPublicationURL: "https://notes.substack.com",
		SubstackSID:            "sid-value",
		SubstackPublish:        true,
	})
	client.http = srv.Client()
	client.baseURL = srv.URL + "/api/v1"
	client.profileURL = srv.URL + "/api/v1/user/profile/self"
	client.publicationURL = "https://notes.substack.com"

	id, url, err := client.Publish(context.Background(), model.SyndicatablePost{
		Slug:        "hello",
		Title:       "Hello",
		Description: "A greeting",
		ContentHTML: "<p>Hi</p>",
	})
	if err != nil {
		t.Fatal(err)
	}
	if id != "123" {
		t.Fatalf("id %q", id)
	}
	if url != "https://notes.substack.com/p/hello" {
		t.Fatalf("url %q", url)
	}
	if !strings.Contains(gotCookie, "substack.sid=sid-value") {
		t.Fatalf("cookie %q", gotCookie)
	}
	if !published {
		t.Fatal("expected publish call")
	}
	var payload map[string]any
	if err := json.Unmarshal([]byte(createBody), &payload); err != nil {
		t.Fatal(err)
	}
	if payload["draft_title"] != "Hello" || payload["draft_slug"] != "hello" {
		t.Fatalf("title/slug %+v", payload)
	}
	body, _ := payload["draft_body"].(string)
	if !strings.Contains(body, "<p>Hi</p>") {
		t.Fatalf("draft_body %q", body)
	}
	bylines, _ := payload["draft_bylines"].([]any)
	if len(bylines) != 1 {
		t.Fatalf("bylines %+v", payload["draft_bylines"])
	}
	byline, _ := bylines[0].(map[string]any)
	if byline["publicationUserId"] != float64(7) {
		t.Fatalf("publicationUserId %+v", byline)
	}
}

func TestSubstackUpdatePutsDraft(t *testing.T) {
	var gotPath string
	mux := http.NewServeMux()
	srv := httptest.NewServer(mux)
	defer srv.Close()
	mux.HandleFunc("/api/v1/user/profile/self", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]any{"id": 1, "publicationUsers": []any{}})
	})
	mux.HandleFunc("/api/v1/drafts/55", func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.Method + " " + r.URL.Path
		json.NewEncoder(w).Encode(map[string]any{"id": "55", "slug": "hello"})
	})

	client := NewSubstackClient(&config.Config{
		SiteURL:                "https://shashwatdixit.com",
		SubstackPublicationURL: "https://notes.substack.com",
		SubstackSID:            "sid",
	})
	client.http = srv.Client()
	client.baseURL = srv.URL + "/api/v1"
	client.profileURL = srv.URL + "/api/v1/user/profile/self"

	id, url, err := client.Update(context.Background(), "55", model.SyndicatablePost{
		Slug: "hello", Title: "Hello", ContentHTML: "<p>Hi</p>",
	})
	if err != nil {
		t.Fatal(err)
	}
	if id != "55" {
		t.Fatalf("id %q", id)
	}
	if !strings.Contains(url, "/p/hello") {
		t.Fatalf("url %q", url)
	}
	if gotPath != "PUT /api/v1/drafts/55" {
		t.Fatalf("path %q", gotPath)
	}
}

func TestNewSubstackClientDisabledWithoutCredentials(t *testing.T) {
	c := NewSubstackClient(&config.Config{SubstackPublicationURL: "https://x.substack.com"})
	if c.Enabled() {
		t.Fatal("should be disabled without cookies")
	}
}
