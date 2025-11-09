package chelru

import (
	"sync"
	"testing"
)

func TestNew(t *testing.T) {
	cache := New[string, int](10)
	if cache.Capacity() != 10 {
		t.Errorf("expected capacity 10, got %d", cache.Capacity())
	}
	if cache.Len() != 0 {
		t.Errorf("expected length 0, got %d", cache.Len())
	}
}

func TestNew_PanicsOnInvalidCapacity(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic for capacity < 1")
		}
	}()
	New[string, int](0)
}

func TestPutAndGet(t *testing.T) {
	cache := New[string, int](2)

	cache.Put("a", 1)
	cache.Put("b", 2)

	if val, ok := cache.Get("a"); !ok || val != 1 {
		t.Errorf("expected a=1, got %d, %v", val, ok)
	}

	if val, ok := cache.Get("b"); !ok || val != 2 {
		t.Errorf("expected b=2, got %d, %v", val, ok)
	}
}

func TestPut_UpdatesExistingKey(t *testing.T) {
	cache := New[string, int](2)

	cache.Put("a", 1)
	cache.Put("a", 10)

	if val, ok := cache.Get("a"); !ok || val != 10 {
		t.Errorf("expected a=10, got %d, %v", val, ok)
	}

	if cache.Len() != 1 {
		t.Errorf("expected length 1, got %d", cache.Len())
	}
}

func TestEviction(t *testing.T) {
	cache := New[string, int](2)

	cache.Put("a", 1)
	cache.Put("b", 2)
	cache.Put("c", 3) // should evict "a"

	if _, ok := cache.Get("a"); ok {
		t.Error("expected 'a' to be evicted")
	}

	if val, ok := cache.Get("b"); !ok || val != 2 {
		t.Errorf("expected b=2, got %d, %v", val, ok)
	}

	if val, ok := cache.Get("c"); !ok || val != 3 {
		t.Errorf("expected c=3, got %d, %v", val, ok)
	}
}

func TestLRUOrder(t *testing.T) {
	cache := New[string, int](3)

	cache.Put("a", 1)
	cache.Put("b", 2)
	cache.Put("c", 3)

	// Access "a" to make it most recently used
	cache.Get("a")

	// Now order should be: a, c, b
	// Adding "d" should evict "b"
	cache.Put("d", 4)

	if _, ok := cache.Get("b"); ok {
		t.Error("expected 'b' to be evicted")
	}

	if _, ok := cache.Get("a"); !ok {
		t.Error("expected 'a' to still exist")
	}
}

func TestContains(t *testing.T) {
	cache := New[string, int](2)

	cache.Put("a", 1)

	if !cache.Contains("a") {
		t.Error("expected cache to contain 'a'")
	}

	if cache.Contains("b") {
		t.Error("expected cache not to contain 'b'")
	}
}

func TestRemove(t *testing.T) {
	cache := New[string, int](2)

	cache.Put("a", 1)
	cache.Put("b", 2)

	if !cache.Remove("a") {
		t.Error("expected Remove to return true")
	}

	if cache.Remove("a") {
		t.Error("expected second Remove to return false")
	}

	if cache.Contains("a") {
		t.Error("expected 'a' to be removed")
	}

	if cache.Len() != 1 {
		t.Errorf("expected length 1, got %d", cache.Len())
	}
}

func TestClear(t *testing.T) {
	cache := New[string, int](2)

	cache.Put("a", 1)
	cache.Put("b", 2)

	cache.Clear()

	if cache.Len() != 0 {
		t.Errorf("expected length 0 after clear, got %d", cache.Len())
	}

	if cache.Contains("a") || cache.Contains("b") {
		t.Error("expected cache to be empty after clear")
	}
}

func TestKeys(t *testing.T) {
	cache := New[string, int](3)

	cache.Put("a", 1)
	cache.Put("b", 2)
	cache.Put("c", 3)

	// Access "b" to make it most recently used
	cache.Get("b")

	keys := cache.Keys()

	// Expected order: b, c, a (b was accessed last, so it's first)
	if len(keys) != 3 {
		t.Errorf("expected 3 keys, got %d", len(keys))
	}

	if keys[0] != "b" {
		t.Errorf("expected first key to be 'b', got '%s'", keys[0])
	}
}

func TestConcurrentAccess(t *testing.T) {
	cache := New[int, int](100)
	var wg sync.WaitGroup

	// Concurrent writes
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(val int) {
			defer wg.Done()
			for j := 0; j < 100; j++ {
				cache.Put(val*100+j, val)
			}
		}(i)
	}

	// Concurrent reads
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 100; j++ {
				cache.Get(j)
			}
		}()
	}

	wg.Wait()

	// Just verify no panic occurred and cache is operational
	if cache.Len() > cache.Capacity() {
		t.Error("cache size exceeded capacity")
	}
}

func TestGetNonExistent(t *testing.T) {
	cache := New[string, int](2)

	val, ok := cache.Get("nonexistent")
	if ok {
		t.Error("expected Get to return false for non-existent key")
	}
	if val != 0 {
		t.Errorf("expected zero value, got %d", val)
	}
}

func TestCapacityOne(t *testing.T) {
	cache := New[string, int](1)

	cache.Put("a", 1)
	if val, ok := cache.Get("a"); !ok || val != 1 {
		t.Error("failed to get value from cache with capacity 1")
	}

	cache.Put("b", 2)
	if _, ok := cache.Get("a"); ok {
		t.Error("expected 'a' to be evicted")
	}
	if val, ok := cache.Get("b"); !ok || val != 2 {
		t.Error("failed to get 'b' from cache")
	}
}

func TestDifferentTypes(t *testing.T) {
	// Test with struct values
	type User struct {
		ID   int
		Name string
	}

	cache := New[int, User](2)
	cache.Put(1, User{ID: 1, Name: "Alice"})
	cache.Put(2, User{ID: 2, Name: "Bob"})

	user, ok := cache.Get(1)
	if !ok || user.Name != "Alice" {
		t.Error("failed to get struct value from cache")
	}
}

func BenchmarkPut(b *testing.B) {
	cache := New[int, int](1000)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		cache.Put(i%1000, i)
	}
}

func BenchmarkGet(b *testing.B) {
	cache := New[int, int](1000)
	for i := 0; i < 1000; i++ {
		cache.Put(i, i)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cache.Get(i % 1000)
	}
}

func BenchmarkPutAndGet(b *testing.B) {
	cache := New[int, int](1000)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		cache.Put(i%1000, i)
		cache.Get(i % 1000)
	}
}
