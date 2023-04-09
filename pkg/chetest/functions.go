package chetest

import (
	"fmt"
	"reflect"
)

// Interfaces

type TestingInterface interface {
	Errorf(format string, args ...any)
}

// Types

type TestOption func(options *testOptions)

// Structs

type testOptions struct {
	message     string
	messageArgs []any
}

// Functions

// RequireEqual Fails the test if "input" and "expected" are not deeply equal. Returns true otherwise. You can also
// pass an extra message to show using the "WithExtraMessage" option.
func RequireEqual[T any](t TestingInterface, input T, expected T, options ...TestOption) {
	if !reflect.DeepEqual(input, expected) {
		extraMessage := prepareExtraMessage(options...)

		t.Errorf("Test Failed - Received input: %v - Expected: %v %s", input, expected, extraMessage)
	}
}

func WithExtraMessage(message string, messageArgs ...any) TestOption {
	return func(testOptions *testOptions) {
		testOptions.message = message
		testOptions.messageArgs = messageArgs
	}
}

func prepareExtraMessage(options ...TestOption) string {
	testOpts := &testOptions{
		message:     "",
		messageArgs: []any{},
	}

	for _, option := range options {
		option(testOpts)
	}

	extraMessage := ""

	if testOpts.message != "" {
		args := []any{testOpts.message}

		args = append(args, testOpts.messageArgs...)

		extraMessage = fmt.Sprintf("%s", args...)
	}

	return extraMessage
}
