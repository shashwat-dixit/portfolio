package handler

import (
	"net/http"

	"gitlab.com/shashwat-dixit/portfolio/backend/internal/service"
)

type SyncHandler struct {
	svc    *service.SyncService
	apiKey string
}

func NewSyncHandler(svc *service.SyncService, apiKey string) *SyncHandler {
	return &SyncHandler{svc: svc, apiKey: apiKey}
}

// Sync handles POST /api/sync
func (h *SyncHandler) Sync(w http.ResponseWriter, r *http.Request) {
	// TODO: call h.svc.Sync, return SyncResult JSON
	w.WriteHeader(http.StatusNotImplemented)
}
