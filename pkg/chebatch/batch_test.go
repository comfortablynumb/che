package chebatch

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/comfortablynumb/che/pkg/chetest"
)

func TestBatcher_MaxSize(t *testing.T) {
	var processed [][]int
	var mu sync.Mutex

	b := NewBatcher(func(ctx context.Context, items []int) error {
		mu.Lock()
		defer mu.Unlock()
		batch := make([]int, len(items))
		copy(batch, items)
		processed = append(processed, batch)
		return nil
	}, WithMaxSize[int](3))
	defer func() { _ = b.Close() }()

	// Add 3 items - should trigger processing
	_ = b.Add(1)
	_ = b.Add(2)
	_ = b.Add(3)

	// Give processor time to run
	time.Sleep(50 * time.Millisecond)

	mu.Lock()
	chetest.RequireEqual(t, len(processed), 1)
	chetest.RequireEqual(t, len(processed[0]), 3)
	mu.Unlock()
}

func TestBatcher_MaxWait(t *testing.T) {
	var processed [][]int
	var mu sync.Mutex

	b := NewBatcher(func(ctx context.Context, items []int) error {
		mu.Lock()
		defer mu.Unlock()
		batch := make([]int, len(items))
		copy(batch, items)
		processed = append(processed, batch)
		return nil
	}, WithMaxSize[int](10), WithMaxWait[int](100*time.Millisecond))
	defer func() { _ = b.Close() }()

	_ = b.Add(1)
	_ = b.Add(2)

	// Should not be processed yet
	mu.Lock()
	chetest.RequireEqual(t, len(processed), 0)
	mu.Unlock()

	// Wait for max wait time
	time.Sleep(150 * time.Millisecond)

	mu.Lock()
	chetest.RequireEqual(t, len(processed), 1)
	chetest.RequireEqual(t, len(processed[0]), 2)
	mu.Unlock()
}

func TestBatcher_Flush(t *testing.T) {
	var processed [][]int
	var mu sync.Mutex

	b := NewBatcher(func(ctx context.Context, items []int) error {
		mu.Lock()
		defer mu.Unlock()
		batch := make([]int, len(items))
		copy(batch, items)
		processed = append(processed, batch)
		return nil
	}, WithMaxSize[int](10))
	defer func() { _ = b.Close() }()

	_ = b.Add(1)
	_ = b.Add(2)
	_ = b.Add(3)

	err := b.Flush()
	chetest.RequireEqual(t, err, nil)

	time.Sleep(50 * time.Millisecond)

	mu.Lock()
	chetest.RequireEqual(t, len(processed), 1)
	chetest.RequireEqual(t, len(processed[0]), 3)
	mu.Unlock()
}

func TestBatcher_Close(t *testing.T) {
	var processed [][]int
	var mu sync.Mutex

	b := NewBatcher(func(ctx context.Context, items []int) error {
		mu.Lock()
		defer mu.Unlock()
		batch := make([]int, len(items))
		copy(batch, items)
		processed = append(processed, batch)
		return nil
	})

	_ = b.Add(1)
	_ = b.Add(2)

	err := b.Close()
	chetest.RequireEqual(t, err, nil)

	time.Sleep(50 * time.Millisecond)

	mu.Lock()
	chetest.RequireEqual(t, len(processed), 1)
	chetest.RequireEqual(t, len(processed[0]), 2)
	mu.Unlock()
}

func TestBatcher_Size(t *testing.T) {
	b := NewBatcher(func(ctx context.Context, items []int) error {
		return nil
	}, WithMaxSize[int](10))
	defer func() { _ = b.Close() }()

	chetest.RequireEqual(t, b.Size(), 0)

	_ = b.Add(1)
	chetest.RequireEqual(t, b.Size(), 1)

	_ = b.Add(2)
	chetest.RequireEqual(t, b.Size(), 2)
}

func TestBatcher_EmptyFlush(t *testing.T) {
	callCount := 0
	var mu sync.Mutex

	b := NewBatcher(func(ctx context.Context, items []int) error {
		mu.Lock()
		callCount++
		mu.Unlock()
		return nil
	})
	defer func() { _ = b.Close() }()

	err := b.Flush()
	chetest.RequireEqual(t, err, nil)

	time.Sleep(50 * time.Millisecond)

	mu.Lock()
	chetest.RequireEqual(t, callCount, 0)
	mu.Unlock()
}

func TestGroup(t *testing.T) {
	items := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

	batches := Group(items, 3)
	chetest.RequireEqual(t, len(batches), 4)
	chetest.RequireEqual(t, len(batches[0]), 3)
	chetest.RequireEqual(t, len(batches[1]), 3)
	chetest.RequireEqual(t, len(batches[2]), 3)
	chetest.RequireEqual(t, len(batches[3]), 1)

	chetest.RequireEqual(t, batches[0][0], 1)
	chetest.RequireEqual(t, batches[3][0], 10)
}

func TestGroup_ExactDivision(t *testing.T) {
	items := []int{1, 2, 3, 4, 5, 6}

	batches := Group(items, 2)
	chetest.RequireEqual(t, len(batches), 3)
	chetest.RequireEqual(t, len(batches[0]), 2)
	chetest.RequireEqual(t, len(batches[1]), 2)
	chetest.RequireEqual(t, len(batches[2]), 2)
}

func TestGroup_InvalidSize(t *testing.T) {
	items := []int{1, 2, 3}

	batches := Group(items, 0)
	chetest.RequireEqual(t, len(batches), 0)

	batches = Group(items, -1)
	chetest.RequireEqual(t, len(batches), 0)
}

