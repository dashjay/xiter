package xslice_test

import (
	"bytes"
	"testing"

	"github.com/samber/lo"

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

	b.Run("benchmark uniq", func(b *testing.B) {
		bytes := bytes.Repeat([]byte("b"), 1024)
		b.Run("xslice", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_ = xslice.Uniq(bytes)
			}
		})

		b.Run("lo", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_ = lo.Uniq(bytes)
			}
		})
	})

	b.Run("benchmark group by", func(b *testing.B) {
		arr := _range(0, 1000)
		fn := func(i int) string {
			if i%2 == 0 {
				return "even"
			}
			return "odd"
		}
		b.Run("xslice", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				xslice.GroupBy(arr, fn)
			}
		})
		b.Run("lo", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				lo.GroupBy(arr, fn)
			}
		})
	})

	b.Run("benchmark group by map", func(b *testing.B) {
		arr := _range(0, 1000)
		fn := func(i int) (string, int) {
			if i%2 == 0 {
				return "even", i * i
			}
			return "odd", i * i
		}
		b.Run("xslice", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				xslice.GroupByMap(arr, fn)
			}
		})
		b.Run("lo", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				lo.GroupByMap(arr, fn)
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
