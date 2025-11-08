package cheset

import (
	"testing"

	"github.com/comfortablynumb/che/pkg/chetest"
)

// TestNewOrdered tests creating a new empty OrderedSet
func TestNewOrdered(t *testing.T) {
	set := NewOrdered[int]()

	chetest.RequireEqual(t, set.IsEmpty(), true)
	chetest.RequireEqual(t, set.Size(), 0)
}

// TestNewOrderedWithCapacity tests creating an OrderedSet with initial capacity
func TestNewOrderedWithCapacity(t *testing.T) {
	set := NewOrderedWithCapacity[int](10)

	chetest.RequireEqual(t, set.IsEmpty(), true)
	chetest.RequireEqual(t, set.Size(), 0)
}

// TestNewOrderedFromSlice tests creating an OrderedSet from a slice
func TestNewOrderedFromSlice(t *testing.T) {
	slice := []int{1, 2, 3, 2, 1}
	set := NewOrderedFromSlice(slice)

	chetest.RequireEqual(t, set.Size(), 3)
	chetest.RequireEqual(t, set.GetAt(0), 1)
	chetest.RequireEqual(t, set.GetAt(1), 2)
	chetest.RequireEqual(t, set.GetAt(2), 3)
}

// TestNewOrderedFromSlice_Empty tests creating an OrderedSet from an empty slice
func TestNewOrderedFromSlice_Empty(t *testing.T) {
	slice := []int{}
	set := NewOrderedFromSlice(slice)

	chetest.RequireEqual(t, set.IsEmpty(), true)
}

// TestOrderedSet_Add tests adding elements preserves order
func TestOrderedSet_Add(t *testing.T) {
	set := NewOrdered[int]()

	added := set.Add(1)
	chetest.RequireEqual(t, added, true)
	chetest.RequireEqual(t, set.Size(), 1)

	added = set.Add(2)
	chetest.RequireEqual(t, added, true)
	chetest.RequireEqual(t, set.Size(), 2)

	// Adding same element returns false
	added = set.Add(1)
	chetest.RequireEqual(t, added, false)
	chetest.RequireEqual(t, set.Size(), 2)

	// Check order
	chetest.RequireEqual(t, set.GetAt(0), 1)
	chetest.RequireEqual(t, set.GetAt(1), 2)
}

// TestOrderedSet_AddMultiple tests adding multiple elements
func TestOrderedSet_AddMultiple(t *testing.T) {
	set := NewOrdered[int]()

	count := set.AddMultiple(1, 2, 3, 2, 4)
	chetest.RequireEqual(t, count, 4) // 2 is duplicate
	chetest.RequireEqual(t, set.Size(), 4)

	// Check order
	chetest.RequireEqual(t, set.GetAt(0), 1)
	chetest.RequireEqual(t, set.GetAt(1), 2)
	chetest.RequireEqual(t, set.GetAt(2), 3)
	chetest.RequireEqual(t, set.GetAt(3), 4)
}

// TestOrderedSet_AddMultiple_Empty tests adding no elements
func TestOrderedSet_AddMultiple_Empty(t *testing.T) {
	set := NewOrdered[int]()

	count := set.AddMultiple()
	chetest.RequireEqual(t, count, 0)
	chetest.RequireEqual(t, set.Size(), 0)
}

// TestOrderedSet_Remove tests removing elements
func TestOrderedSet_Remove(t *testing.T) {
	set := NewOrderedFromSlice([]int{1, 2, 3, 4, 5})

	removed := set.Remove(3)
	chetest.RequireEqual(t, removed, true)
	chetest.RequireEqual(t, set.Size(), 4)
	chetest.RequireEqual(t, set.Contains(3), false)

	// Check order is maintained
	chetest.RequireEqual(t, set.GetAt(0), 1)
	chetest.RequireEqual(t, set.GetAt(1), 2)
	chetest.RequireEqual(t, set.GetAt(2), 4)
	chetest.RequireEqual(t, set.GetAt(3), 5)

	// Removing non-existent element returns false
	removed = set.Remove(10)
	chetest.RequireEqual(t, removed, false)
}

// TestOrderedSet_Remove_First tests removing first element
func TestOrderedSet_Remove_First(t *testing.T) {
	set := NewOrderedFromSlice([]int{1, 2, 3})

	set.Remove(1)
	chetest.RequireEqual(t, set.Size(), 2)
	chetest.RequireEqual(t, set.GetAt(0), 2)
	chetest.RequireEqual(t, set.GetAt(1), 3)
}

