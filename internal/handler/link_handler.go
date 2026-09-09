package handler

import (
	"encoding/json"
	"net/http"

	"github.com/Isc740/url-shortener/internal/service"
)

type LinkHandler struct {
	Service *service.LinkService
}

func (h *LinkHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /links", h.Create)
	mux.HandleFunc("POST /links", h.Create)
}

func NewLinkHandler(service *service.LinkService) *LinkHandler {
	return &LinkHandler{
		Service: service,
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

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(link)
}
