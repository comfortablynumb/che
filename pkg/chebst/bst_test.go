package chebst

import (
	"testing"

	"github.com/comfortablynumb/che/pkg/chetest"
)

func TestNew(t *testing.T) {
	bst := NewOrdered[int]()

	chetest.RequireEqual(t, bst.IsEmpty(), true)
	chetest.RequireEqual(t, bst.Size(), 0)
}

func TestInsert(t *testing.T) {
	bst := NewOrdered[int]()

	bst.Insert(5)
	bst.Insert(3)
	bst.Insert(7)
	bst.Insert(1)
	bst.Insert(9)

	chetest.RequireEqual(t, bst.Size(), 5)
	chetest.RequireEqual(t, bst.IsEmpty(), false)
}

func TestInsert_Duplicates(t *testing.T) {
	bst := NewOrdered[int]()

	bst.Insert(5)
	bst.Insert(5)
	bst.Insert(5)

	// Duplicates should not be inserted
	chetest.RequireEqual(t, bst.Size(), 1)
}

func TestContains(t *testing.T) {
	bst := NewOrdered[int]()

	bst.Insert(5)
	bst.Insert(3)
	bst.Insert(7)

	chetest.RequireEqual(t, bst.Contains(5), true)
	chetest.RequireEqual(t, bst.Contains(3), true)
	chetest.RequireEqual(t, bst.Contains(7), true)
	chetest.RequireEqual(t, bst.Contains(10), false)
}

func TestContains_Empty(t *testing.T) {
	bst := NewOrdered[int]()

	chetest.RequireEqual(t, bst.Contains(1), false)
}

func TestDelete(t *testing.T) {
	bst := NewOrdered[int]()

	bst.Insert(5)
	bst.Insert(3)
	bst.Insert(7)

	deleted := bst.Delete(3)
	chetest.RequireEqual(t, deleted, true)
	chetest.RequireEqual(t, bst.Size(), 2)
	chetest.RequireEqual(t, bst.Contains(3), false)
}

func TestDelete_Root(t *testing.T) {
	bst := NewOrdered[int]()

	bst.Insert(5)
	bst.Insert(3)
	bst.Insert(7)

	deleted := bst.Delete(5)
	chetest.RequireEqual(t, deleted, true)
	chetest.RequireEqual(t, bst.Size(), 2)
	chetest.RequireEqual(t, bst.Contains(5), false)
}

func TestDelete_NotFound(t *testing.T) {
	bst := NewOrdered[int]()

	bst.Insert(5)

	deleted := bst.Delete(10)
	chetest.RequireEqual(t, deleted, false)
	chetest.RequireEqual(t, bst.Size(), 1)
}

func TestDelete_Empty(t *testing.T) {
	bst := NewOrdered[int]()

	deleted := bst.Delete(1)
	chetest.RequireEqual(t, deleted, false)
}

func TestDelete_LeafNode(t *testing.T) {
	bst := NewOrdered[int]()

	bst.Insert(5)
	bst.Insert(3)
	bst.Insert(7)
	bst.Insert(1)

	deleted := bst.Delete(1)
	chetest.RequireEqual(t, deleted, true)
	chetest.RequireEqual(t, bst.Size(), 3)
	chetest.RequireEqual(t, bst.InOrder(), []int{3, 5, 7})
}

func TestDelete_NodeWithOneChild(t *testing.T) {
	bst := NewOrdered[int]()

	bst.Insert(5)
	bst.Insert(3)
	bst.Insert(1)

	deleted := bst.Delete(3)
	chetest.RequireEqual(t, deleted, true)
	chetest.RequireEqual(t, bst.Size(), 2)
	chetest.RequireEqual(t, bst.InOrder(), []int{1, 5})
}

