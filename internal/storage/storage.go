package storage

import (
	"context"
	"errors"
	"github.com/VyacheslavKuzharov/go-url-shortener/internal/config"
	storagecfg "github.com/VyacheslavKuzharov/go-url-shortener/internal/config/storage"
	"github.com/VyacheslavKuzharov/go-url-shortener/internal/entity"
	"github.com/VyacheslavKuzharov/go-url-shortener/internal/storage/infile"
	"github.com/VyacheslavKuzharov/go-url-shortener/internal/storage/inmemory"
	"github.com/VyacheslavKuzharov/go-url-shortener/internal/storage/postgres"
	uuid "github.com/satori/go.uuid"
)

type Storager interface {
	SaveURL(ctx context.Context, originalURL string) (string, error)
	GetURL(ctx context.Context, key string) (string, error)
	GetUserUrls(ctx context.Context, currentUserID uuid.UUID, cfg *config.Config) ([]*entity.CompletedURL, error)
	SaveBatchURLs(ctx context.Context, urls []entity.ShortenURL) error
	DeleteUserUrls(ctx context.Context, currentUserID uuid.UUID, urlKeysBatch []string) error
	Ping(ctx context.Context) error
	Close() error
}

func New(cfg *config.Config) (Storager, error) {
	switch cfg.Storage.Kind {
	case storagecfg.InMemory:
		return inmemory.NewMemoryStorage()
	case storagecfg.InFile:
		return infile.NewFileStorage(cfg.Storage.File.Path)
	case storagecfg.Postgres:
		return postgres.New(cfg.Storage.Postgres.ConnectURL)
	default:
		return nil, errors.New("unknown storage type")
	}
}
