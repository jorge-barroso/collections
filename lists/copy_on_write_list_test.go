package lists

import (
	"sync"
	"testing"
)

func TestCopyOnWriteList_BasicOperations(t *testing.T) {
	list := NewCopyOnWriteList[int]()

	// Test initial state
	if list.Size() != 0 {
		t.Errorf("expected empty list, got size %d", list.Size())
	}

	// Test Add
	list.Add(1)
	if size := list.Size(); size != 1 {
		t.Errorf("expected size 1, got %d", size)
	}

	// Test Get
	value, err := list.Get(0)
	if err != nil || value != 1 {
		t.Errorf("expected 1, got %d with error: %v", value, err)
	}

	// Test Get with invalid index
	_, err = list.Get(-1)
	if err == nil {
		t.Error("expected error for negative index")
	}

	_, err = list.Get(1)
	if err == nil {
		t.Error("expected error for out of bounds index")
	}
}

func TestCopyOnWriteList_Remove(t *testing.T) {
	list := NewCopyOnWriteList[int]()
	list.Add(1)
	list.Add(2)
	list.Add(3)

	// Test Remove middle element
	err := list.Remove(1)
	if err != nil {
		t.Errorf("unexpected error on remove: %v", err)
	}
	if size := list.Size(); size != 2 {
		t.Errorf("expected size 2, got %d", size)
	}
	if v, _ := list.Get(1); v != 3 {
		t.Errorf("expected 3 at index 1, got %d", v)
	}

	// Test Remove first element
	err = list.Remove(0)
	if err != nil {
		t.Errorf("unexpected error on remove: %v", err)
	}
	if v, _ := list.Get(0); v != 3 {
		t.Errorf("expected 3 at index 0, got %d", v)
	}

	// Test Remove last element
	err = list.Remove(0)
	if err != nil {
		t.Errorf("unexpected error on remove: %v", err)
	}
	if size := list.Size(); size != 0 {
		t.Errorf("expected empty list, got size %d", size)
	}

	// Test Remove on empty list
	err = list.Remove(0)
	if err == nil {
		t.Error("expected error removing from empty list")
	}
}

func TestCopyOnWriteList_ConcurrentModification(t *testing.T) {
	list := NewCopyOnWriteList[int]()
	const iterations = 1000
	var wg sync.WaitGroup

	// Concurrent writers
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < iterations; j++ {
				list.Add(j)
			}
		}()
	}

	// Concurrent readers
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < iterations; j++ {
				list.Get(0) // Ignore errors as size changes
			}
		}()
	}

	wg.Wait()
}

func TestCopyOnWriteList_Iterator(t *testing.T) {
	list := NewCopyOnWriteList[int]()

	// Test iterator on empty list
	iter := list.NewIterator()
	if iter.Next() {
		t.Error("Next() should return false for empty list")
	}
	_, err := iter.Value()
	if err == nil {
		t.Error("Value() should return error for empty list")
	}

	// Add elements
	values := []int{1, 2, 3, 4, 5}
	for _, v := range values {
		list.Add(v)
	}

	// Test basic iteration
	iter = list.NewIterator()
	index := 0
	for iter.Next() {
		v, err := iter.Value()
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if v != values[index] {
			t.Errorf("at index %d: expected %d, got %d", index, values[index], v)
		}
		index++
	}
	if index != len(values) {
		t.Errorf("iterator covered %d elements, expected %d", index, len(values))
	}

	// Test iterator snapshot isolation
	iter = list.NewIterator()
	list.Add(6) // Modify list after creating iterator
	count := 0
	for iter.Next() {
		count++
	}
	if count != len(values) {
		t.Errorf("iterator saw %d elements, expected %d", count, len(values))
	}
}

func TestCopyOnWriteList_IteratorReuse(t *testing.T) {
	list := NewCopyOnWriteList[int]()
	list.Add(1)
	list.Add(2)

	iter := list.NewIterator()

	// First iteration
	values := make([]int, 0)
	for iter.Next() {
		v, _ := iter.Value()
		values = append(values, v)
	}
	if len(values) != 2 {
		t.Errorf("expected 2 values, got %d", len(values))
	}

	// Attempting to continue iteration
	if iter.Next() {
		t.Error("Next() should return false after iteration is complete")
	}
	if val, err := iter.Value(); err != nil || val != 2 {
		t.Error("Value() should not return error after iteration is complete but just stay in the last value")
	}
}

func TestCopyOnWriteList_IteratorConcurrency(t *testing.T) {
	list := NewCopyOnWriteList[int]()
	for i := 0; i < 100; i++ {
		list.Add(i)
	}

	var wg sync.WaitGroup
	// Create multiple iterators concurrently
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			iter := list.NewIterator()
			count := 0
			for iter.Next() {
				_, err := iter.Value()
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				count++
			}
			// Each iterator sees exactly the snapshot at creation time
			if count < 100 {
				t.Errorf("iterator saw %d elements, expected at least 100", count)
			}
		}()
	}

	// Modify list while iterators are running
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < 100; i++ {
			list.Add(i)
			if i%10 == 0 {
				list.Remove(0)
			}
		}
	}()

	wg.Wait()
}
