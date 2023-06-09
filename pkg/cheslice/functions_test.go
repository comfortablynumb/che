package cheslice_test

import (
	"fmt"
	"github.com/comfortablynumb/che/pkg/chetest"
	"testing"

	"github.com/comfortablynumb/che/pkg/cheslice"
)

func TestUnion(t *testing.T) {
	cases := []struct {
		input    [][]any
		expected []any
	}{
		{
			[][]any{{1, 2, 3}, {4, 5, 6}, {7, 8, 9}},
			[]any{1, 2, 3, 4, 5, 6, 7, 8, 9},
		},
		{
			[][]any{{1, 2, 3}},
			[]any{1, 2, 3},
		},
		{
			[][]any{},
			[]any{},
		},
	}

	for i, c := range cases {
		t.Run(fmt.Sprintf("TestUnion_Case-%d", i), func(t *testing.T) {
			result := cheslice.Union(c.input...)

			chetest.RequireEqual(t, result, c.expected)
		})
	}
}

func TestForEach(t *testing.T) {
	type ForEachTestHelper struct {
		processed  []int
		iterations int
	}

	cases := []struct {
		input    []int
		expected *ForEachTestHelper
	}{
		{
			[]int{1, 2, 3, 4},
			&ForEachTestHelper{
				processed:  []int{2, 4},
				iterations: 4,
			},
		},
		{
			[]int{1, 2, 3, 4, 5, 0, 3, 1, 2, 3, 4},
			&ForEachTestHelper{
				processed:  []int{2, 4},
				iterations: 5,
			},
		},
	}

	for i, c := range cases {
		t.Run(fmt.Sprintf("TestForEach_Case-%d", i), func(t *testing.T) {
			forEachTestHelper := &ForEachTestHelper{
				processed:  []int{},
				iterations: 0,
			}

			cheslice.ForEach(c.input, func(element int) bool {
				if element == 0 {
					return false
				}

				if (element % 2) == 0 {
					forEachTestHelper.processed = append(forEachTestHelper.processed, element)
				}

				forEachTestHelper.iterations++

				return true
			})

			chetest.RequireEqual(t, forEachTestHelper, c.expected)
		})
	}
}

func TestMap(t *testing.T) {
	cases := []struct {
		input    []any
		mapFunc  cheslice.MapFunc[any]
		expected []any
	}{
		{
			[]any{1, 2, 3},
			func(element any) any {
				return element.(int) + 2
			},
			[]any{3, 4, 5},
		},
		{
			[]any{},
			func(element any) any {
				return element.(int) + 2
			},
			[]any{},
		},
	}

	for i, c := range cases {
		t.Run(fmt.Sprintf("TestMap_Case-%d", i), func(t *testing.T) {
			inputCopy := make([]any, 0, len(c.input))

			inputCopy = append(inputCopy, c.input...)

			result := cheslice.Map(c.input, c.mapFunc)

			chetest.RequireEqual(t, result, c.expected)

			// Make sure original slice is untouched

			chetest.RequireEqual(t, inputCopy, c.input)
		})
	}
}

func TestFilter(t *testing.T) {
	cases := []struct {
		input      []any
		filterFunc cheslice.FilterFunc[any]
		expected   []any
	}{
		{
			[]any{1, 2, 3, 4, 5, 6},
			func(element any) bool {
				return (element.(int) % 2) == 0
			},
			[]any{2, 4, 6},
		},
		{
			[]any{1, 2, 3, 4, 5, 6},
			func(element any) bool {
				return false
			},
			[]any{},
		},
		{
			[]any{1, 2, 3, 4, 5, 6},
			func(element any) bool {
				return true
			},
			[]any{1, 2, 3, 4, 5, 6},
		},
		{
			[]any{},
			func(element any) bool {
				return true
			},
			[]any{},
		},
	}

	for i, c := range cases {
		t.Run(fmt.Sprintf("TestFilter_Case-%d", i), func(t *testing.T) {
			result := cheslice.Filter(c.input, c.filterFunc)

			chetest.RequireEqual(t, result, c.expected)
		})
	}
}

func TestFill(t *testing.T) {
	cases := []struct {
		count    uint
		value    any
		expected []any
	}{
		{0, nil, []any{}},
		{1, nil, []any{nil}},
		{5, 100, []any{100, 100, 100, 100, 100}},
	}

	for i, c := range cases {
		t.Run(fmt.Sprintf("TestFill_Case-%d", i), func(t *testing.T) {
			result := cheslice.Fill(c.count, c.value)

			chetest.RequireEqual(t, result, c.expected)
		})
	}
}

