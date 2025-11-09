# chebatch

Batch processing utilities for efficient data aggregation and processing in Go.

## Features

- **Batcher**: Automatic batching with size and time limits
- **Group**: Split slices into fixed-size batches
- **Process**: Sequential batch processing
- **ProcessParallel**: Parallel batch processing with concurrency control
- Generic support for any type
- Context support for cancellation
- Thread-safe operations
- Zero external dependencies

## Installation

```bash
go get github.com/comfortablynumb/che/pkg/chebatch
```

## Usage

### Auto Batcher

Automatically batch items based on size or time limits:

```go
package main

import (
    "context"
    "fmt"
    "github.com/comfortablynumb/che/pkg/chebatch"
)

func main() {
    // Create batcher that processes batches of 10 or after 1 second
    b := chebatch.NewBatcher(
        func(ctx context.Context, items []string) error {
            fmt.Printf("Processing batch of %d items\n", len(items))
            return nil
        },
        chebatch.WithMaxSize[string](10),
        chebatch.WithMaxWait[string](1*time.Second),
    )
    defer b.Close()

    // Add items - automatically batched
    for i := 0; i < 25; i++ {
        b.Add(fmt.Sprintf("item-%d", i))
    }

    // Flush any remaining items
    b.Flush()
}
```

### Grouping

Split a slice into fixed-size batches:

```go
items := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

// Group into batches of 3
batches := chebatch.Group(items, 3)
// Result: [[1,2,3], [4,5,6], [7,8,9], [10]]

for i, batch := range batches {
    fmt.Printf("Batch %d: %v\n", i, batch)
}
```

### Sequential Processing

Process items in batches sequentially:

```go
import "context"

items := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

ctx := context.Background()
err := chebatch.Process(ctx, items, 3, func(ctx context.Context, batch []int) error {
    fmt.Printf("Processing: %v\n", batch)
    // Process batch
    return nil
})

if err != nil {
    log.Fatal(err)
}
```

### Parallel Processing

Process batches in parallel with concurrency control:

```go
items := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

ctx := context.Background()
err := chebatch.ProcessParallel(
    ctx,
    items,
    2,  // batch size
    3,  // max concurrent batches
    func(ctx context.Context, batch []int) error {
        // This runs in parallel (up to 3 at a time)
        return processItems(batch)
    },
)
```

## Examples

### Database Bulk Insert

```go
type User struct {
    ID   int
    Name string
}

func main() {
    db := connectDB()

    // Create batcher for bulk inserts
    batcher := chebatch.NewBatcher(
        func(ctx context.Context, users []User) error {
            return db.BulkInsert(users)
        },
        chebatch.WithMaxSize[User](100),
        chebatch.WithMaxWait[User](5*time.Second),
    )
    defer batcher.Close()

    // Stream users from API
    for user := range getUsersFromAPI() {
        batcher.Add(user)
    }
}
```

### Log Aggregation

```go
type LogEntry struct {
    Timestamp time.Time
    Level     string
    Message   string
}

type LogAggregator struct {
    batcher *chebatch.Batcher[LogEntry]
}

func NewLogAggregator() *LogAggregator {
    return &LogAggregator{
        batcher: chebatch.NewBatcher(
            func(ctx context.Context, logs []LogEntry) error {
                return writeLogsToFile(logs)
            },
            chebatch.WithMaxSize[LogEntry](1000),
            chebatch.WithMaxWait[LogEntry](10*time.Second),
        ),
    }
}

func (la *LogAggregator) Log(level, message string) {
    la.batcher.Add(LogEntry{
        Timestamp: time.Now(),
        Level:     level,
        Message:   message,
    })
}

func (la *LogAggregator) Close() error {
    return la.batcher.Close()
}

func main() {
    logger := NewLogAggregator()
    defer logger.Close()

    logger.Log("INFO", "Application started")
    logger.Log("DEBUG", "Debug message")
    logger.Log("ERROR", "An error occurred")
}
```

### Metrics Collection

```go
type Metric struct {
    Name  string
    Value float64
    Tags  map[string]string
}

type MetricsCollector struct {
    batcher *chebatch.Batcher[Metric]
}

func NewMetricsCollector() *MetricsCollector {
    return &MetricsCollector{
        batcher: chebatch.NewBatcher(
            func(ctx context.Context, metrics []Metric) error {
                return sendToMonitoringService(metrics)
            },
            chebatch.WithMaxSize[Metric](50),
            chebatch.WithMaxWait[Metric](30*time.Second),
        ),
    }
}

func (mc *MetricsCollector) Record(name string, value float64, tags map[string]string) {
    mc.batcher.Add(Metric{
        Name:  name,
        Value: value,
        Tags:  tags,
    })
}

func (mc *MetricsCollector) Close() error {
    return mc.batcher.Close()
}

func main() {
    collector := NewMetricsCollector()
    defer collector.Close()

    collector.Record("http.requests", 1, map[string]string{
        "method": "GET",
        "status": "200",
    })
}
```

### File Processing

```go
func processLargeFile(filename string) error {
    lines, err := readAllLines(filename)
    if err != nil {
        return err
    }

    ctx := context.Background()
    return chebatch.Process(ctx, lines, 100, func(ctx context.Context, batch []string) error {
        // Process 100 lines at a time
        for _, line := range batch {
            processLine(line)
        }
        return nil
    })
}
```

### API Rate Limiting

```go
type APIRequest struct {
    URL     string
    Payload interface{}
}

func sendBatchedRequests(requests []APIRequest) error {
    batcher := chebatch.NewBatcher(
        func(ctx context.Context, batch []APIRequest) error {
            // Send batch of requests
            return sendToAPI(batch)
        },
        chebatch.WithMaxSize[APIRequest](10),
        chebatch.WithMaxWait[APIRequest](2*time.Second),
    )
    defer batcher.Close()

    for _, req := range requests {
        if err := batcher.Add(req); err != nil {
            return err
        }
    }

    return batcher.Flush()
}
```

