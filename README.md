# Che - An Extension for the Go Standard Library

![Go Report Card](https://goreportcard.com/badge/github.com/comfortablynumb/che)
![Build Status](https://github.com/comfortablynumb/che/actions/workflows/build.yml/badge.svg)
![License](https://img.shields.io/github/license/comfortablynumb/che)

<p align="center">
  <img alt="Che!" width="450" height="482" src="https://github.com/comfortablynumb/che/raw/main/docs/images/gopher.png" />
</p>

---

**:construction_worker: IMPORTANT NOTE: :construction_worker: This is a work in progress. Stay tuned!**

---

## Introduction

This library aims to meet the following requirements:

* It must have all the functions and data structures that we use in our everyday tasks, but that are not present in Golang's standard library
* It must have **zero dependencies**
* It must have **high test coverage**
* It must be **fully generic** using Go 1.20+ generics

## Packages

### Data Structures

#### `cheset` - Set Implementations
- **HashSet** - Unordered set with O(1) operations
- **OrderedSet** - Insertion-ordered set with O(1) lookups

#### `chequeue` - Queue (FIFO)
- Generic Queue implementation with O(1) amortized operations
- Circular buffer for efficient memory usage

#### `chestack` - Stack (LIFO)
- Generic Stack implementation with O(1) amortized operations
- Simple slice-based implementation

### Utilities

#### `cheslice` - Slice Functions
- Union, Diff, Intersect, Unique
- Map, Filter, ForEach
- Chunk, Fill, and more

#### `chemap` - Map Functions
- Keys extraction
- More utilities coming soon

#### `chetest` - Testing Helpers
- RequireEqual with custom messages
- Assertion utilities for tests

## Features Status

- [x] Slice functions: Unique, Diff, Intersect, Map, Filter, etc.
- [x] Map functions: Keys extraction
- [x] Data structures: HashSet, OrderedSet, Queue, Stack
- [ ] More data structures: Multimap, LRU Cache, Trie, etc.
- [ ] File handling functions
- [ ] Networking functions

## Quick Examples

### HashSet
```go
set := cheset.New[int]()
set.Add(1)
set.Add(2)
set.Contains(1) // true
```

### OrderedSet
```go
set := cheset.NewOrdered[string]()
set.Add("first")
set.Add("second")
set.GetAt(0) // "first" - maintains order
```

### Queue
```go
queue := chequeue.New[int]()
queue.Enqueue(1)
queue.Enqueue(2)
value, _ := queue.Dequeue() // 1 (FIFO)
```

### Stack
```go
stack := chestack.New[int]()
stack.Push(1)
stack.Push(2)
value, _ := stack.Pop() // 2 (LIFO)
```

## Credits

* [gopherize.me](https://gopherize.me/): For the excellent gopher image!