package softdelete

import (
	"context"
	"github.com/VyacheslavKuzharov/go-url-shortener/internal/entity"
)

type WorkerResult struct {
	Res bool
	Err error
}

func MemWorker(memUrls map[string]entity.UserShortenURL, ctx context.Context, obj Object) chan WorkerResult {
	workerChan := make(chan WorkerResult)

	go func() {
		defer close(workerChan)

		ok, err := deleteMemUrls(memUrls, ctx, obj)
		workerRes := WorkerResult{
			Res: ok,
			Err: err,
		}

		workerChan <- workerRes
	}()

	return workerChan
}

func deleteMemUrls(memUrls map[string]entity.UserShortenURL, ctx context.Context, delObj Object) (bool, error) {
	for _, key := range delObj.ShortKeys {
		el, ok := memUrls[key]
		if ok {
			el.IsDeleted = true
			memUrls[key] = el
		}
	}

	return true, nil
}
