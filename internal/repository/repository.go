package repository

import (
	"errors"

	"shortify/internal/models"
)

var (
	ErrNotFound        = errors.New("not found")
	ErrShortCodeExists = errors.New("short code already exists")
	ErrLongURLExists   = errors.New("long URL already exists")
)

type URLRepositoryInterface interface {
	Save(url *models.URL) error
	Get(code string) (*models.URL, error)
	GetByLongURL(longURL string) (*models.URL, error)
}
