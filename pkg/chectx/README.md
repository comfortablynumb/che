# chectx - Type-Safe Context Utilities

Generic, type-safe utilities for working with Go contexts. Never worry about type assertions again!

## Features

- **Type-Safe**: Uses Go generics for compile-time type safety
- **Collision-Free**: Keys are typed and isolated, preventing value collisions
- **Convenient**: Helper functions for common patterns (GetOrDefault, MustValue)
- **Zero Dependencies**: Only uses Go standard library

## Installation

```bash
go get github.com/comfortablynumb/che/pkg/chectx
```

## Quick Start

```go
package main

import (
    "context"
    "fmt"

    "github.com/comfortablynumb/che/pkg/chectx"
)

// Create typed keys
var UserIDKey = chectx.Key[int]("userID")
var UsernameKey = chectx.Key[string]("username")

func main() {
    ctx := context.Background()

    // Store values with type safety
    ctx = chectx.WithValue(ctx, UserIDKey, 123)
    ctx = chectx.WithValue(ctx, UsernameKey, "john_doe")

    // Retrieve values with type safety (no type assertions needed!)
    userID, ok := chectx.Value(ctx, UserIDKey)
    if ok {
        fmt.Printf("User ID: %d\n", userID)
    }

    username := chectx.GetOrDefault(ctx, UsernameKey, "anonymous")
    fmt.Printf("Username: %s\n", username)
}
```

## Usage

### Creating Keys

Create typed context keys to ensure type safety:

```go
// Create a key for string values
var RequestIDKey = chectx.Key[string]("request-id")

// Create a key for custom types
type User struct {
    ID   int
    Name string
}
var UserKey = chectx.Key[*User]("user")
```

### Storing Values

Use `WithValue` to store typed values in context:

```go
ctx := context.Background()

// Store a string
ctx = chectx.WithValue(ctx, RequestIDKey, "abc-123")

// Store a custom type
user := &User{ID: 1, Name: "John"}
ctx = chectx.WithValue(ctx, UserKey, user)
```

### Retrieving Values

Several methods for retrieving values:

#### Value() - Safe Retrieval

Returns the value and a boolean indicating if it was found:

```go
requestID, ok := chectx.Value(ctx, RequestIDKey)
if ok {
    fmt.Println("Request ID:", requestID)
} else {
    fmt.Println("Request ID not found")
}
```

#### MustValue() - Panic on Missing

Panics if the value is not found (useful when value is required):

```go
// Will panic if not found
requestID := chectx.MustValue(ctx, RequestIDKey)
fmt.Println("Request ID:", requestID)
```

#### GetOrDefault() - Default Value

Returns a default value if not found:

```go
// Returns "unknown" if not found
requestID := chectx.GetOrDefault(ctx, RequestIDKey, "unknown")
fmt.Println("Request ID:", requestID)
```

## Examples

### HTTP Middleware Example

```go
package main

import (
    "context"
    "net/http"

    "github.com/comfortablynumb/che/pkg/chectx"
)

var UserIDKey = chectx.Key[int]("userID")

func AuthMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Extract user ID from token (simplified)
        userID := extractUserIDFromToken(r)

        // Store in context with type safety
        ctx := chectx.WithValue(r.Context(), UserIDKey, userID)

        // Call next handler with updated context
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}

func Handler(w http.ResponseWriter, r *http.Request) {
    // Retrieve user ID with type safety (no type assertion!)
    userID, ok := chectx.Value(r.Context(), UserIDKey)
    if !ok {
        http.Error(w, "Unauthorized", http.StatusUnauthorized)
        return
    }

    // userID is already an int!
    fmt.Fprintf(w, "Hello, user %d", userID)
}
```

### Request Scoped Data

```go
package main

import (
    "context"

    "github.com/comfortablynumb/che/pkg/chectx"
)

// Define keys for request-scoped data
var (
    RequestIDKey  = chectx.Key[string]("request-id")
    TraceIDKey    = chectx.Key[string]("trace-id")
    UserAgentKey  = chectx.Key[string]("user-agent")
)

func processRequest(ctx context.Context) {
    // All retrievals are type-safe!
    requestID := chectx.GetOrDefault(ctx, RequestIDKey, "no-id")
    traceID := chectx.GetOrDefault(ctx, TraceIDKey, "no-trace")
    userAgent := chectx.GetOrDefault(ctx, UserAgentKey, "unknown")

    log.Printf("Processing request %s (trace: %s) from %s",
        requestID, traceID, userAgent)
}
```

### Complex Types

```go
package main

import (
    "context"

    "github.com/comfortablynumb/che/pkg/chectx"
)

type Config struct {
    Database string
    Cache    string
    Timeout  int
}

var ConfigKey = chectx.Key[*Config]("config")

func main() {
    ctx := context.Background()

    config := &Config{
        Database: "postgres://localhost",
        Cache:    "redis://localhost",
        Timeout:  30,
    }

    // Store complex type
    ctx = chectx.WithValue(ctx, ConfigKey, config)

    // Retrieve with full type safety
    cfg, ok := chectx.Value(ctx, ConfigKey)
    if ok {
        // cfg is *Config, no type assertion needed!
        println("Database:", cfg.Database)
        println("Timeout:", cfg.Timeout)
    }
}
```

## Key Isolation

Keys with the same name are isolated from each other:

```go
key1 := chectx.Key[string]("value")
key2 := chectx.Key[string]("value")

ctx := context.Background()
ctx = chectx.WithValue(ctx, key1, "value1")

// key2 won't retrieve "value1" even though they have the same name
val, ok := chectx.Value(ctx, key2)
// ok will be false

// key1 still works
val, ok = chectx.Value(ctx, key1)
// ok will be true, val will be "value1"
```

This prevents accidental collisions between different parts of your application.

## API Reference

### Functions

- `Key[T](name string) *contextKey[T]` - Creates a new typed context key
- `WithValue[T](ctx, key, value) context.Context` - Stores a typed value in context
- `Value[T](ctx, key) (T, bool)` - Retrieves a typed value from context
- `MustValue[T](ctx, key) T` - Retrieves a value or panics if not found
- `GetOrDefault[T](ctx, key, defaultValue) T` - Retrieves a value or returns default

## Best Practices

1. **Define Keys as Package Variables**: Define keys at package level for reusability
2. **Use Descriptive Names**: Give keys descriptive names for debugging
3. **Pointer Types for Structs**: Use pointer types for struct values to avoid copies
4. **Use GetOrDefault**: Prefer `GetOrDefault` over `Value` when you have a sensible default
5. **Use MustValue Sparingly**: Only use `MustValue` when the value is absolutely required

## Related Packages

- **[chesignal](../chesignal)** - Graceful shutdown utilities
- **[chehttp](../chehttp)** - Ergonomic HTTP client and server utilities

## License

This package is part of the Che library and shares the same license.
