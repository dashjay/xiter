//go:build go1.23

package collect_test

import (
	"testing"

	"github.com/dashjay/xiter/xiter"
	"github.com/dashjay/xiter/xiter/collect"
)

func toSlice[T any](seq xiter.Seq[T]) []T {
	var out []T
	for v := range seq {
		out = append(out, v)
	}
	return out
}

func TestGroupBy(t *testing.T) {
	t.Run("basic", func(t *testing.T) {
		seq := xiter.FromSlice([]string{"apple", "avocado", "banana", "blueberry", "cherry"})
		groups := toSlice(collect.GroupBy(seq, func(s string) string {
			return string(s[0])
		}))

		if len(groups) != 3 {
			t.Fatalf("GroupBy: got %d groups, want 3", len(groups))
		}

		for _, g := range groups {
			switch g.Key {
			case "a":
				if len(g.Items) != 2 {
					t.Fatalf("GroupBy 'a': got %v, want 2 items", g.Items)
				}
			case "b":
				if len(g.Items) != 2 {
					t.Fatalf("GroupBy 'b': got %v, want 2 items", g.Items)
				}
			case "c":
				if len(g.Items) != 1 || g.Items[0] != "cherry" {
					t.Fatalf("GroupBy 'c': got %v, want [cherry]", g.Items)
				}
			default:
				t.Fatalf("GroupBy: unexpected key %q", g.Key)
			}
		}
	})

	t.Run("empty", func(t *testing.T) {
		groups := toSlice(collect.GroupBy(xiter.FromSlice([]int{}), func(v int) int {
			return v % 2
		}))
		if len(groups) != 0 {
			t.Fatalf("GroupBy empty: got %d groups, want 0", len(groups))
		}
	})

	t.Run("single element", func(t *testing.T) {
		groups := toSlice(collect.GroupBy(xiter.FromSlice([]int{42}), func(v int) int {
			return v % 2
		}))
		if len(groups) != 1 || groups[0].Key != 0 || len(groups[0].Items) != 1 || groups[0].Items[0] != 42 {
			t.Fatalf("GroupBy single: got %v", groups)
		}
	})

	t.Run("all same key", func(t *testing.T) {
		seq := xiter.FromSlice([]int{1, 2, 3, 4, 5})
		groups := toSlice(collect.GroupBy(seq, func(v int) int {
			return 1
		}))
		if len(groups) != 1 || len(groups[0].Items) != 5 {
			t.Fatalf("GroupBy same: got %d groups with %d items", len(groups), len(groups[0].Items))
		}
	})
}

func TestSortedGroupBy(t *testing.T) {
	t.Run("basic", func(t *testing.T) {
		seq := xiter.FromSlice([]string{"apple", "avocado", "banana", "blueberry"})
		groups := toSlice(collect.SortedGroupBy(seq, func(s string) string {
			return string(s[0])
		}))

		if len(groups) != 2 {
			t.Fatalf("SortedGroupBy: got %d groups, want 2", len(groups))
		}
		if groups[0].Key != "a" || len(groups[0].Items) != 2 {
			t.Fatalf("SortedGroupBy[0]: key=%q items=%v", groups[0].Key, groups[0].Items)
		}
		if groups[1].Key != "b" || len(groups[1].Items) != 2 {
			t.Fatalf("SortedGroupBy[1]: key=%q items=%v", groups[1].Key, groups[1].Items)
		}
	})

	t.Run("empty", func(t *testing.T) {
		groups := toSlice(collect.SortedGroupBy(xiter.FromSlice([]string{}), func(s string) string {
			return s
		}))
		if len(groups) != 0 {
			t.Fatalf("SortedGroupBy empty: got %d, want 0", len(groups))
		}
	})

	t.Run("single element", func(t *testing.T) {
		groups := toSlice(collect.SortedGroupBy(xiter.FromSlice([]string{"hello"}), func(s string) string {
			return string(s[0])
		}))
		if len(groups) != 1 || groups[0].Key != "h" || len(groups[0].Items) != 1 {
			t.Fatalf("SortedGroupBy single: got %v", groups)
		}
	})

	t.Run("unsorted input splits keys", func(t *testing.T) {
		seq := xiter.FromSlice([]string{"apple", "banana", "avocado"})
		groups := toSlice(collect.SortedGroupBy(seq, func(s string) string {
			return string(s[0])
		}))
		if len(groups) != 3 {
			t.Fatalf("SortedGroupBy unsorted: got %d groups, want 3 (a, b, a)", len(groups))
		}
		// apple (a) -> group 0
		// banana (b) -> group 1
		// avocado (a) -> group 2 (different from group 0!)
		if groups[2].Key != "a" || groups[2].Items[0] != "avocado" {
			t.Fatalf("SortedGroupBy unsorted: last group should be (a, [avocado])")
		}
	})
}

