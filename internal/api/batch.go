package api

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/VyacheslavKuzharov/go-url-shortener/internal/config"
	"github.com/VyacheslavKuzharov/go-url-shortener/internal/entity"
	"github.com/VyacheslavKuzharov/go-url-shortener/internal/lib/httpapi"
	"github.com/VyacheslavKuzharov/go-url-shortener/internal/lib/httpapi/response"
	"github.com/VyacheslavKuzharov/go-url-shortener/internal/lib/random"
	"io"
	"net/http"
)

const batchCapacity = 500

var invalidURLErr string

type urlRequestItem struct {
	CorrelationID string `json:"correlation_id"`
	OriginalURL   string `json:"original_url"`
}

type aliasResponseItem struct {
	CorrelationID string `json:"correlation_id"`
	ShortURL      string `json:"short_url"`
}

type batchURLSaver interface {
	SaveBatchURLs(ctx context.Context, urls *[]entity.ShortenURL) error
}

func batchHandler(storage batchURLSaver, cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Decode request body
		var urlBatch []urlRequestItem

		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&urlBatch)
		if errors.Is(err, io.EOF) {
			response.Err(w, "request is empty", http.StatusBadRequest)
			return
		}
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Make ShortenURLs slice to save in Storage
		shortenURLs := make([]entity.ShortenURL, 0, batchCapacity)

		// Make aliasesBatch slice to api response
		aliasesBatch := make([]aliasResponseItem, 0, len(urlBatch))

		for _, item := range urlBatch {
			if !httpapi.IsURLValid(item.OriginalURL) {
				invalidURLErr = fmt.Sprintf("provided url is invalid: %s", item.OriginalURL)
				break
			}

			// Build batch of ShortenURLs to save in Storage
			su := entity.ShortenURL{
				ShortKey:    random.GenShortKey(),
				OriginalURL: item.OriginalURL,
			}

			shortenURLs = append(shortenURLs, su)
			if len(shortenURLs) == batchCapacity {
				err = storage.SaveBatchURLs(r.Context(), &shortenURLs)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				shortenURLs = shortenURLs[:0]
			}

			// Build batch for api response
			alias := aliasResponseItem{
				CorrelationID: item.CorrelationID,
				ShortURL:      httpapi.FullShortenedURL(su.ShortKey, cfg),
			}

			aliasesBatch = append(aliasesBatch, alias)
		}

		// Save the remaining records
		err = storage.SaveBatchURLs(r.Context(), &shortenURLs)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if invalidURLErr == "" {
			response.OK(w, http.StatusCreated, &aliasesBatch)
		} else {
			response.Err(w, invalidURLErr, http.StatusBadRequest)
			invalidURLErr = ""
		}
	}
}
