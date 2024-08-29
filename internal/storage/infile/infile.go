package infile

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/VyacheslavKuzharov/go-url-shortener/internal/config"
	"github.com/VyacheslavKuzharov/go-url-shortener/internal/entity"
	"github.com/VyacheslavKuzharov/go-url-shortener/internal/lib/httpapi"
	"github.com/VyacheslavKuzharov/go-url-shortener/internal/lib/random"
	uuid "github.com/satori/go.uuid"
	"os"
	"sync"
)

type FileStorage struct {
	mutex   sync.RWMutex
	file    *os.File
	encoder *json.Encoder
}

func NewFileStorage(fileName string) (*FileStorage, error) {
	file, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}

	return &FileStorage{
		file:    file,
		encoder: json.NewEncoder(file),
	}, nil
}

func (s *FileStorage) SaveURL(ctx context.Context, originalURL string) (string, error) {
	if originalURL == "" {
		return "", errors.New("originalURL can't be blank")
	}

	s.mutex.Lock()
	defer s.mutex.Unlock()

	userID, ok := ctx.Value(entity.CurrentUserID).(uuid.UUID)
	if !ok {
		return "", errors.New("invalid uuid type in infile.SaveURL()")
	}

	us := entity.UserShortenURL{
		User: entity.User{
			UUID: userID,
		},
		ShortenURL: entity.ShortenURL{
			ShortKey:    random.GenShortKey(),
			OriginalURL: originalURL},
	}

	err := s.encoder.Encode(&us)
	if err != nil {
		return "", err
	}

	return us.ShortenURL.ShortKey, nil
}

func (s *FileStorage) SaveBatchURLs(ctx context.Context, urls []entity.ShortenURL) error {
	if len(urls) == 0 {
		return nil
	}
	s.mutex.Lock()
	defer s.mutex.Unlock()

	userID, ok := ctx.Value(entity.CurrentUserID).(uuid.UUID)
	if !ok {
		return errors.New("invalid uuid type in infile.SaveBatchURLs()")
	}

	for _, su := range urls {
		us := entity.UserShortenURL{
			User: entity.User{
				UUID: userID,
			},
			ShortenURL: su,
		}

		err := s.encoder.Encode(&us)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *FileStorage) GetURL(ctx context.Context, key string) (string, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	file, err := os.Open(s.file.Name())
	if err != nil {
		return "", err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	us := entity.UserShortenURL{
		User:       entity.User{},
		ShortenURL: entity.ShortenURL{},
	}

	for decoder.More() {
		err = decoder.Decode(&us)
		if err != nil {
			return "", err
		}

		if us.ShortenURL.ShortKey == key && !us.IsDeleted {
			return us.ShortenURL.OriginalURL, nil
		}
	}

	return "", errors.New("shortKey not found")
}

func (s *FileStorage) GetUserUrls(ctx context.Context, currentUserID uuid.UUID, cfg *config.Config) ([]*entity.CompletedURL, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	var userURLs []*entity.CompletedURL

	file, err := os.Open(s.file.Name())
	if err != nil {
		return nil, err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	us := entity.UserShortenURL{
		User:       entity.User{},
		ShortenURL: entity.ShortenURL{},
	}

	for decoder.More() {
		err = decoder.Decode(&us)
		if err != nil {
			return nil, err
		}

		if uuid.Equal(us.User.UUID, currentUserID) && !us.IsDeleted {
			cu := &entity.CompletedURL{
				ShortURL:    httpapi.FullShortenedURL(us.ShortenURL.ShortKey, cfg),
				OriginalURL: us.ShortenURL.OriginalURL,
			}

			userURLs = append(userURLs, cu)
		}
	}

	return userURLs, nil
}

func (s *FileStorage) DeleteUserUrls(ctx context.Context, currentUserID uuid.UUID, urlKeysBatch []string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	var storageURLs []entity.UserShortenURL

	f, err := os.Open(s.file.Name())
	if err != nil {
		return err
	}
	defer f.Close()

	decoder := json.NewDecoder(f)
	us := entity.UserShortenURL{
		User:       entity.User{},
		ShortenURL: entity.ShortenURL{},
		IsDeleted:  false,
	}

	for decoder.More() {
		err = decoder.Decode(&us)
		if err != nil {
			return err
		}

		storageURLs = append(storageURLs, us)
	}
	err = s.file.Truncate(0)
	if err != nil {
		return err
	}
	_, err = s.file.Seek(0, 0)
	if err != nil {
		return err
	}

	for _, item := range storageURLs {
		if contains(urlKeysBatch, item.ShortenURL.ShortKey) && !item.IsDeleted {
			index := indexOf(item, storageURLs)
			if index >= 0 {
				storageURLs = remove(storageURLs, index)

				item.IsDeleted = true
				storageURLs = append(storageURLs, item)
			}
		}
	}

	for _, item := range storageURLs {
		err = s.encoder.Encode(&item)
		if err != nil {
			return err
		}
	}

	return nil
}

func contains(s []string, el string) bool {
	for _, a := range s {
		if a == el {
			return true
		}
	}
	return false
}

func indexOf(el entity.UserShortenURL, data []entity.UserShortenURL) int {
	for k, v := range data {
		if el == v {
			return k
		}
	}
	return -1 //not found.
}

func remove(s []entity.UserShortenURL, i int) []entity.UserShortenURL {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}

func (s *FileStorage) Close() error {
	return s.file.Close()
}

func (s *FileStorage) Ping(ctx context.Context) error {
	return nil
}
