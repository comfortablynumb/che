package chemap_test

import (
	"fmt"
	"github.com/comfortablynumb/che/pkg/chemap"
	"github.com/comfortablynumb/che/pkg/chetest"
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

			chetest.RequireEqual(t, result, c.expectedKeys)
		})
	}
}
