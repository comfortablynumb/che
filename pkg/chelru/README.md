# chelru - LRU Cache

Thread-safe, generic Least Recently Used (LRU) cache with O(1) operations.

## Features

- **Generic**: Works with any comparable key type and any value type
- **O(1) Operations**: Get and Put operations are constant time
- **Thread-Safe**: Safe for concurrent use
- **Fixed Capacity**: Automatic eviction of least recently used items
- **Zero Dependencies**: Only uses Go standard library

## Installation

```bash
go get github.com/comfortablynumb/che/pkg/chelru
```

## Quick Start

```go
package main

import (
    "fmt"
    "github.com/comfortablynumb/che/pkg/chelru"
)

func main() {
    // Create a cache with capacity of 100 items
    cache := chelru.New[string, int](100)

    // Add items
    cache.Put("user:1", 100)
    cache.Put("user:2", 200)

    // Get items
    if val, found := cache.Get("user:1"); found {
        fmt.Println("Value:", val) // Value: 100
    }

    // Check if key exists (doesn't update access time)
    if cache.Contains("user:1") {
        fmt.Println("Key exists")
    }

    // Get cache size
    fmt.Println("Size:", cache.Len()) // Size: 2
}
```

## Usage

### Creating a Cache

```go
// Create cache with capacity of 1000
cache := chelru.New[string, string](1000)

// Create cache with struct values
type User struct {
    ID   int
    Name string
}
cache := chelru.New[int, User](500)
```

### Adding and Updating Items

```go
cache := chelru.New[string, int](3)

// Add new items
cache.Put("a", 1)
cache.Put("b", 2)

// Update existing item (moves it to front as most recently used)
cache.Put("a", 10)
```

### Retrieving Items

```go
// Get returns (value, true) if found, (zero, false) if not found
if val, found := cache.Get("a"); found {
    fmt.Println("Found:", val)
} else {
    fmt.Println("Not found")
}
```

### Checking Existence

```go
// Contains checks existence without updating access time
if cache.Contains("a") {
    fmt.Println("Key exists")
}
```

### Removing Items

```go
// Remove returns true if key was present
if cache.Remove("a") {
    fmt.Println("Removed")
}
```

### Eviction Behavior

When the cache reaches capacity, adding a new item automatically evicts the least recently used item:

```go
cache := chelru.New[string, int](2)

cache.Put("a", 1)
cache.Put("b", 2)
cache.Put("c", 3) // "a" is evicted (least recently used)

_, found := cache.Get("a") // found = false
```

### Updating Access Order

Both `Get` and `Put` update an item's access time, moving it to the front:

```go
cache := chelru.New[string, int](3)

cache.Put("a", 1)
cache.Put("b", 2)
cache.Put("c", 3)

// Access "a" - it becomes most recently used
cache.Get("a")

// Now order is: a (most recent), c, b (least recent)
cache.Put("d", 4) // "b" is evicted, not "a"
```

### Getting All Keys

```go
cache := chelru.New[string, int](3)
cache.Put("a", 1)
cache.Put("b", 2)
cache.Put("c", 3)

// Returns keys in order from most to least recently used
keys := cache.Keys() // ["c", "b", "a"]
```

### Clearing the Cache

```go
cache.Clear() // Removes all items
```

### Cache Information

```go
size := cache.Len()         // Current number of items
capacity := cache.Capacity() // Maximum capacity
```

## Examples

### API Response Cache

```go
type APIResponse struct {
    Data      interface{}
    Timestamp time.Time
}

cache := chelru.New[string, APIResponse](100)

func getUser(userID string) (*User, error) {
    // Check cache first
    if resp, found := cache.Get("user:" + userID); found {
        if time.Since(resp.Timestamp) < 5*time.Minute {
            return resp.Data.(*User), nil
        }
    }

    // Cache miss or expired - fetch from API
    user, err := fetchUserFromAPI(userID)
    if err != nil {
        return nil, err
    }

    // Store in cache
    cache.Put("user:"+userID, APIResponse{
        Data:      user,
        Timestamp: time.Now(),
    })

    return user, nil
}
```