func TestWindow(t *testing.T) {
	t.Run("basic", func(t *testing.T) {
		seq := xiter.FromSlice([]int{1, 2, 3, 4, 5, 6})
		windows := toSlice(collect.Window(seq, 3))
		if len(windows) != 4 {
			t.Fatalf("Window: got %d windows, want 4", len(windows))
		}
		expected := [][]int{{1, 2, 3}, {2, 3, 4}, {3, 4, 5}, {4, 5, 6}}
		for i, w := range windows {
			for j, v := range w {
				if v != expected[i][j] {
					t.Fatalf("Window[%d]: got %v, want %v", i, windows[i], expected[i])
				}
			}
		}
	})

	t.Run("window equals seq length", func(t *testing.T) {
		seq := xiter.FromSlice([]int{1, 2, 3})
		windows := toSlice(collect.Window(seq, 3))
		if len(windows) != 1 || windows[0][0] != 1 || windows[0][1] != 2 || windows[0][2] != 3 {
			t.Fatalf("Window full: got %v, want [[1 2 3]]", windows)
		}
	})

	t.Run("window larger than seq", func(t *testing.T) {
		seq := xiter.FromSlice([]int{1, 2})
		windows := toSlice(collect.Window(seq, 3))
		if len(windows) != 0 {
			t.Fatalf("Window larger: got %d windows, want 0", len(windows))
		}
	})

	t.Run("empty", func(t *testing.T) {
		windows := toSlice(collect.Window(xiter.FromSlice([]int{}), 3))
		if len(windows) != 0 {
			t.Fatalf("Window empty: got %d, want 0", len(windows))
		}
	})

	t.Run("window size 1", func(t *testing.T) {
		seq := xiter.FromSlice([]int{1, 2, 3})
		windows := toSlice(collect.Window(seq, 1))
		if len(windows) != 3 {
			t.Fatalf("Window size 1: got %d windows, want 3", len(windows))
		}
		for i, w := range windows {
			if len(w) != 1 || w[0] != i+1 {
				t.Fatalf("Window[%d]: got %v, want [%d]", i, w, i+1)
			}
		}
	})

	t.Run("early stop", func(t *testing.T) {
		seq := xiter.FromSlice([]int{1, 2, 3, 4, 5, 6, 7, 8})
		count := 0
		for range collect.Window(seq, 3) {
			count++
			if count >= 2 {
				break
			}
		}
		if count != 2 {
			t.Fatalf("Window early stop: count = %d, want 2", count)
		}
	})

	t.Run("zero or negative size", func(t *testing.T) {
		seq := xiter.FromSlice([]int{1, 2, 3})
		windows := toSlice(collect.Window(seq, 0))
		if len(windows) != 0 {
			t.Fatalf("Window zero: got %d, want 0", len(windows))
		}
		windows = toSlice(collect.Window(seq, -1))
		if len(windows) != 0 {
			t.Fatalf("Window neg: got %d, want 0", len(windows))
		}
	})
}
