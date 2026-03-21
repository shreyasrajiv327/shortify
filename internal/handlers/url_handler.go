package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"shortify/internal/services"
)

type URLHandler struct {
	service *services.URLService
}

func NewURLHandler(service *services.URLService) *URLHandler {
	return &URLHandler{
		service: service,
	}
}

// Request struct (clean + reusable)
type CreateURLRequest struct {
	URL string `json:"url"`
}

// POST /shorten
func (h *URLHandler) CreateShortURL(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	defer r.Body.Close()

	var request CreateURLRequest

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if request.URL == "" {
		http.Error(w, "URL is required", http.StatusBadRequest)
		return
	}

	url, err := h.service.CreateShortURL(request.URL)
	if err != nil {
		http.Error(w, "Failed to create short URL", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	response := map[string]string{
		"short_url": "http://localhost:8080/" + url.ShortCode,
	}

	json.NewEncoder(w).Encode(response)
}

// GET /{shortCode}
func (h *URLHandler) RedirectURL(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	code := strings.TrimPrefix(r.URL.Path, "/")

	if code == "" {
		http.Error(w, "Short code is required", http.StatusBadRequest)
		return
	}

	url, err := h.service.GetLongURL(code)
	if err != nil {
		http.Error(w, "URL not found", http.StatusNotFound)
		return
	}

	http.Redirect(w, r, url.LongURL, http.StatusFound)
}