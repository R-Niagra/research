package queue

import (
	"testing"
)

var (
	defaultCapacity = 200000
)

func createLockQueue() *LockQueue {
	return NewBoundedQueue(uint64(defaultCapacity))
}

func createLockFreeQueue() *LockfreeQueue {
	return NewLockfreeQueue(uint64(defaultCapacity))
}

// goos: linux
// goarch: amd64
// pkg: github.com/R-Niagra/research/queue
// cpu: Intel(R) Core(TM) i7-10510U CPU @ 1.80GHz
// BenchmarkLockEnqueueAndDequeue-8   	     409	   2731162 ns/op	  262160 B/op	   16384 allocs/op
// PASS
// ok  	github.com/R-Niagra/research/queue	1.414s

func BenchmarkLockEnqueueAndDequeue(b *testing.B) {
	items := 1 << 14
	q := createLockQueue()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			for j := 0; j < items; j++ {
				q.TryEnqueue(j)
			}
			for j := 0; j < items; j++ {
				q.TryDequeue()
			}
		}
	})
}

// goos: linux
// goarch: amd64
// pkg: github.com/R-Niagra/research/queue
// cpu: Intel(R) Core(TM) i7-10510U CPU @ 1.80GHz
// BenchmarkLockFreeEnqueueAndDequeue-8   	     408	   2702549 ns/op	  262159 B/op	   16384 allocs/op
// PASS
// ok  	github.com/R-Niagra/research/queue	1.401s

func BenchmarkLockFreeEnqueueAndDequeue(b *testing.B) {
	items := 1 << 14
	q := createLockFreeQueue()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			for j := 0; j < items; j++ {
				q.TryEnqueue(j)
			}
			for j := 0; j < items; j++ {
				q.TryDequeue()
			}
		}
	})
}