// TestOrderedSet_Remove_Last tests removing last element
func TestOrderedSet_Remove_Last(t *testing.T) {
	set := NewOrderedFromSlice([]int{1, 2, 3})

	set.Remove(3)
	chetest.RequireEqual(t, set.Size(), 2)
	chetest.RequireEqual(t, set.GetAt(0), 1)
	chetest.RequireEqual(t, set.GetAt(1), 2)
}

// TestOrderedSet_RemoveMultiple tests removing multiple elements
func TestOrderedSet_RemoveMultiple(t *testing.T) {
	set := NewOrderedFromSlice([]int{1, 2, 3, 4, 5})

	count := set.RemoveMultiple(2, 4, 6)
	chetest.RequireEqual(t, count, 2) // 6 doesn't exist
	chetest.RequireEqual(t, set.Size(), 3)

	// Check order
	chetest.RequireEqual(t, set.GetAt(0), 1)
	chetest.RequireEqual(t, set.GetAt(1), 3)
	chetest.RequireEqual(t, set.GetAt(2), 5)
}

// TestOrderedSet_Contains tests checking element existence
func TestOrderedSet_Contains(t *testing.T) {
	set := NewOrderedFromSlice([]int{1, 2, 3})

	chetest.RequireEqual(t, set.Contains(1), true)
	chetest.RequireEqual(t, set.Contains(2), true)
	chetest.RequireEqual(t, set.Contains(4), false)
}

// TestOrderedSet_ContainsAll tests checking if all elements exist
func TestOrderedSet_ContainsAll(t *testing.T) {
	set := NewOrderedFromSlice([]int{1, 2, 3, 4, 5})

	chetest.RequireEqual(t, set.ContainsAll(1, 2, 3), true)
	chetest.RequireEqual(t, set.ContainsAll(1, 6), false)
	chetest.RequireEqual(t, set.ContainsAll(), true)
}

// TestOrderedSet_ContainsAny tests checking if any element exists
func TestOrderedSet_ContainsAny(t *testing.T) {
	set := NewOrderedFromSlice([]int{1, 2, 3})

	chetest.RequireEqual(t, set.ContainsAny(1, 5), true)
	chetest.RequireEqual(t, set.ContainsAny(7, 8), false)
	chetest.RequireEqual(t, set.ContainsAny(), false)
}

// TestOrderedSet_Size tests getting the set size
func TestOrderedSet_Size(t *testing.T) {
	set := NewOrdered[int]()
	chetest.RequireEqual(t, set.Size(), 0)

	set.Add(1)
	chetest.RequireEqual(t, set.Size(), 1)

	set.Add(2)
	chetest.RequireEqual(t, set.Size(), 2)
}

// TestOrderedSet_IsEmpty tests checking if set is empty
func TestOrderedSet_IsEmpty(t *testing.T) {
	set := NewOrdered[int]()
	chetest.RequireEqual(t, set.IsEmpty(), true)

	set.Add(1)
	chetest.RequireEqual(t, set.IsEmpty(), false)
}

// TestOrderedSet_Clear tests clearing all elements
func TestOrderedSet_Clear(t *testing.T) {
	set := NewOrderedFromSlice([]int{1, 2, 3, 4, 5})
	chetest.RequireEqual(t, set.Size(), 5)

	set.Clear()
	chetest.RequireEqual(t, set.Size(), 0)
	chetest.RequireEqual(t, set.IsEmpty(), true)
}

// TestOrderedSet_ToSlice tests converting set to slice
func TestOrderedSet_ToSlice(t *testing.T) {
	set := NewOrderedFromSlice([]int{3, 1, 2})
	slice := set.ToSlice()

	chetest.RequireEqual(t, len(slice), 3)
	chetest.RequireEqual(t, slice[0], 3) // Order preserved
	chetest.RequireEqual(t, slice[1], 1)
	chetest.RequireEqual(t, slice[2], 2)
}

// TestOrderedSet_ToSlice_Empty tests converting empty set to slice
func TestOrderedSet_ToSlice_Empty(t *testing.T) {
	set := NewOrdered[int]()
	slice := set.ToSlice()

	chetest.RequireEqual(t, len(slice), 0)
}