func TestGroup_EmptySlice(t *testing.T) {
	items := []int{}

	batches := Group(items, 10)
	chetest.RequireEqual(t, len(batches), 0)
}

func TestProcess(t *testing.T) {
	items := []int{1, 2, 3, 4, 5}
	var processed [][]int
	var mu sync.Mutex

	ctx := context.Background()
	err := Process(ctx, items, 2, func(ctx context.Context, batch []int) error {
		mu.Lock()
		defer mu.Unlock()
		b := make([]int, len(batch))
		copy(b, batch)
		processed = append(processed, b)
		return nil
	})

	chetest.RequireEqual(t, err, nil)

	mu.Lock()
	chetest.RequireEqual(t, len(processed), 3)
	chetest.RequireEqual(t, len(processed[0]), 2)
	chetest.RequireEqual(t, len(processed[1]), 2)
	chetest.RequireEqual(t, len(processed[2]), 1)
	mu.Unlock()
}

func TestProcess_WithError(t *testing.T) {
	items := []int{1, 2, 3, 4, 5}
	expectedErr := errors.New("processing error")

	ctx := context.Background()
	err := Process(ctx, items, 2, func(ctx context.Context, batch []int) error {
		if batch[0] == 3 {
			return expectedErr
		}
		return nil
	})

	chetest.RequireEqual(t, err, expectedErr)
}

func TestProcess_ContextCancellation(t *testing.T) {
	items := []int{1, 2, 3, 4, 5}

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	err := Process(ctx, items, 2, func(ctx context.Context, batch []int) error {
		return nil
	})

	chetest.RequireEqual(t, err, context.Canceled)
}

func TestProcessParallel(t *testing.T) {
	items := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	var processed [][]int
	var mu sync.Mutex

	ctx := context.Background()
	err := ProcessParallel(ctx, items, 2, 3, func(ctx context.Context, batch []int) error {
		time.Sleep(10 * time.Millisecond)
		mu.Lock()
		defer mu.Unlock()
		b := make([]int, len(batch))
		copy(b, batch)
		processed = append(processed, b)
		return nil
	})

	chetest.RequireEqual(t, err, nil)

	mu.Lock()
	chetest.RequireEqual(t, len(processed), 5)
	mu.Unlock()
}

func TestProcessParallel_WithError(t *testing.T) {
	items := []int{1, 2, 3, 4, 5, 6}
	expectedErr := errors.New("processing error")
	var mu sync.Mutex
	callCount := 0

	ctx := context.Background()
	err := ProcessParallel(ctx, items, 1, 2, func(ctx context.Context, batch []int) error {
		mu.Lock()
		callCount++
		mu.Unlock()

		if batch[0] == 3 {
			return expectedErr
		}
		time.Sleep(100 * time.Millisecond)
		return nil
	})

	// Should get an error (either the processing error or context canceled)
	chetest.RequireEqual(t, err != nil, true)

	// Should stop processing after error
	mu.Lock()
	chetest.RequireEqual(t, callCount < 6, true)
	mu.Unlock()
}

func TestProcessParallel_ContextCancellation(t *testing.T) {
	items := []int{1, 2, 3, 4, 5}

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	err := ProcessParallel(ctx, items, 1, 2, func(ctx context.Context, batch []int) error {
		time.Sleep(100 * time.Millisecond)
		return nil
	})

	chetest.RequireEqual(t, err != nil, true)
}

func TestBatcher_ConcurrentAdd(t *testing.T) {
	var processed [][]int
	var mu sync.Mutex

	b := NewBatcher(func(ctx context.Context, items []int) error {
		mu.Lock()
		defer mu.Unlock()
		batch := make([]int, len(items))
		copy(batch, items)
		processed = append(processed, batch)
		return nil
	}, WithMaxSize[int](100))
	defer func() { _ = b.Close() }()

	var wg sync.WaitGroup
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func(val int) {
			defer wg.Done()
			_ = b.Add(val)
		}(i)
	}

	wg.Wait()
	_ = b.Flush()
	time.Sleep(50 * time.Millisecond)

	mu.Lock()
	totalItems := 0
	for _, batch := range processed {
		totalItems += len(batch)
	}
	chetest.RequireEqual(t, totalItems, 50)
	mu.Unlock()
}

func TestBatcher_MultipleFlushes(t *testing.T) {
	var processed [][]int
	var mu sync.Mutex

	b := NewBatcher(func(ctx context.Context, items []int) error {
		mu.Lock()
		defer mu.Unlock()
		batch := make([]int, len(items))
		copy(batch, items)
		processed = append(processed, batch)
		return nil
	})
	defer func() { _ = b.Close() }()

	_ = b.Add(1)
	_ = b.Flush()

	_ = b.Add(2)
	_ = b.Flush()

	_ = b.Add(3)
	_ = b.Flush()

	time.Sleep(50 * time.Millisecond)

	mu.Lock()
	chetest.RequireEqual(t, len(processed), 3)
	mu.Unlock()
}

func TestBatcher_AddAfterClose(t *testing.T) {
	b := NewBatcher(func(ctx context.Context, items []int) error {
		return nil
	})

	_ = b.Add(1)
	err := b.Close()
	chetest.RequireEqual(t, err, nil)

	// Adding after close should return error
	err = b.Add(2)
	chetest.RequireEqual(t, err != nil, true)
}

func BenchmarkBatcher_Add(b *testing.B) {
	batcher := NewBatcher(func(ctx context.Context, items []int) error {
		return nil
	})
	defer func() { _ = batcher.Close() }()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = batcher.Add(i)
	}
}

func BenchmarkGroup(b *testing.B) {
	items := make([]int, 1000)
	for i := range items {
		items[i] = i
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = Group(items, 10)
	}
}
