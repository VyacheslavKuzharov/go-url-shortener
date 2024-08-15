package api

import (
	"bytes"
	"context"
	"errors"
	"github.com/VyacheslavKuzharov/go-url-shortener/internal/config"
	httpcfg "github.com/VyacheslavKuzharov/go-url-shortener/internal/config/http"
	"github.com/VyacheslavKuzharov/go-url-shortener/internal/entity"
	"github.com/VyacheslavKuzharov/go-url-shortener/internal/storage/mock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestBatchHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	s := mock.NewMockStorager(ctrl)
	cfgs := &config.Config{HTTP: httpcfg.HTTPCfg{Host: "localhost", Port: "8080"}}

	validRequest := `
		[
				{
						"correlation_id": "4DATsj",
						"original_url": "https://practicum.yandex.ru"
				},
				{
						"correlation_id": "3JSTsL",
						"original_url": "https://github.com"
				},
				{
						"correlation_id": "2LDJzP",
						"original_url": "https://google.com"
				}
		]
	`
	inValidRequest := `
		[
				{
						"correlation_id": "4DATsj",
						"original_url": "https://practicum.yandex.ru"
				},
				{
						"correlation_id": "3JSTsL",
						"original_url": "github.com"
				},
				{
						"correlation_id": "2LDJzP",
						"original_url": "https://google.com"
				}
		]
	`

	shortenURLs := []entity.ShortenURL{
		{OriginalURL: "https://practicum.yandex.ru", ShortKey: "test1"},
		{OriginalURL: "https://github.com", ShortKey: "test2"},
		{OriginalURL: "https://google.com", ShortKey: "test3"},
	}

	testCases := []struct {
		name         string
		method       string
		reqBody      io.Reader
		expectedCode int
		mock         *gomock.Call
	}{
		{
			name:         "when happy path: correct response",
			method:       http.MethodPost,
			reqBody:      bytes.NewReader([]byte(validRequest)),
			expectedCode: http.StatusCreated,
			mock:         s.EXPECT().SaveBatchURLs(context.Background(), gomock.AssignableToTypeOf(shortenURLs)).Return(nil),
		},
		{
			name:         "when unhappy path: SaveBatchURLs returns error",
			method:       http.MethodPost,
			reqBody:      bytes.NewReader([]byte(validRequest)),
			expectedCode: http.StatusInternalServerError,
			mock:         s.EXPECT().SaveBatchURLs(context.Background(), gomock.AssignableToTypeOf(shortenURLs)).Return(errors.New("test")),
		},
		{
			name:         "when unhappy path: empty reqBody",
			method:       http.MethodPost,
			reqBody:      bytes.NewReader([]byte("")),
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "when unhappy path: invalid original url",
			method:       http.MethodPost,
			reqBody:      bytes.NewReader([]byte(inValidRequest)),
			expectedCode: http.StatusBadRequest,
			mock:         s.EXPECT().SaveBatchURLs(context.Background(), gomock.AssignableToTypeOf(shortenURLs)).Return(nil),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			r := httptest.NewRequest(tc.method, "/api/shorten/batch", tc.reqBody)
			w := httptest.NewRecorder()

			h := batchHandler(s, cfgs)
			h(w, r)

			res := w.Result()
			// check response code
			assert.Equal(t, tc.expectedCode, w.Code, "Код ответа не совпадает с ожидаемым")
			// check response body
			defer res.Body.Close()
			_, err := io.ReadAll(res.Body)

			require.NoError(t, err)
		})
	}
}
