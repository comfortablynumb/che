package cheratelimit

import (
	"context"
	"sync"
	"testing"
	"time"
)

func TestLimiter_Allow(t *testing.T) {
	// 10 requests per second, burst of 5
	limiter := New(10, 5)

	// Should allow up to burst size immediately
	for i := 0; i < 5; i++ {
		if !limiter.Allow() {
			t.Errorf("request %d should be allowed within burst", i)
		}
	}

	// Next request should be denied (all tokens consumed)
	if limiter.Allow() {
		t.Error("request beyond burst should be denied")
	}
}

func TestLimiter_Refill(t *testing.T) {
	// 10 requests per second, burst of 1
	limiter := New(10, 1)

	// Consume the token
	if !limiter.Allow() {
		t.Error("first request should be allowed")
	}

	// Should be denied immediately after
	if limiter.Allow() {
		t.Error("second request should be denied")
	}

	// Wait for refill (100ms = 1 token at 10 req/s)
	time.Sleep(110 * time.Millisecond)

	// Should be allowed now
	if !limiter.Allow() {
		t.Error("request should be allowed after refill")
	}
}

func TestLimiter_Wait(t *testing.T) {
	limiter := New(10, 1)

	// Consume the token
	limiter.Allow()

	// Wait should block and then succeed
	start := time.Now()
	ctx := context.Background()
	err := limiter.Wait(ctx)

	if err != nil {
		t.Errorf("Wait should succeed: %v", err)
	}

	elapsed := time.Since(start)
	if elapsed < 50*time.Millisecond {
		t.Errorf("Wait should have blocked for ~100ms, only waited %v", elapsed)
	}
}

func TestLimiter_WaitWithContext(t *testing.T) {
	limiter := New(10, 1)

	// Consume the token
	limiter.Allow()

	// Context with short timeout
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	err := limiter.Wait(ctx)
	if err == nil {
		t.Error("Wait should fail with context timeout")
	}

	if err != context.DeadlineExceeded {
		t.Errorf("expected DeadlineExceeded, got %v", err)
	}
}

func TestLimiter_Reserve(t *testing.T) {
	limiter := New(10, 2)

	// First reserve should return 0 (token available)
	waitTime := limiter.Reserve()
	if waitTime != 0 {
		t.Errorf("first reserve should return 0, got %v", waitTime)
	}

	// Second reserve should return 0 (still have tokens)
	waitTime = limiter.Reserve()
	if waitTime != 0 {
		t.Errorf("second reserve should return 0, got %v", waitTime)
	}

	// Third reserve should return wait time
	waitTime = limiter.Reserve()
	if waitTime == 0 {
		t.Error("third reserve should return non-zero wait time")
	}

	// Wait time should be approximately 100ms (1 token at 10 req/s)
	expectedWait := 100 * time.Millisecond
	if waitTime < expectedWait-10*time.Millisecond || waitTime > expectedWait+10*time.Millisecond {
		t.Errorf("expected wait time around %v, got %v", expectedWait, waitTime)
	}
}

func TestLimiter_Tokens(t *testing.T) {
	limiter := New(10, 5)

	// Should start with full burst
	tokens := limiter.Tokens()
	if tokens != 5.0 {
		t.Errorf("expected 5 tokens, got %f", tokens)
	}

	// Consume some tokens
	limiter.Allow()
	limiter.Allow()

	tokens = limiter.Tokens()
	if tokens < 2.99 || tokens > 3.01 {
		t.Errorf("expected ~3 tokens after 2 allows, got %f", tokens)
	}
}

func TestLimiter_SetRate(t *testing.T) {
	limiter := New(10, 1)

	// Consume token
	limiter.Allow()

	// Change rate to 100 req/s (faster refill)
	limiter.SetRate(100)

	// Should refill much faster now
	time.Sleep(15 * time.Millisecond)

	if !limiter.Allow() {
		t.Error("request should be allowed with faster rate")
	}
}

