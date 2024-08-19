package middlewares

import (
	"context"
	"encoding/gob"
	"github.com/VyacheslavKuzharov/go-url-shortener/internal/api/middlewares/cookies"
	"github.com/VyacheslavKuzharov/go-url-shortener/internal/entity"
	"github.com/VyacheslavKuzharov/go-url-shortener/internal/lib/httpapi/response"
	"net/http"
	"strings"
)

func Auth(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		var currentUser entity.User

		val, err := cookies.ReadEncrypted(r, cookies.UserData)
		if err != nil {
			response.Err(w, err.Error(), http.StatusUnauthorized)
			return
		}

		reader := strings.NewReader(val)
		if err = gob.NewDecoder(reader).Decode(&currentUser); err != nil {
			response.Err(w, err.Error(), http.StatusUnauthorized)
			return
		}

		currentUserID := currentUser.UUID

		if currentUserID.String() == "" {
			response.Err(w, "Invalid Current User ID", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), entity.CurrentUserID, currentUserID)
		next.ServeHTTP(w, r.WithContext(ctx))
	}

	return http.HandlerFunc(fn)
}
