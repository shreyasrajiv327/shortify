package main

import (
	
	"net/http"

	"shortify/internal/handlers"
	"shortify/internal/repository"
	"shortify/internal/services"
	"shortify/internal/database"
	"shortify/internal/cache"
	"shortify/internal/logger"
	"shortify/internal/middleware"
)

func main(){
   
	logger.Init()
	db := database.ConnectDB()
	cache := cache.NewRedisClient()
	repository := repository.NewPostgresRepository(db)
	service := services.NewURLService(repository, cache)
	handler := handlers.NewURLHandler(service)

    mux := http.NewServeMux()
    mux.HandleFunc("/shorten", handler.CreateShortURL)
    mux.HandleFunc("/", handler.RedirectURL)


    loggedMux := middleware.LoggingMiddleware(mux)


    logger.Log.Info("Server starting on :8080")
    http.ListenAndServe(":8080", loggedMux)

}