// TestOrderedSet_Clone tests cloning a set
func TestOrderedSet_Clone(t *testing.T) {
	original := NewOrderedFromSlice([]int{1, 2, 3})
	clone := original.Clone()

	chetest.RequireEqual(t, clone.Size(), 3)
	chetest.RequireEqual(t, clone.GetAt(0), 1)
	chetest.RequireEqual(t, clone.GetAt(1), 2)
	chetest.RequireEqual(t, clone.GetAt(2), 3)

	// Modifying clone shouldn't affect original
	clone.Add(4)
	chetest.RequireEqual(t, original.Contains(4), false)
	chetest.RequireEqual(t, clone.Contains(4), true)
}

// TestOrderedSet_Equal tests checking set equality with order
func TestOrderedSet_Equal(t *testing.T) {
	set1 := NewOrderedFromSlice([]int{1, 2, 3})
	set2 := NewOrderedFromSlice([]int{1, 2, 3})
	set3 := NewOrderedFromSlice([]int{3, 2, 1}) // Different order

	chetest.RequireEqual(t, set1.Equal(set2), true)
	chetest.RequireEqual(t, set1.Equal(set3), false) // Order matters
}

// TestOrderedSet_Equal_DifferentSizes tests equality with different sizes
func TestOrderedSet_Equal_DifferentSizes(t *testing.T) {
	set1 := NewOrderedFromSlice([]int{1, 2})
	set2 := NewOrderedFromSlice([]int{1, 2, 3})

	chetest.RequireEqual(t, set1.Equal(set2), false)
}

// TestOrderedSet_Union tests set union operation
func TestOrderedSet_Union(t *testing.T) {
	set1 := NewOrderedFromSlice([]int{1, 2, 3})
	set2 := NewOrderedFromSlice([]int{3, 4, 5})

	result := set1.Union(set2)

	chetest.RequireEqual(t, result.Size(), 5)
	// Check order: elements from set1 first, then new elements from set2
	chetest.RequireEqual(t, result.GetAt(0), 1)
	chetest.RequireEqual(t, result.GetAt(1), 2)
	chetest.RequireEqual(t, result.GetAt(2), 3)
	chetest.RequireEqual(t, result.GetAt(3), 4)
	chetest.RequireEqual(t, result.GetAt(4), 5)
}

// TestOrderedSet_Intersect tests set intersection operation
func TestOrderedSet_Intersect(t *testing.T) {
	set1 := NewOrderedFromSlice([]int{1, 2, 3, 4})
	set2 := NewOrderedFromSlice([]int{3, 4, 5, 6})

	result := set1.Intersect(set2)

	chetest.RequireEqual(t, result.Size(), 2)
	// Order from set1
	chetest.RequireEqual(t, result.GetAt(0), 3)
	chetest.RequireEqual(t, result.GetAt(1), 4)
}

// TestOrderedSet_Diff tests set difference operation
func TestOrderedSet_Diff(t *testing.T) {
	set1 := NewOrderedFromSlice([]int{1, 2, 3, 4})
	set2 := NewOrderedFromSlice([]int{3, 4, 5, 6})

	result := set1.Diff(set2)

	chetest.RequireEqual(t, result.Size(), 2)
	chetest.RequireEqual(t, result.GetAt(0), 1)
	chetest.RequireEqual(t, result.GetAt(1), 2)
}

// TestOrderedSet_SymmetricDiff tests symmetric difference operation
func TestOrderedSet_SymmetricDiff(t *testing.T) {
	set1 := NewOrderedFromSlice([]int{1, 2, 3})
	set2 := NewOrderedFromSlice([]int{3, 4, 5})

	result := set1.SymmetricDiff(set2)

	chetest.RequireEqual(t, result.Size(), 4)
	// Elements from set1 first, then from set2
	chetest.RequireEqual(t, result.GetAt(0), 1)
	chetest.RequireEqual(t, result.GetAt(1), 2)
	chetest.RequireEqual(t, result.GetAt(2), 4)
	chetest.RequireEqual(t, result.GetAt(3), 5)
}

// TestOrderedSet_IsSubset tests subset checking
func TestOrderedSet_IsSubset(t *testing.T) {
	set1 := NewOrderedFromSlice([]int{1, 2})
	set2 := NewOrderedFromSlice([]int{1, 2, 3, 4})

	chetest.RequireEqual(t, set1.IsSubset(set2), true)
	chetest.RequireEqual(t, set2.IsSubset(set1), false)
}

// TestOrderedSet_IsSubset_Empty tests subset with empty set
func TestOrderedSet_IsSubset_Empty(t *testing.T) {
	set1 := NewOrdered[int]()
	set2 := NewOrderedFromSlice([]int{1, 2, 3})

	chetest.RequireEqual(t, set1.IsSubset(set2), true)
}

