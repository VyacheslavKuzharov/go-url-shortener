package api

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/VyacheslavKuzharov/go-url-shortener/internal/entity"
	"github.com/VyacheslavKuzharov/go-url-shortener/internal/lib/httpapi/response"
	uuid "github.com/satori/go.uuid"
	"io"
	"log"
	"net/http"
	"sync"
)

type userURLsDeleter interface {
	DeleteUserUrls(ctx context.Context, currentUserID uuid.UUID, urlKeysBatch []string) error
}

func deleteUserURLsHandler(storage userURLsDeleter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		currentUserID := r.Context().Value(entity.CurrentUserID).(uuid.UUID)

		var urlKeysBatch []string

		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&urlKeysBatch)
		if errors.Is(err, io.EOF) {
			response.Err(w, "request is empty", http.StatusBadRequest)
			return
		}
		if err != nil {
			response.Err(w, err.Error(), http.StatusBadRequest)
			return
		}

		batchSize := 5
		batches := genSliceBatches(urlKeysBatch, batchSize)
		channel := make(chan []string, len(batches))

		for _, shortKeysBatch := range batches {
			channel <- shortKeysBatch
		}

		var wg sync.WaitGroup
		workers := 5

		wg.Add(workers)
		for i := 0; i < workers; i++ {
			go func(ch chan []string, wg *sync.WaitGroup) {
				defer wg.Done()

				for shortKeysBatch := range ch {
					err = storage.DeleteUserUrls(r.Context(), currentUserID, shortKeysBatch)
					if err != nil {
						log.Printf("Delete error: %s", err)
					}
				}

			}(channel, &wg)
		}

		close(channel)
		wg.Wait()

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusAccepted)
	}
}

func genSliceBatches(slice []string, batchSize int) [][]string {
	var batches [][]string

	for i := 0; i < len(slice); i += batchSize {
		end := i + batchSize

		if end > len(slice) {
			end = len(slice)
		}

		batches = append(batches, slice[i:end])
	}

	return batches
}