func TestDelete_NodeWithTwoChildren(t *testing.T) {
	bst := NewOrdered[int]()

	bst.Insert(5)
	bst.Insert(3)
	bst.Insert(7)
	bst.Insert(1)
	bst.Insert(4)

	deleted := bst.Delete(3)
	chetest.RequireEqual(t, deleted, true)
	chetest.RequireEqual(t, bst.Size(), 4)
	chetest.RequireEqual(t, bst.InOrder(), []int{1, 4, 5, 7})
}

func TestMin(t *testing.T) {
	bst := NewOrdered[int]()

	bst.Insert(5)
	bst.Insert(3)
	bst.Insert(7)
	bst.Insert(1)

	min, ok := bst.Min()
	chetest.RequireEqual(t, ok, true)
	chetest.RequireEqual(t, min, 1)
}

func TestMin_Empty(t *testing.T) {
	bst := NewOrdered[int]()

	min, ok := bst.Min()
	chetest.RequireEqual(t, ok, false)
	chetest.RequireEqual(t, min, 0)
}

func TestMax(t *testing.T) {
	bst := NewOrdered[int]()

	bst.Insert(5)
	bst.Insert(3)
	bst.Insert(7)
	bst.Insert(9)

	max, ok := bst.Max()
	chetest.RequireEqual(t, ok, true)
	chetest.RequireEqual(t, max, 9)
}

func TestMax_Empty(t *testing.T) {
	bst := NewOrdered[int]()

	max, ok := bst.Max()
	chetest.RequireEqual(t, ok, false)
	chetest.RequireEqual(t, max, 0)
}

func TestSize(t *testing.T) {
	bst := NewOrdered[int]()

	chetest.RequireEqual(t, bst.Size(), 0)

	bst.Insert(5)
	chetest.RequireEqual(t, bst.Size(), 1)

	bst.Insert(3)
	chetest.RequireEqual(t, bst.Size(), 2)

	bst.Delete(5)
	chetest.RequireEqual(t, bst.Size(), 1)
}

func TestIsEmpty(t *testing.T) {
	bst := NewOrdered[int]()

	chetest.RequireEqual(t, bst.IsEmpty(), true)

	bst.Insert(5)
	chetest.RequireEqual(t, bst.IsEmpty(), false)

	bst.Delete(5)
	chetest.RequireEqual(t, bst.IsEmpty(), true)
}

func TestClear(t *testing.T) {
	bst := NewOrdered[int]()

	bst.Insert(5)
	bst.Insert(3)
	bst.Insert(7)

	bst.Clear()

	chetest.RequireEqual(t, bst.IsEmpty(), true)
	chetest.RequireEqual(t, bst.Size(), 0)
}

func TestHeight(t *testing.T) {
	bst := NewOrdered[int]()

	chetest.RequireEqual(t, bst.Height(), 0)

	bst.Insert(5)
	chetest.RequireEqual(t, bst.Height(), 1)

	bst.Insert(3)
	chetest.RequireEqual(t, bst.Height(), 2)

	bst.Insert(7)
	chetest.RequireEqual(t, bst.Height(), 2)

	bst.Insert(1)
	chetest.RequireEqual(t, bst.Height(), 3)
}

func TestInOrder(t *testing.T) {
	bst := NewOrdered[int]()

	bst.Insert(5)
	bst.Insert(3)
	bst.Insert(7)
	bst.Insert(1)
	bst.Insert(9)

	result := bst.InOrder()

	// InOrder should return sorted values
	chetest.RequireEqual(t, result, []int{1, 3, 5, 7, 9})
}

func TestInOrder_Empty(t *testing.T) {
	bst := NewOrdered[int]()

	result := bst.InOrder()

	chetest.RequireEqual(t, len(result), 0)
}

func TestPreOrder(t *testing.T) {
	bst := NewOrdered[int]()

	bst.Insert(5)
	bst.Insert(3)
	bst.Insert(7)
	bst.Insert(1)
	bst.Insert(9)

	result := bst.PreOrder()

	// PreOrder: root -> left -> right
	chetest.RequireEqual(t, result, []int{5, 3, 1, 7, 9})
}

