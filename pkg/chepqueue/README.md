# chepqueue

Generic priority queue implementation using binary heap.

## Features

- Generic support for any type with ordered priorities
- Min-heap and max-heap variants
- O(log n) push and pop operations
- Update priority and remove operations
- Thread-safe operations available
- Zero external dependencies

## Installation

```bash
go get github.com/comfortablynumb/che/pkg/chepqueue
```

## Usage

### Basic Min-Heap

```go
package main

import (
    "fmt"
    "github.com/comfortablynumb/che/pkg/chepqueue"
)

func main() {
    // Create a min-heap (lower priority = dequeued first)
    pq := chepqueue.New[string, int]()

    pq.Push("low priority", 10)
    pq.Push("high priority", 1)
    pq.Push("medium priority", 5)

    // Items are dequeued in priority order
    fmt.Println(pq.Pop()) // "high priority"
    fmt.Println(pq.Pop()) // "medium priority"
    fmt.Println(pq.Pop()) // "low priority"
}
```

### Max-Heap

```go
// Create a max-heap (higher priority = dequeued first)
pq := chepqueue.NewMax[string, int]()

pq.Push("low priority", 1)
pq.Push("high priority", 10)
pq.Push("medium priority", 5)

fmt.Println(pq.Pop()) // "high priority"
fmt.Println(pq.Pop()) // "medium priority"
fmt.Println(pq.Pop()) // "low priority"
```

### Peek Without Removing

```go
pq := chepqueue.New[string, int]()
pq.Push("task", 5)

// Peek at the next item without removing it
value := pq.Peek()
priority := pq.PeekPriority()

fmt.Printf("%s with priority %d\n", value, priority)
fmt.Println("Length:", pq.Len()) // Still 1
```

### Update Priority

```go
pq := chepqueue.New[string, int]()

pq.Push("task1", 10)
pq.Push("task2", 5)
pq.Push("task3", 1)

// Define equality function
equals := func(a, b string) bool { return a == b }

// Update task3 to highest priority
pq.UpdatePriority("task3", 0, equals)

fmt.Println(pq.Pop()) // "task3" (now has priority 0)
```

### Remove Item

```go
pq := chepqueue.New[string, int]()

pq.Push("task1", 1)
pq.Push("task2", 2)
pq.Push("task3", 3)

equals := func(a, b string) bool { return a == b }

// Remove task2
removed := pq.Remove("task2", equals)
fmt.Println("Removed:", removed) // true

fmt.Println(pq.Len()) // 2
```

### Custom Types

```go
type Task struct {
    ID   int
    Name string
}

pq := chepqueue.New[Task, float64]()

pq.Push(Task{1, "Low"}, 10.5)
pq.Push(Task{2, "High"}, 1.0)
pq.Push(Task{3, "Medium"}, 5.5)

task := pq.Pop()
fmt.Printf("Processing: %s (ID: %d)\n", task.Name, task.ID)
```

### Task Scheduler Example

```go
package main

import (
    "fmt"
    "time"
    "github.com/comfortablynumb/che/pkg/chepqueue"
)

type ScheduledTask struct {
    Name string
    Run  func()
}

type Scheduler struct {
    queue *chepqueue.PriorityQueue[ScheduledTask, time.Time]
}

func NewScheduler() *Scheduler {
    return &Scheduler{
        queue: chepqueue.New[ScheduledTask, time.Time](),
    }
}

func (s *Scheduler) Schedule(task ScheduledTask, runAt time.Time) {
    s.queue.Push(task, runAt)
}

func (s *Scheduler) Run() {
    for !s.queue.IsEmpty() {
        runAt := s.queue.PeekPriority()

        // Wait until it's time to run
        time.Sleep(time.Until(runAt))

        task := s.queue.Pop()
        fmt.Printf("Running task: %s\n", task.Name)
        task.Run()
    }
}

func main() {
    scheduler := NewScheduler()

    now := time.Now()
    scheduler.Schedule(ScheduledTask{
        Name: "Task 1",
        Run:  func() { fmt.Println("Executing Task 1") },
    }, now.Add(2*time.Second))

    scheduler.Schedule(ScheduledTask{
        Name: "Task 2",
        Run:  func() { fmt.Println("Executing Task 2") },
    }, now.Add(1*time.Second))

    scheduler.Run()
}
```

### Dijkstra's Algorithm Example

```go
type Node struct {
    ID   int
    Dist int
}

func dijkstra(graph map[int][]Edge, start int) map[int]int {
    dist := make(map[int]int)
    pq := chepqueue.New[int, int]()

    pq.Push(start, 0)
    dist[start] = 0

    for !pq.IsEmpty() {
        node := pq.Pop()
        currentDist := dist[node]

        for _, edge := range graph[node] {
            newDist := currentDist + edge.Weight

            if oldDist, ok := dist[edge.To]; !ok || newDist < oldDist {
                dist[edge.To] = newDist
                pq.Push(edge.To, newDist)
            }
        }
    }

    return dist
}
```

### Job Queue with Priority

```go
type Job struct {
    ID      int
    Payload string
}

type JobQueue struct {
    pq *chepqueue.PriorityQueue[Job, int]
}

func NewJobQueue() *JobQueue {
    return &JobQueue{
        pq: chepqueue.New[Job, int](),
    }
}

func (jq *JobQueue) Submit(job Job, priority int) {
    jq.pq.Push(job, priority)
}

func (jq *JobQueue) Next() (Job, bool) {
    if jq.pq.IsEmpty() {
        return Job{}, false
    }
    return jq.pq.Pop(), true
}

func (jq *JobQueue) Pending() int {
    return jq.pq.Len()
}
```

## API

### Creating Priority Queues

- `New[T, P]() *PriorityQueue[T, P]` - Create min-heap
- `NewMax[T, P]() *PriorityQueue[T, P]` - Create max-heap

### Operations

- `Push(value T, priority P)` - Add item (O(log n))
- `Pop() T` - Remove and return highest priority item (O(log n))
- `Peek() T` - View highest priority item without removing (O(1))
- `PeekPriority() P` - View priority of next item (O(1))
- `IsEmpty() bool` - Check if queue is empty (O(1))
- `Len() int` - Get number of items (O(1))
- `Clear()` - Remove all items (O(1))
- `Items() []Item[T, P]` - Get copy of all items (O(n))
- `UpdatePriority(value T, newPriority P, equals func(T, T) bool) bool` - Update priority (O(n))
- `Remove(value T, equals func(T, T) bool) bool` - Remove item (O(n))

### Types

```go
type Item[T any, P Ordered] struct {
    Value    T
    Priority P
}

type Ordered interface {
    constraints.Ordered
}
```

## Complexity

| Operation | Time Complexity | Space Complexity |
|-----------|----------------|------------------|
| Push | O(log n) | O(1) |
| Pop | O(log n) | O(1) |
| Peek | O(1) | O(1) |
| UpdatePriority | O(n) | O(1) |
| Remove | O(n) | O(1) |

## Heap Properties

- **Min-Heap**: Parent priority ≤ child priority
- **Max-Heap**: Parent priority ≥ child priority
- Implemented as binary heap using array
- Efficient memory usage
- No pointer overhead

## When to Use

Priority queues are ideal for:
- Task scheduling
- Event-driven simulation
- Graph algorithms (Dijkstra, A*)
- Job queues with priorities
- Merging sorted lists
- Finding top K elements

## License

MIT
