//package storage
//
//import "github.com/VyacheslavKuzharov/go-url-shortener/internal/storage/inmemory"
//
//type Storager interface {
//	StoreURL(originalURL string) (string, error)
//	GetURLByShortKey(key string) (string, bool)
//}
//
//func New() Storager {
//	return inmemory.NewMemoryStorage()
//}

package storage

import "errors"

var (
	ErrURLNotFound = errors.New("url not found")
	ErrURLExists   = errors.New("url exists")
)
