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

#### `chemap` - Map Data Structures
- **Multimap** - One-to-many key-value relationships with O(1) operations

#### `chelinkedlist` - LinkedList
- Generic singly linked list with O(1) prepend/append
- Iterator support and common list operations

#### `chedoublylinkedlist` - DoublyLinkedList
- Generic doubly linked list with O(1) prepend/append/remove
- Bidirectional traversal support

#### `chebst` - Binary Search Tree
- Generic BST with O(log n) average operations
- In-order, pre-order, post-order traversals
- Min, max, height operations

### Utilities

#### `cheslice` - Slice Functions
**Basic Operations:**
- Union, Diff, Intersect, Unique
- Map, Filter, ForEach
- Chunk, Fill, Flatten

**Advanced Functions:**
- Reduce, GroupBy, Partition
- Take, Drop, TakeWhile, DropWhile
- Any, All, None
- Reverse, Find, FindIndex, Count
- Zip (combine two slices)

#### `chemap` - Map Functions
- Keys, Values extraction
- Invert (swap keys and values)
- Filter, MapValues (transform values)
- Merge (combine maps)
- Pick, Omit (select/exclude keys)

#### `chetest` - Testing Helpers
- RequireEqual with custom messages
- Assertion utilities for tests

## Features Status

- [x] Slice functions: Unique, Diff, Intersect, Map, Filter, Reduce, GroupBy, Partition, etc.
- [x] Map functions: Keys, Values, Invert, Filter, MapValues, Merge, Pick, Omit
- [x] Data structures: HashSet, OrderedSet, Queue, Stack, Multimap
- [x] Linked data structures: LinkedList, DoublyLinkedList
- [x] Tree structures: Binary Search Tree
- [ ] More data structures: LRU Cache, Trie, AVL Tree, etc.
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

### Multimap
```go
mm := chemap.NewMultimap[string, int]()
mm.Put("key", 1)
mm.Put("key", 2)
values := mm.Get("key") // [1, 2]
```

### LinkedList
```go
ll := chelinkedlist.New[int]()
ll.Append(1)
ll.Append(2)
ll.Prepend(0) // [0, 1, 2]
```

### DoublyLinkedList
```go
dll := chedoublylinkedlist.New[int]()
dll.Append(1)
dll.Append(2)
dll.RemoveLast() // O(1) removal from both ends
```

### Binary Search Tree
```go
bst := chebst.NewInt()
bst.Insert(5)
bst.Insert(3)
bst.Insert(7)
sorted := bst.InOrder() // [3, 5, 7]
```

### Slice Algorithms
```go
// Reduce
sum := cheslice.Reduce([]int{1, 2, 3}, 0, func(acc, n int) int { return acc + n }) // 6

// GroupBy
groups := cheslice.GroupBy([]int{1, 2, 3, 4}, func(n int) string {
    if n%2 == 0 { return "even" } else { return "odd" }
}) // map[even:[2,4] odd:[1,3]]

// Partition
evens, odds := cheslice.Partition([]int{1, 2, 3, 4}, func(n int) bool { return n%2 == 0 })
```

### Map Algorithms
```go
// Values
m := map[string]int{"a": 1, "b": 2}
values := chemap.Values(m) // [1, 2]

// Invert
inverted := chemap.Invert(m) // map[1:"a" 2:"b"]

// MapValues
doubled := chemap.MapValues(m, func(v int) int { return v * 2 }) // map[a:2 b:4]
```

## Credits

* [gopherize.me](https://gopherize.me/): For the excellent gopher image!