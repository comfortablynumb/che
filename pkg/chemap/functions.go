package chemap

import "golang.org/x/exp/maps"

// Functions

// Keys Returns a slice with the keys found on the given map.
// This function wraps the standard library's maps.Keys.
func Keys[K comparable, T any](m map[K]T) []K {
	return maps.Keys(m)
}

// Values Returns a slice with the values found on the given map.
// This function wraps the standard library's maps.Values.
func Values[K comparable, V any](m map[K]V) []V {
	return maps.Values(m)
}

// Invert swaps keys and values in a map. If multiple keys have the same value,
// only one key-value pair will be retained (last one encountered during iteration).
func Invert[K comparable, V comparable](m map[K]V) map[V]K {
	result := make(map[V]K, len(m))

	for k, v := range m {
		result[v] = k
	}

	return result
}

// Filter returns a new map containing only the key-value pairs for which the predicate returns true.
func Filter[K comparable, V any](m map[K]V, predicate func(K, V) bool) map[K]V {
	result := make(map[K]V)

	for k, v := range m {
		if predicate(k, v) {
			result[k] = v
		}
	}

	return result
}

// MapValues returns a new map with the same keys but with values transformed by the mapper function.
func MapValues[K comparable, V any, R any](m map[K]V, mapper func(V) R) map[K]R {
	result := make(map[K]R, len(m))

	for k, v := range m {
		result[k] = mapper(v)
	}

	return result
}

// Merge returns a new map containing all key-value pairs from all input maps.
// If a key exists in multiple maps, the value from the last map takes precedence.
func Merge[K comparable, V any](maps ...map[K]V) map[K]V {
	result := make(map[K]V)

	for _, m := range maps {
		for k, v := range m {
			result[k] = v
		}
	}

	return result
}

// Pick returns a new map containing only the specified keys from the input map.
func Pick[K comparable, V any](m map[K]V, keys ...K) map[K]V {
	result := make(map[K]V)

	for _, key := range keys {
		if v, ok := m[key]; ok {
			result[key] = v
		}
	}

	return result
}

// Omit returns a new map containing all keys except the specified ones.
func Omit[K comparable, V any](m map[K]V, keys ...K) map[K]V {
	omitSet := make(map[K]struct{}, len(keys))
	for _, key := range keys {
		omitSet[key] = struct{}{}
	}

	result := make(map[K]V)

	for k, v := range m {
		if _, shouldOmit := omitSet[k]; !shouldOmit {
			result[k] = v
		}
	}

	return result
}
