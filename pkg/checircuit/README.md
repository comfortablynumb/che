# checircuit

Circuit breaker pattern implementation for fault tolerance and resilience.

## Features

- Three states: Closed, Open, Half-Open
- Configurable failure threshold and timeout
- State change callbacks
- Thread-safe
- Zero external dependencies

## Installation

```bash
go get github.com/comfortablynumb/che/pkg/checircuit
```

## Usage

### Basic Example

```go
package main

import (
    "fmt"
    "errors"
    "time"
    "github.com/comfortablynumb/che/pkg/checircuit"
)

func main() {
    cb := checircuit.New(&checircuit.Config{
        MaxFailures: 3,                  // Open after 3 failures
        Timeout:     30 * time.Second,   // Try again after 30s
        MaxRequests: 2,                  // Allow 2 requests in half-open
    })

    // Execute a function through the circuit breaker
    err := cb.Execute(func() error {
        // Your risky operation here
        return callExternalService()
    })

    if err == checircuit.ErrCircuitOpen {
        fmt.Println("Circuit is open, request rejected")
    }
}
```

### With State Change Callback

```go
cb := checircuit.New(&checircuit.Config{
    MaxFailures: 5,
    Timeout:     60 * time.Second,
    OnStateChange: func(from, to checircuit.State) {
        fmt.Printf("Circuit breaker: %v -> %v\n", from, to)
    },
})

// Simulate failures
for i := 0; i < 5; i++ {
    cb.Execute(func() error {
        return errors.New("service unavailable")
    })
}
// Output: Circuit breaker: Closed -> Open

// Wait for timeout
time.Sleep(61 * time.Second)

// Next request transitions to half-open
cb.Execute(func() error {
    return nil
})
// Output: Circuit breaker: Open -> HalfOpen
```

### Monitoring State

```go
// Check current state
state := cb.State()
switch state {
case checircuit.StateClosed:
    fmt.Println("Operating normally")
case checircuit.StateOpen:
    fmt.Println("Circuit is open, blocking requests")
case checircuit.StateHalfOpen:
    fmt.Println("Testing if service recovered")
}

// Check failure count
failures := cb.Failures()
fmt.Printf("Current failures: %d\n", failures)

// Manually reset
cb.Reset() // Forces back to Closed state
```

### Practical Example: HTTP Client

```go
package main

import (
    "fmt"
    "net/http"
    "time"
    "github.com/comfortablynumb/che/pkg/checircuit"
)

type ResilientClient struct {
    client  *http.Client
    breaker *checircuit.Breaker
}

func NewResilientClient() *ResilientClient {
    return &ResilientClient{
        client: &http.Client{Timeout: 5 * time.Second},
        breaker: checircuit.New(&checircuit.Config{
            MaxFailures: 5,
            Timeout:     30 * time.Second,
            MaxRequests: 3,
        }),
    }
}

func (rc *ResilientClient) Get(url string) (*http.Response, error) {
    var resp *http.Response

    err := rc.breaker.Execute(func() error {
        var err error
        resp, err = rc.client.Get(url)
        if err != nil {
            return err
        }

        if resp.StatusCode >= 500 {
            return fmt.Errorf("server error: %d", resp.StatusCode)
        }

        return nil
    })

    return resp, err
}
```

## How It Works

The circuit breaker has three states:

1. **Closed** (Normal operation)
   - Requests pass through normally
   - Failures are counted
   - Opens when failure threshold is reached

2. **Open** (Failing fast)
   - All requests are rejected immediately
   - Returns `ErrCircuitOpen`
   - After timeout, transitions to Half-Open

3. **Half-Open** (Testing recovery)
   - Allows limited number of requests through
   - If requests succeed, closes the circuit
   - If any request fails, reopens the circuit

## Configuration

```go
type Config struct {
    // MaxFailures: Number of failures before opening (default: 5)
    MaxFailures uint

    // Timeout: How long to wait before trying half-open (default: 60s)
    Timeout time.Duration

    // MaxRequests: Max requests in half-open state (default: 1)
    MaxRequests uint

    // OnStateChange: Called when state changes (optional)
    OnStateChange func(from, to State)
}
```

## Errors

- `ErrCircuitOpen` - Returned when circuit is open
- `ErrTooManyRequests` - Returned when half-open has too many requests

## Thread Safety

All operations are thread-safe and can be called concurrently.

## License

MIT
