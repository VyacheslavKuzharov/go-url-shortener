package api

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/VyacheslavKuzharov/go-url-shortener/internal/config"
	httpcfg "github.com/VyacheslavKuzharov/go-url-shortener/internal/config/http"
	"github.com/VyacheslavKuzharov/go-url-shortener/internal/entity"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestUserURLsHandler(t *testing.T) {
	cfgs := &config.Config{HTTP: httpcfg.HTTPCfg{Host: "localhost", Port: "8080"}}

	b1, urls1 := validResponse()
	b2, urls2 := emptyResponse()

	testCases := []struct {
		name         string
		method       string
		expectedCode int
		expectedBody string
		mock         *MockStorage
	}{
		{
			name:         "when happy path: correct response",
			method:       http.MethodGet,
			expectedCode: http.StatusOK,
			expectedBody: b1.String(),
			mock: &MockStorage{
				getUserUrls: func(context.Context, uuid.UUID, *config.Config) ([]*entity.CompletedURL, error) {
					return urls1, nil
				},
			},
		},
		{
			name:         "when unhappy path: empty responseBody",
			method:       http.MethodGet,
			expectedCode: http.StatusNoContent,
			expectedBody: b2.String(),
			mock: &MockStorage{
				getUserUrls: func(context.Context, uuid.UUID, *config.Config) ([]*entity.CompletedURL, error) {
					return urls2, nil
				},
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			r := httptest.NewRequest(tc.method, "/api/user/urls", nil)
			w := httptest.NewRecorder()

			ctx := context.WithValue(r.Context(), entity.CurrentUserID, uuid.NewV4())

			h := userURLsHandler(tc.mock, cfgs)
			h(w, r.WithContext(ctx))

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

func validResponse() (*bytes.Buffer, []*entity.CompletedURL) {
	userURLs := []*entity.CompletedURL{
		{OriginalURL: "https://practicum.yandex.ru", ShortURL: "https://localhost:8080/test1"},
		{OriginalURL: "https://github.com", ShortURL: "https://localhost:8080/test2"},
		{OriginalURL: "https://google.com", ShortURL: "https://localhost:8080/test3"},
	}

	buf := &bytes.Buffer{}
	enc := json.NewEncoder(buf)
	enc.Encode(userURLs) //nolint:errcheck

	return buf, userURLs
}

func emptyResponse() (*bytes.Buffer, []*entity.CompletedURL) {
	var emptyUserURLs []*entity.CompletedURL

	buf := &bytes.Buffer{}
	enc := json.NewEncoder(buf)
	enc.Encode(emptyUserURLs) //nolint:errcheck

	return buf, emptyUserURLs
}
