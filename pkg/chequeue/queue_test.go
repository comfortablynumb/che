package chequeue

import (
	"testing"

	"github.com/comfortablynumb/che/pkg/chetest"
)

// TestNew tests creating a new empty Queue
func TestNew(t *testing.T) {
	queue := New[int]()

	chetest.RequireEqual(t, queue.IsEmpty(), true)
	chetest.RequireEqual(t, queue.Size(), 0)
}

// TestNewWithCapacity tests creating a Queue with initial capacity
func TestNewWithCapacity(t *testing.T) {
	queue := NewWithCapacity[int](100)

	chetest.RequireEqual(t, queue.IsEmpty(), true)
	chetest.RequireEqual(t, queue.Size(), 0)
}

// TestNewWithCapacity_Negative tests creating a Queue with negative capacity
func TestNewWithCapacity_Negative(t *testing.T) {
	queue := NewWithCapacity[int](-1)

	chetest.RequireEqual(t, queue.IsEmpty(), true)
	chetest.RequireEqual(t, queue.Size(), 0)
}

// TestNewFromSlice tests creating a Queue from a slice
func TestNewFromSlice(t *testing.T) {
	slice := []int{1, 2, 3, 4, 5}
	queue := NewFromSlice(slice)

	chetest.RequireEqual(t, queue.Size(), 5)

	// Verify FIFO order
	val, ok := queue.Dequeue()
	chetest.RequireEqual(t, ok, true)
	chetest.RequireEqual(t, val, 1)
}

// TestEnqueue tests adding elements to queue
func TestEnqueue(t *testing.T) {
	queue := New[int]()

	queue.Enqueue(1)
	chetest.RequireEqual(t, queue.Size(), 1)

	queue.Enqueue(2)
	chetest.RequireEqual(t, queue.Size(), 2)
}

// TestEnqueueMultiple tests adding multiple elements
func TestEnqueueMultiple(t *testing.T) {
	queue := New[int]()

	queue.EnqueueMultiple(1, 2, 3, 4)
	chetest.RequireEqual(t, queue.Size(), 4)

	// Verify order
	val, _ := queue.Dequeue()
	chetest.RequireEqual(t, val, 1)
}

// TestDequeue tests removing elements from queue
func TestDequeue(t *testing.T) {
	queue := New[int]()
	queue.Enqueue(1)
	queue.Enqueue(2)
	queue.Enqueue(3)

	val, ok := queue.Dequeue()
	chetest.RequireEqual(t, ok, true)
	chetest.RequireEqual(t, val, 1)
	chetest.RequireEqual(t, queue.Size(), 2)

	val, ok = queue.Dequeue()
	chetest.RequireEqual(t, ok, true)
	chetest.RequireEqual(t, val, 2)
}

// TestDequeue_Empty tests dequeuing from empty queue
func TestDequeue_Empty(t *testing.T) {
	queue := New[int]()

	val, ok := queue.Dequeue()
	chetest.RequireEqual(t, ok, false)
	chetest.RequireEqual(t, val, 0)
}

// TestPeek tests peeking at front element
func TestPeek(t *testing.T) {
	queue := New[int]()
	queue.Enqueue(1)
	queue.Enqueue(2)

	val, ok := queue.Peek()
	chetest.RequireEqual(t, ok, true)
	chetest.RequireEqual(t, val, 1)
	chetest.RequireEqual(t, queue.Size(), 2) // Size unchanged
}

// TestPeek_Empty tests peeking at empty queue
func TestPeek_Empty(t *testing.T) {
	queue := New[int]()

	val, ok := queue.Peek()
	chetest.RequireEqual(t, ok, false)
	chetest.RequireEqual(t, val, 0)
}

// TestSize tests getting queue size
func TestSize(t *testing.T) {
	queue := New[int]()
	chetest.RequireEqual(t, queue.Size(), 0)

	queue.Enqueue(1)
	chetest.RequireEqual(t, queue.Size(), 1)

	queue.Enqueue(2)
	chetest.RequireEqual(t, queue.Size(), 2)

	queue.Dequeue()
	chetest.RequireEqual(t, queue.Size(), 1)
}

// TestIsEmpty tests checking if queue is empty
func TestIsEmpty(t *testing.T) {
	queue := New[int]()
	chetest.RequireEqual(t, queue.IsEmpty(), true)

	queue.Enqueue(1)
	chetest.RequireEqual(t, queue.IsEmpty(), false)

	queue.Dequeue()
	chetest.RequireEqual(t, queue.IsEmpty(), true)
}

// TestClear tests clearing all elements
func TestClear(t *testing.T) {
	queue := New[int]()
	queue.EnqueueMultiple(1, 2, 3, 4, 5)

	queue.Clear()
	chetest.RequireEqual(t, queue.IsEmpty(), true)
	chetest.RequireEqual(t, queue.Size(), 0)
}

// TestToSlice tests converting queue to slice
func TestToSlice(t *testing.T) {
	queue := New[int]()
	queue.EnqueueMultiple(1, 2, 3)

	slice := queue.ToSlice()
	chetest.RequireEqual(t, len(slice), 3)
	chetest.RequireEqual(t, slice[0], 1)
	chetest.RequireEqual(t, slice[1], 2)
	chetest.RequireEqual(t, slice[2], 3)
}

