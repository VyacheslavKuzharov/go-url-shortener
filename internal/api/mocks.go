package api

import (
	"context"
	"github.com/VyacheslavKuzharov/go-url-shortener/internal/entity"
)

func NewMockStorage() *MockStorage {
	return &MockStorage{}
}

type MockStorage struct {
	saveURL       func(string) (string, error)
	saveBatchURLs func(ctx context.Context, urls *[]entity.ShortenURL) error
	getURL        func(string) (string, error)
	ping          func() error
}

func (m *MockStorage) SaveURL(originalURL string) (string, error) {
	return m.saveURL(originalURL)
}

func (m *MockStorage) SaveBatchURLs(ctx context.Context, urls *[]entity.ShortenURL) error {
	return m.saveBatchURLs(ctx, urls)
}

func (m *MockStorage) GetURL(key string) (string, error) {
	return m.getURL(key)
}

func (m *MockStorage) Close() error {
	return nil
}

func (m *MockStorage) Ping() error {
	return m.ping()
}
