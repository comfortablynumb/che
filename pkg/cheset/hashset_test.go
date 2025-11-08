package cheset

import (
	"testing"

	"github.com/comfortablynumb/che/pkg/chetest"
)

// TestNew tests creating a new empty HashSet
func TestNew(t *testing.T) {
	set := New[int]()

	chetest.RequireEqual(t, set.IsEmpty(), true)
	chetest.RequireEqual(t, set.Size(), 0)
}

// TestNewWithCapacity tests creating a HashSet with initial capacity
func TestNewWithCapacity(t *testing.T) {
	set := NewWithCapacity[int](10)

	chetest.RequireEqual(t, set.IsEmpty(), true)
	chetest.RequireEqual(t, set.Size(), 0)
}

// TestNewFromSlice tests creating a HashSet from a slice
func TestNewFromSlice(t *testing.T) {
	slice := []int{1, 2, 3, 2, 1}
	set := NewFromSlice(slice)

	chetest.RequireEqual(t, set.Size(), 3)
	chetest.RequireEqual(t, set.Contains(1), true)
	chetest.RequireEqual(t, set.Contains(2), true)
	chetest.RequireEqual(t, set.Contains(3), true)
}

// TestNewFromSlice_Empty tests creating a HashSet from an empty slice
func TestNewFromSlice_Empty(t *testing.T) {
	slice := []int{}
	set := NewFromSlice(slice)

	chetest.RequireEqual(t, set.IsEmpty(), true)
	chetest.RequireEqual(t, set.Size(), 0)
}

// TestAdd tests adding elements to the set
func TestAdd(t *testing.T) {
	set := New[int]()

	added := set.Add(1)
	chetest.RequireEqual(t, added, true)
	chetest.RequireEqual(t, set.Size(), 1)
	chetest.RequireEqual(t, set.Contains(1), true)

	// Adding the same element should return false
	added = set.Add(1)
	chetest.RequireEqual(t, added, false)
	chetest.RequireEqual(t, set.Size(), 1)
}

// TestAdd_String tests adding string elements
func TestAdd_String(t *testing.T) {
	set := New[string]()

	set.Add("hello")
	set.Add("world")

	chetest.RequireEqual(t, set.Size(), 2)
	chetest.RequireEqual(t, set.Contains("hello"), true)
	chetest.RequireEqual(t, set.Contains("world"), true)
	chetest.RequireEqual(t, set.Contains("foo"), false)
}

// TestAddMultiple tests adding multiple elements at once
func TestAddMultiple(t *testing.T) {
	set := New[int]()

	count := set.AddMultiple(1, 2, 3, 2, 4)
	chetest.RequireEqual(t, count, 4) // 2 is duplicate
	chetest.RequireEqual(t, set.Size(), 4)
}

// TestAddMultiple_Empty tests adding no elements
func TestAddMultiple_Empty(t *testing.T) {
	set := New[int]()

	count := set.AddMultiple()
	chetest.RequireEqual(t, count, 0)
	chetest.RequireEqual(t, set.Size(), 0)
}

// TestRemove tests removing elements from the set
func TestRemove(t *testing.T) {
	set := NewFromSlice([]int{1, 2, 3})

	removed := set.Remove(2)
	chetest.RequireEqual(t, removed, true)
	chetest.RequireEqual(t, set.Size(), 2)
	chetest.RequireEqual(t, set.Contains(2), false)

	// Removing non-existent element should return false
	removed = set.Remove(5)
	chetest.RequireEqual(t, removed, false)
	chetest.RequireEqual(t, set.Size(), 2)
}

// TestRemoveMultiple tests removing multiple elements
func TestRemoveMultiple(t *testing.T) {
	set := NewFromSlice([]int{1, 2, 3, 4, 5})

	count := set.RemoveMultiple(2, 4, 6)
	chetest.RequireEqual(t, count, 2) // 6 doesn't exist
	chetest.RequireEqual(t, set.Size(), 3)
}

