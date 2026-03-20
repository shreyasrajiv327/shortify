package services

import(
	"shortify/internal/models"
	"shortify/internal/repository"
	"shortify/internal/utils"
	"sync/atomic"
	"fmt"
)

type URLService struct{
	repo *repository.URLRepository
	nextID int64
}

func NewURLService(repo *repository.URLRepository) *URLService{
	return &URLService{
		repo: repo,
		nextID: 1,
	}
}

func (s *URLService) CreateShortURL(longURL string) *models.URL{
	id := atomic.AddInt64(&s.nextID, 1)


	shortCode := utils.EncodeBase62(id)
	url := &models.URL{
		ShortCode: shortCode,
		LongURL: longURL,
		ID: id,
	}

	s.repo.Save(url)
	return url

}


func (s *URLService) GetLongURL(code string) (*models.URL, error) {
	url, exists := s.repo.Get(code)
	if !exists {
		return nil, fmt.Errorf("URL not found")
	}
	return url, nil
}