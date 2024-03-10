package queue

import (
	"runtime"
	"sync"
)

type Queue interface {
	TryEnqueue(int) error
	TryDequeue() (*Node, error)
	EmptyQueue()
}

func FillQueueToCapacityConcurrently(q Queue, cap int) {
	workers := runtime.NumCPU()
	var wg sync.WaitGroup
	wg.Add(workers)
	for i := 0; i < workers; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < int(cap); j++ {
				err := q.TryEnqueue(j)
				if err != nil {
					return
				}
			}
		}()
	}
	wg.Wait()
}

func RemoveFromQueueConcurrently(q Queue, cap int) {
	workers := runtime.NumCPU()
	var wg sync.WaitGroup
	wg.Add(workers)
	for i := 0; i < workers; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < int(cap); j++ {
				_, err := q.TryDequeue()
				if err == errEmptyQueue {
					return
				}
			}
		}()
	}
	wg.Wait()
}
