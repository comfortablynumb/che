package chestats

import (
	"math"
	"testing"
)

func floatEqual(a, b float64) bool {
	return math.Abs(a-b) < 0.0001
}

func TestMean(t *testing.T) {
	tests := []struct {
		name     string
		values   []float64
		expected float64
	}{
		{"empty slice", []float64{}, 0},
		{"single value", []float64{5}, 5},
		{"multiple values", []float64{1, 2, 3, 4, 5}, 3},
		{"negative values", []float64{-1, -2, -3}, -2},
		{"mixed values", []float64{-10, 0, 10}, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Mean(tt.values)
			if !floatEqual(result, tt.expected) {
				t.Errorf("expected %f, got %f", tt.expected, result)
			}
		})
	}
}

func TestMedian(t *testing.T) {
	tests := []struct {
		name     string
		values   []float64
		expected float64
	}{
		{"empty slice", []float64{}, 0},
		{"single value", []float64{5}, 5},
		{"odd count", []float64{1, 2, 3}, 2},
		{"even count", []float64{1, 2, 3, 4}, 2.5},
		{"unsorted", []float64{3, 1, 2}, 2},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Median(tt.values)
			if !floatEqual(result, tt.expected) {
				t.Errorf("expected %f, got %f", tt.expected, result)
			}
		})
	}
}

func TestMode(t *testing.T) {
	tests := []struct {
		name     string
		values   []int
		expected int
	}{
		{"empty slice", []int{}, 0},
		{"single value", []int{5}, 5},
		{"clear mode", []int{1, 2, 2, 3}, 2},
		{"all same", []int{5, 5, 5}, 5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Mode(tt.values)
			if result != tt.expected {
				t.Errorf("expected %d, got %d", tt.expected, result)
			}
		})
	}
}

func TestVariance(t *testing.T) {
	values := []float64{2, 4, 4, 4, 5, 5, 7, 9}
	expected := 4.0

	result := Variance(values)
	if !floatEqual(result, expected) {
		t.Errorf("expected %f, got %f", expected, result)
	}
}

func TestSampleVariance(t *testing.T) {
	values := []float64{2, 4, 4, 4, 5, 5, 7, 9}
	expected := 4.571428 // approximately

	result := SampleVariance(values)
	if !floatEqual(result, expected) {
		t.Errorf("expected %f, got %f", expected, result)
	}
}

func TestStdDev(t *testing.T) {
	values := []float64{2, 4, 4, 4, 5, 5, 7, 9}
	expected := 2.0

	result := StdDev(values)
	if !floatEqual(result, expected) {
		t.Errorf("expected %f, got %f", expected, result)
	}
}

func TestSampleStdDev(t *testing.T) {
	values := []float64{2, 4, 4, 4, 5, 5, 7, 9}
	expected := 2.138 // approximately

	result := SampleStdDev(values)
	if !floatEqual(result, expected) {
		t.Errorf("expected %f, got %f", expected, result)
	}
}

func TestMin(t *testing.T) {
	values := []int{5, 2, 8, 1, 9}
	expected := 1

	result := Min(values)
	if result != expected {
		t.Errorf("expected %d, got %d", expected, result)
	}
}

func TestMin_Panics(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic for empty slice")
		}
	}()
	Min([]int{})
}

func TestMax(t *testing.T) {
	values := []int{5, 2, 8, 1, 9}
	expected := 9

	result := Max(values)
	if result != expected {
		t.Errorf("expected %d, got %d", expected, result)
	}
}

func TestMax_Panics(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic for empty slice")
		}
	}()
	Max([]int{})
}

func TestSum(t *testing.T) {
	tests := []struct {
		name     string
		values   []int
		expected int
	}{
		{"empty", []int{}, 0},
		{"single", []int{5}, 5},
		{"multiple", []int{1, 2, 3, 4, 5}, 15},
		{"negative", []int{-1, -2, -3}, -6},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Sum(tt.values)
			if result != tt.expected {
				t.Errorf("expected %d, got %d", tt.expected, result)
			}
		})
	}
}

func TestPercentile(t *testing.T) {
	values := []float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

	tests := []struct {
		percentile float64
		expected   float64
	}{
		{0, 1},
		{25, 3.25},
		{50, 5.5},
		{75, 7.75},
		{100, 10},
	}

	for _, tt := range tests {
		result := Percentile(values, tt.percentile)
		if !floatEqual(result, tt.expected) {
			t.Errorf("P%v: expected %f, got %f", tt.percentile, tt.expected, result)
		}
	}
}

func TestPercentile_Panics(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic for invalid percentile")
		}
	}()
	Percentile([]float64{1, 2, 3}, 150)
}

func TestQuartiles(t *testing.T) {
	values := []float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

	q1, q2, q3 := Quartiles(values)

	if !floatEqual(q1, 3.25) {
		t.Errorf("Q1: expected 3.25, got %f", q1)
	}
	if !floatEqual(q2, 5.5) {
		t.Errorf("Q2: expected 5.5, got %f", q2)
	}
	if !floatEqual(q3, 7.75) {
		t.Errorf("Q3: expected 7.75, got %f", q3)
	}
}

func TestIQR(t *testing.T) {
	values := []float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	expected := 4.5 // 7.75 - 3.25

	result := IQR(values)
	if !floatEqual(result, expected) {
		t.Errorf("expected %f, got %f", expected, result)
	}
}

func TestRange(t *testing.T) {
	values := []int{1, 5, 3, 9, 2}
	expected := 8 // 9 - 1

	result := Range(values)
	if result != expected {
		t.Errorf("expected %d, got %d", expected, result)
	}
}

func TestCorrelation(t *testing.T) {
	tests := []struct {
		name     string
		x        []float64
		y        []float64
		expected float64
	}{
		{
			"perfect positive",
			[]float64{1, 2, 3, 4, 5},
			[]float64{2, 4, 6, 8, 10},
			1.0,
		},
		{
			"perfect negative",
			[]float64{1, 2, 3, 4, 5},
			[]float64{10, 8, 6, 4, 2},
			-1.0,
		},
		{
			"no correlation",
			[]float64{1, 2, 3, 4, 5},
			[]float64{5, 5, 5, 5, 5},
			0.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Correlation(tt.x, tt.y)
			if !floatEqual(result, tt.expected) {
				t.Errorf("expected %f, got %f", tt.expected, result)
			}
		})
	}
}

func TestCorrelation_DifferentLengths(t *testing.T) {
	x := []float64{1, 2, 3}
	y := []float64{1, 2}

	result := Correlation(x, y)
	if result != 0 {
		t.Errorf("expected 0 for different length slices, got %f", result)
	}
}

func BenchmarkMean(b *testing.B) {
	values := make([]float64, 1000)
	for i := range values {
		values[i] = float64(i)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Mean(values)
	}
}

func BenchmarkMedian(b *testing.B) {
	values := make([]float64, 1000)
	for i := range values {
		values[i] = float64(i)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Median(values)
	}
}

func BenchmarkStdDev(b *testing.B) {
	values := make([]float64, 1000)
	for i := range values {
		values[i] = float64(i)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		StdDev(values)
	}
}
