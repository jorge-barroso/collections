package lists

// List interface defining common list operations
type List[T any] interface {
	Add(T)
	Remove(T) error
	Get(int) (T, error)
	Size() int
}
