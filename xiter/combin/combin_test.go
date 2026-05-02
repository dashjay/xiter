//go:build go1.23

package combin_test

import (
	"testing"

	"github.com/dashjay/xiter/xiter"
	"github.com/dashjay/xiter/xiter/combin"
)

func toSlice[T any](seq xiter.Seq[T]) []T {
	var out []T
	for v := range seq {
		out = append(out, v)
	}
	return out
}

func TestCombinations(t *testing.T) {
	t.Run("nCk basic", func(t *testing.T) {
		seq := xiter.FromSlice([]string{"a", "b", "c"})
		combs := toSlice(combin.Combinations(seq, 2))
		if len(combs) != 3 {
			t.Fatalf("Combinations 3C2: got %d, want 3", len(combs))
		}
		expected := [][]string{{"a", "b"}, {"a", "c"}, {"b", "c"}}
		for i, c := range combs {
			for j, v := range c {
				if v != expected[i][j] {
					t.Fatalf("Combinations[%d]: got %v, want %v", i, combs, expected)
				}
			}
		}
	})

	t.Run("k=0", func(t *testing.T) {
		combs := toSlice(combin.Combinations(xiter.FromSlice([]int{1, 2, 3}), 0))
		if len(combs) != 0 {
			t.Fatalf("Combinations k=0: got %d, want 0", len(combs))
		}
	})

	t.Run("k=n", func(t *testing.T) {
		combs := toSlice(combin.Combinations(xiter.FromSlice([]int{1, 2, 3}), 3))
		if len(combs) != 1 || len(combs[0]) != 3 {
			t.Fatalf("Combinations 3C3: got %v, want [[1 2 3]]", combs)
		}
	})

	t.Run("k > n", func(t *testing.T) {
		combs := toSlice(combin.Combinations(xiter.FromSlice([]int{1, 2}), 5))
		if len(combs) != 0 {
			t.Fatalf("Combinations k>n: got %d, want 0", len(combs))
		}
	})

	t.Run("empty", func(t *testing.T) {
		combs := toSlice(combin.Combinations(xiter.FromSlice([]int{}), 2))
		if len(combs) != 0 {
			t.Fatalf("Combinations empty: got %d, want 0", len(combs))
		}
	})

	t.Run("k=1", func(t *testing.T) {
		combs := toSlice(combin.Combinations(xiter.FromSlice([]int{10, 20, 30}), 1))
		if len(combs) != 3 {
			t.Fatalf("Combinations k=1: got %d, want 3", len(combs))
		}
	})

	t.Run("negative k", func(t *testing.T) {
		combs := toSlice(combin.Combinations(xiter.FromSlice([]int{1, 2, 3}), -1))
		if len(combs) != 0 {
			t.Fatalf("Combinations k=-1: got %d, want 0", len(combs))
		}
	})

	t.Run("5C3 count", func(t *testing.T) {
		combs := toSlice(combin.Combinations(xiter.FromSlice([]int{1, 2, 3, 4, 5}), 3))
		if len(combs) != 10 {
			t.Fatalf("Combinations 5C3: got %d, want 10", len(combs))
		}
	})
}

