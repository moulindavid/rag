package query

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
)

type queryService interface {
	Query(ctx context.Context, req *Request) (*Response, error)
}

type Handler struct {
	service queryService
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Query(w http.ResponseWriter, r *http.Request) {
	var req Request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	if req.Question == "" {
		http.Error(w, "question is required", http.StatusBadRequest)
		return
	}

	resp, err := h.service.Query(r.Context(), &req)
	if err != nil {
		slog.Error("query failed", "error", err)
		http.Error(w, "query failed", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
