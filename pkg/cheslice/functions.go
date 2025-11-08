package cheslice

// Types

type ForEachFunc[T any] func(element T) bool

type MapFunc[T any] func(element T) T

type FilterFunc[T any] func(element T) bool

// Functions

// Union Returns a new slice with all the elements found in the given slices. It preserves repeated elements.
func Union[T any](slices ...[]T) []T {
	if len(slices) == 0 {
		return []T{}
	}

	result := make([]T, 0, Len(slices))

	for _, slice := range slices {
		result = append(result, slice...)
	}

	return result
}

// ForEach Executes the given "forEachFunc" on each of the elements of the received slice.
func ForEach[T any](slice []T, forEachFunc ForEachFunc[T]) {
	for _, element := range slice {
		if !forEachFunc(element) {
			return
		}
	}
}

// Map Returns a new slice with the result of applying "mapFunc" to each of the elements from the given slice.
func Map[T any](slice []T, mapFunc MapFunc[T]) []T {
	result := make([]T, 0, len(slice))

	for _, element := range slice {
		result = append(result, mapFunc(element))
	}

	return result
}

// Filter Returns a new slice with the elements for which "filterFunc" returned true.
func Filter[T any](slice []T, filterFunc FilterFunc[T]) []T {
	result := make([]T, 0)

	for _, element := range slice {
		if filterFunc(element) {
			result = append(result, element)
		}
	}

	return result
}

// Fill Creates a new slice with the amount of elements determined by "count". Each element will have the value
// determined by "value".
func Fill[T any](count uint, value T) []T {
	result := make([]T, 0, count)

	for i := uint(0); i < count; i++ {
		result = append(result, value)
	}

	return result
}

// Diff Returns a new slice with all the elements found in the first slice that are NOT present in the rest of the
// slices. If no slice is received, it returns an empty slice. If one slice is received, it returns a copy of it.
func Diff[T comparable](slices ...[]T) []T {
	result := make([]T, 0)

	if len(slices) < 1 {
		return result
	}

	if len(slices) == 1 {
		return append(result, slices[0]...)
	}

	checkedElements := make(map[T]struct{})

	for _, element := range slices[0] {
		_, found := checkedElements[element]

		if found {
			continue
		}

		checkedElements[element] = struct{}{}

		if !Exists(element, slices[1:]...) {
			result = append(result, element)
		}
	}

	return result
}

// Chunk Returns a new slice consisting of chunks with the length determined by "length". If slice is empty, then this
// function returns an empty slice.
func Chunk[T any](slice []T, length uint) [][]T {
	result := make([][]T, 0)

	if length < 1 {
		return result
	}

	sliceSize := uint(len(slice))
	elementsLeft := sliceSize
	currentIndex := uint(0)

	for {
		if elementsLeft < 1 {
			break
		}

		chunkSize := length

		if elementsLeft < length {
			chunkSize = elementsLeft
		}

		chunk := make([]T, 0, chunkSize)

		chunk = append(chunk, slice[currentIndex:currentIndex+chunkSize]...)

		result = append(result, chunk)

		elementsLeft -= chunkSize
		currentIndex += chunkSize
	}

	return result
}

// Unique Returns a new slice with all the distinct values found in the given slice.
func Unique[T comparable](slice []T) []T {
	result := make([]T, 0)
	m := make(map[T]struct{})

	for _, element := range slice {
		_, found := m[element]

		if found {
			continue
		}

		result = append(result, element)

		m[element] = struct{}{}
	}

	return result
}

// Intersect Returns a new slice with the elements that are found in ALL the given slices. If no slice is given, then
// it returns an empty slice. If only ne slice is given, it rethrns a copy of the same slice (including repeated
// elements).
func Intersect[T comparable](slices ...[]T) []T {
	result := make([]T, 0)

	if len(slices) == 0 {
		return result
	}

	if len(slices) == 1 {
		return append(result, slices[0]...)
	}

	m := make(map[T]struct{})

	for _, element := range slices[0] {
		if _, found := m[element]; found {
			continue
		}

		if !Exists(element, slices[1:]...) {
			continue
		}

		result = append(result, element)

		m[element] = struct{}{}
	}

	return result
}

// Exists Returns true if the given element is present in ANY of the given slices. Returns false otherwise.
func Exists[T comparable](element T, slices ...[]T) bool {
	for _, slice := range slices {
		for _, e := range slice {
			if e == element {
				return true
			}
		}
	}

	return false
}

// Len Returns the sum of the lengths of all the given slices.
func Len[T any](slices ...[]T) int {
	result := 0

	for _, slice := range slices {
		result += len(slice)
	}

	return result
}

// Reduce reduces a slice to a single value using the given reducer function.
// The accumulator is initialized with the initial value.
func Reduce[T any, R any](slice []T, initial R, reducer func(acc R, element T) R) R {
	acc := initial
	for _, element := range slice {
		acc = reducer(acc, element)
	}
	return acc
}

