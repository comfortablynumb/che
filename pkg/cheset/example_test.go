package cheset_test

import (
	"fmt"

	"github.com/comfortablynumb/che/pkg/cheset"
)

// Example demonstrates basic HashSet operations
func Example() {
	// Create a new set
	set := cheset.New[int]()

	// Add elements
	set.Add(1)
	set.Add(2)
	set.Add(3)

	// Check membership
	fmt.Println("Contains 2:", set.Contains(2))
	fmt.Println("Size:", set.Size())

	// Remove element
	set.Remove(2)
	fmt.Println("Contains 2 after removal:", set.Contains(2))

	// Output:
	// Contains 2: true
	// Size: 3
	// Contains 2 after removal: false
}

// ExampleNew demonstrates creating a new empty HashSet
func ExampleNew() {
	set := cheset.New[string]()
	set.Add("hello")
	set.Add("world")

	fmt.Println("Size:", set.Size())
	fmt.Println("Contains 'hello':", set.Contains("hello"))

	// Output:
	// Size: 2
	// Contains 'hello': true
}

// ExampleNewFromSlice demonstrates creating a HashSet from a slice
func ExampleNewFromSlice() {
	// Create set from slice with duplicates
	slice := []int{1, 2, 3, 2, 1, 4}
	set := cheset.NewFromSlice(slice)

	fmt.Println("Size:", set.Size()) // Duplicates removed

	// Output:
	// Size: 4
}

// ExampleHashSet_Add demonstrates adding elements to a set
func ExampleHashSet_Add() {
	set := cheset.New[int]()

	// First add returns true
	added := set.Add(1)
	fmt.Println("First add:", added)

	// Adding same element returns false
	added = set.Add(1)
	fmt.Println("Duplicate add:", added)

	// Output:
	// First add: true
	// Duplicate add: false
}

// ExampleHashSet_Union demonstrates set union operation
func ExampleHashSet_Union() {
	set1 := cheset.NewFromSlice([]int{1, 2, 3})
	set2 := cheset.NewFromSlice([]int{3, 4, 5})

	union := set1.Union(set2)
	fmt.Println("Union size:", union.Size())
	fmt.Println("Contains 1:", union.Contains(1))
	fmt.Println("Contains 5:", union.Contains(5))

	// Output:
	// Union size: 5
	// Contains 1: true
	// Contains 5: true
}

// ExampleHashSet_Intersect demonstrates set intersection operation
func ExampleHashSet_Intersect() {
	set1 := cheset.NewFromSlice([]int{1, 2, 3, 4})
	set2 := cheset.NewFromSlice([]int{3, 4, 5, 6})

	intersection := set1.Intersect(set2)
	fmt.Println("Intersection size:", intersection.Size())
	fmt.Println("Contains 3:", intersection.Contains(3))
	fmt.Println("Contains 1:", intersection.Contains(1))

	// Output:
	// Intersection size: 2
	// Contains 3: true
	// Contains 1: false
}

// ExampleHashSet_Diff demonstrates set difference operation
func ExampleHashSet_Diff() {
	set1 := cheset.NewFromSlice([]int{1, 2, 3, 4})
	set2 := cheset.NewFromSlice([]int{3, 4, 5, 6})

	diff := set1.Diff(set2)
	fmt.Println("Difference size:", diff.Size())
	fmt.Println("Contains 1:", diff.Contains(1))
	fmt.Println("Contains 3:", diff.Contains(3))

	// Output:
	// Difference size: 2
	// Contains 1: true
	// Contains 3: false
}

// ExampleHashSet_SymmetricDiff demonstrates symmetric difference operation
func ExampleHashSet_SymmetricDiff() {
	set1 := cheset.NewFromSlice([]int{1, 2, 3})
	set2 := cheset.NewFromSlice([]int{3, 4, 5})

	symDiff := set1.SymmetricDiff(set2)
	fmt.Println("Symmetric difference size:", symDiff.Size())
	fmt.Println("Contains 1:", symDiff.Contains(1))
	fmt.Println("Contains 3:", symDiff.Contains(3))
	fmt.Println("Contains 5:", symDiff.Contains(5))

	// Output:
	// Symmetric difference size: 4
	// Contains 1: true
	// Contains 3: false
	// Contains 5: true
}

