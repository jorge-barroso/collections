package maps

type Entry[K comparable, V any] struct {
	key   K
	value V
}

// Key returns the key of the entry.
func (e *Entry[K, V]) Key() K {
	return e.key
}

// Value returns the value of the entry.
func (e *Entry[K, V]) Value() V {
	return e.value
}
