// Package chesemaphore provides weighted semaphore for concurrency control.
package chesemaphore

import (
	"context"
	"errors"
	"sync"
)

var (
	// ErrWeightExceedsLimit is returned when requested weight exceeds semaphore limit.
	ErrWeightExceedsLimit = errors.New("weight exceeds semaphore limit")
)

// Semaphore is a weighted semaphore for limiting concurrent access.
type Semaphore struct {
	size    int64
	current int64
	mu      sync.Mutex
	cond    *sync.Cond
}

// New creates a new weighted semaphore with the given size.
func New(size int64) *Semaphore {
	s := &Semaphore{
		size: size,
	}
	s.cond = sync.NewCond(&s.mu)
	return s
}

// Acquire acquires the semaphore with the specified weight.
// It blocks until resources are available.
func (s *Semaphore) Acquire(ctx context.Context, weight int64) error {
	if weight > s.size {
		return ErrWeightExceedsLimit
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	// Check context before waiting
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	// Wait for resources to be available
	for s.current+weight > s.size {
		// Use a goroutine to wait for context cancellation
		done := make(chan struct{})
		go func() {
			select {
			case <-ctx.Done():
				s.cond.Broadcast()
				close(done)
			case <-done:
			}
		}()

		s.cond.Wait()

		select {
		case <-done:
		default:
			close(done)
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}
	}

	s.current += weight
	return nil
}

// TryAcquire tries to acquire the semaphore without blocking.
// Returns true if successful, false otherwise.
func (s *Semaphore) TryAcquire(weight int64) bool {
	if weight > s.size {
		return false
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if s.current+weight > s.size {
		return false
	}

	s.current += weight
	return true
}

// Release releases the semaphore with the specified weight.
func (s *Semaphore) Release(weight int64) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.current -= weight
	if s.current < 0 {
		s.current = 0
	}

	s.cond.Broadcast()
}

// Available returns the number of available resources.
func (s *Semaphore) Available() int64 {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.size - s.current
}

// Size returns the total size of the semaphore.
func (s *Semaphore) Size() int64 {
	return s.size
}

// Current returns the current usage of the semaphore.
func (s *Semaphore) Current() int64 {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.current
}
