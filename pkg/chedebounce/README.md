# chedebounce

Debounce and throttle utilities for rate-limiting function calls in Go.

## Features

- **Debouncer**: Delays function execution until after a period of inactivity
- **Throttler**: Limits function execution to at most once per interval
- Leading and trailing edge options for throttling
- Flush and cancel operations
- Thread-safe concurrent access
- Zero external dependencies

## Installation

```bash
go get github.com/comfortablynumb/che/pkg/chedebounce
```

## Usage

### Debouncing

Debouncing delays function execution until after a specified period has elapsed since the last call.

```go
package main

import (
    "fmt"
    "time"
    "github.com/comfortablynumb/che/pkg/chedebounce"
)

func main() {
    d := chedebounce.NewDebouncer(200 * time.Millisecond)
    defer d.Close()

    // Call multiple times quickly
    for i := 0; i < 5; i++ {
        d.Call(func() {
            fmt.Println("Executed!")
        })
        time.Sleep(50 * time.Millisecond)
    }

    // Only executes once, 200ms after the last call
    time.Sleep(300 * time.Millisecond)
}
```

### Throttling

Throttling limits function execution to at most once per specified interval.

```go
th := chedebounce.NewThrottler(1 * time.Second)
defer th.Close()

// First call executes immediately
th.Call(func() {
    fmt.Println("Executed!")
})

// Subsequent calls within 1 second are ignored
for i := 0; i < 10; i++ {
    th.Call(func() {
        fmt.Println("This won't execute")
    })
    time.Sleep(100 * time.Millisecond)
}
```

### Function Wrappers

Create debounced or throttled versions of functions:

```go
// Debounced function
original := func() {
    fmt.Println("Searching...")
}
debounced := chedebounce.Debounce(300*time.Millisecond, original)

// Use like a normal function
debounced()
debounced()
debounced() // Only the last call executes

// Throttled function
throttled := chedebounce.Throttle(1*time.Second, original)
throttled() // Executes
throttled() // Ignored
throttled() // Ignored
```

### Throttling Options

```go
// Leading edge only (default)
th := chedebounce.NewThrottler(1*time.Second, chedebounce.WithLeading())

// Trailing edge only
th := chedebounce.NewThrottler(1*time.Second, chedebounce.WithTrailing())

// Both leading and trailing
th := chedebounce.NewThrottler(1*time.Second,
    chedebounce.WithLeading(),
    chedebounce.WithTrailing())
```

### Flush and Cancel

```go
d := chedebounce.NewDebouncer(1 * time.Second)

d.Call(func() {
    fmt.Println("Hello")
})

// Execute immediately instead of waiting
d.Flush()

// Or cancel without executing
d.Cancel()
```

## Examples

### Search Input Debouncing

```go
type SearchBox struct {
    debouncer *chedebounce.Debouncer
}

func NewSearchBox() *SearchBox {
    return &SearchBox{
        debouncer: chedebounce.NewDebouncer(300 * time.Millisecond),
    }
}

func (sb *SearchBox) OnInput(query string) {
    sb.debouncer.Call(func() {
        results := performSearch(query)
        displayResults(results)
    })
}

func (sb *SearchBox) Close() {
    sb.debouncer.Close()
}

func main() {
    searchBox := NewSearchBox()
    defer searchBox.Close()

    // User typing quickly
    searchBox.OnInput("h")
    searchBox.OnInput("he")
    searchBox.OnInput("hel")
    searchBox.OnInput("hell")
    searchBox.OnInput("hello")

    // Search only executes 300ms after user stops typing
    time.Sleep(500 * time.Millisecond)
}
```

### Window Resize Handler

```go
type Window struct {
    debouncer *chedebounce.Debouncer
}

func NewWindow() *Window {
    w := &Window{
        debouncer: chedebounce.NewDebouncer(250 * time.Millisecond),
    }
    return w
}

func (w *Window) OnResize(width, height int) {
    w.debouncer.Call(func() {
        fmt.Printf("Window resized to %dx%d\n", width, height)
        w.recalculateLayout(width, height)
    })
}

func (w *Window) recalculateLayout(width, height int) {
    // Expensive layout calculation
    time.Sleep(100 * time.Millisecond)
}

func (w *Window) Close() {
    w.debouncer.Close()
}
```

### API Rate Limiting with Throttle

```go
type APIClient struct {
    throttler *chedebounce.Throttler
}

func NewAPIClient() *APIClient {
    return &APIClient{
        // Limit to 1 request per second
        throttler: chedebounce.NewThrottler(1*time.Second,
            chedebounce.WithLeading()),
    }
}

func (c *APIClient) SendRequest(data interface{}) {
    c.throttler.Call(func() {
        fmt.Println("Sending request:", data)
        // Actual API call
    })
}

func (c *APIClient) Close() {
    c.throttler.Close()
}

func main() {
    client := NewAPIClient()
    defer client.Close()

    // Rapid requests
    for i := 0; i < 10; i++ {
        client.SendRequest(map[string]int{"id": i})
        time.Sleep(100 * time.Millisecond)
    }
    // Only the first request executes
}
```

### Scroll Event Throttling

```go
type ScrollHandler struct {
    throttler *chedebounce.Throttler
}

func NewScrollHandler() *ScrollHandler {
    return &ScrollHandler{
        // Update at most every 100ms
        throttler: chedebounce.NewThrottler(100*time.Millisecond,
            chedebounce.WithLeading(),
            chedebounce.WithTrailing()),
    }
}

func (sh *ScrollHandler) OnScroll(position int) {
    sh.throttler.Call(func() {
        sh.updateUI(position)
    })
}

func (sh *ScrollHandler) updateUI(position int) {
    fmt.Printf("Updating UI at scroll position: %d\n", position)
    // Update visible elements, lazy load images, etc.
}

func (sh *ScrollHandler) Close() {
    sh.throttler.Close()
}
```

