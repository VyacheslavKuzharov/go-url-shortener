package api

import (
	"net/http"
)

type dbPinger interface {
	Ping() error
}

func pingHandler(storage dbPinger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := storage.Ping()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}

		w.WriteHeader(http.StatusOK)
	}
}
