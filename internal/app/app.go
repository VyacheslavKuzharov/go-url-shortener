package app

import (
	"github.com/VyacheslavKuzharov/go-url-shortener/internal/api"
	"github.com/VyacheslavKuzharov/go-url-shortener/internal/config"
	"github.com/VyacheslavKuzharov/go-url-shortener/internal/httpserver"
	"github.com/VyacheslavKuzharov/go-url-shortener/internal/logger"
	"github.com/VyacheslavKuzharov/go-url-shortener/internal/storage"
	"github.com/go-chi/chi/v5"
	"log"
	"net"
)

func Run(cfg *config.Config) {
	l := logger.New(cfg.Log.Level)
	l.Info("Starting go-url-shortener application...")

	// Storage
	s, err := storage.New(cfg)
	if err != nil {
		log.Fatal(err)
	}
	l.Info("Storage initialized.")
	defer s.Close()

	// Http
	router := chi.NewRouter()
	api.New(router, cfg, s, l)

	httpServer := httpserver.New(router)
	l.Info("Starting server on: %s", net.JoinHostPort(cfg.HTTP.Host, cfg.HTTP.Port))
	httpServer.Start(cfg.HTTP.Host, cfg.HTTP.Port)
}
