package chequeue_test

import (
	"fmt"

	"github.com/comfortablynumb/che/pkg/chequeue"
)

// Example demonstrates basic Queue operations
func Example() {
	queue := chequeue.New[int]()

	queue.Enqueue(1)
	queue.Enqueue(2)
	queue.Enqueue(3)

	// FIFO order
	val, _ := queue.Dequeue()
	fmt.Println("First out:", val)

	val, _ = queue.Dequeue()
	fmt.Println("Second out:", val)

	fmt.Println("Size:", queue.Size())

	// Output:
	// First out: 1
	// Second out: 2
	// Size: 1
}

// ExampleNew demonstrates creating a new queue
func ExampleNew() {
	queue := chequeue.New[string]()
	queue.Enqueue("hello")
	queue.Enqueue("world")

	fmt.Println("Size:", queue.Size())

	// Output:
	// Size: 2
}

// ExampleNewFromSlice demonstrates creating a queue from a slice
func ExampleNewFromSlice() {
	slice := []int{1, 2, 3, 4, 5}
	queue := chequeue.NewFromSlice(slice)

	// First element from slice is at front
	val, _ := queue.Dequeue()
	fmt.Println("First:", val)

	// Output:
	// First: 1
}

// ExampleQueue_Enqueue demonstrates adding elements
func ExampleQueue_Enqueue() {
	queue := chequeue.New[int]()

	queue.Enqueue(10)
	queue.Enqueue(20)
	queue.Enqueue(30)

	fmt.Println("Size:", queue.Size())

	// Output:
	// Size: 3
}

// ExampleQueue_Dequeue demonstrates removing elements
func ExampleQueue_Dequeue() {
	queue := chequeue.NewFromSlice([]int{1, 2, 3})

	val, ok := queue.Dequeue()
	fmt.Println("Value:", val, "Success:", ok)

	val, _ = queue.Dequeue()
	fmt.Println("Next:", val)

	// Output:
	// Value: 1 Success: true
	// Next: 2
}

// ExampleQueue_Peek demonstrates peeking at the front
func ExampleQueue_Peek() {
	queue := chequeue.NewFromSlice([]int{100, 200, 300})

	// Peek doesn't remove the element
	val, ok := queue.Peek()
	fmt.Println("Front:", val, "Success:", ok)
	fmt.Println("Size still:", queue.Size())

	// Output:
	// Front: 100 Success: true
	// Size still: 3
}

// ExampleQueue_ForEach demonstrates iterating over elements
func ExampleQueue_ForEach() {
	queue := chequeue.NewFromSlice([]string{"a", "b", "c"})

	queue.ForEach(func(item string) bool {
		fmt.Println(item)
		return true
	})

	// Output:
	// a
	// b
	// c
}

// ExampleQueue_Clone demonstrates cloning a queue
func ExampleQueue_Clone() {
	original := chequeue.NewFromSlice([]int{1, 2, 3})
	clone := original.Clone()

	clone.Enqueue(4)

	fmt.Println("Original size:", original.Size())
	fmt.Println("Clone size:", clone.Size())

	// Output:
	// Original size: 3
	// Clone size: 4
}

// ExampleQueue_taskQueue demonstrates using Queue for task processing
func ExampleQueue_taskQueue() {
	type Task struct {
		ID   int
		Name string
	}

	taskQueue := chequeue.New[Task]()

	// Add tasks
	taskQueue.Enqueue(Task{ID: 1, Name: "Process data"})
	taskQueue.Enqueue(Task{ID: 2, Name: "Send email"})
	taskQueue.Enqueue(Task{ID: 3, Name: "Update database"})

	// Process tasks in FIFO order
	for !taskQueue.IsEmpty() {
		task, _ := taskQueue.Dequeue()
		fmt.Printf("Processing task %d: %s\n", task.ID, task.Name)
	}

	// Output:
	// Processing task 1: Process data
	// Processing task 2: Send email
	// Processing task 3: Update database
}

// ExampleQueue_bfs demonstrates using Queue for breadth-first search
func ExampleQueue_bfs() {
	type Node struct {
		Value    int
		Children []int
	}

	// Simple tree: 1 -> [2, 3], 2 -> [4, 5]
	nodes := map[int]Node{
		1: {Value: 1, Children: []int{2, 3}},
		2: {Value: 2, Children: []int{4, 5}},
		3: {Value: 3, Children: []int{}},
		4: {Value: 4, Children: []int{}},
		5: {Value: 5, Children: []int{}},
	}

	queue := chequeue.New[int]()
	queue.Enqueue(1) // Start with root

	for !queue.IsEmpty() {
		nodeID, _ := queue.Dequeue()
		node := nodes[nodeID]
		fmt.Println(node.Value)

		for _, childID := range node.Children {
			queue.Enqueue(childID)
		}
	}

	// Output:
	// 1
	// 2
	// 3
	// 4
	// 5
}

// ExampleQueue_ToSlice demonstrates converting queue to slice
func ExampleQueue_ToSlice() {
	queue := chequeue.NewFromSlice([]int{10, 20, 30})

	slice := queue.ToSlice()
	fmt.Println("Slice:", slice)

	// Output:
	// Slice: [10 20 30]
}
