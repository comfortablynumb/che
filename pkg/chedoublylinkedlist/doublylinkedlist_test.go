package chedoublylinkedlist

import (
	"testing"

	"github.com/comfortablynumb/che/pkg/chetest"
)

func TestNew(t *testing.T) {
	dll := New[int]()

	chetest.RequireEqual(t, dll.IsEmpty(), true)
	chetest.RequireEqual(t, dll.Size(), 0)
}

func TestPrepend(t *testing.T) {
	dll := New[int]()

	dll.Prepend(3)
	dll.Prepend(2)
	dll.Prepend(1)

	chetest.RequireEqual(t, dll.Size(), 3)
	chetest.RequireEqual(t, dll.ToSlice(), []int{1, 2, 3})
}

func TestAppend(t *testing.T) {
	dll := New[int]()

	dll.Append(1)
	dll.Append(2)
	dll.Append(3)

	chetest.RequireEqual(t, dll.Size(), 3)
	chetest.RequireEqual(t, dll.ToSlice(), []int{1, 2, 3})
}

func TestInsertAt(t *testing.T) {
	dll := New[int]()

	dll.Append(1)
	dll.Append(3)

	ok := dll.InsertAt(1, 2)
	chetest.RequireEqual(t, ok, true)
	chetest.RequireEqual(t, dll.ToSlice(), []int{1, 2, 3})

	ok = dll.InsertAt(0, 0)
	chetest.RequireEqual(t, ok, true)
	chetest.RequireEqual(t, dll.ToSlice(), []int{0, 1, 2, 3})

	ok = dll.InsertAt(4, 4)
	chetest.RequireEqual(t, ok, true)
	chetest.RequireEqual(t, dll.ToSlice(), []int{0, 1, 2, 3, 4})
}

func TestInsertAt_OutOfBounds(t *testing.T) {
	dll := New[int]()
	dll.Append(1)

	ok := dll.InsertAt(-1, 0)
	chetest.RequireEqual(t, ok, false)

	ok = dll.InsertAt(10, 0)
	chetest.RequireEqual(t, ok, false)
}

func TestRemoveFirst(t *testing.T) {
	dll := New[int]()
	dll.Append(1)
	dll.Append(2)
	dll.Append(3)

	value, ok := dll.RemoveFirst()
	chetest.RequireEqual(t, ok, true)
	chetest.RequireEqual(t, value, 1)
	chetest.RequireEqual(t, dll.Size(), 2)
	chetest.RequireEqual(t, dll.ToSlice(), []int{2, 3})
}

func TestRemoveFirst_Empty(t *testing.T) {
	dll := New[int]()

	value, ok := dll.RemoveFirst()
	chetest.RequireEqual(t, ok, false)
	chetest.RequireEqual(t, value, 0)
}

func TestRemoveFirst_SingleElement(t *testing.T) {
	dll := New[int]()
	dll.Append(1)

	value, ok := dll.RemoveFirst()
	chetest.RequireEqual(t, ok, true)
	chetest.RequireEqual(t, value, 1)
	chetest.RequireEqual(t, dll.IsEmpty(), true)
}

func TestRemoveLast(t *testing.T) {
	dll := New[int]()
	dll.Append(1)
	dll.Append(2)
	dll.Append(3)

	value, ok := dll.RemoveLast()
	chetest.RequireEqual(t, ok, true)
	chetest.RequireEqual(t, value, 3)
	chetest.RequireEqual(t, dll.Size(), 2)
	chetest.RequireEqual(t, dll.ToSlice(), []int{1, 2})
}

func TestRemoveLast_Empty(t *testing.T) {
	dll := New[int]()

	value, ok := dll.RemoveLast()
	chetest.RequireEqual(t, ok, false)
	chetest.RequireEqual(t, value, 0)
}

func TestRemoveLast_SingleElement(t *testing.T) {
	dll := New[int]()
	dll.Append(1)

	value, ok := dll.RemoveLast()
	chetest.RequireEqual(t, ok, true)
	chetest.RequireEqual(t, value, 1)
	chetest.RequireEqual(t, dll.IsEmpty(), true)
}

func TestRemoveAt(t *testing.T) {
	dll := New[int]()
	dll.Append(1)
	dll.Append(2)
	dll.Append(3)
	dll.Append(4)

	value, ok := dll.RemoveAt(1)
	chetest.RequireEqual(t, ok, true)
	chetest.RequireEqual(t, value, 2)
	chetest.RequireEqual(t, dll.ToSlice(), []int{1, 3, 4})

	value, ok = dll.RemoveAt(2)
	chetest.RequireEqual(t, ok, true)
	chetest.RequireEqual(t, value, 4)
	chetest.RequireEqual(t, dll.ToSlice(), []int{1, 3})
}

