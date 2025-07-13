package optional_test

import (
	"testing"

	"github.com/dashjay/xiter/pkg/optional"
	"github.com/stretchr/testify/assert"
)

func TestOptional(t *testing.T) {
	t.Parallel()

	o := optional.FromValue(1)
	assert.True(t, o.Ok())
	assert.Equal(t, 1, o.Must())
	assert.Equal(t, 1, o.ValueOr(2))
	assert.Equal(t, 1, o.ValueOrZero())
	assert.Equal(t, 1, *o.Ptr())

	o = optional.Empty[int]()
	assert.False(t, o.Ok())
	assert.Panics(t, func() {
		o.Must()
	})
	assert.Equal(t, 2, o.ValueOr(2))
	assert.Equal(t, 0, o.ValueOrZero())
	assert.Nil(t, o.Ptr())

	var m = make(map[int]int)
	x, exists := m[1]
	o = optional.FromValue2(x, exists)
	assert.False(t, o.Ok())
}
