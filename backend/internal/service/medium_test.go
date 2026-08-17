package service

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"gitlab.com/shashwat-dixit/portfolio/backend/internal/model"
)

func TestMediumPublishCreatesPost(t *testing.T) {
	var gotAuth, gotPath, gotBody string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotAuth = r.Header.Get("Authorization")
		gotPath = r.URL.Path
		switch {
		case r.Method == http.MethodGet && r.URL.Path == "/me":
			json.NewEncoder(w).Encode(map[string]any{"data": map[string]any{"id": "user-1"}})
		case r.Method == http.MethodPost && r.URL.Path == "/users/user-1/posts":
			b, _ := io.ReadAll(r.Body)
			gotBody = string(b)
			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(map[string]any{
				"data": map[string]any{"id": "post-9", "url": "https://medium.com/@me/hello-post-9"},
			})
		default:
			http.NotFound(w, r)
		}
	}))
	defer srv.Close()

	client := &MediumClient{
		token:         "tok",
		publishStatus: "public",
		siteURL:       "https://shashwatdixit.com",
		baseURL:       srv.URL,
		http:          srv.Client(),
	}

	id, url, err := client.Publish(context.Background(), model.SyndicatablePost{
		Slug:      "hello",
		Title:     "Hello",
		ContentMD: "---\ntitle: Hello\n---\n\nBody.\n",
		Tags:      []string{"go", "blog"},
	})
	if err != nil {
		t.Fatal(err)
	}
	if id != "post-9" || url != "https://medium.com/@me/hello-post-9" {
		t.Fatalf("id=%s url=%s", id, url)
	}
	if gotAuth != "Bearer tok" {
		t.Fatalf("auth %q", gotAuth)
	}
	if gotPath != "/users/user-1/posts" {
		t.Fatalf("path %q", gotPath)
	}
	for _, want := range []string{`"contentFormat":"markdown"`, `"canonicalUrl":"https://shashwatdixit.com/blog/hello"`, `"publishStatus":"public"`} {
		if !strings.Contains(gotBody, want) {
			t.Fatalf("body missing %s: %s", want, gotBody)
		}
	}
}

func TestMediumPublishToPublication(t *testing.T) {
	var gotPath string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		if r.URL.Path == "/me" {
			json.NewEncoder(w).Encode(map[string]any{"data": map[string]any{"id": "user-1"}})
			return
		}
		json.NewEncoder(w).Encode(map[string]any{"data": map[string]any{"id": "p1", "url": "https://medium.com/pub/p1"}})
	}))
	defer srv.Close()

	client := &MediumClient{
		token:         "tok",
		publicationID: "pub-42",
		publishStatus: "draft",
		siteURL:       "https://shashwatdixit.com",
		baseURL:       srv.URL,
		http:          srv.Client(),
	}
	if _, _, err := client.Publish(context.Background(), model.SyndicatablePost{Slug: "x", Title: "X", ContentMD: "hi"}); err != nil {
		t.Fatal(err)
	}
	if gotPath != "/publications/pub-42/posts" {
		t.Fatalf("path %q", gotPath)
	}
}

func TestMediumDisabledWithoutToken(t *testing.T) {
	var c *MediumClient
	if c.Enabled() {
		t.Fatal("nil client should be disabled")
	}
	c = &MediumClient{}
	if c.Enabled() {
		t.Fatal("empty token should be disabled")
	}
}
