package queues

import (
	"errors"
	"sync"
)

// BaseBlockingQueue provides common functionality for blocking queue implementations
type BaseBlockingQueue[T any] struct {
	count    int         // Current number of elements
	capacity int         // Maximum capacity
	mutex    *sync.Mutex // Synchronization lock
	notFull  *sync.Cond  // Signaled when queue becomes not full
	notEmpty *sync.Cond  // Signaled when queue becomes not empty
}

// NewBaseBlockingQueue creates a new abstract queue with the given capacity
func NewBaseBlockingQueue[T any](capacity int) BaseBlockingQueue[T] {
	mutex := &sync.Mutex{}
	return BaseBlockingQueue[T]{
		capacity: capacity,
		mutex:    mutex,
		notFull:  sync.NewCond(mutex),
		notEmpty: sync.NewCond(mutex),
	}
}

// Lock acquires the mutex
func (q *BaseBlockingQueue[T]) Lock() {
	q.mutex.Lock()
}

// Unlock releases the mutex
func (q *BaseBlockingQueue[T]) Unlock() {
	q.mutex.Unlock()
}

// IsFull checks if the queue is at capacity
func (q *BaseBlockingQueue[T]) IsFull() bool {
	return q.count == q.capacity
}

// IsEmpty checks if the queue is empty
func (q *BaseBlockingQueue[T]) IsEmpty() bool {
	return q.count == 0
}

// WaitNotFull waits until the queue is not full
func (q *BaseBlockingQueue[T]) WaitNotFull() {
	for q.IsFull() {
		q.notFull.Wait()
	}
}

// WaitNotEmpty waits until the queue is not empty
func (q *BaseBlockingQueue[T]) WaitNotEmpty() {
	for q.IsEmpty() {
		q.notEmpty.Wait()
	}
}

// IncrementCount increases the count and signals waiting consumers
func (q *BaseBlockingQueue[T]) IncrementCount() {
	q.count++
	q.notEmpty.Signal()
}

// DecrementCount decreases the count and signals waiting producers
func (q *BaseBlockingQueue[T]) DecrementCount() {
	q.count--
	q.notFull.Signal()
}

// SignalAllNotFull broadcasts to all waiting producers
func (q *BaseBlockingQueue[T]) SignalAllNotFull() {
	q.notFull.Broadcast()
}

// GetCount returns the current number of elements
func (q *BaseBlockingQueue[T]) GetCount() int {
	return q.count
}

// GetCapacity returns the maximum capacity
func (q *BaseBlockingQueue[T]) GetCapacity() int {
	return q.capacity
}

// CheckEmpty returns an error if the queue is empty
func (q *BaseBlockingQueue[T]) CheckEmpty() error {
	if q.IsEmpty() {
		return errors.New("queue empty")
	}
	return nil
}

// CheckFull returns an error if the queue is full
func (q *BaseBlockingQueue[T]) CheckFull() error {
	if q.IsFull() {
		return errors.New("queue full")
	}
	return nil
}

// Reset resets the queue to empty state
func (q *BaseBlockingQueue[T]) Reset() {
	q.count = 0
	q.SignalAllNotFull()
}