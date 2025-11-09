package cheratelimit

import (
	"context"
	"sync"
	"time"
)

// Limiter implements the token bucket rate limiting algorithm.
type Limiter struct {
	rate       float64
	burst      int
	tokens     float64
	lastUpdate time.Time
	mu         sync.Mutex
}

// New creates a new rate limiter.
// rate is the number of requests per second.
// burst is the maximum number of requests that can be made at once.
func New(rate float64, burst int) *Limiter {
	return &Limiter{
		rate:       rate,
		burst:      burst,
		tokens:     float64(burst),
		lastUpdate: time.Now(),
	}
}

// Allow returns true if a request is allowed.
func (l *Limiter) Allow() bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.refill()

	if l.tokens >= 1.0 {
		l.tokens--
		return true
	}

	return false
}

// Wait blocks until a request is allowed or the context is cancelled.
func (l *Limiter) Wait(ctx context.Context) error {
	for {
		if l.Allow() {
			return nil
		}

		// Calculate wait time
		waitTime := l.waitDuration()

		select {
		case <-time.After(waitTime):
			continue
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

// Reserve reserves a token and returns the wait duration.
func (l *Limiter) Reserve() time.Duration {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.refill()

	if l.tokens >= 1.0 {
		l.tokens--
		return 0
	}

	// Calculate how long to wait for a token
	tokensNeeded := 1.0 - l.tokens
	waitTime := time.Duration(tokensNeeded/l.rate*float64(time.Second))
	l.tokens = 0
	return waitTime
}

// Tokens returns the current number of available tokens.
func (l *Limiter) Tokens() float64 {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.refill()
	return l.tokens
}

// SetRate updates the rate limit.
func (l *Limiter) SetRate(rate float64) {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.refill()
	l.rate = rate
}

// SetBurst updates the burst size.
func (l *Limiter) SetBurst(burst int) {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.burst = burst
	if l.tokens > float64(burst) {
		l.tokens = float64(burst)
	}
}

func (l *Limiter) refill() {
	now := time.Now()
	elapsed := now.Sub(l.lastUpdate)
	l.lastUpdate = now

	// Add tokens based on elapsed time
	tokensToAdd := l.rate * elapsed.Seconds()
	l.tokens += tokensToAdd

	// Cap at burst
	if l.tokens > float64(l.burst) {
		l.tokens = float64(l.burst)
	}
}

func (l *Limiter) waitDuration() time.Duration {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.refill()

	if l.tokens >= 1.0 {
		return 0
	}

	tokensNeeded := 1.0 - l.tokens
	return time.Duration(tokensNeeded / l.rate * float64(time.Second))
}

// PerKeyLimiter is a rate limiter that maintains separate limits per key.
type PerKeyLimiter struct {
	rate     float64
	burst    int
	limiters sync.Map
	cleanup  time.Duration
	mu       sync.Mutex
}

// NewPerKey creates a new per-key rate limiter.
func NewPerKey(rate float64, burst int) *PerKeyLimiter {
	pkl := &PerKeyLimiter{
		rate:    rate,
		burst:   burst,
		cleanup: 10 * time.Minute,
	}

	// Start cleanup goroutine
	go pkl.cleanupLoop()

	return pkl
}

// Allow checks if a request for the given key is allowed.
func (pkl *PerKeyLimiter) Allow(key string) bool {
	limiter := pkl.getLimiter(key)
	return limiter.Allow()
}

// Wait blocks until a request for the given key is allowed.
func (pkl *PerKeyLimiter) Wait(ctx context.Context, key string) error {
	limiter := pkl.getLimiter(key)
	return limiter.Wait(ctx)
}

func (pkl *PerKeyLimiter) getLimiter(key string) *Limiter {
	if limiter, ok := pkl.limiters.Load(key); ok {
		return limiter.(*Limiter)
	}

	newLimiter := New(pkl.rate, pkl.burst)
	actual, _ := pkl.limiters.LoadOrStore(key, newLimiter)
	return actual.(*Limiter)
}

func (pkl *PerKeyLimiter) cleanupLoop() {
	ticker := time.NewTicker(pkl.cleanup)
	defer ticker.Stop()

	for range ticker.C {
		// Remove limiters that haven't been used recently
		pkl.limiters.Range(func(key, value interface{}) bool {
			limiter := value.(*Limiter)

			// Lock to safely read lastUpdate
			limiter.mu.Lock()
			lastUpdate := limiter.lastUpdate
			limiter.mu.Unlock()

			if time.Since(lastUpdate) > pkl.cleanup {
				pkl.limiters.Delete(key)
			}
			return true
		})
	}
}
