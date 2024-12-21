package maps

import (
	"fmt"
	"testing"
)

// Test Put and Get operations
func TestLinkedHashMap_PutAndGet(t *testing.T) {
	m := NewLinkedHashMap[string, int]()

	// Test inserting a new key-value pair
	m.Put("a", 1)
	value, err := m.Get("a")
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if value != 1 {
		t.Errorf("expected value 1, got %v", value)
	}

	// Test updating an existing key-value pair
	m.Put("a", 2)
	value, err = m.Get("a")
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if value != 2 {
		t.Errorf("expected value 2, got %v", value)
	}

	// Test retrieving a nonexistent key
	_, err = m.Get("b")
	if err == nil {
		t.Errorf("expected error, but got none")
	}
	if err != nil && err.Error() != "key not found" {
		t.Errorf("expected error 'key not found', got %v", err)
	}
}

// Test Remove functionality
func TestLinkedHashMap_Remove(t *testing.T) {
	m := NewLinkedHashMap[string, int]()

	// Add items to the map
	m.Put("a", 1)
	m.Put("b", 2)
	m.Put("c", 3)

	// Test removing an existing key
	err := m.Remove("b")
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	_, err = m.Get("b")
	if err == nil {
		t.Errorf("expected error when retrieving removed key, but got none")
	}

	// Test the order is preserved after removal
	expectedKeys := []string{"a", "c"}
	it := m.NewIterator()
	keys := []string{}
	for it.Next() {
		entry, err := it.Value()
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		keys = append(keys, entry.Key())
	}

	if fmt.Sprintf("%v", keys) != fmt.Sprintf("%v", expectedKeys) {
		t.Errorf("expected keys %v, got %v", expectedKeys, keys)
	}

	// Test removing nonexistent key
	err = m.Remove("d")
	if err == nil {
		t.Errorf("expected error for nonexistent key, got none")
	}

	// Remove remaining elements
	m.Remove("a")
	m.Remove("c")
	if m.Size() != 0 {
		t.Errorf("expected size 0 after removing all elements, got %d", m.Size())
	}
}

// Test the size functionality
func TestLinkedHashMap_Size(t *testing.T) {
	m := NewLinkedHashMap[string, int]()

	m.Put("a", 1)
	m.Put("b", 2)
	if m.Size() != 2 {
		t.Errorf("expected size 2, got %d", m.Size())
	}

	m.Remove("b")
	if m.Size() != 1 {
		t.Errorf("expected size 1 after removal, got %d", m.Size())
	}

	m.Put("a", 5) // Updating a key should not increase size
	if m.Size() != 1 {
		t.Errorf("expected size 1 after updating the same key, got %d", m.Size())
	}
}

// Test the iterator functionality
func TestLinkedHashMap_NewIterator(t *testing.T) {
	m := NewLinkedHashMap[string, int]()

	// Test empty map
	it := m.NewIterator()
	if it.Next() {
		t.Errorf("expected no elements in iterator for empty map")
	}

	// Add elements
	m.Put("a", 1)
	m.Put("b", 2)
	m.Put("c", 3)

	// Iterate through map
	expectedKeys := []string{"a", "b", "c"}
	expectedValues := []int{1, 2, 3}

	keys := []string{}
	values := []int{}
	it = m.NewIterator()
	for it.Next() {
		entry, err := it.Value()
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		keys = append(keys, entry.Key())
		values = append(values, entry.Value())
	}

	if fmt.Sprintf("%v", keys) != fmt.Sprintf("%v", expectedKeys) {
		t.Errorf("expected keys %v, got %v", expectedKeys, keys)
	}
	if fmt.Sprintf("%v", values) != fmt.Sprintf("%v", expectedValues) {
		t.Errorf("expected values %v, got %v", expectedValues, values)
	}
}

// Test iterator with an empty map
func TestLinkedHashMapIterator_Empty(t *testing.T) {
	m := NewLinkedHashMap[string, int]()
	it := m.NewIterator()

	if it.Next() {
		t.Errorf("expected false from iterator.Next() for empty map, got true")
	}

	_, err := it.Value()
	if err == nil || err.Error() != "no more elements" {
		t.Errorf("expected 'no more elements' error, got %v", err)
	}
}

// Test iterator with a single element
func TestLinkedHashMapIterator_SingleElement(t *testing.T) {
	m := NewLinkedHashMap[string, int]()
	m.Put("a", 1)

	it := m.NewIterator()

	if !it.Next() {
		t.Errorf("expected true from iterator.Next(), got false")
	}

	entry, err := it.Value()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if entry.Key() != "a" || entry.Value() != 1 {
		t.Errorf("expected key 'a' with value 1, got key '%s' with value %d", entry.Key(), entry.Value())
	}

	if it.Next() {
		t.Errorf("expected false from iterator.Next() after last element, got true")
	}

	_, err = it.Value()
	if err == nil || err.Error() != "no more elements" {
		t.Errorf("expected 'no more elements' error, got %v", err)
	}
}

// Test iterator with multiple elements
func TestLinkedHashMapIterator_MultipleElements(t *testing.T) {
	m := NewLinkedHashMap[string, int]()
	m.Put("a", 1)
	m.Put("b", 2)
	m.Put("c", 3)

	it := m.NewIterator()

	expected := []struct {
		key   string
		value int
	}{
		{"a", 1},
		{"b", 2},
		{"c", 3},
	}

	for _, e := range expected {
		if !it.Next() {
			t.Errorf("expected true from iterator.Next(), got false")
		}

		entry, err := it.Value()
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if entry.Key() != e.key || entry.Value() != e.value {
			t.Errorf("expected key '%s' with value %d, got key '%s' with value %d",
				e.key, e.value, entry.Key(), entry.Value())
		}
	}

	if it.Next() {
		t.Errorf("expected false from iterator.Next() after last element, got true")
	}
}