// TestRemoveMultiple_Empty tests removing no elements
func TestRemoveMultiple_Empty(t *testing.T) {
	set := NewFromSlice([]int{1, 2, 3})

	count := set.RemoveMultiple()
	chetest.RequireEqual(t, count, 0)
	chetest.RequireEqual(t, set.Size(), 3)
}

// TestContains tests checking element existence
func TestContains(t *testing.T) {
	set := NewFromSlice([]int{1, 2, 3})

	chetest.RequireEqual(t, set.Contains(1), true)
	chetest.RequireEqual(t, set.Contains(2), true)
	chetest.RequireEqual(t, set.Contains(3), true)
	chetest.RequireEqual(t, set.Contains(4), false)
}

// TestContainsAll tests checking if all elements exist
func TestContainsAll(t *testing.T) {
	set := NewFromSlice([]int{1, 2, 3, 4, 5})

	chetest.RequireEqual(t, set.ContainsAll(1, 2, 3), true)
	chetest.RequireEqual(t, set.ContainsAll(1, 6), false)
	chetest.RequireEqual(t, set.ContainsAll(), true) // Empty check
}

// TestContainsAny tests checking if any element exists
func TestContainsAny(t *testing.T) {
	set := NewFromSlice([]int{1, 2, 3})

	chetest.RequireEqual(t, set.ContainsAny(1, 5, 6), true)
	chetest.RequireEqual(t, set.ContainsAny(7, 8, 9), false)
	chetest.RequireEqual(t, set.ContainsAny(), false) // Empty check
}

// TestSize tests getting the set size
func TestSize(t *testing.T) {
	set := New[int]()
	chetest.RequireEqual(t, set.Size(), 0)

	set.Add(1)
	chetest.RequireEqual(t, set.Size(), 1)

	set.Add(2)
	chetest.RequireEqual(t, set.Size(), 2)

	set.Remove(1)
	chetest.RequireEqual(t, set.Size(), 1)
}

// TestIsEmpty tests checking if set is empty
func TestIsEmpty(t *testing.T) {
	set := New[int]()
	chetest.RequireEqual(t, set.IsEmpty(), true)

	set.Add(1)
	chetest.RequireEqual(t, set.IsEmpty(), false)

	set.Remove(1)
	chetest.RequireEqual(t, set.IsEmpty(), true)
}

// TestClear tests clearing all elements
func TestClear(t *testing.T) {
	set := NewFromSlice([]int{1, 2, 3, 4, 5})
	chetest.RequireEqual(t, set.Size(), 5)

	set.Clear()
	chetest.RequireEqual(t, set.Size(), 0)
	chetest.RequireEqual(t, set.IsEmpty(), true)
	chetest.RequireEqual(t, set.Contains(1), false)
}

// TestToSlice tests converting set to slice
func TestToSlice(t *testing.T) {
	set := NewFromSlice([]int{3, 1, 2})
	slice := set.ToSlice()

	chetest.RequireEqual(t, len(slice), 3)

	// Check all elements are present (order doesn't matter)
	sliceSet := NewFromSlice(slice)
	chetest.RequireEqual(t, sliceSet.Contains(1), true)
	chetest.RequireEqual(t, sliceSet.Contains(2), true)
	chetest.RequireEqual(t, sliceSet.Contains(3), true)
}

// TestToSlice_Empty tests converting empty set to slice
func TestToSlice_Empty(t *testing.T) {
	set := New[int]()
	slice := set.ToSlice()

	chetest.RequireEqual(t, len(slice), 0)
}

// TestClone tests cloning a set
func TestClone(t *testing.T) {
	original := NewFromSlice([]int{1, 2, 3})
	clone := original.Clone()

	chetest.RequireEqual(t, clone.Size(), 3)
	chetest.RequireEqual(t, clone.Contains(1), true)
	chetest.RequireEqual(t, clone.Contains(2), true)
	chetest.RequireEqual(t, clone.Contains(3), true)

	// Modifying clone shouldn't affect original
	clone.Add(4)
	chetest.RequireEqual(t, original.Contains(4), false)
	chetest.RequireEqual(t, clone.Contains(4), true)
}

