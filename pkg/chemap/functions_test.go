package chemap_test

import (
	"fmt"
	"github.com/comfortablynumb/che/pkg/chemap"
	"github.com/comfortablynumb/che/pkg/chetest"
	"sort"
	"testing"
)

func TestKeys(t *testing.T) {
	cases := []struct {
		theMap       map[string]struct{}
		expectedKeys []string
	}{
		{
			map[string]struct{}{
				"someKey":      {},
				"someOtherKey": {},
				"aaaa":         {},
				"bbbb":         {},
			},
			[]string{
				"someKey",
				"someOtherKey",
				"aaaa",
				"bbbb",
			},
		},
		{
			map[string]struct{}{},
			[]string{},
		},
		{
			nil,
			[]string{},
		},
	}

	for i, c := range cases {
		t.Run(fmt.Sprintf("TestKeys_Case-%d", i), func(t *testing.T) {
			result := chemap.Keys(c.theMap)

			sort.Strings(result)
			sort.Strings(c.expectedKeys)

			chetest.RequireEqual(t, result, c.expectedKeys)
		})
	}
}

func TestValues(t *testing.T) {
	cases := []struct {
		theMap         map[string]int
		expectedValues []int
	}{
		{
			map[string]int{
				"a": 1,
				"b": 2,
				"c": 3,
			},
			[]int{1, 2, 3},
		},
		{
			map[string]int{},
			[]int{},
		},
		{
			nil,
			[]int{},
		},
	}

	for i, c := range cases {
		t.Run(fmt.Sprintf("TestValues_Case-%d", i), func(t *testing.T) {
			result := chemap.Values(c.theMap)

			sort.Ints(result)
			sort.Ints(c.expectedValues)

			chetest.RequireEqual(t, result, c.expectedValues)
		})
	}
}

func TestInvert(t *testing.T) {
	m := map[string]int{
		"a": 1,
		"b": 2,
		"c": 3,
	}

	result := chemap.Invert(m)

	chetest.RequireEqual(t, len(result), 3)
	chetest.RequireEqual(t, result[1], "a")
	chetest.RequireEqual(t, result[2], "b")
	chetest.RequireEqual(t, result[3], "c")
}

func TestInvert_DuplicateValues(t *testing.T) {
	m := map[string]int{
		"a": 1,
		"b": 1,
		"c": 2,
	}

	result := chemap.Invert(m)

	// Only one of "a" or "b" will be retained for value 1
	chetest.RequireEqual(t, len(result), 2)
	chetest.RequireEqual(t, result[2], "c")
	// result[1] will be either "a" or "b"
}

func TestFilter(t *testing.T) {
	m := map[string]int{
		"a": 1,
		"b": 2,
		"c": 3,
		"d": 4,
	}

	result := chemap.Filter(m, func(k string, v int) bool {
		return v%2 == 0
	})

	chetest.RequireEqual(t, len(result), 2)
	chetest.RequireEqual(t, result["b"], 2)
	chetest.RequireEqual(t, result["d"], 4)
}

func TestFilter_Empty(t *testing.T) {
	m := map[string]int{
		"a": 1,
		"b": 2,
	}

	result := chemap.Filter(m, func(k string, v int) bool {
		return false
	})

	chetest.RequireEqual(t, len(result), 0)
}

func TestMapValues(t *testing.T) {
	m := map[string]int{
		"a": 1,
		"b": 2,
		"c": 3,
	}

	result := chemap.MapValues(m, func(v int) int {
		return v * 2
	})

	chetest.RequireEqual(t, len(result), 3)
	chetest.RequireEqual(t, result["a"], 2)
	chetest.RequireEqual(t, result["b"], 4)
	chetest.RequireEqual(t, result["c"], 6)
}

func TestMapValues_TypeConversion(t *testing.T) {
	m := map[string]int{
		"a": 1,
		"b": 2,
		"c": 3,
	}

	result := chemap.MapValues(m, func(v int) string {
		return fmt.Sprintf("value_%d", v)
	})

	chetest.RequireEqual(t, len(result), 3)
	chetest.RequireEqual(t, result["a"], "value_1")
	chetest.RequireEqual(t, result["b"], "value_2")
	chetest.RequireEqual(t, result["c"], "value_3")
}

func TestMerge(t *testing.T) {
	m1 := map[string]int{
		"a": 1,
		"b": 2,
	}
	m2 := map[string]int{
		"c": 3,
		"d": 4,
	}
	m3 := map[string]int{
		"b": 20,
		"e": 5,
	}

	result := chemap.Merge(m1, m2, m3)

	chetest.RequireEqual(t, len(result), 5)
	chetest.RequireEqual(t, result["a"], 1)
	chetest.RequireEqual(t, result["b"], 20) // Last value wins
	chetest.RequireEqual(t, result["c"], 3)
	chetest.RequireEqual(t, result["d"], 4)
	chetest.RequireEqual(t, result["e"], 5)
}

func TestMerge_Empty(t *testing.T) {
	result := chemap.Merge[string, int]()

	chetest.RequireEqual(t, len(result), 0)
}

func TestPick(t *testing.T) {
	m := map[string]int{
		"a": 1,
		"b": 2,
		"c": 3,
		"d": 4,
	}

	result := chemap.Pick(m, "a", "c", "e")

	chetest.RequireEqual(t, len(result), 2)
	chetest.RequireEqual(t, result["a"], 1)
	chetest.RequireEqual(t, result["c"], 3)
}

func TestPick_Empty(t *testing.T) {
	m := map[string]int{
		"a": 1,
		"b": 2,
	}

	result := chemap.Pick(m, "x", "y")

	chetest.RequireEqual(t, len(result), 0)
}

func TestOmit(t *testing.T) {
	m := map[string]int{
		"a": 1,
		"b": 2,
		"c": 3,
		"d": 4,
	}

	result := chemap.Omit(m, "b", "d")

	chetest.RequireEqual(t, len(result), 2)
	chetest.RequireEqual(t, result["a"], 1)
	chetest.RequireEqual(t, result["c"], 3)
}

func TestOmit_NonExistent(t *testing.T) {
	m := map[string]int{
		"a": 1,
		"b": 2,
	}

	result := chemap.Omit(m, "x", "y")

	chetest.RequireEqual(t, len(result), 2)
	chetest.RequireEqual(t, result["a"], 1)
	chetest.RequireEqual(t, result["b"], 2)
}
