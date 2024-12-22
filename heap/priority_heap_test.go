package heap

import (
	"container/heap"
	"github.com/jorge-barroso/collections"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPriorityHeap_HeapOperations(t *testing.T) {
	h := NewPriorityHeap[string](5)
	h.Initialize()

	// Test Push operation and heap ordering
	heap.Push(h, collections.NewPriorityElement("low", 1))
	heap.Push(h, collections.NewPriorityElement("high", 3))
	heap.Push(h, collections.NewPriorityElement("medium", 2))

	// Verify size
	assert.Equal(t, 3, h.Len())

	// Verify order through Pop operations
	elem1 := heap.Pop(h).(*collections.PriorityElement[string])
	assert.Equal(t, "high", elem1.Value)
	assert.Equal(t, 3, elem1.Priority)
	assert.Equal(t, -1, elem1.Index) // Index should be marked as removed

	elem2 := heap.Pop(h).(*collections.PriorityElement[string])
	assert.Equal(t, "medium", elem2.Value)
	assert.Equal(t, 2, elem2.Priority)

	elem3 := heap.Pop(h).(*collections.PriorityElement[string])
	assert.Equal(t, "low", elem3.Value)
	assert.Equal(t, 1, elem3.Priority)

	// Verify heap is empty
	assert.Equal(t, 0, h.Len())
}

func TestPriorityHeap_Fix(t *testing.T) {
	h := NewPriorityHeap[int](5)
	h.Initialize()

	// Add elements
	elem1 := collections.NewPriorityElement(1, 1)
	elem2 := collections.NewPriorityElement(2, 2)
	elem3 := collections.NewPriorityElement(3, 3)

	heap.Push(h, elem1)
	heap.Push(h, elem2)
	heap.Push(h, elem3)

	// Update priority of lowest element and fix its position
	elem1.Priority = 4
	h.Fix(elem1.Index)

	// Verify new order
	assert.Equal(t, elem1, (*h)[0], "Element with updated priority should be at the top")

	poppedElem := heap.Pop(h).(*collections.PriorityElement[int])
	assert.Equal(t, 1, poppedElem.Value)
	assert.Equal(t, 4, poppedElem.Priority)
}

func TestPriorityHeap_Clear(t *testing.T) {
	h := NewPriorityHeap[int](5)
	h.Initialize()

	// Add elements
	heap.Push(h, collections.NewPriorityElement(1, 1))
	heap.Push(h, collections.NewPriorityElement(2, 2))

	assert.Equal(t, 2, h.Len())

	// Clear heap
	h.Clear()

	assert.Equal(t, 0, h.Len())
	assert.Equal(t, 0, len(*h))

	// Verify we can still add elements after clearing
	heap.Push(h, collections.NewPriorityElement(3, 3))
	assert.Equal(t, 1, h.Len())
}

func TestPriorityHeap_SamePriority(t *testing.T) {
	h := NewPriorityHeap[string](5)
	h.Initialize()

	// Use a map to track remaining values
	remainingValues := map[string]bool{
		"first":  true,
		"second": true,
		"third":  true,
	}

	// Add elements with same priority
	heap.Push(h, collections.NewPriorityElement("first", 1))
	heap.Push(h, collections.NewPriorityElement("second", 1))
	heap.Push(h, collections.NewPriorityElement("third", 1))

	// Verify all elements can be retrieved and have same priority
	for i := 0; i < 3; i++ {
		elem := heap.Pop(h).(*collections.PriorityElement[string])
		assert.Equal(t, 1, elem.Priority)
		assert.True(t, remainingValues[elem.Value], "Value should be in remaining set")
		delete(remainingValues, elem.Value) // Remove the value we just found
	}

	assert.Equal(t, 0, h.Len())
	assert.Empty(t, remainingValues, "All values should have been found")
}