func TestPermutations(t *testing.T) {
	t.Run("3P2 basic", func(t *testing.T) {
		seq := xiter.FromSlice([]string{"a", "b", "c"})
		perms := toSlice(combin.Permutations(seq, 2))
		if len(perms) != 6 {
			t.Fatalf("Permutations 3P2: got %d, want 6", len(perms))
		}
		// Verify all permutations are unique
		seen := make(map[string]bool)
		for _, p := range perms {
			key := p[0] + p[1]
			if seen[key] {
				t.Fatalf("Permutations: duplicate %v", p)
			}
			seen[key] = true
		}
		// Verify expected set
		expected := []string{"ab", "ac", "ba", "bc", "ca", "cb"}
		for _, e := range expected {
			if !seen[e] {
				t.Fatalf("Permutations: missing %s", e)
			}
		}
	})

	t.Run("k=0", func(t *testing.T) {
		perms := toSlice(combin.Permutations(xiter.FromSlice([]int{1, 2, 3}), 0))
		if len(perms) != 0 {
			t.Fatalf("Permutations k=0: got %d, want 0", len(perms))
		}
	})

	t.Run("k=n", func(t *testing.T) {
		perms := toSlice(combin.Permutations(xiter.FromSlice([]int{1, 2, 3}), 3))
		if len(perms) != 6 {
			t.Fatalf("Permutations 3P3: got %d, want 6", len(perms))
		}
	})

	t.Run("k > n", func(t *testing.T) {
		perms := toSlice(combin.Permutations(xiter.FromSlice([]int{1, 2}), 5))
		if len(perms) != 0 {
			t.Fatalf("Permutations k>n: got %d, want 0", len(perms))
		}
	})

	t.Run("empty", func(t *testing.T) {
		perms := toSlice(combin.Permutations(xiter.FromSlice([]int{}), 2))
		if len(perms) != 0 {
			t.Fatalf("Permutations empty: got %d, want 0", len(perms))
		}
	})

	t.Run("single element", func(t *testing.T) {
		perms := toSlice(combin.Permutations(xiter.FromSlice([]int{42}), 1))
		if len(perms) != 1 || perms[0][0] != 42 {
			t.Fatalf("Permutations single: got %v, want [[42]]", perms)
		}
	})

	t.Run("4P3 count", func(t *testing.T) {
		perms := toSlice(combin.Permutations(xiter.FromSlice([]int{1, 2, 3, 4}), 3))
		if len(perms) != 24 {
			t.Fatalf("Permutations 4P3: got %d, want 24", len(perms))
		}
	})
}

func TestProduct(t *testing.T) {
	t.Run("2x3 basic", func(t *testing.T) {
		colors := xiter.FromSlice([]string{"red", "blue"})
		sizes := xiter.FromSlice([]string{"S", "M", "L"})
		products := toSlice(combin.Product(colors, sizes))
		if len(products) != 6 {
			t.Fatalf("Product 2x3: got %d, want 6", len(products))
		}
		expected := [][]string{
			{"red", "S"}, {"red", "M"}, {"red", "L"},
			{"blue", "S"}, {"blue", "M"}, {"blue", "L"},
		}
		for i, p := range products {
			for j, v := range p {
				if v != expected[i][j] {
					t.Fatalf("Product[%d]: got %v, want %v", i, products, expected)
				}
			}
		}
	})

	t.Run("single seq", func(t *testing.T) {
		products := toSlice(combin.Product(xiter.FromSlice([]int{1, 2, 3})))
		if len(products) != 3 {
			t.Fatalf("Product single: got %d, want 3", len(products))
		}
	})

	t.Run("no seqs", func(t *testing.T) {
		products := toSlice(combin.Product[int]())
		if len(products) != 0 {
			t.Fatalf("Product none: got %d, want 0", len(products))
		}
	})

	t.Run("empty seq", func(t *testing.T) {
		products := toSlice(combin.Product(
			xiter.FromSlice([]int{1, 2}),
			xiter.FromSlice([]int{}),
		))
		if len(products) != 0 {
			t.Fatalf("Product empty: got %d, want 0", len(products))
		}
	})

	t.Run("single element seqs", func(t *testing.T) {
		products := toSlice(combin.Product(
			xiter.FromSlice([]int{1}),
			xiter.FromSlice([]int{2}),
			xiter.FromSlice([]int{3}),
		))
		if len(products) != 1 || len(products[0]) != 3 {
			t.Fatalf("Product single elements: got %v, want [[1 2 3]]", products)
		}
	})
}

func TestCombinationsEarlyStop(t *testing.T) {
	seq := xiter.FromSlice([]int{1, 2, 3, 4})
	count := 0
	for range combin.Combinations(seq, 2) {
		count++
		if count >= 2 {
			break
		}
	}
	if count != 2 {
		t.Fatalf("Combinations early stop: count = %d, want 2", count)
	}
}

func TestPermutationsEarlyStop(t *testing.T) {
	seq := xiter.FromSlice([]int{1, 2, 3})
	count := 0
	for range combin.Permutations(seq, 3) {
		count++
		if count >= 2 {
			break
		}
	}
	if count != 2 {
		t.Fatalf("Permutations early stop: count = %d, want 2", count)
	}
}

func TestProductEarlyStop(t *testing.T) {
	colors := xiter.FromSlice([]string{"red", "blue", "green"})
	sizes := xiter.FromSlice([]string{"S", "M", "L"})
	count := 0
	for range combin.Product(colors, sizes) {
		count++
		if count >= 3 {
			break
		}
	}
	if count != 3 {
		t.Fatalf("Product early stop: count = %d, want 3", count)
	}
}
