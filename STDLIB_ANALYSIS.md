# Go Standard Library Gap Analysis

**Date**: 2025-11-08
**Current Go Version**: 1.24.7
**Target Go Version in go.mod**: 1.20

## Executive Summary

Since Go 1.21, the Go standard library has added **`slices`** and **`maps`** packages that provide many of the functions this repository implements. This analysis identifies:
1. What's now redundant with the standard library
2. What still adds value
3. What new features should be added

---

## Current Repository Features

### `cheslice` Package
- `Union()` - Combines multiple slices
- `ForEach()` - Iterates over elements with callback
- `Map()` - Transforms slice elements
- `Filter()` - Filters slice elements
- `Fill()` - Creates slice filled with value
- `Diff()` - Elements in first slice not in others
- `Chunk()` - Splits slice into chunks
- `Unique()` - Returns distinct elements
- `Intersect()` - Elements common to all slices
- `Exists()` - Checks if element exists
- `Len()` - Sum of lengths

### `chemap` Package
- `Keys()` - Returns map keys as slice

### `chetest` Package
- `RequireEqual()` - Testing assertion helper

---

## Analysis: Redundant vs. Standard Library

### ✅ Now in Standard Library (Go 1.21+)

#### Overlapping with `slices` package:
- **`Chunk()`** - `slices.Chunk()` added in Go 1.23 (returns iterator)
- **`Exists()`** - Similar to `slices.Contains()` (Go 1.21+)
- **`Union()`** - Similar to `slices.Concat()` (Go 1.22+)

#### Overlapping with `maps` package:
- **`chemap.Keys()`** - `maps.Keys()` available (Go 1.21+, returns iterator in 1.23+)

### ⚠️ Partially Redundant (Different Behavior)

#### `cheslice` functions with partial overlap:
- **`Map()`** - No direct stdlib equivalent, but can use `for range` loops or iterators (Go 1.23+)
- **`Filter()`** - Can use `slices.DeleteFunc()` or iterators, but not identical API
- **`ForEach()`** - Standard `for range` loop is idiomatic Go

### ✨ Still Valuable (No Direct Standard Library Equivalent)

These functions provide unique value:

1. **`Unique()`** - While `slices.Compact()` removes consecutive duplicates, it requires sorting first. `Unique()` maintains order and removes all duplicates.

2. **`Diff()`** - Set difference operation not in stdlib

3. **`Intersect()`** - Set intersection operation not in stdlib

4. **`Fill()`** - Creating slices with repeated values (though `slices.Repeat()` exists in Go 1.23+)

5. **`Len()`** - Sum of multiple slice lengths (niche use case)

---

## Recommendations

### 🗑️ Consider Deprecating (Now in Stdlib)

If targeting Go 1.21+, these could be deprecated in favor of stdlib:

1. **`cheslice.Exists()`** → Use `slices.Contains()`
2. **`cheslice.Union()`** → Use `slices.Concat()`
3. **`cheslice.Chunk()`** → Use `slices.Chunk()` (Go 1.23+)
4. **`chemap.Keys()`** → Use `maps.Keys()` + `slices.Collect()` for slice

### 🔄 Consider Refactoring

1. **`Map()` and `Filter()`** - Consider deprecating in favor of iterator-based approaches (Go 1.23+)
2. **`ForEach()`** - Standard `for range` is more idiomatic; questionable value

### ✅ Keep and Enhance

These provide unique functionality worth keeping:

1. **`Unique()`** - Order-preserving deduplication
2. **`Diff()`** - Set difference
3. **`Intersect()`** - Set intersection

---

## Suggested Additions

### High Priority: Data Structures

#### 1. **Set** (Generic Set Implementation)
```go
type Set[T comparable] struct {
    // HashSet implementation
}
```
Methods:
- `Add()`, `Remove()`, `Contains()`
- `Union()`, `Intersect()`, `Diff()`, `SymmetricDiff()`
- `IsSubset()`, `IsSuperset()`
- `Equal()`, `Clone()`
- `ToSlice()`, `FromSlice()`

#### 2. **OrderedSet**
Like Set but maintains insertion order

