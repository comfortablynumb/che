package chesemaphore

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/comfortablynumb/che/pkg/chetest"
)

func TestSemaphore_Acquire(t *testing.T) {
	sem := New(10)

	ctx := context.Background()
	err := sem.Acquire(ctx, 5)
	chetest.RequireEqual(t, err, nil)
	chetest.RequireEqual(t, sem.Current(), int64(5))
	chetest.RequireEqual(t, sem.Available(), int64(5))
}

func TestSemaphore_AcquireRelease(t *testing.T) {
	sem := New(10)

	ctx := context.Background()
	err := sem.Acquire(ctx, 5)
	chetest.RequireEqual(t, err, nil)
	chetest.RequireEqual(t, sem.Current(), int64(5))

	sem.Release(5)
	chetest.RequireEqual(t, sem.Current(), int64(0))
	chetest.RequireEqual(t, sem.Available(), int64(10))
}

func TestSemaphore_TryAcquire(t *testing.T) {
	sem := New(10)

	// Should succeed
	ok := sem.TryAcquire(5)
	chetest.RequireEqual(t, ok, true)
	chetest.RequireEqual(t, sem.Current(), int64(5))

	// Should succeed
	ok = sem.TryAcquire(5)
	chetest.RequireEqual(t, ok, true)
	chetest.RequireEqual(t, sem.Current(), int64(10))

	// Should fail (no capacity)
	ok = sem.TryAcquire(1)
	chetest.RequireEqual(t, ok, false)
	chetest.RequireEqual(t, sem.Current(), int64(10))
}

func TestSemaphore_WeightExceedsLimit(t *testing.T) {
	sem := New(10)

	ctx := context.Background()
	err := sem.Acquire(ctx, 11)
	chetest.RequireEqual(t, err, ErrWeightExceedsLimit)

	ok := sem.TryAcquire(11)
	chetest.RequireEqual(t, ok, false)
}

func TestSemaphore_ConcurrentAcquire(t *testing.T) {
	sem := New(10)
	var wg sync.WaitGroup

	// Launch 10 goroutines, each acquiring 1
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			ctx := context.Background()
			err := sem.Acquire(ctx, 1)
			chetest.RequireEqual(t, err, nil)
		}()
	}

	wg.Wait()
	chetest.RequireEqual(t, sem.Current(), int64(10))
	chetest.RequireEqual(t, sem.Available(), int64(0))
}

func TestSemaphore_BlockingAcquire(t *testing.T) {
	sem := New(5)

	ctx := context.Background()
	err := sem.Acquire(ctx, 5)
	chetest.RequireEqual(t, err, nil)

	acquired := false
	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		ctx := context.Background()
		err := sem.Acquire(ctx, 3)
		chetest.RequireEqual(t, err, nil)
		acquired = true
	}()

	// Give goroutine time to block
	time.Sleep(50 * time.Millisecond)
	chetest.RequireEqual(t, acquired, false)

	// Release enough resources
	sem.Release(3)

	wg.Wait()
	chetest.RequireEqual(t, acquired, true)
}

func TestSemaphore_ContextCancellation(t *testing.T) {
	sem := New(5)

	// Acquire all resources
	ctx := context.Background()
	err := sem.Acquire(ctx, 5)
	chetest.RequireEqual(t, err, nil)

	// Try to acquire with cancelled context
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	err = sem.Acquire(ctx, 1)
	chetest.RequireEqual(t, err, context.Canceled)
}

func TestSemaphore_ContextTimeout(t *testing.T) {
	sem := New(5)

	// Acquire all resources
	ctx := context.Background()
	err := sem.Acquire(ctx, 5)
	chetest.RequireEqual(t, err, nil)

	// Try to acquire with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	err = sem.Acquire(ctx, 1)
	chetest.RequireEqual(t, err, context.DeadlineExceeded)
}

func TestSemaphore_MultipleWeights(t *testing.T) {
	sem := New(10)

	ctx := context.Background()

	err := sem.Acquire(ctx, 3)
	chetest.RequireEqual(t, err, nil)
	chetest.RequireEqual(t, sem.Available(), int64(7))

	err = sem.Acquire(ctx, 4)
	chetest.RequireEqual(t, err, nil)
	chetest.RequireEqual(t, sem.Available(), int64(3))

	err = sem.Acquire(ctx, 3)
	chetest.RequireEqual(t, err, nil)
	chetest.RequireEqual(t, sem.Available(), int64(0))
}

func TestSemaphore_ReleaseMoreThanAcquired(t *testing.T) {
	sem := New(10)

	ctx := context.Background()
	err := sem.Acquire(ctx, 5)
	chetest.RequireEqual(t, err, nil)

	// Release more than acquired should reset to 0
	sem.Release(10)
	chetest.RequireEqual(t, sem.Current(), int64(0))
	chetest.RequireEqual(t, sem.Available(), int64(10))
}

func TestSemaphore_Size(t *testing.T) {
	sem := New(100)
	chetest.RequireEqual(t, sem.Size(), int64(100))
}

func TestSemaphore_ConcurrentAcquireRelease(t *testing.T) {
	sem := New(100)
	var wg sync.WaitGroup

	// Launch multiple goroutines that acquire and release
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			ctx := context.Background()

			err := sem.Acquire(ctx, 10)
			chetest.RequireEqual(t, err, nil)

			time.Sleep(10 * time.Millisecond)
			sem.Release(10)
		}()
	}

	wg.Wait()

	// All resources should be released
	chetest.RequireEqual(t, sem.Current(), int64(0))
	chetest.RequireEqual(t, sem.Available(), int64(100))
}

func TestSemaphore_TryAcquirePartial(t *testing.T) {
	sem := New(10)

	ok := sem.TryAcquire(7)
	chetest.RequireEqual(t, ok, true)

	ok = sem.TryAcquire(4)
	chetest.RequireEqual(t, ok, false)

	ok = sem.TryAcquire(3)
	chetest.RequireEqual(t, ok, true)

	chetest.RequireEqual(t, sem.Current(), int64(10))
}

func TestSemaphore_StressTest(t *testing.T) {
	sem := New(50)
	var wg sync.WaitGroup
	iterations := 1000

	for i := 0; i < iterations; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()

			weight := int64((id % 5) + 1)
			ctx := context.Background()

			err := sem.Acquire(ctx, weight)
			if err != nil {
				return
			}

			// Simulate some work
			time.Sleep(time.Microsecond)

			sem.Release(weight)
		}(i)
	}

	wg.Wait()

	// Should be back to initial state
	chetest.RequireEqual(t, sem.Current(), int64(0))
	chetest.RequireEqual(t, sem.Available(), int64(50))
}

func BenchmarkSemaphore_Acquire(b *testing.B) {
	sem := New(int64(b.N))
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = sem.Acquire(ctx, 1)
	}
}

func BenchmarkSemaphore_TryAcquire(b *testing.B) {
	sem := New(int64(b.N))

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sem.TryAcquire(1)
	}
}

func BenchmarkSemaphore_AcquireRelease(b *testing.B) {
	sem := New(100)
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = sem.Acquire(ctx, 1)
		sem.Release(1)
	}
}