func TestPreOrder_Empty(t *testing.T) {
	bst := NewOrdered[int]()

	result := bst.PreOrder()

	chetest.RequireEqual(t, len(result), 0)
}

func TestPostOrder(t *testing.T) {
	bst := NewOrdered[int]()

	bst.Insert(5)
	bst.Insert(3)
	bst.Insert(7)
	bst.Insert(1)
	bst.Insert(9)

	result := bst.PostOrder()

	// PostOrder: left -> right -> root
	chetest.RequireEqual(t, result, []int{1, 3, 9, 7, 5})
}

func TestPostOrder_Empty(t *testing.T) {
	bst := NewOrdered[int]()

	result := bst.PostOrder()

	chetest.RequireEqual(t, len(result), 0)
}

func TestForEach(t *testing.T) {
	bst := NewOrdered[int]()

	bst.Insert(5)
	bst.Insert(3)
	bst.Insert(7)

	values := []int{}
	bst.ForEach(func(value int) bool {
		values = append(values, value)
		return true
	})

	chetest.RequireEqual(t, values, []int{3, 5, 7})
}

func TestForEach_EarlyExit(t *testing.T) {
	bst := NewOrdered[int]()

	bst.Insert(5)
	bst.Insert(3)
	bst.Insert(7)

	count := 0
	bst.ForEach(func(value int) bool {
		count++
		return count < 2
	})

	chetest.RequireEqual(t, count, 2)
}

func TestFind(t *testing.T) {
	bst := NewOrdered[int]()

	bst.Insert(5)
	bst.Insert(3)
	bst.Insert(7)

	value, found := bst.Find(func(v int) bool {
		return v > 4
	})

	chetest.RequireEqual(t, found, true)
	chetest.RequireEqual(t, value, 5)
}

func TestFind_NotFound(t *testing.T) {
	bst := NewOrdered[int]()

	bst.Insert(5)
	bst.Insert(3)

	value, found := bst.Find(func(v int) bool {
		return v > 10
	})

	chetest.RequireEqual(t, found, false)
	chetest.RequireEqual(t, value, 0)
}

func TestClone(t *testing.T) {
	bst := NewOrdered[int]()

	bst.Insert(5)
	bst.Insert(3)
	bst.Insert(7)

	clone := bst.Clone()

	chetest.RequireEqual(t, clone.InOrder(), []int{3, 5, 7})
	chetest.RequireEqual(t, clone.Size(), 3)

	// Modify clone
	clone.Insert(1)

	// Original should be unchanged
	chetest.RequireEqual(t, bst.InOrder(), []int{3, 5, 7})
	chetest.RequireEqual(t, clone.InOrder(), []int{1, 3, 5, 7})
}

func TestString(t *testing.T) {
	bst := NewOrdered[string]()

	bst.Insert("dog")
	bst.Insert("cat")
	bst.Insert("elephant")
	bst.Insert("ant")

	result := bst.InOrder()

	// Should be sorted alphabetically
	chetest.RequireEqual(t, result, []string{"ant", "cat", "dog", "elephant"})
}

func TestComplexOperations(t *testing.T) {
	bst := NewOrdered[int]()

	// Insert multiple values
	values := []int{50, 30, 70, 20, 40, 60, 80}
	for _, v := range values {
		bst.Insert(v)
	}

	chetest.RequireEqual(t, bst.Size(), 7)
	chetest.RequireEqual(t, bst.Height(), 3)

	// Delete some values
	bst.Delete(20)
	bst.Delete(70)

	chetest.RequireEqual(t, bst.Size(), 5)
	chetest.RequireEqual(t, bst.InOrder(), []int{30, 40, 50, 60, 80})

	min, _ := bst.Min()
	max, _ := bst.Max()

	chetest.RequireEqual(t, min, 30)
	chetest.RequireEqual(t, max, 80)
}
