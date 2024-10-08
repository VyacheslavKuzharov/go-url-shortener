package api

import (
	"context"
	"errors"
	"github.com/VyacheslavKuzharov/go-url-shortener/internal/config"
	"github.com/VyacheslavKuzharov/go-url-shortener/internal/lib/httpapi"
	"github.com/VyacheslavKuzharov/go-url-shortener/internal/storage/postgres"
	"io"
	"net/http"
)

type urlSaver interface {
	SaveURL(ctx context.Context, originalURL string) (string, error)
}

func saveURLHandler(storage urlSaver, cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var pgUniqueFieldErr *postgres.UniqueFieldErr

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

		shortKey, err := storage.SaveURL(r.Context(), originalURL)
		if err != nil {
			if errors.As(err, &pgUniqueFieldErr) {
				su := httpapi.FullShortenedURL(pgUniqueFieldErr.Payload, cfg)

				w.Header().Set("Content-Type", "text/plain")
				w.WriteHeader(http.StatusConflict)
				w.Write([]byte(su)) //nolint:errcheck
				return
			}

			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		shortenedURL := httpapi.FullShortenedURL(shortKey, cfg)

		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(shortenedURL)) //nolint:errcheck
	}
}
