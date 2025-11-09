# cheratelimit

Token bucket rate limiting for Go applications.

## Features

- Token bucket algorithm
- Thread-safe
- Per-key rate limiting with automatic cleanup
- Blocking wait with context support
- Dynamic rate and burst adjustment
- Zero external dependencies

## Installation

```bash
go get github.com/comfortablynumb/che/pkg/cheratelimit
```

## Usage

### Basic Rate Limiting

```go
package main

import (
    "fmt"
    "github.com/comfortablynumb/che/pkg/cheratelimit"
)

func main() {
    // 10 requests per second, burst of 5
    limiter := cheratelimit.New(10, 5)

    // Try to acquire a token
    if limiter.Allow() {
        fmt.Println("Request allowed")
        // Process request
    } else {
        fmt.Println("Rate limit exceeded")
    }

    // Check available tokens
    tokens := limiter.Tokens()
    fmt.Printf("Available tokens: %.2f\n", tokens)
}
```

### Blocking Wait

```go
package main

import (
    "context"
    "fmt"
    "time"
    "github.com/comfortablynumb/che/pkg/cheratelimit"
)

func main() {
    limiter := cheratelimit.New(10, 1)

    // Wait blocks until a token is available
    ctx := context.Background()
    err := limiter.Wait(ctx)
    if err != nil {
        fmt.Println("Wait cancelled:", err)
        return
    }

    fmt.Println("Token acquired, processing request")

    // Wait with timeout
    ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
    defer cancel()

    err = limiter.Wait(ctx)
    if err == context.DeadlineExceeded {
        fmt.Println("Timeout waiting for token")
    }
}
```

### Reserve a Token

```go
limiter := cheratelimit.New(10, 1)

// Reserve returns the duration to wait
waitTime := limiter.Reserve()

if waitTime > 0 {
    fmt.Printf("Need to wait %v before processing\n", waitTime)
    time.Sleep(waitTime)
}

fmt.Println("Processing request")
```

### Per-Key Rate Limiting

```go
package main

import (
    "fmt"
    "github.com/comfortablynumb/che/pkg/cheratelimit"
)

func main() {
    // Different rate limit for each API key
    limiter := cheratelimit.NewPerKey(100, 10)

    // Each key has independent limits
    if limiter.Allow("user-123") {
        fmt.Println("Request from user-123 allowed")
    }

    if limiter.Allow("user-456") {
        fmt.Println("Request from user-456 allowed")
    }

    // Wait for a specific key
    ctx := context.Background()
    err := limiter.Wait(ctx, "user-123")
    if err == nil {
        fmt.Println("Token acquired for user-123")
    }
}
```

### Dynamic Rate Adjustment

```go
limiter := cheratelimit.New(10, 5)

// Increase rate to 100 req/s
limiter.SetRate(100)

// Decrease burst to 2
limiter.SetBurst(2)
```

### HTTP Middleware Example

```go
package main

import (
    "net/http"
    "github.com/comfortablynumb/che/pkg/cheratelimit"
)

func RateLimitMiddleware(limiter *cheratelimit.Limiter) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            if !limiter.Allow() {
                http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
                return
            }
            next.ServeHTTP(w, r)
        })
    }
}

func main() {
    limiter := cheratelimit.New(100, 20) // 100 req/s, burst 20

    mux := http.NewServeMux()
    mux.HandleFunc("/api/data", func(w http.ResponseWriter, r *http.Request) {
        w.Write([]byte("OK"))
    })

    handler := RateLimitMiddleware(limiter)(mux)
    http.ListenAndServe(":8080", handler)
}
```

### Per-User HTTP Middleware

```go
func PerUserRateLimit(limiter *cheratelimit.PerKeyLimiter) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            // Extract user ID from request (header, token, etc.)
            userID := r.Header.Get("X-User-ID")
            if userID == "" {
                userID = r.RemoteAddr // Fallback to IP
            }

            if !limiter.Allow(userID) {
                http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
                return
            }

            next.ServeHTTP(w, r)
        })
    }
}
```

## API

### Limiter

- `New(rate float64, burst int) *Limiter` - Create a new limiter
  - `rate`: Requests per second
  - `burst`: Maximum requests that can be made at once

- `Allow() bool` - Try to acquire a token (non-blocking)
- `Wait(ctx context.Context) error` - Wait for a token (blocking)
- `Reserve() time.Duration` - Reserve a token, returns wait duration
- `Tokens() float64` - Get current available tokens
- `SetRate(rate float64)` - Update the rate limit
- `SetBurst(burst int)` - Update the burst size

### PerKeyLimiter

- `NewPerKey(rate float64, burst int) *PerKeyLimiter` - Create per-key limiter
- `Allow(key string) bool` - Try to acquire token for key
- `Wait(ctx context.Context, key string) error` - Wait for token for key

## How It Works

The token bucket algorithm:

1. Tokens are added to a bucket at a constant rate
2. Each request consumes one token
3. Bucket has a maximum capacity (burst)
4. If bucket is empty, requests are denied
5. Tokens accumulate when idle (up to burst size)

This allows for:
- Steady-state rate limiting (rate parameter)
- Handling bursts of traffic (burst parameter)
- Fair distribution over time

## Performance

The limiter is highly efficient:
- O(1) Allow() operation
- Lock-free token refill calculation
- Minimal memory overhead
- Per-key limiter automatically cleans up old entries

## Thread Safety

All operations are thread-safe and can be called concurrently from multiple goroutines.

## License

MIT