func TestDiff(t *testing.T) {
	cases := []struct {
		input    [][]any
		expected []any
	}{
		{[][]any{}, []any{}},
		{[][]any{{1, 2, 3, 3, 4}}, []any{1, 2, 3, 3, 4}},
		{[][]any{{1, 2, 3, 3, 4}, {4, 2, 3, 3, 2, 4, 1, 2, 5, 3, 7}}, []any{}},
		{[][]any{{1, 2, 3, 3, 4}, {2, 3, 3, 2, 1, 2, 5, 3, 7}}, []any{4}},
		{[][]any{{1, 2, 3, 3, 4}, {5, 6, 7, 32, 45, 234, 654, 3453342}}, []any{1, 2, 3, 4}},
	}

	for i, c := range cases {
		t.Run(fmt.Sprintf("TestDiff_Case-%d", i), func(t *testing.T) {
			result := cheslice.Diff(c.input...)

			chetest.RequireEqual(t, result, c.expected)
		})
	}
}

func TestChunk(t *testing.T) {
	cases := []struct {
		input    []any
		length   uint
		expected [][]any
	}{
		{[]any{1, 2, 3}, 5, [][]any{{1, 2, 3}}},
		{[]any{1, 2, 3}, 2, [][]any{{1, 2}, {3}}},
		{[]any{1, 2, 3, 4, 3, 2, 1}, 2, [][]any{{1, 2}, {3, 4}, {3, 2}, {1}}},
		{[]any{1, 2, 3, 4, 3, 2}, 2, [][]any{{1, 2}, {3, 4}, {3, 2}}},
		{[]any{1, 2, 3, 4, 3, 2, 1}, 0, [][]any{}},
	}

	for i, c := range cases {
		t.Run(fmt.Sprintf("TestChunk_Case-%d", i), func(t *testing.T) {
			inputCopy := make([]any, 0, len(c.input))
			inputCopy = append(inputCopy, c.input...)

			result := cheslice.Chunk(c.input, c.length)

			chetest.RequireEqual(t, result, c.expected)

			// Confirm the original slice was not modified

			chetest.RequireEqual(t, c.input, inputCopy)
		})
	}
}

func TestUnique(t *testing.T) {
	cases := []struct {
		input    []any
		expected []any
	}{
		{[]any{1, 2, 3}, []any{1, 2, 3}},
		{[]any{1, 1, 2, 1, 3, 2, 2, 3, 3, 2}, []any{1, 2, 3}},
	}

	for i, c := range cases {
		t.Run(fmt.Sprintf("TestUnique_Case-%d", i), func(t *testing.T) {
			result := cheslice.Unique(c.input)

			chetest.RequireEqual(t, result, c.expected)
		})
	}
}

func TestIntersect(t *testing.T) {
	cases := []struct {
		input    [][]any
		expected []any
	}{
		{[][]any{{1, 2, 3}, {2, 3, 4}}, []any{2, 3}},
		{[][]any{{1, 2, 3}, {3, 4, 5}}, []any{3}},
		{[][]any{{1, 2, 3}}, []any{1, 2, 3}},
		{[][]any{{1, 2, 3, 1, 2, 3, 1, 2, 3, 1, 2, 3}}, []any{1, 2, 3, 1, 2, 3, 1, 2, 3, 1, 2, 3}},
		{[][]any{
			{1, 2, 3, 1, 2, 3, 1, 100, 5, 6, 12, 2, 3, 1, 2, 3, 87},
			{1, 2, 5, 6, 12, 3, 1, 2, 3, 1, 2, 107, 3, 1, 2, 3},
		}, []any{1, 2, 3, 5, 6, 12}},
		{[][]any{{"something", 2, "hi"}}, []any{"something", 2, "hi"}},
		{[][]any{}, []any{}},
		{[][]any{{1, 2, 3}}, []any{1, 2, 3}},
	}

	for i, c := range cases {
		t.Run(fmt.Sprintf("TestIntersect_Case-%d", i), func(t *testing.T) {
			result := cheslice.Intersect(c.input...)

			chetest.RequireEqual(t, result, c.expected)
		})
	}
}

func TestExists(t *testing.T) {
	cases := []struct {
		slicesToCheck [][]any
		value         any
		expected      bool
	}{
		{[][]any{{2, 3}}, 1, false},
		{[][]any{{4, 5, 6, 3}, {}, {2, 3, 1, 5, 1, 1, 1}}, 1, true},
	}

	for i, c := range cases {
		t.Run(fmt.Sprintf("TestExists_Case-%d", i), func(t *testing.T) {
			result := cheslice.Exists(c.value, c.slicesToCheck...)

			chetest.RequireEqual(t, result, c.expected)
		})
	}
}

func TestLen(t *testing.T) {
	cases := []struct {
		slices   [][]any
		expected int
	}{
		{[][]any{{1, 2, 3}, {}, {3, 4, 5}}, 6},
		{[][]any{}, 0},
		{[][]any{{}}, 0},
		{nil, 0},
	}

	for i, c := range cases {
		t.Run(fmt.Sprintf("TestLen_Case-%d", i), func(t *testing.T) {
			result := cheslice.Len(c.slices...)

			chetest.RequireEqual(t, result, c.expected)
		})
	}
}
