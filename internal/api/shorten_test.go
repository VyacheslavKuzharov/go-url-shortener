package api

import (
	"bytes"
	"fmt"
	"github.com/VyacheslavKuzharov/go-url-shortener/internal/config"
	httpcfg "github.com/VyacheslavKuzharov/go-url-shortener/internal/config/http"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestShortenHandler(t *testing.T) {
	shortKey := `shortKey`
	expectedBody := fmt.Sprintf(`{"result":"http://localhost:8080/%s"}`, shortKey)
	cfgs := &config.Config{HTTP: httpcfg.HTTPCfg{Host: "localhost", Port: "8080"}}

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
			reqBody:      bytes.NewReader([]byte(`{"url":"https://practicum.yandex.ru/"}`)),
			expectedCode: http.StatusCreated,
			expectedBody: expectedBody,
			mock:         &MockStorage{saveURL: func(originalURL string) (string, error) { return shortKey, nil }},
		},
		{
			name:         "when unhappy path: empty reqBody",
			method:       http.MethodPost,
			reqBody:      bytes.NewReader([]byte("")),
			expectedCode: http.StatusBadRequest,
			expectedBody: `{"error":"request is empty"}`,
			mock:         &MockStorage{},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			r := httptest.NewRequest(tc.method, "/api/shorten", tc.reqBody)
			w := httptest.NewRecorder()

			h := shortenHandler(tc.mock, cfgs)
			h(w, r)

			res := w.Result()
			// check response code
			assert.Equal(t, tc.expectedCode, w.Code, "Код ответа не совпадает с ожидаемым")
			// check response body
			defer res.Body.Close()
			resBody, err := io.ReadAll(res.Body)

			require.NoError(t, err)
			assert.Equal(t, tc.expectedBody, strings.TrimSuffix(string(resBody), "\n"))
		})
	}
}