// TestClone_Empty tests cloning an empty set
func TestClone_Empty(t *testing.T) {
	original := New[int]()
	clone := original.Clone()

	chetest.RequireEqual(t, clone.IsEmpty(), true)
}

// TestEqual tests checking set equality
func TestEqual(t *testing.T) {
	set1 := NewFromSlice([]int{1, 2, 3})
	set2 := NewFromSlice([]int{3, 2, 1})
	set3 := NewFromSlice([]int{1, 2, 3, 4})

	chetest.RequireEqual(t, set1.Equal(set2), true)
	chetest.RequireEqual(t, set1.Equal(set3), false)
}

// TestEqual_Empty tests equality of empty sets
func TestEqual_Empty(t *testing.T) {
	set1 := New[int]()
	set2 := New[int]()

	chetest.RequireEqual(t, set1.Equal(set2), true)
}

// TestEqual_DifferentSizes tests equality with different sizes
func TestEqual_DifferentSizes(t *testing.T) {
	set1 := NewFromSlice([]int{1, 2})
	set2 := NewFromSlice([]int{1, 2, 3})

	chetest.RequireEqual(t, set1.Equal(set2), false)
}

// TestEqual_SameSizeDifferentElements tests equality with same size but different elements
func TestEqual_SameSizeDifferentElements(t *testing.T) {
	set1 := NewFromSlice([]int{1, 2, 3})
	set2 := NewFromSlice([]int{1, 2, 4})

	chetest.RequireEqual(t, set1.Equal(set2), false)
}

// TestUnion tests set union operation
func TestUnion(t *testing.T) {
	set1 := NewFromSlice([]int{1, 2, 3})
	set2 := NewFromSlice([]int{3, 4, 5})

	result := set1.Union(set2)

	chetest.RequireEqual(t, result.Size(), 5)
	chetest.RequireEqual(t, result.Contains(1), true)
	chetest.RequireEqual(t, result.Contains(2), true)
	chetest.RequireEqual(t, result.Contains(3), true)
	chetest.RequireEqual(t, result.Contains(4), true)
	chetest.RequireEqual(t, result.Contains(5), true)

	// Original sets should be unchanged
	chetest.RequireEqual(t, set1.Size(), 3)
	chetest.RequireEqual(t, set2.Size(), 3)
}

// TestUnion_Empty tests union with empty sets
func TestUnion_Empty(t *testing.T) {
	set1 := NewFromSlice([]int{1, 2})
	set2 := New[int]()

	result := set1.Union(set2)
	chetest.RequireEqual(t, result.Size(), 2)
}

// TestIntersect tests set intersection operation
func TestIntersect(t *testing.T) {
	set1 := NewFromSlice([]int{1, 2, 3, 4})
	set2 := NewFromSlice([]int{3, 4, 5, 6})

	result := set1.Intersect(set2)

	chetest.RequireEqual(t, result.Size(), 2)
	chetest.RequireEqual(t, result.Contains(3), true)
	chetest.RequireEqual(t, result.Contains(4), true)
	chetest.RequireEqual(t, result.Contains(1), false)
	chetest.RequireEqual(t, result.Contains(5), false)
}

// TestIntersect_Empty tests intersection with empty sets
func TestIntersect_Empty(t *testing.T) {
	set1 := NewFromSlice([]int{1, 2, 3})
	set2 := New[int]()

	result := set1.Intersect(set2)
	chetest.RequireEqual(t, result.IsEmpty(), true)
}

