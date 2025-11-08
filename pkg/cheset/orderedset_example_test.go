package cheset_test

import (
	"fmt"

	"github.com/comfortablynumb/che/pkg/cheset"
)

// ExampleNewOrdered demonstrates creating a new OrderedSet
func ExampleNewOrdered() {
	set := cheset.NewOrdered[string]()
	set.Add("first")
	set.Add("second")
	set.Add("third")

	fmt.Println("Size:", set.Size())
	fmt.Println("First:", set.GetAt(0))
	fmt.Println("Last:", set.GetAt(2))

	// Output:
	// Size: 3
	// First: first
	// Last: third
}

// ExampleNewOrderedFromSlice demonstrates creating an OrderedSet from a slice
func ExampleNewOrderedFromSlice() {
	// Duplicates are removed, but order is preserved
	slice := []int{3, 1, 4, 1, 5, 9, 2, 6, 5}
	set := cheset.NewOrderedFromSlice(slice)

	fmt.Println("Size:", set.Size())
	fmt.Println("First:", set.GetAt(0))
	fmt.Println("Second:", set.GetAt(1))

	// Output:
	// Size: 7
	// First: 3
	// Second: 1
}

// ExampleOrderedSet_Add demonstrates adding elements in order
func ExampleOrderedSet_Add() {
	set := cheset.NewOrdered[string]()

	set.Add("apple")
	set.Add("banana")
	set.Add("cherry")

	// Insertion order is maintained
	for i := 0; i < set.Size(); i++ {
		fmt.Println(set.GetAt(i))
	}

	// Output:
	// apple
	// banana
	// cherry
}

// ExampleOrderedSet_GetAt demonstrates accessing elements by index
func ExampleOrderedSet_GetAt() {
	set := cheset.NewOrderedFromSlice([]string{"a", "b", "c", "d"})

	fmt.Println("Index 0:", set.GetAt(0))
	fmt.Println("Index 2:", set.GetAt(2))

	// Output:
	// Index 0: a
	// Index 2: c
}

// ExampleOrderedSet_Index demonstrates finding the index of an element
func ExampleOrderedSet_Index() {
	set := cheset.NewOrderedFromSlice([]string{"red", "green", "blue"})

	fmt.Println("Index of 'green':", set.Index("green"))
	fmt.Println("Index of 'yellow':", set.Index("yellow"))

	// Output:
	// Index of 'green': 1
	// Index of 'yellow': -1
}

// ExampleOrderedSet_First demonstrates getting the first element
func ExampleOrderedSet_First() {
	set := cheset.NewOrderedFromSlice([]int{10, 20, 30})

	first, ok := set.First()
	if ok {
		fmt.Println("First element:", first)
	}

	// Output:
	// First element: 10
}

// ExampleOrderedSet_Last demonstrates getting the last element
func ExampleOrderedSet_Last() {
	set := cheset.NewOrderedFromSlice([]int{10, 20, 30})

	last, ok := set.Last()
	if ok {
		fmt.Println("Last element:", last)
	}

	// Output:
	// Last element: 30
}

// ExampleOrderedSet_PopFirst demonstrates removing and returning the first element
func ExampleOrderedSet_PopFirst() {
	set := cheset.NewOrderedFromSlice([]string{"first", "second", "third"})

	item, ok := set.PopFirst()
	if ok {
		fmt.Println("Popped:", item)
		fmt.Println("New first:", set.GetAt(0))
		fmt.Println("Size:", set.Size())
	}

	// Output:
	// Popped: first
	// New first: second
	// Size: 2
}

// ExampleOrderedSet_PopLast demonstrates removing and returning the last element
func ExampleOrderedSet_PopLast() {
	set := cheset.NewOrderedFromSlice([]string{"first", "second", "third"})

	item, ok := set.PopLast()
	if ok {
		fmt.Println("Popped:", item)
		fmt.Println("Size:", set.Size())
	}

	// Output:
	// Popped: third
	// Size: 2
}

// ExampleOrderedSet_Union demonstrates union with order preservation
func ExampleOrderedSet_Union() {
	set1 := cheset.NewOrderedFromSlice([]int{1, 2, 3})
	set2 := cheset.NewOrderedFromSlice([]int{3, 4, 5})

	union := set1.Union(set2)

	// Elements from set1 come first, then new elements from set2
	fmt.Println("Size:", union.Size())
	for i := 0; i < union.Size(); i++ {
		fmt.Println(union.GetAt(i))
	}

	// Output:
	// Size: 5
	// 1
	// 2
	// 3
	// 4
	// 5
}

