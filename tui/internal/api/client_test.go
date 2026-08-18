package api

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestListPostsPaginates(t *testing.T) {
	page1 := listResponse{
		Posts: []PostSummary{{Slug: "one", Title: "One"}},
		Pagination: pagination{Page: 1, Limit: 100, Total: 2, TotalPages: 2},
	}
	page2 := listResponse{
		Posts: []PostSummary{{Slug: "two", Title: "Two"}},
		Pagination: pagination{Page: 2, Limit: 100, Total: 2, TotalPages: 2},
	}

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/posts" {
			t.Errorf("path = %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		switch r.URL.Query().Get("page") {
		case "2":
			json.NewEncoder(w).Encode(page2)
		default:
			json.NewEncoder(w).Encode(page1)
		}
	}))
	defer srv.Close()

	posts, err := New(srv.URL).ListPosts(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(posts) != 2 || posts[0].Slug != "one" || posts[1].Slug != "two" {
		t.Fatalf("posts = %+v", posts)
	}
}

func TestGetPost(t *testing.T) {
	published := time.Date(2026, 1, 25, 0, 0, 0, 0, time.UTC)
	want := Post{
		Slug:      "hello",
		Title:     "Hello",
		ContentMD: "---\ntitle: Hello\nslug: hello\n---\n\n# Hi\n",
		Date:      &published,
	}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/posts/hello" {
			http.NotFound(w, r)
			return
		}
		json.NewEncoder(w).Encode(want)
	}))
	defer srv.Close()

	got, err := New(srv.URL).GetPost(context.Background(), "hello")
	if err != nil {
		t.Fatal(err)
	}
	if got.Title != "Hello" || got.Slug != "hello" {
		t.Fatalf("got %+v", got)
	}
}

func TestGetPostNotFound(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, `{"error":"not found"}`, http.StatusNotFound)
	}))
	defer srv.Close()

	_, err := New(srv.URL).GetPost(context.Background(), "missing")
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestStripFrontmatter(t *testing.T) {
	in := "---\ntitle: Hello\nslug: hello\n---\n\n# Body\n\nHi."
	got := StripFrontmatter(in)
	if got != "# Body\n\nHi." {
		t.Fatalf("got %q", got)
	}
	if StripFrontmatter("# already body") != "# already body" {
		t.Fatal("should leave body-only markdown alone")
	}
}
