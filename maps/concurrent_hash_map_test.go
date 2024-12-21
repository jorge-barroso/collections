package maps

import (
	"sync"
	"testing"
)

func TestConcurrentHashMap_BasicOperations(t *testing.T) {
	cm := NewConcurrentHashMap[string, int]()

	// Test Put and Get
	cm.Put("one", 1)
	cm.Put("two", 2)

	value, err := cm.Get("one")
	if err != nil || value != 1 {
		t.Errorf("expected 1, got %d", value)
	}

	value, err = cm.Get("two")
	if err != nil || value != 2 {
		t.Errorf("expected 2, got %d", value)
	}

	// Test Size
	if size := cm.Size(); size != 2 {
		t.Errorf("expected size 2, got %d", size)
	}

	// Test Remove
	err = cm.Remove("one")
	if err != nil {
		t.Errorf("unexpected error on remove: %v", err)
	}

	if size := cm.Size(); size != 1 {
		t.Errorf("expected size 1 after remove, got %d", size)
	}

	// Test ContainsKey
	if !cm.ContainsKey("two") {
		t.Error("expected to contain key 'two'")
	}

	if cm.ContainsKey("one") {
		t.Error("expected to not contain key 'one' after removal")
	}
}

func TestConcurrentHashMap_ConcurrentAccess(t *testing.T) {
	cm := NewConcurrentHashMap[int, int]()
	var wg sync.WaitGroup
	numGoroutines := 10
	numOperations := 1000

	// Test concurrent writes
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(base int) {
			defer wg.Done()
			for j := 0; j < numOperations; j++ {
				cm.Put(base+j, base+j)
			}
		}(i * numOperations)
	}
	wg.Wait()

	// Test concurrent reads
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(base int) {
			defer wg.Done()
			for j := 0; j < numOperations; j++ {
				_, _ = cm.Get(base + j)
			}
		}(i * numOperations)
	}
	wg.Wait()

	// Verify size
	expectedSize := numGoroutines * numOperations
	if size := cm.Size(); size != expectedSize {
		t.Errorf("expected size %d, got %d", expectedSize, size)
	}
}

func TestConcurrentHashMap_Iterator(t *testing.T) {
	cm := NewConcurrentHashMap[string, int]()
	testData := map[string]int{
		"one":   1,
		"two":   2,
		"three": 3,
		"four":  4,
	}

	// Add test data
	for k, v := range testData {
		cm.Put(k, v)
	}

	// Test iterator
	iter := cm.NewIterator()
	count := 0
	foundValues := make(map[string]bool)

	for iter.Next() {
		entry, err := iter.Value()
		if err != nil {
			t.Errorf("unexpected error: %v", err)
			continue
		}

		count++
		key := entry.Key()
		value := entry.Value()

		if expectedValue, exists := testData[key]; !exists {
			t.Errorf("unexpected key: %v", key)
		} else if expectedValue != value {
			t.Errorf("for key %v, expected value %v, got %v", key, expectedValue, value)
		}

		foundValues[key] = true
	}

	if count != len(testData) {
		t.Errorf("iterator visited %d elements, expected %d", count, len(testData))
	}

	for k := range testData {
		if !foundValues[k] {
			t.Errorf("iterator missed key: %v", k)
		}
	}
}
