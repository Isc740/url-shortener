package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/Isc740/url-shortener/internal/service"
)

type LinkHandler struct {
	Service *service.LinkService
	Logger  *slog.Logger
}

func (h *LinkHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /links", h.Create)
	mux.HandleFunc("POST /links", h.Create)
}

func NewLinkHandler(service *service.LinkService, logger *slog.Logger) *LinkHandler {
	return &LinkHandler{
		Service: service,
		Logger:  logger,
	}
}

func (h *LinkHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req service.CreateLinkDTO
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	link, err := h.Service.Create(r.Context(), req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.Logger.Info("HTTP REQUEST", "POST", link)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(link)
}
