package lists

import (
	"testing"
)

func TestCopyOnWriteList_Add(t *testing.T) {
	list := NewCopyOnWriteList[int]()
	list.Add(1)
	list.Add(2)

	if list.Size() != 2 {
		t.Errorf("expected size 2, got %d", list.Size())
	}

	value, ok := list.Get(0)
	if !ok || value != 1 {
		t.Errorf("expected 1, got %d", value)
	}

	value, ok = list.Get(1)
	if !ok || value != 2 {
		t.Errorf("expected 2, got %d", value)
	}
}

func TestCopyOnWriteList_Remove(t *testing.T) {
	list := NewCopyOnWriteList[int]()
	list.Add(1)
	list.Add(2)
	list.Add(3)

	list.Remove(1)
	if list.Size() != 2 {
		t.Errorf("expected size 2, got %d", list.Size())
	}

	value, ok := list.Get(1)
	if !ok || value != 3 {
		t.Errorf("expected 3, got %d", value)
	}

	list.Remove(10) // Removing an out-of-bound index
	if list.Size() != 2 {
		t.Errorf("expected size 2 after invalid removal, got %d", list.Size())
	}
}

func TestCopyOnWriteList_Get(t *testing.T) {
	list := NewCopyOnWriteList[int]()
	list.Add(1)

	value, ok := list.Get(0)
	if !ok || value != 1 {
		t.Errorf("expected 1, got %d", value)
	}

	_, ok = list.Get(10)
	if ok {
		t.Errorf("expected false for out-of-bounds index, got true")
	}
}

func TestCopyOnWriteList_Size(t *testing.T) {
	list := NewCopyOnWriteList[int]()
	if list.Size() != 0 {
		t.Errorf("expected size 0, got %d", list.Size())
	}

	list.Add(1)
	if list.Size() != 1 {
		t.Errorf("expected size 1, got %d", list.Size())
	}
}
