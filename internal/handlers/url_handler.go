package handlers

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strings"
     "shortify/internal/logger"
	"shortify/internal/services"
)

type URLHandler struct {
	service *services.URLService
}

type CreateURLRequest struct {
	URL string `json:"url"`
}

func NewURLHandler(service *services.URLService) *URLHandler {
	return &URLHandler{
		service: service,
	}
}


func isValidURL(input string) bool {
	parsed, err := url.ParseRequestURI(input)
	if err!= nil{
		return false
	}

	if parsed.Scheme !="http" && parsed.Scheme != "https"{
		return false
	}

	return true
}


func (h *URLHandler) CreateShortURL(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	defer r.Body.Close()

	var request CreateURLRequest

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		logger.Log.Warn("Invalid request body")
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if request.URL == "" {
		logger.Log.Warn("Empty URL recieved")
		http.Error(w, "URL is required", http.StatusBadRequest)
		return
	}

	if len(request.URL)>2048{
		logger.Log.Warn("URL too long","url",request.URL)
		http.Error(w, "URL too long", http.StatusBadRequest)
		return
	}

	if !isValidURL(request.URL){
		logger.Log.Warn("Invalid URL format", "url", request.URL)
		http.Error(w, "Invalid URL format", http.StatusBadRequest)
		return
	}

	url, err := h.service.CreateShortURL(request.URL)
	if err != nil {
		logger.Log.Error("Failed to create short URL", "error", err)
		http.Error(w, "Failed to create short URL", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	response := map[string]string{
		"short_url": "http://localhost:8080/" + url.ShortCode,
	}

	json.NewEncoder(w).Encode(response)
}


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