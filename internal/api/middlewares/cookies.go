package middlewares

import (
	"bytes"
	"context"
	"encoding/gob"
	"github.com/VyacheslavKuzharov/go-url-shortener/internal/api/middlewares/cookies"
	"github.com/VyacheslavKuzharov/go-url-shortener/internal/entity"
	"github.com/VyacheslavKuzharov/go-url-shortener/internal/logger"
	uuid "github.com/satori/go.uuid"
	"net/http"
	"strings"
)

func Cookies(l *logger.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			var currentUserID uuid.UUID
			var err error

			currentUserID, err = getCookieHandler(r, l)
			if err != nil {
				newUUID := uuid.NewV4()
				user := &entity.User{UUID: newUUID}

				setCookieHandler(w, user, l)
				currentUserID = newUUID
			}

			ctx := context.WithValue(r.Context(), entity.CurrentUserID, currentUserID)
			next.ServeHTTP(w, r.WithContext(ctx))
		}

		return http.HandlerFunc(fn)
	}
}

func setCookieHandler(w http.ResponseWriter, user *entity.User, l *logger.Logger) {
	var buf bytes.Buffer

	err := gob.NewEncoder(&buf).Encode(user)
	if err != nil {
		l.Info("setCookieHandler.gob.NewEncoder error: %v", err)
	}

	cookie := http.Cookie{
		Name:     cookies.UserData,
		Value:    buf.String(),
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
	}

	err = cookies.WriteEncrypted(w, cookie)
	if err != nil {
		l.Info("setCookieHandler.cookies.WriteEncrypted error: %v", err)
	}
}

func getCookieHandler(r *http.Request, l *logger.Logger) (uuid.UUID, error) {
	var user entity.User

	val, err := cookies.ReadEncrypted(r, cookies.UserData)
	if err != nil {
		return uuid.Nil, err
	}

	reader := strings.NewReader(val)

	if err = gob.NewDecoder(reader).Decode(&user); err != nil {
		l.Info("getCookieHandler.gob.NewDecoder error: %v", err)
		return uuid.Nil, err
	}

	return user.UUID, nil
}
