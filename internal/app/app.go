package app

import (
	"github.com/VyacheslavKuzharov/go-url-shortener/internal/api"
	"github.com/VyacheslavKuzharov/go-url-shortener/internal/config"
	"github.com/VyacheslavKuzharov/go-url-shortener/internal/httpserver"
	"github.com/VyacheslavKuzharov/go-url-shortener/internal/storage"
	"github.com/go-chi/chi/v5"
	"log"
)

func Run(cfg *config.Config) {
	log.Println("Starting go-url-shortener application...")

	// Storage
	s := storage.New()
	log.Println("Storage initialized.")

	// Http
	router := chi.NewRouter()
	api.New(router, cfg, s)

	httpServer := httpserver.New(router)
	httpServer.Start(cfg.HTTP.Host, cfg.HTTP.Port)
}