// TestOrderedSet_IsSubset_NotSubset tests when set is not a subset
func TestOrderedSet_IsSubset_NotSubset(t *testing.T) {
	set1 := NewOrderedFromSlice([]int{1, 2, 5})
	set2 := NewOrderedFromSlice([]int{1, 2, 3})

	chetest.RequireEqual(t, set1.IsSubset(set2), false)
}

// TestOrderedSet_IsSuperset tests superset checking
func TestOrderedSet_IsSuperset(t *testing.T) {
	set1 := NewOrderedFromSlice([]int{1, 2, 3, 4})
	set2 := NewOrderedFromSlice([]int{1, 2})

	chetest.RequireEqual(t, set1.IsSuperset(set2), true)
	chetest.RequireEqual(t, set2.IsSuperset(set1), false)
}

// TestOrderedSet_IsProperSubset tests proper subset checking
func TestOrderedSet_IsProperSubset(t *testing.T) {
	set1 := NewOrderedFromSlice([]int{1, 2})
	set2 := NewOrderedFromSlice([]int{1, 2, 3})
	set3 := NewOrderedFromSlice([]int{1, 2})

	chetest.RequireEqual(t, set1.IsProperSubset(set2), true)
	chetest.RequireEqual(t, set1.IsProperSubset(set3), false)
}

// TestOrderedSet_IsProperSuperset tests proper superset checking
func TestOrderedSet_IsProperSuperset(t *testing.T) {
	set1 := NewOrderedFromSlice([]int{1, 2, 3})
	set2 := NewOrderedFromSlice([]int{1, 2})
	set3 := NewOrderedFromSlice([]int{1, 2, 3})

	chetest.RequireEqual(t, set1.IsProperSuperset(set2), true)
	chetest.RequireEqual(t, set1.IsProperSuperset(set3), false)
}

// TestOrderedSet_IsDisjoint tests checking if sets have no common elements
func TestOrderedSet_IsDisjoint(t *testing.T) {
	set1 := NewOrderedFromSlice([]int{1, 2, 3})
	set2 := NewOrderedFromSlice([]int{4, 5, 6})
	set3 := NewOrderedFromSlice([]int{3, 4, 5})

	chetest.RequireEqual(t, set1.IsDisjoint(set2), true)
	chetest.RequireEqual(t, set1.IsDisjoint(set3), false)
}

// TestOrderedSet_IsDisjoint_Empty tests disjoint with empty sets
func TestOrderedSet_IsDisjoint_Empty(t *testing.T) {
	set1 := NewOrderedFromSlice([]int{1, 2, 3})
	set2 := NewOrdered[int]()

	chetest.RequireEqual(t, set1.IsDisjoint(set2), true)
}

// TestOrderedSet_ForEach tests iterating over set elements
func TestOrderedSet_ForEach(t *testing.T) {
	set := NewOrderedFromSlice([]int{1, 2, 3})

	var collected []int
	set.ForEach(func(item int) bool {
		collected = append(collected, item)
		return true
	})

	chetest.RequireEqual(t, len(collected), 3)
	chetest.RequireEqual(t, collected[0], 1)
	chetest.RequireEqual(t, collected[1], 2)
	chetest.RequireEqual(t, collected[2], 3)
}

// TestOrderedSet_ForEach_EarlyExit tests early exit from iteration
func TestOrderedSet_ForEach_EarlyExit(t *testing.T) {
	set := NewOrderedFromSlice([]int{1, 2, 3, 4, 5})

	count := 0
	set.ForEach(func(item int) bool {
		count++
		return count < 3
	})

	chetest.RequireEqual(t, count, 3)
}

// TestOrderedSet_Filter tests filtering set elements
func TestOrderedSet_Filter(t *testing.T) {
	set := NewOrderedFromSlice([]int{1, 2, 3, 4, 5, 6})

	// Filter even numbers
	evens := set.Filter(func(item int) bool {
		return item%2 == 0
	})

	chetest.RequireEqual(t, evens.Size(), 3)
	chetest.RequireEqual(t, evens.GetAt(0), 2)
	chetest.RequireEqual(t, evens.GetAt(1), 4)
	chetest.RequireEqual(t, evens.GetAt(2), 6)
}

// TestOrderedSet_GetAt tests getting element by index
func TestOrderedSet_GetAt(t *testing.T) {
	set := NewOrderedFromSlice([]int{10, 20, 30})

	chetest.RequireEqual(t, set.GetAt(0), 10)
	chetest.RequireEqual(t, set.GetAt(1), 20)
	chetest.RequireEqual(t, set.GetAt(2), 30)
}

