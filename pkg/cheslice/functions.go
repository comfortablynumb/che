package cheslice

// Static Functions

func Intersect[T comparable](slices ...[]T) []T {
	result := make([]T, 0)

	if len(slices) == 0 {
		return result
	}

	for _, element := range slices[0] {
		exists := true

		for _, slice := range slices[1:] {
			if !Contains(slice, element) {
				exists = false

				break
			}
		}

		if exists && !Contains(result, element) {
			result = append(result, element)
		}
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
