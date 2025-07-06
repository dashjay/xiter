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
		res := xiter.ToSlice(xiter.FromSliceReverse(_range(0, 1000)))
		for i := 0; i < 1000; i++ {
			assert.Equal(t, 1000-i-1, res[i])
		}
	})

	t.Run("reverse", func(t *testing.T) {
		res := xiter.ToSlice(xiter.Reverse(xiter.FromSlice(_range(0, 1000))))
		for i := 0; i < 1000; i++ {
			assert.Equal(t, 1000-i-1, res[i])
		}
	})
}
