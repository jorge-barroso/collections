package collections

type Node[T any] struct {
	Item T
	Next *Node[T]
}
