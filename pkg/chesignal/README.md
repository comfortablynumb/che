# chesignal - Graceful Shutdown Utilities

Elegant utilities for handling OS signals and performing graceful shutdowns in Go applications.

## Features

- **Signal Handling**: Listen for SIGINT, SIGTERM, and custom signals
- **Graceful Shutdown**: Execute cleanup functions with configurable timeout
- **Context Support**: Programmatic shutdown via context cancellation
- **Callbacks**: Hooks for shutdown lifecycle events
- **Ordered Execution**: Shutdown functions execute in order with error handling
- **Zero Dependencies**: Only uses Go standard library

## Installation

```bash
go get github.com/comfortablynumb/che/pkg/chesignal
```

## Quick Start

```go
package main

import (
    "context"
    "fmt"
    "net/http"
    "time"

    "github.com/comfortablynumb/che/pkg/chesignal"
)

func main() {
    // Start your application
    server := &http.Server{Addr: ":8080"}
    go server.ListenAndServe()

    fmt.Println("Server started on :8080")
    fmt.Println("Press Ctrl+C to shutdown gracefully")

    // Wait for shutdown signal and cleanup
    chesignal.WaitForShutdown(nil,
        func(ctx context.Context) error {
            fmt.Println("Shutting down server...")
            return server.Shutdown(ctx)
        },
        func(ctx context.Context) error {
            fmt.Println("Closing database connections...")
            // Close database, etc.
            return nil
        },
    )

    fmt.Println("Shutdown complete")
}
```

## Usage

### Basic Shutdown

Use `WaitForShutdown` to wait for SIGINT or SIGTERM:

```go
err := chesignal.WaitForShutdown(nil,
    func(ctx context.Context) error {
        fmt.Println("Cleaning up...")
        return nil
    },
)
if err != nil {
    log.Fatal(err)
}
```

### Custom Configuration

Configure signals, timeout, and callbacks:

```go
config := &chesignal.Config{
    Signals: []os.Signal{syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP},
    Timeout: 15 * time.Second,
    OnShutdownStart: func() {
        fmt.Println("Shutdown initiated...")
    },
    OnShutdownComplete: func() {
        fmt.Println("Shutdown completed successfully")
    },
    OnShutdownTimeout: func() {
        fmt.Println("Warning: Shutdown timed out!")
    },
}

err := chesignal.WaitForShutdown(config, shutdownFuncs...)
```

### Multiple Shutdown Functions

Shutdown functions execute in order. If one fails, execution stops:

```go
chesignal.WaitForShutdown(nil,
    func(ctx context.Context) error {
        // First: stop accepting new requests
        return httpServer.Shutdown(ctx)
    },
    func(ctx context.Context) error {
        // Second: drain message queues
        return drainQueues(ctx)
    },
    func(ctx context.Context) error {
        // Third: close database connections
        return db.Close()
    },
)
```

### Context-Aware Shutdown

Use `WaitForShutdownWithContext` for programmatic shutdown:

```go
ctx, cancel := context.WithCancel(context.Background())

// Start shutdown listener in goroutine
go func() {
    chesignal.WaitForShutdownWithContext(ctx, config, shutdownFuncs...)
}()

// Trigger shutdown programmatically
cancel()
```

### Custom Signal Channel

For custom signal handling logic:

```go
sigChan := chesignal.NotifyOnSignal(syscall.SIGINT, syscall.SIGTERM)

select {
case sig := <-sigChan:
    fmt.Printf("Received signal: %v\n", sig)
    // Perform custom shutdown logic
}
```

## Examples

### HTTP Server with Graceful Shutdown

```go
package main

import (
    "context"
    "fmt"
    "log"
    "net/http"
    "time"

    "github.com/comfortablynumb/che/pkg/chesignal"
)

func main() {
    mux := http.NewServeMux()
    mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        time.Sleep(2 * time.Second) // Simulate work
        fmt.Fprintf(w, "Hello, World!")
    })

    server := &http.Server{
        Addr:    ":8080",
        Handler: mux,
    }

    // Start server
    go func() {
        log.Println("Server starting on :8080")
        if err := server.ListenAndServe(); err != http.ErrServerClosed {
            log.Fatal(err)
        }
    }()

    // Wait for shutdown signal
    config := &chesignal.Config{
        Timeout: 30 * time.Second,
        OnShutdownStart: func() {
            log.Println("Shutdown signal received, stopping server...")
        },
    }

    if err := chesignal.WaitForShutdown(config,
        func(ctx context.Context) error {
            // Gracefully shutdown the server
            return server.Shutdown(ctx)
        },
    ); err != nil {
        log.Fatal(err)
    }

    log.Println("Server stopped")
}
```

### Database and Cache Cleanup

```go
package main

import (
    "context"
    "database/sql"
    "log"
    "time"

    "github.com/comfortablynumb/che/pkg/chesignal"
)

type App struct {
    db    *sql.DB
    cache CacheClient
}

func (app *App) Shutdown() error {
    config := &chesignal.Config{
        Timeout: 15 * time.Second,
        OnShutdownStart: func() {
            log.Println("Starting graceful shutdown...")
        },
        OnShutdownComplete: func() {
            log.Println("All resources cleaned up successfully")
        },
    }

    return chesignal.WaitForShutdown(config,
        func(ctx context.Context) error {
            log.Println("Closing cache connections...")
            return app.cache.Close()
        },
        func(ctx context.Context) error {
            log.Println("Closing database connections...")
            return app.db.Close()
        },
    )
}
```

### Worker Pool Shutdown

