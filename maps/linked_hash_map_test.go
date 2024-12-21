package maps

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLinkedHashMap_PutAndGet(t *testing.T) {
	m := NewLinkedHashMap[string, int]()

	// Test inserting a new key-value pair
	m.Put("a", 1)
	value, err := m.Get("a")
	assert.NoError(t, err, "Unexpected error when retrieving key 'a'")
	assert.Equal(t, 1, value, "Value mismatch for key 'a'")

	// Test updating an existing key-value pair
	m.Put("a", 2)
	value, err = m.Get("a")
	assert.NoError(t, err, "Unexpected error when updating key 'a'")
	assert.Equal(t, 2, value, "Updated value mismatch for key 'a'")

	// Test retrieving a nonexistent key
	_, err = m.Get("b")
	assert.Error(t, err, "Expected error when retrieving nonexistent key 'b'")
	assert.Equal(t, "key not found", err.Error(), "Unexpected error message when retrieving key 'b'")
}

func TestLinkedHashMap_Remove(t *testing.T) {
	m := NewLinkedHashMap[string, int]()

	// Add items to the map
	m.Put("a", 1)
	m.Put("b", 2)
	m.Put("c", 3)

	// Test removing an existing key
	err := m.Remove("b")
	assert.NoError(t, err, "Unexpected error when removing key 'b'")

	_, err = m.Get("b")
	assert.Error(t, err, "Expected error when retrieving removed key 'b'")

	// Test the order is preserved after removal
	expectedKeys := []string{"a", "c"}
	it := m.NewIterator()
	var keys []string
	for it.Next() {
		entry, err := it.Value()
		assert.NoError(t, err, "Unexpected error when iterating")
		keys = append(keys, entry.Key())
	}
	assert.Equal(t, expectedKeys, keys, "Key order mismatch after removal")

	// Test removing nonexistent key
	err = m.Remove("d")
	assert.Error(t, err, "Expected error when removing nonexistent key 'd'")

	// Remove remaining elements and check size
	assert.NoError(t, m.Remove("a"), "Unexpected error when removing remaining elements")
	assert.NoError(t, m.Remove("c"), "Unexpected error when removing remaining elements")
	assert.Equal(t, int64(0), m.Size(), "Expected size 0 after removing all elements")
}

func TestLinkedHashMap_Size(t *testing.T) {
	m := NewLinkedHashMap[string, int]()

	m.Put("a", 1)
	m.Put("b", 2)
	assert.Equal(t, int64(2), m.Size(), "Size mismatch after adding elements")

	assert.NoError(t, m.Remove("b"), "Unexpected error when removing an element")
	assert.Equal(t, int64(1), m.Size(), "Size mismatch after removing an element")

	m.Put("a", 5) // Updating a key should not increase size
	assert.Equal(t, int64(1), m.Size(), "Size mismatch after updating an existing key")
}

func TestLinkedHashMap_NewIterator(t *testing.T) {
	m := NewLinkedHashMap[string, int]()

	// Test empty map
	it := m.NewIterator()
	assert.False(t, it.Next(), "Expected no elements in iterator for empty map")

	// Add elements
	m.Put("a", 1)
	m.Put("b", 2)
	m.Put("c", 3)

	// Iterate through map
	expectedKeys := []string{"a", "b", "c"}
	expectedValues := []int{1, 2, 3}

	var keys []string
	var values []int
	it = m.NewIterator()
	for it.Next() {
		entry, err := it.Value()
		assert.NoError(t, err, "Unexpected error during iteration")
		keys = append(keys, entry.Key())
		values = append(values, entry.Value())
	}

	assert.Equal(t, expectedKeys, keys, "Key order mismatch in iterator")
	assert.Equal(t, expectedValues, values, "Value order mismatch in iterator")
}

func TestLinkedHashMapIterator_Empty(t *testing.T) {
	m := NewLinkedHashMap[string, int]()
	it := m.NewIterator()

	assert.False(t, it.Next(), "Expected Next() to return false for empty map")

	_, err := it.Value()
	assert.Error(t, err, "Expected error when calling Value() on empty map iterator")
	assert.Equal(t, "no more elements", err.Error(), "Unexpected error message for empty map iterator")
}

func TestLinkedHashMapIterator_SingleElement(t *testing.T) {
	m := NewLinkedHashMap[string, int]()
	m.Put("a", 1)

	it := m.NewIterator()

	assert.True(t, it.Next(), "Expected Next() to return true for single element")
	entry, err := it.Value()
	assert.NoError(t, err, "Unexpected error when retrieving single element from iterator")
	assert.Equal(t, "a", entry.Key(), "Key mismatch for single element in iterator")
	assert.Equal(t, 1, entry.Value(), "Value mismatch for single element in iterator")

	assert.False(t, it.Next(), "Expected Next() to return false after the last element")
	_, err = it.Value()
	assert.Error(t, err, "Expected error when retrieving value after iterator end")
	assert.Equal(t, "no more elements", err.Error(), "Unexpected error message after iterator end")
}
