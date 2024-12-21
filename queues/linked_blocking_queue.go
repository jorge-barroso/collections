package queues

import "github.com/jorge-barroso/collections"

// LinkedBlockingQueue is a thread-safe queue with a fixed capacity that uses linked nodes.
type LinkedBlockingQueue[T any] struct {
	BaseBlockingQueue[T]
	head *collections.Node[T] // Points to the first Node in the queue
	tail *collections.Node[T] // Points to the last Node in the queue
}

// Ensure LinkedBlockingQueue implements both Map and Iterable interfaces
var _ BlockingQueue[int] = (*LinkedBlockingQueue[int])(nil)

// NewLinkedBlockingQueue creates a new LinkedBlockingQueue with the specified capacity
func NewLinkedBlockingQueue[T any](capacity int) *LinkedBlockingQueue[T] {
	return &LinkedBlockingQueue[T]{
		BaseBlockingQueue: NewBaseBlockingQueue[T](capacity),
	}
}

// Put inserts the specified element into the queue, blocking if necessary
func (q *LinkedBlockingQueue[T]) Put(item T) error {
	q.Lock()
	defer q.Unlock()

	q.WaitNotFull()
	newNode := &collections.Node[T]{Item: item}
	if q.tail != nil {
		q.tail.Next = newNode
	} else {
		q.head = newNode
	}
	q.tail = newNode
	q.IncrementCount()
	return nil
}

// Take removes and returns the head of the queue, blocking if necessary
func (q *LinkedBlockingQueue[T]) Take() (T, error) {
	q.Lock()
	defer q.Unlock()

	q.WaitNotEmpty()
	item := q.head.Item
	if q.head.Next != nil {
		q.head = q.head.Next
	} else {
		q.head = nil
		q.tail = nil
	}
	q.DecrementCount()
	return item, nil
}

// Offer attempts to add the specified element to the queue without blocking
func (q *LinkedBlockingQueue[T]) Offer(item T) error {
	q.Lock()
	defer q.Unlock()

	if err := q.CheckFull(); err != nil {
		return err
	}

	newNode := &collections.Node[T]{Item: item}
	if q.tail != nil {
		q.tail.Next = newNode
	} else {
		q.head = newNode
	}
	q.tail = newNode
	q.IncrementCount()
	return nil
}

// Poll removes and returns the head of the queue without blocking
func (q *LinkedBlockingQueue[T]) Poll() (T, error) {
	q.Lock()
	defer q.Unlock()

	var zeroValue T
	if err := q.CheckEmpty(); err != nil {
		return zeroValue, err
	}

	item := q.head.Item
	if q.head.Next != nil {
		q.head = q.head.Next
	} else {
		q.head = nil
		q.tail = nil
	}
	q.DecrementCount()
	return item, nil
}

// Peek returns the head of the queue without removing it
func (q *LinkedBlockingQueue[T]) Peek() (T, error) {
	q.Lock()
	defer q.Unlock()

	var zeroValue T
	if err := q.CheckEmpty(); err != nil {
		return zeroValue, err
	}

	return q.head.Item, nil
}

// Dump returns a slice of all elements and clears the queue
func (q *LinkedBlockingQueue[T]) Dump() []T {
	q.Lock()
	defer q.Unlock()

	values := make([]T, 0, q.GetCount())
	for current := q.head; current != nil; current = current.Next {
		values = append(values, current.Item)
	}

	// Reset the queue
	q.head = nil
	q.tail = nil
	q.Reset()

	return values
}
