package lists

import "errors"

// ArrayListIterator struct for ArrayList
type ArrayListIterator[T any] struct {
	index int
	list  *ArrayList[T]
}

func (iter *ArrayListIterator[T]) Next() bool {
	return iter.index < len(iter.list.elements)-1
}

func (iter *ArrayListIterator[T]) Value() (T, error) {
	if !iter.Next() {
		var zero T
		return zero, errors.New("no more elements")
	}
	iter.index++
	return iter.list.elements[iter.index], nil
}
