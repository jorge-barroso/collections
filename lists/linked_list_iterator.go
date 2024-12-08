package lists

import (
	"collections"
	"errors"
)

// LinkedListIterator struct for LinkedList
type LinkedListIterator[T any] struct {
	current *collections.Node[T]
}

func (iter *LinkedListIterator[T]) HasNext() bool {
	return iter.current != nil
}

func (iter *LinkedListIterator[T]) Next() (T, error) {
	if iter.current == nil {
		var zeroValue T
		return zeroValue, errors.New("no more elements")
	}
	value := iter.current.Item
	iter.current = iter.current.Next
	return value, nil
}
