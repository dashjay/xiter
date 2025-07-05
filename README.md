# xiter
A Generic Library for Go with the Power of Iterators

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

## Target
This package is designed to implement iterator utilities which are defined in [proposal: x/exp/xiter: new package with iterator adapters](https://github.com/golang/go/issues/61898).

And this package will also provide some generic utilities with iter.

This package considers Go versions before 1.23, starting from Go 1.18+


## Contribution

- Make an issue to tell us what you want.
- Create pull request and link to issue.

