// Package chedebounce provides debounce and throttle utilities for rate-limiting function calls.
package chedebounce

import (
	"sync"
	"time"
)

// Debouncer delays function execution until after a specified duration has elapsed
// since the last invocation.
type Debouncer struct {
	delay  time.Duration
	timer  *time.Timer
	mu     sync.Mutex
	fn     func()
	last   time.Time
	closed bool
}

// NewDebouncer creates a new debouncer with the specified delay.
func NewDebouncer(delay time.Duration) *Debouncer {
	return &Debouncer{
		delay: delay,
	}
}

// Call schedules the function to be called after the delay.
// If called again before the delay expires, the previous call is cancelled.
func (d *Debouncer) Call(fn func()) {
	d.mu.Lock()
	defer d.mu.Unlock()

	if d.closed {
		return
	}

	d.fn = fn
	d.last = time.Now()

	if d.timer != nil {
		d.timer.Stop()
	}

	d.timer = time.AfterFunc(d.delay, func() {
		d.mu.Lock()
		fn := d.fn
		d.fn = nil
		closed := d.closed
		d.mu.Unlock()

		if fn != nil && !closed {
			fn()
		}
	})
}

// Flush immediately executes any pending function call and cancels the timer.
func (d *Debouncer) Flush() {
	d.mu.Lock()
	if d.timer != nil {
		d.timer.Stop()
		d.timer = nil
	}

	fn := d.fn
	d.fn = nil
	closed := d.closed
	d.mu.Unlock()

	if fn != nil && !closed {
		fn()
	}
}

// Cancel cancels any pending function call without executing it.
func (d *Debouncer) Cancel() {
	d.mu.Lock()
	defer d.mu.Unlock()

	if d.timer != nil {
		d.timer.Stop()
		d.timer = nil
	}

	d.fn = nil
}

// Close cancels any pending calls and marks the debouncer as closed.
func (d *Debouncer) Close() {
	d.mu.Lock()
	defer d.mu.Unlock()

	d.closed = true

	if d.timer != nil {
		d.timer.Stop()
		d.timer = nil
	}

	d.fn = nil
}

// Throttler limits function execution to at most once per specified interval.
type Throttler struct {
	interval time.Duration
	leading  bool
	trailing bool
	mu       sync.Mutex
	last     time.Time
	timer    *time.Timer
	pending  func()
	closed   bool
}

// ThrottleOption is a configuration option for Throttler.
type ThrottleOption func(*Throttler)

// WithLeading enables calling the function on the leading edge.
func WithLeading() ThrottleOption {
	return func(t *Throttler) {
		t.leading = true
	}
}

// WithTrailing enables calling the function on the trailing edge.
func WithTrailing() ThrottleOption {
	return func(t *Throttler) {
		t.trailing = true
	}
}

// NewThrottler creates a new throttler with the specified interval.
// By default, calls on the leading edge only.
// If options are provided, only the specified edges are enabled.
func NewThrottler(interval time.Duration, opts ...ThrottleOption) *Throttler {
	t := &Throttler{
		interval: interval,
		leading:  len(opts) == 0, // Default to leading if no options
		trailing: false,
	}

	for _, opt := range opts {
		opt(t)
	}

	return t
}

// Call attempts to call the function, respecting the throttle interval.
func (t *Throttler) Call(fn func()) {
	t.mu.Lock()

	if t.closed {
		t.mu.Unlock()
		return
	}

	now := time.Now()
	elapsed := now.Sub(t.last)

	// First call
	if t.last.IsZero() {
		if t.leading {
			t.last = now
			t.mu.Unlock()
			fn()
			return
		}
		// Trailing only - schedule it
		if t.trailing {
			t.last = now
			t.pending = fn
			t.scheduleTrailing()
		}
		t.mu.Unlock()
		return
	}

	// Enough time has passed since last execution
	if elapsed >= t.interval {
		if t.leading {
			t.last = now
			t.mu.Unlock()
			fn()
		} else if t.trailing {
			t.pending = fn
			t.scheduleTrailing()
			t.mu.Unlock()
		}
		return
	}

	// Function is throttled
	if t.trailing {
		t.pending = fn
		if t.timer == nil {
			t.scheduleTrailing()
		}
	}
	t.mu.Unlock()
}

func (t *Throttler) scheduleTrailing() {
	if t.timer != nil {
		t.timer.Stop()
	}

	remaining := t.interval - time.Since(t.last)
	t.timer = time.AfterFunc(remaining, func() {
		t.mu.Lock()
		fn := t.pending
		t.pending = nil
		t.timer = nil
		closed := t.closed
		if fn != nil && !closed {
			t.last = time.Now()
		}
		t.mu.Unlock()

		if fn != nil && !closed {
			fn()
		}
	})
}

// Flush immediately executes any pending trailing call.
func (t *Throttler) Flush() {
	t.mu.Lock()
	if t.timer != nil {
		t.timer.Stop()
		t.timer = nil
	}

	fn := t.pending
	t.pending = nil
	closed := t.closed
	if fn != nil && !closed {
		t.last = time.Now()
	}
	t.mu.Unlock()

	if fn != nil && !closed {
		fn()
	}
}

// Cancel cancels any pending trailing call without executing it.
func (t *Throttler) Cancel() {
	t.mu.Lock()
	defer t.mu.Unlock()

	if t.timer != nil {
		t.timer.Stop()
		t.timer = nil
	}

	t.pending = nil
}

// Close cancels any pending calls and marks the throttler as closed.
func (t *Throttler) Close() {
	t.mu.Lock()
	defer t.mu.Unlock()

	t.closed = true

	if t.timer != nil {
		t.timer.Stop()
		t.timer = nil
	}

	t.pending = nil
}

// Debounce returns a debounced version of the function.
func Debounce(delay time.Duration, fn func()) func() {
	d := NewDebouncer(delay)
	return func() {
		d.Call(fn)
	}
}

// Throttle returns a throttled version of the function.
func Throttle(interval time.Duration, fn func(), opts ...ThrottleOption) func() {
	t := NewThrottler(interval, opts...)
	return func() {
		t.Call(fn)
	}
}
