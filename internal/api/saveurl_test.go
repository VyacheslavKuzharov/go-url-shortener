package api

import (
	"bytes"
	"fmt"
	"github.com/VyacheslavKuzharov/go-url-shortener/internal/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSaveUrlHandler(t *testing.T) {
	shortKey := `qwerty`
	expectedBody := fmt.Sprintf("http://localhost:8080/%s", shortKey)
	originalURL := "https://practicum.yandex.ru/"
	cfgs := &config.Config{HTTP: config.HTTPCfg{Host: "localhost", Port: "8080"}}

	testCases := []struct {
		name         string
		method       string
		reqBody      io.Reader
		expectedCode int
		expectedBody string
		mock         *MockStorage
	}{
		{
			name:         "when happy path: correct response",
			method:       http.MethodPost,
			reqBody:      bytes.NewReader([]byte(originalURL)),
			expectedCode: http.StatusCreated,
			expectedBody: expectedBody,
			mock:         &MockStorage{saveURL: func(originalURL string) (string, error) { return shortKey, nil }},
		},
		{
			name:         "when unhappy path: incorrect request method",
			method:       http.MethodGet,
			reqBody:      bytes.NewReader([]byte(originalURL)),
			expectedCode: http.StatusMethodNotAllowed,
			expectedBody: "Only POST requests allowed!\n",
			mock:         &MockStorage{saveURL: func(originalURL string) (string, error) { return shortKey, nil }},
		},
		{
			name:         "when unhappy path: empty reqBody",
			method:       http.MethodPost,
			reqBody:      bytes.NewReader([]byte("")),
			expectedCode: http.StatusBadRequest,
			expectedBody: "URL parameter is missing\n",
			mock:         &MockStorage{saveURL: func(originalURL string) (string, error) { return shortKey, nil }},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			r := httptest.NewRequest(tc.method, "/", tc.reqBody)
			w := httptest.NewRecorder()

			h := saveURLHandler(tc.mock, cfgs)
			h(w, r)

			res := w.Result()
			// check response code
			assert.Equal(t, tc.expectedCode, w.Code, "Код ответа не совпадает с ожидаемым")
			// check response body
			defer res.Body.Close()
			resBody, err := io.ReadAll(res.Body)

			require.NoError(t, err)
			assert.Equal(t, tc.expectedBody, string(resBody))
		})
	}
}
