package xiter

// A Zipped2 is a pair of zipped key-value pairs,
// one of which may be missing, drawn from two different sequences.
type Zipped2[K1, V1, K2, V2 any] struct {
	K1  K1
	V1  V1
	Ok1 bool // whether K1, V1 are present (if not, they will be zero)
	K2  K2
	V2  V2
	Ok2 bool // whether K2, V2 are present (if not, they will be zero)
}

// A Zipped is a pair of zipped values, one of which may be missing,
// drawn from two different sequences.
type Zipped[V1, V2 any] struct {
	V1  V1
	Ok1 bool // whether V1 is present (if not, it will be zero)
	V2  V2
	Ok2 bool // whether V2 is present (if not, it will be zero)
}