// TestIntersect_NoCommon tests intersection with no common elements
func TestIntersect_NoCommon(t *testing.T) {
	set1 := NewFromSlice([]int{1, 2, 3})
	set2 := NewFromSlice([]int{4, 5, 6})

	result := set1.Intersect(set2)
	chetest.RequireEqual(t, result.IsEmpty(), true)
}

// TestDiff tests set difference operation
func TestDiff(t *testing.T) {
	set1 := NewFromSlice([]int{1, 2, 3, 4})
	set2 := NewFromSlice([]int{3, 4, 5, 6})

	result := set1.Diff(set2)

	chetest.RequireEqual(t, result.Size(), 2)
	chetest.RequireEqual(t, result.Contains(1), true)
	chetest.RequireEqual(t, result.Contains(2), true)
	chetest.RequireEqual(t, result.Contains(3), false)
	chetest.RequireEqual(t, result.Contains(4), false)
}

// TestDiff_Empty tests difference with empty sets
func TestDiff_Empty(t *testing.T) {
	set1 := NewFromSlice([]int{1, 2, 3})
	set2 := New[int]()

	result := set1.Diff(set2)
	chetest.RequireEqual(t, result.Size(), 3)
}

// TestSymmetricDiff tests symmetric difference operation
func TestSymmetricDiff(t *testing.T) {
	set1 := NewFromSlice([]int{1, 2, 3, 4})
	set2 := NewFromSlice([]int{3, 4, 5, 6})

	result := set1.SymmetricDiff(set2)

	chetest.RequireEqual(t, result.Size(), 4)
	chetest.RequireEqual(t, result.Contains(1), true)
	chetest.RequireEqual(t, result.Contains(2), true)
	chetest.RequireEqual(t, result.Contains(5), true)
	chetest.RequireEqual(t, result.Contains(6), true)
	chetest.RequireEqual(t, result.Contains(3), false)
	chetest.RequireEqual(t, result.Contains(4), false)
}

// TestSymmetricDiff_Empty tests symmetric difference with empty sets
func TestSymmetricDiff_Empty(t *testing.T) {
	set1 := NewFromSlice([]int{1, 2})
	set2 := New[int]()

	result := set1.SymmetricDiff(set2)
	chetest.RequireEqual(t, result.Size(), 2)
}

// TestIsSubset tests subset checking
func TestIsSubset(t *testing.T) {
	set1 := NewFromSlice([]int{1, 2})
	set2 := NewFromSlice([]int{1, 2, 3, 4})

	chetest.RequireEqual(t, set1.IsSubset(set2), true)
	chetest.RequireEqual(t, set2.IsSubset(set1), false)
}

// TestIsSubset_Equal tests subset with equal sets
func TestIsSubset_Equal(t *testing.T) {
	set1 := NewFromSlice([]int{1, 2, 3})
	set2 := NewFromSlice([]int{1, 2, 3})

	chetest.RequireEqual(t, set1.IsSubset(set2), true)
}

// TestIsSubset_Empty tests subset with empty set
func TestIsSubset_Empty(t *testing.T) {
	set1 := New[int]()
	set2 := NewFromSlice([]int{1, 2, 3})

	chetest.RequireEqual(t, set1.IsSubset(set2), true)
}

// TestIsSubset_NotSubset tests when set is not a subset
func TestIsSubset_NotSubset(t *testing.T) {
	set1 := NewFromSlice([]int{1, 2, 5})
	set2 := NewFromSlice([]int{1, 2, 3})

	chetest.RequireEqual(t, set1.IsSubset(set2), false)
}

// TestIsSuperset tests superset checking
func TestIsSuperset(t *testing.T) {
	set1 := NewFromSlice([]int{1, 2, 3, 4})
	set2 := NewFromSlice([]int{1, 2})

	chetest.RequireEqual(t, set1.IsSuperset(set2), true)
	chetest.RequireEqual(t, set2.IsSuperset(set1), false)
}