### Save File with Debounce

```go
type Editor struct {
    debouncer *chedebounce.Debouncer
    content   string
}

func NewEditor() *Editor {
    return &Editor{
        debouncer: chedebounce.NewDebouncer(2 * time.Second),
    }
}

func (e *Editor) OnChange(newContent string) {
    e.content = newContent

    e.debouncer.Call(func() {
        e.save()
    })
}

func (e *Editor) save() {
    fmt.Println("Saving file...")
    // Write to disk
}

func (e *Editor) ForceSave() {
    // Save immediately on explicit save command
    e.debouncer.Flush()
}

func (e *Editor) Close() {
    // Save any pending changes before closing
    e.debouncer.Flush()
    e.debouncer.Close()
}
```

### Button Click Throttling

```go
type Button struct {
    throttler *chedebounce.Throttler
    onClick   func()
}

func NewButton(onClick func()) *Button {
    return &Button{
        onClick: onClick,
        // Prevent double-clicks within 500ms
        throttler: chedebounce.NewThrottler(500*time.Millisecond,
            chedebounce.WithLeading()),
    }
}

func (b *Button) Click() {
    b.throttler.Call(b.onClick)
}

func (b *Button) Close() {
    b.throttler.Close()
}

func main() {
    button := NewButton(func() {
        fmt.Println("Button clicked!")
        // Submit form, make API call, etc.
    })
    defer button.Close()

    // Simulate rapid clicks
    button.Click() // Executes
    button.Click() // Ignored
    button.Click() // Ignored

    time.Sleep(600 * time.Millisecond)
    button.Click() // Executes
}
```

### Log Aggregation

```go
type Logger struct {
    debouncer *chedebounce.Debouncer
    buffer    []string
    mu        sync.Mutex
}

func NewLogger() *Logger {
    l := &Logger{
        debouncer: chedebounce.NewDebouncer(1 * time.Second),
        buffer:    make([]string, 0),
    }
    return l
}

func (l *Logger) Log(message string) {
    l.mu.Lock()
    l.buffer = append(l.buffer, message)
    l.mu.Unlock()

    l.debouncer.Call(func() {
        l.flush()
    })
}

func (l *Logger) flush() {
    l.mu.Lock()
    messages := make([]string, len(l.buffer))
    copy(messages, l.buffer)
    l.buffer = l.buffer[:0]
    l.mu.Unlock()

    // Write all buffered messages at once
    fmt.Printf("Flushing %d log messages\n", len(messages))
    for _, msg := range messages {
        fmt.Println("  -", msg)
    }
}

func (l *Logger) Close() {
    l.debouncer.Flush() // Flush any pending logs
    l.debouncer.Close()
}
```

## API Reference

### Debouncer

#### Creating
- `NewDebouncer(delay time.Duration) *Debouncer` - Create a new debouncer

#### Methods
- `Call(fn func())` - Schedule function to execute after delay
- `Flush()` - Execute pending function immediately
- `Cancel()` - Cancel pending function without executing
- `Close()` - Close debouncer and cancel pending calls

### Throttler

#### Creating
- `NewThrottler(interval time.Duration, opts ...ThrottleOption) *Throttler` - Create a new throttler

#### Options
- `WithLeading()` - Execute on leading edge (enabled by default)
- `WithTrailing()` - Execute on trailing edge (disabled by default)

#### Methods
- `Call(fn func())` - Attempt to call function (respects throttle interval)
- `Flush()` - Execute pending trailing function immediately
- `Cancel()` - Cancel pending trailing function without executing
- `Close()` - Close throttler and cancel pending calls

### Function Wrappers

- `Debounce(delay time.Duration, fn func()) func()` - Return debounced function
- `Throttle(interval time.Duration, fn func(), opts ...ThrottleOption) func()` - Return throttled function

## Debounce vs Throttle

### Debounce
- Delays execution until after quiet period
- Resets timer on each call
- Useful for: search inputs, resize events, auto-save

```
Calls:     | | | |     | | |
           ^           ^
Executes:            X           X
```

### Throttle (Leading Edge)
- Executes immediately, then blocks for interval
- Useful for: scroll handlers, button clicks

```
Calls:     | | | | | | | | | |
           ^     ^     ^
Executes:  X           X           X
```

### Throttle (Trailing Edge)
- Delays execution to end of interval
- Useful for: API calls with rate limits

```
Calls:     | | | | | | | | | |
               ^       ^
Executes:          X           X
```

### Throttle (Both Edges)
- Executes immediately and at end of interval
- Maximum responsiveness with rate limiting

```
Calls:     | | | | | | | | | |
           ^   ^   ^   ^
Executes:  X       X       X       X
```

## Thread Safety

All operations are thread-safe:
- Safe to call from multiple goroutines
- Uses mutex for synchronization
- Timer management is protected

## Best Practices

1. **Always call Close()** when done to prevent resource leaks
2. **Use Flush()** before shutdown to execute pending calls
3. **Choose appropriate delays**:
   - Search: 200-300ms
   - Resize: 150-250ms
   - Auto-save: 1-3s
   - Scroll: 50-100ms

4. **Debounce for user input**, throttle for events
5. **Use defer** to ensure cleanup:
   ```go
   d := chedebounce.NewDebouncer(300 * time.Millisecond)
   defer d.Close()
   ```

## License

MIT
