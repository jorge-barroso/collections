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
func (cIt *CopyOnWriteListIterator[T]) Next() bool {
	if cIt.index+1 < len(cIt.snapshot) {
		cIt.index++
		return true
	}
	return false
}

// Value returns the current element in the iteration
func (cIt *CopyOnWriteListIterator[T]) Value() (T, error) {
	if cIt.index < 0 {
		var zero T
		return zero, errors.New("no more elements, call Next() first")
	}
	return cIt.snapshot[cIt.index], nil
}
