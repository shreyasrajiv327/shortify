package services

import (
	"errors"
	"fmt"

	"shortify/internal/logger"
	"shortify/internal/models"
	"shortify/internal/repository"
	"shortify/internal/utils"
)

type URLCache interface {
	Set(key string, value string) error
	Get(key string) (string, error)
}

type URLService struct {
	repo  repository.URLRepositoryInterface
	cache URLCache
}

func NewURLService(repo repository.URLRepositoryInterface, cache URLCache) *URLService {
	return &URLService{repo: repo, cache: cache}
}

// CreateShortURL handles URL shortening logic
func (s *URLService) CreateShortURL(longURL string) (*models.URL, error) {

	logger.Log.Info("CreateShortURL called")

	// 1. Check if URL already exists
	existing, err := s.repo.GetByLongURL(longURL)
	if err == nil && existing != nil {
		logger.Log.Info("Duplicate URL found", "short_code", existing.ShortCode)
		return existing, nil
	}
	if err != nil && !errors.Is(err, repository.ErrNotFound) {
		logger.Log.Error("DB error while checking existing URL", "error", err)
		return nil, err
	}

	for attempt := 0; attempt < 5; attempt++ {
		shortCode, err := utils.GenerateShortCode(8)
		if err != nil {
			return nil, fmt.Errorf("generate short code: %w", err)
		}

		url := &models.URL{ShortCode: shortCode, LongURL: longURL}
		err = s.repo.Save(url)
		if err == nil {
			_ = s.cache.Set(shortCode, longURL)
			return url, nil
		}

		if errors.Is(err, repository.ErrShortCodeExists) {
			continue
		}
		if errors.Is(err, repository.ErrLongURLExists) {
			existing, lookupErr := s.repo.GetByLongURL(longURL)
			if lookupErr != nil {
				return nil, lookupErr
			}
			_ = s.cache.Set(existing.ShortCode, existing.LongURL)
			return existing, nil
		}
		return nil, err
	}

	return nil, fmt.Errorf("could not generate a unique short code")
}

// GetLongURL retrieves the original URL from short code
func (s *URLService) GetLongURL(code string) (*models.URL, error) {

	logger.Log.Info("GetLongURL called", "short_code", code)

	// 1. Check Redis cache
	longURL, err := s.cache.Get(code)
	if err == nil {
		logger.Log.Info("Cache HIT", "short_code", code)

		return &models.URL{
			ShortCode: code,
			LongURL:   longURL,
		}, nil
	}

	logger.Log.Info("Cache MISS", "short_code", code)

	// 2. Fetch from DB
	url, err := s.repo.Get(code)
	if err != nil {
		return nil, err
	}

	// 3. Update cache
	logger.Log.Info("Fetched from DB and updating cache",
		"short_code", code,
	)

	_ = s.cache.Set(code, url.LongURL)

	return url, nil
}
