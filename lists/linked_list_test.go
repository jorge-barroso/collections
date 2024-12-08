package lists

import (
	"testing"
)

func TestLinkedList_Add(t *testing.T) {
	list := NewLinkedList[int]()
	list.Add(1)
	list.Add(2)
	list.Add(3)

	if list.Size() != 3 {
		t.Errorf("expected size 3, got %d", list.Size())
	}

	value, err := list.Get(0)
	if err != nil || value != 1 {
		t.Errorf("expected 1, got %d", value)
	}

	value, err = list.Get(1)
	if err != nil || value != 2 {
		t.Errorf("expected 2, got %d", value)
	}

	value, err = list.Get(2)
	if err != nil || value != 3 {
		t.Errorf("expected 3, got %d", value)
	}
}

func TestLinkedList_Remove(t *testing.T) {
	list := NewLinkedList[int]()
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

	err = list.Remove(5) // Attempt out-of-bounds removal
	if err == nil {
		t.Errorf("expected error for out-of-bounds removal, got none")
	}
}

func TestLinkedList_Get(t *testing.T) {
	list := NewLinkedList[int]()
	list.Add(1)
	list.Add(2)

	value, err := list.Get(0)
	if err != nil || value != 1 {
		t.Errorf("expected 1, got %d", value)
	}

	value, err = list.Get(1)
	if err != nil || value != 2 {
		t.Errorf("expected 2, got %d", value)
	}

	_, err = list.Get(2) // Index out of bounds
	if err == nil {
		t.Errorf("expected error for out-of-bounds access, got none")
	}
}

func TestLinkedList_Size(t *testing.T) {
	list := NewLinkedList[int]()
	if list.Size() != 0 {
		t.Errorf("expected size 0, got %d", list.Size())
	}

	list.Add(1)
	if list.Size() != 1 {
		t.Errorf("expected size 1, got %d", list.Size())
	}
}
