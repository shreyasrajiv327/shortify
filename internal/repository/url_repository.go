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

func (r *URLRepository) Save(url *models.URL) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.store[url.ShortCode]; exists {
		return ErrShortCodeExists
	}
	for _, existing := range r.store {
		if existing.LongURL == url.LongURL {
			return ErrLongURLExists
		}
	}

	r.store[url.ShortCode] = url
	return nil
}

func (r *URLRepository) Get(code string) (*models.URL, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	url, exists := r.store[code]
	if !exists {
		return nil, ErrNotFound
	}
	return url, nil
}

func (r *URLRepository) GetByLongURL(longURL string) (*models.URL, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, url := range r.store {
		if url.LongURL == longURL {
			return url, nil
		}
	}
	return nil, ErrNotFound
}
