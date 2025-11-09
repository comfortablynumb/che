// Package chebatch provides utilities for batch processing and aggregation.
package chebatch

import (
	"context"
	"sync"
	"time"
)

// Processor is a function that processes a batch of items.
type Processor[T any] func(ctx context.Context, items []T) error

// Batcher collects items and processes them in batches.
type Batcher[T any] struct {
	maxSize   int
	maxWait   time.Duration
	processor Processor[T]

	mu      sync.Mutex
	items   []T
	timer   *time.Timer
	ctx     context.Context
	cancel  context.CancelFunc
	wg      sync.WaitGroup
	started bool
}

// BatcherOption is a configuration option for Batcher.
type BatcherOption[T any] func(*Batcher[T])

// WithMaxSize sets the maximum batch size.
func WithMaxSize[T any](size int) BatcherOption[T] {
	return func(b *Batcher[T]) {
		b.maxSize = size
	}
}

// WithMaxWait sets the maximum wait time before flushing.
func WithMaxWait[T any](duration time.Duration) BatcherOption[T] {
	return func(b *Batcher[T]) {
		b.maxWait = duration
	}
}

// NewBatcher creates a new batcher with the given processor.
// Default maxSize is 100, maxWait is 1 second.
func NewBatcher[T any](processor Processor[T], opts ...BatcherOption[T]) *Batcher[T] {
	ctx, cancel := context.WithCancel(context.Background())

	b := &Batcher[T]{
		maxSize:   100,
		maxWait:   1 * time.Second,
		processor: processor,
		items:     make([]T, 0),
		ctx:       ctx,
		cancel:    cancel,
	}

	for _, opt := range opts {
		opt(b)
	}

	return b
}

// Start starts the batcher's background processing.
func (b *Batcher[T]) Start() {
	b.mu.Lock()
	defer b.mu.Unlock()

	if b.started {
		return
	}

	b.started = true
}

// Add adds an item to the batch.
// Triggers processing if batch size limit is reached.
func (b *Batcher[T]) Add(item T) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	if !b.started {
		b.started = true
	}

	select {
	case <-b.ctx.Done():
		return b.ctx.Err()
	default:
	}

	b.items = append(b.items, item)

	// Reset timer
	if b.timer != nil {
		b.timer.Stop()
	}

	if len(b.items) >= b.maxSize {
		// Process immediately
		_ = b.processLocked()
	} else {
		// Set timer for max wait
		b.timer = time.AfterFunc(b.maxWait, func() {
			b.mu.Lock()
			defer b.mu.Unlock()
			_ = b.processLocked()
		})
	}

	return nil
}

// Flush processes all pending items immediately.
func (b *Batcher[T]) Flush() error {
	b.mu.Lock()
	defer b.mu.Unlock()

	if b.timer != nil {
		b.timer.Stop()
		b.timer = nil
	}

	return b.processLocked()
}

// Close flushes any pending items and closes the batcher.
func (b *Batcher[T]) Close() error {
	b.mu.Lock()

	if b.timer != nil {
		b.timer.Stop()
		b.timer = nil
	}

	err := b.processLocked()

	b.cancel()
	b.mu.Unlock()

	b.wg.Wait()

	return err
}

func (b *Batcher[T]) processLocked() error {
	if len(b.items) == 0 {
		return nil
	}

	items := make([]T, len(b.items))
	copy(items, b.items)
	b.items = b.items[:0]

	b.wg.Add(1)
	go func() {
		defer b.wg.Done()
		_ = b.processor(b.ctx, items)
	}()

	return nil
}

// Size returns the current number of pending items.
func (b *Batcher[T]) Size() int {
	b.mu.Lock()
	defer b.mu.Unlock()
	return len(b.items)
}

// Group groups items into batches of specified size.
func Group[T any](items []T, size int) [][]T {
	if size <= 0 {
		return [][]T{}
	}

	var batches [][]T
	for i := 0; i < len(items); i += size {
		end := i + size
		if end > len(items) {
			end = len(items)
		}
		batches = append(batches, items[i:end])
	}

	return batches
}

// Process processes items in batches using the provided function.
func Process[T any](ctx context.Context, items []T, batchSize int, fn Processor[T]) error {
	batches := Group(items, batchSize)

	for _, batch := range batches {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		if err := fn(ctx, batch); err != nil {
			return err
		}
	}

	return nil
}

// ProcessParallel processes items in parallel batches with limited concurrency.
func ProcessParallel[T any](ctx context.Context, items []T, batchSize int, maxConcurrent int, fn Processor[T]) error {
	batches := Group(items, batchSize)

	sem := make(chan struct{}, maxConcurrent)
	errCh := make(chan error, 1)
	var wg sync.WaitGroup

	processCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	for _, batch := range batches {
		select {
		case <-processCtx.Done():
			wg.Wait()
			return processCtx.Err()
		case err := <-errCh:
			cancel()
			wg.Wait()
			return err
		default:
		}

		wg.Add(1)
		sem <- struct{}{}

		go func(b []T) {
			defer wg.Done()
			defer func() { <-sem }()

			if err := fn(processCtx, b); err != nil {
				select {
				case errCh <- err:
				default:
				}
				cancel()
			}
		}(batch)
	}

	wg.Wait()

	select {
	case err := <-errCh:
		return err
	default:
		return nil
	}
}
