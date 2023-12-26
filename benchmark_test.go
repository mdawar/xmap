package xmap_test

import (
	"testing"
	"time"

	"github.com/mdawar/xmap"
)

func BenchmarkMapSet(b *testing.B) {
	b.Run("int", func(b *testing.B) {
		b.Run("serial", func(b *testing.B) {
			m := xmap.New[string, int]()
			defer m.Stop()

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				m.Set("keyName", 100, time.Hour)
			}
			b.StopTimer()
		})

		b.Run("parallel", func(b *testing.B) {
			m := xmap.New[string, int]()
			defer m.Stop()

			b.ResetTimer()
			b.RunParallel(func(pb *testing.PB) {
				for pb.Next() {
					m.Set("keyName", 100, time.Hour)
				}
			})
			b.StopTimer()
		})
	})

	b.Run("string", func(b *testing.B) {
		b.Run("serial", func(b *testing.B) {
			m := xmap.New[string, string]()
			defer m.Stop()

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				m.Set("keyName", "test value", time.Hour)
			}
			b.StopTimer()
		})

		b.Run("parallel", func(b *testing.B) {
			m := xmap.New[string, string]()
			defer m.Stop()

			b.ResetTimer()
			b.RunParallel(func(pb *testing.PB) {
				for pb.Next() {
					m.Set("keyName", "test value", time.Hour)
				}
			})
			b.StopTimer()
		})
	})
}

func BenchmarkMapInitialCapacitySet(b *testing.B) {
	b.Run("int", func(b *testing.B) {
		b.Run("serial", func(b *testing.B) {
			m := xmap.NewWithConfig[string, int](xmap.Config{
				InitialCapacity: 10_000_000,
			})
			defer m.Stop()

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				m.Set("keyName", 100, time.Hour)
			}
			b.StopTimer()
		})

		b.Run("parallel", func(b *testing.B) {
			m := xmap.NewWithConfig[string, int](xmap.Config{
				InitialCapacity: 10_000_000,
			})
			defer m.Stop()

			b.ResetTimer()
			b.RunParallel(func(pb *testing.PB) {
				for pb.Next() {
					m.Set("keyName", 100, time.Hour)
				}
			})
			b.StopTimer()
		})
	})

	b.Run("string", func(b *testing.B) {
		b.Run("serial", func(b *testing.B) {
			m := xmap.NewWithConfig[string, string](xmap.Config{
				InitialCapacity: 10_000_000,
			})
			defer m.Stop()

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				m.Set("keyName", "test value", time.Hour)
			}
			b.StopTimer()
		})

		b.Run("parallel", func(b *testing.B) {
			m := xmap.NewWithConfig[string, string](xmap.Config{
				InitialCapacity: 10_000_000,
			})
			defer m.Stop()

			b.ResetTimer()
			b.RunParallel(func(pb *testing.PB) {
				for pb.Next() {
					m.Set("keyName", "test value", time.Hour)
				}
			})
			b.StopTimer()
		})
	})
}

func BenchmarkMapGet(b *testing.B) {
	b.Run("int", func(b *testing.B) {
		b.Run("serial", func(b *testing.B) {
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
		})

		b.Run("parallel", func(b *testing.B) {
			m := xmap.New[string, int]()
			defer m.Stop()

			m.Set("keyName", 100, time.Hour)

			b.ResetTimer()
			b.RunParallel(func(pb *testing.PB) {
				for pb.Next() {
					if _, ok := m.Get("keyName"); !ok {
						b.Fatal("key does not exist")
					}
				}
			})
			b.StopTimer()
		})
	})

	b.Run("string", func(b *testing.B) {
		b.Run("serial", func(b *testing.B) {
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
		})

		b.Run("parallel", func(b *testing.B) {
			m := xmap.New[string, string]()
			defer m.Stop()

			m.Set("keyName", "test value", time.Hour)

			b.ResetTimer()
			b.RunParallel(func(pb *testing.PB) {
				for pb.Next() {
					if _, ok := m.Get("keyName"); !ok {
						b.Fatal("key does not exist")
					}
				}
			})
			b.StopTimer()
		})
	})
}

func BenchmarkMapGetWithExpiration(b *testing.B) {
	b.Run("int", func(b *testing.B) {
		b.Run("serial", func(b *testing.B) {
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
		})

		b.Run("parallel", func(b *testing.B) {
			m := xmap.New[string, int]()
			defer m.Stop()

			m.Set("keyName", 100, time.Hour)

			b.ResetTimer()
			b.RunParallel(func(pb *testing.PB) {
				for pb.Next() {
					if _, _, ok := m.GetWithExpiration("keyName"); !ok {
						b.Fatal("key does not exist")
					}
				}
			})
			b.StopTimer()
		})
	})

	b.Run("string", func(b *testing.B) {
		b.Run("serial", func(b *testing.B) {
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
		})

		b.Run("parallel", func(b *testing.B) {
			m := xmap.New[string, string]()
			defer m.Stop()

			m.Set("keyName", "test value", time.Hour)

			b.ResetTimer()
			b.RunParallel(func(pb *testing.PB) {
				for pb.Next() {
					if _, _, ok := m.GetWithExpiration("keyName"); !ok {
						b.Fatal("key does not exist")
					}
				}
			})
			b.StopTimer()
		})
	})
}

func BenchmarkMapUpdate(b *testing.B) {
	b.Run("int", func(b *testing.B) {
		b.Run("serial", func(b *testing.B) {
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
		})

		b.Run("parallel", func(b *testing.B) {
			m := xmap.New[string, int]()
			defer m.Stop()

			m.Set("keyName", 100, time.Hour)

			b.ResetTimer()
			b.RunParallel(func(pb *testing.PB) {
				for pb.Next() {
					if ok := m.Update("keyName", 1000); !ok {
						b.Fatal("key does not exist")
					}
				}
			})
			b.StopTimer()
		})
	})

	b.Run("string", func(b *testing.B) {
		b.Run("serial", func(b *testing.B) {
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
		})

		b.Run("parallel", func(b *testing.B) {
			m := xmap.New[string, string]()
			defer m.Stop()

			m.Set("keyName", "test value", time.Hour)

			b.ResetTimer()
			b.RunParallel(func(pb *testing.PB) {
				for pb.Next() {
					if ok := m.Update("keyName", "updated value"); !ok {
						b.Fatal("key does not exist")
					}
				}
			})
			b.StopTimer()
		})
	})
}
