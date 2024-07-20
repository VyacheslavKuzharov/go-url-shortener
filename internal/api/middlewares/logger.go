package middlewares

import (
	"github.com/VyacheslavKuzharov/go-url-shortener/internal/logger"
	"github.com/go-chi/chi/v5/middleware"
	"net/http"
	"time"
)

func Logger(l *logger.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		l.Info("logger middleware started")

		logFn := func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

			defer func() {
				l.LogHTTPRequest(ww, r, start)
			}()

			next.ServeHTTP(ww, r)
		}

		return http.HandlerFunc(logFn)
	}
}
