package heap

import (
	"container/heap"
	"github.com/jorge-barroso/collections"
)

// PriorityHeap implements heap.Interface for priority elements
type PriorityHeap[T any] []*collections.PriorityElement[T]

var _ heap.Interface = &PriorityHeap[int]{}

// Len returns the number of elements in the heap
func (h *PriorityHeap[T]) Len() int {
	return len(*h)
}

// Less defines the ordering of elements
func (h *PriorityHeap[T]) Less(i, j int) bool {
	return (*h)[i].Priority > (*h)[j].Priority // Higher priority comes first
}

// Swap exchanges elements at the given indices
func (h *PriorityHeap[T]) Swap(i, j int) {
	(*h)[i], (*h)[j] = (*h)[j], (*h)[i]
	(*h)[i].Index = i
	(*h)[j].Index = j
}

// Push adds an element to the heap
func (h *PriorityHeap[T]) Push(x any) {
	element := x.(*collections.PriorityElement[T])
	element.Index = len(*h)
	*h = append(*h, element)
}

// Pop removes and returns the last element
func (h *PriorityHeap[T]) Pop() any {
	old := *h
	n := len(old)
	element := old[n-1]
	old[n-1] = nil     // avoid memory leak
	element.Index = -1 // mark as removed
	*h = old[:n-1]
	return element
}

// Fix updates the position of the element at index i
func (h *PriorityHeap[T]) Fix(i int) {
	heap.Fix(h, i)
}

// Clear removes all elements from the heap
func (h *PriorityHeap[T]) Clear() {
	*h = (*h)[:0]
}

// Initialize sets up the heap
func (h *PriorityHeap[T]) Initialize() {
	heap.Init(h)
}

// NewPriorityHeap creates a new heap adapter with the given capacity
func NewPriorityHeap[T any](capacity int) *PriorityHeap[T] {
	h := make(PriorityHeap[T], 0, capacity)
	return &h
}
