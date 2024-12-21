package lists

import (
	"github.com/jorge-barroso/collections"
	"sync"
)

// CopyOnWriteList is a thread-safe list implementation that creates a new
// copy of the underlying array whenever the list is modified
type CopyOnWriteList[T any] struct {
	listOps[T]
	lock     sync.Mutex
	elements []T
}

// Ensure CopyOnWriteList implements both Map and Iterable interfaces
var _ List[int] = (*CopyOnWriteList[int])(nil)
var _ collections.Iterable[int] = (*CopyOnWriteList[int])(nil)

// NewCopyOnWriteListWithCapacity creates a new CopyOnWriteList instance with the desired initial capacity
func NewCopyOnWriteListWithCapacity[T any](capacity int) *CopyOnWriteList[T] {
	return &CopyOnWriteList[T]{
		elements: make([]T, 0, capacity),
	}
}

// NewCopyOnWriteList creates a new CopyOnWriteList instance
func NewCopyOnWriteList[T any]() *CopyOnWriteList[T] {
	return NewCopyOnWriteListWithCapacity[T](defaultCapacity)
}

// Add appends an item to the end of the list
func (l *CopyOnWriteList[T]) Add(item T) {
	l.lock.Lock()
	defer l.lock.Unlock()

	// Create a new slice with one more capacity
	newElements := make([]T, len(l.elements)+1)
	// Copy existing elements
	copy(newElements, l.elements)
	// Add new element
	newElements[len(l.elements)] = item
	// Replace old slice with new one
	l.elements = newElements
}

// Remove removes an item at the specified index
func (l *CopyOnWriteList[T]) Remove(index int) error {
	l.lock.Lock()
	defer l.lock.Unlock()

	if err := l.validateIndex(index, l.Size()); err != nil {
		return err
	}

	// Create new slice with one less capacity
	newElements := make([]T, len(l.elements)-1)
	// Copy elements before index
	copy(newElements, l.elements[:index])
	// Copy elements after index
	copy(newElements[index:], l.elements[index+1:])
	// Replace old slice with new one
	l.elements = newElements
	return nil
}

// Get retrieves an element by its index
func (l *CopyOnWriteList[T]) Get(index int) (T, error) {
	// No lock needed for reads
	if err := l.validateIndex(index, l.Size()); err != nil {
		var zeroValue T
		return zeroValue, err
	}

	return l.elements[index], nil
}

// Size returns the number of elements in the list
func (l *CopyOnWriteList[T]) Size() int {
	return len(l.elements)
}

// NewIterator returns a new iterator for the list
func (l *CopyOnWriteList[T]) NewIterator() collections.Iterator[T] {
	// Create a snapshot of current elements
	l.lock.Lock()
	snapshot := make([]T, len(l.elements))
	copy(snapshot, l.elements)
	l.lock.Unlock()

	return &CopyOnWriteListIterator[T]{
		snapshot: snapshot,
		index:    -1,
	}
}
