package api

import (
	"bytes"
	"compress/gzip"
	"context"
	"errors"
	"fmt"
	"github.com/VyacheslavKuzharov/go-url-shortener/internal/api/middlewares"
	"github.com/VyacheslavKuzharov/go-url-shortener/internal/config"
	httpcfg "github.com/VyacheslavKuzharov/go-url-shortener/internal/config/http"
	"github.com/VyacheslavKuzharov/go-url-shortener/internal/logger"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
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

	return resp, strings.TrimSuffix(string(respBody), "\n")
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
				repo.saveURL = func(ctx context.Context, originalURL string) (string, error) { return "qwerty", nil }
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
				repo.getURL = func(ctx context.Context, key string) (string, error) { return "google.com", nil }
			},
			expectedHeader: "google.com",
		},
		{
			url:            "/qwerty",
			reqMethod:      "GET",
			expectedBody:   "shortKey not found",
			expectedStatus: http.StatusGone,
			mockRepo: func() {
				repo.getURL = func(ctx context.Context, key string) (string, error) { return "", errors.New("shortKey not found") }
			},
			expectedHeader: "",
		},
		{
			url:            "/api/shorten",
			reqMethod:      "POST",
			reqBody:        bytes.NewReader([]byte("{\"url\": \"https://practicum.yandex.ru\"}")),
			expectedBody:   `{"result":"http://localhost:8080/fG4oX4"}`,
			expectedStatus: http.StatusCreated,
			mockRepo: func() {
				repo.saveURL = func(ctx context.Context, originalURL string) (string, error) { return "fG4oX4", nil }
			},
		},
		{
			url:            "/api/shorten",
			reqMethod:      "GET",
			reqBody:        bytes.NewReader([]byte("{\"url\": \"https://practicum.yandex.ru\"}")),
			expectedBody:   "",
			expectedStatus: http.StatusMethodNotAllowed,
			mockRepo:       func() {},
		},
		{
			url:            "/api/shorten",
			reqMethod:      "POST",
			reqBody:        bytes.NewReader([]byte("{\"url\": \"https://\"}")),
			expectedBody:   `{"error":"provided url is invalid"}`,
			expectedStatus: http.StatusBadRequest,
			mockRepo:       func() {},
		},
		{
			url:            "/api/shorten",
			reqMethod:      "POST",
			reqBody:        bytes.NewReader([]byte("{}")),
			expectedBody:   `{"error":"provided url is invalid"}`,
			expectedStatus: http.StatusBadRequest,
			mockRepo:       func() {},
		},
		{
			url:            "/api/shorten",
			reqMethod:      "POST",
			reqBody:        bytes.NewReader([]byte("")),
			expectedBody:   `{"error":"request is empty"}`,
			expectedStatus: http.StatusBadRequest,
			mockRepo:       func() {},
		},
		{
			url:            "/ping",
			reqMethod:      "GET",
			reqBody:        bytes.NewReader([]byte("")),
			expectedBody:   "",
			expectedStatus: http.StatusOK,
			mockRepo:       func() { repo.ping = func(ctx context.Context) error { return nil } },
		},
		{
			url:            "/api/user/urls",
			reqMethod:      "GET",
			reqBody:        bytes.NewReader([]byte("")),
			expectedBody:   `{"error":"http: named cookie not present"}`,
			expectedStatus: http.StatusUnauthorized,
			mockRepo:       func() {},
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

func TestGzipCompression(t *testing.T) {
	cfgs := &config.Config{HTTP: httpcfg.HTTPCfg{Host: "localhost", Port: "8080"}}
	mock := &MockStorage{saveURL: func(ctx context.Context, originalURL string) (string, error) { return "NUf6O3", nil }}

	handler := middlewares.Compress(shortenHandler(mock, cfgs))

	srv := httptest.NewServer(handler)
	defer srv.Close()

	requestBody := `{
		"url": "https://practicum.yandex.ru"
	}`

	successBody := `{
			"result": "http://localhost:8080/NUf6O3"
	}`

	t.Run("sends_gzip", func(t *testing.T) {
		buf := bytes.NewBuffer(nil)
		zb := gzip.NewWriter(buf)
		_, err := zb.Write([]byte(requestBody))
		require.NoError(t, err)
		err = zb.Close()
		require.NoError(t, err)

		r := httptest.NewRequest("POST", srv.URL, buf)
		r.RequestURI = ""
		r.Header.Set("Content-Type", "application/json")
		r.Header.Set("Content-Encoding", "gzip")
		r.Header.Set("Accept-Encoding", "")

		resp, err := http.DefaultClient.Do(r)
		require.NoError(t, err)
		require.Equal(t, http.StatusCreated, resp.StatusCode)

		defer resp.Body.Close()

		b, err := io.ReadAll(resp.Body)
		require.NoError(t, err)
		require.JSONEq(t, successBody, string(b))
	})

	t.Run("accepts_gzip", func(t *testing.T) {
		buf := bytes.NewBufferString(requestBody)
		r := httptest.NewRequest("POST", srv.URL, buf)
		r.RequestURI = ""
		r.Header.Set("Content-Type", "text/html")
		r.Header.Set("Accept-Encoding", "gzip")

		resp, err := http.DefaultClient.Do(r)
		require.NoError(t, err)
		require.Equal(t, http.StatusCreated, resp.StatusCode)

		defer resp.Body.Close()

		zr, err := gzip.NewReader(resp.Body)
		require.NoError(t, err)

		b, err := io.ReadAll(zr)
		require.NoError(t, err)

		require.JSONEq(t, successBody, string(b))
	})
}
