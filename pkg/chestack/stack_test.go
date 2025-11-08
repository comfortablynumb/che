package chestack

import (
	"testing"

	"github.com/comfortablynumb/che/pkg/chetest"
)

// TestNew tests creating a new empty Stack
func TestNew(t *testing.T) {
	stack := New[int]()

	chetest.RequireEqual(t, stack.IsEmpty(), true)
	chetest.RequireEqual(t, stack.Size(), 0)
}

// TestNewWithCapacity tests creating a Stack with initial capacity
func TestNewWithCapacity(t *testing.T) {
	stack := NewWithCapacity[int](100)

	chetest.RequireEqual(t, stack.IsEmpty(), true)
	chetest.RequireEqual(t, stack.Size(), 0)
}

// TestNewWithCapacity_Negative tests creating a Stack with negative capacity
func TestNewWithCapacity_Negative(t *testing.T) {
	stack := NewWithCapacity[int](-1)

	chetest.RequireEqual(t, stack.IsEmpty(), true)
	chetest.RequireEqual(t, stack.Size(), 0)
}

// TestNewFromSlice tests creating a Stack from a slice
func TestNewFromSlice(t *testing.T) {
	slice := []int{1, 2, 3, 4, 5}
	stack := NewFromSlice(slice)

	chetest.RequireEqual(t, stack.Size(), 5)

	// Verify LIFO order (last element in slice is on top)
	val, ok := stack.Pop()
	chetest.RequireEqual(t, ok, true)
	chetest.RequireEqual(t, val, 5)
}

// TestPush tests adding elements to stack
func TestPush(t *testing.T) {
	stack := New[int]()

	stack.Push(1)
	chetest.RequireEqual(t, stack.Size(), 1)

	stack.Push(2)
	chetest.RequireEqual(t, stack.Size(), 2)
}

// TestPushMultiple tests adding multiple elements
func TestPushMultiple(t *testing.T) {
	stack := New[int]()

	stack.PushMultiple(1, 2, 3, 4)
	chetest.RequireEqual(t, stack.Size(), 4)

	// Verify LIFO order (4 should be on top)
	val, _ := stack.Pop()
	chetest.RequireEqual(t, val, 4)
}

// TestPop tests removing elements from stack
func TestPop(t *testing.T) {
	stack := New[int]()
	stack.Push(1)
	stack.Push(2)
	stack.Push(3)

	val, ok := stack.Pop()
	chetest.RequireEqual(t, ok, true)
	chetest.RequireEqual(t, val, 3)
	chetest.RequireEqual(t, stack.Size(), 2)

	val, ok = stack.Pop()
	chetest.RequireEqual(t, ok, true)
	chetest.RequireEqual(t, val, 2)
}

// TestPop_Empty tests popping from empty stack
func TestPop_Empty(t *testing.T) {
	stack := New[int]()

	val, ok := stack.Pop()
	chetest.RequireEqual(t, ok, false)
	chetest.RequireEqual(t, val, 0)
}

// TestPeek tests peeking at top element
func TestPeek(t *testing.T) {
	stack := New[int]()
	stack.Push(1)
	stack.Push(2)

	val, ok := stack.Peek()
	chetest.RequireEqual(t, ok, true)
	chetest.RequireEqual(t, val, 2)
	chetest.RequireEqual(t, stack.Size(), 2) // Size unchanged
}

// TestPeek_Empty tests peeking at empty stack
func TestPeek_Empty(t *testing.T) {
	stack := New[int]()

	val, ok := stack.Peek()
	chetest.RequireEqual(t, ok, false)
	chetest.RequireEqual(t, val, 0)
}

// TestSize tests getting stack size
func TestSize(t *testing.T) {
	stack := New[int]()
	chetest.RequireEqual(t, stack.Size(), 0)

	stack.Push(1)
	chetest.RequireEqual(t, stack.Size(), 1)

	stack.Push(2)
	chetest.RequireEqual(t, stack.Size(), 2)

	stack.Pop()
	chetest.RequireEqual(t, stack.Size(), 1)
}

// TestIsEmpty tests checking if stack is empty
func TestIsEmpty(t *testing.T) {
	stack := New[int]()
	chetest.RequireEqual(t, stack.IsEmpty(), true)

	stack.Push(1)
	chetest.RequireEqual(t, stack.IsEmpty(), false)

	stack.Pop()
	chetest.RequireEqual(t, stack.IsEmpty(), true)
}

// TestClear tests clearing all elements
func TestClear(t *testing.T) {
	stack := New[int]()
	stack.PushMultiple(1, 2, 3, 4, 5)

	stack.Clear()
	chetest.RequireEqual(t, stack.IsEmpty(), true)
	chetest.RequireEqual(t, stack.Size(), 0)
}

