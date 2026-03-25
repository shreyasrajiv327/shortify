package services

import (
	"fmt"
	"sync/atomic"

	"shortify/internal/cache"
	"shortify/internal/logger"
	"shortify/internal/models"
	"shortify/internal/repository"
	"shortify/internal/utils"
)

type URLService struct {
	repo   repository.URLRepositoryInterface
	cache  *cache.RedisClient
	nextID int64
}

func NewURLService(repo repository.URLRepositoryInterface, cache *cache.RedisClient) *URLService {
	return &URLService{
		repo:   repo,
		cache:  cache,
		nextID: 1,
	}
}

// CreateShortURL handles URL shortening logic
func (s *URLService) CreateShortURL(longURL string) (*models.URL, error) {

	logger.Log.Info("CreateShortURL called", "long_url", longURL)

	// 1. Check if URL already exists
	existing, err := s.repo.GetByLongURL(longURL)
	if err != nil {
		logger.Log.Error("DB error while checking existing URL", "error", err)
		return nil, err
	}

	if existing != nil {
		logger.Log.Info("Duplicate URL found",
			"long_url", longURL,
			"short_code", existing.ShortCode,
		)
		return existing, nil
	}

	// 2. Generate short code
	id := atomic.AddInt64(&s.nextID, 1)
	shortCode := utils.EncodeBase62(id)

	logger.Log.Info("Generated short code", "short_code", shortCode)

	url := &models.URL{
		ShortCode: shortCode,
		LongURL:   longURL,
	}

	// 3. Save to DB
	err1 := s.repo.Save(url)
	if err1 != nil {
		logger.Log.Warn("Insert failed, attempting fallback fetch", "error", err1)

		// Fallback: fetch existing record (due to race condition)
		existing1, err2 := s.repo.GetByLongURL(longURL)
		if err2 != nil {
			logger.Log.Error("Fallback fetch failed", "error", err2)
			return nil, err1
		}

		logger.Log.Info("Returning existing URL after conflict",
			"short_code", existing1.ShortCode,
		)

		// Cache it
		_ = s.cache.Set(existing1.ShortCode, existing1.LongURL)

		return existing1, nil
	}

	// 4. Cache the new URL
	logger.Log.Info("URL saved successfully",
		"short_code", shortCode,
	)

	_ = s.cache.Set(shortCode, longURL)

	return url, nil
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
		logger.Log.Error("URL not found in DB", "short_code", code)
		return nil, fmt.Errorf("URL not found")
	}

	// 3. Update cache
	logger.Log.Info("Fetched from DB and updating cache",
		"short_code", code,
	)

	_ = s.cache.Set(code, url.LongURL)

	return url, nil
}