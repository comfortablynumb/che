# cheworker - Worker Pool

Concurrent worker pool for processing jobs with graceful shutdown support.

## Features

- **Fixed Worker Pool**: Configure number of concurrent workers
- **Job Queue**: Buffered channel for job queuing
- **Graceful Shutdown**: Wait for jobs to complete or timeout
- **Context Support**: Cancel jobs via context
- **Error Collection**: Gather all errors from job execution
- **Panic Recovery**: Automatic panic recovery with custom handlers
- **Zero Dependencies**: Only uses Go standard library

## Installation

```bash
go get github.com/comfortablynumb/che/pkg/cheworker
```

## Quick Start

```go
package main

import (
    "context"
    "fmt"
    "github.com/comfortablynumb/che/pkg/cheworker"
)

func main() {
    // Create pool with 5 workers
    pool := cheworker.New(&cheworker.Config{
        Workers:   5,
        QueueSize: 100,
    })

    // Start the pool
    pool.Start()

    // Submit jobs
    for i := 0; i < 10; i++ {
        id := i
        pool.Submit(func(ctx context.Context) error {
            fmt.Printf("Processing job %d\n", id)
            return nil
        })
    }

    // Gracefully shutdown
    pool.Shutdown()
}
```

## Usage

### Creating a Pool

```go
// Simple pool with defaults (10 workers, queue size 100)
pool := cheworker.New(nil)

// Custom configuration
pool := cheworker.New(&cheworker.Config{
    Workers:   20,      // Number of concurrent workers
    QueueSize: 500,     // Job queue buffer size
    OnError: func(err error) {
        log.Printf("Job error: %v", err)
    },
})

pool.Start()
```

### Submitting Jobs

```go
// Submit a job
err := pool.Submit(func(ctx context.Context) error {
    // Do work here
    return nil
})
if err != nil {
    log.Printf("Failed to submit: %v", err)
}

// Submit with custom context
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

err = pool.SubmitWithContext(ctx, func(jobCtx context.Context) error {
    // Job will be cancelled if ctx times out
    select {
    case <-jobCtx.Done():
        return jobCtx.Err()
    default:
        // Do work
        return nil
    }
})
```

### Shutdown Options

```go
// Graceful shutdown - waits for all jobs to complete
pool.Shutdown()

// Shutdown with timeout
ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
defer cancel()

err := pool.ShutdownWithContext(ctx)
if err != nil {
    log.Printf("Shutdown timed out: %v", err)
}

// Immediate stop - cancels running jobs
pool.Stop()
```

### Error Handling

```go
// Collect errors after shutdown
pool.Shutdown()
errors := pool.Errors()
for _, err := range errors {
    log.Printf("Job error: %v", err)
}

// Or use callback for real-time error handling
pool := cheworker.New(&cheworker.Config{
    Workers: 5,
    OnError: func(err error) {
        // Called immediately when job returns error
        log.Printf("Job failed: %v", err)
    },
})
```

### Panic Recovery

```go
// Default: panics are converted to errors
pool := cheworker.New(&cheworker.Config{
    Workers: 5,
})

pool.Submit(func(ctx context.Context) error {
    panic("something went wrong")
    // Panic is recovered and added to errors
})

// Custom panic handler
pool := cheworker.New(&cheworker.Config{
    Workers: 5,
    PanicHandler: func(p interface{}) {
        log.Printf("Worker panic: %v", p)
        // Send alert, etc.
    },
})
```

### Pool Information

```go
workerCount := pool.WorkerCount()  // Number of workers
queueSize := pool.QueueSize()      // Queue buffer capacity
pending := pool.PendingJobs()      // Jobs waiting in queue
```

## Examples

### Batch Processing

```go
type Item struct {
    ID   int
    Data string
}

func ProcessBatch(items []Item) error {
    pool := cheworker.New(&cheworker.Config{
        Workers:   10,
        QueueSize: len(items),
    })
    pool.Start()

    for _, item := range items {
        item := item // capture for closure
        err := pool.Submit(func(ctx context.Context) error {
            return processItem(ctx, item)
        })
        if err != nil {
            pool.Stop()
            return fmt.Errorf("failed to submit item %d: %w", item.ID, err)
        }
    }

    pool.Shutdown()

    // Check for errors
    if errs := pool.Errors(); len(errs) > 0 {
        return fmt.Errorf("%d items failed processing", len(errs))
    }

    return nil
}
```

### Web Scraper

```go
func ScrapeURLs(urls []string) ([]Result, error) {
    pool := cheworker.New(&cheworker.Config{
        Workers:   5, // Limit concurrent requests
        QueueSize: len(urls),
    })
    pool.Start()

    var mu sync.Mutex
    results := make([]Result, 0, len(urls))

    for _, url := range urls {
        url := url
        pool.Submit(func(ctx context.Context) error {
            result, err := scrapeURL(ctx, url)
            if err != nil {
                return err
            }

            mu.Lock()
            results = append(results, result)
            mu.Unlock()

            return nil
        })
    }

    pool.Shutdown()

    if errs := pool.Errors(); len(errs) > 0 {
        return results, fmt.Errorf("%d URLs failed", len(errs))
    }

    return results, nil
}
```

### Image Processing

