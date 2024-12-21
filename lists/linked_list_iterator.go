package lists

import (
	"errors"
	"github.com/jorge-barroso/collections"
)

// LinkedListIterator struct for LinkedList
type LinkedListIterator[T any] struct {
	current *collections.Node[T]
	started bool
}

func (iter *LinkedListIterator[T]) Next() bool {
	if !iter.started {
		iter.started = true
		return iter.current != nil
	}

	if iter.current != nil && iter.current.Next != nil {
		iter.current = iter.current.Next
		return true
	}
	return false
}

// Value retrieves the value of the current node.
func (iter *LinkedListIterator[T]) Value() (T, error) {
	if !iter.started {
		var zeroValue T
		return zeroValue, errors.New("iterator is not positioned on a valid element")
	}

	if iter.current == nil {
		var zeroValue T
		return zeroValue, errors.New("no more elements")
	}
	return iter.current.Item, nil
}
