//go:build go1.23

package rt_test

import (
	"testing"
	"time"

	"github.com/dashjay/xiter/xiter"
	"github.com/dashjay/xiter/xiter/rt"
)

func TestTicker(t *testing.T) {
	t.Run("basic", func(t *testing.T) {
		count := 0
		for range rt.Ticker(time.Millisecond) {
			count++
			if count >= 3 {
				break
			}
		}
		if count != 3 {
			t.Fatalf("Ticker: got %d ticks, want 3", count)
		}
	})


	t.Run("early stop", func(t *testing.T) {
		start := time.Now()
		count := 0
		for range rt.Ticker(time.Millisecond) {
			count++
			if count >= 2 {
				break
			}
		}
		elapsed := time.Since(start)
		if count != 2 {
			t.Fatalf("Ticker early stop: got %d ticks, want 2", count)
		}
		if elapsed > 100*time.Millisecond {
			t.Fatalf("Ticker early stop took too long: %v", elapsed)
		}
	})
}

func TestThrottle(t *testing.T) {
	t.Run("basic", func(t *testing.T) {
		seq := xiter.FromSlice([]int{1, 2, 3})
		start := time.Now()
		results := xiter.ToSlice(rt.Throttle(seq, time.Millisecond))
		elapsed := time.Since(start)
		if len(results) != 3 {
			t.Fatalf("Throttle: got %d results, want 3", len(results))
		}
		if elapsed < 2*time.Millisecond {
			t.Fatalf("Throttle: too fast, elapsed = %v", elapsed)
		}
	})

	t.Run("empty", func(t *testing.T) {
		results := xiter.ToSlice(rt.Throttle(xiter.FromSlice([]int{}), time.Millisecond))
		if len(results) != 0 {
			t.Fatalf("Throttle empty: got %d, want 0", len(results))
		}
	})

	t.Run("single element", func(t *testing.T) {
		results := xiter.ToSlice(rt.Throttle(xiter.FromSlice([]int{42}), time.Millisecond))
		if len(results) != 1 || results[0] != 42 {
			t.Fatalf("Throttle single: got %v, want [42]", results)
		}
	})

	t.Run("early stop", func(t *testing.T) {
		seq := xiter.FromSlice([]int{1, 2, 3, 4, 5})
		count := 0
		for range rt.Throttle(seq, time.Millisecond) {
			count++
			if count >= 2 {
				break
			}
		}
		if count != 2 {
			t.Fatalf("Throttle early stop: count = %d, want 2", count)
		}
	})
}

func TestDebounce(t *testing.T) {
	t.Run("single value", func(t *testing.T) {
		seq := xiter.FromSlice([]int{1})
		results := xiter.ToSlice(rt.Debounce(seq, time.Microsecond))
		if len(results) != 1 || results[0] != 1 {
			t.Fatalf("Debounce single: got %v, want [1]", results)
		}
	})

	t.Run("rapid values only emit last", func(t *testing.T) {
		// Values from a slice arrive faster than the debounce period,
		// and verify only the last one is emitted (after the debounce period)
		results := xiter.ToSlice(rt.Debounce(
			xiter.FromSlice([]int{1, 2, 3, 4, 5}),
			50*time.Millisecond,
		))
		if len(results) != 1 || results[0] != 5 {
			t.Fatalf("Debounce rapid: got %v, want [5] (len=%d)", results, len(results))
		}
	})

	t.Run("empty", func(t *testing.T) {
		results := xiter.ToSlice(rt.Debounce(xiter.FromSlice([]int{}), time.Microsecond))
		if len(results) != 0 {
			t.Fatalf("Debounce empty: got %d, want 0", len(results))
		}
	})

		t.Run("early stop", func(t *testing.T) {
			seq := xiter.FromSlice([]int{1, 2, 3})
			count := 0
			for range rt.Debounce(seq, time.Microsecond) {
				count++
				if count >= 1 {
					break
				}
			}
			if count != 1 {
				t.Fatalf("Debounce early stop: count = %d, want 1", count)
			}
		})
}
