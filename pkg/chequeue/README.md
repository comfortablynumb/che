# chequeue - Generic Queue Implementation

A production-ready, generic Queue (FIFO) implementation for Go with zero dependencies and near-100% test coverage.

## Features

- ✅ **Fully Generic** - Works with any type
- ✅ **Zero Dependencies** - Pure Go standard library
- ✅ **O(1) Amortized Performance** - Constant time for Enqueue and Dequeue
- ✅ **Automatic Resizing** - Grows and shrinks as needed
- ✅ **100% Test Coverage** - Comprehensive test suite
- ✅ **Circular Buffer** - Efficient memory usage

## Installation

```bash
go get github.com/comfortablynumb/che/pkg/chequeue
```

## Quick Start

```go
package main

import (
    "fmt"
    "github.com/comfortablynumb/che/pkg/chequeue"
)

func main() {
    // Create a new queue
    queue := chequeue.New[int]()

    // Enqueue elements
    queue.Enqueue(1)
    queue.Enqueue(2)
    queue.Enqueue(3)

    // Dequeue elements (FIFO order)
    value, ok := queue.Dequeue() // 1, true
    fmt.Println(value)

    // Peek at front without removing
    front, ok := queue.Peek() // 2, true

    // Check size
    fmt.Println(queue.Size()) // 2
}
```

## Core Operations

### Creating Queues

```go
// Create empty queue
queue := chequeue.New[string]()

// Create with initial capacity
queue := chequeue.NewWithCapacity[string](100)

// Create from slice
queue := chequeue.NewFromSlice([]string{"a", "b", "c"})
```

### Basic Operations

```go
queue := chequeue.New[int]()

// Enqueue (add to back)
queue.Enqueue(1)
queue.Enqueue(2)

// Enqueue multiple elements
queue.EnqueueMultiple(3, 4, 5)

// Dequeue (remove from front)
value, ok := queue.Dequeue() // 1, true

// Peek at front element
front, ok := queue.Peek() // 2, true

// Check size
size := queue.Size() // 4

// Check if empty
empty := queue.IsEmpty() // false

// Clear all elements
queue.Clear()
```

## Advanced Features

### Iteration

```go
queue := chequeue.NewFromSlice([]int{1, 2, 3, 4, 5})

// Iterate over all elements (FIFO order)
queue.ForEach(func(item int) bool {
    fmt.Println(item)
    return true // return false to stop iteration
})
```

### Cloning

```go
original := chequeue.NewFromSlice([]int{1, 2, 3})
clone := original.Clone()

// Modifications to clone don't affect original
clone.Enqueue(4)
```

### Conversion

```go
queue := chequeue.NewFromSlice([]int{1, 2, 3})

// Convert to slice (FIFO order)
slice := queue.ToSlice() // [1, 2, 3]
```

### Contains (with custom equality)

```go
queue := chequeue.NewFromSlice([]int{1, 2, 3})

equals := func(a, b int) bool { return a == b }
exists := queue.Contains(2, equals) // true
```

## Use Cases

### 1. Task Queue

```go
type Task struct {
    ID   int
    Name string
}

taskQueue := chequeue.New[Task]()

// Producer adds tasks
taskQueue.Enqueue(Task{ID: 1, Name: "Process data"})
taskQueue.Enqueue(Task{ID: 2, Name: "Send email"})

// Consumer processes tasks
for !taskQueue.IsEmpty() {
    task, _ := taskQueue.Dequeue()
    processTask(task)
}
```

### 2. BFS (Breadth-First Search)

```go
type Node struct {
    Value int
    Children []*Node
}

func BFS(root *Node) {
    queue := chequeue.New[*Node]()
    queue.Enqueue(root)

    for !queue.IsEmpty() {
        node, _ := queue.Dequeue()
        fmt.Println(node.Value)

        for _, child := range node.Children {
            queue.Enqueue(child)
        }
    }
}
```

### 3. Request Buffer

