package cheworker

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestNew_DefaultConfig(t *testing.T) {
	pool := New(nil)

	if pool.WorkerCount() != 10 {
		t.Errorf("expected 10 workers, got %d", pool.WorkerCount())
	}

	if pool.QueueSize() != 100 {
		t.Errorf("expected queue size 100, got %d", pool.QueueSize())
	}
}

func TestNew_CustomConfig(t *testing.T) {
	config := &Config{
		Workers:   5,
		QueueSize: 50,
	}

	pool := New(config)

	if pool.WorkerCount() != 5 {
		t.Errorf("expected 5 workers, got %d", pool.WorkerCount())
	}

	if pool.QueueSize() != 50 {
		t.Errorf("expected queue size 50, got %d", pool.QueueSize())
	}
}

func TestPool_SubmitAndExecute(t *testing.T) {
	var counter atomic.Int32
	pool := New(&Config{Workers: 2})
	pool.Start()

	// Submit jobs
	for i := 0; i < 10; i++ {
		err := pool.Submit(func(ctx context.Context) error {
			counter.Add(1)
			return nil
		})
		if err != nil {
			t.Errorf("failed to submit job: %v", err)
		}
	}

	pool.Shutdown()

	if counter.Load() != 10 {
		t.Errorf("expected counter to be 10, got %d", counter.Load())
	}
}

func TestPool_ErrorHandling(t *testing.T) {
	expectedErr := errors.New("job error")
	pool := New(&Config{Workers: 2})
	pool.Start()

	// Submit job that returns error
	err := pool.Submit(func(ctx context.Context) error {
		return expectedErr
	})
	if err != nil {
		t.Fatal(err)
	}

	pool.Shutdown()

	// Give a moment for error collector to process
	time.Sleep(10 * time.Millisecond)

	errs := pool.Errors()
	if len(errs) != 1 {
		t.Fatalf("expected 1 error, got %d", len(errs))
	}

	if errs[0] != expectedErr {
		t.Errorf("expected error %v, got %v", expectedErr, errs[0])
	}
}

func TestPool_OnErrorCallback(t *testing.T) {
	var callbackErrors []error
	var mu sync.Mutex

	config := &Config{
		Workers: 2,
		OnError: func(err error) {
			mu.Lock()
			callbackErrors = append(callbackErrors, err)
			mu.Unlock()
		},
	}

	pool := New(config)
	pool.Start()

	expectedErr := errors.New("test error")
	err := pool.Submit(func(ctx context.Context) error {
		return expectedErr
	})
	if err != nil {
		t.Fatal(err)
	}

	pool.Shutdown()

	// Give a moment for error collector to process
	time.Sleep(10 * time.Millisecond)

	mu.Lock()
	defer mu.Unlock()

	if len(callbackErrors) != 1 {
		t.Errorf("expected 1 callback error, got %d", len(callbackErrors))
	}
}

func TestPool_PanicRecovery(t *testing.T) {
	pool := New(&Config{Workers: 2})
	pool.Start()

	err := pool.Submit(func(ctx context.Context) error {
		panic("test panic")
	})
	if err != nil {
		t.Fatal(err)
	}

	pool.Shutdown()

	errs := pool.Errors()
	if len(errs) != 1 {
		t.Fatalf("expected 1 error from panic, got %d", len(errs))
	}

	if errs[0] == nil || errs[0].Error() != "worker panic: test panic" {
		t.Errorf("unexpected panic error: %v", errs[0])
	}
}

func TestPool_CustomPanicHandler(t *testing.T) {
	var panicValue interface{}
	var mu sync.Mutex

	config := &Config{
		Workers: 2,
		PanicHandler: func(p interface{}) {
			mu.Lock()
			panicValue = p
			mu.Unlock()
		},
	}

	pool := New(config)
	pool.Start()

	err := pool.Submit(func(ctx context.Context) error {
		panic("custom panic")
	})
	if err != nil {
		t.Fatal(err)
	}

	pool.Shutdown()

	mu.Lock()
	defer mu.Unlock()

	if panicValue != "custom panic" {
		t.Errorf("expected panic value 'custom panic', got %v", panicValue)
	}
}

func TestPool_ContextCancellation(t *testing.T) {
	pool := New(&Config{Workers: 2})
	pool.Start()

	var started atomic.Bool
	var cancelled atomic.Bool

	err := pool.Submit(func(ctx context.Context) error {
		started.Store(true)
		<-ctx.Done()
		cancelled.Store(true)
		return ctx.Err()
	})
	if err != nil {
		t.Fatal(err)
	}

	// Wait for job to start
	time.Sleep(50 * time.Millisecond)

	pool.Stop()

	if !started.Load() {
		t.Error("job never started")
	}

	if !cancelled.Load() {
		t.Error("job was not cancelled")
	}
}

