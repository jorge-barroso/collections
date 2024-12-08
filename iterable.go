package collections

// Iterable interface defining methods for iterating over a collection
type Iterable[T any] interface {
	HasNext() bool
	Next() (T, error)
}
