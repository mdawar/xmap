package xmap_test

import (
	"testing"
	"time"

	"github.com/mdawar/xmap"
)

func BenchmarkMapSetIntValue(b *testing.B) {
	m := xmap.New[string, int]()
	defer m.Stop()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m.Set("keyName", 100, time.Hour)
	}
	b.StopTimer()
}

func BenchmarkMapSetIntValueParallel(b *testing.B) {
	m := xmap.New[string, int]()
	defer m.Stop()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			m.Set("keyName", 100, time.Hour)
		}
	})
}

func BenchmarkMapSetStringValue(b *testing.B) {
	m := xmap.New[string, string]()
	defer m.Stop()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m.Set("keyName", "test value", time.Hour)
	}
	b.StopTimer()
}

func BenchmarkMapSetStringValueParallel(b *testing.B) {
	m := xmap.New[string, string]()
	defer m.Stop()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			m.Set("keyName", "test value", time.Hour)
		}
	})
}

func BenchmarkMapGetIntValue(b *testing.B) {
	m := xmap.New[string, int]()
	defer m.Stop()

	m.Set("keyName", 100, time.Hour)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, ok := m.Get("keyName"); !ok {
			b.Fatal("key does not exist")
		}
	}
	b.StopTimer()
}

func BenchmarkMapGetIntValueParallel(b *testing.B) {
	m := xmap.New[string, int]()
	defer m.Stop()

	m.Set("keyName", 100, time.Hour)

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			if _, ok := m.Get("keyName"); !ok {
				b.Fatal("key does not exist")
			}
		}
	})
}

func BenchmarkMapGetStringValue(b *testing.B) {
	m := xmap.New[string, string]()
	defer m.Stop()

	m.Set("keyName", "test value", time.Hour)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, ok := m.Get("keyName"); !ok {
			b.Fatal("key does not exist")
		}
	}
	b.StopTimer()
}

func BenchmarkMapGetStringValueParallel(b *testing.B) {
	m := xmap.New[string, string]()
	defer m.Stop()

	m.Set("keyName", "test value", time.Hour)

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			if _, ok := m.Get("keyName"); !ok {
				b.Fatal("key does not exist")
			}
		}
	})
}
