package lists

import (
	"testing"
)

func TestArrayList_Add(t *testing.T) {
	list := NewArrayList[int]()
	list.Add(1)
	list.Add(2)

	if list.Size() != 2 {
		t.Errorf("expected size 2, got %d", list.Size())
	}

	value, err := list.Get(0)
	if err != nil || value != 1 {
		t.Errorf("expected 1, got %d", value)
	}

	value, err = list.Get(1)
	if err != nil || value != 2 {
		t.Errorf("expected 2, got %d", value)
	}
}

func TestArrayList_RemoveAt(t *testing.T) {
	list := NewArrayList[int]()
	list.Add(1)
	list.Add(2)
	list.Add(3)

	err := list.Remove(1)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if list.Size() != 2 {
		t.Errorf("expected size 2, got %d", list.Size())
	}

	value, err := list.Get(1)
	if err != nil || value != 3 {
		t.Errorf("expected 3, got %d", value)
	}

	err = list.Remove(10)
	if err == nil {
		t.Errorf("expected error for out of bounds, got none")
	}
}

func TestArrayList_Get(t *testing.T) {
	list := NewArrayList[int]()
	list.Add(1)

	value, err := list.Get(0)
	if err != nil || value != 1 {
		t.Errorf("expected 1, got %d", value)
	}

	_, err = list.Get(10)
	if err == nil {
		t.Errorf("expected error for out of bounds, got none")
	}
}

func TestArrayList_Size(t *testing.T) {
	list := NewArrayList[int]()
	if list.Size() != 0 {
		t.Errorf("expected size 0, got %d", list.Size())
	}

	list.Add(1)
	if list.Size() != 1 {
		t.Errorf("expected size 1, got %d", list.Size())
	}
}

func TestArrayListIterator_Next(t *testing.T) {
	list := NewArrayList[int]()
	iter := list.NewIterator()

	if iter.Next() {
		t.Errorf("expected no next element, but found one")
	}

	list.Add(1)
	list.Add(2)
	iter = list.NewIterator()

	if !iter.Next() {
		t.Errorf("expected next element, but found none")
	}

	iter.Next() // move iterator to next

	if !iter.Next() {
		t.Errorf("expected next element, but found none")
	}
}

func TestArrayListIterator_Value(t *testing.T) {
	list := NewArrayList[int]()
	list.Add(1)
	list.Add(2)

	iter := list.NewIterator()

	value, err := iter.Value()
	if err != nil || value != 1 {
		t.Errorf("expected 1, got %d", value)
	}

	value, err = iter.Value()
	if err != nil || value != 2 {
		t.Errorf("expected 2, got %d", value)
	}

	// Attempt to get next element when none exist
	_, err = iter.Value()
	if err == nil {
		t.Errorf("expected error when no more elements, got none")
	}
}
