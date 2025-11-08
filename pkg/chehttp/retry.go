package chehttp

import (
	"math"
	"time"
)

// RetryConfig configures retry behavior for HTTP requests
type RetryConfig struct {
	MaxRetries            int
	RetryableStatusCodes  []StatusCodeRange
	NonRetryableStatusCodes []int
	BackoffStrategy       BackoffStrategy
}

// StatusCodeRange represents a range of HTTP status codes
type StatusCodeRange struct {
	Min int
	Max int
}

// Contains checks if a status code is in the range
func (r StatusCodeRange) Contains(statusCode int) bool {
	return statusCode >= r.Min && statusCode <= r.Max
}

// BackoffStrategy defines how to calculate wait time between retries
type BackoffStrategy interface {
	// NextBackoff returns the duration to wait before the next retry
	// attempt is the retry attempt number (0-indexed)
	NextBackoff(attempt int) time.Duration
}

// FixedBackoff waits a fixed duration between retries
type FixedBackoff struct {
	Delay time.Duration
}

// NextBackoff returns the fixed delay
func (f FixedBackoff) NextBackoff(attempt int) time.Duration {
	return f.Delay
}

// LinearBackoff increases wait time linearly
type LinearBackoff struct {
	BaseDelay time.Duration
}

// NextBackoff returns delay that increases linearly: baseDelay * (attempt + 1)
func (l LinearBackoff) NextBackoff(attempt int) time.Duration {
	return l.BaseDelay * time.Duration(attempt+1)
}

// ExponentialBackoff increases wait time exponentially
type ExponentialBackoff struct {
	BaseDelay  time.Duration
	Multiplier float64
	MaxDelay   time.Duration
}

// NextBackoff returns delay that increases exponentially: baseDelay * (multiplier ^ attempt)
func (e ExponentialBackoff) NextBackoff(attempt int) time.Duration {
	delay := float64(e.BaseDelay) * math.Pow(e.Multiplier, float64(attempt))
	duration := time.Duration(delay)

	if e.MaxDelay > 0 && duration > e.MaxDelay {
		return e.MaxDelay
	}

	return duration
}

// DefaultRetryConfig returns a sensible default retry configuration
func DefaultRetryConfig() *RetryConfig {
	return &RetryConfig{
		MaxRetries: 3,
		RetryableStatusCodes: []StatusCodeRange{
			{Min: 408, Max: 408}, // Request Timeout
			{Min: 429, Max: 429}, // Too Many Requests
			{Min: 500, Max: 599}, // Server Errors
		},
		NonRetryableStatusCodes: []int{
			501, // Not Implemented
			505, // HTTP Version Not Supported
		},
		BackoffStrategy: ExponentialBackoff{
			BaseDelay:  100 * time.Millisecond,
			Multiplier: 2.0,
			MaxDelay:   10 * time.Second,
		},
	}
}

// shouldRetry determines if a request should be retried based on status code
func (r *RetryConfig) shouldRetry(statusCode int, attempt int) bool {
	if attempt >= r.MaxRetries {
		return false
	}

	// Check if status code is in non-retryable list
	for _, code := range r.NonRetryableStatusCodes {
		if statusCode == code {
			return false
		}
	}

	// Check if status code is in retryable ranges
	for _, codeRange := range r.RetryableStatusCodes {
		if codeRange.Contains(statusCode) {
			return true
		}
	}

	return false
}
