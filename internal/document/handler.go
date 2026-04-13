package document

import (
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"strings"
)

type ingester interface {
	Ingest(ctx context.Context, filename string, content io.Reader) (*Document, error)
}

type Handler struct {
	service ingester
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Upload(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid multipart form"})
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "missing file field"})
		return
	}
	defer file.Close()

	doc, err := h.service.Ingest(r.Context(), header.Filename, file)
	if err != nil {
		if strings.HasPrefix(err.Error(), "parse document: unsupported file extension") {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
			return
		}
		slog.Error("document ingestion failed", "filename", header.Filename, "error", err)
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "internal server error"})
		return
	}

	writeJSON(w, http.StatusCreated, map[string]any{
		"id":       doc.ID,
		"filename": doc.Filename,
	})
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}

