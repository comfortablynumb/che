package chemap_test

import (
	"fmt"

	"github.com/comfortablynumb/che/pkg/chemap"
)

// Example demonstrates basic Multimap operations
func Example() {
	mm := chemap.NewMultimap[string, int]()

	// Add multiple values for same key
	mm.Put("numbers", 1)
	mm.Put("numbers", 2)
	mm.Put("numbers", 3)

	// Get all values
	values := mm.Get("numbers")
	fmt.Println("Values:", values)

	fmt.Println("Size:", mm.Size())
	fmt.Println("Key count:", mm.KeyCount())

	// Output:
	// Values: [1 2 3]
	// Size: 3
	// Key count: 1
}

// ExampleNewMultimap demonstrates creating a new multimap
func ExampleNewMultimap() {
	mm := chemap.NewMultimap[string, string]()

	mm.Put("fruits", "apple")
	mm.Put("fruits", "banana")
	mm.Put("vegetables", "carrot")

	fmt.Println("Fruit count:", mm.ValueCount("fruits"))
	fmt.Println("Total keys:", mm.KeyCount())

	// Output:
	// Fruit count: 2
	// Total keys: 2
}

// ExampleMultimap_Put demonstrates adding values
func ExampleMultimap_Put() {
	mm := chemap.NewMultimap[string, int]()

	mm.Put("scores", 85)
	mm.Put("scores", 92)
	mm.Put("scores", 78)

	fmt.Println("Total scores:", mm.ValueCount("scores"))

	// Output:
	// Total scores: 3
}

// ExampleMultimap_PutAll demonstrates adding multiple values at once
func ExampleMultimap_PutAll() {
	mm := chemap.NewMultimap[string, string]()

	mm.PutAll("colors", "red", "green", "blue")

	fmt.Println("Color count:", mm.ValueCount("colors"))

	// Output:
	// Color count: 3
}

// ExampleMultimap_Get demonstrates retrieving values
func ExampleMultimap_Get() {
	mm := chemap.NewMultimap[string, int]()

	mm.PutAll("primes", 2, 3, 5, 7, 11)

	primes := mm.Get("primes")
	fmt.Println("Primes:", primes)

	// Output:
	// Primes: [2 3 5 7 11]
}

// ExampleMultimap_GetFirst demonstrates getting the first value
func ExampleMultimap_GetFirst() {
	mm := chemap.NewMultimap[string, string]()

	mm.PutAll("queue", "task1", "task2", "task3")

	first, ok := mm.GetFirst("queue")
	fmt.Println("First task:", first, "Found:", ok)

	// Output:
	// First task: task1 Found: true
}

// ExampleMultimap_Remove demonstrates removing a specific value
func ExampleMultimap_Remove() {
	mm := chemap.NewMultimap[string, int]()

	mm.PutAll("numbers", 1, 2, 3, 4, 5)

	equals := func(a, b int) bool { return a == b }
	removed := mm.Remove("numbers", 3, equals)

	fmt.Println("Removed:", removed)
	fmt.Println("Remaining:", mm.Get("numbers"))

	// Output:
	// Removed: true
	// Remaining: [1 2 4 5]
}

// ExampleMultimap_RemoveAll demonstrates removing all values for a key
func ExampleMultimap_RemoveAll() {
	mm := chemap.NewMultimap[string, string]()

	mm.PutAll("old", "value1", "value2", "value3")
	mm.Put("keep", "important")

	mm.RemoveAll("old")

	fmt.Println("Keys:", mm.Keys())

	// Output:
	// Keys: [keep]
}

// ExampleMultimap_ForEach demonstrates iterating over all entries
func ExampleMultimap_ForEach() {
	mm := chemap.NewMultimap[string, int]()

	mm.PutAll("a", 1, 2)
	mm.Put("b", 3)

	mm.ForEach(func(key string, value int) bool {
		fmt.Printf("%s: %d\n", key, value)
		return true
	})

	// Output:
	// a: 1
	// a: 2
	// b: 3
}

