package control

import (
	"github.com/polis-interactive/slate-tv/internal/domain"
	"sync"
)

func mergeInputEvents(cs ...chan domain.InputPair) chan domain.InputPair {
	out := make(chan domain.InputPair, 10)
	var wg sync.WaitGroup
	wg.Add(len(cs))
	for _, c := range cs {
		go func(c <-chan domain.InputPair) {
			for v := range c {
				out <- v
			}
			wg.Done()
		}(c)
	}
	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}
