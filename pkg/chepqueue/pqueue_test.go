package chepqueue

import (
	"testing"
)

func TestPriorityQueue_MinHeap(t *testing.T) {
	pq := New[string, int]()

	pq.Push("low", 1)
	pq.Push("high", 10)
	pq.Push("medium", 5)

	if pq.Len() != 3 {
		t.Errorf("expected length 3, got %d", pq.Len())
	}

	// Min heap should return lowest priority first
	if pq.Pop() != "low" {
		t.Error("expected 'low' to be popped first")
	}

	if pq.Pop() != "medium" {
		t.Error("expected 'medium' to be popped second")
	}

	if pq.Pop() != "high" {
		t.Error("expected 'high' to be popped third")
	}

	if !pq.IsEmpty() {
		t.Error("queue should be empty")
	}
}

func TestPriorityQueue_MaxHeap(t *testing.T) {
	pq := NewMax[string, int]()

	pq.Push("low", 1)
	pq.Push("high", 10)
	pq.Push("medium", 5)

	// Max heap should return highest priority first
	if pq.Pop() != "high" {
		t.Error("expected 'high' to be popped first")
	}

	if pq.Pop() != "medium" {
		t.Error("expected 'medium' to be popped second")
	}

	if pq.Pop() != "low" {
		t.Error("expected 'low' to be popped third")
	}

	if !pq.IsEmpty() {
		t.Error("queue should be empty")
	}
}

func TestPriorityQueue_Peek(t *testing.T) {
	pq := New[string, int]()

	pq.Push("task", 5)

	// Peek should not remove the item
	if pq.Peek() != "task" {
		t.Error("expected 'task' from Peek")
	}

	if pq.Len() != 1 {
		t.Error("Peek should not remove items")
	}

	if pq.PeekPriority() != 5 {
		t.Error("expected priority 5")
	}
}

func TestPriorityQueue_Empty(t *testing.T) {
	pq := New[int, int]()

	if !pq.IsEmpty() {
		t.Error("new queue should be empty")
	}

	pq.Push(1, 1)

	if pq.IsEmpty() {
		t.Error("queue should not be empty after push")
	}

	pq.Pop()

	if !pq.IsEmpty() {
		t.Error("queue should be empty after popping all items")
	}
}

func TestPriorityQueue_PopPanics(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic when popping from empty queue")
		}
	}()

	pq := New[int, int]()
	pq.Pop()
}

func TestPriorityQueue_PeekPanics(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic when peeking empty queue")
		}
	}()

	pq := New[int, int]()
	pq.Peek()
}

func TestPriorityQueue_Clear(t *testing.T) {
	pq := New[int, int]()

	pq.Push(1, 1)
	pq.Push(2, 2)
	pq.Push(3, 3)

	pq.Clear()

	if !pq.IsEmpty() {
		t.Error("queue should be empty after clear")
	}

	if pq.Len() != 0 {
		t.Error("length should be 0 after clear")
	}
}

func TestPriorityQueue_Items(t *testing.T) {
	pq := New[string, int]()

	pq.Push("a", 1)
	pq.Push("b", 2)
	pq.Push("c", 3)

	items := pq.Items()

	if len(items) != 3 {
		t.Errorf("expected 3 items, got %d", len(items))
	}

	// Original queue should not be modified
	if pq.Len() != 3 {
		t.Error("Items() should not modify queue")
	}
}

func TestPriorityQueue_UpdatePriority(t *testing.T) {
	pq := New[string, int]()

	pq.Push("low", 1)
	pq.Push("medium", 5)
	pq.Push("high", 10)

	equals := func(a, b string) bool { return a == b }

	// Update "high" to have lowest priority
	updated := pq.UpdatePriority("high", 0, equals)

	if !updated {
		t.Error("UpdatePriority should return true when item is found")
	}

	// Now "high" should be popped first
	if pq.Pop() != "high" {
		t.Error("expected 'high' with updated priority to be popped first")
	}

	// Test updating non-existent item
	updated = pq.UpdatePriority("nonexistent", 0, equals)

	if updated {
		t.Error("UpdatePriority should return false when item not found")
	}
}