// TestIsProperSubset tests proper subset checking
func TestIsProperSubset(t *testing.T) {
	set1 := NewFromSlice([]int{1, 2})
	set2 := NewFromSlice([]int{1, 2, 3})
	set3 := NewFromSlice([]int{1, 2})

	chetest.RequireEqual(t, set1.IsProperSubset(set2), true)
	chetest.RequireEqual(t, set1.IsProperSubset(set3), false) // Equal sets
}

// TestIsProperSuperset tests proper superset checking
func TestIsProperSuperset(t *testing.T) {
	set1 := NewFromSlice([]int{1, 2, 3})
	set2 := NewFromSlice([]int{1, 2})
	set3 := NewFromSlice([]int{1, 2, 3})

	chetest.RequireEqual(t, set1.IsProperSuperset(set2), true)
	chetest.RequireEqual(t, set1.IsProperSuperset(set3), false) // Equal sets
}

// TestIsDisjoint tests checking if sets have no common elements
func TestIsDisjoint(t *testing.T) {
	set1 := NewFromSlice([]int{1, 2, 3})
	set2 := NewFromSlice([]int{4, 5, 6})
	set3 := NewFromSlice([]int{3, 4, 5})

	chetest.RequireEqual(t, set1.IsDisjoint(set2), true)
	chetest.RequireEqual(t, set1.IsDisjoint(set3), false)
}

// TestIsDisjoint_Empty tests disjoint with empty sets
func TestIsDisjoint_Empty(t *testing.T) {
	set1 := NewFromSlice([]int{1, 2, 3})
	set2 := New[int]()

	chetest.RequireEqual(t, set1.IsDisjoint(set2), true)
}

// TestForEach tests iterating over set elements
func TestForEach(t *testing.T) {
	set := NewFromSlice([]int{1, 2, 3, 4, 5})

	sum := 0
	set.ForEach(func(item int) bool {
		sum += item
		return true
	})

	chetest.RequireEqual(t, sum, 15)
}

// TestForEach_EarlyExit tests early exit from iteration
func TestForEach_EarlyExit(t *testing.T) {
	set := NewFromSlice([]int{1, 2, 3, 4, 5})

	count := 0
	set.ForEach(func(item int) bool {
		count++
		return count < 3 // Stop after 3 iterations
	})

	chetest.RequireEqual(t, count, 3)
}

// TestFilter tests filtering set elements
func TestFilter(t *testing.T) {
	set := NewFromSlice([]int{1, 2, 3, 4, 5, 6})

	// Filter even numbers
	evens := set.Filter(func(item int) bool {
		return item%2 == 0
	})

	chetest.RequireEqual(t, evens.Size(), 3)
	chetest.RequireEqual(t, evens.Contains(2), true)
	chetest.RequireEqual(t, evens.Contains(4), true)
	chetest.RequireEqual(t, evens.Contains(6), true)

	// Original set should be unchanged
	chetest.RequireEqual(t, set.Size(), 6)
}

// TestFilter_NoMatches tests filter with no matching elements
func TestFilter_NoMatches(t *testing.T) {
	set := NewFromSlice([]int{1, 3, 5})

	evens := set.Filter(func(item int) bool {
		return item%2 == 0
	})

	chetest.RequireEqual(t, evens.IsEmpty(), true)
}

// TestString tests string representation of set
func TestString(t *testing.T) {
	set := New[int]()
	str := set.String()

	chetest.RequireEqual(t, str, "HashSet{}")

	set.Add(1)
	str = set.String()

	// Check that it contains the element (format may vary)
	chetest.RequireEqual(t, str != "", true)
	chetest.RequireEqual(t, str != "HashSet{}", true)
}

// TestString_MultipleElements tests string representation with multiple elements
func TestString_MultipleElements(t *testing.T) {
	set := NewFromSlice([]string{"a", "b", "c"})
	str := set.String()

	// Just check it's not empty and starts/ends correctly
	chetest.RequireEqual(t, len(str) > len("HashSet{}"), true)
}
