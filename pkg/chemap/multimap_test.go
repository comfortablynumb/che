package chemap

import (
	"testing"

	"github.com/comfortablynumb/che/pkg/chetest"
)

// TestNewMultimap tests creating a new empty Multimap
func TestNewMultimap(t *testing.T) {
	mm := NewMultimap[string, int]()

	chetest.RequireEqual(t, mm.IsEmpty(), true)
	chetest.RequireEqual(t, mm.KeyCount(), 0)
	chetest.RequireEqual(t, mm.Size(), 0)
}

// TestNewMultimapWithCapacity tests creating a Multimap with capacity
func TestNewMultimapWithCapacity(t *testing.T) {
	mm := NewMultimapWithCapacity[string, int](10)

	chetest.RequireEqual(t, mm.IsEmpty(), true)
	chetest.RequireEqual(t, mm.KeyCount(), 0)
}

// TestPut tests adding values
func TestPut(t *testing.T) {
	mm := NewMultimap[string, int]()

	mm.Put("key1", 1)
	mm.Put("key1", 2)
	mm.Put("key2", 3)

	chetest.RequireEqual(t, mm.KeyCount(), 2)
	chetest.RequireEqual(t, mm.Size(), 3)
}

// TestPutAll tests adding multiple values at once
func TestPutAll(t *testing.T) {
	mm := NewMultimap[string, int]()

	mm.PutAll("key1", 1, 2, 3)
	mm.PutAll("key2", 4, 5)

	chetest.RequireEqual(t, mm.KeyCount(), 2)
	chetest.RequireEqual(t, mm.Size(), 5)
	chetest.RequireEqual(t, mm.ValueCount("key1"), 3)
	chetest.RequireEqual(t, mm.ValueCount("key2"), 2)
}

// TestGet tests retrieving values
func TestGet(t *testing.T) {
	mm := NewMultimap[string, int]()

	mm.Put("key1", 1)
	mm.Put("key1", 2)
	mm.Put("key1", 3)

	values := mm.Get("key1")
	chetest.RequireEqual(t, len(values), 3)
	chetest.RequireEqual(t, values[0], 1)
	chetest.RequireEqual(t, values[1], 2)
	chetest.RequireEqual(t, values[2], 3)
}

// TestGet_NonExistent tests getting values for non-existent key
func TestGet_NonExistent(t *testing.T) {
	mm := NewMultimap[string, int]()

	values := mm.Get("nonexistent")
	chetest.RequireEqual(t, len(values), 0)
}

// TestGetFirst tests retrieving first value
func TestGetFirst(t *testing.T) {
	mm := NewMultimap[string, int]()

	mm.PutAll("key1", 10, 20, 30)

	value, ok := mm.GetFirst("key1")
	chetest.RequireEqual(t, ok, true)
	chetest.RequireEqual(t, value, 10)
}

// TestGetFirst_NonExistent tests getting first value for non-existent key
func TestGetFirst_NonExistent(t *testing.T) {
	mm := NewMultimap[string, int]()

	value, ok := mm.GetFirst("nonexistent")
	chetest.RequireEqual(t, ok, false)
	chetest.RequireEqual(t, value, 0)
}

// TestContainsKey tests checking if key exists
func TestContainsKey(t *testing.T) {
	mm := NewMultimap[string, int]()

	mm.Put("key1", 1)

	chetest.RequireEqual(t, mm.ContainsKey("key1"), true)
	chetest.RequireEqual(t, mm.ContainsKey("key2"), false)
}

// TestContainsEntry tests checking if key-value pair exists
func TestContainsEntry(t *testing.T) {
	mm := NewMultimap[string, int]()

	mm.PutAll("key1", 1, 2, 3)

	equals := func(a, b int) bool { return a == b }

	chetest.RequireEqual(t, mm.ContainsEntry("key1", 2, equals), true)
	chetest.RequireEqual(t, mm.ContainsEntry("key1", 5, equals), false)
	chetest.RequireEqual(t, mm.ContainsEntry("key2", 1, equals), false)
}

// TestRemove tests removing specific value
func TestRemove(t *testing.T) {
	mm := NewMultimap[string, int]()

	mm.PutAll("key1", 1, 2, 3)

	equals := func(a, b int) bool { return a == b }

	removed := mm.Remove("key1", 2, equals)
	chetest.RequireEqual(t, removed, true)
	chetest.RequireEqual(t, mm.ValueCount("key1"), 2)

	values := mm.Get("key1")
	chetest.RequireEqual(t, len(values), 2)
	chetest.RequireEqual(t, values[0], 1)
	chetest.RequireEqual(t, values[1], 3)
}

