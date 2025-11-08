# cheset - Generic HashSet Implementation

A production-ready, generic HashSet implementation for Go with zero dependencies and 100% test coverage.

## Features

- ✅ **Fully Generic** - Works with any comparable type
- ✅ **Zero Dependencies** - Pure Go standard library
- ✅ **O(1) Performance** - Average-case constant time for core operations
- ✅ **100% Test Coverage** - Comprehensive test suite with 49+ test cases
- ✅ **Rich API** - Complete set operations and utilities
- ✅ **Well Documented** - Extensive godoc comments and examples

## Installation

```bash
go get github.com/comfortablynumb/che/pkg/cheset
```

## Quick Start

```go
package main

import (
    "fmt"
    "github.com/comfortablynumb/che/pkg/cheset"
)

func main() {
    // Create a new set
    set := cheset.New[int]()

    // Add elements
    set.Add(1)
    set.Add(2)
    set.Add(3)

    // Check membership
    fmt.Println(set.Contains(2)) // true

    // Get size
    fmt.Println(set.Size()) // 3

    // Remove element
    set.Remove(2)
    fmt.Println(set.Contains(2)) // false
}
```

## Core Operations

### Creating Sets

```go
// Create empty set
set := cheset.New[string]()

// Create with initial capacity
set := cheset.NewWithCapacity[string](100)

// Create from slice
set := cheset.NewFromSlice([]string{"a", "b", "c"})
```

### Basic Operations

```go
set := cheset.New[int]()

// Add single element
added := set.Add(1) // returns true if added, false if already exists

// Add multiple elements
count := set.AddMultiple(1, 2, 3, 4) // returns count of elements added

// Remove elements
removed := set.Remove(1) // returns true if removed, false if didn't exist
count := set.RemoveMultiple(2, 3) // returns count of elements removed

// Check membership
exists := set.Contains(1)
allExist := set.ContainsAll(1, 2, 3)
anyExist := set.ContainsAny(1, 5, 9)

// Get size and check if empty
size := set.Size()
empty := set.IsEmpty()

// Clear all elements
set.Clear()
```

## Set Operations

### Union - Elements in Either Set

```go
set1 := cheset.NewFromSlice([]int{1, 2, 3})
set2 := cheset.NewFromSlice([]int{3, 4, 5})

union := set1.Union(set2)
// Result: {1, 2, 3, 4, 5}
```

### Intersection - Elements in Both Sets

```go
set1 := cheset.NewFromSlice([]int{1, 2, 3, 4})
set2 := cheset.NewFromSlice([]int{3, 4, 5, 6})

intersection := set1.Intersect(set2)
// Result: {3, 4}
```

### Difference - Elements in First Set But Not Second

```go
set1 := cheset.NewFromSlice([]int{1, 2, 3, 4})
set2 := cheset.NewFromSlice([]int{3, 4, 5, 6})

diff := set1.Diff(set2)
// Result: {1, 2}
```

### Symmetric Difference - Elements in Either Set But Not Both

```go
set1 := cheset.NewFromSlice([]int{1, 2, 3, 4})
set2 := cheset.NewFromSlice([]int{3, 4, 5, 6})

symDiff := set1.SymmetricDiff(set2)
// Result: {1, 2, 5, 6}
```

## Set Relations

```go
set1 := cheset.NewFromSlice([]int{1, 2})
set2 := cheset.NewFromSlice([]int{1, 2, 3, 4})
set3 := cheset.NewFromSlice([]int{5, 6, 7})

// Check equality
equal := set1.Equal(set2) // false

// Check subset/superset
isSubset := set1.IsSubset(set2)           // true
isSuperset := set2.IsSuperset(set1)       // true
isProperSubset := set1.IsProperSubset(set2) // true

// Check if sets have no common elements
disjoint := set1.IsDisjoint(set3) // true
```

## Utility Methods

### Clone and Convert

```go
set := cheset.NewFromSlice([]int{1, 2, 3})

// Create a copy
clone := set.Clone()

// Convert to slice (order not guaranteed)
slice := set.ToSlice()
```

### Iteration

```go
set := cheset.NewFromSlice([]int{1, 2, 3, 4, 5})

// Iterate over all elements
set.ForEach(func(item int) bool {
    fmt.Println(item)
    return true // return false to stop iteration
})

// Calculate sum
sum := 0
set.ForEach(func(item int) bool {
    sum += item
    return true
})
```

### Filtering

