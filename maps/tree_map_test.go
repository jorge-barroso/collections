package maps

import (
	"testing"
)

// Test TreeMap Put and Get behavior
func TestTreeMap_PutAndGet(t *testing.T) {
	// Simple comparison function for integers
	less := func(a, b int) bool { return a < b }

	tm := NewTreeMap[int, string](less)

	// Test adding a new key-value pair
	tm.Put(1, "A")
	value, err := tm.Get(1)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if value != "A" {
		t.Errorf("expected value 'A', got %v", value)
	}

	// Test updating an existing key-value pair
	tm.Put(1, "B")
	value, err = tm.Get(1)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if value != "B" {
		t.Errorf("expected value 'B', got %v", value)
	}

	// Test retrieving a non-existent key
	_, err = tm.Get(2)
	if err == nil {
		t.Errorf("expected error when retrieving non-existent key, got none")
	}
	if err != nil && err.Error() != "key not found" {
		t.Errorf("expected error 'key not found', got %v", err)
	}
}

// Test TreeMap Remove behavior
func TestTreeMap_Remove(t *testing.T) {
	less := func(a, b int) bool { return a < b }
	tm := NewTreeMap[int, string](less)

	// Add some items to TreeMap
	tm.Put(1, "A")
	tm.Put(2, "B")
	tm.Put(3, "C")

	// Test removing an existing key
	err := tm.Remove(2)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	_, err = tm.Get(2)
	if err == nil {
		t.Errorf("expected error when retrieving removed key, got none")
	}

	// Test removing a non-existent key
	err = tm.Remove(4)
	if err == nil {
		t.Errorf("expected error when removing non-existent key, got none")
	}
	if err != nil && err.Error() != "key not found" {
		t.Errorf("expected error 'key not found', got %v", err)
	}

	// Test removing all elements
	tm.Remove(1)
	tm.Remove(3)
	if tm.Size() != 0 {
		t.Errorf("expected size to be 0 after removing all elements, got %d", tm.Size())
	}
}

// Test TreeMap Size method
func TestTreeMap_Size(t *testing.T) {
	less := func(a, b int) bool { return a < b }
	tm := NewTreeMap[int, string](less)

	// Initially, the size should be zero
	if tm.Size() != 0 {
		t.Errorf("expected size 0, got %d", tm.Size())
	}

	// Add elements and check size
	tm.Put(1, "A")
	tm.Put(2, "B")
	if tm.Size() != 2 {
		t.Errorf("expected size 2, got %d", tm.Size())
	}

	// Remove an element and check size
	tm.Remove(1)
	if tm.Size() != 1 {
		t.Errorf("expected size 1 after removal, got %d", tm.Size())
	}

	// Replacing a value (same key) should not change size
	tm.Put(2, "C")
	if tm.Size() != 1 {
		t.Errorf("expected size 1 after updating an existing key, got %d", tm.Size())
	}
}

// Test TreeMap iterator for in-order traversal
func TestTreeMap_NewIterator(t *testing.T) {
	less := func(a, b int) bool { return a < b }
	tm := NewTreeMap[int, string](less)

	// Add elements (not in sorted order) to the map
	tm.Put(3, "C")
	tm.Put(1, "A")
	tm.Put(2, "B")

	// Check that iterator returns keys in sorted order
	it := tm.NewIterator()
	expectedKeys := []int{1, 2, 3}
	expectedValues := []string{"A", "B", "C"}

	i := 0
	for it.Next() {
		entry, err := it.Value()
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		if entry.Key() != expectedKeys[i] {
			t.Errorf("expected key %d, got %d", expectedKeys[i], entry.Key())
		}
		if entry.Value() != expectedValues[i] {
			t.Errorf("expected value '%s', got '%s'", expectedValues[i], entry.Value())
		}
		i++
	}

	if i != len(expectedKeys) {
		t.Errorf("iterator did not return all elements, expected %d items, got %d", len(expectedKeys), i)
	}
}

// Test TreeMapIterator with an empty map
func TestTreeMapIterator_Empty(t *testing.T) {
	less := func(a, b int) bool { return a < b }
	tm := NewTreeMap[int, string](less)

	it := NewTreeMapIterator[int, string](tm)

	// There should not be any elements to iterate over
	if it.Next() {
		t.Errorf("expected iterator to have no elements, but got some")
	}

	_, err := it.Value()
	if err == nil || err.Error() != "no more elements" {
		t.Errorf("expected 'no more elements' error, got %v", err)
	}
}

// Test TreeMapIterator with a single element
func TestTreeMapIterator_SingleElement(t *testing.T) {
	less := func(a, b int) bool { return a < b }
	tm := NewTreeMap[int, string](less)

	tm.Put(1, "A")
	it := NewTreeMapIterator[int, string](tm)

	if !it.Next() {
		t.Errorf("expected iterator to have one element, but got none")
	}

	entry, err := it.Value()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if entry.Key() != 1 || entry.Value() != "A" {
		t.Errorf("expected key 1 with value 'A', got key %d with value '%s'", entry.Key(), entry.Value())
	}

	if it.Next() {
		t.Errorf("expected iterator to be finished, but got another element")
	}

	_, err = it.Value()
	if err == nil || err.Error() != "no more elements" {
		t.Errorf("expected 'no more elements' error, got %v", err)
	}
}

// Test TreeMapIterator with multiple elements
func TestTreeMapIterator_MultipleElements(t *testing.T) {
	less := func(a, b int) bool { return a < b }
	tm := NewTreeMap[int, string](less)

	tm.Put(3, "C")
	tm.Put(1, "A")
	tm.Put(2, "B")

	it := NewTreeMapIterator[int, string](tm)

	expected := []struct {
		key   int
		value string
	}{
		{1, "A"},
		{2, "B"},
		{3, "C"},
	}

	i := 0
	for it.Next() {
		entry, err := it.Value()
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if entry.Key() != expected[i].key || entry.Value() != expected[i].value {
			t.Errorf("expected key %d with value '%s', got key %d with value '%s'",
				expected[i].key, expected[i].value, entry.Key(), entry.Value())
		}
		i++
	}

	if i != len(expected) {
		t.Errorf("iterator did not return all elements, expected %d items, got %d", len(expected), i)
	}
}
