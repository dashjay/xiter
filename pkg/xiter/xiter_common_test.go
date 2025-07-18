package xiter_test

import (
	"github.com/dashjay/xiter/pkg/xiter"
	"github.com/stretchr/testify/assert"
	"testing"
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
}
