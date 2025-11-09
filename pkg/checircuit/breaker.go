package checircuit

import (
	"errors"
	"sync"
	"time"
)

// State represents the state of the circuit breaker.
type State int

const (
	// StateClosed allows requests through.
	StateClosed State = iota
	// StateOpen blocks requests.
	StateOpen
	// StateHalfOpen allows a limited number of requests to test if the service recovered.
	StateHalfOpen
)

var (
	// ErrCircuitOpen is returned when the circuit breaker is open.
	ErrCircuitOpen = errors.New("circuit breaker is open")
	// ErrTooManyRequests is returned when too many requests are made in half-open state.
	ErrTooManyRequests = errors.New("too many requests")
)

// Config holds the configuration for a circuit breaker.
type Config struct {
	// MaxFailures is the maximum number of failures before opening. Default: 5
	MaxFailures uint

	// Timeout is how long to wait before transitioning from Open to HalfOpen. Default: 60s
	Timeout time.Duration

	// MaxRequests is the maximum requests allowed in HalfOpen state. Default: 1
	MaxRequests uint

	// OnStateChange is called when the state changes. Optional.
	OnStateChange func(from, to State)
}

// Breaker is a circuit breaker implementation.
type Breaker struct {
	config      Config
	state       State
	failures    uint
	successes   uint
	requests    uint
	lastFailure time.Time
	mu          sync.RWMutex
}

// New creates a new circuit breaker.
func New(config *Config) *Breaker {
	if config == nil {
		config = &Config{}
	}

	if config.MaxFailures == 0 {
		config.MaxFailures = 5
	}
	if config.Timeout == 0 {
		config.Timeout = 60 * time.Second
	}
	if config.MaxRequests == 0 {
		config.MaxRequests = 1
	}

	return &Breaker{
		config: *config,
		state:  StateClosed,
	}
}

// Execute runs the function if the circuit is not open.
func (b *Breaker) Execute(fn func() error) error {
	if err := b.beforeRequest(); err != nil {
		return err
	}

	err := fn()
	b.afterRequest(err)
	return err
}

// State returns the current state of the circuit breaker.
func (b *Breaker) State() State {
	b.mu.RLock()
	defer b.mu.RUnlock()

	state := b.state
	b.checkState()
	return state
}

// Failures returns the current failure count.
func (b *Breaker) Failures() uint {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.failures
}

// Reset resets the circuit breaker to closed state.
func (b *Breaker) Reset() {
	b.mu.Lock()
	defer b.mu.Unlock()

	oldState := b.state
	b.state = StateClosed
	b.failures = 0
	b.successes = 0
	b.requests = 0

	if oldState != StateClosed && b.config.OnStateChange != nil {
		b.config.OnStateChange(oldState, StateClosed)
	}
}

func (b *Breaker) beforeRequest() error {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.checkState()

	switch b.state {
	case StateOpen:
		return ErrCircuitOpen
	case StateHalfOpen:
		if b.requests >= b.config.MaxRequests {
			return ErrTooManyRequests
		}
		b.requests++
		return nil
	default: // StateClosed
		return nil
	}
}

func (b *Breaker) afterRequest(err error) {
	b.mu.Lock()
	defer b.mu.Unlock()

	if err != nil {
		b.onFailure()
	} else {
		b.onSuccess()
	}
}

func (b *Breaker) onSuccess() {
	switch b.state {
	case StateClosed:
		b.failures = 0
	case StateHalfOpen:
		b.successes++
		if b.successes >= b.config.MaxRequests {
			b.setState(StateClosed)
			b.failures = 0
			b.successes = 0
			b.requests = 0
		}
	}
}

func (b *Breaker) onFailure() {
	b.failures++
	b.lastFailure = time.Now()

	switch b.state {
	case StateClosed:
		if b.failures >= b.config.MaxFailures {
			b.setState(StateOpen)
		}
	case StateHalfOpen:
		b.setState(StateOpen)
		b.successes = 0
		b.requests = 0
	}
}

func (b *Breaker) checkState() {
	if b.state == StateOpen && time.Since(b.lastFailure) > b.config.Timeout {
		b.setState(StateHalfOpen)
		b.requests = 0
		b.successes = 0
	}
}

func (b *Breaker) setState(newState State) {
	if b.state == newState {
		return
	}

	oldState := b.state
	b.state = newState

	if b.config.OnStateChange != nil {
		b.config.OnStateChange(oldState, newState)
	}
}
