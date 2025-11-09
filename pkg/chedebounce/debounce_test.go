package chedebounce

import (
	"sync"
	"testing"
	"time"

	"github.com/comfortablynumb/che/pkg/chetest"
)

func TestDebouncer_SingleCall(t *testing.T) {
	d := NewDebouncer(50 * time.Millisecond)
	defer d.Close()

	called := false
	var mu sync.Mutex

	d.Call(func() {
		mu.Lock()
		called = true
		mu.Unlock()
	})

	// Should not be called immediately
	mu.Lock()
	chetest.RequireEqual(t, called, false)
	mu.Unlock()

	// Should be called after delay
	time.Sleep(100 * time.Millisecond)
	mu.Lock()
	chetest.RequireEqual(t, called, true)
	mu.Unlock()
}

func TestDebouncer_MultipleCalls(t *testing.T) {
	d := NewDebouncer(50 * time.Millisecond)
	defer d.Close()

	count := 0
	var mu sync.Mutex

	// Call multiple times quickly
	for i := 0; i < 5; i++ {
		d.Call(func() {
			mu.Lock()
			count++
			mu.Unlock()
		})
		time.Sleep(10 * time.Millisecond)
	}

	// Should only execute once after all calls settle
	time.Sleep(100 * time.Millisecond)

	mu.Lock()
	chetest.RequireEqual(t, count, 1)
	mu.Unlock()
}

func TestDebouncer_Flush(t *testing.T) {
	d := NewDebouncer(500 * time.Millisecond)
	defer d.Close()

	called := false
	var mu sync.Mutex

	d.Call(func() {
		mu.Lock()
		called = true
		mu.Unlock()
	})

	// Flush should execute immediately
	d.Flush()
	mu.Lock()
	chetest.RequireEqual(t, called, true)
	mu.Unlock()
}

func TestDebouncer_Cancel(t *testing.T) {
	d := NewDebouncer(50 * time.Millisecond)
	defer d.Close()

	called := false
	var mu sync.Mutex

	d.Call(func() {
		mu.Lock()
		called = true
		mu.Unlock()
	})

	d.Cancel()

	time.Sleep(100 * time.Millisecond)
	mu.Lock()
	chetest.RequireEqual(t, called, false)
	mu.Unlock()
}

func TestDebouncer_Close(t *testing.T) {
	d := NewDebouncer(50 * time.Millisecond)

	called := false
	var mu sync.Mutex

	d.Call(func() {
		mu.Lock()
		called = true
		mu.Unlock()
	})

	d.Close()

	time.Sleep(100 * time.Millisecond)
	mu.Lock()
	chetest.RequireEqual(t, called, false)
	mu.Unlock()

	// Further calls should be ignored
	d.Call(func() {
		mu.Lock()
		called = true
		mu.Unlock()
	})
	time.Sleep(100 * time.Millisecond)
	mu.Lock()
	chetest.RequireEqual(t, called, false)
	mu.Unlock()
}

func TestThrottler_LeadingEdge(t *testing.T) {
	th := NewThrottler(100*time.Millisecond, WithLeading())
	defer th.Close()

	count := 0
	var mu sync.Mutex

	// First call should execute immediately
	th.Call(func() {
		mu.Lock()
		count++
		mu.Unlock()
	})

	mu.Lock()
	chetest.RequireEqual(t, count, 1)
	mu.Unlock()

	// Second call within interval should be ignored
	th.Call(func() {
		mu.Lock()
		count++
		mu.Unlock()
	})

	time.Sleep(50 * time.Millisecond)
	mu.Lock()
	chetest.RequireEqual(t, count, 1)
	mu.Unlock()

	// After interval, should execute again
	time.Sleep(100 * time.Millisecond)
	th.Call(func() {
		mu.Lock()
		count++
		mu.Unlock()
	})

	mu.Lock()
	chetest.RequireEqual(t, count, 2)
	mu.Unlock()
}

func TestThrottler_TrailingEdge(t *testing.T) {
	th := NewThrottler(100*time.Millisecond, WithTrailing())
	defer th.Close()

	count := 0
	var mu sync.Mutex

	// First call should be delayed
	th.Call(func() {
		mu.Lock()
		count++
		mu.Unlock()
	})

	mu.Lock()
	chetest.RequireEqual(t, count, 0)
	mu.Unlock()

	// Should execute after interval
	time.Sleep(150 * time.Millisecond)
	mu.Lock()
	chetest.RequireEqual(t, count, 1)
	mu.Unlock()
}

func TestThrottler_LeadingAndTrailing(t *testing.T) {
	th := NewThrottler(100*time.Millisecond, WithLeading(), WithTrailing())
	defer th.Close()

	count := 0
	var mu sync.Mutex

	// First call - leading edge
	th.Call(func() {
		mu.Lock()
		count++
		mu.Unlock()
	})

	mu.Lock()
	chetest.RequireEqual(t, count, 1)
	mu.Unlock()

	// Calls during throttle period
	for i := 0; i < 3; i++ {
		th.Call(func() {
			mu.Lock()
			count++
			mu.Unlock()
		})
		time.Sleep(20 * time.Millisecond)
	}

	// Should have trailing edge call
	time.Sleep(100 * time.Millisecond)
	mu.Lock()
	chetest.RequireEqual(t, count, 2)
	mu.Unlock()
}

