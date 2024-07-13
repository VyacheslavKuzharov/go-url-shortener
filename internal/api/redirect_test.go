package api

import (
	"github.com/VyacheslavKuzharov/go-url-shortener/internal/storage/inmemory"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRedirectHandler(t *testing.T) {
	newAPI := New(chi.NewRouter())
	newAPI.InitRoutes(inmemory.NewMemoryStorage())

	originalURL := "https://practicum.yandex.ru/"

	testCases := []struct {
		name           string
		method         string
		request        string
		expectedCode   int
		expectedHeader string
		mock           *MockStorage
	}{
		{
			name:           "when happy path: correct response",
			method:         http.MethodGet,
			request:        "/qwerty",
			expectedCode:   http.StatusTemporaryRedirect,
			expectedHeader: originalURL,
			mock:           &MockStorage{getURL: func(key string) (string, bool) { return originalURL, true }},
		},
		{
			name:           "when unhappy path: incorrect request method",
			method:         http.MethodPost,
			request:        "/qwerty",
			expectedCode:   http.StatusMethodNotAllowed,
			expectedHeader: "",
			mock:           &MockStorage{getURL: func(key string) (string, bool) { return originalURL, true }},
		},
		{
			name:           "when unhappy path: short key not found",
			method:         http.MethodGet,
			request:        "/qwerty",
			expectedCode:   http.StatusBadRequest,
			expectedHeader: "",
			mock:           &MockStorage{getURL: func(key string) (string, bool) { return "", false }},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			r := httptest.NewRequest(tc.method, tc.request, nil)
			w := httptest.NewRecorder()

			h := redirectHandler(tc.mock)
			h(w, r)

			res := w.Result()
			defer res.Body.Close()

			// check response code
			assert.Equal(t, tc.expectedCode, w.Code, "Код ответа не совпадает с ожидаемым")
			// check response header
			assert.Equal(t, tc.expectedHeader, res.Header.Get("Location"))
		})
	}
}
