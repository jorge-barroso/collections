package queues

import (
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewPriorityBlockingQueue(t *testing.T) {
	queue := NewPriorityBlockingQueue[int](5)
	assert.NotNil(t, queue)
	assert.NotNil(t, queue.queue)
	assert.NotNil(t, queue.notEmpty)
	assert.NotNil(t, queue.notFull)
	assert.Equal(t, 0, queue.Size())
	assert.True(t, queue.IsEmpty())
}

func TestPriorityBlockingQueue_BasicOperations(t *testing.T) {
	tests := []struct {
		name     string
		capacity int
		items    []int
	}{
		{
			name:     "empty queue",
			capacity: 5,
			items:    []int{},
		},
		{
			name:     "single item",
			capacity: 5,
			items:    []int{1},
		},
		{
			name:     "multiple items",
			capacity: 5,
			items:    []int{1, 2, 3, 4, 5},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			queue := NewPriorityBlockingQueue[int](tt.capacity)

			// Test Offer and Size
			for _, item := range tt.items {
				err := queue.Offer(item)
				require.NoError(t, err)
			}
			assert.Equal(t, len(tt.items), queue.Size())

			// Test Peek
			if len(tt.items) > 0 {
				val, err := queue.Peek()
				require.NoError(t, err)
				assert.Equal(t, tt.items[0], val)
			}

			// Test Poll
			for i := 0; i < len(tt.items); i++ {
				val, err := queue.Poll()
				require.NoError(t, err)
				assert.Equal(t, tt.items[i], val)
			}

			assert.True(t, queue.IsEmpty())
		})
	}
}

func TestPriorityBlockingQueue_PriorityOperations(t *testing.T) {
	queue := NewPriorityBlockingQueue[int](5)

	// Add items with different priorities
	err := queue.OfferWithPriority(3, 3) // Lowest priority
	require.NoError(t, err)
	err = queue.OfferWithPriority(1, 1) // Highest priority
	require.NoError(t, err)
	err = queue.OfferWithPriority(2, 2) // Medium priority
	require.NoError(t, err)

	// Items should come out in priority order
	val, err := queue.Poll()
	require.NoError(t, err)
	assert.Equal(t, 1, val)

	val, err = queue.Poll()
	require.NoError(t, err)
	assert.Equal(t, 2, val)

	val, err = queue.Poll()
	require.NoError(t, err)
	assert.Equal(t, 3, val)
}

func TestPriorityBlockingQueue_BlockingOperations(t *testing.T) {
	queue := NewPriorityBlockingQueue[int](2)
	var wg sync.WaitGroup

	// Test Put blocking when queue is full
	err := queue.Put(1)
	require.NoError(t, err)
	err = queue.Put(2)
	require.NoError(t, err)

	wg.Add(1)
	go func() {
		defer wg.Done()
		// This should block until an item is removed
		err := queue.Put(3)
		assert.NoError(t, err)
	}()

	// Give some time for the goroutine to block
	time.Sleep(100 * time.Millisecond)

	// Remove an item to unblock Put
	val, err := queue.Poll()
	require.NoError(t, err)
	assert.Equal(t, 1, val)

	wg.Wait()
	assert.Equal(t, 2, queue.Size())

	// Test Take blocking when queue is empty
	queue.Clear()

	wg.Add(1)
	go func() {
		defer wg.Done()
		// This should block until an item is added
		val, err := queue.Take()
		assert.NoError(t, err)
		assert.Equal(t, 42, val)
	}()

	// Give some time for the goroutine to block
	time.Sleep(100 * time.Millisecond)

	// Add an item to unblock Take
	err = queue.Put(42)
	require.NoError(t, err)

	wg.Wait()
}

func TestPriorityBlockingQueue_ConcurrentOperations(t *testing.T) {
	queue := NewPriorityBlockingQueue[int](100)
	var wg sync.WaitGroup
	numGoroutines := 10
	itemsPerGoroutine := 100

	// Test concurrent Put operations
	wg.Add(numGoroutines)
	for i := 0; i < numGoroutines; i++ {
		go func(routine int) {
			defer wg.Done()
			for j := 0; j < itemsPerGoroutine; j++ {
				err := queue.Put(routine*itemsPerGoroutine + j)
				assert.NoError(t, err)
			}
		}(i)
	}
	wg.Wait()

	assert.Equal(t, numGoroutines*itemsPerGoroutine, queue.Size())

	// Test concurrent Take operations
	seen := make(map[int]bool)
	var mutex sync.Mutex
	wg.Add(numGoroutines)
	for i := 0; i < numGoroutines; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < itemsPerGoroutine; j++ {
				val, err := queue.Take()
				assert.NoError(t, err)
				mutex.Lock()
				seen[val] = true
				mutex.Unlock()
			}
		}()
	}
	wg.Wait()

	assert.Equal(t, numGoroutines*itemsPerGoroutine, len(seen))
	assert.True(t, queue.IsEmpty())
}

func TestPriorityBlockingQueue_Dump(t *testing.T) {
	queue := NewPriorityBlockingQueue[int](5)
	items := []int{1, 2, 3, 4, 5}

	for _, item := range items {
		err := queue.Put(item)
		require.NoError(t, err)
	}

	dumped := queue.Dump()
	assert.Equal(t, len(items), len(dumped))
	for i, item := range items {
		assert.Equal(t, item, dumped[i])
	}
}

func TestPriorityBlockingQueue_Clear(t *testing.T) {
	queue := NewPriorityBlockingQueue[int](5)
	items := []int{1, 2, 3, 4, 5}

	for _, item := range items {
		err := queue.Put(item)
		require.NoError(t, err)
	}

	queue.Clear()
	assert.True(t, queue.IsEmpty())
	assert.Equal(t, 0, queue.Size())

	// Test that we can still add items after clearing
	err := queue.Put(1)
	require.NoError(t, err)
	assert.Equal(t, 1, queue.Size())
}
