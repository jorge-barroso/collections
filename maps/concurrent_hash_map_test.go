package maps

import (
	"fmt"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConcurrentHashMap_BasicOperations(t *testing.T) {
	cm := NewConcurrentHashMap[string, int]()

	// Test Put and Get
	cm.Put("one", 1)
	cm.Put("two", 2)

	value, err := cm.Get("one")
	assert.NoError(t, err, "Unexpected error when getting key 'one'")
	assert.Equal(t, 1, value, "Value mismatch for key 'one'")

	value, err = cm.Get("two")
	assert.NoError(t, err, "Unexpected error when getting key 'two'")
	assert.Equal(t, 2, value, "Value mismatch for key 'two'")

	// Test Size
	assert.Equal(t, int64(2), cm.Size(), "Size mismatch after adding elements")

	// Test Remove
	err = cm.Remove("one")
	assert.NoError(t, err, "Unexpected error when removing key 'one'")
	assert.Equal(t, int64(1), cm.Size(), "Size mismatch after removing key 'one'")

	// Test ContainsKey
	assert.True(t, cm.ContainsKey("two"), "Expected map to contain key 'two'")
	assert.False(t, cm.ContainsKey("one"), "Expected map to not contain key 'one' after removal")
}

func TestConcurrentHashMap_ConcurrentAccess(t *testing.T) {
	cm := NewConcurrentHashMap[int, int]()
	var wg sync.WaitGroup
	numGoroutines := 10
	numOperations := 1000

	// Concurrent writes
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(base int) {
			defer wg.Done()
			for j := 0; j < numOperations; j++ {
				cm.Put(base+j, base+j)
			}
		}(i * numOperations)
	}
	wg.Wait()

	// Concurrent reads
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(base int) {
			defer wg.Done()
			for j := 0; j < numOperations; j++ {
				_, _ = cm.Get(base + j)
			}
		}(i * numOperations)
	}
	wg.Wait()

	// Verify size
	expectedSize := int64(numGoroutines * numOperations)
	assert.Equal(t, expectedSize, cm.Size(), "Size mismatch after concurrent operations")
}

func TestConcurrentHashMap_Iterator(t *testing.T) {
	cm := NewConcurrentHashMap[string, int]()
	testData := map[string]int{
		"one":   1,
		"two":   2,
		"three": 3,
		"four":  4,
	}

	// Add test data
	for k, v := range testData {
		cm.Put(k, v)
	}

	// Test iterator
	iter := cm.NewIterator()
	count := 0
	foundValues := make(map[string]bool)

	for iter.Next() {
		entry, err := iter.Value()
		assert.NoError(t, err, "Unexpected error during iteration")
		count++
		key := entry.Key()
		value := entry.Value()

		assert.Contains(t, testData, key, "Unexpected key encountered in iteration")
		assert.Equal(t, testData[key], value, "Value mismatch for key %s", key)
		foundValues[key] = true
	}

	assert.Equal(t, len(testData), count, "Iterator visited fewer elements than expected")

	for k := range testData {
		assert.Contains(t, foundValues, k, "Iterator missed expected key: %s", k)
	}
}

func TestConcurrentHashMap_UpdateExistingKey(t *testing.T) {
	cm := NewConcurrentHashMap[string, int]()

	cm.Put("key", 10)
	value, err := cm.Get("key")
	assert.NoError(t, err, "Unexpected error when getting key 'key'")
	assert.Equal(t, 10, value, "Value mismatch before updating key 'key'")

	cm.Put("key", 20)
	value, err = cm.Get("key")
	assert.NoError(t, err, "Unexpected error when getting updated key 'key'")
	assert.Equal(t, 20, value, "Value mismatch after updating key 'key'")

	assert.Equal(t, int64(1), cm.Size(), "Size should remain 1 after overwriting key")
}

func TestConcurrentHashMap_RemoveNonExistentKey(t *testing.T) {
	cm := NewConcurrentHashMap[string, int]()
	err := cm.Remove("nonexistent")
	assert.Error(t, err, "Expected error when removing non-existent key")
}

func TestConcurrentHashMap_ConcurrentReadsAndWrites(t *testing.T) {
	cm := NewConcurrentHashMap[string, int]()
	var wg sync.WaitGroup

	// Concurrent write operations
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(base int) {
			defer wg.Done()
			for j := 0; j < 100; j++ {
				cm.Put(fmt.Sprintf("key%d", base+j), base+j)
			}
		}(i * 100)
	}

	// Concurrent read operations
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(base int) {
			defer wg.Done()
			for j := 0; j < 100; j++ {
				_, _ = cm.Get(fmt.Sprintf("key%d", base+j))
			}
		}(i * 100)
	}

	wg.Wait()

	// Verify size
	expectedSize := int64(1000)
	assert.Equal(t, expectedSize, cm.Size(), "Size mismatch after concurrent reads and writes")
}
