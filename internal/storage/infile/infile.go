package infile

import (
	"encoding/json"
	"errors"
	"github.com/VyacheslavKuzharov/go-url-shortener/internal/entity"
	"github.com/VyacheslavKuzharov/go-url-shortener/internal/lib/random"
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

func (s *FileStorage) SaveURL(originalURL string) (string, error) {
	if originalURL == "" {
		return "", errors.New("originalURL can't be blank")
	}

	s.mutex.Lock()
	defer s.mutex.Unlock()

	shortKey := random.GenShortKey()
	su := entity.ShortenURL{
		ShortKey:    shortKey,
		OriginalURL: originalURL,
	}

	err := s.encoder.Encode(&su)
	if err != nil {
		return "", err
	}

	return su.ShortKey, nil
}

func (s *FileStorage) GetURL(key string) (string, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	file, err := os.Open(s.file.Name())
	if err != nil {
		return "", err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	su := entity.ShortenURL{}

	for decoder.More() {
		err = decoder.Decode(&su)
		if err != nil {
			return "", err
		}

		if su.ShortKey == key {
			return su.OriginalURL, nil
		}
	}

	return "", errors.New("shortKey not found")
}

func (s *FileStorage) Close() error {
	return s.file.Close()
}

func (s *FileStorage) Ping() error {
	return nil
}
