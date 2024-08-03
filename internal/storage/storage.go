package storage

import (
	"errors"
	"github.com/VyacheslavKuzharov/go-url-shortener/internal/config"
	storagecfg "github.com/VyacheslavKuzharov/go-url-shortener/internal/config/storage"
	"github.com/VyacheslavKuzharov/go-url-shortener/internal/storage/infile"
	"github.com/VyacheslavKuzharov/go-url-shortener/internal/storage/inmemory"
	"github.com/VyacheslavKuzharov/go-url-shortener/internal/storage/postgres"
)

type Storager interface {
	SaveURL(originalURL string) (string, error)
	GetURL(key string) (string, error)
	Close() error
	Ping() error
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
