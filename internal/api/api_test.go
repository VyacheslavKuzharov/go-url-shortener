package api

import (
	"bytes"
	"fmt"
	"github.com/VyacheslavKuzharov/go-url-shortener/internal/config"
	"github.com/VyacheslavKuzharov/go-url-shortener/internal/logger"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

var api *API

func initRouter(t *testing.T) (chi.Router, *MockStorage) {
	t.Helper()

	router := chi.NewRouter()
	cfg, _ := config.New()
	storage := NewMockStorage()
	l := logger.New(cfg.Log.Level)

	api = New(router, cfg, storage, l)

	return api.router, storage
}

func TestRouter(t *testing.T) {
	router, repo := initRouter(t)

	var testCases = []struct {
		url            string
		reqMethod      string
		reqBody        io.Reader
		expectedBody   string
		expectedStatus int
		mockRepo       func()
		expectedHeader string
	}{
		{
			url:            "/",
			reqMethod:      "POST",
			reqBody:        bytes.NewReader([]byte("https://practicum.yandex.ru/")),
			expectedBody:   fmt.Sprintf("http://localhost:8080/%s", "qwerty"),
			expectedStatus: http.StatusCreated,
			mockRepo: func() {
				repo.saveURL = func(originalURL string) (string, error) { return "qwerty", nil }
			},
		},
		{
			url:            "/",
			reqMethod:      "GET",
			reqBody:        bytes.NewReader([]byte("https://practicum.yandex.ru/")),
			expectedBody:   "",
			expectedStatus: http.StatusMethodNotAllowed,
			mockRepo:       func() {},
		},
		{
			url:            "/TlHZMa",
			reqMethod:      "GET",
			expectedBody:   "",
			expectedStatus: http.StatusTemporaryRedirect,
			mockRepo: func() {
				repo.getURL = func(key string) (string, bool) { return "google.com", true }
			},
			expectedHeader: "google.com",
		},
		{
			url:            "/qwerty",
			reqMethod:      "GET",
			expectedBody:   "shortKey not found\n",
			expectedStatus: http.StatusBadRequest,
			mockRepo: func() {
				repo.getURL = func(key string) (string, bool) { return "", false }
			},
			expectedHeader: "",
		},
	}
	for _, tc := range testCases {
		tc.mockRepo()
		ts := httptest.NewServer(router)

		resp, resBody := testRequest(t, ts, tc.reqMethod, tc.url, tc.reqBody)

		assert.Equal(t, tc.expectedStatus, resp.StatusCode)
		assert.Equal(t, tc.expectedBody, resBody)
		assert.Equal(t, tc.expectedHeader, resp.Header.Get("Location"))

		resp.Body.Close()
		ts.Close()
	}
}

func testRequest(t *testing.T, ts *httptest.Server, method, path string, body io.Reader) (*http.Response, string) {
	req, err := http.NewRequest(method, ts.URL+path, body)
	require.NoError(t, err)

	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	resp, err := client.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	return resp, string(respBody)
}
