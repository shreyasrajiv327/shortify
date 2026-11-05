package services

import (
	"errors"
	"testing"

	"shortify/internal/logger"
	"shortify/internal/repository"
)

type testURLCache struct {
	values map[string]string
}

func (c *testURLCache) Set(key string, value string) error {
	c.values[key] = value
	return nil
}

func (c *testURLCache) Get(key string) (string, error) {
	value, exists := c.values[key]
	if !exists {
		return "", errors.New("cache miss")
	}
	return value, nil
}

func TestCreateShortURLAndLookup(t *testing.T) {
	logger.Init()
	cache := &testURLCache{values: make(map[string]string)}
	service := NewURLService(repository.NewURLRepository(), cache)
	longURL := "https://example.com/products/shirts?color=white"

	created, err := service.CreateShortURL(longURL)
	if err != nil {
		t.Fatalf("CreateShortURL() error = %v", err)
	}
	if len(created.ShortCode) != 8 {
		t.Fatalf("short code length = %d, want 8", len(created.ShortCode))
	}

	duplicate, err := service.CreateShortURL(longURL)
	if err != nil {
		t.Fatalf("CreateShortURL() for duplicate error = %v", err)
	}
	if duplicate.ShortCode != created.ShortCode {
		t.Fatalf("duplicate short code = %q, want %q", duplicate.ShortCode, created.ShortCode)
	}

	resolved, err := service.GetLongURL(created.ShortCode)
	if err != nil {
		t.Fatalf("GetLongURL() error = %v", err)
	}
	if resolved.LongURL != longURL {
		t.Fatalf("resolved URL = %q, want %q", resolved.LongURL, longURL)
	}
}

func TestGetLongURLNotFound(t *testing.T) {
	logger.Init()
	service := NewURLService(repository.NewURLRepository(), &testURLCache{values: make(map[string]string)})

	_, err := service.GetLongURL("missing")
	if !errors.Is(err, repository.ErrNotFound) {
		t.Fatalf("GetLongURL() error = %v, want %v", err, repository.ErrNotFound)
	}
}
