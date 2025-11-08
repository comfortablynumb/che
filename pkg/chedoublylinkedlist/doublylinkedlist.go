package chedoublylinkedlist

// Node represents a single node in a doubly linked list
type Node[T any] struct {
	Value T
	Next  *Node[T]
	Prev  *Node[T]
}

// DoublyLinkedList is a doubly linked list implementation with O(1) prepend and append operations
type DoublyLinkedList[T any] struct {
	head *Node[T]
	tail *Node[T]
	size int
}

// New creates a new empty DoublyLinkedList
func New[T any]() *DoublyLinkedList[T] {
	return &DoublyLinkedList[T]{
		head: nil,
		tail: nil,
		size: 0,
	}
}

// Prepend adds an element to the beginning of the list - O(1)
func (dll *DoublyLinkedList[T]) Prepend(value T) {
	newNode := &Node[T]{Value: value, Next: dll.head, Prev: nil}

	if dll.head != nil {
		dll.head.Prev = newNode
	}

	dll.head = newNode

	if dll.tail == nil {
		dll.tail = newNode
	}

	dll.size++
}

// Append adds an element to the end of the list - O(1)
func (dll *DoublyLinkedList[T]) Append(value T) {
	newNode := &Node[T]{Value: value, Next: nil, Prev: dll.tail}

	if dll.tail != nil {
		dll.tail.Next = newNode
	}

	dll.tail = newNode

	if dll.head == nil {
		dll.head = newNode
	}

	dll.size++
}

// InsertAt inserts an element at the specified index - O(n)
// Returns false if index is out of bounds
func (dll *DoublyLinkedList[T]) InsertAt(index int, value T) bool {
	if index < 0 || index > dll.size {
		return false
	}

	if index == 0 {
		dll.Prepend(value)
		return true
	}

	if index == dll.size {
		dll.Append(value)
		return true
	}

	current := dll.head
	for i := 0; i < index; i++ {
		current = current.Next
	}

	newNode := &Node[T]{Value: value, Next: current, Prev: current.Prev}
	current.Prev.Next = newNode
	current.Prev = newNode
	dll.size++

	return true
}

// RemoveFirst removes and returns the first element - O(1)
// Returns the element and true if successful, zero value and false if list is empty
func (dll *DoublyLinkedList[T]) RemoveFirst() (T, bool) {
	if dll.head == nil {
		var zero T
		return zero, false
	}

	value := dll.head.Value
	dll.head = dll.head.Next

	if dll.head != nil {
		dll.head.Prev = nil
	} else {
		dll.tail = nil
	}

	dll.size--
	return value, true
}

// RemoveLast removes and returns the last element - O(1)
// Returns the element and true if successful, zero value and false if list is empty
func (dll *DoublyLinkedList[T]) RemoveLast() (T, bool) {
	if dll.tail == nil {
		var zero T
		return zero, false
	}

	value := dll.tail.Value
	dll.tail = dll.tail.Prev

	if dll.tail != nil {
		dll.tail.Next = nil
	} else {
		dll.head = nil
	}

	dll.size--
	return value, true
}

// RemoveAt removes the element at the specified index - O(n)
// Returns the removed element and true if successful, zero value and false if index is out of bounds
func (dll *DoublyLinkedList[T]) RemoveAt(index int) (T, bool) {
	if index < 0 || index >= dll.size {
		var zero T
		return zero, false
	}

	if index == 0 {
		return dll.RemoveFirst()
	}

	if index == dll.size-1 {
		return dll.RemoveLast()
	}

	current := dll.head
	for i := 0; i < index; i++ {
		current = current.Next
	}

	value := current.Value
	current.Prev.Next = current.Next
	current.Next.Prev = current.Prev

	dll.size--
	return value, true
}

