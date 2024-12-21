package maps

import (
	"errors"
	"github.com/jorge-barroso/collections"
	"github.com/jorge-barroso/collections/hashing"
	"sync"
)

// ShardCount determines the number of segments in the concurrent map
const ShardCount = 16

// mapShard represents a single shard of the concurrent map
type mapShard[K comparable, V any] struct {
	items map[K]V
	sync.RWMutex
}

// ConcurrentHashMap implements a thread-safe map using multiple shards
type ConcurrentHashMap[K comparable, V any] struct {
	shards    [ShardCount]mapShard[K, V]
	size      int64
	sizeMutex sync.RWMutex
	hashFunc  hashing.HashFunction[K]
}

// Ensure ConcurrentHashMap implements both Map and Iterable interfaces
var _ Map[string, int] = (*ConcurrentHashMap[string, int])(nil)
var _ collections.Iterable[Entry[string, int]] = (*ConcurrentHashMap[string, int])(nil)

// NewConcurrentHashMap creates a new ConcurrentHashMap with default hash function
func NewConcurrentHashMap[K comparable, V any]() *ConcurrentHashMap[K, V] {
	return NewConcurrentHashMapWithHash[K, V](hashing.NewFNVHash[K]())
}

// NewConcurrentHashMapWithHash creates a new ConcurrentHashMap with a custom hash function
func NewConcurrentHashMapWithHash[K comparable, V any](hashFunc hashing.HashFunction[K]) *ConcurrentHashMap[K, V] {
	cm := &ConcurrentHashMap[K, V]{
		hashFunc: hashFunc,
	}
	for i := 0; i < ShardCount; i++ {
		cm.shards[i] = mapShard[K, V]{
			items: make(map[K]V),
		}
	}
	return cm
}

// getShard returns the appropriate shard for a given key
func (cm *ConcurrentHashMap[K, V]) getShard(key K) *mapShard[K, V] {
	hashCode := cm.hashFunc.Hash(key)
	return &cm.shards[hashCode%ShardCount]
}

// Put adds or updates a key-value pair
func (cm *ConcurrentHashMap[K, V]) Put(key K, value V) {
	shard := cm.getShard(key)
	shard.Lock()
	defer shard.Unlock()

	_, exists := shard.items[key]
	if !exists {
		cm.sizeMutex.Lock()
		cm.size++
		cm.sizeMutex.Unlock()
	}

	shard.items[key] = value
}

// Get retrieves a value by key
func (cm *ConcurrentHashMap[K, V]) Get(key K) (V, error) {
	shard := cm.getShard(key)
	shard.RLock()
	defer shard.RUnlock()

	value, ok := shard.items[key]
	if !ok {
		var zero V
		return zero, errors.New("key not found")
	}
	return value, nil
}

// Remove deletes a key-value pair
func (cm *ConcurrentHashMap[K, V]) Remove(key K) error {
	shard := cm.getShard(key)
	shard.Lock()
	defer shard.Unlock()

	if _, exists := shard.items[key]; !exists {
		return errors.New("key not found")
	}

	delete(shard.items, key)
	cm.sizeMutex.Lock()
	cm.size--
	cm.sizeMutex.Unlock()
	return nil
}

// Size returns the total number of elements across all shards
func (cm *ConcurrentHashMap[K, V]) Size() int {
	cm.sizeMutex.RLock()
	defer cm.sizeMutex.RUnlock()
	return int(cm.size)
}

// Clear removes all elements from the map
func (cm *ConcurrentHashMap[K, V]) Clear() {
	for i := 0; i < ShardCount; i++ {
		shard := &cm.shards[i]
		shard.Lock()
		shard.items = make(map[K]V)
		shard.Unlock()
	}
	cm.sizeMutex.Lock()
	cm.size = 0
	cm.sizeMutex.Unlock()
}

// ContainsKey checks if a key exists in the map
func (cm *ConcurrentHashMap[K, V]) ContainsKey(key K) bool {
	shard := cm.getShard(key)
	shard.RLock()
	defer shard.RUnlock()
	_, exists := shard.items[key]
	return exists
}

// NewIterator returns a new iterator for the concurrent map
func (cm *ConcurrentHashMap[K, V]) NewIterator() collections.Iterator[Entry[K, V]] {
	iterator := &ConcurrentHashMapIterator[K, V]{
		cm:           cm,
		currentShard: 0,
		position:     -1,
	}
	iterator.loadNextShard()
	return iterator
}
