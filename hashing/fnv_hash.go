package hashing

import (
	"fmt"
	"hash/fnv"
)

// FNVHash implements the FNV-1a hash algorithm
type FNVHash[K comparable] struct{}

// NewFNVHash creates a new FNV hash function
func NewFNVHash[K comparable]() *FNVHash[K] {
	return &FNVHash[K]{}
}

// Hash computes the FNV-1a hash of the key
func (f *FNVHash[K]) Hash(key K) uint64 {
	h := fnv.New64a()

	// Include type information in the hashing process
	typeInfo := fmt.Sprintf("%T", key) // Get the type as a string
	if _, err := h.Write([]byte(typeInfo)); err != nil {
		return 0
	}

	// Append the value of the key
	if _, err := h.Write([]byte(fmt.Sprintf("%v", key))); err != nil {
		return 0
	}

	return h.Sum64()
}
