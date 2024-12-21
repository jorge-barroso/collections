package maps

import (
	"errors"
)

// ConcurrentHashMapIterator implements iterator for ConcurrentHashMap
type ConcurrentHashMapIterator[K comparable, V any] struct {
	cm           *ConcurrentHashMap[K, V]
	currentShard int
	entries      []Entry[K, V]
	position     int
}

// loadNextShard loads entries from the next available shard
func (it *ConcurrentHashMapIterator[K, V]) loadNextShard() {
	for it.currentShard < ShardCount {
		shard := &it.cm.shards[it.currentShard]
		shard.RLock()
		entries := make([]Entry[K, V], 0, len(shard.items))
		for k, v := range shard.items {
			entries = append(entries, Entry[K, V]{key: k, value: v})
		}
		shard.RUnlock()

		it.currentShard++
		if len(entries) > 0 {
			it.entries = entries
			it.position = -1
			return
		}
	}
	it.entries = nil
	it.position = -1
}

// Next checks if there are more elements
func (it *ConcurrentHashMapIterator[K, V]) Next() bool {
	if it.entries == nil {
		it.loadNextShard()
		if it.entries == nil {
			return false
		}
	}

	it.position++
	if it.position >= len(it.entries) {
		it.loadNextShard()
		return it.Next()
	}
	return true
}

// Value returns the current key-value pair
func (it *ConcurrentHashMapIterator[K, V]) Value() (Entry[K, V], error) {
	if it.entries == nil || it.position >= len(it.entries) {
		var zero Entry[K, V]
		return zero, errors.New("no more elements")
	}
	return it.entries[it.position], nil
}
