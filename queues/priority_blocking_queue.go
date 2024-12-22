package queues

import (
	"sync"
)

// PriorityBlockingQueue implements BlockingQueue interface using PriorityQueue
type PriorityBlockingQueue[T any] struct {
	queue    *PriorityQueue[T]
	mutex    sync.Mutex
	notEmpty *sync.Cond
	notFull  *sync.Cond
}

var _ BlockingQueue[int] = &PriorityBlockingQueue[int]{}
var _ PriorityQueueInterface[int] = &PriorityBlockingQueue[int]{}

// NewPriorityBlockingQueue creates a new PriorityBlockingQueue with the specified capacity
func NewPriorityBlockingQueue[T any](capacity int) *PriorityBlockingQueue[T] {
	q := &PriorityBlockingQueue[T]{
		queue: NewPriorityQueue[T](capacity),
	}
	q.notEmpty = sync.NewCond(&q.mutex)
	q.notFull = sync.NewCond(&q.mutex)
	return q
}

// Offer adds an item to the queue if space is available, returning immediately
func (q *PriorityBlockingQueue[T]) Offer(item T) error {
	return q.OfferWithPriority(item, 0)
}

// OfferWithPriority adds an item with the specified priority if space is available
func (q *PriorityBlockingQueue[T]) OfferWithPriority(item T, priority int) error {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	err := q.queue.OfferWithPriority(item, priority)
	if err == nil {
		q.notEmpty.Signal()
	}
	return err
}

// Poll retrieves and removes the head of the queue, returning immediately if empty
func (q *PriorityBlockingQueue[T]) Poll() (T, error) {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	val, err := q.queue.Poll()
	if err == nil {
		q.notFull.Signal()
	}
	return val, err
}

// Peek retrieves but does not remove the head of the queue
func (q *PriorityBlockingQueue[T]) Peek() (T, error) {
	q.mutex.Lock()
	defer q.mutex.Unlock()
	return q.queue.Peek()
}

// Put adds an item to the queue, blocking if necessary until space is available
func (q *PriorityBlockingQueue[T]) Put(item T) error {
	return q.PutWithPriority(item, 0)
}

// PutWithPriority adds an item with priority to the queue, blocking if necessary
func (q *PriorityBlockingQueue[T]) PutWithPriority(item T, priority int) error {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	for q.queue.Len() >= q.queue.capacity {
		q.notFull.Wait()
	}

	err := q.queue.OfferWithPriority(item, priority)
	if err == nil {
		q.notEmpty.Signal()
	}
	return err
}

// Take retrieves and removes the head of the queue, blocking if empty
func (q *PriorityBlockingQueue[T]) Take() (T, error) {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	for q.queue.IsEmpty() {
		q.notEmpty.Wait()
	}

	val, err := q.queue.Poll()
	if err == nil {
		q.notFull.Signal()
	}
	return val, err
}

// Dump returns a slice containing all elements in the queue
func (q *PriorityBlockingQueue[T]) Dump() []T {
	q.mutex.Lock()
	defer q.mutex.Unlock()
	return q.queue.Dump()
}

// Size returns the number of elements in the queue
func (q *PriorityBlockingQueue[T]) Size() int {
	q.mutex.Lock()
	defer q.mutex.Unlock()
	return q.queue.Len()
}

// IsEmpty checks whether the queue is empty
func (q *PriorityBlockingQueue[T]) IsEmpty() bool {
	q.mutex.Lock()
	defer q.mutex.Unlock()
	return q.queue.IsEmpty()
}

// Clear removes all elements from the queue
func (q *PriorityBlockingQueue[T]) Clear() {
	q.mutex.Lock()
	defer q.mutex.Unlock()
	q.queue.Clear()
	q.notFull.Broadcast() // Signal that the queue is empty
}