// ExampleOrderedSet_Intersect demonstrates intersection with order preservation
func ExampleOrderedSet_Intersect() {
	set1 := cheset.NewOrderedFromSlice([]int{1, 2, 3, 4})
	set2 := cheset.NewOrderedFromSlice([]int{3, 4, 5, 6})

	intersection := set1.Intersect(set2)

	// Order is based on set1
	fmt.Println("Size:", intersection.Size())
	fmt.Println("First:", intersection.GetAt(0))
	fmt.Println("Second:", intersection.GetAt(1))

	// Output:
	// Size: 2
	// First: 3
	// Second: 4
}

// ExampleOrderedSet_Filter demonstrates filtering with order preservation
func ExampleOrderedSet_Filter() {
	set := cheset.NewOrderedFromSlice([]int{1, 2, 3, 4, 5, 6, 7, 8})

	evens := set.Filter(func(item int) bool {
		return item%2 == 0
	})

	fmt.Println("Size:", evens.Size())
	for i := 0; i < evens.Size(); i++ {
		fmt.Println(evens.GetAt(i))
	}

	// Output:
	// Size: 4
	// 2
	// 4
	// 6
	// 8
}

// ExampleOrderedSet_Equal demonstrates order-sensitive equality
func ExampleOrderedSet_Equal() {
	set1 := cheset.NewOrderedFromSlice([]int{1, 2, 3})
	set2 := cheset.NewOrderedFromSlice([]int{1, 2, 3})
	set3 := cheset.NewOrderedFromSlice([]int{3, 2, 1})

	fmt.Println("set1 == set2:", set1.Equal(set2))
	fmt.Println("set1 == set3:", set1.Equal(set3)) // Different order

	// Output:
	// set1 == set2: true
	// set1 == set3: false
}

// ExampleOrderedSet_ForEach demonstrates ordered iteration
func ExampleOrderedSet_ForEach() {
	set := cheset.NewOrderedFromSlice([]string{"apple", "banana", "cherry"})

	set.ForEach(func(item string) bool {
		fmt.Println(item)
		return true
	})

	// Output:
	// apple
	// banana
	// cherry
}

// ExampleOrderedSet_ToSlice demonstrates converting to a slice
func ExampleOrderedSet_ToSlice() {
	set := cheset.NewOrderedFromSlice([]int{5, 3, 1, 4, 2})

	slice := set.ToSlice()
	fmt.Println("Slice:", slice)

	// Output:
	// Slice: [5 3 1 4 2]
}

// ExampleOrderedSet_queue demonstrates using OrderedSet as a queue
func ExampleOrderedSet_queue() {
	queue := cheset.NewOrdered[string]()

	// Enqueue
	queue.Add("task1")
	queue.Add("task2")
	queue.Add("task3")

	// Dequeue
	for !queue.IsEmpty() {
		item, _ := queue.PopFirst()
		fmt.Println("Processing:", item)
	}

	// Output:
	// Processing: task1
	// Processing: task2
	// Processing: task3
}

// ExampleOrderedSet_recentItems demonstrates maintaining recent items with order
func ExampleOrderedSet_recentItems() {
	// Track recently accessed items in order
	recent := cheset.NewOrdered[string]()

	// Access items
	items := []string{"page1", "page2", "page1", "page3", "page2"}

	for _, item := range items {
		// If already exists, we'd need to remove and re-add to move to end
		// For this example, we just track unique access order
		if !recent.Contains(item) {
			recent.Add(item)
		}
	}

	fmt.Println("Unique pages visited (in order):")
	for i := 0; i < recent.Size(); i++ {
		fmt.Println(recent.GetAt(i))
	}

	// Output:
	// Unique pages visited (in order):
	// page1
	// page2
	// page3
}

// ExampleOrderedSet_Remove demonstrates removal while maintaining order
func ExampleOrderedSet_Remove() {
	set := cheset.NewOrderedFromSlice([]int{1, 2, 3, 4, 5})

	set.Remove(3) // Remove middle element

	fmt.Println("After removing 3:")
	for i := 0; i < set.Size(); i++ {
		fmt.Println(set.GetAt(i))
	}

	// Output:
	// After removing 3:
	// 1
	// 2
	// 4
	// 5
}
