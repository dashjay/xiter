//go:build go1.23

// Package rt provides real-time event stream utilities built on iter.Seq.
//
// These functions enable time-based operations on sequences, such as
// periodic ticks, rate-limiting, and debouncing.
package rt

import (
	"time"

	"github.com/dashjay/xiter/xiter"
)

// Ticker returns a sequence that yields the current time at the specified interval.
// The sequence is unbounded; use with Limit, TakeWhile, etc. to constrain.
//
// Example:
//
//	for t := range rt.Ticker(time.Second) {
//		fmt.Println("tick at", t)
//	}
func Ticker(d time.Duration) xiter.Seq[time.Time] {
	return func(yield func(time.Time) bool) {
		ticker := time.NewTicker(d)
		defer ticker.Stop()
		for t := range ticker.C {
			if !yield(t) {
				return
			}
		}
	}
}

// Throttle limits the rate of values from seq, yielding at most one value
// per duration d.
//
// Example:
//
//	seq := xiter.FromSlice([]int{1, 2, 3, 4, 5})
//	for v := range rt.Throttle(seq, 100*time.Millisecond) {
//		fmt.Println(v) // printed at most every 100ms
//	}
func Throttle[T any](seq xiter.Seq[T], d time.Duration) xiter.Seq[T] {
	return func(yield func(T) bool) {
		ticker := time.NewTicker(d)
		defer ticker.Stop()
		for v := range seq {
			<-ticker.C
			if !yield(v) {
				return
			}
		}
	}
}

// Debounce yields values from seq only after a quiet period of duration d
// has elapsed since the last value was received. If new values arrive before
// the quiet period elapses, the timer resets. The final value is always
// flushed when the input sequence is exhausted.
//
// Example:
//
//	for v := range rt.Debounce(eventStream, 100*time.Millisecond) {
//		fmt.Println("debounced:", v)
//	}
func Debounce[T any](seq xiter.Seq[T], d time.Duration) xiter.Seq[T] {
	return func(yield func(T) bool) {
		done := make(chan struct{})
		defer close(done)

		in := make(chan T, 1)
		go func() {
			defer close(in)
			for v := range seq {
				select {
				case in <- v:
				case <-done:
					return
				}
			}
		}()

		timer := time.NewTimer(d)
		if !timer.Stop() {
			<-timer.C
		}
		defer timer.Stop()

		var lastValue T
		hasPending := false

		for {
			timerC := timer.C
			if !hasPending {
				timerC = nil
			}

			select {
			case v, ok := <-in:
				if !ok {
					if hasPending {
						timer.Stop()
						if !yield(lastValue) {
							return
						}
					}
					return
				}
				lastValue = v
				hasPending = true
				if !timer.Stop() {
					select {
					case <-timer.C:
					default:
					}
				}
				timer.Reset(d)

			case <-timerC:
				if hasPending {
					if !yield(lastValue) {
						return
					}
					hasPending = false
				}
			}
		}
	}
}
