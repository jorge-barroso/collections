package maps

import (
	"errors"
)

// TreeMapIterator implements in-order traversal
type TreeMapIterator[K comparable, V any] struct {
	tree    *TreeMap[K, V]
	current *rbNode[K, V]
}

// NewTreeMapIterator creates a new iterator starting at the leftmost node
func NewTreeMapIterator[K comparable, V any](tree *TreeMap[K, V]) *TreeMapIterator[K, V] {
	var firstNode *rbNode[K, V]
	if tree.root != nil {
		firstNode = tree.root.getMinimum()
	}

	return &TreeMapIterator[K, V]{
		tree:    tree,
		current: firstNode,
	}
}

// Next checks if there are more elements
func (it *TreeMapIterator[K, V]) Next() bool {
	return it.current != nil
}

// Value returns the current element and advances the iterator
func (it *TreeMapIterator[K, V]) Value() (Entry[K, V], error) {
	if it.current == nil {
		var zero Entry[K, V]
		return zero, errors.New("no more elements")
	}

	value := it.current.Node.Item
	it.current = it.current.getSuccessor()
	return value, nil
}
