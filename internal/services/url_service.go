package services

import(
	"shortify/internal/models"
	"shortify/internal/repository"
	"strconv"
	"fmt"
)

type URLService struct{
	repo *repository.URLRepository
	nextID int
}

func NewURLService(repo *repository.URLRepository) *URLService{
	return &URLService{
		repo: repo,
		nextID: 1,
	}
}

func (s *URLService) CreateShortURL(longURL string) *models.URL{
	id := s.nextID
	s.nextID++

	shortCode := strconv.Itoa(id)
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