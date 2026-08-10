package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"gitlab.com/shashwat-dixit/portfolio/backend/internal/service"
)

type PostHandler struct {
	svc *service.PostService
}

func NewPostHandler(svc *service.PostService) *PostHandler {
	return &PostHandler{svc: svc}
}

func (h *PostHandler) List(w http.ResponseWriter, r *http.Request) {
	tag := r.URL.Query().Get("tag")
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	includeDrafts := r.URL.Query().Get("drafts") == "true"

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	resp, err := h.svc.List(r.Context(), tag, page, limit, includeDrafts)
	if err != nil {
		slog.Error("list posts failed", "error", err)
		http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "public, max-age=3600, stale-while-revalidate=86400")
	json.NewEncoder(w).Encode(resp)
}

func (h *PostHandler) GetBySlug(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")
	if slug == "" {
		http.Error(w, `{"error":"slug required"}`, http.StatusBadRequest)
		return
	}

	post, tags, err := h.svc.GetBySlug(r.Context(), slug)
	if err != nil {
		http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
		return
	}
	if post == nil {
		http.Error(w, `{"error":"not found"}`, http.StatusNotFound)
		return
	}

	w.Header().Set("Cache-Control", "public, max-age=86400, stale-while-revalidate=604800")

	// Agents can request raw markdown via Accept: text/markdown or ?format=md
	if wantsMarkdown(r) {
		w.Header().Set("Content-Type", "text/markdown; charset=utf-8")
		w.Write([]byte(post.ContentMD))
		return
	}

	resp := map[string]any{
		"slug":        post.Slug,
		"title":       post.Title,
		"description": post.Description,
		"contentHtml": post.ContentHTML,
		"contentMd":   post.ContentMD,
		"tags":        tags,
		"date":        post.PublishedAt,
		"updated":     post.UpdatedAt,
		"cover":       post.CoverImage,
		"readingTime": post.ReadingTime,
		"author":      post.Author,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func wantsMarkdown(r *http.Request) bool {
	if r.URL.Query().Get("format") == "md" {
		return true
	}
	accept := r.Header.Get("Accept")
	if accept == "" {
		return false
	}
	// Prefer markdown only when it is explicitly requested (and not part of */*).
	for _, part := range strings.Split(accept, ",") {
		mediaType := strings.TrimSpace(strings.SplitN(part, ";", 2)[0])
		if mediaType == "text/markdown" || mediaType == "text/x-markdown" {
			return true
		}
	}
	return false
}
