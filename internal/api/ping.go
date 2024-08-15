package api

import (
	"context"
	"net/http"
)

type dbPinger interface {
	Ping(ctx context.Context) error
}

func pingHandler(storage dbPinger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := storage.Ping(r.Context())
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}

		w.WriteHeader(http.StatusOK)
	}
}
