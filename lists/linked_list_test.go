package lists

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLinkedList_EmptyList(t *testing.T) {
	list := NewLinkedList[int]()

	// Test size
	assert.Equal(t, 0, list.Size(), "Expected size 0 for empty list")

	// Test Get
	_, err := list.Get(0)
	assert.Error(t, err, "Expected error when getting from empty list")

	// Test Remove
	err = list.Remove(0)
	assert.Error(t, err, "Expected error when removing from empty list")

	// Test Iterator
	iter := list.NewIterator()
	assert.False(t, iter.Next(), "Next() should return false for empty list")
	_, err = iter.Value()
	assert.Error(t, err, "Value() should return error for empty list")
}

func TestLinkedList_SingleElement(t *testing.T) {
	list := NewLinkedList[int]()
	list.Add(1)

	// Test size
	assert.Equal(t, 1, list.Size(), "Expected size 1 after adding one element")

	// Test Get
	value, err := list.Get(0)
	assert.NoError(t, err, "Unexpected error when getting value at index 0")
	assert.Equal(t, 1, value, "Value mismatch at index 0")

	// Test Iterator
	iter := list.NewIterator()
	assert.True(t, iter.Next(), "Next() should return true for first element")

	value, err = iter.Value()
	assert.NoError(t, err, "Unexpected error when getting iterator value")
	assert.Equal(t, 1, value, "Iterator value mismatch")

	assert.False(t, iter.Next(), "Next() should return false after last element")
}

func TestLinkedList_MultipleElements(t *testing.T) {
	list := NewLinkedList[int]()
	values := []int{1, 2, 3, 4, 5}

	// Add elements
	for _, v := range values {
		list.Add(v)
	}

	// Test size
	assert.Equal(t, len(values), list.Size(), "List size mismatch after adding elements")

	// Test Get
	for index, expected := range values {
		got, err := list.Get(index)
		assert.NoError(t, err, "Unexpected error when getting value at index %d", index)
		assert.Equal(t, expected, got, "Value mismatch at index %d", index)
	}

	// Test out of bounds
	_, err := list.Get(-1)
	assert.Error(t, err, "Expected error for Get(-1)")

	_, err = list.Get(len(values))
	assert.Error(t, err, "Expected error for Get(len)")
}

func TestLinkedList_Remove(t *testing.T) {
	list := NewLinkedList[int]()
	values := []int{1, 2, 3, 4, 5}
	for _, v := range values {
		list.Add(v)
	}

	// Remove first
	err := list.Remove(0)
	assert.NoError(t, err, "Unexpected error removing first element")
	value, _ := list.Get(0)
	assert.Equal(t, 2, value, "Expected 2 at index 0 after removing first element")

	// Remove last
	err = list.Remove(list.Size() - 1)
	assert.NoError(t, err, "Unexpected error removing last element")
	value, _ = list.Get(list.Size() - 1)
	assert.Equal(t, 4, value, "Expected 4 at last index after removing last element")

	// Remove middle
	err = list.Remove(1)
	assert.NoError(t, err, "Unexpected error removing middle element")
	expected := []int{2, 4}
	for i, want := range expected {
		got, _ := list.Get(i)
		assert.Equal(t, want, got, "Value mismatch at index %d after removing middle element", i)
	}
}

func TestLinkedListIterator_Behavior(t *testing.T) {
	list := NewLinkedList[int]()
	values := []int{1, 2, 3}

	for _, v := range values {
		list.Add(v)
	}

	iter := list.NewIterator()

	// First element
	assert.True(t, iter.Next(), "Next() should return true for the first element")
	value, err := iter.Value()
	assert.NoError(t, err, "Unexpected error at first element")
	assert.Equal(t, 1, value, "Value mismatch at first element")

	// Second element
	assert.True(t, iter.Next(), "Next() should return true for the second element")
	value, err = iter.Value()
	assert.NoError(t, err, "Unexpected error at second element")
	assert.Equal(t, 2, value, "Value mismatch at second element")

	// Third element
	assert.True(t, iter.Next(), "Next() should return true for the third element")
	value, err = iter.Value()
	assert.NoError(t, err, "Unexpected error at third element")
	assert.Equal(t, 3, value, "Value mismatch at third element")

	// After last element
	assert.False(t, iter.Next(), "Next() should return false after last element")
	value, err = iter.Value()
	assert.Equal(t, 3, value, "Value mismatch at third element, Next() should not have moved past the last position")
}

func TestLinkedListIterator_MultipleIterations(t *testing.T) {
	list := NewLinkedList[int]()
	list.Add(1)
	list.Add(2)

	iter := list.NewIterator()

	// First iteration
	count := 0
	for iter.Next() {
		count++
		_, err := iter.Value()
		assert.NoError(t, err, "Unexpected error during iteration")
	}
	assert.Equal(t, 2, count, "Iterator should cover all elements during iteration")

	// Attempting to continue iteration
	assert.False(t, iter.Next(), "Next() should return false after iteration completes")
	val, err := iter.Value()
	assert.Nil(t, err, "Value() should not return error after iteration ends, but remain in the last value")
	assert.Equal(t, 2, val, "Value mismatch for the second iterator position, Next() should not have moved past the last position")
}
