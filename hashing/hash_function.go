package hashing

// HashFunction defines the interface for hash functions
type HashFunction[K comparable] interface {
	// Hash computes a hash value for the given key
	Hash(key K) uint64
}
