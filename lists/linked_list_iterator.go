package lists

import (
	"errors"
	"github.com/jorge-barroso/collections"
)

// LinkedListIterator struct for LinkedList
type LinkedListIterator[T any] struct {
	current *collections.Node[T]
}

func (iter *LinkedListIterator[T]) Next() bool {
	if iter.current != nil && iter.current.Next != nil {
		iter.current = iter.current.Next
		return true
	}
	return false
}

func (iter *LinkedListIterator[T]) Value() (T, error) {
	if iter.current == nil {
		var zeroValue T
		return zeroValue, errors.New("no more elements")
	}
	return iter.current.Item, nil
}
