package cheslice

// Types

type MapFunc[T any] func(element T) T

type FilterFunc[T any] func(element T) bool

// Static Functions

func Union[T any](slices ...[]T) []T {
	if len(slices) == 0 {
		return []T{}
	}

	// Calculate capacity

	capacity := 0

	for _, slice := range slices {
		capacity += len(slice)
	}

	result := make([]T, 0, capacity)

	for _, slice := range slices {
		for _, element := range slice {
			result = append(result, element)
		}
	}

	return result
}

func Map[T any](slice []T, mapFunc MapFunc[T]) {
	for i, element := range slice {
		slice[i] = mapFunc(element)
	}
}

func Filter[T any](slice []T, filterFunc FilterFunc[T]) []T {
	result := make([]T, 0)

	for _, element := range slice {
		if filterFunc(element) {
			result = append(result, element)
		}
	}

	return result
}

func Fill[T any](count uint, value T) []T {
	result := make([]T, 0, count)

	for i := uint(0); i < count; i++ {
		result = append(result, value)
	}

	return result
}

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

		include := true

		for _, slice := range slices[1:] {
			if Contains(slice, element) {
				include = false

				break
			}
		}

		if include {
			result = append(result, element)
		}
	}

	return result
}

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

		exists := true

		for _, slice := range slices[1:] {
			if !Contains(slice, element) {
				exists = false

				break
			}
		}

		if !exists {
			continue
		}

		result = append(result, element)

		m[element] = struct{}{}
	}

	return result
}

func Contains[T comparable](slice []T, element T) bool {
	for _, e := range slice {
		if e == element {
			return true
		}
	}

	return false
}
