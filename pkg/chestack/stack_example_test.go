package chestack_test

import (
	"fmt"

	"github.com/comfortablynumb/che/pkg/chestack"
)

// Example demonstrates basic Stack operations
func Example() {
	stack := chestack.New[int]()

	stack.Push(1)
	stack.Push(2)
	stack.Push(3)

	// LIFO order
	val, _ := stack.Pop()
	fmt.Println("First out:", val)

	val, _ = stack.Pop()
	fmt.Println("Second out:", val)

	fmt.Println("Size:", stack.Size())

	// Output:
	// First out: 3
	// Second out: 2
	// Size: 1
}

// ExampleNew demonstrates creating a new stack
func ExampleNew() {
	stack := chestack.New[string]()
	stack.Push("hello")
	stack.Push("world")

	fmt.Println("Size:", stack.Size())

	// Output:
	// Size: 2
}

// ExampleNewFromSlice demonstrates creating a stack from a slice
func ExampleNewFromSlice() {
	slice := []int{1, 2, 3, 4, 5}
	stack := chestack.NewFromSlice(slice)

	// Last element from slice is on top
	val, _ := stack.Pop()
	fmt.Println("Top:", val)

	// Output:
	// Top: 5
}

// ExampleStack_Push demonstrates adding elements
func ExampleStack_Push() {
	stack := chestack.New[int]()

	stack.Push(10)
	stack.Push(20)
	stack.Push(30)

	fmt.Println("Size:", stack.Size())

	// Output:
	// Size: 3
}

// ExampleStack_Pop demonstrates removing elements
func ExampleStack_Pop() {
	stack := chestack.NewFromSlice([]int{1, 2, 3})

	val, ok := stack.Pop()
	fmt.Println("Value:", val, "Success:", ok)

	val, _ = stack.Pop()
	fmt.Println("Next:", val)

	// Output:
	// Value: 3 Success: true
	// Next: 2
}

// ExampleStack_Peek demonstrates peeking at the top
func ExampleStack_Peek() {
	stack := chestack.NewFromSlice([]int{100, 200, 300})

	// Peek doesn't remove the element
	val, ok := stack.Peek()
	fmt.Println("Top:", val, "Success:", ok)
	fmt.Println("Size still:", stack.Size())

	// Output:
	// Top: 300 Success: true
	// Size still: 3
}

// ExampleStack_ForEach demonstrates iterating from bottom to top
func ExampleStack_ForEach() {
	stack := chestack.NewFromSlice([]string{"a", "b", "c"})

	fmt.Println("Bottom to top:")
	stack.ForEach(func(item string) bool {
		fmt.Println(item)
		return true
	})

	// Output:
	// Bottom to top:
	// a
	// b
	// c
}

// ExampleStack_ForEachReverse demonstrates iterating from top to bottom
func ExampleStack_ForEachReverse() {
	stack := chestack.NewFromSlice([]string{"a", "b", "c"})

	fmt.Println("Top to bottom:")
	stack.ForEachReverse(func(item string) bool {
		fmt.Println(item)
		return true
	})

	// Output:
	// Top to bottom:
	// c
	// b
	// a
}

// ExampleStack_Clone demonstrates cloning a stack
func ExampleStack_Clone() {
	original := chestack.NewFromSlice([]int{1, 2, 3})
	clone := original.Clone()

	clone.Push(4)

	fmt.Println("Original size:", original.Size())
	fmt.Println("Clone size:", clone.Size())

	// Output:
	// Original size: 3
	// Clone size: 4
}

// ExampleStack_undo demonstrates using Stack for undo functionality
func ExampleStack_undo() {
	type Action struct {
		Type string
		Data int
	}

	undoStack := chestack.New[Action]()

	// Perform actions
	undoStack.Push(Action{Type: "insert", Data: 1})
	undoStack.Push(Action{Type: "delete", Data: 2})
	undoStack.Push(Action{Type: "update", Data: 3})

	// Undo last two actions
	fmt.Println("Undoing...")
	for i := 0; i < 2; i++ {
		action, _ := undoStack.Pop()
		fmt.Printf("Undo %s (data: %d)\n", action.Type, action.Data)
	}

	fmt.Println("Remaining actions:", undoStack.Size())

	// Output:
	// Undoing...
	// Undo update (data: 3)
	// Undo delete (data: 2)
	// Remaining actions: 1
}

// ExampleStack_dfs demonstrates using Stack for depth-first search
func ExampleStack_dfs() {
	type Node struct {
		Value    int
		Children []int
	}

	// Simple tree: 1 -> [2, 3], 2 -> [4, 5]
	nodes := map[int]Node{
		1: {Value: 1, Children: []int{2, 3}},
		2: {Value: 2, Children: []int{4, 5}},
		3: {Value: 3, Children: []int{}},
		4: {Value: 4, Children: []int{}},
		5: {Value: 5, Children: []int{}},
	}

	stack := chestack.New[int]()
	stack.Push(1) // Start with root

	for !stack.IsEmpty() {
		nodeID, _ := stack.Pop()
		node := nodes[nodeID]
		fmt.Println(node.Value)

		// Push children in reverse order for correct DFS
		for i := len(node.Children) - 1; i >= 0; i-- {
			stack.Push(node.Children[i])
		}
	}

	// Output:
	// 1
	// 2
	// 4
	// 5
	// 3
}

// ExampleStack_brackets demonstrates using Stack for bracket matching
func ExampleStack_brackets() {
	checkBrackets := func(s string) bool {
		stack := chestack.New[rune]()
		pairs := map[rune]rune{')': '(', ']': '[', '}': '{'}

		for _, char := range s {
			if char == '(' || char == '[' || char == '{' {
				stack.Push(char)
			} else if opening, isClosing := pairs[char]; isClosing {
				if top, ok := stack.Pop(); !ok || top != opening {
					return false
				}
			}
		}
		return stack.IsEmpty()
	}

	fmt.Println("Valid brackets:", checkBrackets("{[()]}"))
	fmt.Println("Invalid brackets:", checkBrackets("{[(])}"))

	// Output:
	// Valid brackets: true
	// Invalid brackets: false
}

// ExampleStack_ToSlice demonstrates converting stack to slice
func ExampleStack_ToSlice() {
	stack := chestack.NewFromSlice([]int{10, 20, 30})

	// Slice shows bottom to top
	slice := stack.ToSlice()
	fmt.Println("Slice:", slice)

	// Output:
	// Slice: [10 20 30]
}

// ExampleStack_expressionEvaluation demonstrates evaluating postfix expressions
func ExampleStack_expressionEvaluation() {
	evaluatePostfix := func(tokens []string) int {
		stack := chestack.New[int]()

		for _, token := range tokens {
			switch token {
			case "+", "-", "*", "/":
				b, _ := stack.Pop()
				a, _ := stack.Pop()
				var result int
				switch token {
				case "+":
					result = a + b
				case "-":
					result = a - b
				case "*":
					result = a * b
				case "/":
					result = a / b
				}
				stack.Push(result)
			default:
				// Parse number
				var num int
				_, _ = fmt.Sscanf(token, "%d", &num)
				stack.Push(num)
			}
		}

		result, _ := stack.Pop()
		return result
	}

	// Evaluate: 5 + ((1 + 2) * 4) - 3 = 5 + 12 - 3 = 14
	// Postfix: 5 1 2 + 4 * + 3 -
	tokens := []string{"5", "1", "2", "+", "4", "*", "+", "3", "-"}
	result := evaluatePostfix(tokens)
	fmt.Println("Result:", result)

	// Output:
	// Result: 14
}
