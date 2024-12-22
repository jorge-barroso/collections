package queues

import (
	"container/heap"
	"github.com/jorge-barroso/collections"
	customheap "github.com/jorge-barroso/collections/heap"
)

// PriorityQueue implements a non-thread-safe priority queue
type PriorityQueue[T any] struct {
	heap     *customheap.PriorityHeap[T]
	capacity int
}

// NewPriorityQueue creates a new priority queue with the specified capacity
func NewPriorityQueue[T any](capacity int) *PriorityQueue[T] {
	pq := &PriorityQueue[T]{
		heap:     customheap.NewPriorityHeap[T](capacity),
		capacity: capacity,
	}
	pq.heap.Initialize()
	return pq
}

// Offer adds an item to the queue with default priority (0)
func (pq *PriorityQueue[T]) Offer(value T) error {
	return pq.OfferWithPriority(value, 0)
}

// OfferWithPriority adds an item with the specified priority
func (pq *PriorityQueue[T]) OfferWithPriority(value T, priority int) error {
	if pq.heap.Len() >= pq.capacity {
		return errQueueFull
	}

	element := collections.NewPriorityElement(value, priority)
	heap.Push(pq.heap, element)
	return nil
}

// Poll removes and returns the highest priority element
func (pq *PriorityQueue[T]) Poll() (T, error) {
	var zero T
	if pq.heap.Len() == 0 {
		return zero, errQueueEmpty
	}

	element := heap.Pop(pq.heap).(*collections.PriorityElement[T])
	return element.Value, nil
}

// Peek returns the highest priority element without removing it
func (pq *PriorityQueue[T]) Peek() (T, error) {
	var zero T
	if pq.heap.Len() == 0 {
		return zero, errQueueEmpty
	}

	return (*pq.heap)[0].Value, nil
}

// Dump returns a slice containing all elements in priority order
func (pq *PriorityQueue[T]) Dump() []T {
	// Create a copy of the heap to not modify the original
	tempHeap := make([]*collections.PriorityElement[T], pq.heap.Len())
	copy(tempHeap, *pq.heap)
	tempPriorityHeap := customheap.PriorityHeap[T](tempHeap)
	tempPriorityHeap.Initialize()

	// Extract elements one by one in priority order
	length := len(tempHeap)
	result := make([]T, length)
	for i := 0; i < length; i++ {
		element := heap.Pop(&tempPriorityHeap).(*collections.PriorityElement[T])
		result[i] = element.Value
	}
	return result
}

// Clear removes all elements from the queue
func (pq *PriorityQueue[T]) Clear() {
	pq.heap.Clear()
}

// Len returns the current number of elements in the queue
func (pq *PriorityQueue[T]) Len() int {
	return pq.heap.Len()
}

// IsEmpty returns true if the queue has no elements
func (pq *PriorityQueue[T]) IsEmpty() bool {
	return pq.Len() == 0
}