func TestThrottler_Flush(t *testing.T) {
	th := NewThrottler(500*time.Millisecond, WithTrailing())
	defer th.Close()

	called := false
	var mu sync.Mutex

	th.Call(func() {
		mu.Lock()
		called = true
		mu.Unlock()
	})

	// Flush should execute immediately
	th.Flush()
	mu.Lock()
	chetest.RequireEqual(t, called, true)
	mu.Unlock()
}

func TestThrottler_Cancel(t *testing.T) {
	th := NewThrottler(100*time.Millisecond, WithTrailing())
	defer th.Close()

	called := false
	var mu sync.Mutex

	th.Call(func() {
		mu.Lock()
		called = true
		mu.Unlock()
	})

	th.Cancel()

	time.Sleep(150 * time.Millisecond)
	mu.Lock()
	chetest.RequireEqual(t, called, false)
	mu.Unlock()
}

func TestThrottler_Close(t *testing.T) {
	th := NewThrottler(100*time.Millisecond, WithTrailing())

	called := false
	var mu sync.Mutex

	th.Call(func() {
		mu.Lock()
		called = true
		mu.Unlock()
	})

	th.Close()

	time.Sleep(150 * time.Millisecond)
	mu.Lock()
	chetest.RequireEqual(t, called, false)
	mu.Unlock()

	// Further calls should be ignored
	th.Call(func() {
		mu.Lock()
		called = true
		mu.Unlock()
	})
	time.Sleep(150 * time.Millisecond)
	mu.Lock()
	chetest.RequireEqual(t, called, false)
	mu.Unlock()
}

func TestDebounce_Function(t *testing.T) {
	count := 0
	var mu sync.Mutex

	debounced := Debounce(50*time.Millisecond, func() {
		mu.Lock()
		count++
		mu.Unlock()
	})

	// Call multiple times
	for i := 0; i < 5; i++ {
		debounced()
		time.Sleep(10 * time.Millisecond)
	}

	time.Sleep(100 * time.Millisecond)

	mu.Lock()
	chetest.RequireEqual(t, count, 1)
	mu.Unlock()
}

func TestThrottle_Function(t *testing.T) {
	count := 0
	var mu sync.Mutex

	throttled := Throttle(100*time.Millisecond, func() {
		mu.Lock()
		count++
		mu.Unlock()
	}, WithLeading())

	// First call executes immediately
	throttled()
	mu.Lock()
	chetest.RequireEqual(t, count, 1)
	mu.Unlock()

	// Subsequent calls within interval are ignored
	for i := 0; i < 5; i++ {
		throttled()
		time.Sleep(10 * time.Millisecond)
	}

	mu.Lock()
	chetest.RequireEqual(t, count, 1)
	mu.Unlock()

	// After interval, can execute again
	time.Sleep(100 * time.Millisecond)
	throttled()

	mu.Lock()
	chetest.RequireEqual(t, count, 2)
	mu.Unlock()
}

func TestDebouncer_ConcurrentCalls(t *testing.T) {
	d := NewDebouncer(50 * time.Millisecond)
	defer d.Close()

	count := 0
	var mu sync.Mutex

	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			d.Call(func() {
				mu.Lock()
				count++
				mu.Unlock()
			})
		}()
	}

	wg.Wait()
	time.Sleep(100 * time.Millisecond)

	mu.Lock()
	chetest.RequireEqual(t, count, 1)
	mu.Unlock()
}

func TestThrottler_ConcurrentCalls(t *testing.T) {
	th := NewThrottler(50*time.Millisecond, WithLeading())
	defer th.Close()

	count := 0
	var mu sync.Mutex

	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			th.Call(func() {
				mu.Lock()
				count++
				mu.Unlock()
			})
		}()
	}

	wg.Wait()
	time.Sleep(100 * time.Millisecond)

	mu.Lock()
	// Should execute at least once (leading edge)
	chetest.RequireEqual(t, count >= 1, true)
	mu.Unlock()
}

func TestDebouncer_FlushWithoutPending(t *testing.T) {
	d := NewDebouncer(50 * time.Millisecond)
	defer d.Close()

	// Flush with no pending call should not panic
	d.Flush()

	called := false
	d.Call(func() {
		called = true
	})
	d.Flush()

	chetest.RequireEqual(t, called, true)

	// Second flush should be safe
	d.Flush()
}

func TestThrottler_FlushWithoutPending(t *testing.T) {
	th := NewThrottler(100*time.Millisecond, WithTrailing())
	defer th.Close()

	// Flush with no pending call should not panic
	th.Flush()
}

func TestDebouncer_MultipleFlush(t *testing.T) {
	d := NewDebouncer(50 * time.Millisecond)
	defer d.Close()

	count := 0
	d.Call(func() {
		count++
	})

	d.Flush()
	chetest.RequireEqual(t, count, 1)

	// Second flush should not call again
	d.Flush()
	chetest.RequireEqual(t, count, 1)
}
