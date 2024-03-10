package queue

import (
	"errors"
	"sync"
	"sync/atomic"
)

var (
	errAtCapacity = errors.New("queue is at capacity")
	errEmptyQueue = errors.New("queue is empty")
)

// concurrent queue implementation using locks
type LockQueue struct {
	capacity uint64
	head     *Node
	tail     *Node
	eLock    sync.Mutex //enqueue lock
	dLock    sync.Mutex //dequeue lock
	size     uint64
}

func NewBoundedQueue(cap uint64) *LockQueue {
	lf := &LockQueue{
		capacity: cap,
		head:     &Node{}, //initial empty node
	}
	lf.tail = lf.head
	return lf
}

func (lq *LockQueue) Capacity() uint64 {
	return lq.capacity
}

func (lq *LockQueue) EmptyQueue() {
	lq.head = &Node{}
	lq.tail = lq.head
	lq.size = 0
}

func (lq *LockQueue) Enqueue(val int) {
	if err := lq.TryEnqueue(val); err == nil {
		return
	}
	//queue must be at capacity
	//make space by Dequeuing

}

// Enqueue: non-blocking; Error out if queue is at capacity
func (lq *LockQueue) TryEnqueue(val int) error {
	lq.eLock.Lock()
	defer lq.eLock.Unlock()

	curSize := atomic.LoadUint64(&lq.size) //curSize of the queue

	if curSize >= lq.capacity {
		return errAtCapacity
	}

	n := NewNode(val)
	lq.tail.Next = n //add the node
	lq.tail = n      //move tail
	//don't need to compare with curSize because only enqueue is allowed to increment size
	atomic.AddUint64(&lq.size, 1)

	return nil
}

// Dequeue: non-Blocking; dequeues only when queue is non-empty
func (lq *LockQueue) TryDequeue() (*Node, error) {
	lq.dLock.Lock()
	defer lq.dLock.Unlock()

	curSize := atomic.LoadUint64(&lq.size) //curSize of the queue
	if curSize == 0 {
		return nil, errEmptyQueue
	}

	lq.head = lq.head.Next
	res := lq.head
	//prev head should automatically be garbage collected once the node is not reachable

	//subtract 1 from size using two's complement
	atomic.AddUint64(&lq.size, ^uint64(1)+1)

	return res, nil
}
