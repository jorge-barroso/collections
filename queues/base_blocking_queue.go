package queues

import (
	"sync"
)

// baseBlockingQueue provides common functionality for blocking queue implementations
type baseBlockingQueue[T any] struct {
	count    int         // Current number of elements
	capacity int         // Maximum capacity
	mutex    *sync.Mutex // Synchronization lock
	notFull  *sync.Cond  // Signaled when queue becomes not full
	notEmpty *sync.Cond  // Signaled when queue becomes not empty
}

// newBaseBlockingQueue creates a new abstract queue with the given capacity
func newBaseBlockingQueue[T any](capacity int) baseBlockingQueue[T] {
	mutex := &sync.Mutex{}
	return baseBlockingQueue[T]{
		capacity: capacity,
		mutex:    mutex,
		notFull:  sync.NewCond(mutex),
		notEmpty: sync.NewCond(mutex),
	}
}

// Lock acquires the mutex
func (q *baseBlockingQueue[T]) Lock() {
	q.mutex.Lock()
}

// Unlock releases the mutex
func (q *baseBlockingQueue[T]) Unlock() {
	q.mutex.Unlock()
}

// IsFull checks if the queue is at capacity
func (q *baseBlockingQueue[T]) IsFull() bool {
	return q.count == q.capacity
}

// IsEmpty checks if the queue is empty
func (q *baseBlockingQueue[T]) IsEmpty() bool {
	return q.count == 0
}

// WaitNotFull waits until the queue is not full
func (q *baseBlockingQueue[T]) WaitNotFull() {
	for q.IsFull() {
		q.notFull.Wait()
	}
}

// WaitNotEmpty waits until the queue is not empty
func (q *baseBlockingQueue[T]) WaitNotEmpty() {
	for q.IsEmpty() {
		q.notEmpty.Wait()
	}
}

// IncrementCount increases the count and signals waiting consumers
func (q *baseBlockingQueue[T]) IncrementCount() {
	q.count++
	q.notEmpty.Signal()
}

// DecrementCount decreases the count and signals waiting producers
func (q *baseBlockingQueue[T]) DecrementCount() {
	q.count--
	q.notFull.Signal()
}

// SignalAllNotFull broadcasts to all waiting producers
func (q *baseBlockingQueue[T]) SignalAllNotFull() {
	q.notFull.Broadcast()
}

// GetCount returns the current number of elements
func (q *baseBlockingQueue[T]) GetCount() int {
	return q.count
}

// GetCapacity returns the maximum capacity
func (q *baseBlockingQueue[T]) GetCapacity() int {
	return q.capacity
}

// CheckEmpty returns an error if the queue is empty
func (q *baseBlockingQueue[T]) CheckEmpty() error {
	if q.IsEmpty() {
		return errQueueEmpty
	}
	return nil
}

// CheckFull returns an error if the queue is full
func (q *baseBlockingQueue[T]) CheckFull() error {
	if q.IsFull() {
		return errQueueFull
	}
	return nil
}

// Reset resets the queue to empty state
func (q *baseBlockingQueue[T]) Reset() {
	q.count = 0
	q.SignalAllNotFull()
}
