// Package xmap provides a generic, thread-safe map with automatic key expiration.
package xmap

import (
	"iter"
	"sync"
	"sync/atomic"
	"time"
)

// entry is the value stored internally in the [Map].
type entry[V any] struct {
	value V         // The actual value stored.
	exp   time.Time // The expiration time of the value.
}

// Config represents the [Map] configuration.
type Config struct {
	// CleanupInterval is the interval at which the expired keys are removed.
	// Default: 5 minutes.
	CleanupInterval time.Duration
	// InitialCapacity is the initial capacity hint passed to make when creating
	// the map. It does not bound the size of the map, It will create a map with
	// an initial space to hold the specified number of elements.
	InitialCapacity int
	// TimeSource is the time source used by the map for key expiration.
	// This is only useful for testing.
	// Default: system time.
	TimeSource Time
}

// setDefaults sets the default values for the [Map] configuration.
func (c *Config) setDefaults() {
	if c.CleanupInterval == 0 {
		c.CleanupInterval = 5 * time.Minute
	}

	if c.TimeSource == nil {
		c.TimeSource = &systemTime{}
	}
}

// Map is a thread-safe map with automatic key expiration.
type Map[K comparable, V any] struct {
	mu       sync.RWMutex    // Mutex to synchronize the map access.
	kv       map[K]*entry[V] // The underlying map.
	interval time.Duration   // Cleanup interval.
	time     Time            // Time source.
	stop     chan struct{}   // Channel closed on stop.
	active   atomic.Int32    // Cleanup active flag.
	stopped  atomic.Int32    // Map stopped flag.
}

// New creates a new [Map] instance with the default configuration.
func New[K comparable, V any]() *Map[K, V] {
	return NewWithConfig[K, V](Config{})
}

// NewWithConfig creates a new [Map] instance with the specified configuration.
func NewWithConfig[K comparable, V any](cfg Config) *Map[K, V] {
	cfg.setDefaults()

	m := &Map[K, V]{
		kv:       make(map[K]*entry[V], cfg.InitialCapacity),
		stop:     make(chan struct{}),
		interval: cfg.CleanupInterval,
		time:     cfg.TimeSource,
	}

	go m.cleanup()

	return m
}

// Stop halts the background cleanup goroutine and clears the [Map].
// It should be called when the [Map] is no longer needed.
//
// This method is safe to be called multiple times.
//
// A stopped [Map] should not be re-used, a new [Map] should be created instead.
func (m *Map[K, V]) Stop() {
	if m.stopped.CompareAndSwap(0, 1) {
		// Stop the cleanup goroutine.
		close(m.stop)

		// Clear the map to free up resources.
		m.mu.Lock()
		m.kv = make(map[K]*entry[V])
		m.mu.Unlock()
	}
}

// Stopped reports whether the [Map] is stopped.
//
// Expired keys are not removed automatically in a stopped [Map].
func (m *Map[K, V]) Stopped() bool {
	return m.stopped.Load() == 1
}

// Len returns the length of the [Map].
//
// The length of the [Map] is the total number of keys, including the expired
// keys that have not been removed yet.
//
// To get the length excluding the number of expired keys, call [Map.RemoveExpired]
// before calling this method.
func (m *Map[K, V]) Len() int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return len(m.kv)
}

// Set creates or replaces a key-value pair in the [Map].
//
// A key can be set to never expire with a ttl value of 0.
func (m *Map[K, V]) Set(key K, value V, ttl time.Duration) {
	var exp time.Time

	if ttl > 0 {
		exp = m.time.Now().Add(ttl)
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

	if entry, ok := m.kv[key]; ok && !m.expired(entry) {
		entry.value = value
		return true
	}
	return false
}

// Get returns the value associated with the key.
//
// The second bool return value reports whether the key exists in the [Map].
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
// The third bool return value reports whether the key exists in the [Map].
func (m *Map[K, V]) GetWithExpiration(key K) (V, time.Time, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if entry, ok := m.kv[key]; ok && !m.expired(entry) {
		return entry.value, entry.exp, true
	}

	var zero V
	return zero, time.Time{}, false
}

// All returns an iterator over key-value pairs from the [Map].
//
// Only the entries that have not expired are produced during the iteration.
//
// Similar to the map type, the iteration order is not guaranteed.
func (m *Map[K, V]) All() iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		m.mu.RLock()
		defer m.mu.RUnlock()

		for key, entry := range m.kv {
			if !m.expired(entry) {
				if !yield(key, entry.value) {
					return
				}
			}
		}
	}
}

// Delete removes a key from the [Map].
func (m *Map[K, V]) Delete(key K) {
	m.mu.Lock()
	delete(m.kv, key)
	m.mu.Unlock()
}

// Clear removes all the entries from the [Map].
func (m *Map[K, V]) Clear() {
	m.mu.Lock()
	clear(m.kv)
	m.mu.Unlock()
}

// cleanup removes expired keys from the [Map] in an interval.
//
// The cleanup is stopped by calling [Map.Stop].
func (m *Map[K, V]) cleanup() {
	ticker := m.time.NewTicker(m.interval)
	defer ticker.Stop()

	// Set as active.
	m.active.Store(1)
	defer m.active.Store(0)

	for {
		select {
		case <-m.stop:
			return
		case <-ticker.C():
			m.RemoveExpired()
		}
	}
}

// CleanupActive reports whether the cleanup goroutine is active.
func (m *Map[K, V]) CleanupActive() bool {
	return m.active.Load() == 1
}

// RemoveExpired checks the [Map] keys and removes the expired ones.
//
// It returns the number of keys that were removed.
func (m *Map[K, V]) RemoveExpired() int {
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

	return len(expired)
}

// expired reports whether an [entry] has expired.
func (m *Map[K, V]) expired(entry *entry[V]) bool {
	return !entry.exp.IsZero() && m.time.Now().After(entry.exp)
}
