package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
	"github.com/joho/godotenv"
)

func ConnectDB() *sql.DB {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Read env variables
	user := os.Getenv("DB_USER")
	dbname := os.Getenv("DB_NAME")
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	sslmode := os.Getenv("DB_SSLMODE")

	// Validate
	if user == "" || dbname == "" {
		log.Fatal("Missing required DB environment variables")
	}

	// Build connection string
	connStr := fmt.Sprintf(
		"user=%s dbname=%s host=%s port=%s sslmode=%s",
		user, dbname, host, port, sslmode,
	)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Error opening database:", err)
	}

	if err := db.Ping(); err != nil {
		log.Fatal("Error connecting to database:", err)
	}

	log.Println("Connected to PostgreSQL!")
	return db
}