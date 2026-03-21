package repository

import "shortify/internal/models"

type URLRepositoryInterface interface {
	Save(url *models.URL) error
	Get(code string) (*models.URL, error)
}