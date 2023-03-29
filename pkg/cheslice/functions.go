package cheslice

// Static Functions

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
