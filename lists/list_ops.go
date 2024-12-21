package lists

import "fmt"

type listOps[T any] struct {
}

func (lo *listOps[T]) validateIndex(index, size int) error {
	if index < 0 || index >= size {
		return fmt.Errorf("index out of bounds, must be between 0 and %d, but %d was provided", size-1, index)
	}
	return nil
}
