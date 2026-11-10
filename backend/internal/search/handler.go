package search

import (
	"encoding/json"
	"net/http"
	"strconv"
)

// Handler serves HTTP endpoints for SearchLite.
type Handler struct {
	service *Service
}

// NewHandler creates a new search handler.
func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// RegisterRoutes registers search endpoints onto the given ServeMux.
func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/search", h.HandleSearch)
	mux.HandleFunc("POST /api/documents", h.HandleCreateDocument)
	mux.HandleFunc("DELETE /api/documents/{id}", h.HandleDeleteDocument)
}

// HandleSearch processes search queries with facets and telemetry.
func (h *Handler) HandleSearch(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("q")
	category := r.URL.Query().Get("category")

	var maxPrice float64
	if mpStr := r.URL.Query().Get("max_price"); mpStr != "" {
		if p, err := strconv.ParseFloat(mpStr, 64); err == nil {
			maxPrice = p
		}
	}

	limit := 20
	if lStr := r.URL.Query().Get("limit"); lStr != "" {
		if l, err := strconv.Atoi(lStr); err == nil && l > 0 {
			limit = l
		}
	}

	params := SearchQueryParams{
		Query:    q,
		Category: category,
		MaxPrice: maxPrice,
		Limit:    limit,
	}

	res, err := h.service.Search(r.Context(), params)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	respondJSON(w, http.StatusOK, res)
}

// HandleCreateDocument indexes a new document.
func (h *Handler) HandleCreateDocument(w http.ResponseWriter, r *http.Request) {
	var input CreateDocumentInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		respondJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request payload: " + err.Error()})
		return
	}

	doc, err := h.service.IndexDocument(r.Context(), input)
	if err != nil {
		respondJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	respondJSON(w, http.StatusCreated, doc)
}

// HandleDeleteDocument deletes an existing document by ID.
func (h *Handler) HandleDeleteDocument(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		respondJSON(w, http.StatusBadRequest, map[string]string{"error": "document id is required"})
		return
	}

	deleted, err := h.service.DeleteDocument(r.Context(), id)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	if !deleted {
		respondJSON(w, http.StatusNotFound, map[string]string{"error": "document not found"})
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{"deleted": true, "id": id})
}

func respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}
