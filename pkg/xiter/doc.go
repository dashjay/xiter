// Package xiter provides the abstraction of map, slice or channel types into iterators for common processing
//
// In most scenarios, we DO NOT need to use the xiter package directly
//Package xiter implements basic adapters for composing iterator sequences:
//
// - [Concat] and [Concat2] concatenate sequences.
// - [Equal], [Equal2], [EqualFunc], and [EqualFunc2] check whether two sequences contain equal values.
// - [Filter] and [Filter2] filter a sequence according to a function f.
// - [Limit] and [Limit2] truncate a sequence after n items.
// - [Map] and [Map2] apply a function f to a sequence.
// - [Merge], [Merge2], [MergeFunc], and [MergeFunc2] merge two ordered sequences.
// - [Reduce] and [Reduce2] combine the values in a sequence.
// - [Zip] and [Zip2] iterate over two sequences in parallel.

package xiter
