package chemap

// Multimap is a generic map implementation where each key can have multiple values.
// It provides O(1) average-case lookups and insertions.
// Multimap is not thread-safe. For concurrent use, external synchronization is required.
type Multimap[K comparable, V any] struct {
	items map[K][]V
}

// NewMultimap creates and returns a new empty Multimap.
func NewMultimap[K comparable, V any]() *Multimap[K, V] {
	return &Multimap[K, V]{
		items: make(map[K][]V),
	}
}

// NewMultimapWithCapacity creates and returns a new empty Multimap with the specified initial capacity.
func NewMultimapWithCapacity[K comparable, V any](capacity int) *Multimap[K, V] {
	return &Multimap[K, V]{
		items: make(map[K][]V, capacity),
	}
}

// Put adds a value to the set of values for the given key.
func (m *Multimap[K, V]) Put(key K, value V) {
	m.items[key] = append(m.items[key], value)
}

// PutAll adds multiple values to the set of values for the given key.
func (m *Multimap[K, V]) PutAll(key K, values ...V) {
	m.items[key] = append(m.items[key], values...)
}

// Get returns all values associated with the given key.
// Returns an empty slice if the key doesn't exist.
func (m *Multimap[K, V]) Get(key K) []V {
	if values, exists := m.items[key]; exists {
		// Return a copy to prevent external modification
		result := make([]V, len(values))
		copy(result, values)
		return result
	}
	return []V{}
}

// GetFirst returns the first value associated with the given key.
// Returns the value and true if the key exists, or zero value and false otherwise.
func (m *Multimap[K, V]) GetFirst(key K) (V, bool) {
	if values, exists := m.items[key]; exists && len(values) > 0 {
		return values[0], true
	}
	var zero V
	return zero, false
}

// ContainsKey checks if the map contains the given key.
func (m *Multimap[K, V]) ContainsKey(key K) bool {
	_, exists := m.items[key]
	return exists
}

// ContainsEntry checks if the map contains the given key-value pair.
func (m *Multimap[K, V]) ContainsEntry(key K, value V, equals func(a, b V) bool) bool {
	if values, exists := m.items[key]; exists {
		for _, v := range values {
			if equals(v, value) {
				return true
			}
		}
	}
	return false
}

// Remove removes a specific value from the given key.
// Returns true if the value was found and removed, false otherwise.
func (m *Multimap[K, V]) Remove(key K, value V, equals func(a, b V) bool) bool {
	if values, exists := m.items[key]; exists {
		for i, v := range values {
			if equals(v, value) {
				// Remove the value at index i
				m.items[key] = append(values[:i], values[i+1:]...)
				// If no values left for this key, remove the key
				if len(m.items[key]) == 0 {
					delete(m.items, key)
				}
				return true
			}
		}
	}
	return false
}

// RemoveAll removes all values associated with the given key.
// Returns true if the key existed, false otherwise.
func (m *Multimap[K, V]) RemoveAll(key K) bool {
	if _, exists := m.items[key]; exists {
		delete(m.items, key)
		return true
	}
	return false
}

// Keys returns all keys in the multimap.
func (m *Multimap[K, V]) Keys() []K {
	keys := make([]K, 0, len(m.items))
	for k := range m.items {
		keys = append(keys, k)
	}
	return keys
}

// Values returns all values in the multimap (flattened across all keys).
func (m *Multimap[K, V]) Values() []V {
	var values []V
	for _, vals := range m.items {
		values = append(values, vals...)
	}
	return values
}

// Size returns the total number of key-value pairs in the multimap.
func (m *Multimap[K, V]) Size() int {
	count := 0
	for _, values := range m.items {
		count += len(values)
	}
	return count
}

// KeyCount returns the number of unique keys in the multimap.
func (m *Multimap[K, V]) KeyCount() int {
	return len(m.items)
}

// ValueCount returns the number of values associated with a specific key.
func (m *Multimap[K, V]) ValueCount(key K) int {
	if values, exists := m.items[key]; exists {
		return len(values)
	}
	return 0
}

// IsEmpty returns true if the multimap contains no key-value pairs.
func (m *Multimap[K, V]) IsEmpty() bool {
	return len(m.items) == 0
}

// Clear removes all key-value pairs from the multimap.
func (m *Multimap[K, V]) Clear() {
	m.items = make(map[K][]V)
}

// ForEach executes a function for each key-value pair in the multimap.
// If the function returns false, iteration stops.
func (m *Multimap[K, V]) ForEach(fn func(key K, value V) bool) {
	for key, values := range m.items {
		for _, value := range values {
			if !fn(key, value) {
				return
			}
		}
	}
}

// ForEachKey executes a function for each key and its associated values.
// If the function returns false, iteration stops.
func (m *Multimap[K, V]) ForEachKey(fn func(key K, values []V) bool) {
	for key, values := range m.items {
		// Make a copy to prevent modification
		valuesCopy := make([]V, len(values))
		copy(valuesCopy, values)
		if !fn(key, valuesCopy) {
			return
		}
	}
}

// Clone returns a shallow copy of the multimap.
func (m *Multimap[K, V]) Clone() *Multimap[K, V] {
	clone := NewMultimapWithCapacity[K, V](len(m.items))
	for key, values := range m.items {
		clone.items[key] = make([]V, len(values))
		copy(clone.items[key], values)
	}
	return clone
}

// Merge adds all entries from another multimap into this one.
func (m *Multimap[K, V]) Merge(other *Multimap[K, V]) {
	for key, values := range other.items {
		m.items[key] = append(m.items[key], values...)
	}
}

// ReplaceValues replaces all values for a given key with new values.
func (m *Multimap[K, V]) ReplaceValues(key K, values ...V) {
	if len(values) == 0 {
		delete(m.items, key)
	} else {
		m.items[key] = make([]V, len(values))
		copy(m.items[key], values)
	}
}