#### 3. **Queue** (Generic Queue)
```go
type Queue[T any] struct {
    // Ring buffer or linked list implementation
}
```
Methods: `Enqueue()`, `Dequeue()`, `Peek()`, `IsEmpty()`, `Len()`

#### 4. **Stack** (Generic Stack)
```go
type Stack[T any] struct {}
```
Methods: `Push()`, `Pop()`, `Peek()`, `IsEmpty()`, `Len()`

#### 5. **PriorityQueue** (Generic Heap-based)
Already exists in `container/heap` but generic wrapper would improve ergonomics

#### 6. **LRU Cache**
```go
type LRUCache[K comparable, V any] struct {}
```
Methods: `Get()`, `Put()`, `Remove()`, `Clear()`, `Len()`, `Cap()`

#### 7. **Trie** (Prefix Tree)
For string operations, autocomplete, etc.

#### 8. **LinkedList** (Doubly Linked)
`container/list` exists but is not generic

#### 9. **CircularBuffer** / **RingBuffer**
```go
type RingBuffer[T any] struct {}
```

#### 10. **Multimap**
```go
type Multimap[K comparable, V any] struct {
    // map[K][]V wrapper
}
```

### Medium Priority: Algorithm Functions

#### Slice Operations
1. **`Partition()`** - Split slice by predicate into two slices
2. **`GroupBy()`** - Group elements by key function
3. **`Reduce()` / `Fold()`** - Reduce slice to single value
4. **`Flatten()`** - Flatten nested slices
5. **`Zip()` / `Unzip()`** - Combine/split multiple slices
6. **`Window()`** - Sliding window iterator
7. **`Take()` / `Drop()` / `TakeWhile()` / `DropWhile()`**
8. **`Sample()` / `Shuffle()`** - Random sampling and shuffling
9. **`Permutations()` / `Combinations()`**
10. **`Any()` / `All()` / `None()`** - Predicate checks

#### Map Operations
1. **`Values()`** - Extract values as slice
2. **`Invert()`** - Swap keys and values
3. **`Merge()`** - Merge multiple maps with conflict resolution
4. **`Filter()`** - Filter map by predicate
5. **`Map()`** - Transform map values
6. **`GroupBy()`** - Create map from slice using key function

#### Numeric Operations
1. **`Sum()` / `Product()`** - For numeric slices
2. **`Mean()` / `Median()` / `Mode()`** - Statistics
3. **`Min()` / `Max()`** - Already in stdlib (slices.Min/Max) but could add for multiple values
4. **`Range()`** - Generate numeric ranges
5. **`Clamp()`** - Constrain value to range

#### String Operations (High-level helpers)
1. **`StringSet`** - Common set operations on strings
2. **`Levenshtein()`** - Edit distance
3. **`CommonPrefix()` / `CommonSuffix()`**
4. **`CamelCase()` / `SnakeCase()` / `KebabCase()`** - Case conversions

### Lower Priority: Functional Utilities

1. **`Curry()` / `Partial()`** - Function currying
2. **`Compose()` / `Pipe()`** - Function composition
3. **`Memoize()`** - Function result caching
4. **`Debounce()` / `Throttle()`** - Rate limiting

### File & I/O Utilities

1. **`ReadLines()`** / **`WriteLines()`** - Simple line-based file I/O
2. **`ReadJSON()` / **`WriteJSON()`** - JSON file helpers
3. **`CopyFile()` / **`MoveFile()`** - File operations
4. **`TempDir()` / **`TempFile()`** with auto-cleanup - Resource management
5. **`WalkDir()`** wrapper - More ergonomic directory traversal
6. **`PathExists()` / `IsDir()` / `IsFile()`** - Path checking helpers

### Concurrency Utilities

1. **`Parallel()`** - Parallel execution of functions
2. **`WorkerPool`** - Generic worker pool
3. **`RateLimiter`** - Rate limiting
4. **`Semaphore`** - Counting semaphore (generic wrapper)
5. **`Once`** - Generic sync.Once wrapper
6. **`ErrGroup`** wrapper - Improve ergonomics of golang.org/x/sync/errgroup

### Network Utilities

