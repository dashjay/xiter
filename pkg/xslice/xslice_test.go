package xslice_test

import (
	"strconv"
	"testing"

	"github.com/dashjay/xiter/pkg/internal/constraints"
	"github.com/dashjay/xiter/pkg/optional"
	"github.com/dashjay/xiter/pkg/xiter"
	"github.com/dashjay/xiter/pkg/xslice"
	"github.com/stretchr/testify/assert"
)

func _range(a, b int) []int {
	var res []int
	for i := a; i < b; i++ {
		res = append(res, i)
	}
	return res
}

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

func TestSlices(t *testing.T) {
	t.Run("test all", func(t *testing.T) {
		assert.True(t, xslice.All([]int{1, 2, 3}, func(x int) bool { return x > 0 }))
		assert.False(t, xslice.All([]int{-1, 1, 2, 3}, func(x int) bool { return x > 0 }))
		assert.True(t, xslice.All(_range(1, 9999), func(x int) bool { return x > 0 }))
		assert.False(t, xslice.All(_range(0, 9999), func(x int) bool { return x > 0 }))
	})

	t.Run("test any", func(t *testing.T) {
		assert.True(t, xslice.Any([]int{0, 1, 2, 3}, func(x int) bool { return x == 0 }))
		assert.False(t, xslice.Any([]int{0, 1, 2, 3}, func(x int) bool { return x == -1 }))
		assert.True(t, xslice.Any(_range(1, 9999), func(x int) bool { return x == 5000 }))
		assert.False(t, xslice.Any(_range(0, 9999), func(x int) bool { return x < 0 }))
	})

	t.Run("test avg", func(t *testing.T) {
		assert.Equal(t, avg(_range(1, 101)), xslice.Avg(_range(1, 101)))
		assert.Equal(t, float64(0), xslice.Avg([]int{}))
		assert.Equal(t, float64(0), xslice.Avg(_range(-50, 51)))

		assert.Equal(t, avg(_range(1, 101)), xslice.AvgN(_range(1, 101)...))
		assert.Equal(t, float64(0), xslice.AvgN([]int{}...))
		assert.Equal(t, float64(0), xslice.AvgN(_range(-50, 51)...))
	})

	t.Run("test avg by", func(t *testing.T) {
		assert.Equal(t, float64(2), xslice.AvgBy([]string{"1", "2", "3"}, func(x string) int {
			i, _ := strconv.Atoi(x)
			return i
		}))
		assert.Equal(t, float64(0), xslice.AvgBy([]string{"0"}, func(x string) int {
			i, _ := strconv.Atoi(x)
			return i
		}))
		assert.Equal(t, float64(0), xslice.AvgBy([]string{}, func(x string) int {
			i, _ := strconv.Atoi(x)
			return i
		}))
	})

	t.Run("test contains", func(t *testing.T) {
		// contains
		assert.True(t, xslice.Contains([]int{1, 2, 3}, 1))
		assert.False(t, xslice.Contains([]int{-1, 2, 3}, 1))

		// contains by
		assert.True(t, xslice.ContainsBy([]string{"1", "2", "3"}, func(x string) bool {
			i, _ := strconv.Atoi(x)
			return i == 1
		}))
		assert.False(t, xslice.ContainsBy([]string{"1", "2", "3"}, func(x string) bool {
			i, _ := strconv.Atoi(x)
			return i == -1
		}))

		// contains any
		assert.True(t, xslice.ContainsAny([]string{"1", "2", "3"}, []string{"1", "99", "1000"}))
		assert.False(t, xslice.ContainsAny([]string{"1", "2", "3"}, []string{"-1"}))
		assert.False(t, xslice.ContainsAny([]string{"1", "2", "3"}, []string{}))

		// contains all
		assert.True(t, xslice.ContainsAll([]string{"1", "2", "3"}, []string{"1", "2", "3"}))
		assert.False(t, xslice.ContainsAll([]string{"1", "2", "3"}, []string{"1", "99", "1000"}))
		assert.True(t, xslice.ContainsAll([]string{"1", "2", "3"}, []string{}))
	})

	t.Run("test count", func(t *testing.T) {
		assert.Equal(t, 3, xslice.Count([]int{1, 2, 3}))
		assert.Equal(t, 0, xslice.Count([]int{}))
		assert.Equal(t, 10000, xslice.Count(_range(0, 10000)))
	})

	t.Run("test find", func(t *testing.T) {
		assert.Equal(t, 1,
			optional.FromValue2(xslice.Find(_range(0, 10), func(x int) bool { return x == 1 })).Must())
		assert.False(t, optional.FromValue2(xslice.Find(_range(0, 10), func(x int) bool { return x == -1 })).Ok())

		assert.Equal(t, 1,
			xslice.FindO(_range(0, 10), func(x int) bool { return x == 1 }).Must())
		assert.False(t,
			xslice.FindO(_range(0, 10), func(x int) bool { return x == -1 }).Ok())
	})

	t.Run("test index", func(t *testing.T) {
		assert.Equal(t, -1, xslice.Index([]byte{}, 'x'))
		assert.Equal(t, 0, xslice.Index([]byte("aabb"), 'a'))
	})

	t.Run("test foreach", func(t *testing.T) {
		var res []int
		xslice.ForEach(_range(0, 10), func(i int) bool {
			if i == 5 {
				return false
			}
			res = append(res, i)
			return true
		})
		assert.Equal(t, _range(0, 5), res)

		var idxs []int
		var res2 []int

		xslice.ForEachIdx(_range(0, 10), func(idx int, v int) bool {
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
			optional.FromValue2(xslice.Head(_range(0, 10))).Must())
		assert.False(t,
			optional.FromValue2(xslice.Head(_range(0, 0))).Ok())

		assert.Equal(t, 0,
			xslice.HeadO(_range(0, 10)).Must())
		assert.False(t,
			xslice.HeadO(_range(0, 0)).Ok())
	})

	t.Run("test join", func(t *testing.T) {
		assert.Equal(t, "1.2.3", xslice.Join([]string{"1", "2", "3"}, "."))
		assert.Equal(t, "", xslice.Join([]string{}, "."))
	})

	t.Run("min max", func(t *testing.T) {
		assert.Equal(t, 1, xslice.Min([]int{1, 2, 3}).Must())
		assert.Equal(t, 1, xslice.MinN([]int{1, 2, 3}...).Must())
		assert.Equal(t, 3, xslice.Max([]int{1, 2, 3}).Must())
		assert.Equal(t, 3, xslice.MaxN([]int{1, 2, 3}...).Must())

		assert.False(t, xslice.Min([]int{}).Ok())
		assert.False(t, xslice.Max([]int{}).Ok())

		assert.Equal(t, 3,
			xslice.MinBy([]int{3, 2, 1} /*less = */, func(a, b int) bool { return a > b }).Must())
		assert.Equal(t, 1,
			xslice.MaxBy([]int{1, 2, 3} /*less = */, func(a, b int) bool { return a > b }).Must())
	})

	t.Run("clone", func(t *testing.T) {
		assert.Equal(t, []int{1, 2, 3}, xslice.Clone([]int{1, 2, 3}))
		assert.Len(t, xslice.Clone([]int{}), 0)
		assert.Len(t, xslice.Clone([]int(nil)), 0)
		assert.Equal(t, []string{"1", "2", "3"}, xslice.CloneBy([]int{1, 2, 3}, strconv.Itoa))
		assert.Len(t, xslice.CloneBy([]int{}, strconv.Itoa), 0)
		assert.Len(t, xslice.CloneBy([]int(nil), strconv.Itoa), 0)
	})

	t.Run("concat", func(t *testing.T) {
		assert.Equal(t, []int{1, 2, 3, 4, 5}, xslice.Concat([]int{1, 2, 3}, []int{4, 5}))
		assert.Equal(t, []int{1, 2, 3}, xslice.Concat([]int{1, 2, 3}))
		assert.Len(t, xslice.Concat([]int{}, []int{}, []int{}), 0)
	})

	t.Run("subset", func(t *testing.T) {
		assert.Len(t, xslice.Subset([]int{1, 2, 3}, 0, -1), 0)
		assert.Len(t, xslice.Subset([]int{1, 2, 3}, 3, 0), 0)
		assert.Len(t, xslice.Subset([]int{1, 2, 3}, -3, 0), 0)
		assert.Len(t, xslice.Subset([]int{1, 2, 3}, 0, 0), 0)
		assert.Equal(t, []int{1}, xslice.Subset([]int{1, 2, 3}, 0, 1))
		assert.Equal(t, []int{1, 2}, xslice.Subset([]int{1, 2, 3}, 0, 2))
		assert.Equal(t, []int{1, 2, 3}, xslice.Subset([]int{1, 2, 3}, 0, 3))

		for i := 0; i < 100; i++ {
			assert.Equal(t, _range(i, i+10), xslice.Subset(_range(0, 200), i, 10))
			assert.Subset(t, _range(0, 200), xslice.Subset(_range(0, 200), i, 10))
		}

		assert.Equal(t, []int{3}, xslice.SubsetInPlace([]int{1, 2, 3}, -1, 1))
		assert.Equal(t, []int{3}, xslice.SubsetInPlace([]int{1, 2, 3}, -1, 2))
		assert.Equal(t, []int{3}, xslice.SubsetInPlace([]int{1, 2, 3}, -1, 3))
		assert.Len(t, xslice.SubsetInPlace([]int{1, 2, 3}, -999, 3), 0)
		assert.Len(t, xslice.SubsetInPlace([]int{1, 2, 3}, 999, 3), 0)

		for i := 1; i < 100; i++ {
			assert.Equal(t, _range(200-i, xslice.MinN(200-i+10, 200).Must()), xslice.Subset(_range(0, 200), -i, 10))
		}

		original := []int{1, 2, 3}
		clonedSubset := xslice.Subset(original, 0, 2)
		clonedSubset[0] = 100
		assert.Equal(t, []int{1, 2, 3}, original)
		assert.Equal(t, []int{100, 2}, clonedSubset)
	})

	t.Run("test replace", func(t *testing.T) {
		assert.Equal(t, append([]int{10}, _range(1, 10)...), xslice.ReplaceAll(_range(0, 10), 0, 10))
		assert.Equal(t, append([]int{10}, _range(1, 10)...), xslice.Replace(_range(0, 10), 0, 10, 1))
		assert.Equal(t, append([]int{10}, _range(1, 10)...), xslice.Replace(_range(0, 10), 0, 10, 5))

		// replace nothing
		assert.Equal(t, _range(0, 10), xiter.ToSlice(xiter.Replace(xiter.FromSlice(_range(0, 10)), 0, 100, 0)))
	})

	t.Run("test reverse", func(t *testing.T) {
		assert.Equal(t, []int{3, 2, 1}, xslice.ReverseClone([]int{1, 2, 3}))
		assert.Equal(t, []int{5, 4, 3, 2, 1}, xslice.ReverseClone([]int{1, 2, 3, 4, 5}))

		arr := []int{1, 2, 3}
		xslice.Reverse(arr)
		assert.Equal(t, []int{3, 2, 1}, arr)
	})

	t.Run("test repeat and repeat_by", func(t *testing.T) {
		assert.Equal(t, []int{1, 1, 1}, xslice.Repeat([]int{1}, 3))
		assert.Equal(t, []int{1, 2, 3, 1, 2, 3, 1, 2, 3}, xslice.Repeat([]int{1, 2, 3}, 3))

		assert.Equal(t, []int{1, 2, 3}, xslice.RepeatBy(3, func(idx int) int {
			return idx + 1
		}))

		assert.Equal(t, []string{"0", "1", "2"}, xslice.RepeatBy(3, func(idx int) string {
			return strconv.Itoa(idx)
		}))
	})

	t.Run("test shuffle", func(t *testing.T) {
		assert.Len(t, xslice.Shuffle([]int{1, 2, 3}), 3)
		arr := _range(1, 100)
		xslice.ShuffleInPlace(arr)
		assert.Len(t, arr, 99)
	})

	t.Run("chunk and chunk inplace", func(t *testing.T) {
		assert.Equal(t, [][]int{{1, 2, 3}, {4, 5, 6}}, xslice.Chunk([]int{1, 2, 3, 4, 5, 6}, 3))
		assert.Equal(t, [][]int{{1, 2, 3}, {4, 5, 6}}, xslice.ChunkInPlace([]int{1, 2, 3, 4, 5, 6}, 3))

		assert.Equal(t, [][]int{{1, 2, 3}, {4, 5, 6}, {7, 8, 9}}, xslice.Chunk([]int{1, 2, 3, 4, 5, 6, 7, 8, 9}, 3))
		assert.Equal(t, [][]int{{1, 2, 3}, {4, 5, 6}, {7, 8, 9}}, xslice.ChunkInPlace([]int{1, 2, 3, 4, 5, 6, 7, 8, 9}, 3))

		assert.Equal(t, [][]int{{1, 2, 3}, {4, 5, 6}, {7, 8, 9}, {10}}, xslice.Chunk([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}, 3))
		assert.Equal(t, [][]int{{1, 2, 3}, {4, 5, 6}, {7, 8, 9}, {10}}, xslice.ChunkInPlace([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}, 3))

		assert.Len(t, xslice.ChunkInPlace([]int{}, 1), 0)
		assert.Len(t, xslice.Chunk([]int{}, 1), 0)
	})

	t.Run("index", func(t *testing.T) {
		assert.Equal(t, 50, xslice.Index(_range(0, 101), 50))
		assert.Equal(t, -1, xslice.Index(_range(0, 101), 6666))
	})

	t.Run("sum", func(t *testing.T) {
		assert.Equal(t, 5050, xslice.Sum(_range(0, 101)))
		assert.Equal(t, 5050, xslice.SumN(_range(0, 101)...))

		assert.Equal(t, 55, xslice.SumBy([]string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "10"}, func(t string) int {
			r, _ := strconv.Atoi(t)
			return r
		}))
	})
	t.Run("uniq", func(t *testing.T) {
		assert.Equal(t, []int{1, 2, 3, 4}, xslice.Uniq([]int{1, 2, 3, 2, 4}))
	})
	t.Run("group by", func(t *testing.T) {
		groupedBy := xslice.GroupBy([]int{0, 1, 2, 3, 4}, func(i int) string {
			if i%2 == 0 {
				return "even"
			}
			return "odd"
		})
		assert.Contains(t, groupedBy, "even")
		assert.Contains(t, groupedBy, "odd")
		assert.Contains(t, groupedBy["even"], 0)
		assert.Contains(t, groupedBy["even"], 2)
		assert.Contains(t, groupedBy["even"], 4)
		assert.Contains(t, groupedBy["odd"], 1)
		assert.Contains(t, groupedBy["odd"], 3)
	})
	t.Run("group by map", func(t *testing.T) {
		groupedBy := xslice.GroupByMap([]int{0, 1, 2, 3, 4}, func(i int) (string, int) {
			if i%2 == 0 {
				return "even", i * 2
			}
			return "odd", i*2 + 1
		})
		assert.Contains(t, groupedBy, "even")
		assert.Contains(t, groupedBy, "odd")
		assert.Contains(t, groupedBy["even"], 0)
		assert.Contains(t, groupedBy["even"], 4)
		assert.Contains(t, groupedBy["even"], 8)
		assert.Contains(t, groupedBy["odd"], 3)
		assert.Contains(t, groupedBy["odd"], 7)
	})
}
