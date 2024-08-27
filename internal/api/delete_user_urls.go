package api

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/VyacheslavKuzharov/go-url-shortener/internal/entity"
	"github.com/VyacheslavKuzharov/go-url-shortener/internal/lib/httpapi/response"
	uuid "github.com/satori/go.uuid"
	"io"
	"net/http"
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

		err = storage.DeleteUserUrls(r.Context(), currentUserID, urlKeysBatch)
		if err != nil {
			response.Err(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusAccepted)
	}
}
