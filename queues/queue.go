package queues

// Queue defines standard queue operations
type Queue[T any] interface {
	Offer(item T) error
	Poll() (T, error)
	Peek() (T, error)
	Dump() []T
}
