package chelinkedlist

import (
	"testing"

	"github.com/comfortablynumb/che/pkg/chetest"
)

func TestNew(t *testing.T) {
	ll := New[int]()

	chetest.RequireEqual(t, ll.IsEmpty(), true)
	chetest.RequireEqual(t, ll.Size(), 0)
}

func TestPrepend(t *testing.T) {
	ll := New[int]()

	ll.Prepend(3)
	ll.Prepend(2)
	ll.Prepend(1)

	chetest.RequireEqual(t, ll.Size(), 3)
	chetest.RequireEqual(t, ll.ToSlice(), []int{1, 2, 3})
}

func TestAppend(t *testing.T) {
	ll := New[int]()

	ll.Append(1)
	ll.Append(2)
	ll.Append(3)

	chetest.RequireEqual(t, ll.Size(), 3)
	chetest.RequireEqual(t, ll.ToSlice(), []int{1, 2, 3})
}

func TestInsertAt(t *testing.T) {
	ll := New[int]()

	ll.Append(1)
	ll.Append(3)

	ok := ll.InsertAt(1, 2)
	chetest.RequireEqual(t, ok, true)
	chetest.RequireEqual(t, ll.ToSlice(), []int{1, 2, 3})

	ok = ll.InsertAt(0, 0)
	chetest.RequireEqual(t, ok, true)
	chetest.RequireEqual(t, ll.ToSlice(), []int{0, 1, 2, 3})

	ok = ll.InsertAt(4, 4)
	chetest.RequireEqual(t, ok, true)
	chetest.RequireEqual(t, ll.ToSlice(), []int{0, 1, 2, 3, 4})
}

func TestInsertAt_OutOfBounds(t *testing.T) {
	ll := New[int]()
	ll.Append(1)

	ok := ll.InsertAt(-1, 0)
	chetest.RequireEqual(t, ok, false)

	ok = ll.InsertAt(10, 0)
	chetest.RequireEqual(t, ok, false)
}

func TestRemoveFirst(t *testing.T) {
	ll := New[int]()
	ll.Append(1)
	ll.Append(2)
	ll.Append(3)

	value, ok := ll.RemoveFirst()
	chetest.RequireEqual(t, ok, true)
	chetest.RequireEqual(t, value, 1)
	chetest.RequireEqual(t, ll.Size(), 2)
	chetest.RequireEqual(t, ll.ToSlice(), []int{2, 3})
}

func TestRemoveFirst_Empty(t *testing.T) {
	ll := New[int]()

	value, ok := ll.RemoveFirst()
	chetest.RequireEqual(t, ok, false)
	chetest.RequireEqual(t, value, 0)
}

func TestRemoveFirst_SingleElement(t *testing.T) {
	ll := New[int]()
	ll.Append(1)

	value, ok := ll.RemoveFirst()
	chetest.RequireEqual(t, ok, true)
	chetest.RequireEqual(t, value, 1)
	chetest.RequireEqual(t, ll.IsEmpty(), true)
}

func TestRemoveLast(t *testing.T) {
	ll := New[int]()
	ll.Append(1)
	ll.Append(2)
	ll.Append(3)

	value, ok := ll.RemoveLast()
	chetest.RequireEqual(t, ok, true)
	chetest.RequireEqual(t, value, 3)
	chetest.RequireEqual(t, ll.Size(), 2)
	chetest.RequireEqual(t, ll.ToSlice(), []int{1, 2})
}

func TestRemoveLast_Empty(t *testing.T) {
	ll := New[int]()

	value, ok := ll.RemoveLast()
	chetest.RequireEqual(t, ok, false)
	chetest.RequireEqual(t, value, 0)
}

func TestRemoveLast_SingleElement(t *testing.T) {
	ll := New[int]()
	ll.Append(1)

	value, ok := ll.RemoveLast()
	chetest.RequireEqual(t, ok, true)
	chetest.RequireEqual(t, value, 1)
	chetest.RequireEqual(t, ll.IsEmpty(), true)
}

