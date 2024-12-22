package queues

// ArrayBlockingQueue is a thread-safe fixed size queue that blocks on full or
// empty conditions when adding or removing elements respectively.
type ArrayBlockingQueue[T any] struct {
	baseBlockingQueue[T]
	items []T
	head  int
	tail  int
}

var _ BlockingQueue[string] = &ArrayBlockingQueue[string]{}

// NewArrayBlockingQueue creates a new ArrayBlockingQueue with the specified capacity
func NewArrayBlockingQueue[T any](capacity int) *ArrayBlockingQueue[T] {
	return &ArrayBlockingQueue[T]{
		baseBlockingQueue: newBaseBlockingQueue[T](capacity),
		items:             make([]T, capacity),
	}
}

// Put adds an item to the tail of the queue, blocking if the queue is full
func (q *ArrayBlockingQueue[T]) Put(item T) error {
	q.Lock()
	defer q.Unlock()

	q.WaitNotFull()
	q.items[q.tail] = item
	q.tail = (q.tail + 1) % q.GetCapacity()
	q.IncrementCount()
	return nil
}

// Take retrieves and removes the item at the head of the queue, blocking if empty
func (q *ArrayBlockingQueue[T]) Take() (T, error) {
	q.Lock()
	defer q.Unlock()

	q.WaitNotEmpty()
	item := q.items[q.head]
	var zeroValue T
	q.items[q.head] = zeroValue
	q.head = (q.head + 1) % q.GetCapacity()
	q.DecrementCount()
	return item, nil
}

// Offer attempts to add an item to the queue without blocking
func (q *ArrayBlockingQueue[T]) Offer(item T) error {
	q.Lock()
	defer q.Unlock()

	if err := q.CheckFull(); err != nil {
		return err
	}

	q.items[q.tail] = item
	q.tail = (q.tail + 1) % q.GetCapacity()
	q.IncrementCount()
	return nil
}

// Poll retrieves and removes the head item without blocking
func (q *ArrayBlockingQueue[T]) Poll() (T, error) {
	q.Lock()
	defer q.Unlock()

	var zeroValue T
	if err := q.CheckEmpty(); err != nil {
		return zeroValue, err
	}

	item := q.items[q.head]
	q.items[q.head] = zeroValue
	q.head = (q.head + 1) % q.GetCapacity()
	q.DecrementCount()
	return item, nil
}

// Peek returns the head item without removing it
func (q *ArrayBlockingQueue[T]) Peek() (T, error) {
	q.Lock()
	defer q.Unlock()

	var zeroValue T
	if err := q.CheckEmpty(); err != nil {
		return zeroValue, err
	}

	return q.items[q.head], nil
}

// Dump returns a slice of all elements and clears the queue
func (q *ArrayBlockingQueue[T]) Dump() []T {
	q.Lock()
	defer q.Unlock()

	dump := make([]T, q.GetCount())
	if q.GetCount() > 0 {
		if q.head < q.tail {
			copy(dump, q.items[q.head:q.tail])
		} else {
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
	q.Reset()

	return dump
}

// Size returns the current number of elements in the queue
func (q *ArrayBlockingQueue[T]) Size() int {
	q.Lock()
	defer q.Unlock()

	return q.GetCount()
}

// Clear removes all elements from the queue
func (q *ArrayBlockingQueue[T]) Clear() {
	q.Lock()
	defer q.Unlock()

	// Clear the queue
	var zeroValue T
	for i := range q.items {
		q.items[i] = zeroValue
	}
	q.head = 0
	q.tail = 0
	q.Reset()
}
