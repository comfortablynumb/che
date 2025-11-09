# chepubsub

Simple in-memory publish-subscribe event system for Go.

## Features

- Topic-based message routing
- Multiple subscribers per topic
- Asynchronous message delivery
- Synchronous message delivery option
- Buffered message channels
- Thread-safe operations
- Context support for cancellation
- Zero external dependencies

## Installation

```bash
go get github.com/comfortablynumb/che/pkg/chepubsub
```

## Usage

### Basic Pub/Sub

```go
package main

import (
    "fmt"
    "github.com/comfortablynumb/che/pkg/chepubsub"
)

func main() {
    // Create a new pub/sub instance with buffer size of 10
    ps := chepubsub.New(10)
    defer ps.Close()

    // Subscribe to a topic
    sub := ps.Subscribe("events", func(msg chepubsub.Message) {
        fmt.Printf("Received: %v\n", msg.Data)
    })
    defer ps.Unsubscribe(sub)

    // Publish a message
    count := ps.Publish("events", "Hello, World!")
    fmt.Printf("Delivered to %d subscribers\n", count)
}
```

### Multiple Subscribers

```go
ps := chepubsub.New(10)
defer ps.Close()

// First subscriber
sub1 := ps.Subscribe("notifications", func(msg chepubsub.Message) {
    fmt.Println("Subscriber 1:", msg.Data)
})
defer ps.Unsubscribe(sub1)

// Second subscriber
sub2 := ps.Subscribe("notifications", func(msg chepubsub.Message) {
    fmt.Println("Subscriber 2:", msg.Data)
})
defer ps.Unsubscribe(sub2)

// Both subscribers will receive this message
ps.Publish("notifications", "Important update!")
```

### Multiple Topics

```go
ps := chepubsub.New(10)
defer ps.Close()

// Subscribe to different topics
ps.Subscribe("errors", func(msg chepubsub.Message) {
    fmt.Println("Error:", msg.Data)
})

ps.Subscribe("info", func(msg chepubsub.Message) {
    fmt.Println("Info:", msg.Data)
})

// Publish to specific topics
ps.Publish("errors", "Connection failed")
ps.Publish("info", "Server started")
```

### Synchronous Publishing

```go
import "context"

ps := chepubsub.New(10)
defer ps.Close()

ps.Subscribe("events", func(msg chepubsub.Message) {
    // Handler is called synchronously
    processEvent(msg.Data)
})

ctx := context.Background()
err := ps.PublishSync(ctx, "events", "Sync message")
if err != nil {
    fmt.Println("Error:", err)
}
// All handlers have completed when this returns
```

### With Context Cancellation

```go
ps := chepubsub.New(10)
defer ps.Close()

ps.Subscribe("slow", func(msg chepubsub.Message) {
    time.Sleep(5 * time.Second)
    fmt.Println("Processed:", msg.Data)
})

ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
defer cancel()

err := ps.PublishSync(ctx, "slow", "message")
if err == context.DeadlineExceeded {
    fmt.Println("Publishing timed out")
}
```

### Structured Messages

```go
type UserEvent struct {
    UserID string
    Action string
}

ps := chepubsub.New(10)
defer ps.Close()

ps.Subscribe("user.events", func(msg chepubsub.Message) {
    event := msg.Data.(UserEvent)
    fmt.Printf("User %s performed %s\n", event.UserID, event.Action)
})

ps.Publish("user.events", UserEvent{
    UserID: "user123",
    Action: "login",
})
```

### Managing Subscriptions

```go
ps := chepubsub.New(10)
defer ps.Close()

sub := ps.Subscribe("events", func(msg chepubsub.Message) {
    fmt.Println("Received:", msg.Data)
})

// Get topic from subscriber
fmt.Println("Subscribed to:", sub.Topic())

// Unsubscribe when done
ps.Unsubscribe(sub)

// Or unsubscribe all subscribers from a topic
ps.UnsubscribeAll("events")
```

### Introspection

```go
ps := chepubsub.New(10)
defer ps.Close()

ps.Subscribe("topic1", func(msg chepubsub.Message) {})
ps.Subscribe("topic1", func(msg chepubsub.Message) {})
ps.Subscribe("topic2", func(msg chepubsub.Message) {})

// Get all active topics
topics := ps.Topics()
fmt.Println("Active topics:", topics)

// Get subscriber count for a topic
count := ps.SubscriberCount("topic1")
fmt.Printf("topic1 has %d subscribers\n", count)
```

## Examples

### Event Bus

```go
type EventBus struct {
    ps *chepubsub.PubSub
}

func NewEventBus() *EventBus {
    return &EventBus{
        ps: chepubsub.New(100),
    }
}

func (eb *EventBus) On(event string, handler func(data interface{})) {
    eb.ps.Subscribe(event, func(msg chepubsub.Message) {
        handler(msg.Data)
    })
}

func (eb *EventBus) Emit(event string, data interface{}) {
    eb.ps.Publish(event, data)
}

func (eb *EventBus) Close() {
    eb.ps.Close()
}

func main() {
    bus := NewEventBus()
    defer bus.Close()

    bus.On("user.login", func(data interface{}) {
        fmt.Println("User logged in:", data)
    })

    bus.On("user.logout", func(data interface{}) {
        fmt.Println("User logged out:", data)
    })

    bus.Emit("user.login", "alice")
    bus.Emit("user.logout", "bob")
}
```

