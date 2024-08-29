package api

import (
	"context"
	"github.com/VyacheslavKuzharov/go-url-shortener/internal/config"
	"github.com/VyacheslavKuzharov/go-url-shortener/internal/entity"
	uuid "github.com/satori/go.uuid"
)

func NewMockStorage() *MockStorage {
	return &MockStorage{}
}

type MockStorage struct {
	saveURL        func(context.Context, string) (string, error)
	saveBatchURLs  func(context.Context, []entity.ShortenURL) error
	getURL         func(context.Context, string) (string, error)
	getUserUrls    func(context.Context, uuid.UUID, *config.Config) ([]*entity.CompletedURL, error)
	deleteUserUrls func(context.Context, uuid.UUID, []string) error
	ping           func(context.Context) error
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

func (m *MockStorage) GetUserUrls(ctx context.Context, currentUserID uuid.UUID, cfg *config.Config) ([]*entity.CompletedURL, error) {
	return m.getUserUrls(ctx, currentUserID, cfg)
}

func (m *MockStorage) DeleteUserUrls(ctx context.Context, currentUserID uuid.UUID, urlKeysBatch []string) error {
	return m.deleteUserUrls(ctx, currentUserID, urlKeysBatch)
}

func (m *MockStorage) Ping(ctx context.Context) error {
	return m.ping(ctx)
}

func (m *MockStorage) Close() error {
	return nil
}