```go
set := cheset.NewFromSlice([]int{1, 2, 3, 4, 5, 6})

// Filter even numbers
evens := set.Filter(func(item int) bool {
    return item%2 == 0
})
// Result: {2, 4, 6}
```

### String Representation

```go
set := cheset.NewFromSlice([]int{1, 2, 3})
fmt.Println(set.String()) // "HashSet{1, 2, 3}"
```

## Advanced Examples

### Working with Custom Types

```go
type User struct {
    ID   int
    Name string
}

// User must be comparable (no slices, maps, or functions as fields)
users := cheset.New[User]()
users.Add(User{ID: 1, Name: "Alice"})
users.Add(User{ID: 2, Name: "Bob"})

// Check if user exists
exists := users.Contains(User{ID: 1, Name: "Alice"})
```

### Set-Based Deduplication

```go
// Remove duplicates from slice while preserving a set of unique values
data := []string{"apple", "banana", "apple", "cherry", "banana"}
uniqueSet := cheset.NewFromSlice(data)
unique := uniqueSet.ToSlice()
// unique contains: ["apple", "banana", "cherry"] (order not guaranteed)
```

### Finding Common Elements Across Multiple Slices

```go
slice1 := []int{1, 2, 3, 4, 5}
slice2 := []int{3, 4, 5, 6, 7}
slice3 := []int{4, 5, 6, 7, 8}

set1 := cheset.NewFromSlice(slice1)
set2 := cheset.NewFromSlice(slice2)
set3 := cheset.NewFromSlice(slice3)

// Find elements common to all three
common := set1.Intersect(set2).Intersect(set3)
// Result: {4, 5}
```

### Tag/Category Management

```go
// Managing tags or categories
userTags := cheset.NewFromSlice([]string{"golang", "python", "javascript"})
requiredTags := cheset.NewFromSlice([]string{"golang", "docker"})

// Check if user has all required tags
hasAllTags := requiredTags.IsSubset(userTags) // false

// Find missing tags
missingTags := requiredTags.Diff(userTags)
// Result: {"docker"}
```

### Permission System

```go
type Permission string

userPerms := cheset.NewFromSlice([]Permission{
    "read", "write", "execute",
})

adminPerms := cheset.NewFromSlice([]Permission{
    "read", "write", "execute", "admin", "delete",
})

// Check if user is admin (has all admin permissions)
isAdmin := userPerms.Equal(adminPerms) // false

// Find additional permissions needed for admin
additionalPerms := adminPerms.Diff(userPerms)
// Result: {"admin", "delete"}
```

## Performance Characteristics

| Operation | Average Case | Worst Case |
|-----------|--------------|------------|
| Add | O(1) | O(n) |
| Remove | O(1) | O(n) |
| Contains | O(1) | O(n) |
| Size | O(1) | O(1) |
| Union | O(n+m) | O(n+m) |
| Intersect | O(min(n,m)) | O(n*m) |
| Diff | O(n) | O(n*m) |
| Equal | O(n) | O(n) |

Where:
- n = size of first set
- m = size of second set

## Thread Safety

**HashSet is not thread-safe.** For concurrent access, you must provide external synchronization:

```go
import "sync"

var (
    set = cheset.New[int]()
    mu  sync.RWMutex
)

// For writes
mu.Lock()
set.Add(1)
mu.Unlock()

// For reads
mu.RLock()
exists := set.Contains(1)
mu.RUnlock()
```

## Comparison with Alternatives

### vs. map[T]bool

HashSet provides:
- ✅ Cleaner API with set semantics
- ✅ Set operations (Union, Intersect, Diff)
- ✅ Set relations (IsSubset, IsDisjoint, etc.)
- ✅ Better memory efficiency (uses empty struct{})

### vs. map[T]struct{}

HashSet provides:
- ✅ Type-safe operations
- ✅ Rich set operations out of the box
- ✅ Better readability and maintainability
- ✅ Comprehensive testing

## Contributing

Contributions are welcome! Please ensure:
- All tests pass: `go test ./pkg/cheset/...`
- Coverage remains at 100%: `go test ./pkg/cheset/... -coverprofile=coverage.out`
- Code follows Go best practices and conventions

## License

See the LICENSE file in the repository root.

## Related Packages

- `cheslice` - Slice utility functions
- `chemap` - Map utility functions
- `chetest` - Testing utilities

## Why "che"?

The package name follows the repository's naming convention. The library aims to extend Go's standard library with commonly needed functionality.
