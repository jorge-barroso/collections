package queues

import (
	"errors"
	"sync"
)

// ArrayBlockingQueue is a thread-safe fixed size queue that blocks on full or
// empty conditions when adding or removing elements respectively.
type ArrayBlockingQueue[T any] struct {
	items    []T
	head     int
	tail     int
	count    int
	capacity int
	mutex    *sync.Mutex
	notFull  *sync.Cond
	notEmpty *sync.Cond
}

// Ensure ArrayBlockingQueue implements both Map and Iterable interfaces
var _ BlockingQueue[int] = (*ArrayBlockingQueue[int])(nil)

// NewArrayBlockingQueue creates a new ArrayBlockingQueue with the specified capacity.
// The queue blocks operations when it is full or empty.
func NewArrayBlockingQueue[T any](capacity int) *ArrayBlockingQueue[T] {
	mutex := &sync.Mutex{}
	return &ArrayBlockingQueue[T]{
		items:    make([]T, capacity),
		capacity: capacity,
		mutex:    mutex,
		notFull:  sync.NewCond(mutex),
		notEmpty: sync.NewCond(mutex),
	}
}

// Put adds an item to the tail of the queue, blocking if the queue is full
// until space becomes available.
func (q *ArrayBlockingQueue[T]) Put(item T) error {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	for q.count == q.capacity {
		q.notFull.Wait()
	}

	q.items[q.tail] = item
	q.tail = (q.tail + 1) % q.capacity
	q.count++
	q.notEmpty.Signal()
	return nil
}

// Take retrieves and removes the item at the head of the queue, blocking
// if the queue is empty until an item becomes available.
func (q *ArrayBlockingQueue[T]) Take() (T, error) {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	for q.count == 0 {
		q.notEmpty.Wait()
	}

	item := q.items[q.head]
	var zeroValue T
	q.items[q.head] = zeroValue
	q.head = (q.head + 1) % q.capacity
	q.count--
	q.notFull.Signal()
	return item, nil
}

// Offer attempts to add an item to the queue without blocking. Returns an
// error if the queue is full.
func (q *ArrayBlockingQueue[T]) Offer(item T) error {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	if q.count == q.capacity {
		return errors.New("queue full")
	}

	q.items[q.tail] = item
	q.tail = (q.tail + 1) % q.capacity
	q.count++
	q.notEmpty.Signal()
	return nil
}

// Poll retrieves and removes the item at the head of the queue without
// blocking. Returns an error if the queue is empty.
func (q *ArrayBlockingQueue[T]) Poll() (T, error) {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	var zeroValue T
	if q.count == 0 {
		return zeroValue, errors.New("queue empty")
	}

	item := q.items[q.head]
	q.head = (q.head + 1) % q.capacity
	q.count--
	q.notFull.Signal()
	return item, nil
}

// Peek returns the item at the head of the queue without removing it.
// Returns an error if the queue is empty.
func (q *ArrayBlockingQueue[T]) Peek() (T, error) {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	var zeroValue T
	if q.count == 0 {
		return zeroValue, errors.New("queue empty")
	}

	return q.items[q.head], nil
}

// Dump returns a slice of all the current elements in the queue and clears it.
// This operation is thread-safe.
func (q *ArrayBlockingQueue[T]) Dump() []T {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	dump := make([]T, q.count)
	if q.count > 0 {
		if q.head < q.tail {
			copy(dump, q.items[q.head:q.tail])
		} else {
			n := copy(dump, q.items[q.head:])
			copy(dump[n:], q.items[:q.tail])
		}
	}

	var zeroValue T
	for i := range q.items {
		q.items[i] = zeroValue
	}
	q.head = 0
	q.tail = 0
	q.count = 0

	q.notFull.Broadcast()

	return dump
}
