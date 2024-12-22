package queues

// Queue represents a generic queue with basic operations for elements of any type T.
// Offer inserts an item into the queue, returning an error if the operation fails.
// Poll retrieves and removes the head of the queue, returning an error if the queue is empty.
// Peek retrieves but does not remove the head of the queue, returning an error if the queue is empty.
// Dump returns a slice containing all elements in the queue.
// Size returns the current number of elements in the queue.
// IsEmpty returns true if the queue contains no elements.
// Clear removes all elements from the queue.
type Queue[T any] interface {
	Offer(item T) error
	Poll() (T, error)
	Peek() (T, error)
	Dump() []T
	Size() int
	IsEmpty() bool
	Clear()
}
