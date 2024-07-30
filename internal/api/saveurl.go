package api

import (
	"github.com/VyacheslavKuzharov/go-url-shortener/internal/config"
	"github.com/VyacheslavKuzharov/go-url-shortener/internal/lib/httpapi"
	"io"
	"net/http"
)

type urlSaver interface {
	SaveURL(originalURL string) (string, error)
}

func saveURLHandler(storage urlSaver, cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Only POST requests allowed!", http.StatusMethodNotAllowed)
			return
		}

		b, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		originalURL := string(b)
		if originalURL == "" {
			http.Error(w, "URL parameter is missing", http.StatusBadRequest)
			return
		}

		if !httpapi.IsURLValid(originalURL) {
			http.Error(w, "provided url is invalid", http.StatusBadRequest)
			return
		}

		shortKey, err := storage.SaveURL(originalURL)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		shortenedURL := httpapi.FullShortenedURL(shortKey, cfg)

		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusCreated)
		if _, err = w.Write([]byte(shortenedURL)); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
	}
}
