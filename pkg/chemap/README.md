# chemap - Map Utilities and Multimap

A collection of map utilities and the Multimap data structure for Go with zero dependencies and 100% test coverage.

## Features

- ✅ **Fully Generic** - Works with any comparable key type
- ✅ **Zero Dependencies** - Pure Go standard library
- ✅ **100% Test Coverage** - Comprehensive test suite
- ✅ **Production Ready** - Battle-tested implementations

## Components

### Map Utility Functions

#### `Keys()`
Extracts all keys from a map as a slice.

```go
m := map[string]int{"a": 1, "b": 2, "c": 3}
keys := chemap.Keys(m) // []string{"a", "b", "c"} (order not guaranteed)
```

### Multimap Data Structure

A map where each key can have multiple values. Perfect for representing one-to-many relationships.

## Installation

```bash
go get github.com/comfortablynumb/che/pkg/chemap
```

## Multimap Quick Start

```go
package main

import (
    "fmt"
    "github.com/comfortablynumb/che/pkg/chemap"
)

func main() {
    mm := chemap.NewMultimap[string, int]()

    // Add multiple values for same key
    mm.Put("numbers", 1)
    mm.Put("numbers", 2)
    mm.Put("numbers", 3)

    // Get all values
    values := mm.Get("numbers") // [1, 2, 3]
    fmt.Println(values)

    // Check key existence
    exists := mm.ContainsKey("numbers") // true

    // Get count
    count := mm.ValueCount("numbers") // 3
}
```

## Multimap API

### Creating Multimaps

```go
// Create empty multimap
mm := chemap.NewMultimap[string, int]()

// Create with initial capacity
mm := chemap.NewMultimapWithCapacity[string, int](100)
```

### Adding Values

```go
mm := chemap.NewMultimap[string, string]()

// Add single value
mm.Put("fruits", "apple")

// Add multiple values at once
mm.PutAll("fruits", "banana", "cherry", "date")

// Add to existing key
mm.Put("fruits", "elderberry") // Now has 5 values
```

### Retrieving Values

```go
mm := chemap.NewMultimap[string, int]()
mm.PutAll("scores", 85, 92, 78, 95)

// Get all values (returns a copy)
values := mm.Get("scores") // [85, 92, 78, 95]

// Get first value only
first, ok := mm.GetFirst("scores") // 85, true

// Get for non-existent key
empty := mm.Get("nonexistent") // []int{} (empty slice)
```

### Checking Existence

```go
mm := chemap.NewMultimap[string, int]()
mm.PutAll("data", 1, 2, 3)

// Check if key exists
hasKey := mm.ContainsKey("data") // true

// Check if specific entry exists
equals := func(a, b int) bool { return a == b }
hasEntry := mm.ContainsEntry("data", 2, equals) // true
```

### Removing Values

```go
mm := chemap.NewMultimap[string, int]()
mm.PutAll("numbers", 1, 2, 3, 4, 5)

// Remove specific value
equals := func(a, b int) bool { return a == b }
removed := mm.Remove("numbers", 3, equals) // true
// Now has [1, 2, 4, 5]

// Remove all values for a key
mm.RemoveAll("numbers") // Removes entire key
```

### Replacing Values

```go
mm := chemap.NewMultimap[string, string]()
mm.PutAll("tags", "old1", "old2", "old3")

// Replace all values for a key
mm.ReplaceValues("tags", "new1", "new2")
// Now has ["new1", "new2"]

// Remove key by replacing with empty
mm.ReplaceValues("tags") // Key removed
```

## Size and Counts

```go
mm := chemap.NewMultimap[string, int]()
mm.PutAll("a", 1, 2, 3)
mm.PutAll("b", 4, 5)

// Total number of values
total := mm.Size() // 5

// Number of unique keys
keys := mm.KeyCount() // 2

// Number of values for specific key
count := mm.ValueCount("a") // 3

// Check if empty
empty := mm.IsEmpty() // false
```

## Iteration

### Iterate Over All Entries

```go
mm := chemap.NewMultimap[string, int]()
mm.PutAll("a", 1, 2)
mm.Put("b", 3)

// Iterate over each key-value pair
mm.ForEach(func(key string, value int) bool {
    fmt.Printf("%s: %d\n", key, value)
    return true // continue iteration
})
// Output:
// a: 1
// a: 2
// b: 3
```

### Iterate Over Keys

```go
mm := chemap.NewMultimap[string, int]()
mm.PutAll("evens", 2, 4, 6)
mm.PutAll("odds", 1, 3, 5)

// Iterate over each key and its values
mm.ForEachKey(func(key string, values []int) bool {
    fmt.Printf("%s: %v\n", key, values)
    return true // continue iteration
})
```

## Advanced Operations

### Cloning

```go
original := chemap.NewMultimap[string, int]()
original.PutAll("data", 1, 2, 3)

// Create independent copy
clone := original.Clone()
clone.Put("data", 4) // Doesn't affect original
```

### Merging

```go
mm1 := chemap.NewMultimap[string, int]()
mm2 := chemap.NewMultimap[string, int]()

mm1.PutAll("shared", 1, 2)
mm2.PutAll("shared", 3, 4)
mm2.Put("unique", 5)

// Merge mm2 into mm1
mm1.Merge(mm2)
// mm1 now has: shared:[1,2,3,4], unique:[5]
```

### Get All Keys and Values