```go
func ProcessImages(inputDir, outputDir string) error {
    files, err := os.ReadDir(inputDir)
    if err != nil {
        return err
    }

    pool := cheworker.New(&cheworker.Config{
        Workers:   runtime.NumCPU(),
        QueueSize: len(files),
        OnError: func(err error) {
            log.Printf("Image processing error: %v", err)
        },
    })
    pool.Start()

    for _, file := range files {
        if file.IsDir() {
            continue
        }

        filename := file.Name()
        pool.Submit(func(ctx context.Context) error {
            return processImage(
                filepath.Join(inputDir, filename),
                filepath.Join(outputDir, filename),
            )
        })
    }

    // Shutdown with timeout
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
    defer cancel()

    return pool.ShutdownWithContext(ctx)
}
```

### Database Migration

```go
func MigrateRecords(db *sql.DB, limit int) error {
    pool := cheworker.New(&cheworker.Config{
        Workers:   5, // Limit DB connections
        QueueSize: 1000,
    })
    pool.Start()

    offset := 0
    batchSize := 100

    for {
        rows, err := db.Query(
            "SELECT id, data FROM old_table LIMIT ? OFFSET ?",
            batchSize, offset,
        )
        if err != nil {
            pool.Stop()
            return err
        }

        count := 0
        for rows.Next() {
            var id int
            var data string
            rows.Scan(&id, &data)

            id, data := id, data
            pool.Submit(func(ctx context.Context) error {
                return migrateRecord(db, id, data)
            })
            count++
        }
        rows.Close()

        if count == 0 {
            break
        }

        offset += batchSize

        if offset >= limit {
            break
        }
    }

    pool.Shutdown()

    if errs := pool.Errors(); len(errs) > 0 {
        return fmt.Errorf("%d records failed migration", len(errs))
    }

    return nil
}
```

### API Rate Limiting

```go
type RateLimitedPool struct {
    pool    *cheworker.Pool
    limiter *rate.Limiter
}

func NewRateLimitedPool(requestsPerSecond int, workers int) *RateLimitedPool {
    return &RateLimitedPool{
        pool: cheworker.New(&cheworker.Config{
            Workers:   workers,
            QueueSize: 100,
        }),
        limiter: rate.NewLimiter(rate.Limit(requestsPerSecond), requestsPerSecond),
    }
}

func (p *RateLimitedPool) Start() {
    p.pool.Start()
}

func (p *RateLimitedPool) SubmitAPICall(ctx context.Context, call func(context.Context) error) error {
    // Wait for rate limit
    if err := p.limiter.Wait(ctx); err != nil {
        return err
    }

    return p.pool.Submit(call)
}

func (p *RateLimitedPool) Shutdown() {
    p.pool.Shutdown()
}
```

### Integration with chesignal

```go
import (
    "github.com/comfortablynumb/che/pkg/cheworker"
    "github.com/comfortablynumb/che/pkg/chesignal"
)

func main() {
    pool := cheworker.New(&cheworker.Config{
        Workers:   10,
        QueueSize: 100,
    })
    pool.Start()

    // Submit background jobs
    go func() {
        for {
            pool.Submit(func(ctx context.Context) error {
                // Do work
                return nil
            })
            time.Sleep(time.Second)
        }
    }()

    // Graceful shutdown on signal
    config := &chesignal.Config{
        Signals: []os.Signal{os.Interrupt, syscall.SIGTERM},
        Timeout: 30 * time.Second,
    }

    chesignal.WaitForShutdown(config, func(ctx context.Context) error {
        log.Println("Shutting down worker pool...")
        return pool.ShutdownWithContext(ctx)
    })
}
```

## API Reference

### Types

```go
type Pool struct { ... }

type Config struct {
    Workers      int             // Number of workers (default: 10)
    QueueSize    int             // Queue buffer size (default: 100)
    OnError      func(error)     // Error callback
    PanicHandler func(interface{}) // Panic handler
}

type Job func(context.Context) error
```

### Methods

- `New(config *Config) *Pool` - Create new pool
- `Start()` - Start workers
- `Submit(job Job) error` - Submit job
- `SubmitWithContext(ctx context.Context, job Job) error` - Submit with context
- `Shutdown()` - Graceful shutdown
- `ShutdownWithContext(ctx context.Context) error` - Shutdown with timeout
- `Stop()` - Immediate stop
- `Errors() []error` - Get all errors
- `WorkerCount() int` - Get worker count
- `QueueSize() int` - Get queue capacity
- `PendingJobs() int` - Get pending job count

## Best Practices

1. **Right-size the pool**: Match workers to available CPU cores for CPU-bound tasks
2. **Set appropriate queue size**: Prevent memory issues with very large job counts
3. **Handle errors**: Always check `pool.Errors()` after shutdown
4. **Use contexts**: Submit jobs with contexts for cancellation support
5. **Graceful shutdown**: Prefer `Shutdown()` over `Stop()` to avoid data loss
6. **Monitor pending jobs**: Use `PendingJobs()` to detect backlogs

## Performance Considerations

- Workers are goroutines, so overhead is minimal
- Queue is a buffered channel - O(1) enqueue/dequeue
- Context cancellation has minimal overhead
- Error collection uses mutex - consider using `OnError` callback for high-error scenarios

## Related Packages

- **[chesignal](../chesignal)** - Graceful shutdown utilities
- **[chectx](../chectx)** - Type-safe context utilities

## License

This package is part of the Che library and shares the same license.
