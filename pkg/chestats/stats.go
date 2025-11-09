package chestats

import (
	"fmt"
	"math"
	"sort"

	"golang.org/x/exp/constraints"
)

// Number represents numeric types that can be used with statistical functions.
type Number interface {
	constraints.Integer | constraints.Float
}

// Mean calculates the arithmetic mean (average) of a slice of numbers.
// Returns 0 for empty slice.
func Mean[T Number](values []T) float64 {
	if len(values) == 0 {
		return 0
	}

	var sum float64
	for _, v := range values {
		sum += float64(v)
	}

	return sum / float64(len(values))
}

// Median calculates the median (middle value) of a slice of numbers.
// Returns 0 for empty slice.
func Median[T Number](values []T) float64 {
	if len(values) == 0 {
		return 0
	}

	// Create a copy to avoid modifying the original
	sorted := make([]T, len(values))
	copy(sorted, values)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i] < sorted[j]
	})

	mid := len(sorted) / 2

	if len(sorted)%2 == 0 {
		// Even number of elements - average the two middle values
		return (float64(sorted[mid-1]) + float64(sorted[mid])) / 2
	}

	// Odd number of elements
	return float64(sorted[mid])
}

// Mode calculates the mode (most frequent value) of a slice of numbers.
// Returns 0 if the slice is empty or if there's no unique mode.
func Mode[T Number](values []T) T {
	var zero T
	if len(values) == 0 {
		return zero
	}

	freq := make(map[T]int)
	for _, v := range values {
		freq[v]++
	}

	var mode T
	maxCount := 0

	for v, count := range freq {
		if count > maxCount {
			maxCount = count
			mode = v
		}
	}

	return mode
}

// Variance calculates the population variance of a slice of numbers.
// Returns 0 for empty slice.
func Variance[T Number](values []T) float64 {
	if len(values) == 0 {
		return 0
	}

	mean := Mean(values)
	var sumSquares float64

	for _, v := range values {
		diff := float64(v) - mean
		sumSquares += diff * diff
	}

	return sumSquares / float64(len(values))
}

// SampleVariance calculates the sample variance of a slice of numbers.
// Returns 0 for slices with less than 2 elements.
func SampleVariance[T Number](values []T) float64 {
	if len(values) < 2 {
		return 0
	}

	mean := Mean(values)
	var sumSquares float64

	for _, v := range values {
		diff := float64(v) - mean
		sumSquares += diff * diff
	}

	return sumSquares / float64(len(values)-1)
}

// StdDev calculates the population standard deviation of a slice of numbers.
// Returns 0 for empty slice.
func StdDev[T Number](values []T) float64 {
	return math.Sqrt(Variance(values))
}

// SampleStdDev calculates the sample standard deviation of a slice of numbers.
// Returns 0 for slices with less than 2 elements.
func SampleStdDev[T Number](values []T) float64 {
	return math.Sqrt(SampleVariance(values))
}

// Min returns the minimum value from a slice of numbers.
// Panics if the slice is empty.
func Min[T Number](values []T) T {
	if len(values) == 0 {
		panic("chestats: Min called on empty slice")
	}

	min := values[0]
	for _, v := range values[1:] {
		if v < min {
			min = v
		}
	}

	return min
}

// Max returns the maximum value from a slice of numbers.
// Panics if the slice is empty.
func Max[T Number](values []T) T {
	if len(values) == 0 {
		panic("chestats: Max called on empty slice")
	}

	max := values[0]
	for _, v := range values[1:] {
		if v > max {
			max = v
		}
	}

	return max
}

// Sum calculates the sum of a slice of numbers.
func Sum[T Number](values []T) T {
	var sum T
	for _, v := range values {
		sum += v
	}
	return sum
}

// Percentile calculates the nth percentile of a slice of numbers.
// p should be between 0 and 100.
// Returns 0 for empty slice.
func Percentile[T Number](values []T, p float64) float64 {
	if len(values) == 0 {
		return 0
	}

	if p < 0 || p > 100 {
		panic(fmt.Sprintf("chestats: percentile must be between 0 and 100, got %f", p))
	}

	// Create a copy to avoid modifying the original
	sorted := make([]T, len(values))
	copy(sorted, values)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i] < sorted[j]
	})

	if p == 0 {
		return float64(sorted[0])
	}

	if p == 100 {
		return float64(sorted[len(sorted)-1])
	}

	// Linear interpolation between closest ranks
	rank := (p / 100) * float64(len(sorted)-1)
	lowerIndex := int(math.Floor(rank))
	upperIndex := int(math.Ceil(rank))

	if lowerIndex == upperIndex {
		return float64(sorted[lowerIndex])
	}

	weight := rank - float64(lowerIndex)
	return float64(sorted[lowerIndex])*(1-weight) + float64(sorted[upperIndex])*weight
}

// Quartiles calculates the first (Q1), second (Q2/median), and third (Q3) quartiles.
// Returns (0, 0, 0) for empty slice.
func Quartiles[T Number](values []T) (q1, q2, q3 float64) {
	if len(values) == 0 {
		return 0, 0, 0
	}

	return Percentile(values, 25), Percentile(values, 50), Percentile(values, 75)
}

// IQR calculates the Interquartile Range (Q3 - Q1).
// Returns 0 for empty slice.
func IQR[T Number](values []T) float64 {
	if len(values) == 0 {
		return 0
	}

	q1, _, q3 := Quartiles(values)
	return q3 - q1
}

// Range calculates the range (max - min) of a slice of numbers.
// Panics if the slice is empty.
func Range[T Number](values []T) T {
	return Max(values) - Min(values)
}

// Correlation calculates the Pearson correlation coefficient between two slices.
// Returns 0 if slices are empty or of different lengths.
// Returns value between -1 and 1, where:
//   -1 indicates perfect negative correlation
//    0 indicates no correlation
//    1 indicates perfect positive correlation
func Correlation[T Number](x, y []T) float64 {
	if len(x) == 0 || len(y) == 0 || len(x) != len(y) {
		return 0
	}

	meanX := Mean(x)
	meanY := Mean(y)

	var numerator, sumSqX, sumSqY float64

	for i := range x {
		diffX := float64(x[i]) - meanX
		diffY := float64(y[i]) - meanY
		numerator += diffX * diffY
		sumSqX += diffX * diffX
		sumSqY += diffY * diffY
	}

	denominator := math.Sqrt(sumSqX * sumSqY)
	if denominator == 0 {
		return 0
	}

	return numerator / denominator
}
