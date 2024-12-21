package lists

import (
	"fmt"
	"github.com/jorge-barroso/collections"
)

type LinkedList[T any] struct {
	dummy *collections.Node[T] // Dummy first node
	size  int
}

func NewLinkedList[T any]() *LinkedList[T] {
	var zero T
	return &LinkedList[T]{
		dummy: &collections.Node[T]{Item: zero},
	}
}

func (ll *LinkedList[T]) Add(value T) {
	newNode := &collections.Node[T]{Item: value}
	current := ll.dummy
	for current.Next != nil {
		current = current.Next
	}
	current.Next = newNode
	ll.size++
}

func (ll *LinkedList[T]) Remove(index int) error {
	if index < 0 || index >= ll.size {
		return fmt.Errorf("index out of bounds: %d", index)
	}

	current := ll.dummy
	for i := 0; i < index; i++ {
		current = current.Next
	}
	current.Next = current.Next.Next
	ll.size--
	return nil
}

func (ll *LinkedList[T]) Get(index int) (T, error) {
	if index < 0 || index >= ll.size {
		var zero T
		return zero, fmt.Errorf("index out of bounds: %d", index)
	}

	current := ll.dummy.Next
	for i := 0; i < index; i++ {
		current = current.Next
	}
	return current.Item, nil
}

func (ll *LinkedList[T]) Size() int {
	return ll.size
}

func (ll *LinkedList[T]) NewIterator() collections.Iterator[T] {
	return &LinkedListIterator[T]{
		current: ll.dummy,
	}
}
