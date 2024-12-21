package lists

import (
	"testing"
)

func TestLinkedList_EmptyList(t *testing.T) {
	list := NewLinkedList[int]()

	if list.Size() != 0 {
		t.Errorf("expected size 0, got %d", list.Size())
	}

	_, err := list.Get(0)
	if err == nil {
		t.Error("expected error getting from empty list")
	}

	err = list.Remove(0)
	if err == nil {
		t.Error("expected error removing from empty list")
	}

	iter := list.NewIterator()
	if iter.Next() {
		t.Error("Next() should return false for empty list")
	}
	_, err = iter.Value()
	if err == nil {
		t.Error("Value() should error on empty list")
	}
}

func TestLinkedList_SingleElement(t *testing.T) {
	list := NewLinkedList[int]()
	list.Add(1)

	if list.Size() != 1 {
		t.Errorf("expected size 1, got %d", list.Size())
	}

	value, err := list.Get(0)
	if err != nil || value != 1 {
		t.Errorf("Get(0) = %v, %v; want 1, nil", value, err)
	}

	iter := list.NewIterator()

	if !iter.Next() {
		t.Errorf("Next() should return true for first element")
	}

	value, err = iter.Value()
	if err != nil || value != 1 {
		t.Errorf("Value() = %v, %v; want 1, nil", value, err)
	}
	if iter.Next() {
		t.Error("Next() should return false after last element")
	}
}

func TestLinkedList_MultipleElements(t *testing.T) {
	list := NewLinkedList[int]()
	values := []int{1, 2, 3, 4, 5}

	for _, v := range values {
		list.Add(v)
	}

	if list.Size() != len(values) {
		t.Errorf("expected size %d, got %d", len(values), list.Size())
	}

	// Test Get
	for i, want := range values {
		got, err := list.Get(i)
		if err != nil || got != want {
			t.Errorf("Get(%d) = %v, %v; want %v, nil", i, got, err, want)
		}
	}

	// Test out of bounds
	_, err := list.Get(-1)
	if err == nil {
		t.Error("expected error for Get(-1)")
	}
	_, err = list.Get(len(values))
	if err == nil {
		t.Error("expected error for Get(len)")
	}
}

func TestLinkedList_Remove(t *testing.T) {
	list := NewLinkedList[int]()
	values := []int{1, 2, 3, 4, 5}
	for _, v := range values {
		list.Add(v)
	}

	// Remove first
	err := list.Remove(0)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if v, _ := list.Get(0); v != 2 {
		t.Errorf("after removing first, Get(0) = %v; want 2", v)
	}

	// Remove last
	err = list.Remove(list.Size() - 1)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if v, _ := list.Get(list.Size() - 1); v != 4 {
		t.Errorf("after removing last, Get(last) = %v; want 4", v)
	}

	// Remove middle
	err = list.Remove(1)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	expected := []int{2, 4}
	for i, want := range expected {
		if v, _ := list.Get(i); v != want {
			t.Errorf("Get(%d) = %v; want %v", i, v, want)
		}
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
	if !iter.Next() {
		t.Fatal("Next() should return true for first element")
	}
	v, err := iter.Value()
	if err != nil || v != 1 {
		t.Errorf("first value = %v, %v; want 1, nil", v, err)
	}

	// Second element
	if !iter.Next() {
		t.Fatal("Next() should return true for second element")
	}
	v, err = iter.Value()
	if err != nil || v != 2 {
		t.Errorf("second value = %v, %v; want 2, nil", v, err)
	}

	// Third element
	if !iter.Next() {
		t.Fatal("Next() should return true for third element")
	}
	v, err = iter.Value()
	if err != nil || v != 3 {
		t.Errorf("third value = %v, %v; want 3, nil", v, err)
	}

	// After last element
	if iter.Next() {
		t.Error("Next() should return false after last element")
	}

	if v, err = iter.Value(); v != 3 || err != nil {
		t.Error("Value() should error after iteration end")
	}
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
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	}
	if count != 2 {
		t.Errorf("first iteration count = %d; want 2", count)
	}

	// Try to continue iteration
	if iter.Next() {
		t.Error("Next() should return false after iteration completed")
	}

	if v, err := iter.Value(); v != 2 || err != nil {
		t.Error("Value() should error after iteration completed")
	}
}

func TestLinkedListIterator_ValueBeforeNext(t *testing.T) {
	list := NewLinkedList[int]()
	list.Add(1)

	iter := list.NewIterator()
	_, err := iter.Value()
	if err == nil {
		t.Error("Value() should error when called before Next()")
	}
}

func TestLinkedListIterator_ModificationDuringIteration(t *testing.T) {
	list := NewLinkedList[int]()
	list.Add(1)
	list.Add(2)

	iter := list.NewIterator()
	iter.Next()
	list.Add(3) // Add during iteration

	value, err := iter.Value()
	if err != nil || value != 1 {
		t.Errorf("Value() = %v, %v; want 1, nil", value, err)
	}

	if !iter.Next() {
		t.Error("Next() should return true for second element")
	}
	value, err = iter.Value()
	if err != nil || value != 2 {
		t.Errorf("Value() = %v, %v; want 2, nil", value, err)
	}
}
