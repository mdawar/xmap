// Package xmap provides a generic, thread-safe map with automatic key expiration.
package xmap

import (
	"sync"
	"sync/atomic"
	"time"
)

// entry is the value stored internally in the map.
type entry[V any] struct {
	value V         // The actual value stored.
	exp   time.Time // The expiration time.
}

// Config represents the Map configuration.
type Config struct {
	// CleanupInterval is the interval at which the expired keys are removed.
	CleanupInterval time.Duration
}

// Map is a thread-safe map with automatic key expiration.
type Map[K comparable, V any] struct {
	mu       sync.RWMutex    // Mutex to synchronize the map access.
	kv       map[K]*entry[V] // The underlying map.
	interval time.Duration   // Cleanup interval.
	stop     chan struct{}   // Channel closed on stop.
	stopped  atomic.Int32    // Map stopped flag.
}

// New creates a new Map instance with the default configuration.
func New[K comparable, V any]() *Map[K, V] {
	return NewWithConfig[K, V](Config{
		CleanupInterval: 5 * time.Minute,
	})
}

// NewWithConfig creates a new Map instance with the specified configuration.
func NewWithConfig[K comparable, V any](cfg Config) *Map[K, V] {
	m := &Map[K, V]{
		kv:       make(map[K]*entry[V]),
		stop:     make(chan struct{}),
		interval: cfg.CleanupInterval,
	}

	go m.cleanup()

	return m
}

// Stop halts the background cleanup goroutine and clears the map.
// It should be called when the map is no longer needed.
//
// This method is safe to be called multiple times.
//
// A stopped map should not be re-used, a new map should be created instead.
func (m *Map[K, V]) Stop() {
	if m.stopped.CompareAndSwap(0, 1) {
		close(m.stop)
	}
}

// Stopped reports whether the Map is stopped.
//
// Expired keys are not removed automatically in a stopped map.
func (m *Map[K, V]) Stopped() bool {
	return m.stopped.Load() == 1
}

// Len returns the length of the map.
//
// The length of the map is the total number of keys, including the
// expired keys that have not been removed yet.
func (m *Map[K, V]) Len() int {
	m.mu.RLock()
	defer m.mu.RUnlock()
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

	m.mu.Lock()
	m.kv[key] = &entry[V]{value, exp}
	m.mu.Unlock()
}

// Update changes the value of the key while preserving the expiration time.
//
// The return value reports whether there was an update (Key exists).
func (m *Map[K, V]) Update(key K, value V) bool {
	m.mu.Lock()
	defer m.mu.Unlock()

	if entry, ok := m.kv[key]; ok {
		entry.value = value
		return true
	}
	return false
}

// Get returns the value associated with the key.
//
// The second bool return value reports whether the key exists in the map.
func (m *Map[K, V]) Get(key K) (V, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if entry, ok := m.kv[key]; ok && !m.expired(entry) {
		return entry.value, true
	}

	var zero V
	return zero, false
}

// GetWithExpiration returns the value and expiration time of the key.
//
// The third bool return value reports whether the key exists in the map.
func (m *Map[K, V]) GetWithExpiration(key K) (V, time.Time, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if entry, ok := m.kv[key]; ok && !m.expired(entry) {
		return entry.value, entry.exp, true
	}

	var zero V
	return zero, time.Time{}, false
}

// Delete removes a key from the map.
func (m *Map[K, V]) Delete(key K) {
	m.mu.Lock()
	delete(m.kv, key)
	m.mu.Unlock()
}

// Clear removes all the entries from the map.
func (m *Map[K, V]) Clear() {
	m.mu.Lock()
	clear(m.kv)
	m.mu.Unlock()
}

// cleanup removes expired keys from the map in an interval.
//
// The cleanup is stopped by calling Stop.
func (m *Map[K, V]) cleanup() {
	ticker := time.NewTicker(m.interval)
	defer ticker.Stop()

	for {
		select {
		case <-m.stop:
			return
		case <-ticker.C:
			m.removeExpired()
		}
	}
}

// removeExpired checks the keys and removes the expired ones.
func (m *Map[K, V]) removeExpired() {
	// Expired keys.
	var expired []K

	// Find the expired keys.
	m.mu.RLock()
	for key, entry := range m.kv {
		if m.expired(entry) {
			expired = append(expired, key)
		}
	}
	m.mu.RUnlock()

	// Remove the expired keys.
	m.mu.Lock()
	for _, key := range expired {
		delete(m.kv, key)
	}
	m.mu.Unlock()
}

// expired reports whether an entry has expired.
func (m *Map[K, V]) expired(entry *entry[V]) bool {
	return !entry.exp.IsZero() && time.Now().After(entry.exp)
}
