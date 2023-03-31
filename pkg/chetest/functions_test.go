package chetest_test

import (
	"fmt"
	"github.com/comfortablynumb/che/pkg/chetest"
	"testing"
)

// Structs

type SimpleTestingMock struct {
	HasError bool
}

func (s *SimpleTestingMock) Errorf(_ string, _ ...any) {
	s.HasError = true
}

// Tests

func TestRequireEqual(t *testing.T) {
	cases := []struct {
		arg1          any
		arg2          any
		expectedError bool
		message       string
		messageArgs   []any
	}{
		{
			[]any{1, 2, 3},
			[]any{1, 2, 3},
			false,
			"",
			[]any{},
		},
		{
			[]any{"some", "string", "some", "other", "string"},
			[]any{"some", "string", "other"},
			true,
			"Some Message Without Args",
			[]any{},
		},
		{
			[]any{1, 1, 2, 1, 3, 2, 2, 3, 3, 2},
			[]any{1, 2, 3},
			true,
			"Some Message With Args - %v",
			[]any{1},
		},
	}

	for i, c := range cases {
		t.Run(fmt.Sprintf("TestRequireEqual_Case-%d", i), func(t *testing.T) {
			simpleTestingMock := &SimpleTestingMock{}
			opts := make([]chetest.TestOption, 0)

			if c.message != "" {
				opts = append(opts, chetest.WithExtraMessage(c.message, c.messageArgs...))
			}

			chetest.RequireEqual(simpleTestingMock, c.arg1, c.arg2, opts...)

			if simpleTestingMock.HasError != c.expectedError {
				t.Errorf(
					"Test Failed - Arg 1: %v - Arg 2: %v - Had Error: %v - Expected: %v",
					c.arg1,
					c.arg2,
					simpleTestingMock.HasError,
					c.expectedError,
				)
			}
		})
	}
}