```go
type Request struct {
    ID      string
    Payload []byte
}

requestBuffer := chequeue.NewWithCapacity[Request](1000)

// Buffer incoming requests
func handleRequest(req Request) {
    requestBuffer.Enqueue(req)
}

// Process buffered requests
func processRequests() {
    for !requestBuffer.IsEmpty() {
        req, _ := requestBuffer.Dequeue()
        handleRequestLogic(req)
    }
}
```

### 4. Event Queue

```go
type Event struct {
    Type      string
    Timestamp int64
    Data      interface{}
}

eventQueue := chequeue.New[Event]()

// Enqueue events as they occur
eventQueue.Enqueue(Event{Type: "click", Timestamp: time.Now().Unix()})

// Process events in order
for !eventQueue.IsEmpty() {
    event, _ := eventQueue.Dequeue()
    dispatchEvent(event)
}
```

### 5. Round-Robin Scheduler

```go
type Worker struct {
    ID   int
    Load int
}

workers := chequeue.NewFromSlice([]Worker{
    {ID: 1, Load: 0},
    {ID: 2, Load: 0},
    {ID: 3, Load: 0},
})

func assignTask() Worker {
    // Get next worker
    worker, _ := workers.Dequeue()

    // Assign task to worker
    worker.Load++

    // Put worker back at end of queue
    workers.Enqueue(worker)

    return worker
}
```

## Performance Characteristics

| Operation | Time Complexity | Notes |
|-----------|----------------|-------|
| Enqueue | O(1) amortized | May resize occasionally |
| Dequeue | O(1) amortized | May shrink occasionally |
| Peek | O(1) | |
| Size | O(1) | |
| IsEmpty | O(1) | |
| Clear | O(1) | |
| ToSlice | O(n) | |
| Clone | O(n) | |
| ForEach | O(n) | |
| Contains | O(n) | Linear search |

**Space Complexity**: O(n) where n is the number of elements

### Memory Management

- Initial capacity: 8 elements
- Growth: Doubles when full
- Shrink: Halves when less than 25% full (minimum capacity: 8)
- Circular buffer implementation for efficient memory usage

## Thread Safety

**Queue is not thread-safe.** For concurrent use, provide external synchronization:

```go
import "sync"

var (
    queue = chequeue.New[int]()
    mu    sync.Mutex
)

// For writes
mu.Lock()
queue.Enqueue(1)
mu.Unlock()

// For reads
mu.Lock()
value, ok := queue.Dequeue()
mu.Unlock()
```

## Comparison with Alternatives

### vs. Slice

**Queue advantages:**
- ✅ O(1) dequeue (slice requires O(n) for removing first element)
- ✅ Automatic memory management
- ✅ Circular buffer efficiency
- ✅ Clear FIFO semantics

### vs. container/list

**Queue advantages:**
- ✅ Better cache locality (contiguous memory)
- ✅ Simpler, more intuitive API
- ✅ Fully generic (not interface{})
- ✅ Better performance for typical use cases

### vs. Channel

**Queue advantages:**
- ✅ Non-blocking operations
- ✅ Can inspect size without blocking
- ✅ Can peek without removing
- ✅ No goroutine required

**Channel advantages:**
- ✅ Thread-safe by default
- ✅ Blocking/buffered semantics
- ✅ Select statement support

## FIFO Guarantee

Queue strictly maintains First-In-First-Out order:

```go
queue := chequeue.New[int]()
queue.EnqueueMultiple(1, 2, 3, 4, 5)

// Elements come out in exact order they went in
val1, _ := queue.Dequeue() // 1
val2, _ := queue.Dequeue() // 2
val3, _ := queue.Dequeue() // 3
// ...
```

## Contributing

Contributions are welcome! Please ensure:
- All tests pass
- Coverage remains high
- Code follows Go best practices

## Related Packages

- `chestack` - Stack (LIFO) implementation
- `cheset` - Set implementations (HashSet, OrderedSet)
- `cheslice` - Slice utility functions
