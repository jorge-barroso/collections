package lists

import (
	"errors"
)

// CopyOnWriteListIterator provides iteration over a snapshot of CopyOnWriteList elements
type CopyOnWriteListIterator[T any] struct {
	snapshot []T // Immutable snapshot of the list at iterator creation time
	index    int // Current position in the snapshot
}

// Next advances the iterator to the next element
func (it *CopyOnWriteListIterator[T]) Next() bool {
	if it.index+1 < len(it.snapshot) {
		it.index++
		return true
	}
	return false
}

// Value returns the current element in the iteration
func (it *CopyOnWriteListIterator[T]) Value() (T, error) {
	if it.index < 0 {
		var zero T
		return zero, errors.New("no more elements")
	}
	return it.snapshot[it.index], nil
}
