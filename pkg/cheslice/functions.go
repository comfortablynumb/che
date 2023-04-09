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