func TestRemoveAt(t *testing.T) {
	ll := New[int]()
	ll.Append(1)
	ll.Append(2)
	ll.Append(3)
	ll.Append(4)

	value, ok := ll.RemoveAt(1)
	chetest.RequireEqual(t, ok, true)
	chetest.RequireEqual(t, value, 2)
	chetest.RequireEqual(t, ll.ToSlice(), []int{1, 3, 4})

	value, ok = ll.RemoveAt(2)
	chetest.RequireEqual(t, ok, true)
	chetest.RequireEqual(t, value, 4)
	chetest.RequireEqual(t, ll.ToSlice(), []int{1, 3})
}

func TestRemoveAt_OutOfBounds(t *testing.T) {
	ll := New[int]()
	ll.Append(1)

	value, ok := ll.RemoveAt(-1)
	chetest.RequireEqual(t, ok, false)
	chetest.RequireEqual(t, value, 0)

	value, ok = ll.RemoveAt(10)
	chetest.RequireEqual(t, ok, false)
	chetest.RequireEqual(t, value, 0)
}

func TestGet(t *testing.T) {
	ll := New[int]()
	ll.Append(10)
	ll.Append(20)
	ll.Append(30)

	value, ok := ll.Get(0)
	chetest.RequireEqual(t, ok, true)
	chetest.RequireEqual(t, value, 10)

	value, ok = ll.Get(1)
	chetest.RequireEqual(t, ok, true)
	chetest.RequireEqual(t, value, 20)

	value, ok = ll.Get(2)
	chetest.RequireEqual(t, ok, true)
	chetest.RequireEqual(t, value, 30)
}

func TestGet_OutOfBounds(t *testing.T) {
	ll := New[int]()
	ll.Append(1)

	value, ok := ll.Get(-1)
	chetest.RequireEqual(t, ok, false)
	chetest.RequireEqual(t, value, 0)

	value, ok = ll.Get(10)
	chetest.RequireEqual(t, ok, false)
	chetest.RequireEqual(t, value, 0)
}

func TestFirst(t *testing.T) {
	ll := New[int]()
	ll.Append(10)
	ll.Append(20)

	value, ok := ll.First()
	chetest.RequireEqual(t, ok, true)
	chetest.RequireEqual(t, value, 10)
}

func TestFirst_Empty(t *testing.T) {
	ll := New[int]()

	value, ok := ll.First()
	chetest.RequireEqual(t, ok, false)
	chetest.RequireEqual(t, value, 0)
}

func TestLast(t *testing.T) {
	ll := New[int]()
	ll.Append(10)
	ll.Append(20)

	value, ok := ll.Last()
	chetest.RequireEqual(t, ok, true)
	chetest.RequireEqual(t, value, 20)
}

func TestLast_Empty(t *testing.T) {
	ll := New[int]()

	value, ok := ll.Last()
	chetest.RequireEqual(t, ok, false)
	chetest.RequireEqual(t, value, 0)
}

func TestSize(t *testing.T) {
	ll := New[int]()

	chetest.RequireEqual(t, ll.Size(), 0)

	ll.Append(1)
	chetest.RequireEqual(t, ll.Size(), 1)

	ll.Append(2)
	chetest.RequireEqual(t, ll.Size(), 2)

	ll.RemoveFirst()
	chetest.RequireEqual(t, ll.Size(), 1)
}

func TestIsEmpty(t *testing.T) {
	ll := New[int]()

	chetest.RequireEqual(t, ll.IsEmpty(), true)

	ll.Append(1)
	chetest.RequireEqual(t, ll.IsEmpty(), false)

	ll.RemoveFirst()
	chetest.RequireEqual(t, ll.IsEmpty(), true)
}

