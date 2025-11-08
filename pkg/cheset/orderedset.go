package cheset

import "fmt"

// OrderedSet is a generic set implementation that maintains insertion order.
// It combines O(1) average-case lookups (via a map) with ordered iteration (via a slice).
// OrderedSet is not thread-safe. For concurrent use, external synchronization is required.
type OrderedSet[T comparable] struct {
	items   map[T]int // maps item to its index in the order slice
	order   []T       // maintains insertion order
}

// NewOrdered creates and returns a new empty OrderedSet.
func NewOrdered[T comparable]() *OrderedSet[T] {
	return &OrderedSet[T]{
		items: make(map[T]int),
		order: make([]T, 0),
	}
}

// NewOrderedWithCapacity creates and returns a new empty OrderedSet with the specified initial capacity.
// This can improve performance when the expected size is known in advance.
func NewOrderedWithCapacity[T comparable](capacity int) *OrderedSet[T] {
	return &OrderedSet[T]{
		items: make(map[T]int, capacity),
		order: make([]T, 0, capacity),
	}
}

// NewOrderedFromSlice creates and returns a new OrderedSet containing all elements from the given slice.
// Insertion order is preserved based on first occurrence in the slice.
func NewOrderedFromSlice[T comparable](slice []T) *OrderedSet[T] {
	set := NewOrderedWithCapacity[T](len(slice))
	for _, item := range slice {
		set.Add(item)
	}
	return set
}

// Add adds an element to the set. If the element already exists, this is a no-op.
// Returns true if the element was added, false if it already existed.
func (s *OrderedSet[T]) Add(item T) bool {
	if _, exists := s.items[item]; exists {
		return false
	}
	s.items[item] = len(s.order)
	s.order = append(s.order, item)
	return true
}

// AddMultiple adds multiple elements to the set in order.
// Returns the number of elements that were actually added (excluding duplicates).
func (s *OrderedSet[T]) AddMultiple(items ...T) int {
	count := 0
	for _, item := range items {
		if s.Add(item) {
			count++
		}
	}
	return count
}

// Remove removes an element from the set.
// Returns true if the element was removed, false if it didn't exist.
// Note: This operation is O(n) because it requires updating indices.
func (s *OrderedSet[T]) Remove(item T) bool {
	idx, exists := s.items[item]
	if !exists {
		return false
	}

	// Remove from map
	delete(s.items, item)

	// Remove from order slice
	s.order = append(s.order[:idx], s.order[idx+1:]...)

	// Update indices for all items after the removed one
	for i := idx; i < len(s.order); i++ {
		s.items[s.order[i]] = i
	}

	return true
}

// RemoveMultiple removes multiple elements from the set.
// Returns the number of elements that were actually removed.
func (s *OrderedSet[T]) RemoveMultiple(items ...T) int {
	count := 0
	for _, item := range items {
		if s.Remove(item) {
			count++
		}
	}
	return count
}

// Contains checks if an element exists in the set.
func (s *OrderedSet[T]) Contains(item T) bool {
	_, exists := s.items[item]
	return exists
}

// ContainsAll checks if all given elements exist in the set.
func (s *OrderedSet[T]) ContainsAll(items ...T) bool {
	for _, item := range items {
		if !s.Contains(item) {
			return false
		}
	}
	return true
}

// ContainsAny checks if any of the given elements exist in the set.
func (s *OrderedSet[T]) ContainsAny(items ...T) bool {
	for _, item := range items {
		if s.Contains(item) {
			return true
		}
	}
	return false
}

// Size returns the number of elements in the set.
func (s *OrderedSet[T]) Size() int {
	return len(s.order)
}

// IsEmpty returns true if the set contains no elements.
func (s *OrderedSet[T]) IsEmpty() bool {
	return s.Size() == 0
}

// Clear removes all elements from the set.
func (s *OrderedSet[T]) Clear() {
	s.items = make(map[T]int)
	s.order = make([]T, 0)
}

// ToSlice returns a slice containing all elements in insertion order.
func (s *OrderedSet[T]) ToSlice() []T {
	result := make([]T, len(s.order))
	copy(result, s.order)
	return result
}

// Clone returns a shallow copy of the set with insertion order preserved.
func (s *OrderedSet[T]) Clone() *OrderedSet[T] {
	clone := NewOrderedWithCapacity[T](s.Size())
	for _, item := range s.order {
		clone.Add(item)
	}
	return clone
}

// Equal checks if this set contains exactly the same elements as another set in the same order.
func (s *OrderedSet[T]) Equal(other *OrderedSet[T]) bool {
	if s.Size() != other.Size() {
		return false
	}
	for i, item := range s.order {
		if other.order[i] != item {
			return false
		}
	}
	return true
}

// Union returns a new ordered set containing all elements from both sets.
// Elements from this set appear first in insertion order, followed by new elements from the other set.
func (s *OrderedSet[T]) Union(other *OrderedSet[T]) *OrderedSet[T] {
	result := s.Clone()
	for _, item := range other.order {
		result.Add(item)
	}
	return result
}

