package api

import (
	"context"
	"github.com/VyacheslavKuzharov/go-url-shortener/internal/entity"
)

func NewMockStorage() *MockStorage {
	return &MockStorage{}
}

type MockStorage struct {
	saveURL       func(context.Context, string) (string, error)
	saveBatchURLs func(context.Context, []entity.ShortenURL) error
	getURL        func(context.Context, string) (string, error)
	ping          func(context.Context) error
}

func (m *MockStorage) SaveURL(ctx context.Context, originalURL string) (string, error) {
	return m.saveURL(ctx, originalURL)
}

func (m *MockStorage) SaveBatchURLs(ctx context.Context, urls []entity.ShortenURL) error {
	return m.saveBatchURLs(ctx, urls)
}

func (m *MockStorage) GetURL(ctx context.Context, key string) (string, error) {
	return m.getURL(ctx, key)
}

func (m *MockStorage) Ping(ctx context.Context) error {
	return m.ping(ctx)
}

func (m *MockStorage) Close() error {
	return nil
}
