package hashing

import (
	"testing"
)

// Hypothetical FNVHash function or method
// Assuming a placeholder method called 'Hash' is added to FNVHash for hashing input values.
func TestFNVHash_IntValues(t *testing.T) {
	hasher := FNVHash[int]{}

	// Test hashing integer values
	input := []int{1, 2, 3, 4, 5}
	hashes := make(map[int]uint64)

	for _, val := range input {
		hash := hasher.Hash(val)
		if _, exists := hashes[val]; exists {
			t.Errorf("Hash for input %d already exists", val)
		}
		hashes[val] = hash
	}

	// Ensure no two different inputs have the same hash
	for i, h1 := range hashes {
		for j, h2 := range hashes {
			if i != j && h1 == h2 {
				t.Errorf("Hashes collision detected between %d and %d (hash: %d)", i, j, h1)
			}
		}
	}
}

func TestFNVHash_StringValues(t *testing.T) {
	hasher := FNVHash[string]{}

	// Test hashing string values
	strings := []string{"hello", "world", "hash", "test", "openai"}
	hashes := make(map[string]uint64)

	for _, str := range strings {
		hash := hasher.Hash(str)
		if _, exists := hashes[str]; exists {
			t.Errorf("Hash for input %q already exists", str)
		}
		hashes[str] = hash
	}

	// Ensure no two different inputs have the same hash
	for i, h1 := range hashes {
		for j, h2 := range hashes {
			if i != j && h1 == h2 {
				t.Errorf("Hashes collision detected between %q and %q (hash: %d)", i, j, h1)
			}
		}
	}
}

func TestFNVHash_EmptyValue(t *testing.T) {
	hasher := FNVHash[string]{}

	// Hash an empty string
	emptyHash := hasher.Hash("")
	anotherEmptyHash := hasher.Hash("")

	// Ensure hashing the same empty input gives the same result
	if emptyHash != anotherEmptyHash {
		t.Errorf("Hash mismatch for empty values: %d, %d", emptyHash, anotherEmptyHash)
	}
}

func TestFNVHash_DifferentTypes(t *testing.T) {
	// Create different instances for different key types
	intHasher := FNVHash[int]{}
	stringHasher := FNVHash[string]{}

	// Hash integer and string with the same logical value but ensure types don't clash
	intHash := intHasher.Hash(123)
	stringHash := stringHasher.Hash("123")

	if intHash == stringHash {
		t.Errorf("Hashes between int '123' and string '123' unexpectedly match: %d", intHash)
	}
}
