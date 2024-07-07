package main

import (
	"github.com/VyacheslavKuzharov/go-url-shortener/internal/app"
	"github.com/VyacheslavKuzharov/go-url-shortener/internal/config"
	"log"
)

func main() {
	cfg, err := config.New()
	if err != nil {
		log.Fatalf("config error: %s", err)
	}

	app.Run(cfg)
}
