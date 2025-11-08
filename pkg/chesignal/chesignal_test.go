package chesignal

import (
	"context"
	"errors"
	"os"
	"testing"
	"time"

	"github.com/comfortablynumb/che/pkg/chetest"
)

func TestDefaultConfig(t *testing.T) {
	config := DefaultConfig()

	chetest.RequireEqual(t, config != nil, true)
	chetest.RequireEqual(t, len(config.Signals) > 0, true)
	chetest.RequireEqual(t, config.Timeout, 30*time.Second)
}

func TestWaitForShutdownWithContext_ImmediateCancel(t *testing.T) {
	config := &Config{
		Timeout: 5 * time.Second,
	}

	called := false
	shutdownFunc := func(ctx context.Context) error {
		called = true
		return nil
	}

	// Create a context that's already canceled
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	err := WaitForShutdownWithContext(ctx, config, shutdownFunc)

	chetest.RequireEqual(t, err, nil)
	chetest.RequireEqual(t, called, true)
}

func TestWaitForShutdownWithContext_ContextCancelation(t *testing.T) {
	config := &Config{
		Timeout: 5 * time.Second,
	}

	called := false
	shutdownFunc := func(ctx context.Context) error {
		called = true
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	err := WaitForShutdownWithContext(ctx, config, shutdownFunc)

	chetest.RequireEqual(t, err, nil)
	chetest.RequireEqual(t, called, true)
}

func TestWaitForShutdownWithContext_ShutdownFuncError(t *testing.T) {
	config := &Config{
		Timeout: 5 * time.Second,
	}

	expectedErr := errors.New("shutdown failed")
	shutdownFunc := func(ctx context.Context) error {
		return expectedErr
	}

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	err := WaitForShutdownWithContext(ctx, config, shutdownFunc)

	chetest.RequireEqual(t, err, expectedErr)
}

func TestWaitForShutdownWithContext_MultipleShutdownFuncs(t *testing.T) {
	config := &Config{
		Timeout: 5 * time.Second,
	}

	order := []int{}

	shutdownFunc1 := func(ctx context.Context) error {
		order = append(order, 1)
		return nil
	}

	shutdownFunc2 := func(ctx context.Context) error {
		order = append(order, 2)
		return nil
	}

	shutdownFunc3 := func(ctx context.Context) error {
		order = append(order, 3)
		return nil
	}

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	err := WaitForShutdownWithContext(ctx, config, shutdownFunc1, shutdownFunc2, shutdownFunc3)

	chetest.RequireEqual(t, err, nil)
	chetest.RequireEqual(t, len(order), 3)
	chetest.RequireEqual(t, order[0], 1)
	chetest.RequireEqual(t, order[1], 2)
	chetest.RequireEqual(t, order[2], 3)
}

func TestWaitForShutdownWithContext_ShutdownTimeout(t *testing.T) {
	config := &Config{
		Timeout: 50 * time.Millisecond,
	}

	timeoutCalled := false
	config.OnShutdownTimeout = func() {
		timeoutCalled = true
	}

	shutdownFunc := func(ctx context.Context) error {
		// Sleep longer than timeout
		time.Sleep(200 * time.Millisecond)
		return nil
	}

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	err := WaitForShutdownWithContext(ctx, config, shutdownFunc)

	chetest.RequireEqual(t, err != nil, true)
	chetest.RequireEqual(t, timeoutCalled, true)
}

func TestWaitForShutdownWithContext_Callbacks(t *testing.T) {
	config := &Config{
		Timeout: 5 * time.Second,
	}

	startCalled := false
	completeCalled := false

	config.OnShutdownStart = func() {
		startCalled = true
	}

	config.OnShutdownComplete = func() {
		completeCalled = true
	}

	shutdownFunc := func(ctx context.Context) error {
		return nil
	}

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	err := WaitForShutdownWithContext(ctx, config, shutdownFunc)

	chetest.RequireEqual(t, err, nil)
	chetest.RequireEqual(t, startCalled, true)
	chetest.RequireEqual(t, completeCalled, true)
}

func TestWaitForShutdownWithContext_StopOnFirstError(t *testing.T) {
	config := &Config{
		Timeout: 5 * time.Second,
	}

	called := []int{}
	expectedErr := errors.New("second function failed")

	shutdownFunc1 := func(ctx context.Context) error {
		called = append(called, 1)
		return nil
	}

	shutdownFunc2 := func(ctx context.Context) error {
		called = append(called, 2)
		return expectedErr
	}

	shutdownFunc3 := func(ctx context.Context) error {
		called = append(called, 3)
		return nil
	}

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	err := WaitForShutdownWithContext(ctx, config, shutdownFunc1, shutdownFunc2, shutdownFunc3)

	// Should stop at the second function
	chetest.RequireEqual(t, err, expectedErr)
	chetest.RequireEqual(t, len(called), 2)
	chetest.RequireEqual(t, called[0], 1)
	chetest.RequireEqual(t, called[1], 2)
}

func TestWaitForShutdownWithContext_NilConfig(t *testing.T) {
	called := false
	shutdownFunc := func(ctx context.Context) error {
		called = true
		return nil
	}

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	// Should use default config
	err := WaitForShutdownWithContext(ctx, nil, shutdownFunc)

	chetest.RequireEqual(t, err, nil)
	chetest.RequireEqual(t, called, true)
}

func TestWaitForShutdownWithContext_EmptySignals(t *testing.T) {
	config := &Config{
		Signals: []os.Signal{}, // Empty signals
		Timeout: 5 * time.Second,
	}

	called := false
	shutdownFunc := func(ctx context.Context) error {
		called = true
		return nil
	}

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	// Should use default signals
	err := WaitForShutdownWithContext(ctx, config, shutdownFunc)

	chetest.RequireEqual(t, err, nil)
	chetest.RequireEqual(t, called, true)
}

func TestNotifyOnSignal(t *testing.T) {
	// Test that we can create a signal channel
	sigChan := NotifyOnSignal()
	chetest.RequireEqual(t, sigChan != nil, true)

	// Note: We can't easily test actual signal delivery without sending real signals
	// which would be problematic in a test environment
}

func TestExecuteShutdownFuncs_Success(t *testing.T) {
	order := []int{}

	funcs := []ShutdownFunc{
		func(ctx context.Context) error {
			order = append(order, 1)
			return nil
		},
		func(ctx context.Context) error {
			order = append(order, 2)
			return nil
		},
	}

	ctx := context.Background()
	err := executeShutdownFuncs(ctx, funcs)

	chetest.RequireEqual(t, err, nil)
	chetest.RequireEqual(t, len(order), 2)
	chetest.RequireEqual(t, order[0], 1)
	chetest.RequireEqual(t, order[1], 2)
}

func TestExecuteShutdownFuncs_Error(t *testing.T) {
	expectedErr := errors.New("test error")
	order := []int{}

	funcs := []ShutdownFunc{
		func(ctx context.Context) error {
			order = append(order, 1)
			return nil
		},
		func(ctx context.Context) error {
			order = append(order, 2)
			return expectedErr
		},
		func(ctx context.Context) error {
			order = append(order, 3)
			return nil
		},
	}

	ctx := context.Background()
	err := executeShutdownFuncs(ctx, funcs)

	chetest.RequireEqual(t, err, expectedErr)
	chetest.RequireEqual(t, len(order), 2) // Should stop at error
}
