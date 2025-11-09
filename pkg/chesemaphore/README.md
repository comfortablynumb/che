# chesemaphore

Weighted semaphore for concurrency control in Go.

## Features

- Weighted semaphore (acquire different amounts)
- Blocking and non-blocking acquire operations
- Context support for cancellation and timeouts
- Thread-safe operations
- Simple API
- Zero external dependencies

## Installation

```bash
go get github.com/comfortablynumb/che/pkg/chesemaphore
```

## Usage

### Basic Semaphore

```go
package main

import (
    "context"
    "fmt"
    "github.com/comfortablynumb/che/pkg/chesemaphore"
)

func main() {
    // Create a semaphore with size 10
    sem := chesemaphore.New(10)

    ctx := context.Background()

    // Acquire 5 resources
    err := sem.Acquire(ctx, 5)
    if err != nil {
        fmt.Println("Error:", err)
        return
    }

    fmt.Printf("Current: %d, Available: %d\n", sem.Current(), sem.Available())

    // Release 5 resources
    sem.Release(5)

    fmt.Printf("Current: %d, Available: %d\n", sem.Current(), sem.Available())
}
```

### Limiting Concurrency

```go
import (
    "context"
    "sync"
)

func main() {
    // Limit to 5 concurrent operations
    sem := chesemaphore.New(5)
    var wg sync.WaitGroup

    for i := 0; i < 100; i++ {
        wg.Add(1)
        go func(id int) {
            defer wg.Done()

            // Acquire semaphore
            ctx := context.Background()
            if err := sem.Acquire(ctx, 1); err != nil {
                return
            }
            defer sem.Release(1)

            // Do work (max 5 concurrent)
            processJob(id)
        }(i)
    }

    wg.Wait()
}
```

### Non-blocking Acquire

```go
sem := chesemaphore.New(10)

// Try to acquire without blocking
if sem.TryAcquire(5) {
    defer sem.Release(5)

    // Resources acquired, do work
    doWork()
} else {
    // Resources not available
    fmt.Println("Resources not available")
}
```

### With Timeout

```go
import (
    "context"
    "time"
)

sem := chesemaphore.New(5)

// Try to acquire with timeout
ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
defer cancel()

err := sem.Acquire(ctx, 3)
if err == context.DeadlineExceeded {
    fmt.Println("Timeout waiting for resources")
} else if err != nil {
    fmt.Println("Error:", err)
} else {
    defer sem.Release(3)
    doWork()
}
```

### Weighted Acquire

```go
sem := chesemaphore.New(100)

// Small job requires 10 resources
if err := sem.Acquire(ctx, 10); err == nil {
    defer sem.Release(10)
    processSmallJob()
}

// Large job requires 50 resources
if err := sem.Acquire(ctx, 50); err == nil {
    defer sem.Release(50)
    processLargeJob()
}
```

## Examples

### Database Connection Pool

```go
type DBPool struct {
    sem *chesemaphore.Semaphore
}

func NewDBPool(maxConnections int) *DBPool {
    return &DBPool{
        sem: chesemaphore.New(int64(maxConnections)),
    }
}

func (p *DBPool) Query(ctx context.Context, query string) error {
    // Acquire connection
    if err := p.sem.Acquire(ctx, 1); err != nil {
        return err
    }
    defer p.sem.Release(1)

    // Execute query
    return executeQuery(query)
}

func (p *DBPool) Available() int64 {
    return p.sem.Available()
}

func main() {
    pool := NewDBPool(10)

    ctx := context.Background()
    err := pool.Query(ctx, "SELECT * FROM users")
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("Available connections: %d\n", pool.Available())
}
```

### Rate Limiting with Resources

```go
type ResourceLimiter struct {
    sem *chesemaphore.Semaphore
}

func NewResourceLimiter(capacity int64) *ResourceLimiter {
    return &ResourceLimiter{
        sem: chesemaphore.New(capacity),
    }
}

func (rl *ResourceLimiter) Do(ctx context.Context, cost int64, fn func() error) error {
    if err := rl.sem.Acquire(ctx, cost); err != nil {
        return err
    }
    defer rl.sem.Release(cost)

    return fn()
}

func main() {
    // 100 units of capacity
    limiter := NewResourceLimiter(100)

    // Small operation costs 10 units
    ctx := context.Background()
    err := limiter.Do(ctx, 10, func() error {
        return performSmallOperation()
    })

    // Large operation costs 50 units
    err = limiter.Do(ctx, 50, func() error {
        return performLargeOperation()
    })

    if err != nil {
        log.Fatal(err)
    }
}
```

### API Request Throttling

