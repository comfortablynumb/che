package checircuit

import (
	"errors"
	"testing"
	"time"
)

func TestBreaker_Closed(t *testing.T) {
	cb := New(nil)

	if cb.State() != StateClosed {
		t.Error("initial state should be closed")
	}

	// Should allow requests through
	err := cb.Execute(func() error {
		return nil
	})

	if err != nil {
		t.Errorf("closed circuit should allow requests: %v", err)
	}
}

func TestBreaker_OpensAfterFailures(t *testing.T) {
	cb := New(&Config{
		MaxFailures: 3,
		Timeout:     100 * time.Millisecond,
	})

	// First 2 failures should keep circuit closed
	for i := 0; i < 2; i++ {
		_ = cb.Execute(func() error {
			return errors.New("failure")
		})

		if cb.State() != StateClosed {
			t.Errorf("circuit should remain closed after %d failures", i+1)
		}
	}

	// Third failure should open the circuit
	_ = cb.Execute(func() error {
		return errors.New("failure")
	})

	if cb.State() != StateOpen {
		t.Error("circuit should be open after max failures")
	}

	if cb.Failures() != 3 {
		t.Errorf("expected 3 failures, got %d", cb.Failures())
	}
}

func TestBreaker_Open(t *testing.T) {
	cb := New(&Config{
		MaxFailures: 1,
		Timeout:     1 * time.Second,
	})

	// Trigger open state
	_ = cb.Execute(func() error {
		return errors.New("failure")
	})

	// Should reject requests
	err := cb.Execute(func() error {
		t.Error("should not execute when circuit is open")
		return nil
	})

	if err != ErrCircuitOpen {
		t.Errorf("expected ErrCircuitOpen, got %v", err)
	}
}

func TestBreaker_HalfOpen(t *testing.T) {
	cb := New(&Config{
		MaxFailures: 1,
		Timeout:     50 * time.Millisecond,
		MaxRequests: 2,
	})

	// Open the circuit
	_ = cb.Execute(func() error {
		return errors.New("failure")
	})

	if cb.State() != StateOpen {
		t.Error("circuit should be open")
	}

	// Wait for timeout
	time.Sleep(60 * time.Millisecond)

	// Should transition to half-open
	err := cb.Execute(func() error {
		return nil
	})

	if err != nil {
		t.Errorf("half-open circuit should allow request: %v", err)
	}
}

func TestBreaker_HalfOpenToClosedOnSuccess(t *testing.T) {
	cb := New(&Config{
		MaxFailures: 1,
		Timeout:     50 * time.Millisecond,
		MaxRequests: 2,
	})

	// Open the circuit
	_ = cb.Execute(func() error {
		return errors.New("failure")
	})

	// Wait for timeout to transition to half-open
	time.Sleep(60 * time.Millisecond)

	// Successful requests should close the circuit
	for i := 0; i < 2; i++ {
		err := cb.Execute(func() error {
			return nil
		})
		if err != nil {
			t.Errorf("request %d failed: %v", i, err)
		}
	}

	if cb.State() != StateClosed {
		t.Error("circuit should be closed after successful requests in half-open")
	}
}

func TestBreaker_HalfOpenToOpenOnFailure(t *testing.T) {
	cb := New(&Config{
		MaxFailures: 1,
		Timeout:     50 * time.Millisecond,
		MaxRequests: 2,
	})

	// Open the circuit
	_ = cb.Execute(func() error {
		return errors.New("failure")
	})

	// Wait for timeout to transition to half-open
	time.Sleep(60 * time.Millisecond)

	// Failure in half-open should reopen
	_ = cb.Execute(func() error {
		return errors.New("failure")
	})

	if cb.State() != StateOpen {
		t.Error("circuit should reopen on failure in half-open state")
	}
}

