package services

import (
	"fmt"
	"shortify/internal/models"
	"shortify/internal/repository"
	"shortify/internal/utils"
	"sync/atomic"
)

type URLService struct {
	repo   repository.URLRepositoryInterface
	nextID int64
}

func NewURLService(repo repository.URLRepositoryInterface) *URLService {
	return &URLService{
		repo:   repo,
		nextID: 1,
	}
}

func (s *URLService) CreateShortURL(longURL string) (*models.URL, error) {
	// generate ID only for shortCode
	id := atomic.AddInt64(&s.nextID, 1)

	shortCode := utils.EncodeBase62(id)

	url := &models.URL{
		ShortCode: shortCode,
		LongURL:   longURL,
	}

	err := s.repo.Save(url)
	if err != nil {
		return nil, err
	}

	return url, nil
}

func (s *URLService) GetLongURL(code string) (*models.URL, error) {
	url, err := s.repo.Get(code)
	if err != nil {
		return nil, fmt.Errorf("URL not found")
	}
	return url, nil
}