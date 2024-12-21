package queues

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestArrayBlockingQueue_PutAndTake(t *testing.T) {
	queue := NewArrayBlockingQueue[int](3)

	// Test basic put and take operations
	go func() {
		assert.NoError(t, queue.Put(1), "Unexpected error on Put")
		assert.NoError(t, queue.Put(2), "Unexpected error on Put")
		assert.NoError(t, queue.Put(3), "Unexpected error on Put")
	}()

	time.Sleep(100 * time.Millisecond) // Allow goroutine to execute

	value, err := queue.Take()
	assert.NoError(t, err, "Unexpected error on Take")
	assert.Equal(t, 1, value, "Value mismatch on Take")

	value, err = queue.Take()
	assert.NoError(t, err, "Unexpected error on Take")
	assert.Equal(t, 2, value, "Value mismatch on Take")

	value, err = queue.Take()
	assert.NoError(t, err, "Unexpected error on Take")
	assert.Equal(t, 3, value, "Value mismatch on Take")
}

func TestArrayBlockingQueue_OfferAndPoll(t *testing.T) {
	queue := NewArrayBlockingQueue[int](2)

	assert.NoError(t, queue.Offer(1), "Unexpected error on Offer")
	assert.NoError(t, queue.Offer(2), "Unexpected error on Offer")

	// Should return an error since the queue is full
	err := queue.Offer(3)
	assert.Error(t, err, "Expected error on Offer when full")

	value, err := queue.Poll()
	assert.NoError(t, err, "Unexpected error on Poll")
	assert.Equal(t, 1, value, "Value mismatch on Poll")

	value, err = queue.Poll()
	assert.NoError(t, err, "Unexpected error on Poll")
	assert.Equal(t, 2, value, "Value mismatch on Poll")

	// Should return an error since the queue is empty
	_, err = queue.Poll()
	assert.Error(t, err, "Expected error on Poll when empty")
}

func TestArrayBlockingQueue_Peek(t *testing.T) {
	queue := NewArrayBlockingQueue[int](2)

	_, err := queue.Peek()
	assert.Error(t, err, "Expected error on Peek when empty")

	assert.NoError(t, queue.Offer(1), "Unexpected error on Offer")

	value, err := queue.Peek()
	assert.NoError(t, err, "Unexpected error on Peek")
	assert.Equal(t, 1, value, "Value mismatch on Peek")

	assert.NoError(t, queue.Offer(2), "Unexpected error on Offer")

	value, err = queue.Peek()
	assert.NoError(t, err, "Unexpected error on Peek")
	assert.Equal(t, 1, value, "Expected Peek to return the same first value")
}

func TestArrayBlockingQueue_Dump(t *testing.T) {
	queue := NewArrayBlockingQueue[int](5)

	for i := 1; i <= 3; i++ {
		assert.NoError(t, queue.Put(i), "Unexpected error on Put")
	}

	dumpedItems := queue.Dump()
	expected := []int{1, 2, 3}

	assert.Equal(t, expected, dumpedItems, "Dumped items mismatch")

	_, err := queue.Poll()
	assert.Error(t, err, "Expected error on Poll after dump")
}
