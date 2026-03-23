package services

import (
	"fmt"
	"shortify/internal/models"
	"shortify/internal/repository"
	"shortify/internal/utils"
	"shortify/internal/cache"
	"sync/atomic"
)

type URLService struct {
	repo   repository.URLRepositoryInterface
	cache *cache.RedisClient
	nextID int64
}

func NewURLService(repo repository.URLRepositoryInterface, cache *cache.RedisClient) *URLService {
	return &URLService{
		repo:   repo,
		cache: cache,
		nextID: 1,
	}
}

func (s *URLService) CreateShortURL(longURL string) (*models.URL, error) {
	
	existing, err := s.repo.GetByLongURL(longURL)
	if err != nil{
		return nil, err
	}

	if existing != nil{
		return existing, nil
	}
	

	
	// generate ID only for shortCode
	id := atomic.AddInt64(&s.nextID, 1)

	shortCode := utils.EncodeBase62(id)

	url := &models.URL{
		ShortCode: shortCode,
		LongURL:   longURL,
	}

	err1 := s.repo.Save(url)
	if err1 != nil {
		existing1, err2 := s.repo.GetByLongURL(longURL)
		if err2 != nil{
			return nil,err1
		}

		_ = s.cache.Set(existing.ShortCode, existing1.LongURL)
		return existing1, nil
	}
    

	_ = s.cache.Set(shortCode, longURL)
	return url, nil
}

func (s *URLService) GetLongURL(code string) (*models.URL, error) {

	//checking Redis
	longURL, err1 := s.cache.Get(code)
	if err1 == nil {
		return &models.URL{
			ShortCode: code,
			LongURL: longURL,
		}, nil
	}


	//Check DB if not stored in Redis
	url, err := s.repo.Get(code)
	if err != nil {
		return nil, fmt.Errorf("URL not found")
	}

	_ = s.cache.Set(code, url.LongURL)
	return url, nil
}