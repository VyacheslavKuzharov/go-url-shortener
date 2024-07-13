package app

import (
	"fmt"
	"io"
	"math/rand"
	"net/http"
)

var urls = make(map[string]string)

func Run() error {
	fmt.Println("Starting go-url-shortener application...")

	mux := http.NewServeMux()
	mux.HandleFunc(`/`, shortenHandler)
	mux.HandleFunc(`/{shortKey}`, redirectHandler)

	err := http.ListenAndServe(`:8080`, mux)
	if err != nil {
		return err
	}

	return nil
}

func shortenHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST requests allowed!", http.StatusMethodNotAllowed)
		return
	}

	b, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	originalURL := string(b)
	if originalURL == "" {
		http.Error(w, "URL parameter is missing", http.StatusBadRequest)
	}

	// Generate a unique shortKey
	shortKey := generateShortKey()
	urls[shortKey] = originalURL

	// Construct the full shortened URL
	shortenedURL := fmt.Sprintf("http://localhost:8080/%s", shortKey)

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusCreated)
	if _, err = w.Write([]byte(shortenedURL)); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
}

func redirectHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Only GET requests allowed!", http.StatusMethodNotAllowed)
		return
	}

	shortKey := r.PathValue("shortKey")

	originalURL, ok := urls[shortKey]
	if !ok {
		http.Error(w, "shortKey not found", http.StatusBadRequest)
	}

	w.Header().Set("Location", originalURL)
	w.WriteHeader(http.StatusTemporaryRedirect)
}

func generateShortKey() string {
	charset := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	keyLength := 6

	shortKey := make([]byte, keyLength)
	for i := range shortKey {
		shortKey[i] = charset[rand.Intn(len(charset))]
	}
	return string(shortKey)
}
