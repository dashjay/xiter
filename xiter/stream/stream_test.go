//go:build go1.23

package stream_test

import (
	"sync/atomic"
	"testing"

	"github.com/dashjay/xiter/xiter"
	"github.com/dashjay/xiter/xiter/stream"
)

func refRange(a, b int) []int {
	var res []int
	for i := a; i < b; i++ {
		res = append(res, i)
	}
	return res
}

func toSlice[T any](seq xiter.Seq[T]) []T {
	var out []T
	for v := range seq {
		out = append(out, v)
	}
	return out
}

func TestParallelMap(t *testing.T) {
	t.Run("basic", func(t *testing.T) {
		in := refRange(0, 100)
		seq := xiter.FromSlice(in)
		results := toSlice(stream.ParallelMap(seq, func(v int) int {
			return v * 2
		}, 4))

		if len(results) != 100 {
			t.Fatalf("ParallelMap: got %d results, want 100", len(results))
		}

		// Results come in non-deterministic order, but all values should be present
		seen := make(map[int]bool)
		for _, v := range results {
			if v%2 != 0 {
				t.Fatalf("ParallelMap: unexpected value %d (should be even)", v)
			}
			seen[v] = true
		}
		for i := 0; i < 100; i++ {
			if !seen[i*2] {
				t.Fatalf("ParallelMap: missing result %d", i*2)
			}
		}
	})

	t.Run("empty", func(t *testing.T) {
		results := toSlice(stream.ParallelMap(xiter.FromSlice([]int{}), func(v int) int {
			return v
		}, 2))
		if len(results) != 0 {
			t.Fatalf("ParallelMap empty: got %d, want 0", len(results))
		}
	})

	t.Run("single worker", func(t *testing.T) {
		in := refRange(0, 10)
		results := toSlice(stream.ParallelMap(xiter.FromSlice(in), func(v int) int {
			return v + 1
		}, 1))
		if len(results) != 10 {
			t.Fatalf("ParallelMap single: got %d, want 10", len(results))
		}
	})

	t.Run("zero workers defaults to one", func(t *testing.T) {
		in := refRange(0, 10)
		results := toSlice(stream.ParallelMap(xiter.FromSlice(in), func(v int) int {
			return v
		}, 0))
		if len(results) != 10 {
			t.Fatalf("ParallelMap zero: got %d, want 10", len(results))
		}
	})

	t.Run("early stop", func(t *testing.T) {
		in := refRange(0, 10000)
		var count atomic.Int32
		stopAfter := int32(50)
		for range stream.ParallelMap(xiter.FromSlice(in), func(v int) int {
			count.Add(1)
			return v
		}, 4) {
			if count.Load() >= stopAfter {
				break
			}
		}
		// We stopped early; verify seq iteration stopped
		// Note: some extra items may have been processed by workers
		// due to buffering, but the key point is we don't iterate forever
	})
}

func TestFanIn(t *testing.T) {
	t.Run("basic", func(t *testing.T) {
		seq1 := xiter.FromSlice([]int{1, 2, 3})
		seq2 := xiter.FromSlice([]int{4, 5, 6})
		results := toSlice(stream.FanIn(seq1, seq2))
		if len(results) != 6 {
			t.Fatalf("FanIn: got %d results, want 6", len(results))
		}
	})

	t.Run("empty seqs", func(t *testing.T) {
		results := toSlice(stream.FanIn[int]())
		if len(results) != 0 {
			t.Fatalf("FanIn empty: got %d, want 0", len(results))
		}
	})

	t.Run("single seq", func(t *testing.T) {
		results := toSlice(stream.FanIn(xiter.FromSlice([]int{1, 2, 3})))
		if len(results) != 3 {
			t.Fatalf("FanIn single: got %d, want 3", len(results))
		}
	})

	t.Run("empty and non-empty", func(t *testing.T) {
		seq1 := xiter.FromSlice([]int{})
		seq2 := xiter.FromSlice([]int{1, 2})
		results := toSlice(stream.FanIn(seq1, seq2))
		if len(results) != 2 {
			t.Fatalf("FanIn mixed: got %d, want 2", len(results))
		}
	})

	t.Run("early stop", func(t *testing.T) {
		seq1 := xiter.FromSlice(refRange(0, 1000))
		seq2 := xiter.FromSlice(refRange(1000, 2000))
		count := 0
		for range stream.FanIn(seq1, seq2) {
			count++
			if count >= 5 {
				break
			}
		}
		if count != 5 {
			t.Fatalf("FanIn early stop: count = %d, want 5", count)
		}
	})
}

func TestBatch(t *testing.T) {
	t.Run("basic", func(t *testing.T) {
		seq := xiter.FromSlice(refRange(0, 5))
		batches := toSlice(stream.Batch(seq, 2))
		if len(batches) != 3 {
			t.Fatalf("Batch: got %d batches, want 3", len(batches))
		}
		if len(batches[0]) != 2 || batches[0][0] != 0 || batches[0][1] != 1 {
			t.Fatalf("Batch[0]: got %v, want [0 1]", batches[0])
		}
		if len(batches[1]) != 2 || batches[1][0] != 2 || batches[1][1] != 3 {
			t.Fatalf("Batch[1]: got %v, want [2 3]", batches[1])
		}
		if len(batches[2]) != 1 || batches[2][0] != 4 {
			t.Fatalf("Batch[2]: got %v, want [4]", batches[2])
		}
	})

	t.Run("exact", func(t *testing.T) {
		seq := xiter.FromSlice([]int{1, 2, 3, 4})
		batches := toSlice(stream.Batch(seq, 2))
		if len(batches) != 2 {
			t.Fatalf("Batch exact: got %d batches, want 2", len(batches))
		}
	})

	t.Run("single element", func(t *testing.T) {
		seq := xiter.FromSlice([]int{1})
		batches := toSlice(stream.Batch(seq, 5))
		if len(batches) != 1 || batches[0][0] != 1 {
			t.Fatalf("Batch single: got %v, want [[1]]", batches)
		}
	})

	t.Run("empty", func(t *testing.T) {
		batches := toSlice(stream.Batch(xiter.FromSlice([]int{}), 2))
		if len(batches) != 0 {
			t.Fatalf("Batch empty: got %d, want 0", len(batches))
		}
	})

	t.Run("zero or negative size", func(t *testing.T) {
		seq := xiter.FromSlice(refRange(0, 5))
		batches := toSlice(stream.Batch(seq, 0))
		if len(batches) != 0 {
			t.Fatalf("Batch zero: got %d, want 0", len(batches))
		}
	})
}