```go
type APIThrottler struct {
    sem *chesemaphore.Semaphore
}

func NewAPIThrottler(maxConcurrent int) *APIThrottler {
    return &APIThrottler{
        sem: chesemaphore.New(int64(maxConcurrent)),
    }
}

func (t *APIThrottler) Request(ctx context.Context, url string) (*Response, error) {
    // Acquire slot for request
    if err := t.sem.Acquire(ctx, 1); err != nil {
        return nil, err
    }
    defer t.sem.Release(1)

    // Make API request
    return http.Get(url)
}

func main() {
    // Limit to 5 concurrent API requests
    throttler := NewAPIThrottler(5)

    urls := []string{
        "https://api.example.com/users",
        "https://api.example.com/posts",
        // ... more URLs
    }

    var wg sync.WaitGroup
    for _, url := range urls {
        wg.Add(1)
        go func(u string) {
            defer wg.Done()

            ctx := context.Background()
            resp, err := throttler.Request(ctx, u)
            if err != nil {
                log.Println("Error:", err)
                return
            }

            processResponse(resp)
        }(url)
    }

    wg.Wait()
}
```

### Download Manager

```go
type DownloadManager struct {
    bandwidth *chesemaphore.Semaphore // KB/s
}

func NewDownloadManager(maxBandwidth int64) *DownloadManager {
    return &DownloadManager{
        bandwidth: chesemaphore.New(maxBandwidth),
    }
}

func (dm *DownloadManager) Download(ctx context.Context, url string, size int64) error {
    // Calculate required bandwidth (in KB)
    requiredBandwidth := size / 1024

    // Acquire bandwidth
    if err := dm.bandwidth.Acquire(ctx, requiredBandwidth); err != nil {
        return err
    }
    defer dm.bandwidth.Release(requiredBandwidth)

    // Perform download
    return performDownload(url)
}

func main() {
    // 1000 KB/s max bandwidth
    dm := NewDownloadManager(1000)

    files := []struct {
        url  string
        size int64
    }{
        {"https://example.com/small.zip", 100 * 1024},  // 100 KB
        {"https://example.com/large.zip", 500 * 1024},  // 500 KB
    }

    var wg sync.WaitGroup
    for _, file := range files {
        wg.Add(1)
        go func(url string, size int64) {
            defer wg.Done()

            ctx := context.Background()
            err := dm.Download(ctx, url, size)
            if err != nil {
                log.Printf("Download failed: %v\n", err)
            }
        }(file.url, file.size)
    }

    wg.Wait()
}
```

### Worker Pool with Priorities

```go
type PriorityPool struct {
    sem *chesemaphore.Semaphore
}

func NewPriorityPool(size int) *PriorityPool {
    return &PriorityPool{
        sem: chesemaphore.New(int64(size)),
    }
}

const (
    LowPriority  = 1
    MedPriority  = 3
    HighPriority = 5
)

func (pp *PriorityPool) Execute(ctx context.Context, priority int64, fn func()) error {
    if err := pp.sem.Acquire(ctx, priority); err != nil {
        return err
    }
    defer pp.sem.Release(priority)

    fn()
    return nil
}

func main() {
    pool := NewPriorityPool(10)

    ctx := context.Background()

    // High priority task
    pool.Execute(ctx, HighPriority, func() {
        processHighPriorityTask()
    })

    // Low priority task
    pool.Execute(ctx, LowPriority, func() {
        processLowPriorityTask()
    })
}
```

## API Reference

### Creating

- `New(size int64) *Semaphore` - Create a new semaphore with the given size

### Acquiring

- `Acquire(ctx context.Context, weight int64) error` - Acquire with blocking
  - Blocks until resources are available
  - Returns error if context is cancelled or weight exceeds limit

- `TryAcquire(weight int64) bool` - Try to acquire without blocking
  - Returns true if successful, false otherwise
  - Never blocks

### Releasing

- `Release(weight int64)` - Release the specified weight
  - If releasing more than acquired, current resets to 0
  - Broadcasts to waiting goroutines

### Inspection

- `Available() int64` - Get available resources
- `Size() int64` - Get total semaphore size
- `Current() int64` - Get current usage

### Errors

- `ErrWeightExceedsLimit` - Requested weight exceeds semaphore size

## Behavior

### Blocking Acquire

When calling `Acquire()`:
- If resources are available, acquires immediately
- If not enough resources, blocks until available
- Respects context cancellation and timeout
- Multiple goroutines waiting are woken when resources released

### Non-blocking Acquire

When calling `TryAcquire()`:
- Never blocks
- Returns immediately with success/failure
- Does not respect context (no context parameter)

### Weight Limits

- Cannot acquire weight greater than semaphore size
- Returns `ErrWeightExceedsLimit` if weight too large
- Weight must be positive (not enforced, but undefined behavior if negative)

### Release Behavior

- Releasing more than acquired resets current to 0
- Broadcasts to all waiting goroutines
- Goroutines compete for available resources

### Thread Safety

All operations are thread-safe:
- Safe to acquire/release from multiple goroutines
- Uses mutex and condition variable for synchronization

## Comparison with Standard Library

Go's `golang.org/x/sync/semaphore` is similar but has some differences:

| Feature | chesemaphore | x/sync/semaphore |
|---------|-------------|------------------|
| Weighted | ✓ | ✓ |
| Context support | ✓ | ✓ |
| TryAcquire | ✓ | ✓ |
| Available() | ✓ | ✗ |
| Current() | ✓ | ✗ |
| Size() | ✓ | ✗ |

chesemaphore provides additional introspection methods that can be useful for monitoring and debugging.

## License

MIT
