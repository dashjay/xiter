package xassert

import (
	"fmt"

	"github.com/dashjay/xiter/pkg/internal/constraints"
)

func MustBePositive[T constraints.Number](in T) {
	if in < 0 {
		panic(fmt.Sprintf("number %v[%T] must be positive", in, in))
	}
}
