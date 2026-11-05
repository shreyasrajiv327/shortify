package repository

import (
	"database/sql"
	"errors"
	"shortify/internal/models"

	"github.com/lib/pq"
)

type PostgresRepository struct {
	db *sql.DB
}

func NewPostgresRepository(db *sql.DB) *PostgresRepository {
	return &PostgresRepository{db: db}
}

func (r *PostgresRepository) Save(url *models.URL) error {
	query := `
		INSERT INTO urls (short_code, long_url)
		VALUES ($1, $2)
		RETURNING id
	`

	err := r.db.QueryRow(query, url.ShortCode, url.LongURL).Scan(&url.ID)
	if err == nil {
		return nil
	}

	var postgresErr *pq.Error
	if errors.As(err, &postgresErr) && postgresErr.Code == "23505" {
		switch postgresErr.Constraint {
		case "urls_short_code_key":
			return ErrShortCodeExists
		case "unique_long_url":
			return ErrLongURLExists
		}
	}

	return err
}

func (r *PostgresRepository) Get(code string) (*models.URL, error) {
	query := `
		SELECT id, short_code, long_url
		FROM urls
		WHERE short_code = $1
	`

	var url models.URL
	err := r.db.QueryRow(query, code).
		Scan(&url.ID, &url.ShortCode, &url.LongURL)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	return &url, nil
}

func (r *PostgresRepository) GetByLongURL(longURL string) (*models.URL, error) {
	query := `
	  SELECT id, short_code, long_url
	  FROM urls
	  WHERE long_url = $1
	`

	var url models.URL

	err := r.db.QueryRow(query, longURL).Scan(&url.ID, &url.ShortCode, &url.LongURL)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}

	if err != nil {
		return nil, err
	}

	return &url, nil
}
