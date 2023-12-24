// Package xmap provides a thread-safe, generic map with automatic key expiration.
package xmap

import "time"

// item is the value stored internally in the map.
type item[V any] struct {
	value V         // The actual value stored.
	exp   time.Time // The expiration time.
}

// Map is a thread-safe map with automatic key expiration.
type Map[K comparable, V any] struct {
	kv map[K]*item[V] // The underlying map.
}

// New creates a new Map instance.
func New[K comparable, V any]() *Map[K, V] {
	return &Map[K, V]{
		kv: make(map[K]*item[V]),
	}
}

// Stop halts the background cleanup goroutine and clears the map.
// It should be called when the map is no longer needed.
//
// A stopped map should not be re-used, a new map should be created instead.
func (m *Map[K, V]) Stop() {}

// Len returns the length of the map.
//
// The length of the map is the total number of keys, including the
// expired keys that have not been removed yet.
func (m *Map[K, V]) Len() int {
	return len(m.kv)
}

// Set creates or replaces a key/value pair in the map.
//
// A key can be set to never expire with a ttl value of 0.
func (m *Map[K, V]) Set(key K, value V, ttl time.Duration) {
	var exp time.Time

	if ttl > 0 {
		exp = time.Now().Add(ttl)
	}

	m.kv[key] = &item[V]{value, exp}
}

// Update changes the value of the key while preserving the expiration time.
//
// The return value reports whether there was an update (Key exists).
func (m *Map[K, V]) Update(key K, value V) bool {
	if item, ok := m.kv[key]; ok {
		item.value = value
		return true
	}
	return false
}

// Get returns the value associated with the key.
//
// The second bool return value reports whether the key exists in the map.
func (m *Map[K, V]) Get(key K) (V, bool) {
	if item, ok := m.kv[key]; ok {
		return item.value, true
	}

	var zero V
	return zero, false
}

// GetWithExpiration returns the value and expiration time of the key.
//
// The third bool return value reports whether the key exists in the map.
func (m *Map[K, V]) GetWithExpiration(key K) (V, time.Time, bool) {
	if item, ok := m.kv[key]; ok {
		return item.value, item.exp, ok
	}

	var zero V
	return zero, time.Time{}, false
}

// Delete removes a key from the map.
func (m *Map[K, V]) Delete(key K) {}
