package cheslice_test

import (
	"fmt"
	"github.com/comfortablynumb/che/pkg/cheslice"
	"reflect"
	"testing"
)

func TestIntersect(t *testing.T) {
	cases := []struct {
		input    [][]interface{}
		expected []interface{}
	}{
		{[][]interface{}{{1, 2, 3}, {2, 3, 4}}, []interface{}{2, 3}},
		{[][]interface{}{{1, 2, 3}, {3, 4, 5}}, []interface{}{3}},
		{[][]interface{}{{1, 2, 3}}, []interface{}{1, 2, 3}},
		{[][]interface{}{{1, 2, 3, 1, 2, 3, 1, 2, 3, 1, 2, 3}}, []interface{}{1, 2, 3}},
		{[][]interface{}{{"something", 2, "hi"}}, []interface{}{"something", 2, "hi"}},
		{[][]interface{}{}, []interface{}{}},
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
		sliceToCheck []interface{}
		value        interface{}
		expected     bool
	}{
		{[]interface{}{2, 3}, 1, false},
		{[]interface{}{2, 3, 1, 5, 1, 1, 1}, 1, true},
	}

	for i, c := range cases {
		t.Run(fmt.Sprintf("TestIntersect_Case-%d", i), func(t *testing.T) {
			if res := cheslice.Contains(c.sliceToCheck, c.value); res != c.expected {
				t.Errorf("Contains(%v, %v) == %v, expected %v", c.sliceToCheck, res, c.value, c.expected)
			}
		})

	}
}
