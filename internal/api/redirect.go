package api

import (
	"github.com/go-chi/chi/v5"
	"net/http"
)

type urlGetter interface {
	GetURL(key string) (string, error)
}

func redirectHandler(storage urlGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Only GET requests allowed!", http.StatusMethodNotAllowed)
			return
		}

		shortKey := chi.URLParam(r, "shortKey")

		originalURL, err := storage.GetURL(shortKey)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.Header().Set("Location", originalURL)
		w.WriteHeader(http.StatusTemporaryRedirect)
	}
}
