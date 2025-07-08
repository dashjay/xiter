package xiter_test

import (
	"bytes"
	"fmt"
	"math"
	"strconv"
	"testing"

	"github.com/dashjay/xiter/pkg/internal/constraints"
	"github.com/dashjay/xiter/pkg/optional"
	"github.com/dashjay/xiter/pkg/xiter"
	"github.com/stretchr/testify/assert"
)

func testLimit[T any](t *testing.T, seq xiter.Seq[T], n int) {
	assert.Len(t, xiter.ToSlice(xiter.Limit(seq, 1)), n)
}

func testLimit2[K, V any](t *testing.T, seq xiter.Seq2[K, V], n int) {
	assert.Len(t, xiter.ToSliceSeq2Key(xiter.Limit2(seq, 1)), n)
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

func _range(a, b int) []int {
	var res []int
	for i := a; i < b; i++ {
		res = append(res, i)
	}
	return res
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

		assert.Len(t, xiter.ToSlice(xiter.Limit(xiter.FromMapKeys(m), 1)), 1)
		assert.Len(t, xiter.ToSlice(xiter.Limit(xiter.FromMapValues(m), 1)), 1)

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

func TestXIter61898(t *testing.T) {
	const (
		range1Start = 0
		range1End   = 500
		range2Start = 500
		range2End   = 1000
	)
	range1 := _range(range1Start, range1End)
	range2 := _range(range2Start, range2End)

	t.Run("concat", func(t *testing.T) {
		seq := xiter.Concat(xiter.FromSlice(range1), xiter.FromSlice(range2))
		testLimit(t, seq, 1)
		next, stop := xiter.Pull(seq)
		defer stop()

		for i := 0; i < 1000; i++ {
			k, ok := next()
			assert.Equal(t, i, k)
			assert.True(t, ok)
		}

		_, ok := next()
		assert.False(t, ok)
	})

	t.Run("concat2", func(t *testing.T) {
		mapA := make(map[int]string)
		mapB := make(map[int]string)
		for i := 0; i < 1000; i++ {
			mapA[i] = fmt.Sprintf("%d", i)
		}
		for i := 1000; i < 2000; i++ {
			mapB[i] = fmt.Sprintf("%d", i)
		}

		xiterMapa := xiter.FromMapKeyAndValues(mapA)
		xiterMapb := xiter.FromMapKeyAndValues(mapB)
		xiterMap := xiter.Concat2(xiterMapa, xiterMapb)
		testLimit2(t, xiterMap, 1)

		mapAB := xiter.ToMap(xiterMap)

		for k, v := range mapA {
			assert.Contains(t, mapAB, k)
			assert.True(t, v == mapAB[k])
		}
		for k, v := range mapB {
			assert.Contains(t, mapAB, k)
			assert.True(t, v == mapAB[k])
		}
	})

	t.Run("equal", func(t *testing.T) {
		seq1 := xiter.FromSlice(range1)
		assert.True(t, xiter.Equal(seq1, xiter.FromSlice(_range(range1Start, range1End))))
		// not equal
		assert.False(t, xiter.Equal(seq1, xiter.Limit(xiter.FromSlice(_range(range1Start, range1End)), 1)))
	})

	t.Run("equal2", func(t *testing.T) {
		seq1 := xiter.FromSliceIdx(range1)
		assert.True(t, xiter.Equal2(seq1, xiter.FromSliceIdx(_range(range1Start, range1End))))
		// not equal
		assert.False(t, xiter.Equal2(seq1, xiter.Limit2(xiter.FromSliceIdx(_range(range1Start, range1End)), 2)))
	})

	t.Run("equal_func", func(t *testing.T) {
		seq1 := xiter.FromSlice(range1)
		assert.True(t, xiter.EqualFunc(seq1, xiter.FromSlice(_range(range1Start, range1End)), func(a int, b int) bool {
			return a == b
		}))
		// not equal
		assert.False(t, xiter.EqualFunc(xiter.Limit(seq1, 1), xiter.FromSlice(_range(range1Start, range1End)), func(a int, b int) bool {
			return a == b
		}))
	})

	t.Run("equal_func2", func(t *testing.T) {
		seq1 := xiter.FromSliceIdx(range1)
		assert.True(t, xiter.EqualFunc2(seq1, xiter.FromSliceIdx(_range(range1Start, range1End)), func(k1, k2, v1, v2 int) bool {
			return k1 == k2 && v1 == v2
		}))

		// not equal
		assert.False(t, xiter.EqualFunc2(xiter.Limit2(seq1, 2), xiter.FromSliceIdx(_range(range1Start, range1End)), func(k1, k2, v1, v2 int) bool {
			return k1 == k2 && v1 == v2
		}))
	})

	t.Run("filter", func(t *testing.T) {
		seq1 := xiter.Filter(func(v int) bool {
			return v%2 == 0
		}, xiter.FromSlice(range1))

		var expectedSeq = xiter.Seq[int](func(yield func(v int) bool) {
			for i := range1Start; i < range1End; i++ {
				if i%2 == 0 {
					if !yield(i) {
						return
					}
				}
			}
		})
		assert.True(t, xiter.Equal(seq1, expectedSeq))
		assert.False(t, xiter.Equal(xiter.Limit(seq1, 1), expectedSeq))
	})

	t.Run("filter2", func(t *testing.T) {
		seq1 := xiter.Filter2(func(k, v int) bool {
			return v%2 == 0
		}, xiter.FromSliceIdx(range1))

		var expectedSeq = xiter.Seq2[int, int](func(yield func(k, v int) bool) {
			for i := range1Start; i < range1End; i++ {
				if i%2 == 0 {
					if !yield(i, i) {
						return
					}
				}
			}
		})
		assert.True(t, xiter.Equal2(seq1, expectedSeq))
		assert.False(t, xiter.Equal2(xiter.Limit2(seq1, 1), expectedSeq))
	})

	t.Run("limit", func(t *testing.T) {
		t.Run("limit zero", func(t *testing.T) {
			seq1 := xiter.FromSlice(range1)
			seq1Limited := xiter.Limit(seq1, 0)
			assert.Len(t, xiter.ToSlice(seq1Limited), 0)
		})

		t.Run("limit one", func(t *testing.T) {
			seq1 := xiter.FromSlice(range1)
			seq1Limited := xiter.Limit(seq1, 1)
			assert.Len(t, xiter.ToSlice(seq1Limited), 1)
		})

		t.Run("limit large", func(t *testing.T) {
			seq1 := xiter.FromSlice(range1)
			seq1Limited := xiter.Limit(seq1, math.MaxInt64)
			assert.Len(t, xiter.ToSlice(seq1Limited), range1End-range1Start)
		})

		t.Run("limit limit", func(t *testing.T) {
			seq1 := xiter.FromSlice(range1)
			seq1Limited := xiter.Limit(seq1, math.MaxInt64)
			seq1Limited = xiter.Limit(seq1Limited, 1)
			assert.Len(t, xiter.ToSlice(seq1Limited), 1)
		})
	})

	t.Run("limit2", func(t *testing.T) {
		t.Run("limit2 zero", func(t *testing.T) {
			seq1 := xiter.FromSliceIdx(range1)
			seq1Limit0 := xiter.Limit2(seq1, 0)
			assert.Len(t, xiter.ToSliceSeq2Key(seq1Limit0), 0)
		})

		t.Run("limit2 one", func(t *testing.T) {
			seq1 := xiter.FromSliceIdx(range1)
			seq1Limit0 := xiter.Limit2(seq1, 1)
			assert.Len(t, xiter.ToSliceSeq2Key(seq1Limit0), 1)
		})

		t.Run("limit2 large", func(t *testing.T) {
			seq1 := xiter.FromSliceIdx(range1)
			seq1Limit0 := xiter.Limit2(seq1, math.MaxInt64)
			assert.Len(t, xiter.ToSliceSeq2Key(seq1Limit0), range1End-range1Start)
			assert.Equal(t, 0, xiter.ToSliceSeq2Key(seq1Limit0)[0])
			assert.Equal(t, range1End-1, xiter.ToSliceSeq2Key(seq1Limit0)[range1End-1])
		})

		t.Run("limit2 limit2", func(t *testing.T) {
			seq1 := xiter.FromSliceIdx(range1)
			seq1Limited := xiter.Limit2(seq1, math.MaxInt64)
			seq1Limited = xiter.Limit2(seq1Limited, 1)
			assert.Len(t, xiter.ToSliceSeq2Key(seq1Limited), 1)
		})
	})

	t.Run("map", func(t *testing.T) {
		strSeq := xiter.Map(func(in int) string {
			return strconv.Itoa(in)
		}, xiter.FromSlice(_range(range1Start, range1End)))

		strArr := xiter.ToSlice(strSeq)
		for i := range1Start; i < range1End; i++ {
			assert.Equal(t, fmt.Sprintf("%d", i), strArr[i])
		}
	})

	t.Run("map2", func(t *testing.T) {
		strSeq := xiter.Map2(func(keyIn int, valueIn int) (string, string) {
			return strconv.Itoa(keyIn), strconv.Itoa(valueIn)
		}, xiter.FromSliceIdx(_range(range1Start, range1End)))
		strArr := xiter.ToSliceSeq2Key(strSeq)
		strArr2 := xiter.ToSliceSeq2Value(strSeq)
		for i := range1Start; i < range1End; i++ {
			assert.Equal(t, fmt.Sprintf("%d", i), strArr[i])
			assert.Equal(t, fmt.Sprintf("%d", i), strArr2[i])
		}

		for i := 0; i < 10; i++ {
			assert.Len(t, xiter.ToSliceSeq2Key(xiter.Limit2(strSeq, i)), i)
		}
	})

	t.Run("merge", func(t *testing.T) {
		mergedSeq := xiter.Merge(xiter.FromSlice(range1), xiter.FromSlice(range2))
		assert.Equal(t, xiter.ToSlice(mergedSeq), xiter.ToSlice(xiter.Concat(xiter.FromSlice(range1), xiter.FromSlice(range2))))
		var odds []int
		var evens []int
		for i := 0; i < 1000; i++ {
			if i%2 == 0 {
				evens = append(evens, i)
			} else {
				odds = append(odds, i)
			}
		}
		mergedSeq = xiter.Merge(xiter.FromSlice(odds), xiter.FromSlice(evens))
		arr := xiter.ToSlice(mergedSeq)
		assert.Equal(t, _range(0, 1000), arr)

		for i := 0; i < 100; i++ {
			assert.Len(t, xiter.ToSlice(xiter.Limit(mergedSeq, i)), i)
		}

		shortSeq := xiter.FromSlice([]int{})
		longSeq := xiter.FromSlice([]int{1, 2, 3})
		shortLongSeq := xiter.Merge(shortSeq, longSeq)
		longShortSeq := xiter.Merge(longSeq, shortSeq)
		assert.True(t, xiter.Equal(shortLongSeq, longSeq))
		assert.True(t, xiter.Equal(longShortSeq, longSeq))

		assert.Len(t, xiter.ToSlice(xiter.Limit(shortLongSeq, 1)), 1)
		assert.Len(t, xiter.ToSlice(xiter.Limit(shortLongSeq, 2)), 2)
		assert.Len(t, xiter.ToSlice(xiter.Limit(shortLongSeq, 3)), 3)
		assert.Len(t, xiter.ToSlice(xiter.Limit(longShortSeq, 1)), 1)
		assert.Len(t, xiter.ToSlice(xiter.Limit(longShortSeq, 2)), 2)
		assert.Len(t, xiter.ToSlice(xiter.Limit(longShortSeq, 3)), 3)

	})

	t.Run("merge2", func(t *testing.T) {
		var seq1 = xiter.Seq2[int, string](func(yield func(int, string) bool) {
			for i := 0; i < 10; i++ {
				if !yield(i, fmt.Sprintf("%d", i)) {
					return
				}
			}
			for i := 20; i < 30; i++ {
				if !yield(i, fmt.Sprintf("%d", i)) {
					return
				}
			}
			for i := 40; i < 50; i++ {
				if !yield(i, fmt.Sprintf("%d", i)) {
					return
				}
			}
		})
		var seq2 = xiter.Seq2[int, string](func(yield func(int, string) bool) {
			for i := 10; i < 20; i++ {
				if !yield(i, fmt.Sprintf("%d", i)) {
					return
				}
			}
			for i := 30; i < 40; i++ {
				if !yield(i, fmt.Sprintf("%d", i)) {
					return
				}
			}
		})
		mergedSeq := xiter.Merge2(seq1, seq2)
		var expectedSeq = xiter.Seq2[int, string](func(yield func(int, string) bool) {
			for i := 0; i < 50; i++ {
				if !yield(i, fmt.Sprintf("%d", i)) {
					break
				}
			}
		})

		assert.True(t, xiter.Equal2(mergedSeq, expectedSeq))
		for i := 0; i < 50; i++ {
			assert.Len(t, xiter.ToSliceSeq2Key(xiter.Limit2(mergedSeq, i)), i)
		}
		mergedSeq = xiter.Merge2(seq2, seq1)
		assert.True(t, xiter.Equal2(mergedSeq, expectedSeq))
		for i := 0; i < 50; i++ {
			assert.Len(t, xiter.ToSliceSeq2Key(xiter.Limit2(mergedSeq, i)), i)
		}
	})

	t.Run("reduce", func(t *testing.T) {
		seq := xiter.FromSlice(range1) // 0+1+2+...+499
		sum := xiter.Reduce[int, int](func(sum int, v int) int {
			return sum + v
		}, 0, seq)

		expectedSum := 0
		for i := range1Start; i < range1End; i++ {
			expectedSum += i
		}
		assert.Equal(t, expectedSum, sum)
	})

	t.Run("reduce2", func(t *testing.T) {
		sum := xiter.Reduce2(func(sum int, k int, v int) int {
			return sum + k
		}, 0, xiter.FromSliceIdx(range1))

		expectedSum := 0
		for i := range1Start; i < range1End; i++ {
			expectedSum += i
		}
		assert.Equal(t, expectedSum, sum)
	})

	t.Run("zip", func(t *testing.T) {
		t.Run("case1-base", func(t *testing.T) {
			zipped := xiter.Zip(xiter.FromSlice(range1), xiter.FromSlice(range2))
			next, stop := xiter.Pull(zipped)
			defer stop()

			for i := 0; i < range1End; i++ {
				v, ok := next()
				assert.True(t, ok)
				assert.True(t, v.Ok1)
				assert.True(t, v.Ok2)
				assert.Equal(t, i, v.V1)
				assert.Equal(t, i+range2Start, v.V2)
			}
		})
		t.Run("case2-str", func(t *testing.T) {
			strRange1 := xiter.Map(func(in int) string {
				return strconv.Itoa(in)
			}, xiter.FromSlice(range1))

			zipped := xiter.Zip(xiter.FromSlice(range1), strRange1)
			next, stop := xiter.Pull(zipped)
			defer stop()
			for i := 0; i < range1End; i++ {
				v, ok := next()
				assert.True(t, ok)
				assert.True(t, v.Ok1)
				assert.True(t, v.Ok2)
				assert.Equal(t, i, v.V1)
				assert.Equal(t, fmt.Sprintf("%d", i), v.V2)
			}
		})

		t.Run("case3-empty-one", func(t *testing.T) {
			zipped := xiter.Zip(xiter.FromSlice([]int{}), xiter.FromSlice(range1))
			next, stop := xiter.Pull(zipped)
			defer stop()

			for i := 0; i < range1End; i++ {
				v, ok := next()
				assert.True(t, ok)
				assert.False(t, v.Ok1)
				assert.True(t, v.Ok2)
				assert.Equal(t, i, v.V2)
			}
		})
	})

	t.Run("zip2", func(t *testing.T) {
		t.Run("case1", func(t *testing.T) {
			zipped := xiter.Zip2(xiter.FromSliceIdx(range1), xiter.FromSliceIdx(range2))
			next, stop := xiter.Pull(zipped)
			defer stop()
			for i := 0; i < range1End; i++ {
				v, ok := next()
				assert.True(t, ok)
				assert.True(t, v.Ok1)
				assert.True(t, v.Ok2)
				assert.Equal(t, i, v.V1)
				assert.Equal(t, i+range2Start, v.V2)
			}
		})

		t.Run("case2", func(t *testing.T) {
			zipped := xiter.Zip2(xiter.FromSliceIdx([]int{}), xiter.FromSliceIdx(range1))
			next, stop := xiter.Pull(zipped)
			defer stop()
			for i := 0; i < range1End; i++ {
				v, ok := next()
				assert.True(t, ok)
				assert.False(t, v.Ok1)
				assert.True(t, v.Ok2)
				assert.Equal(t, i, v.V2)
			}
		})
	})

	t.Run("replace", func(t *testing.T) {
		t.Run("case_replace_1", func(t *testing.T) {
			replacedSeq := xiter.Replace(xiter.FromSlice(range1), 1, 2, 1)
			arr := xiter.ToSlice(replacedSeq)
			assert.Equal(t, 2, arr[1])

			for i := 0; i < 10; i++ {
				assert.Len(t, xiter.ToSlice(xiter.Limit(replacedSeq, i)), i)
			}
		})

		t.Run("case_replace_0", func(t *testing.T) {
			replacedSeq := xiter.Replace(xiter.FromSlice(range1), 1, 2, 0)
			arr := xiter.ToSlice(replacedSeq)
			assert.Equal(t, range1, arr)
		})

		t.Run("case_replace_all", func(t *testing.T) {
			arr := xiter.ToSlice(xiter.Replace(xiter.FromSlice(bytes.Repeat([]byte("b"), 1024)), 'b', 'a', -1))
			arr1 := xiter.ToSlice(xiter.ReplaceAll(xiter.FromSlice(bytes.Repeat([]byte("b"), 1024)), 'b', 'a'))
			assert.Equal(t, bytes.Repeat([]byte("a"), 1024), arr)
			assert.Equal(t, bytes.Repeat([]byte("a"), 1024), arr1)
		})
		t.Run("replace limit", func(t *testing.T) {
			replacedSeq := xiter.Replace(xiter.FromSlice(bytes.Repeat([]byte("b"), 1024)), 'b', 'a', -1)
			for i := 1; i < 100; i++ {

				assert.Equal(t, bytes.Repeat([]byte("a"), i), xiter.ToSlice(xiter.Limit(replacedSeq, i)))
			}
		})
	})
}