### Application Events

```go
type App struct {
    events *chepubsub.PubSub
}

func NewApp() *App {
    return &App{
        events: chepubsub.New(50),
    }
}

func (a *App) setupHandlers() {
    // Log all events
    a.events.Subscribe("*", func(msg chepubsub.Message) {
        log.Printf("[%s] %v", msg.Topic, msg.Data)
    })

    // Handle specific events
    a.events.Subscribe("app.started", func(msg chepubsub.Message) {
        fmt.Println("Application started")
    })

    a.events.Subscribe("app.shutdown", func(msg chepubsub.Message) {
        fmt.Println("Application shutting down")
    })

    a.events.Subscribe("error", func(msg chepubsub.Message) {
        fmt.Println("Error occurred:", msg.Data)
    })
}

func (a *App) Start() {
    a.setupHandlers()
    a.events.Publish("app.started", nil)
}

func (a *App) Shutdown() {
    a.events.Publish("app.shutdown", nil)
    a.events.Close()
}
```

### Request/Response Pattern

```go
type Request struct {
    ID       string
    Data     interface{}
    Response chan interface{}
}

func main() {
    ps := chepubsub.New(10)
    defer ps.Close()

    // Handler processes requests and sends responses
    ps.Subscribe("requests", func(msg chepubsub.Message) {
        req := msg.Data.(Request)

        // Process request
        result := processRequest(req.Data)

        // Send response
        req.Response <- result
    })

    // Make a request
    responseChan := make(chan interface{})
    ps.Publish("requests", Request{
        ID:       "req-123",
        Data:     "some data",
        Response: responseChan,
    })

    // Wait for response
    response := <-responseChan
    fmt.Println("Response:", response)
}

func processRequest(data interface{}) interface{} {
    // Process the request
    return "processed: " + data.(string)
}
```

### Fan-out Pattern

```go
func main() {
    ps := chepubsub.New(100)
    defer ps.Close()

    // Multiple workers subscribe to the same topic
    for i := 0; i < 5; i++ {
        workerID := i
        ps.Subscribe("jobs", func(msg chepubsub.Message) {
            job := msg.Data.(string)
            fmt.Printf("Worker %d processing: %s\n", workerID, job)
            time.Sleep(100 * time.Millisecond)
        })
    }

    // Publish jobs
    for i := 0; i < 20; i++ {
        ps.Publish("jobs", fmt.Sprintf("job-%d", i))
    }

    time.Sleep(1 * time.Second)
}
```

### Topic Hierarchies

```go
func main() {
    ps := chepubsub.New(10)
    defer ps.Close()

    // Subscribe to all user events
    ps.Subscribe("user.*", func(msg chepubsub.Message) {
        fmt.Println("User event:", msg.Topic, msg.Data)
    })

    // Subscribe to specific user events
    ps.Subscribe("user.login", func(msg chepubsub.Message) {
        fmt.Println("Login event:", msg.Data)
    })

    ps.Subscribe("user.logout", func(msg chepubsub.Message) {
        fmt.Println("Logout event:", msg.Data)
    })

    // Note: Wildcards are not automatically matched
    // You need to publish to exact topics
    ps.Publish("user.login", "alice")
    ps.Publish("user.logout", "bob")
}
```

## API Reference

### PubSub

#### Creating

- `New(bufferSize int) *PubSub` - Create a new PubSub instance

#### Publishing

- `Publish(topic string, data interface{}) int` - Publish message asynchronously, returns subscriber count
- `PublishSync(ctx context.Context, topic string, data interface{}) error` - Publish message synchronously

#### Subscribing

- `Subscribe(topic string, handler Handler) *Subscriber` - Subscribe to a topic
- `Unsubscribe(sub *Subscriber)` - Remove a subscription
- `UnsubscribeAll(topic string)` - Remove all subscriptions from a topic

#### Introspection

- `Topics() []string` - Get all active topics
- `SubscriberCount(topic string) int` - Get number of subscribers for a topic

#### Cleanup

- `Close()` - Close all subscriptions and clean up

### Subscriber

- `Topic() string` - Get the topic this subscriber is subscribed to

### Message

```go
type Message struct {
    Topic string
    Data  interface{}
}
```

### Handler

```go
type Handler func(msg Message)
```

## Behavior

### Asynchronous Publishing

When using `Publish()`, messages are sent to subscriber channels asynchronously:
- Non-blocking send to each subscriber's channel
- If a subscriber's buffer is full, the message is skipped for that subscriber
- Returns the count of subscribers that received the message

### Synchronous Publishing

When using `PublishSync()`:
- Handlers are called synchronously in sequence
- Blocks until all handlers complete or context is cancelled
- Use for critical events where you need guaranteed processing

### Buffer Size

The buffer size determines how many messages can be queued per subscriber:
- Larger buffers prevent message loss when handlers are slow
- Smaller buffers use less memory
- Zero buffer means messages must be consumed immediately

### Thread Safety

All operations are thread-safe:
- Safe to publish and subscribe from multiple goroutines
- Internal locking protects subscriber management
- Each subscriber runs in its own goroutine

## License

MIT
