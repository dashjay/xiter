package xmap_test

import (
	"testing"

	"github.com/dashjay/xiter/xmap"
)

var _ = _map // suppress unused lint

func BenchmarkXmap(b *testing.B) {
	const size = 10_000

	m := _map(0, size)

	b.Run("keys", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = xmap.Keys(m)
		}
	})

	b.Run("values", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = xmap.Values(m)
		}
	})

	b.Run("to union slice", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = xmap.ToUnionSlice(m)
		}
	})

	b.Run("filter", func(b *testing.B) {
		fn := func(k int, v string) bool { return k%2 == 0 }
		for i := 0; i < b.N; i++ {
			_ = xmap.Filter(m, fn)
		}
	})

	b.Run("map values", func(b *testing.B) {
		fn := func(k int, v string) string { return v + "_x" }
		for i := 0; i < b.N; i++ {
			_ = xmap.MapValues(m, fn)
		}
	})

	b.Run("map keys", func(b *testing.B) {
		fn := func(k int, v string) int { return k * 2 }
		for i := 0; i < b.N; i++ {
			_ = xmap.MapKeys(m, fn)
		}
	})

	b.Run("coalesce maps", func(b *testing.B) {
		maps := make([]map[int]string, 10)
		for i := 0; i < 10; i++ {
			maps[i] = _map(i*1000, (i+1)*1000)
		}
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = xmap.CoalesceMaps(maps...)
		}
	})

	b.Run("max key", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = xmap.MaxKey(m)
		}
	})

	b.Run("min key", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = xmap.MinKey(m)
		}
	})

	b.Run("max value", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = xmap.MaxValue(m)
		}
	})

	b.Run("min value", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = xmap.MinValue(m)
		}
	})
}