### Database Query Cache

```go
type QueryResult struct {
    Rows []map[string]interface{}
}

queryCache := chelru.New[string, QueryResult](50)

func executeQuery(sql string) (QueryResult, error) {
    // Check cache
    if result, found := queryCache.Get(sql); found {
        return result, nil
    }

    // Execute query
    result, err := db.Query(sql)
    if err != nil {
        return QueryResult{}, err
    }

    // Cache result
    queryCache.Put(sql, result)
    return result, nil
}
```

### Session Storage

```go
type Session struct {
    UserID    int
    Data      map[string]interface{}
    ExpiresAt time.Time
}

sessions := chelru.New[string, Session](1000)

func getSession(sessionID string) (*Session, bool) {
    session, found := sessions.Get(sessionID)
    if !found {
        return nil, false
    }

    // Check expiration
    if time.Now().After(session.ExpiresAt) {
        sessions.Remove(sessionID)
        return nil, false
    }

    return &session, true
}

func createSession(userID int) string {
    sessionID := generateSessionID()
    sessions.Put(sessionID, Session{
        UserID:    userID,
        Data:      make(map[string]interface{}),
        ExpiresAt: time.Now().Add(24 * time.Hour),
    })
    return sessionID
}
```

### Computed Value Cache

```go
// Cache expensive computations
cache := chelru.New[string, *big.Int](100)

func fibonacci(n int) *big.Int {
    key := fmt.Sprintf("fib:%d", n)

    if result, found := cache.Get(key); found {
        return result
    }

    // Compute
    result := computeFibonacci(n)
    cache.Put(key, result)
    return result
}
```

### Template Cache

```go
type CompiledTemplate struct {
    Template *template.Template
    ModTime  time.Time
}

templates := chelru.New[string, CompiledTemplate](50)

func getTemplate(path string) (*template.Template, error) {
    info, err := os.Stat(path)
    if err != nil {
        return nil, err
    }

    // Check cache
    if cached, found := templates.Get(path); found {
        // Return cached if file hasn't been modified
        if cached.ModTime.Equal(info.ModTime()) {
            return cached.Template, nil
        }
    }

    // Compile template
    tmpl, err := template.ParseFiles(path)
    if err != nil {
        return nil, err
    }

    templates.Put(path, CompiledTemplate{
        Template: tmpl,
        ModTime:  info.ModTime(),
    })

    return tmpl, nil
}
```

## API Reference

### Types

```go
type LRU[K comparable, V any] struct { ... }
```

### Functions

- `New[K comparable, V any](capacity int) *LRU[K, V]` - Create new LRU cache
- `Get(key K) (V, bool)` - Get value and mark as recently used
- `Put(key K, value V)` - Add or update value
- `Contains(key K) bool` - Check if key exists (no access update)
- `Remove(key K) bool` - Remove key from cache
- `Len() int` - Get current size
- `Capacity() int` - Get maximum capacity
- `Clear()` - Remove all items
- `Keys() []K` - Get all keys in LRU order

## Thread Safety

All operations are thread-safe and can be called concurrently from multiple goroutines.

```go
cache := chelru.New[string, int](100)

// Safe to use from multiple goroutines
go cache.Put("a", 1)
go cache.Get("a")
```

## Performance

- **Get**: O(1)
- **Put**: O(1)
- **Remove**: O(1)
- **Contains**: O(1)
- **Keys**: O(n) where n is the number of items

## Best Practices

1. **Choose appropriate capacity**: Set capacity based on memory constraints and expected access patterns
2. **Use Contains for existence checks**: Avoid `Get` if you don't need the value, as it updates access time
3. **Consider TTL separately**: This LRU doesn't handle expiration; implement TTL logic in your value type if needed
4. **Monitor size**: Use `Len()` to monitor cache utilization

## Related Packages

- **[cheset](../cheset)** - Set data structures (HashSet, OrderedSet)
- **[chemap](../chemap)** - Map utilities and Multimap

## License

This package is part of the Che library and shares the same license.
