package repository

import "shortify/internal/models"
import "sync"
import "fmt"

type URLRepository struct {
	store map[string]*models.URL
	mu    sync.RWMutex
}

func NewURLRepository() *URLRepository {
	return &URLRepository{
		store: make(map[string]*models.URL),
	}
}

func (r *URLRepository) Save(url *models.URL) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.store[url.ShortCode] = url
	return nil
}

func (r *URLRepository) Get(code string) (*models.URL, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	url, exists := r.store[code]
	if !exists {
		return nil, fmt.Errorf("not found")
	}
	return url, nil
}