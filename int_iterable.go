package collections

// Iterable interface defining methods for iterating over a collection
type Iterable[T any] interface {
	NewIterator() Iterator[T]
}
