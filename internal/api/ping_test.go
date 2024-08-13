package api

import (
	"context"
	"errors"
	"github.com/VyacheslavKuzharov/go-url-shortener/internal/storage/mock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPingHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	s := mock.NewMockStorager(ctrl)

	testCases := []struct {
		name         string // добавим название тестов
		method       string
		body         string // добавим тело запроса в табличные тесты
		expectedCode int
		expectedBody string
		mocks        *gomock.Call
	}{
		{
			name:         "when ok",
			method:       http.MethodGet,
			expectedCode: http.StatusOK,
			expectedBody: "",
			mocks:        s.EXPECT().Ping(context.Background()).Return(nil),
		},
		{
			name:         "when InternalServerError",
			method:       http.MethodGet,
			expectedCode: http.StatusInternalServerError,
			expectedBody: "",
			mocks:        s.EXPECT().Ping(context.Background()).Return(errors.New("test")),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			r := httptest.NewRequest(tc.method, "/ping", nil)
			w := httptest.NewRecorder()

			h := pingHandler(s)
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
