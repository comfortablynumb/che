# chestack - Generic Stack Implementation

A production-ready, generic Stack (LIFO) implementation for Go with zero dependencies and 100% test coverage.

## Features

- ✅ **Fully Generic** - Works with any type
- ✅ **Zero Dependencies** - Pure Go standard library
- ✅ **O(1) Amortized Performance** - Constant time for Push and Pop
- ✅ **100% Test Coverage** - Comprehensive test suite
- ✅ **Automatic Growth** - Grows as needed
- ✅ **Simple API** - Clean, intuitive interface

## Installation

```bash
go get github.com/comfortablynumb/che/pkg/chestack
```

## Quick Start

```go
package main

import (
    "fmt"
    "github.com/comfortablynumb/che/pkg/chestack"
)

func main() {
    // Create a new stack
    stack := chestack.New[int]()

    // Push elements
    stack.Push(1)
    stack.Push(2)
    stack.Push(3)

    // Pop elements (LIFO order)
    value, ok := stack.Pop() // 3, true
    fmt.Println(value)

    // Peek at top without removing
    top, ok := stack.Peek() // 2, true

    // Check size
    fmt.Println(stack.Size()) // 2
}
```

## Core Operations

### Creating Stacks

```go
// Create empty stack
stack := chestack.New[string]()

// Create with initial capacity
stack := chestack.NewWithCapacity[string](100)

// Create from slice
stack := chestack.NewFromSlice([]string{"a", "b", "c"})
// Note: Last element in slice ("c") will be on top
```

### Basic Operations

```go
stack := chestack.New[int]()

// Push (add to top)
stack.Push(1)
stack.Push(2)

// Push multiple elements
stack.PushMultiple(3, 4, 5)

// Pop (remove from top)
value, ok := stack.Pop() // 5, true

// Peek at top element
top, ok := stack.Peek() // 4, true

// Check size
size := stack.Size() // 4

// Check if empty
empty := stack.IsEmpty() // false

// Clear all elements
stack.Clear()
```

## Advanced Features

### Iteration

```go
stack := chestack.NewFromSlice([]int{1, 2, 3, 4, 5})

// Iterate from bottom to top
stack.ForEach(func(item int) bool {
    fmt.Println(item) // 1, 2, 3, 4, 5
    return true // return false to stop iteration
})

// Iterate from top to bottom
stack.ForEachReverse(func(item int) bool {
    fmt.Println(item) // 5, 4, 3, 2, 1
    return true
})
```

### Cloning

```go
original := chestack.NewFromSlice([]int{1, 2, 3})
clone := original.Clone()

// Modifications to clone don't affect original
clone.Push(4)
```

### Conversion

```go
stack := chestack.NewFromSlice([]int{1, 2, 3})

// Convert to slice (bottom to top)
slice := stack.ToSlice() // [1, 2, 3]
```

### Contains (with custom equality)

```go
stack := chestack.NewFromSlice([]int{1, 2, 3})

equals := func(a, b int) bool { return a == b }
exists := stack.Contains(2, equals) // true
```

## Use Cases

### 1. Expression Evaluation

```go
func evaluatePostfix(expr []string) int {
    stack := chestack.New[int]()

    for _, token := range expr {
        if isOperator(token) {
            b, _ := stack.Pop()
            a, _ := stack.Pop()
            result := apply(token, a, b)
            stack.Push(result)
        } else {
            stack.Push(parseInt(token))
        }
    }

    result, _ := stack.Pop()
    return result
}
```

### 2. Undo/Redo Functionality

```go
type Action struct {
    Type string
    Data interface{}
}

undoStack := chestack.New[Action]()
redoStack := chestack.New[Action]()

func doAction(action Action) {
    // Perform action
    performAction(action)

    // Save to undo stack
    undoStack.Push(action)

    // Clear redo stack
    redoStack.Clear()
}

func undo() {
    if action, ok := undoStack.Pop(); ok {
        // Undo the action
        revertAction(action)

        // Move to redo stack
        redoStack.Push(action)
    }
}

func redo() {
    if action, ok := redoStack.Pop(); ok {
        // Redo the action
        performAction(action)

        // Move back to undo stack
        undoStack.Push(action)
    }
}
```

### 3. DFS (Depth-First Search)

```go
type Node struct {
    Value int
    Children []*Node
}

func DFS(root *Node) {
    stack := chestack.New[*Node]()
    stack.Push(root)

    for !stack.IsEmpty() {
        node, _ := stack.Pop()
        fmt.Println(node.Value)

        // Push children in reverse order for correct DFS order
        for i := len(node.Children) - 1; i >= 0; i-- {
            stack.Push(node.Children[i])
        }
    }
}
```

### 4. Parentheses/Bracket Matching

