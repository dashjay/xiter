package xassert_test

import (
	"testing"

	"github.com/dashjay/xiter/pkg/internal/xassert"
	"github.com/stretchr/testify/assert"
)

func TestAssert(t *testing.T) {
	assert.NotPanics(t, func() {
		xassert.MustBePositive(1)
	})

	assert.Panics(t, func() {
		xassert.MustBePositive(-1)
	})
}