func TestBreaker_HalfOpenMaxRequests(t *testing.T) {
	cb := New(&Config{
		MaxFailures: 1,
		Timeout:     50 * time.Millisecond,
		MaxRequests: 2,
	})

	// Open the circuit
	_ = cb.Execute(func() error {
		return errors.New("failure")
	})

	// Wait for timeout
	time.Sleep(60 * time.Millisecond)

	// First request should succeed
	err := cb.Execute(func() error {
		return nil
	})
	if err != nil {
		t.Errorf("first request failed: %v", err)
	}

	// Second request should succeed (still in half-open, within MaxRequests)
	err = cb.Execute(func() error {
		return nil
	})
	if err != nil {
		t.Errorf("second request failed: %v", err)
	}

	// After 2 successes with MaxRequests=2, circuit should be closed
	if cb.State() != StateClosed {
		t.Error("circuit should be closed after max successful requests in half-open")
	}

	// Third request should succeed because circuit is now closed
	err = cb.Execute(func() error {
		return nil
	})

	if err != nil {
		t.Errorf("request should succeed when circuit is closed: %v", err)
	}
}

func TestBreaker_Reset(t *testing.T) {
	cb := New(&Config{
		MaxFailures: 1,
	})

	// Open the circuit
	_ = cb.Execute(func() error {
		return errors.New("failure")
	})

	if cb.State() != StateOpen {
		t.Error("circuit should be open")
	}

	// Reset
	cb.Reset()

	if cb.State() != StateClosed {
		t.Error("circuit should be closed after reset")
	}

	if cb.Failures() != 0 {
		t.Errorf("failures should be 0 after reset, got %d", cb.Failures())
	}

	// Should allow requests
	err := cb.Execute(func() error {
		return nil
	})

	if err != nil {
		t.Errorf("reset circuit should allow requests: %v", err)
	}
}

func TestBreaker_OnStateChange(t *testing.T) {
	var transitions []struct {
		from, to State
	}

	cb := New(&Config{
		MaxFailures: 1,
		Timeout:     50 * time.Millisecond,
		OnStateChange: func(from, to State) {
			transitions = append(transitions, struct{ from, to State }{from, to})
		},
	})

	// Closed -> Open
	_ = cb.Execute(func() error {
		return errors.New("failure")
	})

	if len(transitions) != 1 {
		t.Fatalf("expected 1 transition, got %d", len(transitions))
	}
	if transitions[0].from != StateClosed || transitions[0].to != StateOpen {
		t.Errorf("expected Closed->Open, got %v->%v", transitions[0].from, transitions[0].to)
	}

	// Wait for timeout to transition to half-open
	time.Sleep(60 * time.Millisecond)

	// Trigger state check by making a request
	_ = cb.Execute(func() error {
		return nil
	})

	if len(transitions) < 2 {
		t.Fatalf("expected at least 2 transitions, got %d", len(transitions))
	}
	if transitions[1].from != StateOpen || transitions[1].to != StateHalfOpen {
		t.Errorf("expected Open->HalfOpen, got %v->%v", transitions[1].from, transitions[1].to)
	}
}

func TestBreaker_SuccessResetsFailures(t *testing.T) {
	cb := New(&Config{
		MaxFailures: 3,
	})

	// 2 failures
	for i := 0; i < 2; i++ {
		_ = cb.Execute(func() error {
			return errors.New("failure")
		})
	}

	if cb.Failures() != 2 {
		t.Errorf("expected 2 failures, got %d", cb.Failures())
	}

	// Success should reset failures in closed state
	_ = cb.Execute(func() error {
		return nil
	})

	if cb.Failures() != 0 {
		t.Errorf("failures should be reset to 0 after success, got %d", cb.Failures())
	}
}

func TestBreaker_DefaultConfig(t *testing.T) {
	cb := New(nil)

	if cb.config.MaxFailures != 5 {
		t.Errorf("expected default MaxFailures 5, got %d", cb.config.MaxFailures)
	}

	if cb.config.Timeout != 60*time.Second {
		t.Errorf("expected default Timeout 60s, got %v", cb.config.Timeout)
	}

	if cb.config.MaxRequests != 1 {
		t.Errorf("expected default MaxRequests 1, got %d", cb.config.MaxRequests)
	}
}
