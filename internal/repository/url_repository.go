package repository

import "shortify/internal/models"
import "sync"

type URLRepository struct {
	store map[string]*models.URL
	mu    sync.RWMutex
}

func NewURLRepository() *URLRepository {
	return &URLRepository{
		store: make(map[string]*models.URL),
	}
}

func (r *URLRepository) Save(url *models.URL) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.store[url.ShortCode] = url
}

func (r *URLRepository) Get(code string) (url *models.URL, exists bool){
	r.mu.RLock()
	defer r.mu.RUnlock()
	
	url , exists = r.store[code]
	return url, exists
}