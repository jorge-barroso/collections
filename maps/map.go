package maps

type Map[K comparable, V any] interface {
	Put(key K, value V)   // Inserts or updates a key-value pair
	Get(key K) (V, error) // Retrieves the value associated with a key
	Remove(key K) error   // Removes a key-value pair
	Size() int64          // Returns the number of key-value pairs
}
