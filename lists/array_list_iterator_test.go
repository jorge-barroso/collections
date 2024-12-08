package lists

import (
	"testing"
)

func TestArrayListIterator_HasNext(t *testing.T) {
	list := NewArrayList[int]()
	iter := list.NewIterator()

	if iter.HasNext() {
		t.Errorf("expected no next element, but found one")
	}

	list.Add(1)
	list.Add(2)
	iter = list.NewIterator()

	if !iter.HasNext() {
		t.Errorf("expected next element, but found none")
	}

	iter.Next() // move iterator to next

	if !iter.HasNext() {
		t.Errorf("expected next element, but found none")
	}
}

func TestArrayListIterator_Next(t *testing.T) {
	list := NewArrayList[int]()
	list.Add(1)
	list.Add(2)

	iter := list.NewIterator()

	value, err := iter.Next()
	if err != nil || value != 1 {
		t.Errorf("expected 1, got %d", value)
	}

	value, err = iter.Next()
	if err != nil || value != 2 {
		t.Errorf("expected 2, got %d", value)
	}

	// Attempt to get next element when none exist
	_, err = iter.Next()
	if err == nil {
		t.Errorf("expected error when no more elements, got none")
	}
}