// TestToSlice tests converting stack to slice
func TestToSlice(t *testing.T) {
	stack := New[int]()
	stack.PushMultiple(1, 2, 3)

	slice := stack.ToSlice()
	chetest.RequireEqual(t, len(slice), 3)
	chetest.RequireEqual(t, slice[0], 1)
	chetest.RequireEqual(t, slice[1], 2)
	chetest.RequireEqual(t, slice[2], 3)
}

// TestToSlice_Empty tests converting empty stack
func TestToSlice_Empty(t *testing.T) {
	stack := New[int]()

	slice := stack.ToSlice()
	chetest.RequireEqual(t, len(slice), 0)
}

// TestClone tests cloning a stack
func TestClone(t *testing.T) {
	original := New[int]()
	original.PushMultiple(1, 2, 3)

	clone := original.Clone()
	chetest.RequireEqual(t, clone.Size(), 3)

	// Modifying clone shouldn't affect original
	clone.Push(4)
	chetest.RequireEqual(t, original.Size(), 3)
	chetest.RequireEqual(t, clone.Size(), 4)
}

// TestForEach tests iterating over stack elements
func TestForEach(t *testing.T) {
	stack := New[int]()
	stack.PushMultiple(1, 2, 3)

	var collected []int
	stack.ForEach(func(item int) bool {
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
	stack := New[int]()
	stack.PushMultiple(1, 2, 3, 4, 5)

	count := 0
	stack.ForEach(func(item int) bool {
		count++
		return count < 3
	})

	chetest.RequireEqual(t, count, 3)
}

// TestForEachReverse tests iterating in reverse order
func TestForEachReverse(t *testing.T) {
	stack := New[int]()
	stack.PushMultiple(1, 2, 3)

	var collected []int
	stack.ForEachReverse(func(item int) bool {
		collected = append(collected, item)
		return true
	})

	chetest.RequireEqual(t, len(collected), 3)
	chetest.RequireEqual(t, collected[0], 3)
	chetest.RequireEqual(t, collected[1], 2)
	chetest.RequireEqual(t, collected[2], 1)
}

// TestForEachReverse_EarlyExit tests early exit from reverse iteration
func TestForEachReverse_EarlyExit(t *testing.T) {
	stack := New[int]()
	stack.PushMultiple(1, 2, 3, 4, 5)

	count := 0
	stack.ForEachReverse(func(item int) bool {
		count++
		return count < 3
	})

	chetest.RequireEqual(t, count, 3)
}

// TestContains tests checking if element exists
func TestContains(t *testing.T) {
	stack := New[int]()
	stack.PushMultiple(1, 2, 3)

	equals := func(a, b int) bool { return a == b }

	chetest.RequireEqual(t, stack.Contains(2, equals), true)
	chetest.RequireEqual(t, stack.Contains(5, equals), false)
}

// TestString tests string representation
func TestString(t *testing.T) {
	stack := New[int]()
	str := stack.String()
	chetest.RequireEqual(t, str, "Stack[]")

	stack.Push(1)
	stack.Push(2)
	str = stack.String()
	chetest.RequireEqual(t, str, "Stack[1, 2]")
}

// TestLIFOOrder tests Last-In-First-Out order
func TestLIFOOrder(t *testing.T) {
	stack := New[string]()

	stack.Push("first")
	stack.Push("second")
	stack.Push("third")

	val, _ := stack.Pop()
	chetest.RequireEqual(t, val, "third")

	val, _ = stack.Pop()
	chetest.RequireEqual(t, val, "second")

	val, _ = stack.Pop()
	chetest.RequireEqual(t, val, "first")

	chetest.RequireEqual(t, stack.IsEmpty(), true)
}

// TestMixedOperations tests mixing push and pop
func TestMixedOperations(t *testing.T) {
	stack := New[int]()

	stack.Push(1)
	stack.Push(2)

	val, _ := stack.Pop()
	chetest.RequireEqual(t, val, 2)

	stack.Push(3)
	stack.Push(4)

	val, _ = stack.Pop()
	chetest.RequireEqual(t, val, 4)

	val, _ = stack.Pop()
	chetest.RequireEqual(t, val, 3)

	chetest.RequireEqual(t, stack.Size(), 1)
}

// TestGrowth tests automatic growth
func TestGrowth(t *testing.T) {
	stack := New[int]()

	// Add many elements to trigger growth
	for i := 0; i < 100; i++ {
		stack.Push(i)
	}

	chetest.RequireEqual(t, stack.Size(), 100)

	// Verify LIFO order
	val, _ := stack.Pop()
	chetest.RequireEqual(t, val, 99)
}
