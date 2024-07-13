package api

import (
	"fmt"
	"io"
	"net/http"
)

type URLSaver interface {
	SaveURL(originalURL string) (string, error)
}

func saveURLHandler(storage URLSaver) http.HandlerFunc {
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

		shortKey, err := storage.SaveURL(originalURL)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Construct the full shortened URL
		shortenedURL := fmt.Sprintf("http://localhost:8080/%s", shortKey)

		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusCreated)
		if _, err = w.Write([]byte(shortenedURL)); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
	}
}
