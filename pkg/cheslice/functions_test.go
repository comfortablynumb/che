package cheslice_test

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/comfortablynumb/che/pkg/cheslice"
)

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

			if !reflect.DeepEqual(result, c.expected) {
				t.Errorf("Intersect(%v) == %v, expected %v", c.input, result, c.expected)
			}
		})
	}
}

func TestContains(t *testing.T) {
	cases := []struct {
		sliceToCheck []any
		value        any
		expected     bool
	}{
		{[]any{2, 3}, 1, false},
		{[]any{2, 3, 1, 5, 1, 1, 1}, 1, true},
	}

	for i, c := range cases {
		t.Run(fmt.Sprintf("TestContains_Case-%d", i), func(t *testing.T) {
			if res := cheslice.Contains(c.sliceToCheck, c.value); res != c.expected {
				t.Errorf("Contains(%v, %v) == %v, expected %v", c.sliceToCheck, res, c.value, c.expected)
			}
		})
	}
}
