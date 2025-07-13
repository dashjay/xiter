// Package xiter provides the abstraction of map, slice or channel types into iterators for common processing
// In most scenarios, we DO NOT need to use the xiter package directly
//
// WARNING: golang 1.23 has higher performance on iterating Seq/Seq2 which boost by coroutine
// we provide for low golang version just for compatible usage
package xiter