### Data Migration

```go
func migrateData(sourceDB, targetDB *sql.DB) error {
    rows, err := sourceDB.Query("SELECT * FROM users")
    if err != nil {
        return err
    }
    defer rows.Close()

    batcher := chebatch.NewBatcher(
        func(ctx context.Context, users []User) error {
            return targetDB.BulkInsert(users)
        },
        chebatch.WithMaxSize[User](500),
    )
    defer batcher.Close()

    for rows.Next() {
        var user User
        if err := rows.Scan(&user.ID, &user.Name); err != nil {
            return err
        }
        batcher.Add(user)
    }

    return batcher.Flush()
}
```

### Parallel File Upload

```go
type File struct {
    Path string
    Data []byte
}

func uploadFiles(files []File) error {
    ctx := context.Background()

    // Upload in batches of 5, with 3 concurrent uploads
    return chebatch.ProcessParallel(
        ctx,
        files,
        5,  // 5 files per batch
        3,  // 3 concurrent uploads
        func(ctx context.Context, batch []File) error {
            for _, file := range batch {
                if err := uploadFile(file); err != nil {
                    return err
                }
            }
            return nil
        },
    )
}
```

### Event Notification Batching

```go
type Event struct {
    Type string
    Data interface{}
}

type EventNotifier struct {
    batcher *chebatch.Batcher[Event]
}

func NewEventNotifier() *EventNotifier {
    return &EventNotifier{
        batcher: chebatch.NewBatcher(
            func(ctx context.Context, events []Event) error {
                // Send batch notification
                return notifySubscribers(events)
            },
            chebatch.WithMaxSize[Event](20),
            chebatch.WithMaxWait[Event](5*time.Second),
        ),
    }
}

func (en *EventNotifier) Notify(eventType string, data interface{}) {
    en.batcher.Add(Event{
        Type: eventType,
        Data: data,
    })
}

func (en *EventNotifier) Close() error {
    return en.batcher.Close()
}
```

### Elasticsearch Bulk Indexing

```go
type Document struct {
    ID      string
    Content string
}

func indexDocuments(docs []Document) error {
    batcher := chebatch.NewBatcher(
        func(ctx context.Context, batch []Document) error {
            return esBulkIndex(batch)
        },
        chebatch.WithMaxSize[Document](1000),
        chebatch.WithMaxWait[Document](10*time.Second),
    )
    defer batcher.Close()

    for _, doc := range docs {
        if err := batcher.Add(doc); err != nil {
            return err
        }
    }

    return batcher.Flush()
}
```

## API Reference

### Batcher

#### Creating
```go
NewBatcher[T any](processor Processor[T], opts ...BatcherOption[T]) *Batcher[T]
```

#### Options
- `WithMaxSize[T](size int)` - Set maximum batch size (default: 100)
- `WithMaxWait[T](duration time.Duration)` - Set maximum wait time (default: 1s)

#### Methods
- `Add(item T) error` - Add item to batch
- `Flush() error` - Process pending items immediately
- `Close() error` - Flush and close batcher
- `Size() int` - Get number of pending items
- `Start()` - Start background processing

### Functions

#### Group
```go
Group[T any](items []T, size int) [][]T
```
Split slice into batches of specified size.

#### Process
```go
Process[T any](ctx context.Context, items []T, batchSize int, fn Processor[T]) error
```
Process items in batches sequentially.

#### ProcessParallel
```go
ProcessParallel[T any](ctx context.Context, items []T, batchSize int, maxConcurrent int, fn Processor[T]) error
```
Process items in batches with parallel execution.

### Processor Type
```go
type Processor[T any] func(ctx context.Context, items []T) error
```

## Behavior

### Batcher Triggers

A batch is processed when:
1. **Size limit reached**: `len(items) >= maxSize`
2. **Time limit reached**: `time.Since(lastAdd) >= maxWait`
3. **Explicit flush**: `Flush()` called
4. **Batcher closed**: `Close()` called

### Error Handling

- `Process`: Stops on first error, returns error
- `ProcessParallel`: Cancels remaining batches on first error
- `Batcher`: Processes in background, errors are not returned

### Context Cancellation

All processing functions respect context cancellation:
- Stops processing remaining batches
- Returns context error
- Cleans up resources

### Thread Safety

- `Batcher.Add()` is thread-safe
- Multiple goroutines can add items concurrently
- Processing happens in separate goroutines

## Best Practices

1. **Always close batchers**:
   ```go
   defer batcher.Close()
   ```

2. **Choose appropriate batch sizes**:
   - Database: 100-1000 items
   - API calls: 10-100 items
   - File operations: 50-500 items

3. **Use context for cancellation**:
   ```go
   ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
   defer cancel()
   ```

4. **Handle errors in processors**:
   ```go
   processor := func(ctx context.Context, items []T) error {
       if err := process(items); err != nil {
           log.Printf("Batch processing failed: %v", err)
           return err
       }
       return nil
   }
   ```

5. **Flush before shutdown**:
   ```go
   func shutdown() {
       batcher.Flush()
       batcher.Close()
   }
   ```

## Performance

### Benchmarks

- `Batcher.Add()`: ~50ns per operation
- `Group()`: O(n) time complexity
- `Process()`: Sequential, predictable performance
- `ProcessParallel()`: Can achieve near-linear speedup with proper concurrency

### Memory Usage

- Batcher holds pending items in memory
- Consider memory constraints when setting `maxSize`
- Use streaming for very large datasets

## License

MIT
