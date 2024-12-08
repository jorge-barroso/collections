package lists

import (
	"errors"
)

type CopyOnWriteListIterator[T any] struct {
	index int
	list  *CopyOnWriteList[T]
}

// HasNext checks if there are more elements to iterate over
func (iter *CopyOnWriteListIterator[T]) HasNext() bool {
	iter.list.mutex.RLock()
	defer iter.list.mutex.RUnlock()
	return iter.index < len(iter.list.elements)
}

// Next returns the next element in the iteration
func (iter *CopyOnWriteListIterator[T]) Next() (T, error) {
	iter.list.mutex.RLock()
	defer iter.list.mutex.RUnlock()
	if iter.index >= len(iter.list.elements) {
		var zeroValue T
		return zeroValue, errors.New("no more elements")
	}
	element := iter.list.elements[iter.index]
	iter.index++
	return element, nil
}