// TestToSlice_Empty tests converting empty queue
func TestToSlice_Empty(t *testing.T) {
	queue := New[int]()

	slice := queue.ToSlice()
	chetest.RequireEqual(t, len(slice), 0)
}

// TestClone tests cloning a queue
func TestClone(t *testing.T) {
	original := New[int]()
	original.EnqueueMultiple(1, 2, 3)

	clone := original.Clone()
	chetest.RequireEqual(t, clone.Size(), 3)

	// Modifying clone shouldn't affect original
	clone.Enqueue(4)
	chetest.RequireEqual(t, original.Size(), 3)
	chetest.RequireEqual(t, clone.Size(), 4)
}

// TestForEach tests iterating over queue elements
func TestForEach(t *testing.T) {
	queue := New[int]()
	queue.EnqueueMultiple(1, 2, 3)

	var collected []int
	queue.ForEach(func(item int) bool {
		collected = append(collected, item)
		return true
	})

	chetest.RequireEqual(t, len(collected), 3)
	chetest.RequireEqual(t, collected[0], 1)
	chetest.RequireEqual(t, collected[1], 2)
	chetest.RequireEqual(t, collected[2], 3)
}

// TestForEach_EarlyExit tests early exit from iteration
func TestForEach_EarlyExit(t *testing.T) {
	queue := New[int]()
	queue.EnqueueMultiple(1, 2, 3, 4, 5)

	count := 0
	queue.ForEach(func(item int) bool {
		count++
		return count < 3
	})

	chetest.RequireEqual(t, count, 3)
}

// TestContains tests checking if element exists
func TestContains(t *testing.T) {
	queue := New[int]()
	queue.EnqueueMultiple(1, 2, 3)

	equals := func(a, b int) bool { return a == b }

	chetest.RequireEqual(t, queue.Contains(2, equals), true)
	chetest.RequireEqual(t, queue.Contains(5, equals), false)
}

// TestString tests string representation
func TestString(t *testing.T) {
	queue := New[int]()
	str := queue.String()
	chetest.RequireEqual(t, str, "Queue[]")

	queue.Enqueue(1)
	queue.Enqueue(2)
	str = queue.String()
	chetest.RequireEqual(t, str, "Queue[1, 2]")
}

// TestResize tests automatic resizing
func TestResize(t *testing.T) {
	queue := New[int]()

	// Add more than initial capacity to trigger resize
	for i := 0; i < 20; i++ {
		queue.Enqueue(i)
	}

	chetest.RequireEqual(t, queue.Size(), 20)

	// Verify order is maintained
	for i := 0; i < 20; i++ {
		val, ok := queue.Dequeue()
		chetest.RequireEqual(t, ok, true)
		chetest.RequireEqual(t, val, i)
	}
}

// TestShrink tests automatic shrinking
func TestShrink(t *testing.T) {
	queue := New[int]()

	// Fill queue to trigger growth
	for i := 0; i < 100; i++ {
		queue.Enqueue(i)
	}

	// Remove most elements to trigger shrink
	for i := 0; i < 95; i++ {
		queue.Dequeue()
	}

	// Should have shrunk but still work correctly
	chetest.RequireEqual(t, queue.Size(), 5)

	val, ok := queue.Dequeue()
	chetest.RequireEqual(t, ok, true)
	chetest.RequireEqual(t, val, 95)
}

// TestCircularBehavior tests the circular buffer behavior
func TestCircularBehavior(t *testing.T) {
	queue := NewWithCapacity[int](4)

	// Fill queue
	queue.EnqueueMultiple(1, 2, 3)

	// Dequeue some
	queue.Dequeue()
	queue.Dequeue()

	// Enqueue more to wrap around
	queue.EnqueueMultiple(4, 5, 6)

	// Verify correct order
	val, _ := queue.Dequeue()
	chetest.RequireEqual(t, val, 3)
	val, _ = queue.Dequeue()
	chetest.RequireEqual(t, val, 4)
	val, _ = queue.Dequeue()
	chetest.RequireEqual(t, val, 5)
	val, _ = queue.Dequeue()
	chetest.RequireEqual(t, val, 6)
}

// TestFIFOOrder tests First-In-First-Out order
func TestFIFOOrder(t *testing.T) {
	queue := New[string]()

	queue.Enqueue("first")
	queue.Enqueue("second")
	queue.Enqueue("third")

	val, _ := queue.Dequeue()
	chetest.RequireEqual(t, val, "first")

	val, _ = queue.Dequeue()
	chetest.RequireEqual(t, val, "second")

	val, _ = queue.Dequeue()
	chetest.RequireEqual(t, val, "third")

	chetest.RequireEqual(t, queue.IsEmpty(), true)
}

// TestMixedOperations tests mixing enqueue and dequeue
func TestMixedOperations(t *testing.T) {
	queue := New[int]()

	queue.Enqueue(1)
	queue.Enqueue(2)

	val, _ := queue.Dequeue()
	chetest.RequireEqual(t, val, 1)

	queue.Enqueue(3)
	queue.Enqueue(4)

	val, _ = queue.Dequeue()
	chetest.RequireEqual(t, val, 2)

	val, _ = queue.Dequeue()
	chetest.RequireEqual(t, val, 3)

	chetest.RequireEqual(t, queue.Size(), 1)
}
