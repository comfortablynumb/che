package cheset

import "fmt"

// HashSet is a generic set implementation backed by a map.
// It provides O(1) average-case performance for basic operations (Add, Remove, Contains).
// HashSet is not thread-safe. For concurrent use, external synchronization is required.
type HashSet[T comparable] struct {
	items map[T]struct{}
}

// New creates and returns a new empty HashSet.
func New[T comparable]() *HashSet[T] {
	return &HashSet[T]{
		items: make(map[T]struct{}),
	}
}

// NewWithCapacity creates and returns a new empty HashSet with the specified initial capacity.
// This can improve performance when the expected size is known in advance.
func NewWithCapacity[T comparable](capacity int) *HashSet[T] {
	return &HashSet[T]{
		items: make(map[T]struct{}, capacity),
	}
}

// NewFromSlice creates and returns a new HashSet containing all elements from the given slice.
func NewFromSlice[T comparable](slice []T) *HashSet[T] {
	set := NewWithCapacity[T](len(slice))
	for _, item := range slice {
		set.Add(item)
	}
	return set
}

// Add adds an element to the set. If the element already exists, this is a no-op.
// Returns true if the element was added, false if it already existed.
func (s *HashSet[T]) Add(item T) bool {
	if s.Contains(item) {
		return false
	}
	s.items[item] = struct{}{}
	return true
}

// AddMultiple adds multiple elements to the set.
// Returns the number of elements that were actually added (excluding duplicates).
func (s *HashSet[T]) AddMultiple(items ...T) int {
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
func (s *HashSet[T]) Remove(item T) bool {
	if !s.Contains(item) {
		return false
	}
	delete(s.items, item)
	return true
}

// RemoveMultiple removes multiple elements from the set.
// Returns the number of elements that were actually removed.
func (s *HashSet[T]) RemoveMultiple(items ...T) int {
	count := 0
	for _, item := range items {
		if s.Remove(item) {
			count++
		}
	}
	return count
}

// Contains checks if an element exists in the set.
func (s *HashSet[T]) Contains(item T) bool {
	_, exists := s.items[item]
	return exists
}

// ContainsAll checks if all given elements exist in the set.
func (s *HashSet[T]) ContainsAll(items ...T) bool {
	for _, item := range items {
		if !s.Contains(item) {
			return false
		}
	}
	return true
}

// ContainsAny checks if any of the given elements exist in the set.
func (s *HashSet[T]) ContainsAny(items ...T) bool {
	for _, item := range items {
		if s.Contains(item) {
			return true
		}
	}
	return false
}

// Size returns the number of elements in the set.
func (s *HashSet[T]) Size() int {
	return len(s.items)
}

// IsEmpty returns true if the set contains no elements.
func (s *HashSet[T]) IsEmpty() bool {
	return s.Size() == 0
}

// Clear removes all elements from the set.
func (s *HashSet[T]) Clear() {
	s.items = make(map[T]struct{})
}

// ToSlice returns a slice containing all elements in the set.
// The order of elements is not guaranteed.
func (s *HashSet[T]) ToSlice() []T {
	result := make([]T, 0, len(s.items))
	for item := range s.items {
		result = append(result, item)
	}
	return result
}

// Clone returns a shallow copy of the set.
func (s *HashSet[T]) Clone() *HashSet[T] {
	clone := NewWithCapacity[T](s.Size())
	for item := range s.items {
		clone.Add(item)
	}
	return clone
}

// Equal checks if this set contains exactly the same elements as another set.
func (s *HashSet[T]) Equal(other *HashSet[T]) bool {
	if s.Size() != other.Size() {
		return false
	}
	for item := range s.items {
		if !other.Contains(item) {
			return false
		}
	}
	return true
}

// Union returns a new set containing all elements that are in either this set or the other set.
func (s *HashSet[T]) Union(other *HashSet[T]) *HashSet[T] {
	result := s.Clone()
	for item := range other.items {
		result.Add(item)
	}
	return result
}

// Intersect returns a new set containing only elements that are in both this set and the other set.
func (s *HashSet[T]) Intersect(other *HashSet[T]) *HashSet[T] {
	result := New[T]()
	// Iterate over the smaller set for efficiency
	smaller, larger := s, other
	if other.Size() < s.Size() {
		smaller, larger = other, s
	}
	for item := range smaller.items {
		if larger.Contains(item) {
			result.Add(item)
		}
	}
	return result
}

// Diff returns a new set containing elements that are in this set but not in the other set.
func (s *HashSet[T]) Diff(other *HashSet[T]) *HashSet[T] {
	result := New[T]()
	for item := range s.items {
		if !other.Contains(item) {
			result.Add(item)
		}
	}
	return result
}

// SymmetricDiff returns a new set containing elements that are in either set but not in both.
func (s *HashSet[T]) SymmetricDiff(other *HashSet[T]) *HashSet[T] {
	result := New[T]()
	for item := range s.items {
		if !other.Contains(item) {
			result.Add(item)
		}
	}
	for item := range other.items {
		if !s.Contains(item) {
			result.Add(item)
		}
	}
	return result
}

// IsSubset checks if this set is a subset of another set (all elements in this set are in the other set).
func (s *HashSet[T]) IsSubset(other *HashSet[T]) bool {
	if s.Size() > other.Size() {
		return false
	}
	for item := range s.items {
		if !other.Contains(item) {
			return false
		}
	}
	return true
}

// IsSuperset checks if this set is a superset of another set (all elements in the other set are in this set).
func (s *HashSet[T]) IsSuperset(other *HashSet[T]) bool {
	return other.IsSubset(s)
}

// IsProperSubset checks if this set is a proper subset of another set
// (all elements in this set are in the other set, and the other set has more elements).
func (s *HashSet[T]) IsProperSubset(other *HashSet[T]) bool {
	return s.Size() < other.Size() && s.IsSubset(other)
}

// IsProperSuperset checks if this set is a proper superset of another set
// (all elements in the other set are in this set, and this set has more elements).
func (s *HashSet[T]) IsProperSuperset(other *HashSet[T]) bool {
	return other.IsProperSubset(s)
}

// IsDisjoint checks if this set has no elements in common with another set.
func (s *HashSet[T]) IsDisjoint(other *HashSet[T]) bool {
	// Iterate over the smaller set for efficiency
	smaller, larger := s, other
	if other.Size() < s.Size() {
		smaller, larger = other, s
	}
	for item := range smaller.items {
		if larger.Contains(item) {
			return false
		}
	}
	return true
}

// ForEach executes a function for each element in the set.
// If the function returns false, iteration stops.
func (s *HashSet[T]) ForEach(fn func(item T) bool) {
	for item := range s.items {
		if !fn(item) {
			return
		}
	}
}

// Filter returns a new set containing only elements that satisfy the predicate.
func (s *HashSet[T]) Filter(predicate func(item T) bool) *HashSet[T] {
	result := New[T]()
	for item := range s.items {
		if predicate(item) {
			result.Add(item)
		}
	}
	return result
}

// String returns a string representation of the set.
func (s *HashSet[T]) String() string {
	if s.IsEmpty() {
		return "HashSet{}"
	}

	result := "HashSet{"
	first := true
	for item := range s.items {
		if !first {
			result += ", "
		}
		result += fmt.Sprint(item)
		first = false
	}
	result += "}"
	return result
}