// TestOrderedSet_GetAt_Panic tests GetAt panics on invalid index
func TestOrderedSet_GetAt_Panic(t *testing.T) {
	set := NewOrderedFromSlice([]int{1, 2, 3})

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("GetAt should panic on out of bounds index")
		}
	}()

	set.GetAt(10) // Should panic
}

// TestOrderedSet_GetAt_NegativePanic tests GetAt panics on negative index
func TestOrderedSet_GetAt_NegativePanic(t *testing.T) {
	set := NewOrderedFromSlice([]int{1, 2, 3})

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("GetAt should panic on negative index")
		}
	}()

	set.GetAt(-1) // Should panic
}

// TestOrderedSet_Index tests getting index of element
func TestOrderedSet_Index(t *testing.T) {
	set := NewOrderedFromSlice([]int{10, 20, 30})

	chetest.RequireEqual(t, set.Index(10), 0)
	chetest.RequireEqual(t, set.Index(20), 1)
	chetest.RequireEqual(t, set.Index(30), 2)
	chetest.RequireEqual(t, set.Index(40), -1) // Not found
}

// TestOrderedSet_First tests getting first element
func TestOrderedSet_First(t *testing.T) {
	set := NewOrderedFromSlice([]int{1, 2, 3})

	first, ok := set.First()
	chetest.RequireEqual(t, ok, true)
	chetest.RequireEqual(t, first, 1)
}

// TestOrderedSet_First_Empty tests First on empty set
func TestOrderedSet_First_Empty(t *testing.T) {
	set := NewOrdered[int]()

	_, ok := set.First()
	chetest.RequireEqual(t, ok, false)
}

// TestOrderedSet_Last tests getting last element
func TestOrderedSet_Last(t *testing.T) {
	set := NewOrderedFromSlice([]int{1, 2, 3})

	last, ok := set.Last()
	chetest.RequireEqual(t, ok, true)
	chetest.RequireEqual(t, last, 3)
}

// TestOrderedSet_Last_Empty tests Last on empty set
func TestOrderedSet_Last_Empty(t *testing.T) {
	set := NewOrdered[int]()

	_, ok := set.Last()
	chetest.RequireEqual(t, ok, false)
}

// TestOrderedSet_PopFirst tests removing and returning first element
func TestOrderedSet_PopFirst(t *testing.T) {
	set := NewOrderedFromSlice([]int{1, 2, 3})

	first, ok := set.PopFirst()
	chetest.RequireEqual(t, ok, true)
	chetest.RequireEqual(t, first, 1)
	chetest.RequireEqual(t, set.Size(), 2)
	chetest.RequireEqual(t, set.GetAt(0), 2)
}

// TestOrderedSet_PopFirst_Empty tests PopFirst on empty set
func TestOrderedSet_PopFirst_Empty(t *testing.T) {
	set := NewOrdered[int]()

	_, ok := set.PopFirst()
	chetest.RequireEqual(t, ok, false)
}

// TestOrderedSet_PopLast tests removing and returning last element
func TestOrderedSet_PopLast(t *testing.T) {
	set := NewOrderedFromSlice([]int{1, 2, 3})

	last, ok := set.PopLast()
	chetest.RequireEqual(t, ok, true)
	chetest.RequireEqual(t, last, 3)
	chetest.RequireEqual(t, set.Size(), 2)
	chetest.RequireEqual(t, set.GetAt(1), 2)
}

// TestOrderedSet_PopLast_Empty tests PopLast on empty set
func TestOrderedSet_PopLast_Empty(t *testing.T) {
	set := NewOrdered[int]()

	_, ok := set.PopLast()
	chetest.RequireEqual(t, ok, false)
}

// TestOrderedSet_String tests string representation
func TestOrderedSet_String(t *testing.T) {
	set := NewOrdered[int]()
	str := set.String()

	chetest.RequireEqual(t, str, "OrderedSet[]")

	set.Add(1)
	set.Add(2)
	str = set.String()

	chetest.RequireEqual(t, str, "OrderedSet[1, 2]")
}

// TestOrderedSet_String_Strings tests string representation with strings
func TestOrderedSet_String_Strings(t *testing.T) {
	set := NewOrderedFromSlice([]string{"a", "b", "c"})
	str := set.String()

	chetest.RequireEqual(t, str, "OrderedSet[a, b, c]")
}