// ExampleMultimap_ForEachKey demonstrates iterating over keys
func ExampleMultimap_ForEachKey() {
	mm := chemap.NewMultimap[string, int]()

	mm.PutAll("evens", 2, 4, 6)
	mm.PutAll("odds", 1, 3, 5)

	mm.ForEachKey(func(key string, values []int) bool {
		fmt.Printf("%s has %d values\n", key, len(values))
		return true
	})

	// Output:
	// evens has 3 values
	// odds has 3 values
}

// ExampleMultimap_Clone demonstrates cloning a multimap
func ExampleMultimap_Clone() {
	original := chemap.NewMultimap[string, int]()
	original.PutAll("data", 1, 2, 3)

	clone := original.Clone()
	clone.Put("data", 4)

	fmt.Println("Original:", original.Get("data"))
	fmt.Println("Clone:", clone.Get("data"))

	// Output:
	// Original: [1 2 3]
	// Clone: [1 2 3 4]
}

// ExampleMultimap_Merge demonstrates merging multimaps
func ExampleMultimap_Merge() {
	mm1 := chemap.NewMultimap[string, int]()
	mm2 := chemap.NewMultimap[string, int]()

	mm1.PutAll("shared", 1, 2)
	mm2.PutAll("shared", 3, 4)
	mm2.Put("unique", 5)

	mm1.Merge(mm2)

	fmt.Println("Shared values:", mm1.Get("shared"))
	fmt.Println("Total size:", mm1.Size())

	// Output:
	// Shared values: [1 2 3 4]
	// Total size: 5
}

// ExampleMultimap_ReplaceValues demonstrates replacing values
func ExampleMultimap_ReplaceValues() {
	mm := chemap.NewMultimap[string, string]()

	mm.PutAll("tags", "old1", "old2", "old3")
	fmt.Println("Before:", mm.Get("tags"))

	mm.ReplaceValues("tags", "new1", "new2")
	fmt.Println("After:", mm.Get("tags"))

	// Output:
	// Before: [old1 old2 old3]
	// After: [new1 new2]
}

// ExampleMultimap_headers demonstrates HTTP header management
func ExampleMultimap_headers() {
	headers := chemap.NewMultimap[string, string]()

	// HTTP headers can have multiple values
	headers.Put("Accept", "application/json")
	headers.Put("Accept", "text/html")
	headers.Put("Cache-Control", "no-cache")
	headers.Put("Cache-Control", "no-store")

	fmt.Println("Accept headers:", headers.Get("Accept"))
	fmt.Println("Cache-Control:", headers.Get("Cache-Control"))

	// Output:
	// Accept headers: [application/json text/html]
	// Cache-Control: [no-cache no-store]
}

// ExampleMultimap_index demonstrates building an inverted index
func ExampleMultimap_index() {
	// Document ID -> words index
	index := chemap.NewMultimap[string, int]()

	// Document 1 contains words: "go", "programming"
	index.Put("go", 1)
	index.Put("programming", 1)

	// Document 2 contains words: "go", "language"
	index.Put("go", 2)
	index.Put("language", 2)

	// Find documents containing "go"
	docs := index.Get("go")
	fmt.Println("Documents with 'go':", docs)

	// Output:
	// Documents with 'go': [1 2]
}

// ExampleMultimap_graph demonstrates representing a graph
func ExampleMultimap_graph() {
	// Adjacency list representation
	graph := chemap.NewMultimap[int, int]()

	// Node 1 connects to nodes 2 and 3
	graph.PutAll(1, 2, 3)

	// Node 2 connects to nodes 3 and 4
	graph.PutAll(2, 3, 4)

	// Get neighbors of node 1
	neighbors := graph.Get(1)
	fmt.Println("Node 1 neighbors:", neighbors)

	// Output:
	// Node 1 neighbors: [2 3]
}
