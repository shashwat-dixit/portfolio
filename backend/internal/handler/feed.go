package handler

import (
	"net/http"

	"gitlab.com/shashwat-dixit/portfolio/backend/internal/config"
	"gitlab.com/shashwat-dixit/portfolio/backend/internal/service"
)

type FeedHandler struct {
	svc *service.PostService
	cfg *config.Config
}

func NewFeedHandler(svc *service.PostService, cfg *config.Config) *FeedHandler {
	return &FeedHandler{svc: svc, cfg: cfg}
}

// RSS handles GET /api/feed.xml
func (h *FeedHandler) RSS(w http.ResponseWriter, r *http.Request) {
	// TODO: fetch published posts, generate RSS 2.0 XML
	w.WriteHeader(http.StatusNotImplemented)
}

// Health handles GET /api/health
func Health(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"status":"ok"}`))
}
