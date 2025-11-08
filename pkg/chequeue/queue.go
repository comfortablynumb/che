// Package chequeue provides a generic Queue implementation for Go.
//
// Queue is a FIFO (First-In-First-Out) data structure that provides
// O(1) amortized performance for Enqueue and Dequeue operations.
//
// # Basic Usage
//
//	queue := chequeue.New[int]()
//	queue.Enqueue(1)
//	queue.Enqueue(2)
//	value, ok := queue.Dequeue() // 1, true
//	fmt.Println(queue.Size())    // 1
//
// # Thread Safety
//
// Queue is not thread-safe. For concurrent use, external synchronization
// is required (e.g., using sync.Mutex).
//
// # Performance
//
// Enqueue and Dequeue operations have O(1) amortized time complexity.
// The implementation uses a dynamic slice with automatic resizing.
package chequeue

import "fmt"

// Queue is a generic FIFO (First-In-First-Out) queue implementation.
// It provides O(1) amortized time complexity for Enqueue and Dequeue operations.
// Queue is not thread-safe. For concurrent use, external synchronization is required.
type Queue[T any] struct {
	items []T
	head  int // Index of the front element
	tail  int // Index where next element will be inserted
	count int // Number of elements in the queue
}

// New creates and returns a new empty Queue.
func New[T any]() *Queue[T] {
	return &Queue[T]{
		items: make([]T, 8), // Start with small capacity
		head:  0,
		tail:  0,
		count: 0,
	}
}

// NewWithCapacity creates and returns a new empty Queue with the specified initial capacity.
// This can improve performance when the expected size is known in advance.
func NewWithCapacity[T any](capacity int) *Queue[T] {
	if capacity < 1 {
		capacity = 8
	}
	return &Queue[T]{
		items: make([]T, capacity),
		head:  0,
		tail:  0,
		count: 0,
	}
}

// NewFromSlice creates and returns a new Queue containing all elements from the given slice.
// Elements are enqueued in the order they appear in the slice.
func NewFromSlice[T any](slice []T) *Queue[T] {
	queue := NewWithCapacity[T](len(slice))
	for _, item := range slice {
		queue.Enqueue(item)
	}
	return queue
}

// Enqueue adds an element to the back of the queue.
func (q *Queue[T]) Enqueue(item T) {
	if q.count == len(q.items) {
		q.resize()
	}
	q.items[q.tail] = item
	q.tail = (q.tail + 1) % len(q.items)
	q.count++
}

// EnqueueMultiple adds multiple elements to the back of the queue.
func (q *Queue[T]) EnqueueMultiple(items ...T) {
	for _, item := range items {
		q.Enqueue(item)
	}
}

// Dequeue removes and returns the element at the front of the queue.
// Returns the element and true if successful, or zero value and false if the queue is empty.
func (q *Queue[T]) Dequeue() (T, bool) {
	if q.IsEmpty() {
		var zero T
		return zero, false
	}
	item := q.items[q.head]
	var zero T
	q.items[q.head] = zero // Clear reference to help GC
	q.head = (q.head + 1) % len(q.items)
	q.count--

	// Shrink if capacity is much larger than needed
	if len(q.items) > 8 && q.count < len(q.items)/4 {
		q.shrink()
	}

	return item, true
}

// Peek returns the element at the front of the queue without removing it.
// Returns the element and true if successful, or zero value and false if the queue is empty.
func (q *Queue[T]) Peek() (T, bool) {
	if q.IsEmpty() {
		var zero T
		return zero, false
	}
	return q.items[q.head], true
}

// Size returns the number of elements in the queue.
func (q *Queue[T]) Size() int {
	return q.count
}

// IsEmpty returns true if the queue contains no elements.
func (q *Queue[T]) IsEmpty() bool {
	return q.count == 0
}

// Clear removes all elements from the queue.
func (q *Queue[T]) Clear() {
	q.items = make([]T, 8)
	q.head = 0
	q.tail = 0
	q.count = 0
}

// ToSlice returns a slice containing all elements in FIFO order.
func (q *Queue[T]) ToSlice() []T {
	result := make([]T, q.count)
	for i := 0; i < q.count; i++ {
		idx := (q.head + i) % len(q.items)
		result[i] = q.items[idx]
	}
	return result
}

// Clone returns a shallow copy of the queue.
func (q *Queue[T]) Clone() *Queue[T] {
	clone := NewWithCapacity[T](q.count)
	for i := 0; i < q.count; i++ {
		idx := (q.head + i) % len(q.items)
		clone.Enqueue(q.items[idx])
	}
	return clone
}

// ForEach executes a function for each element in the queue in FIFO order.
// If the function returns false, iteration stops.
func (q *Queue[T]) ForEach(fn func(item T) bool) {
	for i := 0; i < q.count; i++ {
		idx := (q.head + i) % len(q.items)
		if !fn(q.items[idx]) {
			return
		}
	}
}

// Contains checks if an element exists in the queue.
// This operation is O(n) as it requires scanning all elements.
func (q *Queue[T]) Contains(item T, equals func(a, b T) bool) bool {
	for i := 0; i < q.count; i++ {
		idx := (q.head + i) % len(q.items)
		if equals(q.items[idx], item) {
			return true
		}
	}
	return false
}

// String returns a string representation of the queue.
func (q *Queue[T]) String() string {
	if q.IsEmpty() {
		return "Queue[]"
	}

	result := "Queue["
	for i := 0; i < q.count; i++ {
		if i > 0 {
			result += ", "
		}
		idx := (q.head + i) % len(q.items)
		result += fmt.Sprint(q.items[idx])
	}
	result += "]"
	return result
}

// resize doubles the capacity of the queue.
func (q *Queue[T]) resize() {
	newCapacity := len(q.items) * 2
	newItems := make([]T, newCapacity)

	// Copy elements in order from head to tail
	for i := 0; i < q.count; i++ {
		idx := (q.head + i) % len(q.items)
		newItems[i] = q.items[idx]
	}

	q.items = newItems
	q.head = 0
	q.tail = q.count
}

// shrink reduces the capacity of the queue to save memory.
func (q *Queue[T]) shrink() {
	newCapacity := len(q.items) / 2
	if newCapacity < 8 {
		newCapacity = 8
	}

	newItems := make([]T, newCapacity)

	// Copy elements in order from head to tail
	for i := 0; i < q.count; i++ {
		idx := (q.head + i) % len(q.items)
		newItems[i] = q.items[idx]
	}

	q.items = newItems
	q.head = 0
	q.tail = q.count
}
