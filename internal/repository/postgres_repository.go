package repository

import (
	"database/sql"
	"shortify/internal/models"
	"fmt"
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

	return r.db.QueryRow(query, url.ShortCode, url.LongURL).Scan(&url.ID)
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

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("not found")
	}
	if err != nil {
		return nil, err
	}

	return &url, nil
}