package chebst

import "golang.org/x/exp/constraints"

// Node represents a single node in a binary search tree
type Node[T constraints.Ordered] struct {
	Value T
	Left  *Node[T]
	Right *Node[T]
}

// BST is a binary search tree implementation
type BST[T constraints.Ordered] struct {
	root *Node[T]
	size int
}

// New creates a new empty Binary Search Tree
func New[T constraints.Ordered]() *BST[T] {
	return &BST[T]{
		root: nil,
		size: 0,
	}
}

// Insert adds a value to the tree - O(log n) average, O(n) worst case
func (bst *BST[T]) Insert(value T) {
	bst.root = bst.insertNode(bst.root, value)
}

func (bst *BST[T]) insertNode(node *Node[T], value T) *Node[T] {
	if node == nil {
		bst.size++
		return &Node[T]{Value: value}
	}

	if value < node.Value {
		node.Left = bst.insertNode(node.Left, value)
	} else if value > node.Value {
		node.Right = bst.insertNode(node.Right, value)
	}
	// If value == node.Value, do nothing (no duplicates)

	return node
}

// Delete removes a value from the tree - O(log n) average, O(n) worst case
// Returns true if the value was found and deleted, false otherwise
func (bst *BST[T]) Delete(value T) bool {
	var deleted bool
	bst.root, deleted = bst.deleteNode(bst.root, value)
	if deleted {
		bst.size--
	}
	return deleted
}

func (bst *BST[T]) deleteNode(node *Node[T], value T) (*Node[T], bool) {
	if node == nil {
		return nil, false
	}

	var deleted bool

	if value < node.Value {
		node.Left, deleted = bst.deleteNode(node.Left, value)
	} else if value > node.Value {
		node.Right, deleted = bst.deleteNode(node.Right, value)
	} else {
		// Node to delete found
		deleted = true

		// Case 1: Node has no children
		if node.Left == nil && node.Right == nil {
			return nil, true
		}

		// Case 2: Node has one child
		if node.Left == nil {
			return node.Right, true
		}
		if node.Right == nil {
			return node.Left, true
		}

		// Case 3: Node has two children
		// Find min value in right subtree (successor)
		minRight := bst.findMin(node.Right)
		node.Value = minRight.Value
		node.Right, _ = bst.deleteNode(node.Right, minRight.Value)
	}

	return node, deleted
}

func (bst *BST[T]) findMin(node *Node[T]) *Node[T] {
	current := node
	for current.Left != nil {
		current = current.Left
	}
	return current
}

// Contains returns true if the tree contains the value - O(log n) average, O(n) worst case
func (bst *BST[T]) Contains(value T) bool {
	return bst.search(bst.root, value)
}

func (bst *BST[T]) search(node *Node[T], value T) bool {
	if node == nil {
		return false
	}

	if value == node.Value {
		return true
	} else if value < node.Value {
		return bst.search(node.Left, value)
	} else {
		return bst.search(node.Right, value)
	}
}

// Min returns the minimum value in the tree - O(log n) average
// Returns the value and true if found, zero value and false if tree is empty
func (bst *BST[T]) Min() (T, bool) {
	if bst.root == nil {
		var zero T
		return zero, false
	}

	node := bst.findMin(bst.root)
	return node.Value, true
}

// Max returns the maximum value in the tree - O(log n) average
// Returns the value and true if found, zero value and false if tree is empty
func (bst *BST[T]) Max() (T, bool) {
	if bst.root == nil {
		var zero T
		return zero, false
	}

	current := bst.root
	for current.Right != nil {
		current = current.Right
	}
	return current.Value, true
}

// Size returns the number of elements in the tree - O(1)
func (bst *BST[T]) Size() int {
	return bst.size
}

// IsEmpty returns true if the tree is empty - O(1)
func (bst *BST[T]) IsEmpty() bool {
	return bst.size == 0
}

// Clear removes all elements from the tree - O(1)
func (bst *BST[T]) Clear() {
	bst.root = nil
	bst.size = 0
}

// Height returns the height of the tree - O(n)
func (bst *BST[T]) Height() int {
	return bst.height(bst.root)
}

func (bst *BST[T]) height(node *Node[T]) int {
	if node == nil {
		return 0
	}

	leftHeight := bst.height(node.Left)
	rightHeight := bst.height(node.Right)

	if leftHeight > rightHeight {
		return leftHeight + 1
	}
	return rightHeight + 1
}

// InOrder returns a slice of values in ascending order - O(n)
func (bst *BST[T]) InOrder() []T {
	result := make([]T, 0, bst.size)
	bst.inOrderTraversal(bst.root, &result)
	return result
}

func (bst *BST[T]) inOrderTraversal(node *Node[T], result *[]T) {
	if node == nil {
		return
	}

	bst.inOrderTraversal(node.Left, result)
	*result = append(*result, node.Value)
	bst.inOrderTraversal(node.Right, result)
}

// PreOrder returns a slice of values in pre-order - O(n)
func (bst *BST[T]) PreOrder() []T {
	result := make([]T, 0, bst.size)
	bst.preOrderTraversal(bst.root, &result)
	return result
}

func (bst *BST[T]) preOrderTraversal(node *Node[T], result *[]T) {
	if node == nil {
		return
	}

	*result = append(*result, node.Value)
	bst.preOrderTraversal(node.Left, result)
	bst.preOrderTraversal(node.Right, result)
}

// PostOrder returns a slice of values in post-order - O(n)
func (bst *BST[T]) PostOrder() []T {
	result := make([]T, 0, bst.size)
	bst.postOrderTraversal(bst.root, &result)
	return result
}

func (bst *BST[T]) postOrderTraversal(node *Node[T], result *[]T) {
	if node == nil {
		return
	}

	bst.postOrderTraversal(node.Left, result)
	bst.postOrderTraversal(node.Right, result)
	*result = append(*result, node.Value)
}

// ForEach iterates over each element in order - O(n)
// The function receives the value and returns true to continue, false to stop
func (bst *BST[T]) ForEach(fn func(T) bool) {
	bst.forEach(bst.root, fn)
}

func (bst *BST[T]) forEach(node *Node[T], fn func(T) bool) bool {
	if node == nil {
		return true
	}

	if !bst.forEach(node.Left, fn) {
		return false
	}

	if !fn(node.Value) {
		return false
	}

	return bst.forEach(node.Right, fn)
}

// Find returns the first element for which the predicate returns true - O(n)
// Returns the element and true if found, zero value and false otherwise
func (bst *BST[T]) Find(predicate func(T) bool) (T, bool) {
	return bst.find(bst.root, predicate)
}

func (bst *BST[T]) find(node *Node[T], predicate func(T) bool) (T, bool) {
	if node == nil {
		var zero T
		return zero, false
	}

	// Search left subtree
	if value, found := bst.find(node.Left, predicate); found {
		return value, true
	}

	// Check current node
	if predicate(node.Value) {
		return node.Value, true
	}

	// Search right subtree
	return bst.find(node.Right, predicate)
}

// Clone creates a deep copy of the tree - O(n)
func (bst *BST[T]) Clone() *BST[T] {
	newBst := New[T]()
	newBst.root = bst.cloneNode(bst.root)
	newBst.size = bst.size
	return newBst
}

func (bst *BST[T]) cloneNode(node *Node[T]) *Node[T] {
	if node == nil {
		return nil
	}

	return &Node[T]{
		Value: node.Value,
		Left:  bst.cloneNode(node.Left),
		Right: bst.cloneNode(node.Right),
	}
}
