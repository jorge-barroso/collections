package queues

import (
	"sync"
	"testing"
	"time"
)

// TestBaseBlockingQueue_New tests the constructor
func TestBaseBlockingQueue_New(t *testing.T) {
	tests := []struct {
		name     string
		capacity int
	}{
		{"zero capacity", 0},
		{"positive capacity", 5},
		{"large capacity", 1000},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			queue := NewBaseBlockingQueue[int](tt.capacity)
			if queue.capacity != tt.capacity {
				t.Errorf("NewBaseBlockingQueue(%d) capacity = %d; want %d",
					tt.capacity, queue.capacity, tt.capacity)
			}
			if queue.count != 0 {
				t.Errorf("NewBaseBlockingQueue(%d) count = %d; want 0",
					tt.capacity, queue.count)
			}
			if queue.mutex == nil {
				t.Error("NewBaseBlockingQueue() mutex is nil")
			}
			if queue.notFull == nil {
				t.Error("NewBaseBlockingQueue() notFull condition is nil")
			}
			if queue.notEmpty == nil {
				t.Error("NewBaseBlockingQueue() notEmpty condition is nil")
			}
		})
	}
}

// TestBaseBlockingQueue_Lock tests mutex operations
func TestBaseBlockingQueue_Lock(t *testing.T) {
	queue := NewBaseBlockingQueue[int](5)

	// Test basic lock/unlock
	queue.Lock()
	queue.Unlock()

	// Test concurrent access
	var wg sync.WaitGroup
	concurrent := 10
	counter := 0

	for i := 0; i < concurrent; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			queue.Lock()
			counter++
			queue.Unlock()
		}()
	}

	wg.Wait()
	if counter != concurrent {
		t.Errorf("Counter = %d; want %d", counter, concurrent)
	}
}

// TestBaseBlockingQueue_IsFull tests capacity checks
func TestBaseBlockingQueue_IsFull(t *testing.T) {
	tests := []struct {
		name     string
		capacity int
		count    int
		wantFull bool
		isEmpty  bool
	}{
		{"empty queue", 5, 0, false, true},
		{"partially filled", 5, 3, false, false},
		{"full queue", 5, 5, true, false},
		{"zero capacity queue", 0, 0, true, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			queue := NewBaseBlockingQueue[int](tt.capacity)
			queue.count = tt.count

			if got := queue.IsFull(); got != tt.wantFull {
				t.Errorf("IsFull() = %v; want %v", got, tt.wantFull)
			}

			if got := queue.IsEmpty(); got != tt.isEmpty {
				t.Errorf("IsEmpty() = %v; want %v", got, tt.isEmpty)
			}
		})
	}
}

// TestBaseBlockingQueue_WaitNotFull tests waiting for space
func TestBaseBlockingQueue_WaitNotFull(t *testing.T) {
	queue := NewBaseBlockingQueue[int](1)
	done := make(chan bool)

	queue.count = queue.capacity
	go func() {
		queue.Lock()
		queue.WaitNotFull()
		queue.Unlock()
		done <- true
	}()

	// Signal after a short delay
	time.Sleep(10 * time.Millisecond)
	queue.Lock()
	queue.count--
	queue.notFull.Signal()
	queue.Unlock()

	select {
	case <-done:
		// Test passed
	case <-time.After(100 * time.Millisecond):
		t.Error("WaitNotFull() timed out")
	}
}

// TestBaseBlockingQueue_WaitNotEmpty tests waiting for items
func TestBaseBlockingQueue_WaitNotEmpty(t *testing.T) {
	queue := NewBaseBlockingQueue[int](1)
	done := make(chan bool)

	queue.count = 0
	go func() {
		queue.Lock()
		queue.WaitNotEmpty()
		queue.Unlock()
		done <- true
	}()

	// Signal after a short delay
	time.Sleep(10 * time.Millisecond)
	queue.Lock()
	queue.count++
	queue.notEmpty.Signal()
	queue.Unlock()

	select {
	case <-done:
		// Test passed
	case <-time.After(100 * time.Millisecond):
		t.Error("WaitNotEmpty() timed out")
	}
}

// TestBaseBlockingQueue_IncrementCount tests increment operation
func TestBaseBlockingQueue_IncrementCount(t *testing.T) {
	queue := NewBaseBlockingQueue[int](5)
	initial := queue.count

	queue.Lock()
	queue.IncrementCount()
	queue.Unlock()

	if queue.count != initial+1 {
		t.Errorf("IncrementCount() count = %d; want %d", queue.count, initial+1)
	}
}

// TestBaseBlockingQueue_DecrementCount tests decrement operation
func TestBaseBlockingQueue_DecrementCount(t *testing.T) {
	queue := NewBaseBlockingQueue[int](5)
	queue.Lock()
	queue.count = 2
	queue.DecrementCount()
	queue.Unlock()

	if queue.count != 1 {
		t.Errorf("DecrementCount() count = %d; want 1", queue.count)
	}
}

// TestBaseBlockingQueue_CheckEmpty tests empty state validation
func TestBaseBlockingQueue_CheckEmpty(t *testing.T) {
	queue := NewBaseBlockingQueue[int](2)

	if err := queue.CheckEmpty(); err == nil {
		t.Error("CheckEmpty() on empty queue should return error")
	}

	queue.count = 1
	if err := queue.CheckEmpty(); err != nil {
		t.Error("CheckEmpty() on non-empty queue should not return error")
	}
}

// TestBaseBlockingQueue_CheckFull tests full state validation
func TestBaseBlockingQueue_CheckFull(t *testing.T) {
	queue := NewBaseBlockingQueue[int](2)
	queue.count = queue.capacity

	if err := queue.CheckFull(); err == nil {
		t.Error("CheckFull() on full queue should return error")
	}

	queue.count = 1
	if err := queue.CheckFull(); err != nil {
		t.Error("CheckFull() on non-full queue should not return error")
	}
}

// TestBaseBlockingQueue_Reset tests queue reset
func TestBaseBlockingQueue_Reset(t *testing.T) {
	queue := NewBaseBlockingQueue[int](5)
	queue.count = 3

	queue.Reset()

	if queue.count != 0 {
		t.Errorf("Reset() count = %d; want 0", queue.count)
	}

	// Test that Reset signals waiting producers
	done := make(chan bool)
	go func() {
		queue.Lock()
		queue.WaitNotFull()
		queue.Unlock()
		done <- true
	}()

	select {
	case <-done:
		// Test passed
	case <-time.After(100 * time.Millisecond):
		t.Error("Reset() did not signal waiting producers")
	}
}

// TestBaseBlockingQueue_GetCount tests count retrieval
func TestBaseBlockingQueue_GetCount(t *testing.T) {
	queue := NewBaseBlockingQueue[int](5)
	queue.count = 3

	if got := queue.GetCount(); got != 3 {
		t.Errorf("GetCount() = %d; want 3", got)
	}
}

// TestBaseBlockingQueue_GetCapacity tests capacity retrieval
func TestBaseBlockingQueue_GetCapacity(t *testing.T) {
	capacity := 5
	queue := NewBaseBlockingQueue[int](capacity)

	if got := queue.GetCapacity(); got != capacity {
		t.Errorf("GetCapacity() = %d; want %d", got, capacity)
	}
}
