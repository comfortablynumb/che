package cheoption

import "fmt"

// Optional represents a value that may or may not be present.
type Optional[T any] struct {
	value   T
	present bool
}

// Some creates an Optional with a value.
func Some[T any](value T) Optional[T] {
	return Optional[T]{value: value, present: true}
}

// None creates an empty Optional.
func None[T any]() Optional[T] {
	return Optional[T]{present: false}
}

// IsPresent returns true if the Optional contains a value.
func (o Optional[T]) IsPresent() bool {
	return o.present
}

// IsEmpty returns true if the Optional is empty.
func (o Optional[T]) IsEmpty() bool {
	return !o.present
}

// Get returns the value if present, panics otherwise.
func (o Optional[T]) Get() T {
	if !o.present {
		panic("cheoption: Get called on None")
	}
	return o.value
}

// GetOr returns the value if present, otherwise returns the default value.
func (o Optional[T]) GetOr(defaultValue T) T {
	if o.present {
		return o.value
	}
	return defaultValue
}

// GetOrElse returns the value if present, otherwise calls the function and returns its result.
func (o Optional[T]) GetOrElse(fn func() T) T {
	if o.present {
		return o.value
	}
	return fn()
}

// OrElse returns this Optional if present, otherwise returns the other Optional.
func (o Optional[T]) OrElse(other Optional[T]) Optional[T] {
	if o.present {
		return o
	}
	return other
}

// Map transforms the value inside the Optional if present.
func (o Optional[T]) Map(fn func(T) T) Optional[T] {
	if !o.present {
		return None[T]()
	}
	return Some(fn(o.value))
}

// FlatMap transforms the value and flattens the result.
func (o Optional[T]) FlatMap(fn func(T) Optional[T]) Optional[T] {
	if !o.present {
		return None[T]()
	}
	return fn(o.value)
}

// Filter returns this Optional if present and the predicate returns true, otherwise None.
func (o Optional[T]) Filter(predicate func(T) bool) Optional[T] {
	if !o.present {
		return None[T]()
	}
	if predicate(o.value) {
		return o
	}
	return None[T]()
}

// IfPresent calls the function if the Optional has a value.
func (o Optional[T]) IfPresent(fn func(T)) {
	if o.present {
		fn(o.value)
	}
}

// IfPresentOrElse calls the first function if present, otherwise calls the second.
func (o Optional[T]) IfPresentOrElse(presentFn func(T), emptyFn func()) {
	if o.present {
		presentFn(o.value)
	} else {
		emptyFn()
	}
}

// Unpack returns the value and a boolean indicating presence.
func (o Optional[T]) Unpack() (T, bool) {
	return o.value, o.present
}

// String returns a string representation of the Optional.
func (o Optional[T]) String() string {
	if o.present {
		return fmt.Sprintf("Some(%v)", o.value)
	}
	return "None"
}

// Result represents the result of an operation that can succeed or fail.
type Result[T any] struct {
	value T
	err   error
}

// Ok creates a successful Result.
func Ok[T any](value T) Result[T] {
	return Result[T]{value: value, err: nil}
}

// Err creates a failed Result.
func Err[T any](err error) Result[T] {
	var zero T
	return Result[T]{value: zero, err: err}
}

// IsOk returns true if the Result is successful.
func (r Result[T]) IsOk() bool {
	return r.err == nil
}

// IsErr returns true if the Result is a failure.
func (r Result[T]) IsErr() bool {
	return r.err != nil
}

// Unwrap returns the value if successful, panics otherwise.
func (r Result[T]) Unwrap() T {
	if r.err != nil {
		panic(fmt.Sprintf("cheoption: Unwrap called on Err: %v", r.err))
	}
	return r.value
}

// UnwrapOr returns the value if successful, otherwise returns the default value.
func (r Result[T]) UnwrapOr(defaultValue T) T {
	if r.err != nil {
		return defaultValue
	}
	return r.value
}

// UnwrapOrElse returns the value if successful, otherwise calls the function.
func (r Result[T]) UnwrapOrElse(fn func(error) T) T {
	if r.err != nil {
		return fn(r.err)
	}
	return r.value
}

// Error returns the error if failed, nil otherwise.
func (r Result[T]) Error() error {
	return r.err
}

// Map transforms the value inside the Result if successful.
func (r Result[T]) Map(fn func(T) T) Result[T] {
	if r.err != nil {
		return Err[T](r.err)
	}
	return Ok(fn(r.value))
}

// MapErr transforms the error inside the Result if failed.
func (r Result[T]) MapErr(fn func(error) error) Result[T] {
	if r.err == nil {
		return r
	}
	return Err[T](fn(r.err))
}

// FlatMap transforms the value and flattens the result.
func (r Result[T]) FlatMap(fn func(T) Result[T]) Result[T] {
	if r.err != nil {
		return Err[T](r.err)
	}
	return fn(r.value)
}

// AndThen is an alias for FlatMap.
func (r Result[T]) AndThen(fn func(T) Result[T]) Result[T] {
	return r.FlatMap(fn)
}

// Or returns this Result if successful, otherwise returns the other Result.
func (r Result[T]) Or(other Result[T]) Result[T] {
	if r.err == nil {
		return r
	}
	return other
}

// Unpack returns the value and error.
func (r Result[T]) Unpack() (T, error) {
	return r.value, r.err
}

// String returns a string representation of the Result.
func (r Result[T]) String() string {
	if r.err == nil {
		return fmt.Sprintf("Ok(%v)", r.value)
	}
	return fmt.Sprintf("Err(%v)", r.err)
}
