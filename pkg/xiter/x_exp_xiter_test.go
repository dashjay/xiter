package xiter_test

import (
	"fmt"
	"math"
	"strconv"
	"testing"

	"github.com/dashjay/xiter/pkg/xiter"
	"github.com/stretchr/testify/assert"
)

func _range(a, b int) []int {
	var res []int
	for i := a; i < b; i++ {
		res = append(res, i)
	}
	return res
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
	})

	t.Run("equal2", func(t *testing.T) {
		seq1 := xiter.FromSliceIdx(range1)
		assert.True(t, xiter.Equal2(seq1, xiter.FromSliceIdx(_range(range1Start, range1End))))
	})

	t.Run("equal_func", func(t *testing.T) {
		seq1 := xiter.FromSlice(range1)
		assert.True(t, xiter.EqualFunc(seq1, xiter.FromSlice(_range(range1Start, range1End)), func(a int, b int) bool {
			return a == b
		}))
	})

	t.Run("equal_func2", func(t *testing.T) {
		seq1 := xiter.FromSliceIdx(range1)
		assert.True(t, xiter.EqualFunc2(seq1, xiter.FromSliceIdx(_range(range1Start, range1End)), func(k1, k2, v1, v2 int) bool {
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
	})

	t.Run("limit", func(t *testing.T) {
		t.Run("limit zero", func(t *testing.T) {
			seq1 := xiter.FromSlice(range1)
			seq1Limit0 := xiter.Limit(seq1, 0)
			assert.Len(t, xiter.ToSlice(seq1Limit0), 0)
		})

		t.Run("limit one", func(t *testing.T) {
			seq1 := xiter.FromSlice(range1)
			seq1Limit0 := xiter.Limit(seq1, 1)
			assert.Len(t, xiter.ToSlice(seq1Limit0), 1)
		})

		t.Run("limit large", func(t *testing.T) {
			seq1 := xiter.FromSlice(range1)
			seq1Limit0 := xiter.Limit(seq1, math.MaxInt64)
			assert.Len(t, xiter.ToSlice(seq1Limit0), range1End-range1Start)
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
	})

	t.Run("merge2", func(t *testing.T) {
		var seq1 = xiter.Seq2[int, string](func(yield func(int, string) bool) {
			for i := 0; i < 10; i++ {
				if !yield(i, fmt.Sprintf("%d", i)) {
					break
				}
			}
			for i := 20; i < 30; i++ {
				if !yield(i, fmt.Sprintf("%d", i)) {
					break
				}
			}
			for i := 40; i < 50; i++ {
				if !yield(i, fmt.Sprintf("%d", i)) {
					break
				}
			}
		})
		var seq2 = xiter.Seq2[int, string](func(yield func(int, string) bool) {
			for i := 10; i < 20; i++ {
				if !yield(i, fmt.Sprintf("%d", i)) {
					break
				}
			}
			for i := 30; i < 40; i++ {
				if !yield(i, fmt.Sprintf("%d", i)) {
					break
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
}
