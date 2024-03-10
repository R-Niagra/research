package queue

import (
	"sync/atomic"
	"unsafe"
)

type LockfreeQueue struct {
	capacity uint64
	head     *Node
	tail     *Node
	size     uint64
}

func NewLockfreeQueue(capacity uint64) *LockfreeQueue {
	lfq := &LockfreeQueue{
		capacity: capacity,
		head:     new(Node),
	}
	lfq.tail = lfq.head
	return lfq
}

func (lfq *LockfreeQueue) Capacity() uint64 {
	return lfq.capacity
}

func (lfq *LockfreeQueue) EmptyQueue() {
	lfq.head = new(Node)
	lfq.tail = lfq.head
	lfq.size = 0
}

// Enqueue: Non-blocking
func (lfq *LockfreeQueue) TryEnqueue(val int) error {

	curSize := atomic.LoadUint64(&lfq.size) //curSize of the queue
	if curSize >= lfq.capacity {
		return errAtCapacity
	}
	var (
		last      *Node
		lastsNext *Node
	)

	n := NewNode(val)
	// fmt.Println(n)
	for {
		last = lfq.getTail() //atomic operation
		lastsNext = last.Next

		if last == lfq.getTail() {
			if lastsNext == nil {
				if atomic.CompareAndSwapPointer((*unsafe.Pointer)(unsafe.Pointer(&last.Next)), unsafe.Pointer(lastsNext), unsafe.Pointer(n)) {

					//if atomic swap was a success, set the tail as well
					atomic.CompareAndSwapPointer((*unsafe.Pointer)(unsafe.Pointer(&lfq.tail)), unsafe.Pointer(last), unsafe.Pointer(n))
					atomic.AddUint64(&lfq.size, 1)
					break
				}
			} else {
				atomic.CompareAndSwapPointer((*unsafe.Pointer)(unsafe.Pointer(&lfq.tail)), unsafe.Pointer(last), unsafe.Pointer(lastsNext))
			}
		}
	}
	return nil
}

// Dequeue: Non-blocking
func (lfq *LockfreeQueue) TryDequeue() (*Node, error) {
	curSize := atomic.LoadUint64(&lfq.size) //curSize of the queue
	if curSize == 0 {
		return nil, errEmptyQueue
	}

	var (
		first      *Node
		firstsNext *Node
		last       *Node
		result     *Node
	)

	for {
		first = lfq.getHead()
		firstsNext = first.Next
		last = lfq.getTail()

		if first == lfq.getHead() { //if first is still the head
			if first == last { //if it is also the last item then the queue is empty
				//queue has 1 dummy item by default
				if firstsNext == nil {
					return nil, errEmptyQueue
				}
				atomic.CompareAndSwapPointer((*unsafe.Pointer)(unsafe.Pointer(&lfq.tail)), unsafe.Pointer(last), unsafe.Pointer(firstsNext))
			} else {
				result = firstsNext
				//move head forward
				if atomic.CompareAndSwapPointer((*unsafe.Pointer)(unsafe.Pointer(&lfq.head)), unsafe.Pointer(first), unsafe.Pointer(firstsNext)) {
					atomic.AddUint64(&lfq.size, ^uint64(1)+1) //subtract 1 from size using two's complement
					return result, nil
				}
			}
		}

	}

}

func (lfq *LockfreeQueue) getTail() *Node {
	return (*Node)(atomic.LoadPointer((*unsafe.Pointer)(unsafe.Pointer(&lfq.tail))))
}

func (lfq *LockfreeQueue) getHead() *Node {
	return (*Node)(atomic.LoadPointer((*unsafe.Pointer)(unsafe.Pointer(&lfq.head))))
}
