package queues

// PriorityQueueInterface represents a generic priority queue allowing operations with priority.
// OfferWithPriority adds an item with the specified priority to the priority queue.
// Embeds Queue to inherit common queue operations.
type PriorityQueueInterface[T any] interface {
	// OfferWithPriority adds an item with the specified priority
	OfferWithPriority(item T, priority int) error
	Queue[T]
}
