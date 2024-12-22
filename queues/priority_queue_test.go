package queues

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPriorityQueue_Offer(t *testing.T) {
	q := NewPriorityQueue[string](3)

	// Test basic offer
	assert.NoError(t, q.Offer("first"))
	assert.Equal(t, 1, q.Len())
	assert.False(t, q.IsEmpty())

	// Test offer with priority
	assert.NoError(t, q.OfferWithPriority("second", priorityHigh))
	assert.Equal(t, 2, q.Len())

	// Test capacity limit
	assert.NoError(t, q.Offer("third"))
	assert.ErrorIs(t, q.Offer("fourth"), errQueueFull)
}

func TestPriorityQueue_Poll(t *testing.T) {
	q := NewPriorityQueue[string](5)

	// Test poll empty queue
	_, err := q.Poll()
	assert.ErrorIs(t, err, errQueueEmpty)

	// Add items with different priorities
	assert.NoError(t, q.OfferWithPriority("low", priorityLow))
	assert.NoError(t, q.OfferWithPriority("high", priorityHigh))
	assert.NoError(t, q.OfferWithPriority("medium", priorityMiddle))

	// Test poll order
	val, err := q.Poll()
	assert.NoError(t, err)
	assert.Equal(t, "high", val)

	val, err = q.Poll()
	assert.NoError(t, err)
	assert.Equal(t, "medium", val)

	val, err = q.Poll()
	assert.NoError(t, err)
	assert.Equal(t, "low", val)

	assert.True(t, q.IsEmpty())
}

func TestPriorityQueue_Peek(t *testing.T) {
	q := NewPriorityQueue[string](5)

	// Test peek empty queue
	_, err := q.Peek()
	assert.ErrorIs(t, err, errQueueEmpty)

	// Add items
	assert.NoError(t, q.OfferWithPriority("low", priorityLow))
	assert.NoError(t, q.OfferWithPriority("high", priorityHigh))

	// Test peek
	val, err := q.Peek()
	assert.NoError(t, err)
	assert.Equal(t, "high", val)

	// Verify peek doesn't remove element
	assert.Equal(t, 2, q.Len())
}

func TestPriorityQueue_Clear(t *testing.T) {
	q := NewPriorityQueue[int](5)

	// Add items
	assert.NoError(t, q.Offer(1))
	assert.NoError(t, q.Offer(2))
	assert.Equal(t, 2, q.Len())

	// Clear queue
	q.Clear()
	assert.Equal(t, 0, q.Len())
	assert.True(t, q.IsEmpty())

	// Verify can still add items after clearing
	assert.NoError(t, q.Offer(3))
	assert.Equal(t, 1, q.Len())
}

func TestPriorityQueue_Dump(t *testing.T) {
	q := NewPriorityQueue[int](5)

	// Test empty dump
	dump := q.Dump()
	assert.Empty(t, dump)

	// Add items with different priorities
	assert.NoError(t, q.OfferWithPriority(1, priorityLow))
	assert.NoError(t, q.OfferWithPriority(2, priorityHigh))
	assert.NoError(t, q.OfferWithPriority(3, priorityMiddle))

	// Test dump
	dump = q.Dump()
	assert.Len(t, dump, 3)

	// Verify original queue is unchanged
	assert.Equal(t, 3, q.Len())

	// Verify the items are in priority order
	assert.Equal(t, 2, dump[0]) // Highest priority
	assert.Equal(t, 3, dump[1])
	assert.Equal(t, 1, dump[2]) // Lowest priority
}
