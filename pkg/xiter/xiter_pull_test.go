// Copyright 2023 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package xiter_test

import (
	"fmt"
	"github.com/dashjay/xiter/pkg/xiter"
	"github.com/stretchr/testify/assert"
	"runtime"
	"testing"
	"time"
)

func count(n int) xiter.Seq[int] {
	return func(yield func(int) bool) {
		for i := 0; i < n; i++ {
			if !yield(i) {
				break
			}
		}
	}
}

func squares(n int) xiter.Seq2[int, int64] {
	return func(yield func(int, int64) bool) {
		for i := 0; i < n; i++ {
			if !yield(i, int64(i)*int64(i)) {
				break
			}
		}
	}
}

func TestPull(t *testing.T) {
	//ver := runtime.Version()
	//var major, minor int
	//_, err := fmt.Sscanf(ver, "go%d.%d", &major, &minor)
	//assert.Nil(t, err)
	//if minor < 23 {
	//	t.Skipf("skil due to go version is lower than 1.23, goversion: %s", ver)
	//	return
	//}
	for end := 0; end <= 3; end++ {
		t.Run(fmt.Sprint(end), func(t *testing.T) {
			ng := stableNumGoroutine()
			wantNG := func(want int) {
				// old version goroutine num will count down later
				time.Sleep(50 * time.Millisecond)
				if xg := runtime.NumGoroutine() - ng; xg != want {
					t.Helper()
					t.Errorf("have %d extra goroutines, want %d", xg, want)
				}
			}
			wantNG(0)
			next, stop := xiter.Pull(count(3))
			wantNG(1)
			for i := 0; i < end; i++ {
				v, ok := next()
				if v != i || ok != true {
					t.Fatalf("next() = %d, %v, want %d, %v", v, ok, i, true)
				}
				wantNG(1)
			}
			wantNG(1)
			if end < 3 {
				stop()
				wantNG(0)
			}
			for i := 0; i < 2; i++ {
				v, ok := next()
				if v != 0 || ok != false {
					t.Fatalf("next() = %d, %v, want %d, %v", v, ok, 0, false)
				}
				wantNG(0)
			}
			wantNG(0)

			stop()
			stop()
			stop()
			wantNG(0)
		})
	}
}

func TestPull2(t *testing.T) {
	//ver := runtime.Version()
	//var major, minor int
	//_, err := fmt.Sscanf(ver, "go%d.%d", &major, &minor)
	//assert.Nil(t, err)
	//if minor < 23 {
	//	t.Skipf("skil due to go version is lower than 1.23, goversion: %s", ver)
	//	return
	//}
	for end := 0; end <= 3; end++ {
		t.Run(fmt.Sprint(end), func(t *testing.T) {
			ng := stableNumGoroutine()
			wantNG := func(want int) {
				time.Sleep(50 * time.Millisecond)
				if xg := runtime.NumGoroutine() - ng; xg != want {
					t.Helper()
					t.Errorf("have %d extra goroutines, want %d", xg, want)
				}
			}
			wantNG(0)
			next, stop := xiter.Pull2(squares(3))
			wantNG(1)
			for i := 0; i < end; i++ {
				k, v, ok := next()
				if k != i || v != int64(i*i) || ok != true {
					t.Fatalf("next() = %d, %d, %v, want %d, %d, %v", k, v, ok, i, i*i, true)
				}
				wantNG(1)
			}
			wantNG(1)
			if end < 3 {
				stop()
				wantNG(0)
			}
			for i := 0; i < 2; i++ {
				k, v, ok := next()
				if v != 0 || ok != false {
					t.Fatalf("next() = %d, %d, %v, want %d, %d, %v", k, v, ok, 0, 0, false)
				}
				wantNG(0)
			}
			wantNG(0)

			stop()
			stop()
			stop()
			wantNG(0)
		})
	}
}

// stableNumGoroutine is like NumGoroutine but tries to ensure stability of
// the value by letting any exiting goroutines finish exiting.
func stableNumGoroutine() int {
	// The idea behind stablizing the value of NumGoroutine is to
	// see the same value enough times in a row in between calls to
	// runtime.Gosched. With GOMAXPROCS=1, we're trying to make sure
	// that other goroutines run, so that they reach a stable point.
	// It's not guaranteed, because it is still possible for a goroutine
	// to Gosched back into itself, so we require NumGoroutine to be
	// the same 100 times in a row. This should be more than enough to
	// ensure all goroutines get a chance to run to completion (or to
	// some block point) for a small group of test goroutines.
	defer runtime.GOMAXPROCS(runtime.GOMAXPROCS(1))

	c := 0
	ng := runtime.NumGoroutine()
	for i := 0; i < 1000; i++ {
		nng := runtime.NumGoroutine()
		if nng == ng {
			c++
		} else {
			c = 0
			ng = nng
		}
		if c >= 100 {
			// The same value 100 times in a row is good enough.
			return ng
		}
		runtime.Gosched()
	}
	panic("failed to stabilize NumGoroutine after 1000 iterations")
}

