package handler

import (
	"net/http"

	"gitlab.com/shashwat-dixit/portfolio/backend/internal/service"
)

type PostHandler struct {
	svc *service.PostService
}

func NewPostHandler(svc *service.PostService) *PostHandler {
	return &PostHandler{svc: svc}
}

// List handles GET /api/posts
// Query params: ?tag=, ?page=, ?limit=, ?q=
func (h *PostHandler) List(w http.ResponseWriter, r *http.Request) {
	// TODO: parse query params, call h.svc.List, return JSON
	w.WriteHeader(http.StatusNotImplemented)
}

// GetBySlug handles GET /api/posts/{slug}
func (h *PostHandler) GetBySlug(w http.ResponseWriter, r *http.Request) {
	// TODO: extract slug from path, call h.svc.GetBySlug, return JSON
	w.WriteHeader(http.StatusNotImplemented)
}
