package queues

// BlockingQueue is a thread-safe generic queue with blocking operations for adding and removing elements.
// Put inserts an item into the queue, blocking if full.
// Take retrieves and removes the head item, blocking if empty.
// Embeds Queue to inherit common queue operations.
type BlockingQueue[T any] interface {
	Put(item T) error
	Take() (T, error)
	Queue[T] // Embedding the Queue interface
}
