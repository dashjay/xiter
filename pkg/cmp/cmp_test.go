// Copy from https://pkg.go.dev/cmp
// Copyright 2023 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cmp_test

import (
	"fmt"
	"math"
	"sort"
	"strings"
	"testing"
	"unsafe"

	"github.com/dashjay/xiter/pkg/cmp"
)

var negzero = math.Copysign(0, -1)
var nonnilptr uintptr = uintptr(unsafe.Pointer(&negzero))
var nilptr uintptr = uintptr(unsafe.Pointer(nil))

var tests = []struct {
	x, y    any
	compare int
}{
	{1, 2, -1},
	{1, 1, 0},
	{2, 1, +1},
	{"a", "aa", -1},
	{"a", "a", 0},
	{"aa", "a", +1},
	{1.0, 1.1, -1},
	{1.1, 1.1, 0},
	{1.1, 1.0, +1},
	{math.Inf(1), math.Inf(1), 0},
	{math.Inf(-1), math.Inf(-1), 0},
	{math.Inf(-1), 1.0, -1},
	{1.0, math.Inf(-1), +1},
	{math.Inf(1), 1.0, +1},
	{1.0, math.Inf(1), -1},
	{math.NaN(), math.NaN(), 0},
	{0.0, math.NaN(), +1},
	{math.NaN(), 0.0, -1},
	{math.NaN(), math.Inf(-1), -1},
	{math.Inf(-1), math.NaN(), +1},
	{0.0, 0.0, 0},
	{negzero, negzero, 0},
	{negzero, 0.0, 0},
	{0.0, negzero, 0},
	{negzero, 1.0, -1},
	{negzero, -1.0, +1},
	{nilptr, nonnilptr, -1},
	{nonnilptr, nilptr, 1},
	{nonnilptr, nonnilptr, 0},
}

func TestLess(t *testing.T) {
	for _, test := range tests {
		var b bool
		switch test.x.(type) {
		case int:
			b = cmp.Less(test.x.(int), test.y.(int))
		case string:
			b = cmp.Less(test.x.(string), test.y.(string))
		case float64:
			b = cmp.Less(test.x.(float64), test.y.(float64))
		case uintptr:
			b = cmp.Less(test.x.(uintptr), test.y.(uintptr))
		}
		if b != (test.compare < 0) {
			t.Errorf("Less(%v, %v) == %t, want %t", test.x, test.y, b, test.compare < 0)
		}
	}
}

func TestCompare(t *testing.T) {
	for _, test := range tests {
		var c int
		switch test.x.(type) {
		case int:
			c = cmp.Compare(test.x.(int), test.y.(int))
		case string:
			c = cmp.Compare(test.x.(string), test.y.(string))
		case float64:
			c = cmp.Compare(test.x.(float64), test.y.(float64))
		case uintptr:
			c = cmp.Compare(test.x.(uintptr), test.y.(uintptr))
		}
		if c != test.compare {
			t.Errorf("Compare(%v, %v) == %d, want %d", test.x, test.y, c, test.compare)
		}
	}
}

func TestSort(t *testing.T) {
	// Test that our comparison function is consistent with
	// sort.Float64s.
	input := []float64{1.0, 0.0, negzero, math.Inf(1), math.Inf(-1), math.NaN()}
	sort.Float64s(input)
	for i := 0; i < len(input)-1; i++ {
		if cmp.Less(input[i+1], input[i]) {
			t.Errorf("Less sort mismatch at %d in %v", i, input)
		}
		if cmp.Compare(input[i], input[i+1]) > 0 {
			t.Errorf("Compare sort mismatch at %d in %v", i, input)
		}
	}
}

func TestOr(t *testing.T) {
	cases := []struct {
		in   []int
		want int
	}{
		{nil, 0},
		{[]int{0}, 0},
		{[]int{1}, 1},
		{[]int{0, 2}, 2},
		{[]int{3, 0}, 3},
		{[]int{4, 5}, 4},
		{[]int{0, 6, 7}, 6},
	}
	for _, tc := range cases {
		if got := cmp.Or(tc.in...); got != tc.want {
			t.Errorf("cmp.Or(%v) = %v; want %v", tc.in, got, tc.want)
		}
	}
}

func ExampleOr() {
	// Suppose we have some user input
	// that may or may not be an empty string
	userInput1 := ""
	userInput2 := "some text"

	fmt.Println(cmp.Or(userInput1, "default"))
	fmt.Println(cmp.Or(userInput2, "default"))
	fmt.Println(cmp.Or(userInput1, userInput2, "default"))
	// Output:
	// default
	// some text
	// some text
}

type Order struct {
	Product  string
	Customer string
	Price    float64
}

type Orders []Order

func (o Orders) Len() int {
	return len(o)
}
func (o Orders) Less(i, j int) bool {
	a, b := o[i], o[j]
	if cmp.Or(
		strings.Compare(a.Customer, b.Customer),
		strings.Compare(a.Product, b.Product),
		cmp.Compare(b.Price, a.Price)) < 0 {
		return true
	} else {
		return false
	}
}

func (o Orders) Swap(i, j int) {
	o[i], o[j] = o[j], o[i]
}

func ExampleOr_sort() {
	orders := []Order{
		{"foo", "alice", 1.00},
		{"bar", "bob", 3.00},
		{"baz", "carol", 4.00},
		{"foo", "alice", 2.00},
		{"bar", "carol", 1.00},
		{"foo", "bob", 4.00},
	}
	//Sort by customer first, product second, and last by higher price
	sort.Sort(Orders(orders))

	// wait for the  implement of slices.SortFunc
	//SortFunc(orders, func(a, b Order) int {
	//	return cmp.Or(
	//		strings.Compare(a.Customer, b.Customer),
	//		strings.Compare(a.Product, b.Product),
	//		cmp.Compare(b.Price, a.Price),
	//	)
	//})
	for _, order := range orders {
		fmt.Printf("%s %s %.2f\n", order.Product, order.Customer, order.Price)
	}

	// Output:
	// foo alice 2.00
	// foo alice 1.00
	// bar bob 3.00
	// foo bob 4.00
	// bar carol 1.00
	// baz carol 4.00
}

func ExampleLess() {
	fmt.Println(cmp.Less(1, 2))
	fmt.Println(cmp.Less("a", "aa"))
	fmt.Println(cmp.Less(1.0, math.NaN()))
	fmt.Println(cmp.Less(math.NaN(), 1.0))
	// Output:
	// true
	// true
	// false
	// true
}

func ExampleCompare() {
	fmt.Println(cmp.Compare(1, 2))
	fmt.Println(cmp.Compare("a", "aa"))
	fmt.Println(cmp.Compare(1.5, 1.5))
	fmt.Println(cmp.Compare(math.NaN(), 1.0))
	// Output:
	// -1
	// -1
	// 0
	// -1
}
