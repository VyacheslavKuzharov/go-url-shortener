package inmemory

import (
	"context"
	"errors"
	"github.com/VyacheslavKuzharov/go-url-shortener/internal/config"
	"github.com/VyacheslavKuzharov/go-url-shortener/internal/entity"
	"github.com/VyacheslavKuzharov/go-url-shortener/internal/lib/httpapi"
	"github.com/VyacheslavKuzharov/go-url-shortener/internal/lib/random"
	uuid "github.com/satori/go.uuid"
	"sync"
)

type MemStorage struct {
	mutex sync.RWMutex
	urls  map[string]entity.UserShortenURL
}

func NewMemoryStorage() (*MemStorage, error) {
	return &MemStorage{
		urls: make(map[string]entity.UserShortenURL),
	}, nil
}

func (s *MemStorage) SaveURL(ctx context.Context, originalURL string) (string, error) {
	if originalURL == "" {
		return "", errors.New("originalURL can't be blank")
	}

	s.mutex.Lock()
	defer s.mutex.Unlock()

	userID, ok := ctx.Value(entity.CurrentUserID).(uuid.UUID)
	if !ok {
		return "", errors.New("invalid uuid type in inmemory.SaveURL()")
	}

	shortKey := random.GenShortKey()

	item := entity.UserShortenURL{
		User: entity.User{
			UUID: userID,
		},
		ShortenURL: entity.ShortenURL{
			ShortKey:    shortKey,
			OriginalURL: originalURL},
	}
	s.urls[shortKey] = item

	return shortKey, nil
}

func (s *MemStorage) SaveBatchURLs(ctx context.Context, urls []entity.ShortenURL) error {
	if len(urls) == 0 {
		return nil
	}
	s.mutex.Lock()
	defer s.mutex.Unlock()

	userID, ok := ctx.Value(entity.CurrentUserID).(uuid.UUID)
	if !ok {
		return errors.New("invalid uuid type in inmemory.SaveBatchURLs()")
	}

	for _, su := range urls {
		item := entity.UserShortenURL{
			User: entity.User{
				UUID: userID,
			},
			ShortenURL: su,
		}

		s.urls[su.ShortKey] = item
	}

	return nil
}

func (s *MemStorage) GetURL(ctx context.Context, key string) (string, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	item, ok := s.urls[key]
	if !ok {
		return "", errors.New("shortKey not found")
	}

	return item.OriginalURL, nil
}

func (s *MemStorage) GetUserUrls(ctx context.Context, currentUserID uuid.UUID, cfg *config.Config) ([]*entity.CompletedURL, error) {
	var userURLs []*entity.CompletedURL

	for _, v := range s.urls {
		if uuid.Equal(v.User.UUID, currentUserID) {
			urlItem := &entity.CompletedURL{
				ShortURL:    httpapi.FullShortenedURL(v.ShortenURL.ShortKey, cfg),
				OriginalURL: v.ShortenURL.OriginalURL,
			}

			userURLs = append(userURLs, urlItem)
		}
	}

	return userURLs, nil
}

func (s *MemStorage) Close() error {
	return nil
}

func (s *MemStorage) Ping(ctx context.Context) error {
	return nil
}
