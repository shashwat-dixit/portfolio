package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"gitlab.com/shashwat-dixit/portfolio/backend/internal/cache"
	"gitlab.com/shashwat-dixit/portfolio/backend/internal/repository"
)

type TagHandler struct {
	repo  *repository.TagRepo
	cache *cache.RedisCache
}

func NewTagHandler(repo *repository.TagRepo, cache *cache.RedisCache) *TagHandler {
	return &TagHandler{repo: repo, cache: cache}
}

func (h *TagHandler) List(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	cacheKey := cache.TagsKey()

	cached, err := h.cache.Get(ctx, cacheKey)
	if err != nil {
		slog.Warn("cache get error", "key", cacheKey, "error", err)
	}
	if cached != "" {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Cache-Control", "public, max-age=3600, stale-while-revalidate=86400")
		w.Write([]byte(cached))
		return
	}

	tags, err := h.repo.ListWithCounts(ctx)
	if err != nil {
		http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
		return
	}

	resp := map[string]any{"tags": tags}
	data, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
		return
	}

	if err := h.cache.Set(ctx, cacheKey, string(data)); err != nil {
		slog.Warn("cache set error", "key", cacheKey, "error", err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "public, max-age=3600, stale-while-revalidate=86400")
	w.Write(data)
}