func TestRemoveAt_OutOfBounds(t *testing.T) {
	dll := New[int]()
	dll.Append(1)

	value, ok := dll.RemoveAt(-1)
	chetest.RequireEqual(t, ok, false)
	chetest.RequireEqual(t, value, 0)

	value, ok = dll.RemoveAt(10)
	chetest.RequireEqual(t, ok, false)
	chetest.RequireEqual(t, value, 0)
}

func TestGet(t *testing.T) {
	dll := New[int]()
	dll.Append(10)
	dll.Append(20)
	dll.Append(30)

	value, ok := dll.Get(0)
	chetest.RequireEqual(t, ok, true)
	chetest.RequireEqual(t, value, 10)

	value, ok = dll.Get(1)
	chetest.RequireEqual(t, ok, true)
	chetest.RequireEqual(t, value, 20)

	value, ok = dll.Get(2)
	chetest.RequireEqual(t, ok, true)
	chetest.RequireEqual(t, value, 30)
}

func TestGet_OutOfBounds(t *testing.T) {
	dll := New[int]()
	dll.Append(1)

	value, ok := dll.Get(-1)
	chetest.RequireEqual(t, ok, false)
	chetest.RequireEqual(t, value, 0)

	value, ok = dll.Get(10)
	chetest.RequireEqual(t, ok, false)
	chetest.RequireEqual(t, value, 0)
}

func TestFirst(t *testing.T) {
	dll := New[int]()
	dll.Append(10)
	dll.Append(20)

	value, ok := dll.First()
	chetest.RequireEqual(t, ok, true)
	chetest.RequireEqual(t, value, 10)
}

func TestFirst_Empty(t *testing.T) {
	dll := New[int]()

	value, ok := dll.First()
	chetest.RequireEqual(t, ok, false)
	chetest.RequireEqual(t, value, 0)
}

func TestLast(t *testing.T) {
	dll := New[int]()
	dll.Append(10)
	dll.Append(20)

	value, ok := dll.Last()
	chetest.RequireEqual(t, ok, true)
	chetest.RequireEqual(t, value, 20)
}

func TestLast_Empty(t *testing.T) {
	dll := New[int]()

	value, ok := dll.Last()
	chetest.RequireEqual(t, ok, false)
	chetest.RequireEqual(t, value, 0)
}

func TestSize(t *testing.T) {
	dll := New[int]()

	chetest.RequireEqual(t, dll.Size(), 0)

	dll.Append(1)
	chetest.RequireEqual(t, dll.Size(), 1)

	dll.Append(2)
	chetest.RequireEqual(t, dll.Size(), 2)

	dll.RemoveFirst()
	chetest.RequireEqual(t, dll.Size(), 1)
}

func TestIsEmpty(t *testing.T) {
	dll := New[int]()

	chetest.RequireEqual(t, dll.IsEmpty(), true)

	dll.Append(1)
	chetest.RequireEqual(t, dll.IsEmpty(), false)

	dll.RemoveFirst()
	chetest.RequireEqual(t, dll.IsEmpty(), true)
}

func TestClear(t *testing.T) {
	dll := New[int]()
	dll.Append(1)
	dll.Append(2)
	dll.Append(3)

	dll.Clear()

	chetest.RequireEqual(t, dll.IsEmpty(), true)
	chetest.RequireEqual(t, dll.Size(), 0)
}

func TestToSlice(t *testing.T) {
	dll := New[int]()
	dll.Append(1)
	dll.Append(2)
	dll.Append(3)

	slice := dll.ToSlice()

	chetest.RequireEqual(t, slice, []int{1, 2, 3})
}

func TestToSlice_Empty(t *testing.T) {
	dll := New[int]()

	slice := dll.ToSlice()

	chetest.RequireEqual(t, len(slice), 0)
}

func TestToSliceReverse(t *testing.T) {
	dll := New[int]()
	dll.Append(1)
	dll.Append(2)
	dll.Append(3)

	slice := dll.ToSliceReverse()

	chetest.RequireEqual(t, slice, []int{3, 2, 1})
}

func TestToSliceReverse_Empty(t *testing.T) {
	dll := New[int]()

	slice := dll.ToSliceReverse()

	chetest.RequireEqual(t, len(slice), 0)
}