func TestPool_SubmitWithContext(t *testing.T) {
	pool := New(&Config{Workers: 2})
	pool.Start()

	ctx, cancel := context.WithCancel(context.Background())
	var jobCancelled atomic.Bool

	err := pool.SubmitWithContext(ctx, func(jobCtx context.Context) error {
		<-jobCtx.Done()
		jobCancelled.Store(true)
		return jobCtx.Err()
	})
	if err != nil {
		t.Fatal(err)
	}

	time.Sleep(50 * time.Millisecond)
	cancel()

	pool.Shutdown()

	if !jobCancelled.Load() {
		t.Error("job was not cancelled")
	}
}

func TestPool_ShutdownWithContext(t *testing.T) {
	pool := New(&Config{Workers: 2})
	pool.Start()

	// Submit a slow job
	var completed atomic.Bool
	err := pool.Submit(func(ctx context.Context) error {
		time.Sleep(100 * time.Millisecond)
		completed.Store(true)
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}

	// Shutdown with generous timeout
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()

	err = pool.ShutdownWithContext(ctx)
	if err != nil {
		t.Errorf("unexpected shutdown error: %v", err)
	}

	if !completed.Load() {
		t.Error("job did not complete")
	}
}

func TestPool_ShutdownWithContext_Timeout(t *testing.T) {
	pool := New(&Config{Workers: 2})
	pool.Start()

	// Submit a slow job
	err := pool.Submit(func(ctx context.Context) error {
		select {
		case <-time.After(500 * time.Millisecond):
		case <-ctx.Done():
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}

	// Shutdown with short timeout
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	err = pool.ShutdownWithContext(ctx)
	if err != context.DeadlineExceeded {
		t.Errorf("expected DeadlineExceeded, got %v", err)
	}
}

func TestPool_SubmitAfterShutdown(t *testing.T) {
	pool := New(&Config{Workers: 2})
	pool.Start()
	pool.Shutdown()

	err := pool.Submit(func(ctx context.Context) error {
		return nil
	})

	if err == nil {
		t.Error("expected error when submitting after shutdown")
	}
}

func TestPool_ConcurrentSubmit(t *testing.T) {
	var counter atomic.Int32
	pool := New(&Config{Workers: 10, QueueSize: 100})
	pool.Start()

	var wg sync.WaitGroup
	jobCount := 100

	for i := 0; i < jobCount; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			err := pool.Submit(func(ctx context.Context) error {
				counter.Add(1)
				return nil
			})
			if err != nil {
				t.Errorf("failed to submit: %v", err)
			}
		}()
	}

	wg.Wait()
	pool.Shutdown()

	if counter.Load() != int32(jobCount) {
		t.Errorf("expected %d jobs executed, got %d", jobCount, counter.Load())
	}
}

func TestPool_PendingJobs(t *testing.T) {
	pool := New(&Config{Workers: 1, QueueSize: 10})
	pool.Start()

	// Submit jobs that will block the worker
	blocker := make(chan struct{})
	err := pool.Submit(func(ctx context.Context) error {
		<-blocker
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}

	// Submit more jobs
	for i := 0; i < 5; i++ {
		err := pool.Submit(func(ctx context.Context) error {
			return nil
		})
		if err != nil {
			t.Fatal(err)
		}
	}

	// Wait for worker to start processing
	time.Sleep(50 * time.Millisecond)

	pending := pool.PendingJobs()
	if pending != 5 {
		t.Errorf("expected 5 pending jobs, got %d", pending)
	}

	close(blocker)
	pool.Shutdown()
}

func BenchmarkPool_Submit(b *testing.B) {
	pool := New(&Config{Workers: 10, QueueSize: 1000})
	pool.Start()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err := pool.Submit(func(ctx context.Context) error {
			return nil
		})
		if err != nil {
			b.Fatal(err)
		}
	}

	pool.Shutdown()
}

func BenchmarkPool_SubmitAndWait(b *testing.B) {
	pool := New(&Config{Workers: 10, QueueSize: 100})
	pool.Start()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err := pool.Submit(func(ctx context.Context) error {
			// Simulate some work
			time.Sleep(time.Microsecond)
			return nil
		})
		if err != nil {
			b.Fatal(err)
		}
	}

	pool.Shutdown()
}

func ExamplePool() {
	// Create a worker pool with 5 workers
	pool := New(&Config{
		Workers:   5,
		QueueSize: 100,
	})

	// Start the pool
	pool.Start()

	// Submit jobs
	for i := 0; i < 10; i++ {
		id := i
		err := pool.Submit(func(ctx context.Context) error {
			fmt.Printf("Processing job %d\n", id)
			return nil
		})
		if err != nil {
			fmt.Printf("Failed to submit job: %v\n", err)
		}
	}

	// Gracefully shutdown
	pool.Shutdown()
}
