// Package chestack provides a generic Stack implementation for Go.
//
// Stack is a LIFO (Last-In-First-Out) data structure that provides
// O(1) amortized performance for Push and Pop operations.
//
// # Basic Usage
//
//	stack := chestack.New[int]()
//	stack.Push(1)
//	stack.Push(2)
//	value, ok := stack.Pop() // 2, true
//	fmt.Println(stack.Size()) // 1
//
// # Thread Safety
//
// Stack is not thread-safe. For concurrent use, external synchronization
// is required (e.g., using sync.Mutex).
//
// # Performance
//
// Push and Pop operations have O(1) amortized time complexity.
// The implementation uses a dynamic slice with automatic resizing.
package chestack

import "fmt"

// Stack is a generic LIFO (Last-In-First-Out) stack implementation.
// It provides O(1) amortized time complexity for Push and Pop operations.
// Stack is not thread-safe. For concurrent use, external synchronization is required.
type Stack[T any] struct {
	items []T
}

// New creates and returns a new empty Stack.
func New[T any]() *Stack[T] {
	return &Stack[T]{
		items: make([]T, 0, 8), // Start with small capacity
	}
}

// NewWithCapacity creates and returns a new empty Stack with the specified initial capacity.
// This can improve performance when the expected size is known in advance.
func NewWithCapacity[T any](capacity int) *Stack[T] {
	if capacity < 0 {
		capacity = 8
	}
	return &Stack[T]{
		items: make([]T, 0, capacity),
	}
}

// NewFromSlice creates and returns a new Stack containing all elements from the given slice.
// Elements are pushed in the order they appear in the slice, so the last element
// in the slice will be at the top of the stack.
func NewFromSlice[T any](slice []T) *Stack[T] {
	stack := &Stack[T]{
		items: make([]T, len(slice)),
	}
	copy(stack.items, slice)
	return stack
}

// Push adds an element to the top of the stack.
func (s *Stack[T]) Push(item T) {
	s.items = append(s.items, item)
}

// PushMultiple adds multiple elements to the top of the stack.
// Elements are pushed in the order they appear in the arguments.
func (s *Stack[T]) PushMultiple(items ...T) {
	s.items = append(s.items, items...)
}

// Pop removes and returns the element at the top of the stack.
// Returns the element and true if successful, or zero value and false if the stack is empty.
func (s *Stack[T]) Pop() (T, bool) {
	if s.IsEmpty() {
		var zero T
		return zero, false
	}
	lastIdx := len(s.items) - 1
	item := s.items[lastIdx]

	// Clear reference to help GC
	var zero T
	s.items[lastIdx] = zero

	s.items = s.items[:lastIdx]
	return item, true
}

// Peek returns the element at the top of the stack without removing it.
// Returns the element and true if successful, or zero value and false if the stack is empty.
func (s *Stack[T]) Peek() (T, bool) {
	if s.IsEmpty() {
		var zero T
		return zero, false
	}
	return s.items[len(s.items)-1], true
}

// Size returns the number of elements in the stack.
func (s *Stack[T]) Size() int {
	return len(s.items)
}

// IsEmpty returns true if the stack contains no elements.
func (s *Stack[T]) IsEmpty() bool {
	return len(s.items) == 0
}

// Clear removes all elements from the stack.
func (s *Stack[T]) Clear() {
	s.items = make([]T, 0, 8)
}

// ToSlice returns a slice containing all elements from bottom to top.
// The first element in the slice is the bottom of the stack.
func (s *Stack[T]) ToSlice() []T {
	result := make([]T, len(s.items))
	copy(result, s.items)
	return result
}

// Clone returns a shallow copy of the stack.
func (s *Stack[T]) Clone() *Stack[T] {
	clone := &Stack[T]{
		items: make([]T, len(s.items), cap(s.items)),
	}
	copy(clone.items, s.items)
	return clone
}

// ForEach executes a function for each element in the stack from bottom to top.
// If the function returns false, iteration stops.
func (s *Stack[T]) ForEach(fn func(item T) bool) {
	for _, item := range s.items {
		if !fn(item) {
			return
		}
	}
}

// ForEachReverse executes a function for each element in the stack from top to bottom.
// If the function returns false, iteration stops.
func (s *Stack[T]) ForEachReverse(fn func(item T) bool) {
	for i := len(s.items) - 1; i >= 0; i-- {
		if !fn(s.items[i]) {
			return
		}
	}
}

// Contains checks if an element exists in the stack.
// This operation is O(n) as it requires scanning all elements.
func (s *Stack[T]) Contains(item T, equals func(a, b T) bool) bool {
	for _, element := range s.items {
		if equals(element, item) {
			return true
		}
	}
	return false
}

// String returns a string representation of the stack.
// Elements are shown from bottom to top, with the top element on the right.
func (s *Stack[T]) String() string {
	if s.IsEmpty() {
		return "Stack[]"
	}

	result := "Stack["
	for i, item := range s.items {
		if i > 0 {
			result += ", "
		}
		result += fmt.Sprint(item)
	}
	result += "]"
	return result
}
