package softdelete

import (
	"encoding/json"
	"github.com/VyacheslavKuzharov/go-url-shortener/internal/entity"
	"os"
	"sync"
)

func FileWorker(file *os.File, encoder *json.Encoder, mu *sync.RWMutex, obj Object) chan WorkerResult {
	workerChan := make(chan WorkerResult)

	go func() {
		defer close(workerChan)

		ok, err := deleteFileUrls(file, encoder, mu, obj)
		workerRes := WorkerResult{
			Res: ok,
			Err: err,
		}

		workerChan <- workerRes
	}()

	return workerChan
}

func deleteFileUrls(file *os.File, encoder *json.Encoder, mu *sync.RWMutex, delObj Object) (bool, error) {
	mu.Lock()
	defer mu.Unlock()

	var storageURLs []entity.UserShortenURL

	f, err := os.Open(file.Name())
	if err != nil {
		return false, err
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
			return false, err
		}

		storageURLs = append(storageURLs, us)
	}
	err = file.Truncate(0)
	if err != nil {
		return false, err
	}
	_, err = file.Seek(0, 0)
	if err != nil {
		return false, err
	}

	for _, item := range storageURLs {
		if contains(delObj.ShortKeys, item.ShortenURL.ShortKey) && !item.IsDeleted {
			index := indexOf(item, storageURLs)
			if index >= 0 {
				storageURLs = remove(storageURLs, index)

				item.IsDeleted = true
				storageURLs = append(storageURLs, item)
			}
		}
	}

	for _, item := range storageURLs {
		err = encoder.Encode(&item)
		if err != nil {
			return false, err
		}
	}

	return true, nil
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
