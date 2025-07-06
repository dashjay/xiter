package xiter_test

import (
	"github.com/dashjay/xiter/pkg/internal/constraints"
	"github.com/dashjay/xiter/pkg/optional"
	"github.com/dashjay/xiter/pkg/xiter"
	"github.com/stretchr/testify/assert"
	"strconv"
	"testing"
)

func avg[T constraints.Number](in []T) float64 {
	if len(in) == 0 {
		return 0
	}
	var sum T
	for i := 0; i < len(in); i++ {
		sum += in[i]
	}
	return float64(sum) / float64(len(in))
}

func TestPull(t *testing.T) {
	seq := xiter.FromSlice(_range(0, 100))
	next, stop := xiter.Pull(seq)
	defer stop()
	for i := 0; i < 100; i++ {
		v, ok := next()
		assert.True(t, ok)
		assert.Equal(t, i, v)
	}

	_, ok := next()
	assert.False(t, ok)
}

func TestPull2(t *testing.T) {
	seq := xiter.FromSliceIdx(_range(0, 100))
	next, stop := xiter.Pull2(seq)
	defer stop()
	for i := 0; i < 100; i++ {
		k, v, ok := next()
		assert.True(t, ok)
		assert.Equal(t, i, k)
		assert.Equal(t, i, v)
	}

	_, _, ok := next()
	assert.False(t, ok)
}

