# xiter

A Generic Library for Go with the Power of Iterators

[![codecov](https://codecov.io/gh/dashjay/xiter/graph/badge.svg?token=GTTJNP1MHT)](https://codecov.io/gh/dashjay/xiter)
![ci.yml](https://github.com/dashjay/xiter/actions/workflows/ci.yml/badge.svg?branch=main)

## Introduction

TL;DR:
The Go team ultimately did not implement the iterator utility due to various considerations, which led to the creation of this package.

----

About two years ago, rsc (Russ Cox) proposed [x/exp/xiter: new package with iterator adapters](https://github.com/golang/go/issues/61898), which defined a new package but was ultimately declined.

This was not the first attempt by the Go community to address iterator patterns:
- discussion: standard iterator interface in go [golang/go#61898](https://github.com/golang/go/issues/61898)
- user-defined iteration using range over func values (https://github.com/golang/go/discussions/56413)

These proposals introduced new approaches for native iterator implementation in Go, but were eventually rejected.

The reason can be summarized as below:

- Overabstraction: Many participants felt that the proposed iterator adapters encouraged overabstraction, which could lead to code that is harder to read and maintain.
- Lack of Clear Use Cases: There was a general sentiment that the proposed functions did not have clear, compelling use cases that justified their inclusion in the standard library.

The Go team recommended this functionality be implemented through third-party packages, which is why this project exists.

## Philosophy

After some exploration, xiter found its true north: **lazy evaluation with `iter.Seq`**.

The core principle: **only put things in xiter that benefit from lazy evaluation**. Functions like `Sum`, `Mean`, `Uniq` might *work* with Seq, but they don't benefit from it — a plain slice version is more intuitive and faster. Those belong in companion packages like `xslice`/`xmap`.

**xiter shines where lazy evaluation matters:**

- **Sources** — `FromSlice`, `FromChan`, `Range`, `Cycle`, `Generate`
- **Lazy transformers** — `Map`, `Filter`, `Chunk`, `WithIndex` (no intermediate allocations)
- **Short-circuit consumers** — `All`, `Any`, `Find`, `First` (don't read the whole sequence)
- **Combinators** — `Concat`, `Zip`, `Merge`, `Equal`
- **I/O streams** — `Lines`, `ReadFileByChunk` (process files without loading into memory)
- **Concurrent processing** — `ParallelMap`, `FanIn` (orchestrate goroutines)
- **Time-based** — `Ticker`, `Throttle`, `Debounce` (event streams as sequences)
- **Combinatorial** — `Combinations`, `Permutations`, `Product` (explosive output, lazy iteration)
- **Windowing & grouping** — `Window`, `SortedGroupBy` (streaming over batches)

## Sub-packages

| Package | Description | Requires Go 1.23+ |
|---------|-------------|-------------------|
| `xiter` | Core iterator utilities (Map, Filter, Concat, Zip, etc.) | Yes |
| `xiter/io` | Lazy I/O — `Lines`, `ReadDir`, `ReadFileByChunk` | Yes |
| `xiter/stream` | Concurrent stream processing — `ParallelMap`, `FanIn`, `Batch` | Yes |
| `xiter/rt` | Time-based sequences — `Ticker`, `Throttle`, `Debounce` | Yes |
| `xiter/collect` | Windowing & grouping — `Window`, `GroupBy`, `SortedGroupBy` | Yes |
| `xiter/combin` | Combinatorial generators — `Combinations`, `Permutations`, `Product` | Yes |
| `xsync` | Concurrency-safe wrappers for sync.Map, sync.Pool, sync.Mutex | No |
| `xslice` | Slice utilities (CountBy, KeyBy, Partition, Chunk, etc.) | No |
| `xmap` | Map utilities (Merge, Filter, Difference, etc.) | No |
| `xstl/list` | Doubly-linked list (generics port of container/list) | No |
| `xstl/lockedmap` | Concurrency-safe map with sync.RWMutex | No |
| `optional` | Optional value type | No |
| `xcmp` | Comparison utilities | No |
| `union` | Algebraic union types | No |

## Target

This package is designed to implement iterator utilities which are defined in [proposal: x/exp/xiter: new package with iterator adapters](https://github.com/golang/go/issues/61898).

And this package will also provide some generic utilities with iter.

This package considers Go versions before 1.23, starting from Go 1.18+

## Docs

[Full Docs](./doc/doc.md)

Docs For Packages:

- [xcmp](./xcmp/README.md)
- [optional](./optional/README.md)
- [union](./union/README.md)
- [xiter](./xiter/README.md)
- [xslice](./xslice/README.md)
- [xmap](./xmap/README.md)

## Contribution

- Make an issue to tell us what you want.
- Create pull request and link to issue.