func TestPullDoubleNext(t *testing.T) {
	next, _ := xiter.Pull(doDoubleNext())
	nextSlot = next
	next()
	if nextSlot != nil {
		t.Fatal("double next did not fail")
	}
}

var nextSlot func() (int, bool)

func doDoubleNext() xiter.Seq[int] {
	return func(_ func(int) bool) {
		defer func() {
			if recover() != nil {
				nextSlot = nil
			}
		}()
		nextSlot()
	}
}

func TestPullDoubleNext2(t *testing.T) {
	next, _ := xiter.Pull2(doDoubleNext2())
	nextSlot2 = next
	next()
	if nextSlot2 != nil {
		t.Fatal("double next did not fail")
	}
}

var nextSlot2 func() (int, int, bool)

func doDoubleNext2() xiter.Seq2[int, int] {
	return func(_ func(int, int) bool) {
		defer func() {
			if recover() != nil {
				nextSlot2 = nil
			}
		}()
		nextSlot2()
	}
}

func TestPullDoubleYield(t *testing.T) {
	_, stop := xiter.Pull(storeYield())
	defer func() {
		if recover() != nil {
			yieldSlot = nil
		}
		stop()
	}()
	yieldSlot(5)
	if yieldSlot != nil {
		t.Fatal("double yield did not fail")
	}
}

func storeYield() xiter.Seq[int] {
	return func(yield func(int) bool) {
		yieldSlot = yield
		if !yield(5) {
			return
		}
	}
}

var yieldSlot func(int) bool

func TestPullDoubleYield2(t *testing.T) {
	_, stop := xiter.Pull2(storeYield2())
	defer func() {
		if recover() != nil {
			yieldSlot2 = nil
		}
		stop()
	}()
	yieldSlot2(23, 77)
	if yieldSlot2 != nil {
		t.Fatal("double yield did not fail")
	}
}

func storeYield2() xiter.Seq2[int, int] {
	return func(yield func(int, int) bool) {
		yieldSlot2 = yield
		if !yield(23, 77) {
			return
		}
	}
}

var yieldSlot2 func(int, int) bool

func TestPullPanic(t *testing.T) {
	t.Run("next", func(t *testing.T) {
		next, stop := xiter.Pull(panicSeq())
		if !panicsWith("boom", func() { next() }) {
			t.Fatal("failed to propagate panic on first next")
		}
		// Make sure we don't panic again if we try to call next or stop.
		if _, ok := next(); ok {
			t.Fatal("next returned true after iterator panicked")
		}
		// Calling stop again should be a no-op.
		stop()
	})
	t.Run("stop", func(t *testing.T) {
		if testSkip(t) {
			return
		}
		next, stop := xiter.Pull(panicCleanupSeq())
		x, ok := next()
		if !ok || x != 55 {
			t.Fatalf("expected (55, true) from next, got (%d, %t)", x, ok)
		}
		if !panicsWith("boom", func() { stop() }) {
			t.Fatal("failed to propagate panic on stop")
		}
		// Make sure we don't panic again if we try to call next or stop.
		if _, ok := next(); ok {
			t.Fatal("next returned true after iterator panicked")
		}
		// Calling stop again should be a no-op.
		stop()
	})
}

func panicSeq() xiter.Seq[int] {
	return func(yield func(int) bool) {
		panic("boom")
	}
}

func panicCleanupSeq() xiter.Seq[int] {
	return func(yield func(int) bool) {
		for {
			if !yield(55) {
				panic("boom")
			}
		}
	}
}

func TestPull2Panic(t *testing.T) {
	t.Run("next", func(t *testing.T) {
		next, stop := xiter.Pull2(panicSeq2())
		if !panicsWith("boom", func() { next() }) {
			t.Fatal("failed to propagate panic on first next")
		}
		// Make sure we don't panic again if we try to call next or stop.
		if _, _, ok := next(); ok {
			t.Fatal("next returned true after iterator panicked")
		}
		// Calling stop again should be a no-op.
		stop()
	})
	t.Run("stop", func(t *testing.T) {
		if testSkip(t) {
			return
		}
		next, stop := xiter.Pull2(panicCleanupSeq2())
		x, y, ok := next()
		if !ok || x != 55 || y != 100 {
			t.Fatalf("expected (55, 100, true) from next, got (%d, %d, %t)", x, y, ok)
		}
		if !panicsWith("boom", func() { stop() }) {
			t.Fatal("failed to propagate panic on stop")
		}
		// Make sure we don't panic again if we try to call next or stop.
		if _, _, ok := next(); ok {
			t.Fatal("next returned true after iterator panicked")
		}
		// Calling stop again should be a no-op.
		stop()
	})
}