```go
func isBalanced(s string) bool {
    stack := chestack.New[rune]()
    pairs := map[rune]rune{')': '(', ']': '[', '}': '{'}

    for _, char := range s {
        if char == '(' || char == '[' || char == '{' {
            stack.Push(char)
        } else if opening, isClosing := pairs[char]; isClosing {
            if top, ok := stack.Pop(); !ok || top != opening {
                return false
            }
        }
    }

    return stack.IsEmpty()
}
```

### 5. Browser History

```go
type Page struct {
    URL   string
    Title string
}

backStack := chestack.New[Page]()
forwardStack := chestack.New[Page]()
currentPage := Page{URL: "home", Title: "Home"}

func navigate(page Page) {
    backStack.Push(currentPage)
    currentPage = page
    forwardStack.Clear()
}

func goBack() {
    if page, ok := backStack.Pop(); ok {
        forwardStack.Push(currentPage)
        currentPage = page
    }
}

func goForward() {
    if page, ok := forwardStack.Pop(); ok {
        backStack.Push(currentPage)
        currentPage = page
    }
}
```

### 6. Function Call Stack Simulation

```go
type CallFrame struct {
    FunctionName string
    Arguments    []interface{}
    LocalVars    map[string]interface{}
}

callStack := chestack.New[CallFrame]()

func callFunction(name string, args ...interface{}) {
    frame := CallFrame{
        FunctionName: name,
        Arguments:    args,
        LocalVars:    make(map[string]interface{}),
    }
    callStack.Push(frame)

    // Execute function...

    // When function returns
    callStack.Pop()
}

func printStackTrace() {
    callStack.ForEachReverse(func(frame CallFrame) bool {
        fmt.Printf("  at %s\n", frame.FunctionName)
        return true
    })
}
```

### 7. Backtracking Algorithms

```go
func solveMaze(maze [][]int, start, end Point) []Point {
    stack := chestack.New[Point]()
    visited := make(map[Point]bool)

    stack.Push(start)

    for !stack.IsEmpty() {
        current, _ := stack.Pop()

        if current == end {
            return reconstructPath(current)
        }

        if visited[current] {
            continue
        }
        visited[current] = true

        // Push unvisited neighbors
        for _, neighbor := range getNeighbors(maze, current) {
            if !visited[neighbor] {
                stack.Push(neighbor)
            }
        }
    }

    return nil // No path found
}
```

## Performance Characteristics

| Operation | Time Complexity | Notes |
|-----------|----------------|-------|
| Push | O(1) amortized | May resize occasionally |
| Pop | O(1) | |
| Peek | O(1) | |
| Size | O(1) | |
| IsEmpty | O(1) | |
| Clear | O(1) | |
| ToSlice | O(n) | |
| Clone | O(n) | |
| ForEach | O(n) | |
| ForEachReverse | O(n) | |
| Contains | O(n) | Linear search |

**Space Complexity**: O(n) where n is the number of elements

### Memory Management

- Initial capacity: 8 elements
- Growth: Automatically grows using Go's slice growth strategy
- Efficient amortized performance
- No automatic shrinking (call Clear() to reset)

## Thread Safety

**Stack is not thread-safe.** For concurrent use, provide external synchronization:

```go
import "sync"

var (
    stack = chestack.New[int]()
    mu    sync.Mutex
)

// For writes
mu.Lock()
stack.Push(1)
mu.Unlock()

// For reads
mu.Lock()
value, ok := stack.Pop()
mu.Unlock()
```

## Comparison with Alternatives

### vs. Slice

**Stack advantages:**
- ✅ Clear LIFO semantics
- ✅ Type-safe operations
- ✅ Intuitive API (Push/Pop vs append/slice)
- ✅ Additional utilities (Peek, ForEachReverse, etc.)

### vs. container/list

**Stack advantages:**
- ✅ Better cache locality (contiguous memory)
- ✅ Simpler API
- ✅ Fully generic (not interface{})
- ✅ Better performance

## LIFO Guarantee

Stack strictly maintains Last-In-First-Out order:

```go
stack := chestack.New[int]()
stack.PushMultiple(1, 2, 3, 4, 5)

// Elements come out in reverse order
val1, _ := stack.Pop() // 5
val2, _ := stack.Pop() // 4
val3, _ := stack.Pop() // 3
// ...
```

## String Representation

```go
stack := chestack.New[int]()
stack.PushMultiple(1, 2, 3)

fmt.Println(stack.String()) // "Stack[1, 2, 3]"
// Elements shown from bottom to top, top element on right
```

## Contributing

Contributions are welcome! Please ensure:
- All tests pass
- Coverage remains at 100%
- Code follows Go best practices

## Related Packages

- `chequeue` - Queue (FIFO) implementation
- `cheset` - Set implementations (HashSet, OrderedSet)
- `cheslice` - Slice utility functions
