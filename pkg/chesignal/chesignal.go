package chesignal

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// ShutdownFunc is a function that performs cleanup on shutdown
// It receives a context that will be canceled when the shutdown timeout expires
type ShutdownFunc func(ctx context.Context) error

// Config configures the graceful shutdown behavior
type Config struct {
	// Signals to listen for (defaults to SIGINT and SIGTERM)
	Signals []os.Signal

	// Timeout for graceful shutdown (defaults to 30 seconds)
	Timeout time.Duration

	// OnShutdownStart is called when shutdown begins (optional)
	OnShutdownStart func()

	// OnShutdownComplete is called when shutdown completes successfully (optional)
	OnShutdownComplete func()

	// OnShutdownTimeout is called if shutdown times out (optional)
	OnShutdownTimeout func()
}

// DefaultConfig returns a sensible default configuration
func DefaultConfig() *Config {
	return &Config{
		Signals: []os.Signal{syscall.SIGINT, syscall.SIGTERM},
		Timeout: 30 * time.Second,
	}
}

// WaitForShutdown blocks until a shutdown signal is received, then executes shutdown functions
// Returns an error if any shutdown function fails or if timeout is exceeded
func WaitForShutdown(config *Config, shutdownFuncs ...ShutdownFunc) error {
	if config == nil {
		config = DefaultConfig()
	}

	// Set defaults
	if len(config.Signals) == 0 {
		config.Signals = []os.Signal{syscall.SIGINT, syscall.SIGTERM}
	}
	if config.Timeout == 0 {
		config.Timeout = 30 * time.Second
	}

	// Create signal channel
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, config.Signals...)

	// Wait for signal
	<-sigChan

	// Call onShutdownStart callback
	if config.OnShutdownStart != nil {
		config.OnShutdownStart()
	}

	// Create context with timeout for shutdown
	ctx, cancel := context.WithTimeout(context.Background(), config.Timeout)
	defer cancel()

	// Execute shutdown functions
	errChan := make(chan error, 1)
	go func() {
		errChan <- executeShutdownFuncs(ctx, shutdownFuncs)
	}()

	// Wait for shutdown to complete or timeout
	select {
	case err := <-errChan:
		if err != nil {
			return err
		}
		if config.OnShutdownComplete != nil {
			config.OnShutdownComplete()
		}
		return nil
	case <-ctx.Done():
		if config.OnShutdownTimeout != nil {
			config.OnShutdownTimeout()
		}
		return ctx.Err()
	}
}

// executeShutdownFuncs executes all shutdown functions in order
func executeShutdownFuncs(ctx context.Context, funcs []ShutdownFunc) error {
	for _, fn := range funcs {
		if err := fn(ctx); err != nil {
			return err
		}
	}
	return nil
}

// WaitForShutdownWithContext is similar to WaitForShutdown but also listens to context cancellation
// This is useful when you want to trigger shutdown programmatically
func WaitForShutdownWithContext(ctx context.Context, config *Config, shutdownFuncs ...ShutdownFunc) error {
	if config == nil {
		config = DefaultConfig()
	}

	// Set defaults
	if len(config.Signals) == 0 {
		config.Signals = []os.Signal{syscall.SIGINT, syscall.SIGTERM}
	}
	if config.Timeout == 0 {
		config.Timeout = 30 * time.Second
	}

	// Create signal channel
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, config.Signals...)
	defer signal.Stop(sigChan)

	// Wait for signal or context cancellation
	select {
	case <-sigChan:
	case <-ctx.Done():
	}

	// Call onShutdownStart callback
	if config.OnShutdownStart != nil {
		config.OnShutdownStart()
	}

	// Create context with timeout for shutdown
	shutdownCtx, cancel := context.WithTimeout(context.Background(), config.Timeout)
	defer cancel()

	// Execute shutdown functions
	errChan := make(chan error, 1)
	go func() {
		errChan <- executeShutdownFuncs(shutdownCtx, shutdownFuncs)
	}()

	// Wait for shutdown to complete or timeout
	select {
	case err := <-errChan:
		if err != nil {
			return err
		}
		if config.OnShutdownComplete != nil {
			config.OnShutdownComplete()
		}
		return nil
	case <-shutdownCtx.Done():
		if config.OnShutdownTimeout != nil {
			config.OnShutdownTimeout()
		}
		return shutdownCtx.Err()
	}
}

// NotifyOnSignal creates a channel that receives a value when any of the specified signals are received
// This is useful for custom signal handling logic
func NotifyOnSignal(signals ...os.Signal) <-chan os.Signal {
	if len(signals) == 0 {
		signals = []os.Signal{syscall.SIGINT, syscall.SIGTERM}
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, signals...)
	return sigChan
}
