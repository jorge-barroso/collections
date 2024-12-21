package queues

import (
	"errors"
	"github.com/jorge-barroso/collections"
	"sync"
)

// LinkedBlockingQueue is a thread-safe queue with a fixed capacity that uses linked nodes.
type LinkedBlockingQueue[T any] struct {
	head     *collections.Node[T] // Points to the first Node in the queue
	tail     *collections.Node[T] // Points to the last Node in the queue
	count    int                  // Number of items currently in the queue
	capacity int                  // Maximum queue capacity
	mutex    *sync.Mutex          // Mutex for synchronizing access
	notFull  *sync.Cond           // Condition variable for signaling space availability
	notEmpty *sync.Cond           // Condition variable for signaling item availability
}

// Ensure LinkedBlockingQueue implements both Map and Iterable interfaces
var _ BlockingQueue[int] = (*LinkedBlockingQueue[int])(nil)

// NewLinkedBlockingQueue creates a new LinkedBlockingQueue with the specified capacity.
func NewLinkedBlockingQueue[T any](capacity int) *LinkedBlockingQueue[T] {
	mutex := &sync.Mutex{}
	return &LinkedBlockingQueue[T]{
		capacity: capacity,
		mutex:    mutex,
		notFull:  sync.NewCond(mutex),
		notEmpty: sync.NewCond(mutex),
	}
}

// Put inserts the specified element into the queue, blocking if necessary for space to become available.
// Blocks if the queue is full.
func (q *LinkedBlockingQueue[T]) Put(item T) error {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	// Wait until space is available in the queue
	for q.count == q.capacity {
		q.notFull.Wait()
	}

	// Add item to the queue
	newNode := &collections.Node[T]{Item: item}
	if q.tail != nil {
		q.tail.Next = newNode
	} else {
		q.head = newNode
	}
	q.tail = newNode
	q.count++
	q.notEmpty.Signal()
	return nil
}

// Take removes and returns the head of the queue, blocking if necessary until an element becomes available.
// Blocks if the queue is empty.
func (q *LinkedBlockingQueue[T]) Take() (T, error) {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	// Wait until an item is available in the queue
	for q.count == 0 {
		q.notEmpty.Wait()
	}

	// Remove and return the head item
	removedNode := q.head
	if q.head.Next != nil {
		q.head = q.head.Next
	} else {
		q.head = nil
		q.tail = nil
	}
	q.count--
	q.notFull.Signal()
	return removedNode.Item, nil
}

// Offer attempts to add the specified element to the queue without blocking.
// Returns an error if the queue is full.
func (q *LinkedBlockingQueue[T]) Offer(item T) error {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	// Check if space is available in the queue
	if q.count == q.capacity {
		return errors.New("queue full")
	}

	// Add item to the queue
	newNode := &collections.Node[T]{Item: item}
	if q.tail != nil {
		q.tail.Next = newNode
	} else {
		q.head = newNode
	}
	q.tail = newNode
	q.count++
	q.notEmpty.Signal()
	return nil
}

// Poll removes and returns the head of the queue, or returns an error if the queue is empty.
// Does not block if the queue is empty.
func (q *LinkedBlockingQueue[T]) Poll() (T, error) {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	var zeroValue T
	// Check if the queue is empty
	if q.count == 0 {
		return zeroValue, errors.New("queue empty")
	}

	// Remove and return the head item
	removedNode := q.head
	if q.head.Next != nil {
		q.head = q.head.Next
	} else {
		q.head = nil
		q.tail = nil
	}
	q.count--
	q.notFull.Signal()
	return removedNode.Item, nil
}

// Peek returns the head of the queue without removing it, or returns an error if the queue is empty.
// Does not block if the queue is empty.
func (q *LinkedBlockingQueue[T]) Peek() (T, error) {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	var zeroValue T
	// Check if the queue is empty
	if q.count == 0 {
		return zeroValue, errors.New("queue empty")
	}

	// Return the head item without removing it
	return q.head.Item, nil
}

// Dump clears the queue and returns all its elements in order.
// Notifies all waiting goroutines that space is now available.
func (q *LinkedBlockingQueue[T]) Dump() []T {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	// Collect current values in the queue
	values := make([]T, 0, q.count)
	for current := q.head; current != nil; current = current.Next {
		values = append(values, current.Item)
	}

	// Reset the queue
	q.head = nil
	q.tail = nil
	q.count = 0

	// Notify waiting goroutines that space is now available
	q.notFull.Broadcast()

	return values
}