// ExampleHashSet_IsSubset demonstrates subset checking
func ExampleHashSet_IsSubset() {
	small := cheset.NewFromSlice([]int{1, 2})
	large := cheset.NewFromSlice([]int{1, 2, 3, 4})

	fmt.Println("Small is subset of large:", small.IsSubset(large))
	fmt.Println("Large is subset of small:", large.IsSubset(small))

	// Output:
	// Small is subset of large: true
	// Large is subset of small: false
}

// ExampleHashSet_IsDisjoint demonstrates checking for disjoint sets
func ExampleHashSet_IsDisjoint() {
	set1 := cheset.NewFromSlice([]int{1, 2, 3})
	set2 := cheset.NewFromSlice([]int{4, 5, 6})
	set3 := cheset.NewFromSlice([]int{3, 4, 5})

	fmt.Println("Set1 and Set2 disjoint:", set1.IsDisjoint(set2))
	fmt.Println("Set1 and Set3 disjoint:", set1.IsDisjoint(set3))

	// Output:
	// Set1 and Set2 disjoint: true
	// Set1 and Set3 disjoint: false
}

// ExampleHashSet_Filter demonstrates filtering set elements
func ExampleHashSet_Filter() {
	set := cheset.NewFromSlice([]int{1, 2, 3, 4, 5, 6})

	// Filter even numbers
	evens := set.Filter(func(item int) bool {
		return item%2 == 0
	})

	fmt.Println("Evens size:", evens.Size())
	fmt.Println("Contains 2:", evens.Contains(2))
	fmt.Println("Contains 3:", evens.Contains(3))

	// Output:
	// Evens size: 3
	// Contains 2: true
	// Contains 3: false
}

// ExampleHashSet_ForEach demonstrates iterating over set elements
func ExampleHashSet_ForEach() {
	set := cheset.NewFromSlice([]int{1, 2, 3})

	sum := 0
	set.ForEach(func(item int) bool {
		sum += item
		return true
	})

	fmt.Println("Sum:", sum)

	// Output:
	// Sum: 6
}

// ExampleHashSet_Clone demonstrates cloning a set
func ExampleHashSet_Clone() {
	original := cheset.NewFromSlice([]int{1, 2, 3})
	clone := original.Clone()

	// Modify clone
	clone.Add(4)

	fmt.Println("Original size:", original.Size())
	fmt.Println("Clone size:", clone.Size())

	// Output:
	// Original size: 3
	// Clone size: 4
}

// ExampleHashSet_deduplication demonstrates using HashSet for deduplication
func ExampleHashSet_deduplication() {
	// Input slice with duplicates
	input := []string{"apple", "banana", "apple", "cherry", "banana"}

	// Create set to remove duplicates
	set := cheset.NewFromSlice(input)

	fmt.Println("Original length:", len(input))
	fmt.Println("Unique count:", set.Size())

	// Output:
	// Original length: 5
	// Unique count: 3
}

// ExampleHashSet_commonElements demonstrates finding common elements
func ExampleHashSet_commonElements() {
	users1 := cheset.NewFromSlice([]string{"alice", "bob", "charlie"})
	users2 := cheset.NewFromSlice([]string{"bob", "charlie", "dave"})
	users3 := cheset.NewFromSlice([]string{"charlie", "dave", "eve"})

	// Find users present in all three sets
	common := users1.Intersect(users2).Intersect(users3)

	fmt.Println("Common users:", common.Size())
	fmt.Println("Contains charlie:", common.Contains("charlie"))

	// Output:
	// Common users: 1
	// Contains charlie: true
}

// ExampleHashSet_tags demonstrates using HashSet for tag management
func ExampleHashSet_tags() {
	userTags := cheset.NewFromSlice([]string{"golang", "python", "javascript"})
	requiredTags := cheset.NewFromSlice([]string{"golang", "docker"})

	// Check if user has all required tags
	hasAll := requiredTags.IsSubset(userTags)
	fmt.Println("Has all required tags:", hasAll)

	// Find missing tags
	missing := requiredTags.Diff(userTags)
	fmt.Println("Missing tags count:", missing.Size())

	// Output:
	// Has all required tags: false
	// Missing tags count: 1
}