func TestClear(t *testing.T) {
	ll := New[int]()
	ll.Append(1)
	ll.Append(2)
	ll.Append(3)

	ll.Clear()

	chetest.RequireEqual(t, ll.IsEmpty(), true)
	chetest.RequireEqual(t, ll.Size(), 0)
}

func TestToSlice(t *testing.T) {
	ll := New[int]()
	ll.Append(1)
	ll.Append(2)
	ll.Append(3)

	slice := ll.ToSlice()

	chetest.RequireEqual(t, slice, []int{1, 2, 3})
}

func TestToSlice_Empty(t *testing.T) {
	ll := New[int]()

	slice := ll.ToSlice()

	chetest.RequireEqual(t, len(slice), 0)
}

func TestForEach(t *testing.T) {
	ll := New[int]()
	ll.Append(1)
	ll.Append(2)
	ll.Append(3)

	sum := 0
	ll.ForEach(func(value int) bool {
		sum += value
		return true
	})

	chetest.RequireEqual(t, sum, 6)
}

func TestForEach_EarlyExit(t *testing.T) {
	ll := New[int]()
	ll.Append(1)
	ll.Append(2)
	ll.Append(3)

	count := 0
	ll.ForEach(func(value int) bool {
		count++
		return count < 2
	})

	chetest.RequireEqual(t, count, 2)
}

func TestFind(t *testing.T) {
	ll := New[int]()
	ll.Append(1)
	ll.Append(2)
	ll.Append(3)

	value, found := ll.Find(func(v int) bool {
		return v == 2
	})

	chetest.RequireEqual(t, found, true)
	chetest.RequireEqual(t, value, 2)
}

func TestFind_NotFound(t *testing.T) {
	ll := New[int]()
	ll.Append(1)
	ll.Append(2)

	value, found := ll.Find(func(v int) bool {
		return v == 10
	})

	chetest.RequireEqual(t, found, false)
	chetest.RequireEqual(t, value, 0)
}

func TestContains(t *testing.T) {
	ll := New[int]()
	ll.Append(1)
	ll.Append(2)
	ll.Append(3)

	contains := ll.Contains(func(v int) bool {
		return v == 2
	})

	chetest.RequireEqual(t, contains, true)

	contains = ll.Contains(func(v int) bool {
		return v == 10
	})

	chetest.RequireEqual(t, contains, false)
}

func TestReverse(t *testing.T) {
	ll := New[int]()
	ll.Append(1)
	ll.Append(2)
	ll.Append(3)

	ll.Reverse()

	chetest.RequireEqual(t, ll.ToSlice(), []int{3, 2, 1})
}

func TestReverse_Empty(t *testing.T) {
	ll := New[int]()

	ll.Reverse()

	chetest.RequireEqual(t, ll.IsEmpty(), true)
}

func TestReverse_SingleElement(t *testing.T) {
	ll := New[int]()
	ll.Append(1)

	ll.Reverse()

	chetest.RequireEqual(t, ll.ToSlice(), []int{1})
}

func TestClone(t *testing.T) {
	ll := New[int]()
	ll.Append(1)
	ll.Append(2)
	ll.Append(3)

	clone := ll.Clone()

	chetest.RequireEqual(t, clone.ToSlice(), []int{1, 2, 3})

	// Modify clone
	clone.Append(4)

	// Original should be unchanged
	chetest.RequireEqual(t, ll.ToSlice(), []int{1, 2, 3})
	chetest.RequireEqual(t, clone.ToSlice(), []int{1, 2, 3, 4})
}

func TestMixedOperations(t *testing.T) {
	ll := New[string]()

	ll.Append("a")
	ll.Prepend("b")
	ll.Append("c")
	ll.InsertAt(1, "d")

	chetest.RequireEqual(t, ll.ToSlice(), []string{"b", "d", "a", "c"})

	value, _ := ll.RemoveAt(2)
	chetest.RequireEqual(t, value, "a")
	chetest.RequireEqual(t, ll.ToSlice(), []string{"b", "d", "c"})
}
