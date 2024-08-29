package api

import (
	"context"
	"github.com/VyacheslavKuzharov/go-url-shortener/internal/config"
	"github.com/VyacheslavKuzharov/go-url-shortener/internal/entity"
	"github.com/VyacheslavKuzharov/go-url-shortener/internal/lib/httpapi/response"
	uuid "github.com/satori/go.uuid"
	"net/http"
)

type userURLsGetter interface {
	GetUserUrls(ctx context.Context, currentUserID uuid.UUID, cfg *config.Config) ([]*entity.CompletedURL, error)
}

func userURLsHandler(storage userURLsGetter, cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		currentUserID := r.Context().Value(entity.CurrentUserID).(uuid.UUID)

		userURLs, err := storage.GetUserUrls(r.Context(), currentUserID, cfg)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if len(userURLs) == 0 {
			response.OK(w, http.StatusNoContent, userURLs)
			return
		}

		response.OK(w, http.StatusOK, userURLs)
	}
}