// Get returns the element at the specified index - O(n)
// Returns the element and true if found, zero value and false if index is out of bounds
func (dll *DoublyLinkedList[T]) Get(index int) (T, bool) {
	if index < 0 || index >= dll.size {
		var zero T
		return zero, false
	}

	// Optimize: traverse from the closer end
	if index < dll.size/2 {
		current := dll.head
		for i := 0; i < index; i++ {
			current = current.Next
		}
		return current.Value, true
	}

	current := dll.tail
	for i := dll.size - 1; i > index; i-- {
		current = current.Prev
	}
	return current.Value, true
}

// First returns the first element without removing it - O(1)
// Returns the element and true if found, zero value and false if list is empty
func (dll *DoublyLinkedList[T]) First() (T, bool) {
	if dll.head == nil {
		var zero T
		return zero, false
	}

	return dll.head.Value, true
}

// Last returns the last element without removing it - O(1)
// Returns the element and true if found, zero value and false if list is empty
func (dll *DoublyLinkedList[T]) Last() (T, bool) {
	if dll.tail == nil {
		var zero T
		return zero, false
	}

	return dll.tail.Value, true
}

// Size returns the number of elements in the list - O(1)
func (dll *DoublyLinkedList[T]) Size() int {
	return dll.size
}

// IsEmpty returns true if the list is empty - O(1)
func (dll *DoublyLinkedList[T]) IsEmpty() bool {
	return dll.size == 0
}

// Clear removes all elements from the list - O(1)
func (dll *DoublyLinkedList[T]) Clear() {
	dll.head = nil
	dll.tail = nil
	dll.size = 0
}

// ToSlice converts the doubly linked list to a slice - O(n)
func (dll *DoublyLinkedList[T]) ToSlice() []T {
	result := make([]T, 0, dll.size)

	current := dll.head
	for current != nil {
		result = append(result, current.Value)
		current = current.Next
	}

	return result
}

// ToSliceReverse converts the doubly linked list to a slice in reverse order - O(n)
func (dll *DoublyLinkedList[T]) ToSliceReverse() []T {
	result := make([]T, 0, dll.size)

	current := dll.tail
	for current != nil {
		result = append(result, current.Value)
		current = current.Prev
	}

	return result
}

// ForEach iterates over each element in the list forward - O(n)
// The function receives the value and returns true to continue, false to stop
func (dll *DoublyLinkedList[T]) ForEach(fn func(T) bool) {
	current := dll.head
	for current != nil {
		if !fn(current.Value) {
			return
		}
		current = current.Next
	}
}

// ForEachReverse iterates over each element in the list in reverse order - O(n)
// The function receives the value and returns true to continue, false to stop
func (dll *DoublyLinkedList[T]) ForEachReverse(fn func(T) bool) {
	current := dll.tail
	for current != nil {
		if !fn(current.Value) {
			return
		}
		current = current.Prev
	}
}

// Find returns the first element for which the predicate returns true - O(n)
// Returns the element and true if found, zero value and false otherwise
func (dll *DoublyLinkedList[T]) Find(predicate func(T) bool) (T, bool) {
	current := dll.head
	for current != nil {
		if predicate(current.Value) {
			return current.Value, true
		}
		current = current.Next
	}

	var zero T
	return zero, false
}

// Contains returns true if the list contains an element matching the predicate - O(n)
func (dll *DoublyLinkedList[T]) Contains(predicate func(T) bool) bool {
	_, found := dll.Find(predicate)
	return found
}

// Reverse reverses the list in place - O(n)
func (dll *DoublyLinkedList[T]) Reverse() {
	if dll.head == nil || dll.head.Next == nil {
		return
	}

	current := dll.head
	dll.head, dll.tail = dll.tail, dll.head

	for current != nil {
		current.Next, current.Prev = current.Prev, current.Next
		current = current.Prev
	}
}

// Clone creates a deep copy of the list - O(n)
func (dll *DoublyLinkedList[T]) Clone() *DoublyLinkedList[T] {
	newList := New[T]()

	current := dll.head
	for current != nil {
		newList.Append(current.Value)
		current = current.Next
	}

	return newList
}
