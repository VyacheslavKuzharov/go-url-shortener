package api

import (
	"github.com/go-chi/chi/v5"
	"net/http"
)

type urlGetter interface {
	GetURL(key string) (string, bool)
}

func redirectHandler(storage urlGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Only GET requests allowed!", http.StatusMethodNotAllowed)
			return
		}

		shortKey := chi.URLParam(r, "shortKey")

		originalURL, ok := storage.GetURL(shortKey)
		if !ok {
			http.Error(w, "shortKey not found", http.StatusBadRequest)
			return
		}

		w.Header().Set("Location", originalURL)
		w.WriteHeader(http.StatusTemporaryRedirect)
	}
}
