package collections

// PriorityElement wraps a value with its priority information for use in priority queues
type PriorityElement[T any] struct {
	Value    T
	Priority int
	Index    int // used by heap implementation
}

// NewPriorityElement creates a new PriorityElement with the given value and priority
func NewPriorityElement[T any](value T, priority int) *PriorityElement[T] {
	return &PriorityElement[T]{
		Value:    value,
		Priority: priority,
		Index:    -1,
	}
}
