package app

import (
	"fmt"
	"github.com/VyacheslavKuzharov/go-url-shortener/internal/api"
	"github.com/VyacheslavKuzharov/go-url-shortener/internal/httpserver"
	"github.com/VyacheslavKuzharov/go-url-shortener/internal/storage"
	"github.com/go-chi/chi/v5"
)

func Run() error {
	fmt.Println("Starting go-url-shortener application...")

	// Storage
	s := storage.New()
	fmt.Println("Storage initialized.")

	// Http
	router := chi.NewRouter()
	newAPI := api.New(router)
	routes := newAPI.InitRoutes(s)

	httpServer := httpserver.New(routes)
	httpServer.Start()

	return nil
}
