package api

import (
	"bytes"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

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

func TestRouter(t *testing.T) {
	newAPI := New(chi.NewRouter())

	var testCases = []struct {
		url            string
		reqMethod      string
		reqBody        io.Reader
		expectedBody   string
		expectedStatus int
		mock           *MockStorage
		expectedHeader string
	}{
		{
			url:            "/",
			reqMethod:      "POST",
			reqBody:        bytes.NewReader([]byte("https://practicum.yandex.ru/")),
			expectedBody:   fmt.Sprintf("http://localhost:8080/%s", "qwerty"),
			expectedStatus: http.StatusCreated,
			mock:           &MockStorage{saveURL: func(originalURL string) (string, error) { return "qwerty", nil }},
		},
		{
			url:            "/",
			reqMethod:      "GET",
			reqBody:        bytes.NewReader([]byte("https://practicum.yandex.ru/")),
			expectedBody:   "",
			expectedStatus: http.StatusMethodNotAllowed,
			mock:           &MockStorage{},
		},
		{
			url:            "/TlHZMa",
			reqMethod:      "GET",
			expectedBody:   "",
			expectedStatus: http.StatusTemporaryRedirect,
			mock:           &MockStorage{getURL: func(key string) (string, bool) { return "google.com", true }},
			expectedHeader: "google.com",
		},
		{
			url:            "/qwerty",
			reqMethod:      "GET",
			expectedBody:   "shortKey not found\n",
			expectedStatus: http.StatusBadRequest,
			mock:           &MockStorage{getURL: func(key string) (string, bool) { return "", false }},
			expectedHeader: "",
		},
	}
	for _, tc := range testCases {
		ts := httptest.NewServer(newAPI.InitRoutes(tc.mock))

		resp, resBody := testRequest(t, ts, tc.reqMethod, tc.url, tc.reqBody)

		assert.Equal(t, tc.expectedStatus, resp.StatusCode)
		assert.Equal(t, tc.expectedBody, resBody)
		assert.Equal(t, tc.expectedHeader, resp.Header.Get("Location"))

		resp.Body.Close()
		ts.Close()
	}
}