// Intersect returns a new ordered set containing only elements that are in both sets.
// The order is determined by the first set's insertion order.
func (s *OrderedSet[T]) Intersect(other *OrderedSet[T]) *OrderedSet[T] {
	result := NewOrdered[T]()
	for _, item := range s.order {
		if other.Contains(item) {
			result.Add(item)
		}
	}
	return result
}

// Diff returns a new ordered set containing elements that are in this set but not in the other set.
// The order is preserved from the first set.
func (s *OrderedSet[T]) Diff(other *OrderedSet[T]) *OrderedSet[T] {
	result := NewOrdered[T]()
	for _, item := range s.order {
		if !other.Contains(item) {
			result.Add(item)
		}
	}
	return result
}

// SymmetricDiff returns a new ordered set containing elements that are in either set but not in both.
// Elements from this set appear first, followed by elements unique to the other set.
func (s *OrderedSet[T]) SymmetricDiff(other *OrderedSet[T]) *OrderedSet[T] {
	result := NewOrdered[T]()
	for _, item := range s.order {
		if !other.Contains(item) {
			result.Add(item)
		}
	}
	for _, item := range other.order {
		if !s.Contains(item) {
			result.Add(item)
		}
	}
	return result
}

// IsSubset checks if this set is a subset of another set (all elements in this set are in the other set).
func (s *OrderedSet[T]) IsSubset(other *OrderedSet[T]) bool {
	if s.Size() > other.Size() {
		return false
	}
	for _, item := range s.order {
		if !other.Contains(item) {
			return false
		}
	}
	return true
}

// IsSuperset checks if this set is a superset of another set (all elements in the other set are in this set).
func (s *OrderedSet[T]) IsSuperset(other *OrderedSet[T]) bool {
	return other.IsSubset(s)
}

// IsProperSubset checks if this set is a proper subset of another set.
func (s *OrderedSet[T]) IsProperSubset(other *OrderedSet[T]) bool {
	return s.Size() < other.Size() && s.IsSubset(other)
}

// IsProperSuperset checks if this set is a proper superset of another set.
func (s *OrderedSet[T]) IsProperSuperset(other *OrderedSet[T]) bool {
	return other.IsProperSubset(s)
}

// IsDisjoint checks if this set has no elements in common with another set.
func (s *OrderedSet[T]) IsDisjoint(other *OrderedSet[T]) bool {
	smaller, larger := s, other
	if other.Size() < s.Size() {
		smaller, larger = other, s
	}
	for _, item := range smaller.order {
		if larger.Contains(item) {
			return false
		}
	}
	return true
}

// ForEach executes a function for each element in the set in insertion order.
// If the function returns false, iteration stops.
func (s *OrderedSet[T]) ForEach(fn func(item T) bool) {
	for _, item := range s.order {
		if !fn(item) {
			return
		}
	}
}

// Filter returns a new ordered set containing only elements that satisfy the predicate.
// The insertion order is preserved.
func (s *OrderedSet[T]) Filter(predicate func(item T) bool) *OrderedSet[T] {
	result := NewOrdered[T]()
	for _, item := range s.order {
		if predicate(item) {
			result.Add(item)
		}
	}
	return result
}

// GetAt returns the element at the specified index in insertion order.
// Panics if the index is out of bounds.
func (s *OrderedSet[T]) GetAt(index int) T {
	if index < 0 || index >= len(s.order) {
		panic("index out of bounds")
	}
	return s.order[index]
}

// Index returns the insertion order index of the element, or -1 if not found.
func (s *OrderedSet[T]) Index(item T) int {
	if idx, exists := s.items[item]; exists {
		return idx
	}
	return -1
}

// First returns the first element in insertion order and true.
// Returns zero value and false if the set is empty.
func (s *OrderedSet[T]) First() (T, bool) {
	if s.IsEmpty() {
		var zero T
		return zero, false
	}
	return s.order[0], true
}

// Last returns the last element in insertion order and true.
// Returns zero value and false if the set is empty.
func (s *OrderedSet[T]) Last() (T, bool) {
	if s.IsEmpty() {
		var zero T
		return zero, false
	}
	return s.order[len(s.order)-1], true
}

// PopFirst removes and returns the first element in insertion order and true.
// Returns zero value and false if the set is empty.
func (s *OrderedSet[T]) PopFirst() (T, bool) {
	if s.IsEmpty() {
		var zero T
		return zero, false
	}
	first := s.order[0]
	s.Remove(first)
	return first, true
}

// PopLast removes and returns the last element in insertion order and true.
// Returns zero value and false if the set is empty.
func (s *OrderedSet[T]) PopLast() (T, bool) {
	if s.IsEmpty() {
		var zero T
		return zero, false
	}
	last := s.order[len(s.order)-1]
	s.Remove(last)
	return last, true
}

// String returns a string representation of the set in insertion order.
func (s *OrderedSet[T]) String() string {
	if s.IsEmpty() {
		return "OrderedSet[]"
	}

	result := "OrderedSet["
	for i, item := range s.order {
		if i > 0 {
			result += ", "
		}
		result += fmt.Sprint(item)
	}
	result += "]"
	return result
}
