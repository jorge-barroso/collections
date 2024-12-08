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

	err := list.RemoveAt(1)
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

	err = list.RemoveAt(10)
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
