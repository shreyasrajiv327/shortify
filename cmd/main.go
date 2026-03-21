package main

import (
	"log"
	"net/http"

	"shortify/internal/handlers"
	"shortify/internal/repository"
	"shortify/internal/services"
	"shortify/internal/database"
)

func main(){

	db := database.ConnectDB()
	repository := repository.NewPostgresRepository(db)
	service := services.NewURLService(repository)
	handler := handlers.NewURLHandler(service)
	http.HandleFunc("/shorten", handler.CreateShortURL)
	http.HandleFunc("/", handler.RedirectURL)

	log.Println("Server running on port 8080...")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}

}