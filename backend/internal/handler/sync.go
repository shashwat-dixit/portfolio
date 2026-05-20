package handler

import (
	"encoding/json"
	"log/slog"
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

func (h *SyncHandler) Sync(w http.ResponseWriter, r *http.Request) {
	result, err := h.svc.Sync(r.Context())
	if err != nil {
		slog.Error("sync failed", "error", err)
		http.Error(w, `{"error":"sync failed"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}
