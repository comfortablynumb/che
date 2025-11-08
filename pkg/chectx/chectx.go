package chectx

import (
	"context"
)

// contextKey is a type for context keys to avoid collisions
type contextKey[T any] struct {
	name string
}

// Key creates a new typed context key
// This ensures type safety when storing and retrieving values from context
func Key[T any](name string) *contextKey[T] {
	return &contextKey[T]{name: name}
}

// WithValue returns a copy of parent context with the value associated with key
// This is a type-safe wrapper around context.WithValue
func WithValue[T any](ctx context.Context, key *contextKey[T], value T) context.Context {
	return context.WithValue(ctx, key, value)
}

// Value retrieves the value associated with key from the context
// Returns the value and true if found, zero value and false otherwise
func Value[T any](ctx context.Context, key *contextKey[T]) (T, bool) {
	val := ctx.Value(key)
	if val == nil {
		var zero T
		return zero, false
	}

	typedVal, ok := val.(T)
	return typedVal, ok
}

// MustValue retrieves the value associated with key from the context
// Panics if the value is not found or has the wrong type
func MustValue[T any](ctx context.Context, key *contextKey[T]) T {
	val, ok := Value(ctx, key)
	if !ok {
		panic("chectx: value not found in context for key: " + key.name)
	}
	return val
}

// GetOrDefault retrieves the value associated with key from the context
// Returns the default value if not found
func GetOrDefault[T any](ctx context.Context, key *contextKey[T], defaultValue T) T {
	val, ok := Value(ctx, key)
	if !ok {
		return defaultValue
	}
	return val
}
