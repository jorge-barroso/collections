package maps

import (
	"fmt"
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
	expectedSize := int64(numGoroutines * numOperations)
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

func TestConcurrentHashMap_UpdateExistingKey(t *testing.T) {
	cm := NewConcurrentHashMap[string, int]()

	cm.Put("key", 10)
	value, _ := cm.Get("key")
	if value != 10 {
		t.Errorf("expected 10, got %d", value)
	}

	cm.Put("key", 20)
	value, _ = cm.Get("key")
	if value != 20 {
		t.Errorf("expected updated value 20, got %d", value)
	}

	if size := cm.Size(); size != 1 {
		t.Errorf("size should remain 1 after overwriting key, got %d", size)
	}
}

func TestConcurrentHashMap_RemoveNonExistentKey(t *testing.T) {
	cm := NewConcurrentHashMap[string, int]()
	err := cm.Remove("nonexistent")
	if err == nil {
		t.Error("expected error when removing non-existent key, got nil")
	}
}

func TestConcurrentHashMap_ConcurrentReadsAndWrites(t *testing.T) {
	cm := NewConcurrentHashMap[string, int]()
	var wg sync.WaitGroup

	// Concurrent write operations
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(base int) {
			defer wg.Done()
			for j := 0; j < 100; j++ {
				cm.Put(fmt.Sprintf("key%d", base+j), base+j)
			}
		}(i * 100)
	}

	// Concurrent read operations
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(base int) {
			defer wg.Done()
			for j := 0; j < 100; j++ {
				_, _ = cm.Get(fmt.Sprintf("key%d", base+j))
			}
		}(i * 100)
	}

	wg.Wait()

	// Verify size
	expectedSize := int64(1000)
	if size := cm.Size(); size != expectedSize {
		t.Errorf("expected size %d, got %d", expectedSize, size)
	}
}

func TestConcurrentHashMap_EmptyKeyOrValue(t *testing.T) {
	cm := NewConcurrentHashMap[string, string]()

	// Test empty key
	cm.Put("", "emptyKey")
	value, _ := cm.Get("")
	if value != "emptyKey" {
		t.Errorf("expected value 'emptyKey', got %v", value)
	}

	// Test empty value
	cm.Put("emptyValueKey", "")
	value, _ = cm.Get("emptyValueKey")
	if value != "" {
		t.Errorf("expected empty value, got %v", value)
	}
}

func TestConcurrentHashMap_StressTest(t *testing.T) {
	cm := NewConcurrentHashMap[int, int]()
	var wg sync.WaitGroup
	numEntries := int64(1_000_000)

	// Concurrent writes
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; int64(i) < numEntries; i++ {
			cm.Put(i, i)
		}
	}()

	// Concurrent reads
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; int64(i) < numEntries; i++ {
			_, _ = cm.Get(i)
		}
	}()

	wg.Wait()

	// Verify size
	if size := cm.Size(); size != numEntries {
		t.Errorf("expected size %d, got %d", numEntries, size)
	}
}
func TestConcurrentHashMap_IterateWhileModifying(t *testing.T) {
	cm := NewConcurrentHashMap[int, int]()
	for i := 0; i < 100; i++ {
		cm.Put(i, i)
	}

	var wg sync.WaitGroup

	// Concurrent iteration
	wg.Add(1)
	go func() {
		defer wg.Done()
		iter := cm.NewIterator()
		for iter.Next() {
			_, _ = iter.Value()
		}
	}()

	// Concurrent modifications
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 100; i < 200; i++ {
			cm.Put(i, i)
		}
	}()

	wg.Wait()
}

func TestConcurrentHashMap_Clear(t *testing.T) {
	cm := NewConcurrentHashMap[string, string]()

	cm.Put("key1", "value1")
	cm.Put("key2", "value2")

	cm.Clear()

	if cm.Size() != 0 {
		t.Errorf("expected size 0 after Clear, got %d", cm.Size())
	}

	if cm.ContainsKey("key1") || cm.ContainsKey("key2") {
		t.Error("map should not contain any keys after Clear")
	}
}

func TestConcurrentHashMap_HighlyConcurrentReadsAndWrites(t *testing.T) {
	cm := NewConcurrentHashMap[int, int]()
	var wg sync.WaitGroup
	numGoroutines := 50              // Number of concurrent goroutines
	numOperationsPerGoroutine := 500 // Operations performed by each goroutine

	// Concurrent writes
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(base int) {
			defer wg.Done()
			for j := 0; j < numOperationsPerGoroutine; j++ {
				cm.Put(base+j, base+j)
			}
		}(i * numOperationsPerGoroutine)
	}

	// Concurrent reads while writes are ongoing
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(base int) {
			defer wg.Done()
			for j := 0; j < numOperationsPerGoroutine; j++ {
				_, _ = cm.Get(base + j)
			}
		}(i * numOperationsPerGoroutine)
	}

	wg.Wait()

	// Verify the total size
	expectedSize := int64(numGoroutines * numOperationsPerGoroutine)
	if size := cm.Size(); size != expectedSize {
		t.Errorf("expected size %d, got %d", expectedSize, size)
	}
}

func TestConcurrentHashMap_ConcurrentIteration(t *testing.T) {
	cm := NewConcurrentHashMap[int, int]()
	numElements := 1000 // Number of elements added to the map
	numReaders := 10    // Number of concurrent iterators

	// Add elements to the map
	for i := 0; i < numElements; i++ {
		cm.Put(i, i)
	}

	// Concurrent iteration
	var wg sync.WaitGroup
	for i := 0; i < numReaders; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			iter := cm.NewIterator()
			count := 0
			for iter.Next() {
				_, err := iter.Value()
				if err != nil {
					t.Errorf("unexpected iterator error: %v", err)
					return
				}
				count++
			}

			// Ensure count is not greater than total size
			if count > numElements {
				t.Errorf("iterator read too many elements: %d > %d", count, numElements)
			}
		}()
	}
	wg.Wait()
}

func TestConcurrentHashMap_ClearWhileAccessing(t *testing.T) {
	cm := NewConcurrentHashMap[int, int]()
	numElements := 1000 // Number of elements added to the map

	// Add elements to the map
	for i := 0; i < numElements; i++ {
		cm.Put(i, i)
	}

	// Concurrent clear while accessing
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		cm.Clear()
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < numElements; i++ {
			_, _ = cm.Get(i)
		}
	}()
	wg.Wait()

	// Validate map is cleared
	if cm.Size() != 0 {
		t.Errorf("expected size 0 after clear, got %d", cm.Size())
	}
}

func TestConcurrentHashMap_RemoveWhileAccessing(t *testing.T) {
	cm := NewConcurrentHashMap[int, int]()
	numElements := 100 // Number of elements added to the map

	// Add elements to the map
	for i := 0; i < numElements; i++ {
		cm.Put(i, i)
	}

	// Concurrent removes while reading
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < numElements; i++ {
			_ = cm.Remove(i)
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < numElements; i++ {
			_, _ = cm.Get(i)
		}
	}()
	wg.Wait()

	// Validate map is empty
	if cm.Size() != 0 {
		t.Errorf("expected size 0 after removals, got %d", cm.Size())
	}
}
