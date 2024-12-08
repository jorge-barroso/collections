package lists

import (
	"errors"
)

// ArrayList implementation using generics
type ArrayList[T any] struct {
	elements []T
}

// NewArrayList creates and returns a new instance of ArrayList
func NewArrayList[T any]() *ArrayList[T] {
	return &ArrayList[T]{elements: []T{}}
}

// NewIterator for ArrayList
func (a *ArrayList[T]) NewIterator() *ArrayListIterator[T] {
	return &ArrayListIterator[T]{
		index: -1,
		list:  a,
	}
}

// Add appends an item to the end of the list
func (a *ArrayList[T]) Add(item T) {
	a.elements = append(a.elements, item)
}

// RemoveAt removes element at specified index
func (a *ArrayList[T]) RemoveAt(index int) error {
	if index < 0 || index >= len(a.elements) {
		return errors.New("index out of bounds")
	}
	a.elements = append(a.elements[:index], a.elements[index+1:]...)
	return nil
}

// Get retrieves an element by its index
func (a *ArrayList[T]) Get(index int) (T, error) {
	var zeroValue T
	if index < 0 || index >= len(a.elements) {
		return zeroValue, errors.New("index out of bounds")
	}
	return a.elements[index], nil
}

// Size returns the number of elements in the list
func (a *ArrayList[T]) Size() int {
	return len(a.elements)
}