func TestPriorityQueue_Remove(t *testing.T) {
	pq := New[string, int]()

	pq.Push("a", 1)
	pq.Push("b", 2)
	pq.Push("c", 3)

	equals := func(a, b string) bool { return a == b }

	// Remove middle item
	removed := pq.Remove("b", equals)

	if !removed {
		t.Error("Remove should return true when item is found")
	}

	if pq.Len() != 2 {
		t.Errorf("expected length 2 after remove, got %d", pq.Len())
	}

	// Verify "b" is gone
	if pq.Pop() != "a" {
		t.Error("expected 'a' to be first")
	}

	if pq.Pop() != "c" {
		t.Error("expected 'c' to be second")
	}

	// Test removing non-existent item
	removed = pq.Remove("nonexistent", equals)

	if removed {
		t.Error("Remove should return false when item not found")
	}
}

func TestPriorityQueue_RemoveLast(t *testing.T) {
	pq := New[int, int]()

	pq.Push(1, 1)
	pq.Push(2, 2)

	equals := func(a, b int) bool { return a == b }

	// Remove last item
	removed := pq.Remove(2, equals)

	if !removed {
		t.Error("Remove should return true")
	}

	if pq.Len() != 1 {
		t.Errorf("expected length 1, got %d", pq.Len())
	}

	if pq.Pop() != 1 {
		t.Error("expected 1 to remain")
	}
}

func TestPriorityQueue_SamePriority(t *testing.T) {
	pq := New[string, int]()

	// All items have same priority
	pq.Push("first", 5)
	pq.Push("second", 5)
	pq.Push("third", 5)

	// Should still be able to pop all items
	count := 0
	for !pq.IsEmpty() {
		pq.Pop()
		count++
	}

	if count != 3 {
		t.Errorf("expected to pop 3 items, popped %d", count)
	}
}

func TestPriorityQueue_FloatPriority(t *testing.T) {
	pq := New[string, float64]()

	pq.Push("low", 1.5)
	pq.Push("high", 10.5)
	pq.Push("medium", 5.5)

	if pq.Pop() != "low" {
		t.Error("expected 'low' first")
	}

	if pq.Pop() != "medium" {
		t.Error("expected 'medium' second")
	}

	if pq.Pop() != "high" {
		t.Error("expected 'high' third")
	}
}

func TestPriorityQueue_StringPriority(t *testing.T) {
	pq := New[int, string]()

	pq.Push(1, "a")
	pq.Push(2, "c")
	pq.Push(3, "b")

	// Should be ordered alphabetically
	if pq.Pop() != 1 {
		t.Error("expected item with priority 'a' first")
	}

	if pq.Pop() != 3 {
		t.Error("expected item with priority 'b' second")
	}

	if pq.Pop() != 2 {
		t.Error("expected item with priority 'c' third")
	}
}

func TestPriorityQueue_ManyItems(t *testing.T) {
	pq := New[int, int]()

	// Push items in random order
	priorities := []int{50, 10, 80, 20, 90, 5, 30, 70, 40, 60}

	for i, p := range priorities {
		pq.Push(i, p)
	}

	// Pop all and verify they come out sorted
	lastPriority := -1

	for !pq.IsEmpty() {
		priority := pq.PeekPriority()

		if priority < lastPriority {
			t.Errorf("priorities out of order: %d came after %d", priority, lastPriority)
		}

		pq.Pop()
		lastPriority = priority
	}
}

func TestPriorityQueue_String(t *testing.T) {
	pq := New[string, int]()

	str := pq.String()
	if str != "PriorityQueue[]" {
		t.Errorf("empty queue string should be 'PriorityQueue[]', got %s", str)
	}

	pq.Push("task", 5)
	str = pq.String()

	if str == "" {
		t.Error("string representation should not be empty")
	}
}

func BenchmarkPriorityQueue_Push(b *testing.B) {
	pq := New[int, int]()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		pq.Push(i, i)
	}
}

func BenchmarkPriorityQueue_Pop(b *testing.B) {
	pq := New[int, int]()

	// Pre-fill the queue
	for i := 0; i < b.N; i++ {
		pq.Push(i, i)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		pq.Pop()
	}
}

func BenchmarkPriorityQueue_PushPop(b *testing.B) {
	pq := New[int, int]()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		pq.Push(i, i)
		if i%2 == 0 && !pq.IsEmpty() {
			pq.Pop()
		}
	}
}