// GroupBy groups slice elements by a key function.
// Returns a map where each key maps to a slice of elements that produced that key.
func GroupBy[T any, K comparable](slice []T, keyFunc func(T) K) map[K][]T {
	result := make(map[K][]T)
	for _, element := range slice {
		key := keyFunc(element)
		result[key] = append(result[key], element)
	}
	return result
}

// Partition splits a slice into two slices based on a predicate.
// The first slice contains elements for which the predicate returns true,
// the second contains elements for which it returns false.
func Partition[T any](slice []T, predicate func(T) bool) ([]T, []T) {
	truthy := make([]T, 0)
	falsy := make([]T, 0)

	for _, element := range slice {
		if predicate(element) {
			truthy = append(truthy, element)
		} else {
			falsy = append(falsy, element)
		}
	}

	return truthy, falsy
}

// Flatten flattens a slice of slices into a single slice.
func Flatten[T any](slices [][]T) []T {
	result := make([]T, 0)
	for _, slice := range slices {
		result = append(result, slice...)
	}
	return result
}

// Zip combines two slices into a slice of pairs.
// The resulting slice length is the minimum of the two input slices.
func Zip[T any, U any](slice1 []T, slice2 []U) [][2]interface{} {
	minLen := len(slice1)
	if len(slice2) < minLen {
		minLen = len(slice2)
	}

	result := make([][2]interface{}, minLen)
	for i := 0; i < minLen; i++ {
		result[i] = [2]interface{}{slice1[i], slice2[i]}
	}
	return result
}

// Take returns the first n elements from the slice.
// If n is greater than the slice length, returns the entire slice.
func Take[T any](slice []T, n int) []T {
	if n <= 0 {
		return []T{}
	}
	if n >= len(slice) {
		result := make([]T, len(slice))
		copy(result, slice)
		return result
	}
	result := make([]T, n)
	copy(result, slice[:n])
	return result
}

// Drop returns a slice with the first n elements removed.
// If n is greater than or equal to the slice length, returns an empty slice.
func Drop[T any](slice []T, n int) []T {
	if n <= 0 {
		result := make([]T, len(slice))
		copy(result, slice)
		return result
	}
	if n >= len(slice) {
		return []T{}
	}
	result := make([]T, len(slice)-n)
	copy(result, slice[n:])
	return result
}

// TakeWhile returns elements from the slice while the predicate returns true.
// Stops at the first element for which the predicate returns false.
func TakeWhile[T any](slice []T, predicate func(T) bool) []T {
	result := make([]T, 0)
	for _, element := range slice {
		if !predicate(element) {
			break
		}
		result = append(result, element)
	}
	return result
}

// DropWhile drops elements from the slice while the predicate returns true.
// Returns the remaining elements starting from the first element for which the predicate returns false.
func DropWhile[T any](slice []T, predicate func(T) bool) []T {
	for i, element := range slice {
		if !predicate(element) {
			result := make([]T, len(slice)-i)
			copy(result, slice[i:])
			return result
		}
	}
	return []T{}
}

// Any returns true if the predicate returns true for any element in the slice.
func Any[T any](slice []T, predicate func(T) bool) bool {
	for _, element := range slice {
		if predicate(element) {
			return true
		}
	}
	return false
}

// All returns true if the predicate returns true for all elements in the slice.
func All[T any](slice []T, predicate func(T) bool) bool {
	for _, element := range slice {
		if !predicate(element) {
			return false
		}
	}
	return true
}

// None returns true if the predicate returns false for all elements in the slice.
func None[T any](slice []T, predicate func(T) bool) bool {
	return !Any(slice, predicate)
}

// Reverse returns a new slice with elements in reverse order.
func Reverse[T any](slice []T) []T {
	result := make([]T, len(slice))
	for i, element := range slice {
		result[len(slice)-1-i] = element
	}
	return result
}

// Find returns the first element for which the predicate returns true.
// Returns the element and true if found, or zero value and false otherwise.
func Find[T any](slice []T, predicate func(T) bool) (T, bool) {
	for _, element := range slice {
		if predicate(element) {
			return element, true
		}
	}
	var zero T
	return zero, false
}

// FindIndex returns the index of the first element for which the predicate returns true.
// Returns the index and true if found, or -1 and false otherwise.
func FindIndex[T any](slice []T, predicate func(T) bool) (int, bool) {
	for i, element := range slice {
		if predicate(element) {
			return i, true
		}
	}
	return -1, false
}

// Count returns the number of elements for which the predicate returns true.
func Count[T any](slice []T, predicate func(T) bool) int {
	count := 0
	for _, element := range slice {
		if predicate(element) {
			count++
		}
	}
	return count
}
