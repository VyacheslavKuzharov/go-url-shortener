package inmemory

import (
	"context"
	"errors"
	"github.com/VyacheslavKuzharov/go-url-shortener/internal/entity"
	"github.com/VyacheslavKuzharov/go-url-shortener/internal/lib/random"
	"sync"
)

type MemStorage struct {
	mutex sync.RWMutex
	urls  map[string]string
}

func NewMemoryStorage() (*MemStorage, error) {
	return &MemStorage{
		urls: make(map[string]string),
	}, nil
}

func (s *MemStorage) SaveURL(ctx context.Context, originalURL string) (string, error) {
	if originalURL == "" {
		return "", errors.New("originalURL can't be blank")
	}

	s.mutex.Lock()
	defer s.mutex.Unlock()

	shortKey := random.GenShortKey()

	s.urls[shortKey] = originalURL
	return shortKey, nil
}

func (s *MemStorage) SaveBatchURLs(ctx context.Context, urls []entity.ShortenURL) error {
	if len(urls) == 0 {
		return nil
	}
	s.mutex.Lock()
	defer s.mutex.Unlock()

	for _, su := range urls {
		s.urls[su.ShortKey] = su.OriginalURL
	}

	return nil
}

func (s *MemStorage) GetURL(ctx context.Context, key string) (string, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	originalURL, ok := s.urls[key]
	if !ok {
		return "", errors.New("shortKey not found")
	}

	return originalURL, nil
}

func (s *MemStorage) Close() error {
	return nil
}

func (s *MemStorage) Ping(ctx context.Context) error {
	return nil
}
