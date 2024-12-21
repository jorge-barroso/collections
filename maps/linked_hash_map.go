package maps

import (
	"errors"
	"github.com/jorge-barroso/collections"
)

// LinkedHashMap implements both Map and collections.Iterable interfaces
type LinkedHashMap[K comparable, V any] struct {
	items map[K]*collections.Node[Entry[K, V]]
	head  *collections.Node[Entry[K, V]]
	tail  *collections.Node[Entry[K, V]]
	size  int64
}

// Ensure LinkedHashMap implements both Map and Iterable interfaces
var _ Map[string, int] = (*LinkedHashMap[string, int])(nil)
var _ collections.Iterable[Entry[string, int]] = (*LinkedHashMap[string, int])(nil)

// NewLinkedHashMap creates a new LinkedHashMap
func NewLinkedHashMap[K comparable, V any]() *LinkedHashMap[K, V] {
	return &LinkedHashMap[K, V]{
		items: make(map[K]*collections.Node[Entry[K, V]]),
	}
}

// Put inserts or updates a key-value pair
func (m *LinkedHashMap[K, V]) Put(key K, value V) {
	entry := Entry[K, V]{
		key:   key,
		value: value,
	}

	if existingNode, ok := m.items[key]; ok {
		// Update existing node value
		existingNode.Item = entry
		return
	}

	// Create new node
	newNode := &collections.Node[Entry[K, V]]{
		Item: entry,
	}

	// Add to hashmap
	m.items[key] = newNode

	// Handle linked list insertion
	if m.tail == nil {
		// First element
		m.head = newNode
		m.tail = newNode
	} else {
		// Add to end
		m.tail.Next = newNode
		m.tail = newNode
	}

	m.size++
}

// Get retrieves the value associated with a key
func (m *LinkedHashMap[K, V]) Get(key K) (V, error) {
	if node, ok := m.items[key]; ok {
		return node.Item.Value(), nil
	}
	var zero V
	return zero, errors.New("key not found")
}

// Remove removes a key-value pair
func (m *LinkedHashMap[K, V]) Remove(key K) error {
	target, ok := m.items[key]
	if !ok {
		return errors.New("key not found")
	}

	// Remove from hashmap
	delete(m.items, key)

	// Special case: single element
	if m.head == m.tail {
		m.head = nil
		m.tail = nil
		m.size--
		return nil
	}

	// Special case: remove head
	if m.head == target {
		m.head = m.head.Next
		m.size--
		return nil
	}

	// Find the previous node
	current := m.head
	for current != nil && current.Next != target {
		current = current.Next
	}

	if current != nil {
		// Remove the target node
		current.Next = target.Next
		// Update tail if needed
		if target == m.tail {
			m.tail = current
		}
	}

	m.size--
	return nil
}

// Size returns the number of key-value pairs
func (m *LinkedHashMap[K, V]) Size() int64 {
	return m.size
}

// NewIterator returns a new iterator for the LinkedHashMap
func (m *LinkedHashMap[K, V]) NewIterator() collections.Iterator[Entry[K, V]] {
	return &LinkedHashMapIterator[K, V]{
		current: m.head,
	}
}
