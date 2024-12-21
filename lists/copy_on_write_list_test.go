package lists

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCopyOnWriteList_BasicOperations(t *testing.T) {
	list := NewCopyOnWriteList[int]()

	// Test initial state
	assert.Equal(t, 0, list.Size(), "Expected empty list")

	// Test Add
	list.Add(1)
	assert.Equal(t, 1, list.Size(), "Expected size 1 after adding one element")

	// Test Get
	value, err := list.Get(0)
	assert.NoError(t, err, "Unexpected error while getting value at index 0")
	assert.Equal(t, 1, value, "Value mismatch at index 0")

	// Test Get with invalid index
	_, err = list.Get(-1)
	assert.Error(t, err, "Expected error for negative index")

	_, err = list.Get(1)
	assert.Error(t, err, "Expected error for out of bounds index")
}

func TestCopyOnWriteList_Remove(t *testing.T) {
	list := NewCopyOnWriteList[int]()
	list.Add(1)
	list.Add(2)
	list.Add(3)

	// Test Remove middle element
	err := list.Remove(1)
	assert.NoError(t, err, "Unexpected error removing middle element")
	assert.Equal(t, 2, list.Size(), "List size mismatch after removal")
	value, _ := list.Get(1)
	assert.Equal(t, 3, value, "Expected 3 at index 1 after removal")

	// Test Remove first element
	err = list.Remove(0)
	assert.NoError(t, err, "Unexpected error removing first element")
	value, _ = list.Get(0)
	assert.Equal(t, 3, value, "Expected 3 at index 0 after removal")

	// Test Remove last element
	err = list.Remove(0)
	assert.NoError(t, err, "Unexpected error removing last element")
	assert.Equal(t, 0, list.Size(), "Expected empty list after last removal")

	// Test Remove on empty list
	err = list.Remove(0)
	assert.Error(t, err, "Expected error removing from empty list")
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
	assert.GreaterOrEqual(t, list.Size(), iterations, "Expected at least one thread to successfully write elements")
}

func TestCopyOnWriteList_Iterator(t *testing.T) {
	list := NewCopyOnWriteList[int]()

	// Test iterator on empty list
	iter := list.NewIterator()
	assert.False(t, iter.Next(), "Next() should return false for empty list")
	_, err := iter.Value()
	assert.Error(t, err, "Value() should return an error for an empty list")

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
		assert.NoError(t, err, "Unexpected error during iteration")
		assert.Equal(t, values[index], v, "Value mismatch during iteration")
		index++
	}
	assert.Equal(t, len(values), index, "Iterator covered fewer elements than expected")

	// Test iterator snapshot isolation
	iter = list.NewIterator()
	list.Add(6) // Modify list after creating iterator
	count := 0
	for iter.Next() {
		count++
	}
	assert.Equal(t, len(values), count, "Iterator should see only the snapshot from when it was created")
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
	assert.Equal(t, 2, len(values), "Expected iteration to cover two values")

	// Attempting to continue iteration
	assert.False(t, iter.Next(), "Next() should return false after iteration is complete")
	val, err := iter.Value()
	assert.NoError(t, err, "Value() should not return an error after iteration but hold the last value")
	assert.Equal(t, 2, val, "Value() mismatch after iteration completion")
}

func TestCopyOnWriteList_IteratorConcurrency(t *testing.T) {
	list := NewCopyOnWriteList[int]()

	// Add elements before testing
	for i := 0; i < 100; i++ {
		list.Add(i)
	}

	iter := list.NewIterator()
	var wg sync.WaitGroup

	// Concurrent iteration
	wg.Add(1)
	go func() {
		defer wg.Done()
		for iter.Next() {
			_, _ = iter.Value() // Ignore errors in this test, focus on concurrency
		}
	}()

	// Concurrent modification
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < 50; i++ {
			list.Add(i)
		}
	}()

	wg.Wait()
	assert.GreaterOrEqual(t, list.Size(), 100, "List size should reflect concurrent modifications")
	// Iterator behavior under concurrency is undefined for most cases so no strict assertions
}
