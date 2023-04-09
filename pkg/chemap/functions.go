package chemap

// Functions

// Keys Returns a slice with the keys found on the given map.
func Keys[K comparable, T any](m map[K]T) []K {
	result := make([]K, 0, len(m))

	for k := range m {
		result = append(result, k)
	}

	return result
}