```go
mm := chemap.NewMultimap[string, int]()
mm.PutAll("a", 1, 2)
mm.PutAll("b", 3, 4)

// Get all keys
keys := mm.Keys() // ["a", "b"]

// Get all values (flattened)
values := mm.Values() // [1, 2, 3, 4]
```

## Use Cases

### 1. HTTP Headers

```go
headers := chemap.NewMultimap[string, string]()

// HTTP headers can have multiple values
headers.Put("Accept", "application/json")
headers.Put("Accept", "text/html")
headers.Put("Set-Cookie", "session=abc123")
headers.Put("Set-Cookie", "user=john")

// Get all Accept headers
accepts := headers.Get("Accept")
```

### 2. Inverted Index

```go
// Word -> Document IDs
index := chemap.NewMultimap[string, int]()

// Document 1 contains: "go", "programming"
index.Put("go", 1)
index.Put("programming", 1)

// Document 2 contains: "go", "language"
index.Put("go", 2)
index.Put("language", 2)

// Find all documents containing "go"
docs := index.Get("go") // [1, 2]
```

### 3. Graph Representation (Adjacency List)

```go
// Node -> Adjacent nodes
graph := chemap.NewMultimap[int, int]()

// Node 1 connects to 2 and 3
graph.PutAll(1, 2, 3)

// Node 2 connects to 3 and 4
graph.PutAll(2, 3, 4)

// Get neighbors of node 1
neighbors := graph.Get(1)
```

### 4. Form Data / Query Parameters

```go
// Query parameter -> values
params := chemap.NewMultimap[string, string]()

// URL: ?tag=go&tag=programming&category=tech
params.Put("tag", "go")
params.Put("tag", "programming")
params.Put("category", "tech")

// Get all tags
tags := params.Get("tag") // ["go", "programming"]
```

### 5. Event Listeners / Observers

```go
type Handler func(event interface{})

listeners := chemap.NewMultimap[string, Handler]()

// Multiple listeners for same event
listeners.Put("click", handleClick1)
listeners.Put("click", handleClick2)
listeners.Put("submit", handleSubmit)

// Trigger all click handlers
handlers := listeners.Get("click")
for _, handler := range handlers {
    handler(clickEvent)
}
```

### 6. Grouping / Categorization

```go
// Category -> Items
catalog := chemap.NewMultimap[string, string]()

catalog.Put("fruits", "apple")
catalog.Put("fruits", "banana")
catalog.Put("vegetables", "carrot")
catalog.Put("vegetables", "broccoli")

// Get all fruits
fruits := catalog.Get("fruits")
```

### 7. Role-Based Access Control

```go
// User -> Roles
userRoles := chemap.NewMultimap[string, string]()

userRoles.PutAll("john", "admin", "developer")
userRoles.PutAll("jane", "developer", "reviewer")
userRoles.Put("bob", "viewer")

// Check if user has admin role
equals := func(a, b string) bool { return a == b }
isAdmin := userRoles.ContainsEntry("john", "admin", equals)
```

## Performance Characteristics

| Operation | Time Complexity | Notes |
|-----------|----------------|-------|
| Put | O(1) amortized | |
| PutAll | O(k) | k = number of values |
| Get | O(n) | n = values for key (copy) |
| GetFirst | O(1) | |
| ContainsKey | O(1) | |
| ContainsEntry | O(n) | n = values for key |
| Remove | O(n) | n = values for key |
| RemoveAll | O(1) | |
| Size | O(k) | k = number of keys |
| KeyCount | O(1) | |
| ValueCount | O(1) | |
| ForEach | O(n) | n = total entries |
| Clone | O(n) | n = total entries |
| Merge | O(m) | m = entries in other |

**Space Complexity**: O(n) where n is the total number of key-value pairs

## Thread Safety

**Multimap is not thread-safe.** For concurrent use, provide external synchronization:

```go
import "sync"

var (
    mm = chemap.NewMultimap[string, int]()
    mu sync.RWMutex
)

// For writes
mu.Lock()
mm.Put("key", 1)
mu.Unlock()

// For reads
mu.RLock()
values := mm.Get("key")
mu.RUnlock()
```

## Comparison with Alternatives

### vs. map[K][]V

**Multimap advantages:**
- ✅ Cleaner, more intuitive API
- ✅ Safe Get operations (returns copies)
- ✅ Rich utility methods
- ✅ Better encapsulation

### vs. Multiple Maps

**Multimap advantages:**
- ✅ Single data structure to manage
- ✅ Atomic operations on related data
- ✅ Easier to reason about
- ✅ Built-in size tracking

## Design Decisions

### Why Get Returns Copies?

The `Get()` method returns a copy of the values slice to prevent external code from modifying the internal state. This ensures:
- Data integrity
- Thread-safe reads (when combined with proper locking)
- Predictable behavior

### Why Requires Equality Function?

Operations like `ContainsEntry()` and `Remove()` require an equality function to support:
- Custom comparison logic
- Pointer comparisons
- Struct field comparisons
- Complex types

```go
// Simple comparison
equals := func(a, b int) bool { return a == b }

// Custom struct comparison
type User struct { ID int; Name string }
equalsUser := func(a, b User) bool { return a.ID == b.ID }
```

## Contributing

Contributions are welcome! Please ensure:
- All tests pass
- Coverage remains at 100%
- Code follows Go best practices

## Related Packages

- `cheset` - Set implementations (HashSet, OrderedSet)
- `chequeue` - Queue (FIFO) implementation
- `chestack` - Stack (LIFO) implementation
- `cheslice` - Slice utility functions
