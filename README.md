# Che - An Extension for the Go Standard Library

![Go Report Card](https://goreportcard.com/badge/github.com/comfortablynumb/che)
![Build Status](https://github.com/comfortablynumb/che/actions/workflows/build.yml/badge.svg)
![License](https://img.shields.io/github/license/comfortablynumb/che)

<p align="center">
  <img alt="Che!" width="450" height="482" src="https://github.com/comfortablynumb/che/raw/main/docs/images/gopher.png" />
</p>

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

#### `chehttp` - HTTP Client
- Ergonomic HTTP client with builder pattern
- Automatic JSON marshalling/unmarshalling for requests and responses
- Request options for headers, timeouts, body
- Connection timeout vs request timeout distinction
- Request lifecycle hooks (pre-request, post-request, on-success, on-error, on-complete)
- Retry configuration with exponential, linear, and fixed backoff strategies
- Context-aware methods for cancellation support
- Response body streaming
- Interface-based design for easy mocking

#### `chestring` - String Utilities
- Case conversions: ToCamelCase, ToPascalCase, ToSnakeCase, ToKebabCase, ToScreamingSnakeCase
- Transformations: Capitalize, Uncapitalize, Reverse
- Validation: IsEmpty, IsBlank, IsNotEmpty, IsNotBlank
- Truncation: Truncate by length or words
- Search: ContainsAny, ContainsAll
- Other utilities: Repeat, RemoveWhitespace, DefaultIfEmpty, DefaultIfBlank, SplitAndTrim

#### `cheenv` - Environment Variables
- Type-safe environment variable access (string, int, int64, float64, bool, duration)
- Default values and Must* variants for required config
- List support with custom separators (GetStringList, GetIntList)
- Flexible boolean parsing (true/false, yes/no, on/off, 1/0, y/n, t/f)
- Batch operations: GetAll, GetWithPrefix
- Variable management: Set, Unset, Has

#### `chectx` - Context Utilities
- Type-safe context key/value pairs using generics
- Eliminates type assertions and key collisions
- WithValue, Value, MustValue, GetOrDefault functions

#### `chesignal` - Graceful Shutdown
- Signal handling utilities for graceful application shutdown
- Configurable signals and timeout
- Ordered shutdown function execution
- Lifecycle callbacks (OnShutdownStart, OnShutdownComplete, OnShutdownTimeout)
- Context-aware shutdown

#### `chetest` - Testing Helpers
- RequireEqual with custom messages
- Assertion utilities for tests

## Features Status

- [x] Slice functions: Unique, Diff, Intersect, Map, Filter, Reduce, GroupBy, Partition, etc.
- [x] Map functions: Keys, Values, Invert, Filter, MapValues, Merge, Pick, Omit
- [x] Data structures: HashSet, OrderedSet, Queue, Stack, Multimap
- [x] Linked data structures: LinkedList, DoublyLinkedList
- [x] Tree structures: Binary Search Tree
- [x] HTTP client: Ergonomic client with hooks, retries, and context support
- [x] String utilities: Case conversions, validation, truncation, search
- [x] Environment utilities: Type-safe env var access with defaults
- [x] Context utilities: Type-safe context values with generics
- [x] Signal utilities: Graceful shutdown handling
- [ ] More data structures: LRU Cache, Trie, AVL Tree, etc.
- [ ] File handling functions

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

### HTTP Client
```go
// Create client with builder
client := chehttp.NewBuilder().
    WithBaseURL("https://api.example.com").
    WithDefaultHeader("Authorization", "Bearer token").
    WithRequestTimeout(30 * time.Second).
    WithConnectionTimeout(10 * time.Second).
    WithRetry(&chehttp.RetryConfig{
        MaxRetries: 3,
        BackoffStrategy: &chehttp.ExponentialBackoff{
            BaseDelay: 1 * time.Second,
            Multiplier: 2.0,
        },
    }).
    WithPreRequestHook(func(ctx *chehttp.HookContext) error {
        log.Printf("Request: %s %s", ctx.Method, ctx.URL)
        return nil
    }).
    Build()

// Make requests with automatic JSON handling
type User struct {
    ID   int    `json:"id"`
    Name string `json:"name"`
}

var user User
resp, err := client.Get("/users/1", chehttp.WithSuccess(&user))
if resp.IsSuccess() {
    fmt.Println("User:", user.Name)
}

// POST with JSON body and context
ctx := context.WithTimeout(context.Background(), 5*time.Second)
newUser := User{Name: "John"}
resp, err = client.PostWithCtx(ctx, "/users", chehttp.WithJSONBody(newUser))
```

### String Utilities
```go
// Case conversions
chestring.ToCamelCase("hello_world")      // "helloWorld"
chestring.ToSnakeCase("HelloWorld")       // "hello_world"
chestring.ToKebabCase("helloWorld")       // "hello-world"

// Validation
chestring.IsBlank("  ")                   // true
chestring.IsNotEmpty("hello")             // true

// Truncation
chestring.Truncate("Hello World", 5)      // "Hello..."
chestring.TruncateWords("one two three", 2) // "one two..."

// Search
chestring.ContainsAny("hello", "x", "e")  // true
```

### Environment Variables
```go
// String values with defaults
dbHost := cheenv.Get("DB_HOST", "localhost")
dbName := cheenv.MustGet("DB_NAME") // panics if not set

// Typed values
port := cheenv.GetInt("PORT", 8080)
debug := cheenv.GetBool("DEBUG", false) // accepts: true, yes, on, 1, etc.
timeout := cheenv.GetDuration("TIMEOUT", 30*time.Second)

// Lists
hosts := cheenv.GetStringList("ALLOWED_HOSTS", ",", []string{"localhost"})
ports := cheenv.GetIntList("PORTS", ",", []int{8080})

// Batch operations
appConfig := cheenv.GetWithPrefix("APP_") // Get all APP_* variables
```

### Context Utilities
```go
// Type-safe context values
userIDKey := chectx.Key[int]("userID")
requestIDKey := chectx.Key[string]("requestID")

// Set values
ctx := chectx.WithValue(context.Background(), userIDKey, 42)
ctx = chectx.WithValue(ctx, requestIDKey, "abc123")

// Get values (type-safe, no assertions needed)
userID, ok := chectx.Value(ctx, userIDKey)       // userID is int
requestID := chectx.MustValue(ctx, requestIDKey) // panics if not found
```

### Graceful Shutdown
```go
// Define shutdown functions
shutdownFuncs := []chesignal.ShutdownFunc{
    func(ctx context.Context) error {
        log.Println("Closing database connection...")
        return db.Close()
    },
    func(ctx context.Context) error {
        log.Println("Shutting down HTTP server...")
        return server.Shutdown(ctx)
    },
}

// Wait for shutdown signal (SIGINT, SIGTERM)
config := &chesignal.Config{
    Signals: []os.Signal{os.Interrupt, syscall.SIGTERM},
    Timeout: 30 * time.Second,
    OnShutdownStart: func() {
        log.Println("Shutdown initiated...")
    },
}

err := chesignal.WaitForShutdown(config, shutdownFuncs...)
```

## Credits

* [gopherize.me](https://gopherize.me/): For the excellent gopher image!