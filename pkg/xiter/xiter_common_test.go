package xiter_test

import (
	"sort"
	"strconv"
	"testing"

	"github.com/dashjay/xiter/pkg/xiter"
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
}
