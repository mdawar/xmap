package xmap_test

import (
	"testing"
	"time"

	"github.com/mdawar/xmap"
)

func BenchmarkMapSetInt(b *testing.B) {
	m := xmap.New[string, int]()
	defer m.Stop()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m.Set("keyName", 100, time.Hour)
	}
	b.StopTimer()
}

func BenchmarkMapSetIntInitialCapacity(b *testing.B) {
	m := xmap.NewWithConfig[string, int](xmap.Config{
		InitialCapacity: 10_000_000,
	})
	defer m.Stop()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m.Set("keyName", 100, time.Hour)
	}
	b.StopTimer()
}

func BenchmarkMapSetIntParallel(b *testing.B) {
	m := xmap.New[string, int]()
	defer m.Stop()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			m.Set("keyName", 100, time.Hour)
		}
	})
}

func BenchmarkMapSetIntParallelInitialCapacity(b *testing.B) {
	m := xmap.NewWithConfig[string, int](xmap.Config{
		InitialCapacity: 10_000_000,
	})
	defer m.Stop()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			m.Set("keyName", 100, time.Hour)
		}
	})
}

func BenchmarkMapSetString(b *testing.B) {
	m := xmap.New[string, string]()
	defer m.Stop()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m.Set("keyName", "test value", time.Hour)
	}
	b.StopTimer()
}

func BenchmarkMapSetStringInitialCapacity(b *testing.B) {
	m := xmap.NewWithConfig[string, string](xmap.Config{
		InitialCapacity: 10_000_000,
	})
	defer m.Stop()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m.Set("keyName", "test value", time.Hour)
	}
	b.StopTimer()
}

func BenchmarkMapSetStringParallel(b *testing.B) {
	m := xmap.New[string, string]()
	defer m.Stop()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			m.Set("keyName", "test value", time.Hour)
		}
	})
}

func BenchmarkMapSetStringParallelInitialCapacity(b *testing.B) {
	m := xmap.NewWithConfig[string, string](xmap.Config{
		InitialCapacity: 10_000_000,
	})
	defer m.Stop()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			m.Set("keyName", "test value", time.Hour)
		}
	})
}

func BenchmarkMapGetInt(b *testing.B) {
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

func BenchmarkMapGetIntParallel(b *testing.B) {
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

func BenchmarkMapGetString(b *testing.B) {
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

func BenchmarkMapGetStringParallel(b *testing.B) {
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

func BenchmarkMapGetWithExpirationInt(b *testing.B) {
	m := xmap.New[string, int]()
	defer m.Stop()

	m.Set("keyName", 100, time.Hour)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, _, ok := m.GetWithExpiration("keyName"); !ok {
			b.Fatal("key does not exist")
		}
	}
	b.StopTimer()
}

func BenchmarkMapGetWithExpirationIntParallel(b *testing.B) {
	m := xmap.New[string, int]()
	defer m.Stop()

	m.Set("keyName", 100, time.Hour)

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			if _, _, ok := m.GetWithExpiration("keyName"); !ok {
				b.Fatal("key does not exist")
			}
		}
	})
}

func BenchmarkMapGetWithExpirationString(b *testing.B) {
	m := xmap.New[string, string]()
	defer m.Stop()

	m.Set("keyName", "test value", time.Hour)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, _, ok := m.GetWithExpiration("keyName"); !ok {
			b.Fatal("key does not exist")
		}
	}
	b.StopTimer()
}

func BenchmarkMapGetWithExpirationStringParallel(b *testing.B) {
	m := xmap.New[string, string]()
	defer m.Stop()

	m.Set("keyName", "test value", time.Hour)

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			if _, _, ok := m.GetWithExpiration("keyName"); !ok {
				b.Fatal("key does not exist")
			}
		}
	})
}

func BenchmarkMapUpdateInt(b *testing.B) {
	m := xmap.New[string, int]()
	defer m.Stop()

	m.Set("keyName", 100, time.Hour)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if ok := m.Update("keyName", 1000); !ok {
			b.Fatal("key does not exist")
		}
	}
	b.StopTimer()
}

func BenchmarkMapUpdateIntParallel(b *testing.B) {
	m := xmap.New[string, int]()
	defer m.Stop()

	m.Set("keyName", 100, time.Hour)

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			if ok := m.Update("keyName", 1000); !ok {
				b.Fatal("key does not exist")
			}
		}
	})
}

func BenchmarkMapUpdateString(b *testing.B) {
	m := xmap.New[string, string]()
	defer m.Stop()

	m.Set("keyName", "test value", time.Hour)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if ok := m.Update("keyName", "updated value"); !ok {
			b.Fatal("key does not exist")
		}
	}
	b.StopTimer()
}

func BenchmarkMapUpdateStringParallel(b *testing.B) {
	m := xmap.New[string, string]()
	defer m.Stop()

	m.Set("keyName", "test value", time.Hour)

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			if ok := m.Update("keyName", "updated value"); !ok {
				b.Fatal("key does not exist")
			}
		}
	})
}