```go
package main

import (
    "context"
    "log"
    "sync"
    "time"

    "github.com/comfortablynumb/che/pkg/chesignal"
)

type WorkerPool struct {
    workers int
    jobs    chan Job
    wg      sync.WaitGroup
    stop    chan struct{}
}

func (wp *WorkerPool) Shutdown(ctx context.Context) error {
    close(wp.stop) // Signal workers to stop

    // Wait for workers with context timeout
    done := make(chan struct{})
    go func() {
        wp.wg.Wait()
        close(done)
    }()

    select {
    case <-done:
        log.Println("All workers stopped")
        return nil
    case <-ctx.Done():
        log.Println("Timeout waiting for workers")
        return ctx.Err()
    }
}

func main() {
    pool := NewWorkerPool(10)
    pool.Start()

    config := &chesignal.Config{
        Timeout: 30 * time.Second,
    }

    chesignal.WaitForShutdown(config,
        func(ctx context.Context) error {
            return pool.Shutdown(ctx)
        },
    )
}
```

### Health Check Server

```go
package main

import (
    "context"
    "fmt"
    "log"
    "net/http"
    "sync"
    "time"

    "github.com/comfortablynumb/che/pkg/chesignal"
)

type HealthCheckServer struct {
    server  *http.Server
    healthy bool
    mu      sync.RWMutex
}

func (h *HealthCheckServer) SetHealthy(healthy bool) {
    h.mu.Lock()
    defer h.mu.Unlock()
    h.healthy = healthy
}

func (h *HealthCheckServer) Start() {
    mux := http.NewServeMux()
    mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
        h.mu.RLock()
        defer h.mu.RUnlock()

        if h.healthy {
            w.WriteHeader(http.StatusOK)
            fmt.Fprint(w, "OK")
        } else {
            w.WriteHeader(http.StatusServiceUnavailable)
            fmt.Fprint(w, "Shutting down")
        }
    })

    h.server = &http.Server{Addr: ":8081", Handler: mux}
    go h.server.ListenAndServe()
}

func (h *HealthCheckServer) Shutdown(ctx context.Context) error {
    // Mark as unhealthy first
    h.SetHealthy(false)

    // Wait a bit for load balancers to detect unhealthy state
    time.Sleep(2 * time.Second)

    // Now shutdown the server
    return h.server.Shutdown(ctx)
}

func main() {
    healthCheck := &HealthCheckServer{healthy: true}
    healthCheck.Start()

    // Main application
    app := startApp()

    config := &chesignal.Config{
        Timeout: 30 * time.Second,
    }

    chesignal.WaitForShutdown(config,
        func(ctx context.Context) error {
            // Shutdown health check server first
            return healthCheck.Shutdown(ctx)
        },
        func(ctx context.Context) error {
            // Then shutdown main app
            return app.Shutdown(ctx)
        },
    )

    log.Println("All services stopped")
}
```

## Configuration

### Config struct

```go
type Config struct {
    // Signals to listen for (defaults to SIGINT and SIGTERM)
    Signals []os.Signal

    // Timeout for graceful shutdown (defaults to 30 seconds)
    Timeout time.Duration

    // OnShutdownStart is called when shutdown begins
    OnShutdownStart func()

    // OnShutdownComplete is called when shutdown completes successfully
    OnShutdownComplete func()

    // OnShutdownTimeout is called if shutdown times out
    OnShutdownTimeout func()
}
```

### Default Configuration

```go
config := chesignal.DefaultConfig()
// Listens for: SIGINT, SIGTERM
// Timeout: 30 seconds
// No callbacks
```

## API Reference

### Functions

- `WaitForShutdown(config, ...shutdownFuncs) error` - Wait for signal and execute shutdown functions
- `WaitForShutdownWithContext(ctx, config, ...shutdownFuncs) error` - Wait for signal or context cancellation
- `NotifyOnSignal(...signals) <-chan os.Signal` - Create a signal notification channel
- `DefaultConfig() *Config` - Returns default configuration

### Types

- `ShutdownFunc func(context.Context) error` - Function to execute on shutdown
- `Config` - Configuration for shutdown behavior

## Best Practices

1. **Order Matters**: Shutdown functions execute in order - structure them logically
2. **Set Appropriate Timeouts**: Choose timeout based on your longest-running operation
3. **Use Context**: Respect the context deadline in shutdown functions
4. **Handle Errors**: Check and handle errors from shutdown functions
5. **Stop Accepting First**: Stop accepting new work before cleanup
6. **Drain Gracefully**: Allow in-flight requests to complete before shutdown
7. **Close External Connections Last**: Close database/cache connections after internal cleanup

## Common Patterns

### Stop Accept → Drain → Close

```go
chesignal.WaitForShutdown(config,
    func(ctx context.Context) error {
        // 1. Stop accepting new requests
        return server.Shutdown(ctx)
    },
    func(ctx context.Context) error {
        // 2. Drain work queues
        return queue.Drain(ctx)
    },
    func(ctx context.Context) error {
        // 3. Close connections
        return db.Close()
    },
)
```

### Parallel Cleanup

```go
func(ctx context.Context) error {
    var wg sync.WaitGroup
    errors := make(chan error, 2)

    // Close multiple resources in parallel
    wg.Add(2)
    go func() {
        defer wg.Done()
        errors <- cache.Close()
    }()
    go func() {
        defer wg.Done()
        errors <- queue.Close()
    }()

    wg.Wait()
    close(errors)

    // Check for errors
    for err := range errors {
        if err != nil {
            return err
        }
    }
    return nil
}
```

## Related Packages

- **[chectx](../chectx)** - Type-safe context utilities
- **[chehttp](../chehttp)** - Ergonomic HTTP client and server utilities

## License

This package is part of the Che library and shares the same license.