func TestXIter(t *testing.T) {
	t.Run("to slice", func(t *testing.T) {
		assert.Equal(t, _range(0, 1000), xiter.ToSlice(xiter.FromSlice(_range(0, 1000))))
	})

	t.Run("to slice 2", func(t *testing.T) {
		m := map[int]string{1: "1", 2: "2", 3: "3"}

		assert.Contains(t, xiter.ToSliceSeq2Key(xiter.FromMapKeyAndValues(m)), 1)
		assert.Contains(t, xiter.ToSliceSeq2Key(xiter.FromMapKeyAndValues(m)), 2)
		assert.Contains(t, xiter.ToSliceSeq2Key(xiter.FromMapKeyAndValues(m)), 3)

		assert.Contains(t, xiter.ToSliceSeq2Value(xiter.FromMapKeyAndValues(m)), "1")
		assert.Contains(t, xiter.ToSliceSeq2Value(xiter.FromMapKeyAndValues(m)), "2")
		assert.Contains(t, xiter.ToSliceSeq2Value(xiter.FromMapKeyAndValues(m)), "3")
	})

	t.Run("to map", func(t *testing.T) {
		m := map[int]string{1: "1", 2: "2", 3: "3"}
		newMap := xiter.ToMap(xiter.FromMapKeyAndValues(m))

		for k, v := range newMap {
			assert.Equal(t, v, m[k])
		}
	})

	t.Run("from map", func(t *testing.T) {
		m := map[int]string{1: "1", 2: "2", 3: "3"}
		keys := xiter.ToSlice(xiter.FromMapKeys(m))
		values := xiter.ToSlice(xiter.FromMapValues(m))

		assert.Contains(t, keys, 1)
		assert.Contains(t, keys, 2)
		assert.Contains(t, keys, 3)
		assert.Contains(t, values, "1")
		assert.Contains(t, values, "2")
		assert.Contains(t, values, "3")

		kvPair := xiter.FromMapKeyAndValues(m)
		newMap := xiter.ToMap(kvPair)

		for k, v := range newMap {
			assert.Equal(t, v, m[k])
		}
	})

	t.Run("test all", func(t *testing.T) {
		assert.True(t, xiter.AllFromSeq(xiter.FromSlice([]int{1, 2, 3}), func(x int) bool { return x > 0 }))
		assert.False(t, xiter.AllFromSeq(xiter.FromSlice([]int{-1, 1, 2, 3}), func(x int) bool { return x > 0 }))
		assert.True(t, xiter.AllFromSeq(xiter.FromSlice(_range(1, 9999)), func(x int) bool { return x > 0 }))
		assert.False(t, xiter.AllFromSeq(xiter.FromSlice(_range(0, 9999)), func(x int) bool { return x > 0 }))
	})

	t.Run("test any", func(t *testing.T) {
		assert.True(t, xiter.AnyFromSeq(xiter.FromSlice([]int{0, 1, 2, 3}), func(x int) bool { return x == 0 }))
		assert.False(t, xiter.AnyFromSeq(xiter.FromSlice([]int{0, 1, 2, 3}), func(x int) bool { return x == -1 }))
		assert.True(t, xiter.AnyFromSeq(xiter.FromSlice(_range(1, 9999)), func(x int) bool { return x == 5000 }))
		assert.False(t, xiter.AnyFromSeq(xiter.FromSlice(_range(0, 9999)), func(x int) bool { return x < 0 }))
	})

	t.Run("test avg & avg by", func(t *testing.T) {
		assert.Equal(t, avg(_range(1, 101)), xiter.AvgFromSeq(xiter.FromSlice(_range(1, 101))))
		assert.Equal(t, float64(0), xiter.AvgFromSeq(xiter.FromSlice([]int{})))
		assert.Equal(t, float64(0), xiter.AvgFromSeq(xiter.FromSlice(_range(-50, 51))))

		assert.Equal(t, float64(2), xiter.AvgByFromSeq(xiter.FromSlice([]string{"1", "2", "3"}), func(x string) int {
			i, _ := strconv.Atoi(x)
			return i
		}))
		assert.Equal(t, float64(0), xiter.AvgByFromSeq(xiter.FromSlice([]string{"0"}), func(x string) int {
			i, _ := strconv.Atoi(x)
			return i
		}))
		assert.Equal(t, float64(0), xiter.AvgByFromSeq(xiter.FromSlice([]string{}), func(x string) int {
			i, _ := strconv.Atoi(x)
			return i
		}))
	})

	t.Run("test contains", func(t *testing.T) {
		// contains
		assert.True(t, xiter.Contains(xiter.FromSlice([]int{1, 2, 3}), 1))
		assert.False(t, xiter.Contains(xiter.FromSlice([]int{-1, 2, 3}), 1))

		// contains by
		assert.True(t, xiter.ContainsBy(xiter.FromSlice([]string{"1", "2", "3"}), func(x string) bool {
			i, _ := strconv.Atoi(x)
			return i == 1
		}))
		assert.False(t, xiter.ContainsBy(xiter.FromSlice([]string{"1", "2", "3"}), func(x string) bool {
			i, _ := strconv.Atoi(x)
			return i == -1
		}))

		// contains any
		assert.True(t, xiter.ContainsAny(xiter.FromSlice([]string{"1", "2", "3"}), []string{"1", "99", "1000"}))
		assert.False(t, xiter.ContainsAny(xiter.FromSlice([]string{"1", "2", "3"}), []string{"-1"}))
		assert.False(t, xiter.ContainsAny(xiter.FromSlice([]string{"1", "2", "3"}), []string{}))

		// contains all
		assert.True(t, xiter.ContainsAll(xiter.FromSlice([]string{"1", "2", "3"}), []string{"1", "2", "3"}))
		assert.False(t, xiter.ContainsAll(xiter.FromSlice([]string{"1", "2", "3"}), []string{"1", "99", "1000"}))
		assert.True(t, xiter.ContainsAll(xiter.FromSlice([]string{"1", "2", "3"}), []string{}))
	})

	t.Run("test count", func(t *testing.T) {
		assert.Equal(t, len(_range(0, 10)), xiter.Count(xiter.FromSlice(_range(0, 10))))
	})

	t.Run("test find", func(t *testing.T) {
		assert.Equal(t, 1,
			optional.FromValue2(xiter.Find(xiter.FromSlice(_range(0, 10)), func(x int) bool { return x == 1 })).Must())
		assert.False(t, optional.FromValue2(xiter.Find(xiter.FromSlice(_range(0, 10)), func(x int) bool { return x == -1 })).Ok())

		assert.Equal(t, 1,
			xiter.FindO(xiter.FromSlice(_range(0, 10)), func(x int) bool { return x == 1 }).Must())
		assert.False(t,
			xiter.FindO(xiter.FromSlice(_range(0, 10)), func(x int) bool { return x == -1 }).Ok())
	})

	t.Run("test foreach", func(t *testing.T) {
		var res []int
		xiter.ForEach(xiter.FromSlice(_range(0, 10)), func(i int) bool {
			if i == 5 {
				return false
			}
			res = append(res, i)
			return true
		})
		assert.Equal(t, _range(0, 5), res)

		var idxs []int
		var res2 []int

		xiter.ForEachIdx(xiter.FromSlice(_range(0, 10)), func(idx int, v int) bool {
			if idx == 5 {
				return false
			}
			idxs = append(idxs, idx)
			res2 = append(res2, v)
			return true
		})
		assert.Equal(t, _range(0, 5), idxs)
		assert.Equal(t, _range(0, 5), res2)
	})

	t.Run("test head", func(t *testing.T) {
		assert.Equal(t, 0,
			optional.FromValue2(xiter.Head(xiter.FromSlice(_range(0, 10)))).Must())
		assert.False(t,
			optional.FromValue2(xiter.Head(xiter.FromSlice(_range(0, 0)))).Ok())

		assert.Equal(t, 0,
			xiter.HeadO(xiter.FromSlice(_range(0, 10))).Must())
		assert.False(t,
			xiter.HeadO(xiter.FromSlice(_range(0, 0))).Ok())
	})

	t.Run("test join", func(t *testing.T) {
		assert.Equal(t, "1.2.3", xiter.Join(xiter.FromSlice([]string{"1", "2", "3"}), "."))
		assert.Equal(t, "", xiter.Join(xiter.FromSlice([]string{}), "."))
	})

	t.Run("min max", func(t *testing.T) {
		assert.Equal(t, 1, xiter.Min(xiter.FromSlice([]int{3, 2, 1})).Must())
		assert.Equal(t, 3, xiter.Max(xiter.FromSlice([]int{1, 2, 3})).Must())

		assert.False(t, xiter.Min(xiter.FromSlice([]int{})).Ok())
		assert.False(t, xiter.Max(xiter.FromSlice([]int{})).Ok())

		assert.Equal(t, 3,
			xiter.MinBy(xiter.FromSlice([]int{1, 3, 2}) /*less = */, func(a, b int) bool { return a > b }).Must())
		assert.Equal(t, 1,
			xiter.MaxBy(xiter.FromSlice([]int{3, 1, 2}) /*less = */, func(a, b int) bool { return a > b }).Must())

		assert.False(t, xiter.MinBy(xiter.FromSlice([]int{}) /*less = */, func(a, b int) bool { return a > b }).Ok())
		assert.False(t, xiter.MaxBy(xiter.FromSlice([]int{}) /*less = */, func(a, b int) bool { return a > b }).Ok())
	})

	t.Run("to slice", func(t *testing.T) {
		assert.Equal(t, _range(0, 10), xiter.ToSlice(xiter.FromSlice(_range(0, 10))))
	})

	t.Run("concat and filter", func(t *testing.T) {
		assert.Equal(t, _range(0, 10), xiter.ToSlice(xiter.Concat(xiter.FromSlice(_range(0, 5)), xiter.FromSlice(_range(5, 10)))))
		assert.Equal(t, []int{0, 1, 2, 3, 4 /* 5 is filtered */, 6, 7, 8, 9},
			xiter.ToSlice(xiter.Concat(
				xiter.FromSlice(_range(0, 5)),
				xiter.Filter(func(v int) bool { return v != 5 }, xiter.FromSlice(_range(5, 10))),
			)))
	})

	t.Run("test pullout", func(t *testing.T) {
		for i := 0; i < 1000; i++ {
			if i < 100 {
				assert.Len(t, xiter.PullOut(xiter.FromSlice(_range(0, 100)), i), i)
			} else {
				assert.Len(t, xiter.PullOut(xiter.FromSlice(_range(0, 100)), i), 100)
			}
		}

		assert.Len(t, xiter.PullOut(xiter.FromSlice(_range(0, 100)), -1), 100)
	})

	t.Run("test at", func(t *testing.T) {
		for i := 0; i < 1000; i++ {
			if i < 100 {
				assert.Equal(t, i, xiter.At(xiter.FromSlice(_range(0, 100)), i).Must())
			} else {
				assert.False(t, xiter.At(xiter.FromSlice(_range(0, 100)), i).Ok())
			}
		}

		cc := xiter.Concat(
			xiter.FromSlice(_range(0, 100)),
			xiter.FromSlice(_range(100, 200)),
		)
		assert.Equal(t, 150, xiter.At(cc, 150).Must())
		cc = xiter.Filter(func(v int) bool { return v%5 == 0 }, xiter.FromSlice(_range(0, 100)))
		assert.Equal(t, 25, xiter.At(cc, 5).Must())
	})

	t.Run("skip and limit", func(t *testing.T) {
		// skip
		assert.Equal(t, _range(10, 30), xiter.ToSlice(xiter.Skip(xiter.FromSlice(_range(0, 30)), 10)))
		assert.Equal(t, _range(0, 30), xiter.ToSlice(xiter.Skip(xiter.FromSlice(_range(0, 30)), 0)))
		assert.Equal(t, _range(10, 20), xiter.ToSlice(xiter.Limit(xiter.Skip(xiter.FromSlice(_range(0, 30)), 10), 10)))

		// limit
		assert.Equal(t, _range(0, 10), xiter.ToSlice(xiter.Limit(xiter.FromSlice(_range(0, 30)), 10)))
		assert.Equal(t, _range(0, 10)[0:1], xiter.ToSlice(xiter.Limit(xiter.Limit(xiter.FromSlice(_range(0, 30)), 10), 1)))
		assert.Equal(t, _range(0, 10), xiter.ToSlice(xiter.Limit(xiter.FromSlice(_range(0, 10)), 10)))
		assert.Equal(t, _range(0, 0), xiter.ToSlice(xiter.Limit(xiter.FromSlice(_range(0, 0)), 10)))
		assert.Equal(t, _range(0, 0), xiter.ToSlice(xiter.Limit(xiter.FromSlice(_range(0, 10)), 0)))
	})

	t.Run("test repeat", func(t *testing.T) {
		assert.Equal(t, []int{1, 2, 3, 1, 2, 3, 1, 2, 3}, xiter.ToSlice(xiter.Repeat(xiter.FromSlice([]int{1, 2, 3}), 3)))
	})

	t.Run("test shuffle", func(t *testing.T) {
		assert.Len(t, xiter.ToSlice(xiter.Limit(xiter.FromSliceShuffle(_range(0, 10)), 5)), 5)
		assert.Len(t, xiter.ToSlice(xiter.FromSliceShuffle(_range(0, 10))), 10)
	})
}