// TestRemove_LastValue tests removing the last value for a key
func TestRemove_LastValue(t *testing.T) {
	mm := NewMultimap[string, int]()

	mm.Put("key1", 1)

	equals := func(a, b int) bool { return a == b }

	removed := mm.Remove("key1", 1, equals)
	chetest.RequireEqual(t, removed, true)
	chetest.RequireEqual(t, mm.ContainsKey("key1"), false)
	chetest.RequireEqual(t, mm.KeyCount(), 0)
}

// TestRemove_NonExistent tests removing non-existent value
func TestRemove_NonExistent(t *testing.T) {
	mm := NewMultimap[string, int]()

	mm.Put("key1", 1)

	equals := func(a, b int) bool { return a == b }

	removed := mm.Remove("key1", 5, equals)
	chetest.RequireEqual(t, removed, false)
	chetest.RequireEqual(t, mm.ValueCount("key1"), 1)
}

// TestRemoveAll tests removing all values for a key
func TestRemoveAll(t *testing.T) {
	mm := NewMultimap[string, int]()

	mm.PutAll("key1", 1, 2, 3)
	mm.PutAll("key2", 4, 5)

	removed := mm.RemoveAll("key1")
	chetest.RequireEqual(t, removed, true)
	chetest.RequireEqual(t, mm.ContainsKey("key1"), false)
	chetest.RequireEqual(t, mm.KeyCount(), 1)
	chetest.RequireEqual(t, mm.Size(), 2)
}

// TestRemoveAll_NonExistent tests removing all values for non-existent key
func TestRemoveAll_NonExistent(t *testing.T) {
	mm := NewMultimap[string, int]()

	removed := mm.RemoveAll("nonexistent")
	chetest.RequireEqual(t, removed, false)
}

// TestKeys tests getting all keys
func TestKeys(t *testing.T) {
	mm := NewMultimap[string, int]()

	mm.Put("key1", 1)
	mm.Put("key2", 2)
	mm.Put("key3", 3)

	keys := mm.Keys()
	chetest.RequireEqual(t, len(keys), 3)

	// Create a set to verify all keys are present
	keySet := make(map[string]bool)
	for _, k := range keys {
		keySet[k] = true
	}
	chetest.RequireEqual(t, keySet["key1"], true)
	chetest.RequireEqual(t, keySet["key2"], true)
	chetest.RequireEqual(t, keySet["key3"], true)
}

// TestValues tests getting all values
func TestValues(t *testing.T) {
	mm := NewMultimap[string, int]()

	mm.PutAll("key1", 1, 2)
	mm.PutAll("key2", 3, 4)

	values := mm.Values()
	chetest.RequireEqual(t, len(values), 4)
}

// TestSize tests getting total size
func TestSize(t *testing.T) {
	mm := NewMultimap[string, int]()

	chetest.RequireEqual(t, mm.Size(), 0)

	mm.Put("key1", 1)
	chetest.RequireEqual(t, mm.Size(), 1)

	mm.PutAll("key1", 2, 3)
	chetest.RequireEqual(t, mm.Size(), 3)

	mm.Put("key2", 4)
	chetest.RequireEqual(t, mm.Size(), 4)
}

// TestKeyCount tests getting key count
func TestKeyCount(t *testing.T) {
	mm := NewMultimap[string, int]()

	chetest.RequireEqual(t, mm.KeyCount(), 0)

	mm.Put("key1", 1)
	chetest.RequireEqual(t, mm.KeyCount(), 1)

	mm.Put("key1", 2)
	chetest.RequireEqual(t, mm.KeyCount(), 1)

	mm.Put("key2", 3)
	chetest.RequireEqual(t, mm.KeyCount(), 2)
}

// TestValueCount tests getting value count for specific key
func TestValueCount(t *testing.T) {
	mm := NewMultimap[string, int]()

	mm.PutAll("key1", 1, 2, 3)

	chetest.RequireEqual(t, mm.ValueCount("key1"), 3)
	chetest.RequireEqual(t, mm.ValueCount("nonexistent"), 0)
}

// TestIsEmpty tests checking if multimap is empty
func TestIsEmpty(t *testing.T) {
	mm := NewMultimap[string, int]()

	chetest.RequireEqual(t, mm.IsEmpty(), true)

	mm.Put("key1", 1)
	chetest.RequireEqual(t, mm.IsEmpty(), false)

	mm.RemoveAll("key1")
	chetest.RequireEqual(t, mm.IsEmpty(), true)
}

