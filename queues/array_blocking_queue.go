package queues

import (
	"errors"
	"sync"
)

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

func (q *ArrayBlockingQueue[T]) Take() (T, error) {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	for q.count == 0 {
		q.notEmpty.Wait()
	}

	// Get the item
	item := q.items[q.head]

	// Clear the item from the queue for GC
	var zeroValue T // Initialize zero value for this type
	q.items[q.head] = zeroValue

	// Move the head
	q.head = (q.head + 1) % q.capacity
	q.count--

	// Notify that there is space available
	q.notFull.Signal()

	return item, nil
}

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

func (q *ArrayBlockingQueue[T]) Peek() (T, error) {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	var zeroValue T
	if q.count == 0 {
		return zeroValue, errors.New("queue empty")
	}

	return q.items[q.head], nil
}

func (q *ArrayBlockingQueue[T]) Dump() []T {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	// Prepare a slice to dump existing elements in order
	dump := make([]T, q.count)
	if q.count > 0 {
		if q.head < q.tail {
			copy(dump, q.items[q.head:q.tail])
		} else {
			// Case where the circular buffer wrapped around
			n := copy(dump, q.items[q.head:])
			copy(dump[n:], q.items[:q.tail])
		}
	}

	// Clear the queue
	var zeroValue T
	for i := range q.items {
		q.items[i] = zeroValue
	}
	q.head = 0
	q.tail = 0
	q.count = 0

	// Notify waiting goroutines that space is now available
	q.notFull.Broadcast()

	return dump
}
