package xmap

import "time"

type Map[K comparable, V any] struct{}

func New[K comparable, V any]() *Map[K, V] {
	return &Map[K, V]{}
}

func (m *Map[K, V]) Len() int {
	return 0
}

func (m *Map[K, V]) Set(key K, value V, ttl time.Duration) {}

func (m *Map[K, V]) Update(key K, value V) bool {
	return false
}

func (m *Map[K, V]) Get(key K) (V, bool) {
	var zero V
	return zero, false
}

func (m *Map[K, V]) GetWithExpiration(key K) (V, time.Time, bool) {
	var zero V
	return zero, time.Time{}, false
}

func (m *Map[K, V]) Delete(key K) {}
