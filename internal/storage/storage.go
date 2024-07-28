package storage

import (
	"errors"
	"github.com/VyacheslavKuzharov/go-url-shortener/internal/config"
	storagecfg "github.com/VyacheslavKuzharov/go-url-shortener/internal/config/storage"
	"github.com/VyacheslavKuzharov/go-url-shortener/internal/storage/infile"
	"github.com/VyacheslavKuzharov/go-url-shortener/internal/storage/inmemory"
)

type Storager interface {
	SaveURL(originalURL string) (string, error)
	GetURL(key string) (string, bool)
	Close() error
}

func New(cfg *config.Config) (Storager, error) {
	switch cfg.Storage.Kind {
	case storagecfg.InMemory:
		return inmemory.NewMemoryStorage()
	case storagecfg.InFile:
		return infile.NewFileStorage(cfg.Storage.File.Path)
	default:
		return nil, errors.New("unknown storage type")
	}
}
