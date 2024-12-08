package lists

import (
	"sync"
)

// CopyOnWriteList represents a generic, thread-safe list with copy-on-write semantics.
type CopyOnWriteList[T any] struct {
	mutex    sync.RWMutex
	elements []T
}

// NewCopyOnWriteList creates and returns a new instance of CopyOnWriteList
func NewCopyOnWriteList[T any]() *CopyOnWriteList[T] {
	return &CopyOnWriteList[T]{elements: []T{}}
}

// NewIterator creates a new iterator for CopyOnWriteList
func (c *CopyOnWriteList[T]) NewIterator() *CopyOnWriteListIterator[T] {
	return &CopyOnWriteListIterator[T]{
		index: 0,
		list:  c,
	}
}

// Add adds a new element to the list, creating a new copy for modification
func (c *CopyOnWriteList[T]) Add(element T) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	// Create a new slice with an additional capacity
	newElements := make([]T, len(c.elements)+1)
	copy(newElements, c.elements)
	newElements[len(c.elements)] = element
	c.elements = newElements
}

// Remove removes the element at the specified index, creating a new copy for modification
func (c *CopyOnWriteList[T]) Remove(index int) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	if index < 0 || index >= len(c.elements) {
		return
	}
	// Create a new slice with the current elements except the one being removed
	newElements := append([]T{}, c.elements[:index]...)
	newElements = append(newElements, c.elements[index+1:]...)
	c.elements = newElements
}

// Get returns the element at the specified index
func (c *CopyOnWriteList[T]) Get(index int) (T, bool) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	if index < 0 || index >= len(c.elements) {
		var zeroValue T
		return zeroValue, false
	}
	return c.elements[index], true
}

// Size returns the number of elements in the list
func (c *CopyOnWriteList[T]) Size() int {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	return len(c.elements)
}
