package chelru

import "sync"

// entry represents a key-value pair in the LRU cache
type entry[K comparable, V any] struct {
	key   K
	value V
	prev  *entry[K, V]
	next  *entry[K, V]
}

// LRU is a generic Least Recently Used cache with fixed capacity.
// It provides O(1) get and put operations.
type LRU[K comparable, V any] struct {
	capacity int
	cache    map[K]*entry[K, V]
	head     *entry[K, V] // most recently used
	tail     *entry[K, V] // least recently used
	mu       *sync.RWMutex
}

// New creates a new LRU cache with the specified capacity.
// Panics if capacity is less than 1.
func New[K comparable, V any](capacity int) *LRU[K, V] {
	if capacity < 1 {
		panic("chelru: capacity must be at least 1")
	}

	return &LRU[K, V]{
		capacity: capacity,
		cache:    make(map[K]*entry[K, V], capacity),
		mu:       &sync.RWMutex{},
	}
}

// Get retrieves a value from the cache.
// Returns the value and true if found, or zero value and false otherwise.
// Updates the item as most recently used.
func (l *LRU[K, V]) Get(key K) (V, bool) {
	l.mu.Lock()
	defer l.mu.Unlock()

	if node, found := l.cache[key]; found {
		l.moveToFront(node)
		return node.value, true
	}

	var zero V
	return zero, false
}

// Put adds or updates a value in the cache.
// If the key already exists, updates the value and marks it as most recently used.
// If the cache is at capacity, evicts the least recently used item.
func (l *LRU[K, V]) Put(key K, value V) {
	l.mu.Lock()
	defer l.mu.Unlock()

	// If key exists, update value and move to front
	if node, found := l.cache[key]; found {
		node.value = value
		l.moveToFront(node)
		return
	}

	// Create new entry
	node := &entry[K, V]{
		key:   key,
		value: value,
	}

	// Add to cache
	l.cache[key] = node
	l.addToFront(node)

	// Evict if over capacity
	if len(l.cache) > l.capacity {
		l.removeTail()
	}
}

// Contains checks if a key exists in the cache without updating its access time.
func (l *LRU[K, V]) Contains(key K) bool {
	l.mu.RLock()
	defer l.mu.RUnlock()

	_, found := l.cache[key]
	return found
}

// Remove removes a key from the cache.
// Returns true if the key was present, false otherwise.
func (l *LRU[K, V]) Remove(key K) bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	if node, found := l.cache[key]; found {
		l.removeNode(node)
		delete(l.cache, key)
		return true
	}

	return false
}

// Len returns the current number of items in the cache.
func (l *LRU[K, V]) Len() int {
	l.mu.RLock()
	defer l.mu.RUnlock()

	return len(l.cache)
}

// Capacity returns the maximum capacity of the cache.
func (l *LRU[K, V]) Capacity() int {
	return l.capacity
}

// Clear removes all items from the cache.
func (l *LRU[K, V]) Clear() {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.cache = make(map[K]*entry[K, V], l.capacity)
	l.head = nil
	l.tail = nil
}

// Keys returns all keys in the cache in order from most to least recently used.
func (l *LRU[K, V]) Keys() []K {
	l.mu.RLock()
	defer l.mu.RUnlock()

	keys := make([]K, 0, len(l.cache))
	for node := l.head; node != nil; node = node.next {
		keys = append(keys, node.key)
	}

	return keys
}

// moveToFront moves a node to the front of the list (most recently used)
func (l *LRU[K, V]) moveToFront(node *entry[K, V]) {
	if node == l.head {
		return
	}

	l.removeNode(node)
	l.addToFront(node)
}

// addToFront adds a node to the front of the list
func (l *LRU[K, V]) addToFront(node *entry[K, V]) {
	node.next = l.head
	node.prev = nil

	if l.head != nil {
		l.head.prev = node
	}

	l.head = node

	if l.tail == nil {
		l.tail = node
	}
}

// removeNode removes a node from the list
func (l *LRU[K, V]) removeNode(node *entry[K, V]) {
	if node.prev != nil {
		node.prev.next = node.next
	} else {
		l.head = node.next
	}

	if node.next != nil {
		node.next.prev = node.prev
	} else {
		l.tail = node.prev
	}
}

// removeTail removes the tail node (least recently used)
func (l *LRU[K, V]) removeTail() {
	if l.tail == nil {
		return
	}

	delete(l.cache, l.tail.key)
	l.removeNode(l.tail)
}
