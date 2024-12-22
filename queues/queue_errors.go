package queues

import "errors"

var (
	// errQueueFull is returned when attempting to add to a queue that has reached its capacity
	errQueueFull = errors.New("queue is full")

	// errQueueEmpty is returned when attempting to retrieve from an empty queue
	errQueueEmpty = errors.New("queue is empty")
)
