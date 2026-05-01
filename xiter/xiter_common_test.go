package xiter_test

import (
	"sort"
	"strconv"
	"testing"

	"github.com/dashjay/xiter/xiter"
	"github.com/stretchr/testify/assert"
)

func TestXIterCommon(t *testing.T) {
	t.Run("from slice", func(t *testing.T) {
		input := xiter.FromSlice(_range(0, 1000))
		assert.Equal(t, _range(0, 1000), xiter.ToSlice(input))
	})

	t.Run("from slice idx", func(t *testing.T) {
		input := xiter.FromSliceIdx(_range(0, 1000))
		assert.Equal(t, _range(0, 1000), xiter.ToSliceSeq2Key(input))
		assert.Equal(t, _range(0, 1000), xiter.ToSliceSeq2Value(input))
	})

	t.Run("from slice reverse", func(t *testing.T) {
		reversedSeq := xiter.FromSliceReverse(_range(0, 1000))
		res := xiter.ToSlice(reversedSeq)
		for i := 0; i < 1000; i++ {
			assert.Equal(t, 1000-i-1, res[i])
		}

		res = xiter.ToSlice(xiter.Limit(reversedSeq, 1))
		assert.Len(t, res, 1)
		assert.Equal(t, res[0], 999)
	})

	t.Run("reverse", func(t *testing.T) {
		reversedSeq := xiter.Reverse(xiter.FromSlice(_range(0, 1000)))
		res := xiter.ToSlice(reversedSeq)
		for i := 0; i < 1000; i++ {
			assert.Equal(t, 1000-i-1, res[i])
		}
		res = xiter.ToSlice(xiter.Limit(reversedSeq, 1))
		assert.Len(t, res, 1)
		assert.Equal(t, res[0], 999)

		source := _range(0, 1000)
		reversedSeq = xiter.Reverse(xiter.FromSlice(source)) // r
		for i := 0; i < 10; i++ {
			reversedSeq = xiter.Reverse(reversedSeq)
		}
		res = xiter.ToSlice(reversedSeq)
		for i := 0; i < 1000; i++ {
			assert.Equal(t, 1000-i-1, res[i])
		}
	})

	t.Run("from chan", func(t *testing.T) {
		ch := make(chan int, 10)
		for i := 0; i < 10; i++ {
			ch <- i
		}
		close(ch)
		seq := xiter.FromChan(ch)
		res := xiter.ToSlice(seq)
		assert.Len(t, res, 10)
		assert.Equal(t, _range(0, 10), res)
	})

	t.Run("from chan limit", func(t *testing.T) {
		ch := make(chan int, 10)
		for i := 0; i < 10; i++ {
			ch <- i
		}
		close(ch)
		seq := xiter.FromChan(ch)
		testLimit(t, seq, 1)
	})

	t.Run("difference", func(t *testing.T) {
		left := xiter.FromSlice(_range(0, 10))
		right := xiter.FromSlice(_range(5, 15))
		onlyLeft, onlyRight := xiter.Difference(left, right)
		assert.Equal(t, _range(0, 5), xiter.ToSlice(onlyLeft))
		assert.Equal(t, _range(10, 15), xiter.ToSlice(onlyRight))
	})

	t.Run("intersect", func(t *testing.T) {
		left := xiter.FromSlice(_range(0, 10))
		right := xiter.FromSlice(_range(5, 15))
		assert.Equal(t, _range(5, 10), xiter.ToSlice(xiter.Intersect(left, right)))
		assert.True(t, xiter.Equal(left, xiter.Intersect(left, left)))
	})

	t.Run("mean", func(t *testing.T) {
		// 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10
		seq := xiter.FromSlice(_range(0, 10))
		r := xiter.Mean(seq)
		assert.Equal(t, xiter.Sum(seq)/len(_range(0, 10)), r)
	})

	t.Run("mean by", func(t *testing.T) {
		strSeq := xiter.Map(func(in int) string {
			return strconv.Itoa(in)
		}, xiter.FromSlice(_range(0, 10)))

		r := xiter.MeanBy(strSeq, func(t string) int {
			v, _ := strconv.Atoi(t)
			return v
		})
		assert.Equal(t, xiter.Sum(xiter.FromSlice(_range(0, 10)))/len(_range(0, 10)), r)
	})
	t.Run("moderate", func(t *testing.T) {
		moderate := xiter.ModerateO(xiter.FromSlice([]int{1, 2, 3, 4, 5, 5, 5, 6, 6, 6, 6}))
		assert.True(t, moderate.Ok())
		assert.Equal(t, 6, moderate.Must())
	})

	t.Run("union", func(t *testing.T) {
		left := xiter.FromSlice(_range(0, 10))
		right := xiter.FromSlice(_range(5, 15))
		x := xiter.ToSlice(xiter.Union(left, right))
		sort.Sort(sort.IntSlice(x))
		assert.Equal(t, _range(0, 15), x)

		ut := func(int) struct{} {
			return struct{}{}
		}
		assert.Equal(t, xiter.ToMapFromSeq(xiter.Union(left, right), ut),
			xiter.ToMapFromSeq(xiter.Union(right, left), ut))
		})

		t.Run("cycle", func(t *testing.T) {
			seq := xiter.Cycle(xiter.FromSlice([]int{1, 2, 3}))
			first9 := xiter.ToSlice(xiter.Limit(seq, 9))
			assert.Equal(t, []int{1, 2, 3, 1, 2, 3, 1, 2, 3}, first9)

			emptyCycle := xiter.Cycle(xiter.FromSlice([]int{}))
			assert.Len(t, xiter.ToSlice(emptyCycle), 0)

			singleCycle := xiter.Cycle(xiter.FromSlice([]int{42}))
			first5 := xiter.ToSlice(xiter.Limit(singleCycle, 5))
			assert.Equal(t, []int{42, 42, 42, 42, 42}, first5)
		})

		t.Run("generate", func(t *testing.T) {
			i := 0
			gen := xiter.Generate(func() int {
				i++
				return i
			})
			first5 := xiter.ToSlice(xiter.Limit(gen, 5))
			assert.Equal(t, []int{1, 2, 3, 4, 5}, first5)

			assert.Len(t, xiter.ToSlice(xiter.Limit(gen, 0)), 0)
		})

		t.Run("to chan", func(t *testing.T) {
			seq := xiter.FromSlice(_range(0, 10))
			ch := xiter.ToChan(seq)
			var result []int
			for v := range ch {
				result = append(result, v)
			}
			assert.Equal(t, _range(0, 10), result)

			emptyCh := xiter.ToChan(xiter.FromSlice([]int{}))
			var emptyResult []int
			for v := range emptyCh {
				emptyResult = append(emptyResult, v)
			}
			assert.Len(t, emptyResult, 0)
		})

		t.Run("range", func(t *testing.T) {
			r := xiter.ToSlice(xiter.Range(0, 10, 2))
			assert.Equal(t, []int{0, 2, 4, 6, 8}, r)

			r2 := xiter.ToSlice(xiter.Range(10, 0, -3))
			assert.Equal(t, []int{10, 7, 4, 1}, r2)

			assert.Len(t, xiter.ToSlice(xiter.Range(0, 0, 1)), 0)
			assert.Len(t, xiter.ToSlice(xiter.Range(5, 0, 1)), 0)
			assert.Len(t, xiter.ToSlice(xiter.Range(0, 5, -1)), 0)
			assert.Len(t, xiter.ToSlice(xiter.Range(0, 10, 0)), 0)

			r3 := xiter.ToSlice(xiter.Range(0, 100, 1000))
			assert.Equal(t, []int{0}, r3)

			r4 := xiter.ToSlice(xiter.Range(0, 5, 1))
			assert.Equal(t, []int{0, 1, 2, 3, 4}, r4)
		})

		t.Run("with index", func(t *testing.T) {
			seq := xiter.FromSlice([]string{"a", "b", "c"})
			idxSeq := xiter.WithIndex(seq)

			collected := xiter.ToSliceSeq2Key(idxSeq)
			assert.Equal(t, []int{0, 1, 2}, collected)
			values := xiter.ToSliceSeq2Value(idxSeq)
			assert.Equal(t, []string{"a", "b", "c"}, values)

			emptyIdx := xiter.WithIndex(xiter.FromSlice([]int{}))
			assert.Len(t, xiter.ToSliceSeq2Key(emptyIdx), 0)
		})
}
