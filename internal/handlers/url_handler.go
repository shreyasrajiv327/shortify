package handlers
import (
	"encoding/json"
	"net/http"
	"strings"

	"shortify/internal/services"
)

type URLHandler struct{
	service *services.URLService
}

func NewURLHandelr(service *services.URLService) *URLHandler{
	return &URLHandler{
		service: service,
	}
}

func (h *URLHandler) CreateShortURL(w http.ResponseWriter, r *http.Request){
		if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var request struct{
		LongURL string `json:"long_url"`
	}
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	url := h.service.CreateShortURL(request.LongURL)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(url)

}

func (h *URLHandler) RedirectURL(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	code := strings.TrimPrefix(r.URL.Path, "/")

	url, err := h.service.GetLongURL(code)
	if err != nil {
		http.Error(w, "URL not found", http.StatusNotFound)
		return
	}

	http.Redirect(w, r, url.LongURL, http.StatusFound)
}