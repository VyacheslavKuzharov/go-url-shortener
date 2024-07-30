package response

import (
	"bytes"
	"encoding/json"
	"net/http"
)

type Response struct {
	Error string `json:"error,omitempty"`
}

func OK(w http.ResponseWriter, code int, payload any) {
	buf := &bytes.Buffer{}
	enc := json.NewEncoder(buf)

	if err := enc.Encode(payload); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(buf.Bytes()) //nolint:errcheck
}

func Err(w http.ResponseWriter, error string, code int) {
	buf := &bytes.Buffer{}
	enc := json.NewEncoder(buf)

	resp := Response{
		Error: error,
	}

	if err := enc.Encode(resp); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(buf.Bytes()) //nolint:errcheck
}