func TestForEach(t *testing.T) {
	dll := New[int]()
	dll.Append(1)
	dll.Append(2)
	dll.Append(3)

	sum := 0
	dll.ForEach(func(value int) bool {
		sum += value
		return true
	})

	chetest.RequireEqual(t, sum, 6)
}

func TestForEach_EarlyExit(t *testing.T) {
	dll := New[int]()
	dll.Append(1)
	dll.Append(2)
	dll.Append(3)

	count := 0
	dll.ForEach(func(value int) bool {
		count++
		return count < 2
	})

	chetest.RequireEqual(t, count, 2)
}

func TestForEachReverse(t *testing.T) {
	dll := New[int]()
	dll.Append(1)
	dll.Append(2)
	dll.Append(3)

	values := []int{}
	dll.ForEachReverse(func(value int) bool {
		values = append(values, value)
		return true
	})

	chetest.RequireEqual(t, values, []int{3, 2, 1})
}

func TestForEachReverse_EarlyExit(t *testing.T) {
	dll := New[int]()
	dll.Append(1)
	dll.Append(2)
	dll.Append(3)

	count := 0
	dll.ForEachReverse(func(value int) bool {
		count++
		return count < 2
	})

	chetest.RequireEqual(t, count, 2)
}

func TestFind(t *testing.T) {
	dll := New[int]()
	dll.Append(1)
	dll.Append(2)
	dll.Append(3)

	value, found := dll.Find(func(v int) bool {
		return v == 2
	})

	chetest.RequireEqual(t, found, true)
	chetest.RequireEqual(t, value, 2)
}

func TestFind_NotFound(t *testing.T) {
	dll := New[int]()
	dll.Append(1)
	dll.Append(2)

	value, found := dll.Find(func(v int) bool {
		return v == 10
	})

	chetest.RequireEqual(t, found, false)
	chetest.RequireEqual(t, value, 0)
}

func TestContains(t *testing.T) {
	dll := New[int]()
	dll.Append(1)
	dll.Append(2)
	dll.Append(3)

	contains := dll.Contains(func(v int) bool {
		return v == 2
	})

	chetest.RequireEqual(t, contains, true)

	contains = dll.Contains(func(v int) bool {
		return v == 10
	})

	chetest.RequireEqual(t, contains, false)
}

func TestReverse(t *testing.T) {
	dll := New[int]()
	dll.Append(1)
	dll.Append(2)
	dll.Append(3)

	dll.Reverse()

	chetest.RequireEqual(t, dll.ToSlice(), []int{3, 2, 1})
}

func TestReverse_Empty(t *testing.T) {
	dll := New[int]()

	dll.Reverse()

	chetest.RequireEqual(t, dll.IsEmpty(), true)
}

func TestReverse_SingleElement(t *testing.T) {
	dll := New[int]()
	dll.Append(1)

	dll.Reverse()

	chetest.RequireEqual(t, dll.ToSlice(), []int{1})
}

func TestClone(t *testing.T) {
	dll := New[int]()
	dll.Append(1)
	dll.Append(2)
	dll.Append(3)

	clone := dll.Clone()

	chetest.RequireEqual(t, clone.ToSlice(), []int{1, 2, 3})

	// Modify clone
	clone.Append(4)

	// Original should be unchanged
	chetest.RequireEqual(t, dll.ToSlice(), []int{1, 2, 3})
	chetest.RequireEqual(t, clone.ToSlice(), []int{1, 2, 3, 4})
}

func TestMixedOperations(t *testing.T) {
	dll := New[string]()

	dll.Append("a")
	dll.Prepend("b")
	dll.Append("c")
	dll.InsertAt(1, "d")

	chetest.RequireEqual(t, dll.ToSlice(), []string{"b", "d", "a", "c"})

	value, _ := dll.RemoveAt(2)
	chetest.RequireEqual(t, value, "a")
	chetest.RequireEqual(t, dll.ToSlice(), []string{"b", "d", "c"})
}

func TestBidirectionalTraversal(t *testing.T) {
	dll := New[int]()
	dll.Append(1)
	dll.Append(2)
	dll.Append(3)

	// Forward traversal
	forward := dll.ToSlice()
	chetest.RequireEqual(t, forward, []int{1, 2, 3})

	// Backward traversal
	backward := dll.ToSliceReverse()
	chetest.RequireEqual(t, backward, []int{3, 2, 1})
}
