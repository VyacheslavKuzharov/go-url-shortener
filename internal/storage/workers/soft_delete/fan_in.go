package softdelete

import "sync"

func FanIn(channels []<-chan WorkerResult) <-chan WorkerResult {
	finalCh := make(chan WorkerResult)
	var wg sync.WaitGroup

	wg.Add(len(channels))

	for _, channel := range channels {
		ch := channel

		go func() {
			defer wg.Done()

			for res := range ch {
				finalCh <- res
			}

		}()
	}

	go func() {
		wg.Wait()
		close(finalCh)
	}()

	return finalCh
}
