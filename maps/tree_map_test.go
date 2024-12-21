package maps

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTreeMap_PutAndGet(t *testing.T) {
	// Simple comparison function for integers
	less := func(a, b int) bool { return a < b }

	tm := NewTreeMap[int, string](less)

	// Test adding a new key-value pair
	tm.Put(1, "A")
	value, err := tm.Get(1)
	assert.NoError(t, err, "Unexpected error when getting key 1")
	assert.Equal(t, "A", value, "Value mismatch for key 1")

	// Test updating an existing key-value pair
	tm.Put(1, "B")
	value, err = tm.Get(1)
	assert.NoError(t, err, "Unexpected error when updating key 1")
	assert.Equal(t, "B", value, "Updated value mismatch for key 1")

	// Test retrieving a non-existent key
	_, err = tm.Get(2)
	assert.Error(t, err, "Expected error when getting non-existent key 2")
	assert.Equal(t, "key not found", err.Error(), "Unexpected error message for non-existent key")
}

func TestTreeMap_Remove(t *testing.T) {
	less := func(a, b int) bool { return a < b }
	tm := NewTreeMap[int, string](less)

	// Add items to TreeMap
	tm.Put(1, "A")
	tm.Put(2, "B")
	tm.Put(3, "C")

	// Test removing an existing key
	err := tm.Remove(2)
	assert.NoError(t, err, "Unexpected error when removing key 2")

	_, err = tm.Get(2)
	assert.Error(t, err, "Expected error when accessing removed key 2")

	// Test removing a non-existent key
	err = tm.Remove(4)
	assert.Error(t, err, "Expected error when removing non-existent key 4")
	assert.Equal(t, "key not found", err.Error(), "Unexpected error message for non-existent key")

	// Test removing all elements
	assert.NoError(t, tm.Remove(1), "Unexpected error when removing first element")
	assert.NoError(t, tm.Remove(3), "Unexpected error when removing last element")
	assert.Equal(t, int64(0), tm.Size(), "Expected size 0 after removing all elements")
}

func TestTreeMap_Size(t *testing.T) {
	less := func(a, b int) bool { return a < b }
	tm := NewTreeMap[int, string](less)

	// Initially, the size should be zero
	assert.Equal(t, int64(0), tm.Size(), "Expected initial size to be 0")

	// Add elements and check size
	tm.Put(1, "A")
	tm.Put(2, "B")
	assert.Equal(t, int64(2), tm.Size(), "Size mismatch after adding elements")

	// Remove an element and check size
	assert.NoError(t, tm.Remove(1), "Unexpected error when removing an element")
	assert.Equal(t, int64(1), tm.Size(), "Size mismatch after removing an element")

	// Replacing a value (same key) should not change size
	tm.Put(2, "C")
	assert.Equal(t, int64(1), tm.Size(), "Size mismatch after updating existing key")
}

func TestTreeMap_NewIterator(t *testing.T) {
	less := func(a, b int) bool { return a < b }
	tm := NewTreeMap[int, string](less)

	// Add elements (not in sorted order) to the map
	tm.Put(3, "C")
	tm.Put(1, "A")
	tm.Put(2, "B")

	// Check that iterator returns keys in sorted order
	it := tm.NewIterator()
	expectedKeys := []int{1, 2, 3}
	expectedValues := []string{"A", "B", "C"}

	i := 0
	for it.Next() {
		entry, err := it.Value()
		assert.NoError(t, err, "Unexpected error during iteration")
		assert.Equal(t, expectedKeys[i], entry.Key(), "Key mismatch during iteration")
		assert.Equal(t, expectedValues[i], entry.Value(), "Value mismatch during iteration")
		i++
	}

	assert.Equal(t, len(expectedKeys), i, "Iterator did not return all elements")
}

func TestTreeMapIterator_Empty(t *testing.T) {
	less := func(a, b int) bool { return a < b }
	tm := NewTreeMap[int, string](less)

	it := NewTreeMapIterator[int, string](tm)

	// Check that iterator has no elements
	assert.False(t, it.Next(), "Expected no elements in iterator for empty map")

	_, err := it.Value()
	assert.Error(t, err, "Expected error when calling Value() on empty iterator")
	assert.Equal(t, "no more elements", err.Error(), "Unexpected error message for empty map iterator")
}
