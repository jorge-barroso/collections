package queues

// BlockingQueue extends Queue with blocking operations
type BlockingQueue[T any] interface {
	Put(item T) error
	Take() (T, error)
	Queue[T] // Embedding the Queue interface
}
