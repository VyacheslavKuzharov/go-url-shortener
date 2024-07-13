package api

import "net/http"

type URLGetter interface {
	GetURL(key string) (string, bool)
}

func redirectHandler(storage URLGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Only GET requests allowed!", http.StatusMethodNotAllowed)
			return
		}

		shortKey := r.PathValue("shortKey")

		originalURL, ok := storage.GetURL(shortKey)
		if !ok {
			http.Error(w, "shortKey not found", http.StatusBadRequest)
			return
		}

		w.Header().Set("Location", originalURL)
		w.WriteHeader(http.StatusTemporaryRedirect)
	}
}
