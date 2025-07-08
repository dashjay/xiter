// Package constraints defined constraints for generics tools
package constraints

type Float interface {
	~float32 | ~float64
}

type Integer interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 | ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr
}

type Number interface {
	Integer | Float
}

type Ordered interface {
	Integer | Float | ~string
}
