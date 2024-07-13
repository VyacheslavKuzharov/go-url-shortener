package storage

import "github.com/VyacheslavKuzharov/go-url-shortener/internal/storage/inmemory"

type Storager interface {
	SaveURL(originalURL string) (string, error)
	GetURL(key string) (string, bool)
}

func New() Storager {
	return inmemory.NewMemoryStorage()
}
