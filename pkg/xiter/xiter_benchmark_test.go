package xiter_test

import (
	"testing"

	"github.com/dashjay/xiter/pkg/xiter"
)

func BenchmarkCycle(b *testing.B) {
	seq := xiter.Cycle(xiter.FromSlice(_range(0, 1000)))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		xiter.ToSlice(xiter.Limit(seq, 10000))
	}
}

func BenchmarkGenerate(b *testing.B) {
	i := 0
	gen := xiter.Generate(func() int {
		i++
		return i
	})
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		xiter.ToSlice(xiter.Limit(gen, 10000))
	}
}

func BenchmarkToChan(b *testing.B) {
	seq := xiter.FromSlice(_range(0, 10000))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ch := xiter.ToChan(seq)
		for range ch {
		}
	}
}

func BenchmarkRange(b *testing.B) {
	for i := 0; i < b.N; i++ {
		xiter.ToSlice(xiter.Range(0, 10000, 1))
	}
}

func BenchmarkWithIndex(b *testing.B) {
	seq := xiter.FromSlice(_range(0, 10000))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		xiter.WithIndex(seq)(func(_ int, v int) bool {
			_ = v
			return true
		})
	}
}
