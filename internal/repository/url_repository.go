package repository

import "shortify/internal/models"

type URLRepository struct {
	store map[string]*models.URL
}

func NewURLRepository() *URLRepository {
	return &URLRepository{
		store: make(map[string]*models.URL),
	}
}

func (r *URLRepository) Save(url *models.URL) {
	r.store[url.ShortCode] = url
}

func (r *URLRepository) Get(code string) (url *models.URL, exists bool){
	url , exists = r.store[code]
	return url, exists
}