package xslice_test

import (
	"github.com/samber/lo"
	"testing"

	"github.com/dashjay/xiter/pkg/xslice"
)

func BenchmarkSlice(b *testing.B) {
	const length = 1_000_000
	b.Run("benchmark all", func(b *testing.B) {
		seq := _range(1, length)
		fn := func(i int) bool {
			return i != length-1
		}

		b.Run("baseline", func(b *testing.B) {
			for j := 0; j < b.N; j++ {
				for i := 0; i < len(seq); i++ {
					if !fn(seq[i]) {
						break
					}
				}
			}
		})
		b.Run("xslice", func(b *testing.B) {
			for j := 0; j < b.N; j++ {
				_ = xslice.All(seq, fn)
			}
		})

	})

	b.Run("benchmark any", func(b *testing.B) {
		seq := _range(1, length)
		fn := func(i int) bool {
			return i == length-1
		}

		b.Run("baseline", func(b *testing.B) {
			for j := 0; j < b.N; j++ {
				for i := 0; i < len(seq); i++ {
					_ = fn(seq[i])
				}
			}
		})

		b.Run("xslice", func(b *testing.B) {
			for j := 0; j < b.N; j++ {
				_ = xslice.Any(seq, fn)
			}
		})
	})

	b.Run("benchmark avg", func(b *testing.B) {
		seq := _range(1, length)
		b.Run("baseline", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_ = avg(seq)
			}

		})
		b.Run("xslice", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_ = xslice.Avg(seq)
			}
		})
	})

	b.Run("benchmark contain", func(b *testing.B) {
		seq := _range(1, length)
		b.Run("xslice", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				xslice.Contains(seq, length/2)
			}
		})
		b.Run("lo", func(b *testing.B) {
			for j := 0; j < b.N; j++ {
				lo.Contains(seq, length/2)
			}
		})
	})

	b.Run("benchmark sum", func(b *testing.B) {
		seq := _range(1, length)
		b.Run("xslice", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_ = xslice.Sum(seq)
			}
		})

		b.Run("lo", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_ = lo.Sum(seq)
			}
		})
	})
}

func BenchmarkChunk(b *testing.B) {
	b.Run("xslice", func(b *testing.B) {
		arr := _range(0, 1000)
		for i := 1; i < b.N; i++ {
			xslice.Chunk(arr, i)
		}
	})

	b.Run("xslice-inplace", func(b *testing.B) {
		arr := _range(0, 1000)
		for i := 1; i < b.N; i++ {
			xslice.ChunkInPlace(arr, i)
		}
	})

	b.Run("lo", func(b *testing.B) {
		arr := _range(0, 1000)
		for i := 1; i < b.N; i++ {
			lo.Chunk(arr, i)
		}
	})
}
