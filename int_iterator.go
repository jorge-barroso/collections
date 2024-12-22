package collections

type Iterator[T any] interface {
	Next() bool
	Value() (T, error)
}
