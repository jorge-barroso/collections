package lists

import "errors"

// ArrayListIterator struct for ArrayList
type ArrayListIterator[T any] struct {
	index int
	list  *ArrayList[T]
}

func (aIt *ArrayListIterator[T]) Next() bool {
	if aIt.index < len(aIt.list.elements)-1 {
		aIt.index++
		return true
	}

	return false
}

func (aIt *ArrayListIterator[T]) Value() (T, error) {
	if aIt.index < 0 {
		var zero T
		return zero, errors.New("out of bounds, call Next() first")
	}

	return aIt.list.elements[aIt.index], nil
}
