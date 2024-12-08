package queues

import (
	"testing"
	"time"
)

func TestArrayBlockingQueue_PutAndTake(t *testing.T) {
	queue := NewArrayBlockingQueue[int](3)

	// Test basic put and take operations
	go func() {
		if err := queue.Put(1); err != nil {
			t.Errorf("unexpected error on Put: %v", err)
		}
		if err := queue.Put(2); err != nil {
			t.Errorf("unexpected error on Put: %v", err)
		}
		if err := queue.Put(3); err != nil {
			t.Errorf("unexpected error on Put: %v", err)
		}
	}()

	time.Sleep(100 * time.Millisecond) // Allow some time for goroutine to execute

	value, err := queue.Take()
	if err != nil || value != 1 {
		t.Errorf("expected 1, got %d", value)
	}

	value, err = queue.Take()
	if err != nil || value != 2 {
		t.Errorf("expected 2, got %d", value)
	}

	value, err = queue.Take()
	if err != nil || value != 3 {
		t.Errorf("expected 3, got %d", value)
	}
}

func TestArrayBlockingQueue_OfferAndPoll(t *testing.T) {
	queue := NewArrayBlockingQueue[int](2)

	if err := queue.Offer(1); err != nil {
		t.Errorf("unexpected error on Offer: %v", err)
	}

	if err := queue.Offer(2); err != nil {
		t.Errorf("unexpected error on Offer: %v", err)
	}

	// Should return an error since the queue is full
	if err := queue.Offer(3); err == nil {
		t.Errorf("expected error on Offer when full, got none")
	}

	value, err := queue.Poll()
	if err != nil || value != 1 {
		t.Errorf("expected 1, got %d", value)
	}

	value, err = queue.Poll()
	if err != nil || value != 2 {
		t.Errorf("expected 2, got %d", value)
	}

	// Should return an error since the queue is empty
	if _, err := queue.Poll(); err == nil {
		t.Errorf("expected error on Poll when empty, got none")
	}
}

func TestArrayBlockingQueue_Peek(t *testing.T) {
	queue := NewArrayBlockingQueue[int](2)

	if _, err := queue.Peek(); err == nil {
		t.Errorf("expected error on Peek when empty, got none")
	}

	if err := queue.Offer(1); err != nil {
		t.Errorf("unexpected error on Offer: %v", err)
	}

	value, err := queue.Peek()
	if err != nil || value != 1 {
		t.Errorf("expected 1, got %d", value)
	}

	if err := queue.Offer(2); err != nil {
		t.Errorf("unexpected error on Offer: %v", err)
	}

	value, err = queue.Peek()
	if err != nil || value != 1 {
		t.Errorf("expected 1 again, got %d", value)
	}
}

func TestArrayBlockingQueue_Dump(t *testing.T) {
	queue := NewArrayBlockingQueue[int](5)

	for i := 1; i <= 3; i++ {
		if err := queue.Put(i); err != nil {
			t.Errorf("unexpected error on Put: %v", err)
		}
	}

	dumpedItems := queue.Dump()
	expected := []int{1, 2, 3}

	if len(dumpedItems) != len(expected) {
		t.Errorf("expected dump length %d, got %d", len(expected), len(dumpedItems))
	}

	for i, v := range dumpedItems {
		if v != expected[i] {
			t.Errorf("expected %d at index %d, got %d", expected[i], i, v)
		}
	}

	if _, err := queue.Poll(); err == nil {
		t.Errorf("expected error on Poll after dump, got none")
	}
}
