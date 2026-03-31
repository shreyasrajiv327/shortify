package main

import (
	"net/http"

	"shortify/internal/cache"
	"shortify/internal/database"
	"shortify/internal/handlers"
	"shortify/internal/logger"
	"shortify/internal/middleware"
	"shortify/internal/repository"
	"shortify/internal/services"
)

func main() {

	logger.Init()

	db := database.ConnectDB()
	redisClient := cache.NewRedisClient() // ✅ FIXED

	repo := repository.NewPostgresRepository(db)
	service := services.NewURLService(repo, redisClient)
	handler := handlers.NewURLHandler(service)

	mux := http.NewServeMux()
	mux.HandleFunc("/shorten", handler.CreateShortURL)
	mux.HandleFunc("/", handler.RedirectURL)

	loggedMux := middleware.LoggingMiddleware(mux)
	rateLimitedMux := middleware.RateLimitMiddleware(redisClient, 10)(loggedMux)

	logger.Log.Info("Server starting on :8080")

	http.ListenAndServe(":8080", rateLimitedMux)
}