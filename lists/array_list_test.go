package lists

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestArrayList_Add(t *testing.T) {
	list := NewArrayList[int]()
	list.Add(1)
	list.Add(2)

	assert.Equal(t, 2, list.Size(), "List size mismatch")

	value, err := list.Get(0)
	assert.NoError(t, err, "Unexpected error when getting value at index 0")
	assert.Equal(t, 1, value, "Value mismatch at index 0")

	value, err = list.Get(1)
	assert.NoError(t, err, "Unexpected error when getting value at index 1")
	assert.Equal(t, 2, value, "Value mismatch at index 1")
}

func TestArrayList_RemoveAt(t *testing.T) {
	list := NewArrayList[int]()
	list.Add(1)
	list.Add(2)
	list.Add(3)

	err := list.Remove(1)
	assert.NoError(t, err, "Unexpected error when removing element at index 1")
	assert.Equal(t, 2, list.Size(), "List size mismatch after removal")

	value, err := list.Get(1)
	assert.NoError(t, err, "Unexpected error when getting value at index 1")
	assert.Equal(t, 3, value, "Value mismatch at index 1 after removal")

	err = list.Remove(10)
	assert.Error(t, err, "Expected error for out-of-bounds removal, but got none")
}

func TestArrayList_Get(t *testing.T) {
	list := NewArrayList[int]()
	list.Add(1)

	value, err := list.Get(0)
	assert.NoError(t, err, "Unexpected error when getting value at index 0")
	assert.Equal(t, 1, value, "Value mismatch at index 0")

	_, err = list.Get(10)
	assert.Error(t, err, "Expected error for out-of-bounds Get, but got none")
}

func TestArrayList_Size(t *testing.T) {
	list := NewArrayList[int]()
	assert.Equal(t, 0, list.Size(), "Initial list size mismatch")

	list.Add(1)
	assert.Equal(t, 1, list.Size(), "List size mismatch after adding an element")
}

func TestArrayListIterator_Next(t *testing.T) {
	list := NewArrayList[int]()
	iter := list.NewIterator()

	assert.False(t, iter.Next(), "Iterator should have no next element in an empty list")

	list.Add(1)
	list.Add(2)
	iter = list.NewIterator()

	assert.True(t, iter.Next(), "Iterator should have a next element for the first entry")
	assert.True(t, iter.Next(), "Iterator should have a next element for the second entry")
	assert.False(t, iter.Next(), "Iterator should have no next element after the list ends")
}

func TestArrayListIterator_Value(t *testing.T) {
	list := NewArrayList[int]()
	list.Add(1)
	list.Add(2)

	iter := list.NewIterator()

	iter.Next()
	value, err := iter.Value()
	assert.NoError(t, err, "Unexpected error when getting the value of the iterator")
	assert.Equal(t, 1, value, "Value mismatch for the first iterator position")

	iter.Next() // Move to the next position
	value, err = iter.Value()
	assert.NoError(t, err, "Unexpected error when getting the value of the iterator")
	assert.Equal(t, 2, value, "Value mismatch for the second iterator position")

	assert.False(t, iter.Next()) // Move past the last position
	_, err = iter.Value()
	assert.Equal(t, 2, value, "Value mismatch for the second iterator position, Next() should not have moved past the last position")
}
