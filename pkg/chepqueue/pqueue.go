package chepqueue

import (
	"fmt"

	"golang.org/x/exp/constraints"
)

// Ordered represents types that can be ordered.
type Ordered interface {
	constraints.Ordered
}

// Item represents an item in the priority queue with a value and priority.
type Item[T any, P Ordered] struct {
	Value    T
	Priority P
}

// PriorityQueue is a generic priority queue implementation using a binary heap.
// Lower priority values are dequeued first (min-heap by default).
type PriorityQueue[T any, P Ordered] struct {
	items   []Item[T, P]
	maxHeap bool
}

// New creates a new min-heap priority queue.
func New[T any, P Ordered]() *PriorityQueue[T, P] {
	return &PriorityQueue[T, P]{
		items:   make([]Item[T, P], 0),
		maxHeap: false,
	}
}

// NewMax creates a new max-heap priority queue.
// Higher priority values are dequeued first.
func NewMax[T any, P Ordered]() *PriorityQueue[T, P] {
	return &PriorityQueue[T, P]{
		items:   make([]Item[T, P], 0),
		maxHeap: true,
	}
}

// Push adds an item to the queue with the given priority.
func (pq *PriorityQueue[T, P]) Push(value T, priority P) {
	item := Item[T, P]{Value: value, Priority: priority}
	pq.items = append(pq.items, item)
	pq.heapifyUp(len(pq.items) - 1)
}

// Pop removes and returns the item with the highest priority.
// Panics if the queue is empty.
func (pq *PriorityQueue[T, P]) Pop() T {
	if pq.IsEmpty() {
		panic("chepqueue: Pop called on empty queue")
	}

	item := pq.items[0]
	lastIdx := len(pq.items) - 1

	pq.items[0] = pq.items[lastIdx]
	pq.items = pq.items[:lastIdx]

	if len(pq.items) > 0 {
		pq.heapifyDown(0)
	}

	return item.Value
}

// Peek returns the item with the highest priority without removing it.
// Panics if the queue is empty.
func (pq *PriorityQueue[T, P]) Peek() T {
	if pq.IsEmpty() {
		panic("chepqueue: Peek called on empty queue")
	}
	return pq.items[0].Value
}

// PeekPriority returns the priority of the item at the front of the queue.
// Panics if the queue is empty.
func (pq *PriorityQueue[T, P]) PeekPriority() P {
	if pq.IsEmpty() {
		panic("chepqueue: PeekPriority called on empty queue")
	}
	return pq.items[0].Priority
}

// IsEmpty returns true if the queue is empty.
func (pq *PriorityQueue[T, P]) IsEmpty() bool {
	return len(pq.items) == 0
}

// Len returns the number of items in the queue.
func (pq *PriorityQueue[T, P]) Len() int {
	return len(pq.items)
}

// Clear removes all items from the queue.
func (pq *PriorityQueue[T, P]) Clear() {
	pq.items = make([]Item[T, P], 0)
}

// Items returns a slice of all items in the queue (not in priority order).
func (pq *PriorityQueue[T, P]) Items() []Item[T, P] {
	result := make([]Item[T, P], len(pq.items))
	copy(result, pq.items)
	return result
}

// UpdatePriority finds an item by value and updates its priority.
// Returns true if the item was found and updated, false otherwise.
// This is an O(n) operation.
func (pq *PriorityQueue[T, P]) UpdatePriority(value T, newPriority P, equals func(T, T) bool) bool {
	for i, item := range pq.items {
		if equals(item.Value, value) {
			oldPriority := item.Priority
			pq.items[i].Priority = newPriority

			if pq.shouldSwap(newPriority, oldPriority) {
				pq.heapifyUp(i)
			} else {
				pq.heapifyDown(i)
			}
			return true
		}
	}
	return false
}

// Remove finds and removes an item by value.
// Returns true if the item was found and removed, false otherwise.
// This is an O(n) operation.
func (pq *PriorityQueue[T, P]) Remove(value T, equals func(T, T) bool) bool {
	for i, item := range pq.items {
		if equals(item.Value, value) {
			lastIdx := len(pq.items) - 1

			if i == lastIdx {
				pq.items = pq.items[:lastIdx]
				return true
			}

			pq.items[i] = pq.items[lastIdx]
			pq.items = pq.items[:lastIdx]

			if i < len(pq.items) {
				pq.heapifyDown(i)
				pq.heapifyUp(i)
			}

			return true
		}
	}
	return false
}

// heapifyUp restores the heap property by moving an item up.
func (pq *PriorityQueue[T, P]) heapifyUp(index int) {
	for index > 0 {
		parentIdx := (index - 1) / 2

		if !pq.shouldSwap(pq.items[index].Priority, pq.items[parentIdx].Priority) {
			break
		}

		pq.items[index], pq.items[parentIdx] = pq.items[parentIdx], pq.items[index]
		index = parentIdx
	}
}

// heapifyDown restores the heap property by moving an item down.
func (pq *PriorityQueue[T, P]) heapifyDown(index int) {
	for {
		leftChild := 2*index + 1
		rightChild := 2*index + 2
		smallest := index

		if leftChild < len(pq.items) && pq.shouldSwap(pq.items[leftChild].Priority, pq.items[smallest].Priority) {
			smallest = leftChild
		}

		if rightChild < len(pq.items) && pq.shouldSwap(pq.items[rightChild].Priority, pq.items[smallest].Priority) {
			smallest = rightChild
		}

		if smallest == index {
			break
		}

		pq.items[index], pq.items[smallest] = pq.items[smallest], pq.items[index]
		index = smallest
	}
}

// shouldSwap returns true if p1 should be swapped with p2 based on heap type.
func (pq *PriorityQueue[T, P]) shouldSwap(p1, p2 P) bool {
	if pq.maxHeap {
		return p1 > p2
	}
	return p1 < p2
}

// String returns a string representation of the priority queue.
func (pq *PriorityQueue[T, P]) String() string {
	if pq.IsEmpty() {
		return "PriorityQueue[]"
	}

	result := "PriorityQueue["
	for i, item := range pq.items {
		if i > 0 {
			result += ", "
		}
		result += fmt.Sprintf("{%v:%v}", item.Value, item.Priority)
	}
	result += "]"
	return result
}
