package lists

import (
	"fmt"
	"github.com/jorge-barroso/collections"
)

// LinkedList represents a singly linked list
type LinkedList[T any] struct {
	head *collections.Node[T]
	size int
}

// NewLinkedList creates and returns a new instance of LinkedList
func NewLinkedList[T any]() *LinkedList[T] {
	return &LinkedList[T]{}
}

// NewIterator for LinkedList
func (ll *LinkedList[T]) NewIterator() collections.Iterator[T] {
	return &LinkedListIterator[T]{
		current: ll.head,
	}
}

// Add appends an item to the end of the list
func (ll *LinkedList[T]) Add(value T) {
	newNode := &collections.Node[T]{Item: value}
	if ll.head == nil {
		ll.head = newNode
	} else {
		current := ll.head
		for current.Next != nil {
			current = current.Next
		}
		current.Next = newNode
	}
	ll.size++
}

// Remove removes an element at the specified index
func (ll *LinkedList[T]) Remove(index int) error {
	if index < 0 || index >= ll.size {
		return fmt.Errorf("index out of bounds, must be between 0 and %d, but was %d", ll.size-1, index)
	}

	if index == 0 {
		ll.head = ll.head.Next
	} else {
		prev := ll.head
		for i := 0; i < index-1; i++ {
			prev = prev.Next
		}
		prev.Next = prev.Next.Next
	}
	ll.size--
	return nil
}

// Get retrieves an element by its index
func (ll *LinkedList[T]) Get(index int) (T, error) {
	if index < 0 || index >= ll.size {
		var zeroValue T
		return zeroValue, fmt.Errorf("index out of bounds, must be between 0 and %d, but was %d", ll.size-1, index)
	}

	current := ll.head
	for i := 0; i < index; i++ {
		current = current.Next
	}
	return current.Item, nil
}

// Size returns the number of elements in the list
func (ll *LinkedList[T]) Size() int {
	return ll.size
}
