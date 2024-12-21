package maps

import (
	"errors"
	"github.com/jorge-barroso/collections"
)

// LinkedHashMapIterator implements the Iterator interface
type LinkedHashMapIterator[K comparable, V any] struct {
	current *collections.Node[Entry[K, V]]
}

// Next checks if there are more elements
func (it *LinkedHashMapIterator[K, V]) Next() bool {
	return it.current != nil
}

// Value returns the current element and advances the iterator
func (it *LinkedHashMapIterator[K, V]) Value() (Entry[K, V], error) {
	if it.current == nil {
		var zero Entry[K, V]
		return zero, errors.New("no more elements")
	}
	value := it.current.Item
	it.current = it.current.Next
	return value, nil
}
