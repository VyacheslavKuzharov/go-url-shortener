package api

import (
	"encoding/json"
	"errors"
	"github.com/VyacheslavKuzharov/go-url-shortener/internal/config"
	"github.com/VyacheslavKuzharov/go-url-shortener/internal/lib/httpapi"
	"github.com/VyacheslavKuzharov/go-url-shortener/internal/lib/httpapi/response"
	"github.com/VyacheslavKuzharov/go-url-shortener/internal/storage/postgres"
	"io"
	"net/http"
)

type Request struct {
	URL string `json:"url"`
}

type Response struct {
	Result string `json:"result,omitempty"`
}

func shortenHandler(storage urlSaver, cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req Request
		var pgUniqueFieldErr *postgres.UniqueFieldErr

		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&req)
		if errors.Is(err, io.EOF) {
			response.Err(w, "request is empty", http.StatusBadRequest)
			return
		}
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if !httpapi.IsURLValid(req.URL) {
			response.Err(w, "provided url is invalid", http.StatusBadRequest)
			return
		}

		shortKey, err := storage.SaveURL(r.Context(), req.URL)
		if err != nil {
			if errors.As(err, &pgUniqueFieldErr) {
				response.OK(w, http.StatusConflict, Response{
					Result: httpapi.FullShortenedURL(pgUniqueFieldErr.Payload, cfg),
				})

				return
			}

			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		shortenedURL := httpapi.FullShortenedURL(shortKey, cfg)

		response.OK(w, http.StatusCreated, Response{
			Result: shortenedURL,
		})
	}
}
