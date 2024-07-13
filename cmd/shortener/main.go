package main

import (
	"github.com/VyacheslavKuzharov/go-url-shortener/internal/app"
	"log"
)

func main() {
	if err := app.Run(); err != nil {
		log.Fatalf("app error: %s", err)
	}
}
