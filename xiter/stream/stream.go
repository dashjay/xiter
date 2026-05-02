//go:build go1.23

// Package stream provides concurrent stream processing utilities for iter.Seq.
//
// These functions enable parallel and batched processing of sequences,
// useful for CPU-bound or I/O-bound workloads where goroutines can
// improve throughput.
package stream

import (
	"sync"

	"github.com/dashjay/xiter/xiter"
)

// ParallelMap applies fn to each element in seq concurrently using n worker
// goroutines. Results are yielded in non-deterministic order as workers
// complete.
//
// Example:
//
//	results := stream.ParallelMap(
//		xiter.FromSlice([]int{1, 2, 3, 4, 5}),
//		func(v int) int { return v * 2 },
//		3,
//	)
func ParallelMap[T, R any](seq xiter.Seq[T], fn func(T) R, n int) xiter.Seq[R] {
	return func(yield func(R) bool) {
		if n <= 0 {
			n = 1
		}

		done := make(chan struct{})
		defer close(done)

		in := make(chan T)
		out := make(chan R)

		// Feed input sequence into channel
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

		// Start worker pool
		var wg sync.WaitGroup
		for i := 0; i < n; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				for v := range in {
					r := fn(v)
					select {
					case out <- r:
					case <-done:
						return
					}
				}
			}()
		}

		// Close output channel when all workers finish
		go func() {
			wg.Wait()
			close(out)
		}()

		// Yield results
		for r := range out {
			if !yield(r) {
				return
			}
		}
	}
}

// FanIn merges multiple sequences into a single sequence.
// Values from all input sequences are interleaved as they arrive.
// The returned sequence ends when all input sequences are exhausted.
//
// Example:
//
//	seq1 := xiter.FromSlice([]int{1, 2, 3})
//	seq2 := xiter.FromSlice([]int{4, 5, 6})
//	for v := range stream.FanIn(seq1, seq2) {
//		fmt.Println(v) // may print 1,4,2,5,3,6 in any interleaving
//	}
func FanIn[T any](seqs ...xiter.Seq[T]) xiter.Seq[T] {
	return func(yield func(T) bool) {
		if len(seqs) == 0 {
			return
		}

		done := make(chan struct{})
		defer close(done)

		out := make(chan T)
		var wg sync.WaitGroup

		// Start a goroutine for each input sequence
		for _, seq := range seqs {
			wg.Add(1)
			go func(s xiter.Seq[T]) {
				defer wg.Done()
				for v := range s {
					select {
					case out <- v:
					case <-done:
						return
					}
				}
			}(seq)
		}

		go func() {
			wg.Wait()
			close(out)
		}()

		for v := range out {
			if !yield(v) {
				return
			}
		}
	}
}

// Batch groups elements from seq into slices of at most n elements.
// The last batch may contain fewer than n elements.
//
// Example:
//
//	seq := xiter.FromSlice([]int{1, 2, 3, 4, 5})
//	for batch := range stream.Batch(seq, 2) {
//		fmt.Println(batch) // [1 2] [3 4] [5]
//	}
func Batch[T any](seq xiter.Seq[T], n int) xiter.Seq[[]T] {
	return func(yield func([]T) bool) {
		if n <= 0 {
			return
		}
		batch := make([]T, 0, n)
		for v := range seq {
			batch = append(batch, v)
			if len(batch) == n {
				if !yield(batch) {
					return
				}
				batch = make([]T, 0, n)
			}
		}
		if len(batch) > 0 {
			yield(batch)
		}
	}
}