func panicSeq2() xiter.Seq2[int, int] {
	return func(yield func(int, int) bool) {
		panic("boom")
	}
}

func panicCleanupSeq2() xiter.Seq2[int, int] {
	return func(yield func(int, int) bool) {
		for {
			if !yield(55, 100) {
				panic("boom")
			}
		}
	}
}

func panicsWith(v any, f func()) (panicked bool) {
	defer func() {
		if r := recover(); r != nil {
			if r != v {
				panic(r)
			}
			panicked = true
		}
	}()
	f()
	return
}

func TestPullGoexit(t *testing.T) {
	if testSkip(t) {
		return
	}
	t.Run("next", func(t *testing.T) {
		var next func() (int, bool)
		var stop func()
		if !goexits(t, func() {
			next, stop = xiter.Pull(goexitSeq())
			next()
		}) {
			t.Fatal("failed to Goexit from next")
		}
		if x, ok := next(); x != 0 || ok {
			t.Fatal("iterator returned valid value after iterator Goexited")
		}
		stop()
	})
	t.Run("stop", func(t *testing.T) {
		next, stop := xiter.Pull(goexitCleanupSeq())
		x, ok := next()
		if !ok || x != 55 {
			t.Fatalf("expected (55, true) from next, got (%d, %t)", x, ok)
		}
		if !goexits(t, func() { stop() }) {
			t.Fatal("failed to Goexit from stop")
		}
		// Make sure we don't panic again if we try to call next or stop.
		if x, ok := next(); x != 0 || ok {
			t.Fatal("next returned true or non-zero value after iterator Goexited")
		}
		// Calling stop again should be a no-op.
		stop()
	})
}

func goexitSeq() xiter.Seq[int] {
	return func(yield func(int) bool) {
		runtime.Goexit()
	}
}

func goexitCleanupSeq() xiter.Seq[int] {
	return func(yield func(int) bool) {
		for {
			if !yield(55) {
				runtime.Goexit()
			}
		}
	}
}

func TestPull2Goexit(t *testing.T) {
	if testSkip(t) {
		return
	}
	t.Run("next", func(t *testing.T) {
		var next func() (int, int, bool)
		var stop func()
		if !goexits(t, func() {
			next, stop = xiter.Pull2(goexitSeq2())
			next()
		}) {
			t.Fatal("failed to Goexit from next")
		}
		if x, y, ok := next(); x != 0 || y != 0 || ok {
			t.Fatal("iterator returned valid value after iterator Goexited")
		}
		stop()
	})
	t.Run("stop", func(t *testing.T) {
		next, stop := xiter.Pull2(goexitCleanupSeq2())
		x, y, ok := next()
		if !ok || x != 55 || y != 100 {
			t.Fatalf("expected (55, 100, true) from next, got (%d, %d, %t)", x, y, ok)
		}
		if !goexits(t, func() {
			stop()
		}) {
			t.Fatal("failed to Goexit from stop")
		}
		// Make sure we don't panic again if we try to call next or stop.
		if x, y, ok := next(); x != 0 || y != 0 || ok {
			t.Fatal("next returned true or non-zero after iterator Goexited")
		}
		// Calling stop again should be a no-op.
		stop()
	})
}

func goexitSeq2() xiter.Seq2[int, int] {
	return func(yield func(int, int) bool) {
		runtime.Goexit()
	}
}

func goexitCleanupSeq2() xiter.Seq2[int, int] {
	return func(yield func(int, int) bool) {
		for {
			if !yield(55, 100) {
				runtime.Goexit()
			}
		}
	}
}

func goexits(t *testing.T, f func()) bool {
	t.Helper()

	exit := make(chan bool)
	go func() {
		cleanExit := false
		defer func() {
			err := recover()
			if err != nil {
				t.Errorf("go exits error: %s", err)
			}
			exit <- err == nil && !cleanExit
		}()
		f()
		cleanExit = true
	}()
	return <-exit
}

func TestPullImmediateStop(t *testing.T) {
	if testSkip(t) {
		return
	}
	next, stop := xiter.Pull(panicSeq())
	stop()
	// Make sure we don't panic if we try to call next or stop.
	if _, ok := next(); ok {
		t.Fatal("next returned true after iterator was stopped")
	}
}

func TestPull2ImmediateStop(t *testing.T) {
	if testSkip(t) {
		return
	}
	next, stop := xiter.Pull2(panicSeq2())
	stop()
	// Make sure we don't panic if we try to call next or stop.
	if _, _, ok := next(); ok {
		t.Fatal("next returned true after iterator was stopped")
	}
}

func testSkip(t *testing.T) bool {
	ver := runtime.Version()
	var major, minor int
	_, err := fmt.Sscanf(ver, "go%d.%d", &major, &minor)
	assert.Nil(t, err)
	if minor < 23 {
		t.Skipf("skil due to go version is lower than 1.23, goversion: %s", ver)
		return true
	}
	return false
}
