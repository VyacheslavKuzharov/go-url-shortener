package app

import (
	"fmt"
	"github.com/VyacheslavKuzharov/go-url-shortener/internal/api"
	"github.com/VyacheslavKuzharov/go-url-shortener/internal/httpserver"
	"github.com/VyacheslavKuzharov/go-url-shortener/internal/storage/inmemory"
	"net/http"
)

func Run() error {
	fmt.Println("Starting go-url-shortener application...")

	// Storage
	storage := inmemory.NewMemoryStorage()
	fmt.Println("Storage initialized.")

	// Http
	mux := http.NewServeMux()
	newAPI := api.New(mux)
	newAPI.InitRoutes(storage)

	httpServer := httpserver.New(mux)
	httpServer.Start()

	return nil
}
