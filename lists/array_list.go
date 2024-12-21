package lists

import (
	"github.com/jorge-barroso/collections"
)

// ArrayList implementation using generics
type ArrayList[T any] struct {
	listOps[T]
	elements []T
}

// Ensure ArrayList implements both Map and Iterable interfaces
var _ List[int] = (*ArrayList[int])(nil)
var _ collections.Iterable[int] = (*ArrayList[int])(nil)

// NewArrayListWithCapacity creates and returns a new instance of ArrayList with the desired initial capacity
func NewArrayListWithCapacity[T any](capacity int) *ArrayList[T] {
	return &ArrayList[T]{
		elements: make([]T, 0, capacity),
	}
}

// NewArrayList creates and returns a new instance of ArrayList
func NewArrayList[T any]() *ArrayList[T] {
	return NewArrayListWithCapacity[T](defaultCapacity)
}

// Add appends an item to the end of the list
func (a *ArrayList[T]) Add(item T) {
	a.elements = append(a.elements, item)
}

// Remove removes element at specified index
func (a *ArrayList[T]) Remove(index int) error {
	if err := a.validateIndex(index, a.Size()); err != nil {
		return err
	}
	a.elements = append(a.elements[:index], a.elements[index+1:]...)
	return nil
}

// Get retrieves an element by its index
func (a *ArrayList[T]) Get(index int) (T, error) {
	if err := a.validateIndex(index, a.Size()); err != nil {
		var zeroValue T
		return zeroValue, err
	}

	return a.elements[index], nil
}

// Size returns the number of elements in the list
func (a *ArrayList[T]) Size() int {
	return len(a.elements)
}

// NewIterator creates and returns a new iterator for the ArrayList.
func (a *ArrayList[T]) NewIterator() collections.Iterator[T] {
	return &ArrayListIterator[T]{
		index: -1,
		list:  a,
	}
}
