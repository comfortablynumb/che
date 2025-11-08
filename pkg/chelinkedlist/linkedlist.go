package chelinkedlist

// Node represents a single node in a singly linked list
type Node[T any] struct {
	Value T
	Next  *Node[T]
}

// LinkedList is a singly linked list implementation with O(1) prepend
// and O(n) append operations
type LinkedList[T any] struct {
	head *Node[T]
	tail *Node[T]
	size int
}

// New creates a new empty LinkedList
func New[T any]() *LinkedList[T] {
	return &LinkedList[T]{
		head: nil,
		tail: nil,
		size: 0,
	}
}

// Prepend adds an element to the beginning of the list - O(1)
func (ll *LinkedList[T]) Prepend(value T) {
	newNode := &Node[T]{Value: value, Next: ll.head}

	ll.head = newNode

	if ll.tail == nil {
		ll.tail = newNode
	}

	ll.size++
}

// Append adds an element to the end of the list - O(1)
func (ll *LinkedList[T]) Append(value T) {
	newNode := &Node[T]{Value: value, Next: nil}

	if ll.tail == nil {
		ll.head = newNode
		ll.tail = newNode
	} else {
		ll.tail.Next = newNode
		ll.tail = newNode
	}

	ll.size++
}

// InsertAt inserts an element at the specified index - O(n)
// Returns false if index is out of bounds
func (ll *LinkedList[T]) InsertAt(index int, value T) bool {
	if index < 0 || index > ll.size {
		return false
	}

	if index == 0 {
		ll.Prepend(value)
		return true
	}

	if index == ll.size {
		ll.Append(value)
		return true
	}

	current := ll.head
	for i := 0; i < index-1; i++ {
		current = current.Next
	}

	newNode := &Node[T]{Value: value, Next: current.Next}
	current.Next = newNode
	ll.size++

	return true
}

// RemoveFirst removes and returns the first element - O(1)
// Returns the element and true if successful, zero value and false if list is empty
func (ll *LinkedList[T]) RemoveFirst() (T, bool) {
	if ll.head == nil {
		var zero T
		return zero, false
	}

	value := ll.head.Value
	ll.head = ll.head.Next

	if ll.head == nil {
		ll.tail = nil
	}

	ll.size--
	return value, true
}

// RemoveLast removes and returns the last element - O(n)
// Returns the element and true if successful, zero value and false if list is empty
func (ll *LinkedList[T]) RemoveLast() (T, bool) {
	if ll.head == nil {
		var zero T
		return zero, false
	}

	if ll.head == ll.tail {
		value := ll.head.Value
		ll.head = nil
		ll.tail = nil
		ll.size--
		return value, true
	}

	current := ll.head
	for current.Next != ll.tail {
		current = current.Next
	}

	value := ll.tail.Value
	ll.tail = current
	ll.tail.Next = nil
	ll.size--

	return value, true
}

// RemoveAt removes the element at the specified index - O(n)
// Returns the removed element and true if successful, zero value and false if index is out of bounds
func (ll *LinkedList[T]) RemoveAt(index int) (T, bool) {
	if index < 0 || index >= ll.size {
		var zero T
		return zero, false
	}

	if index == 0 {
		return ll.RemoveFirst()
	}

	current := ll.head
	for i := 0; i < index-1; i++ {
		current = current.Next
	}

	value := current.Next.Value
	current.Next = current.Next.Next

	if current.Next == nil {
		ll.tail = current
	}

	ll.size--
	return value, true
}

// Get returns the element at the specified index - O(n)
// Returns the element and true if found, zero value and false if index is out of bounds
func (ll *LinkedList[T]) Get(index int) (T, bool) {
	if index < 0 || index >= ll.size {
		var zero T
		return zero, false
	}

	current := ll.head
	for i := 0; i < index; i++ {
		current = current.Next
	}

	return current.Value, true
}

// First returns the first element without removing it - O(1)
// Returns the element and true if found, zero value and false if list is empty
func (ll *LinkedList[T]) First() (T, bool) {
	if ll.head == nil {
		var zero T
		return zero, false
	}

	return ll.head.Value, true
}

// Last returns the last element without removing it - O(1)
// Returns the element and true if found, zero value and false if list is empty
func (ll *LinkedList[T]) Last() (T, bool) {
	if ll.tail == nil {
		var zero T
		return zero, false
	}

	return ll.tail.Value, true
}

// Size returns the number of elements in the list - O(1)
func (ll *LinkedList[T]) Size() int {
	return ll.size
}

// IsEmpty returns true if the list is empty - O(1)
func (ll *LinkedList[T]) IsEmpty() bool {
	return ll.size == 0
}

// Clear removes all elements from the list - O(1)
func (ll *LinkedList[T]) Clear() {
	ll.head = nil
	ll.tail = nil
	ll.size = 0
}

// ToSlice converts the linked list to a slice - O(n)
func (ll *LinkedList[T]) ToSlice() []T {
	result := make([]T, 0, ll.size)

	current := ll.head
	for current != nil {
		result = append(result, current.Value)
		current = current.Next
	}

	return result
}

// ForEach iterates over each element in the list - O(n)
// The function receives the value and returns true to continue, false to stop
func (ll *LinkedList[T]) ForEach(fn func(T) bool) {
	current := ll.head
	for current != nil {
		if !fn(current.Value) {
			return
		}
		current = current.Next
	}
}

// Find returns the first element for which the predicate returns true - O(n)
// Returns the element and true if found, zero value and false otherwise
func (ll *LinkedList[T]) Find(predicate func(T) bool) (T, bool) {
	current := ll.head
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
func (ll *LinkedList[T]) Contains(predicate func(T) bool) bool {
	_, found := ll.Find(predicate)
	return found
}

// Reverse reverses the list in place - O(n)
func (ll *LinkedList[T]) Reverse() {
	if ll.head == nil || ll.head.Next == nil {
		return
	}

	var prev *Node[T]
	current := ll.head
	ll.tail = ll.head

	for current != nil {
		next := current.Next
		current.Next = prev
		prev = current
		current = next
	}

	ll.head = prev
}

// Clone creates a deep copy of the list - O(n)
func (ll *LinkedList[T]) Clone() *LinkedList[T] {
	newList := New[T]()

	current := ll.head
	for current != nil {
		newList.Append(current.Value)
		current = current.Next
	}

	return newList
}
