package handler

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

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
	ctx, cancel := context.WithTimeout(r.Context(), 4*time.Minute)
	defer cancel()

	result, err := h.svc.Sync(ctx)
	if err != nil {
		slog.Error("sync failed", "error", err)
		http.Error(w, `{"error":"sync failed"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}