// TestClear tests clearing the multimap
func TestClear(t *testing.T) {
	mm := NewMultimap[string, int]()

	mm.PutAll("key1", 1, 2, 3)
	mm.PutAll("key2", 4, 5)

	mm.Clear()
	chetest.RequireEqual(t, mm.IsEmpty(), true)
	chetest.RequireEqual(t, mm.KeyCount(), 0)
	chetest.RequireEqual(t, mm.Size(), 0)
}

// TestForEach tests iterating over all entries
func TestForEach(t *testing.T) {
	mm := NewMultimap[string, int]()

	mm.PutAll("key1", 1, 2)
	mm.Put("key2", 3)

	count := 0
	mm.ForEach(func(key string, value int) bool {
		count++
		return true
	})

	chetest.RequireEqual(t, count, 3)
}

// TestForEach_EarlyExit tests early exit from iteration
func TestForEach_EarlyExit(t *testing.T) {
	mm := NewMultimap[string, int]()

	mm.PutAll("key1", 1, 2, 3, 4, 5)

	count := 0
	mm.ForEach(func(key string, value int) bool {
		count++
		return count < 3
	})

	chetest.RequireEqual(t, count, 3)
}

// TestForEachKey tests iterating over keys and their values
func TestForEachKey(t *testing.T) {
	mm := NewMultimap[string, int]()

	mm.PutAll("key1", 1, 2)
	mm.Put("key2", 3)

	keyCount := 0
	valueCount := 0
	mm.ForEachKey(func(key string, values []int) bool {
		keyCount++
		valueCount += len(values)
		return true
	})

	chetest.RequireEqual(t, keyCount, 2)
	chetest.RequireEqual(t, valueCount, 3)
}

// TestForEachKey_EarlyExit tests early exit from key iteration
func TestForEachKey_EarlyExit(t *testing.T) {
	mm := NewMultimap[string, int]()

	mm.Put("key1", 1)
	mm.Put("key2", 2)
	mm.Put("key3", 3)

	count := 0
	mm.ForEachKey(func(key string, values []int) bool {
		count++
		return count < 2
	})

	chetest.RequireEqual(t, count, 2)
}

// TestClone tests cloning a multimap
func TestClone(t *testing.T) {
	original := NewMultimap[string, int]()

	original.PutAll("key1", 1, 2, 3)
	original.Put("key2", 4)

	clone := original.Clone()

	chetest.RequireEqual(t, clone.KeyCount(), 2)
	chetest.RequireEqual(t, clone.Size(), 4)
	chetest.RequireEqual(t, clone.ValueCount("key1"), 3)

	// Modifying clone shouldn't affect original
	clone.Put("key3", 5)
	chetest.RequireEqual(t, original.ContainsKey("key3"), false)
	chetest.RequireEqual(t, clone.ContainsKey("key3"), true)
}

// TestMerge tests merging two multimaps
func TestMerge(t *testing.T) {
	mm1 := NewMultimap[string, int]()
	mm2 := NewMultimap[string, int]()

	mm1.PutAll("key1", 1, 2)
	mm2.PutAll("key1", 3, 4)
	mm2.Put("key2", 5)

	mm1.Merge(mm2)

	chetest.RequireEqual(t, mm1.KeyCount(), 2)
	chetest.RequireEqual(t, mm1.ValueCount("key1"), 4)
	chetest.RequireEqual(t, mm1.ValueCount("key2"), 1)
	chetest.RequireEqual(t, mm1.Size(), 5)
}

// TestReplaceValues tests replacing all values for a key
func TestReplaceValues(t *testing.T) {
	mm := NewMultimap[string, int]()

	mm.PutAll("key1", 1, 2, 3)

	mm.ReplaceValues("key1", 10, 20)

	values := mm.Get("key1")
	chetest.RequireEqual(t, len(values), 2)
	chetest.RequireEqual(t, values[0], 10)
	chetest.RequireEqual(t, values[1], 20)
}

// TestReplaceValues_Empty tests replacing with empty values
func TestReplaceValues_Empty(t *testing.T) {
	mm := NewMultimap[string, int]()

	mm.PutAll("key1", 1, 2, 3)

	mm.ReplaceValues("key1")

	chetest.RequireEqual(t, mm.ContainsKey("key1"), false)
}

// TestMultimap_String tests with string keys and values
func TestMultimap_String(t *testing.T) {
	mm := NewMultimap[string, string]()

	mm.Put("fruits", "apple")
	mm.Put("fruits", "banana")
	mm.Put("vegetables", "carrot")

	chetest.RequireEqual(t, mm.KeyCount(), 2)
	chetest.RequireEqual(t, mm.ValueCount("fruits"), 2)
	chetest.RequireEqual(t, mm.ValueCount("vegetables"), 1)
}
