package handler

import (
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

// List handles GET /api/tags
func (h *TagHandler) List(w http.ResponseWriter, r *http.Request) {
	// TODO: check cache, fall back to repo, return JSON
	w.WriteHeader(http.StatusNotImplemented)
}