1. **`RetryHTTP()`** - HTTP client with retry logic
2. **`DownloadFile()`** - Simple file download
3. **`IsPortOpen()` / `GetFreePort()`** - Port utilities
4. **`ParseURL()` with validation** - Enhanced URL parsing

### Error Handling

1. **`Must()`** - Panic on error (common pattern)
2. **`Try()` / `Catch()`** - Try/catch pattern
3. **`MultiError`** - Collect multiple errors
4. **`Unwrap()` helpers** - Better error unwrapping

### Validation

1. **`Validate()`** - Struct validation framework
2. **`IsEmail()` / `IsURL()` / `IsIP()`** - Common validators
3. **`InRange()` / `OneOf()`** - Value validators

---

## Migration Strategy

### Phase 1: Update Compatibility
1. Update `go.mod` to Go 1.21+ (currently 1.20)
2. Add deprecation warnings to functions now in stdlib
3. Add migration guide in documentation

### Phase 2: Core Data Structures
1. Implement `Set` (highest value add)
2. Implement `Queue`, `Stack`
3. Implement `LRU Cache`

### Phase 3: Enhanced Algorithms
1. Add missing slice/map operations
2. Add numeric/statistical functions
3. Add functional utilities

### Phase 4: Utilities
1. File/IO helpers
2. Concurrency utilities
3. Network helpers

---

## Architecture Recommendations

### 1. Package Organization
```
pkg/
├── cheslice/      # Slice utilities (consider renaming)
├── chemap/        # Map utilities
├── cheset/        # NEW: Set implementation
├── checollections/ # NEW: Other data structures
├── chemath/       # NEW: Math/statistics
├── chefile/       # NEW: File utilities
├── chehttp/       # NEW: HTTP utilities
├── checoncurrent/ # NEW: Concurrency utilities
└── chetest/       # Testing helpers
```

### 2. Embrace Iterators (Go 1.23+)
Consider iterator-based APIs for:
- `Set.All()` → `iter.Seq[T]`
- `Map()`, `Filter()` → Chain with iterators
- `Window()`, `Chunk()` → Return iterators

### 3. Options Pattern
Use functional options for complex functions:
```go
func NewSet[T comparable](opts ...SetOption[T]) *Set[T]
```

### 4. Thread Safety Options
Offer both concurrent-safe and unsafe versions:
```go
type Set[T comparable] struct { ... }        // Not thread-safe
type ConcurrentSet[T comparable] struct { ... } // Thread-safe
```

---

## Competitive Analysis

Similar libraries in the ecosystem:
- **`golang.org/x/exp/slices`** (precursor to stdlib slices)
- **`samber/lo`** - Lodash-style Go library (very popular, 18k+ stars)
- **`thoas/go-funk`** - Functional utilities
- **`emirpasic/gods`** - Data structures (7k+ stars)

**Differentiation Strategy**:
- Focus on **zero dependencies** (unlike many alternatives)
- **100% test coverage** guarantee
- **Stdlib-compatible** APIs where overlap exists
- **Modern Go idioms** (generics, iterators)
- **Production-grade** performance and safety

---

## Testing Recommendations

1. **Maintain 100% coverage** (current goal)
2. **Add benchmarks** for all operations
3. **Fuzz testing** for complex algorithms
4. **Thread safety tests** for concurrent structures
5. **Comparison tests** vs stdlib where applicable

---

## Documentation Improvements

1. **Add godoc examples** for all public functions
2. **Performance characteristics** (Big-O notation)
3. **Thread safety** guarantees
4. **Migration guide** from stdlib (when to use che vs stdlib)
5. **Comparison table** with popular alternatives

---

## Summary

**Immediate Actions**:
1. Update to Go 1.21+ in go.mod
2. Add deprecation notices for functions now in stdlib
3. Implement `Set` data structure (highest ROI)
4. Add comprehensive slice/map utilities not in stdlib

**Long-term Vision**:
Position as a **production-ready, zero-dependency, comprehensive standard library extension** focusing on data structures and utilities that Go's stdlib intentionally omits.

**Core Value Props**:
- ✅ Zero dependencies
- ✅ 100% test coverage
- ✅ Generic-first design
- ✅ Stdlib-compatible APIs
- ✅ Production-grade quality