func TestLimiter_SetBurst(t *testing.T) {
	limiter := New(10, 5)

	// Change burst to 2
	limiter.SetBurst(2)

	tokens := limiter.Tokens()
	if tokens > 2.0 {
		t.Errorf("tokens should be capped at new burst size, got %f", tokens)
	}

	// Change burst to 10
	limiter.SetBurst(10)

	// Should allow more tokens to accumulate
	time.Sleep(200 * time.Millisecond)
	tokens = limiter.Tokens()
	if tokens > 10.0 {
		t.Errorf("tokens should be capped at 10, got %f", tokens)
	}
}

func TestLimiter_Concurrent(t *testing.T) {
	limiter := New(100, 10)
	var wg sync.WaitGroup

	allowed := 0
	denied := 0
	var mu sync.Mutex

	// Spawn many goroutines trying to acquire tokens
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if limiter.Allow() {
				mu.Lock()
				allowed++
				mu.Unlock()
			} else {
				mu.Lock()
				denied++
				mu.Unlock()
			}
		}()
	}

	wg.Wait()

	// Should allow exactly burst size
	if allowed != 10 {
		t.Errorf("expected 10 allowed, got %d", allowed)
	}

	if denied != 90 {
		t.Errorf("expected 90 denied, got %d", denied)
	}
}

func TestPerKeyLimiter_Allow(t *testing.T) {
	pkl := NewPerKey(10, 1)

	// Different keys should have independent limits
	if !pkl.Allow("key1") {
		t.Error("first request for key1 should be allowed")
	}

	if !pkl.Allow("key2") {
		t.Error("first request for key2 should be allowed")
	}

	// Same key should be limited
	if pkl.Allow("key1") {
		t.Error("second request for key1 should be denied")
	}

	// Different key should still work
	if !pkl.Allow("key3") {
		t.Error("first request for key3 should be allowed")
	}
}

func TestPerKeyLimiter_Wait(t *testing.T) {
	pkl := NewPerKey(10, 1)

	// Consume token for key1
	pkl.Allow("key1")

	// Wait should succeed
	ctx := context.Background()
	err := pkl.Wait(ctx, "key1")

	if err != nil {
		t.Errorf("Wait should succeed: %v", err)
	}
}

func TestPerKeyLimiter_Cleanup(t *testing.T) {
	// Short cleanup interval for testing
	pkl := &PerKeyLimiter{
		rate:    10,
		burst:   1,
		cleanup: 100 * time.Millisecond,
	}

	// Start cleanup goroutine
	go pkl.cleanupLoop()

	// Create limiters
	pkl.Allow("key1")
	pkl.Allow("key2")

	// Check they exist
	if _, ok := pkl.limiters.Load("key1"); !ok {
		t.Error("key1 should exist")
	}

	// Wait for cleanup
	time.Sleep(250 * time.Millisecond)

	// Old limiters should be cleaned up (or might have already been removed)
	// This test just verifies the cleanup goroutine runs without panicking
}

func TestPerKeyLimiter_ConcurrentKeys(t *testing.T) {
	// Use a very low rate (0.01 tokens/sec) to avoid refill during the test
	// At 0.01 tokens/sec, it takes 100 seconds to get 1 token, so refill won't affect the burst
	pkl := NewPerKey(0.01, 10)
	var wg sync.WaitGroup

	allowed := make(map[string]int)
	var mu sync.Mutex

	keys := []string{"key1", "key2", "key3"}

	// Each key should have independent limits
	for _, key := range keys {
		for i := 0; i < 20; i++ {
			wg.Add(1)
			go func(k string) {
				defer wg.Done()
				if pkl.Allow(k) {
					mu.Lock()
					allowed[k]++
					mu.Unlock()
				}
			}(key)
		}
	}

	wg.Wait()

	// Each key should allow up to burst size
	for _, key := range keys {
		if allowed[key] != 10 {
			t.Errorf("key %s: expected 10 allowed, got %d", key, allowed[key])
		}
	}
}

func BenchmarkLimiter_Allow(b *testing.B) {
	limiter := New(1000000, 1000)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		limiter.Allow()
	}
}

func BenchmarkLimiter_AllowParallel(b *testing.B) {
	limiter := New(1000000, 1000)

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			limiter.Allow()
		}
	})
}

func BenchmarkPerKeyLimiter_Allow(b *testing.B) {
	pkl := NewPerKey(1000000, 1000)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		pkl.Allow("benchmark-key")
	}
}
