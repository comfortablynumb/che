package chestrsim

import (
	"math"
	"testing"
)

func floatEqual(a, b, tolerance float64) bool {
	return math.Abs(a-b) < tolerance
}

func TestLevenshtein(t *testing.T) {
	tests := []struct {
		s1       string
		s2       string
		expected int
	}{
		{"", "", 0},
		{"", "abc", 3},
		{"abc", "", 3},
		{"abc", "abc", 0},
		{"abc", "abd", 1},
		{"kitten", "sitting", 3},
		{"Saturday", "Sunday", 3},
		{"Hello", "Hallo", 1},
		{"café", "cafe", 1}, // Unicode test
	}

	for _, tt := range tests {
		result := Levenshtein(tt.s1, tt.s2)
		if result != tt.expected {
			t.Errorf("Levenshtein(%q, %q) = %d, expected %d", tt.s1, tt.s2, result, tt.expected)
		}
	}
}

func TestLevenshteinSimilarity(t *testing.T) {
	tests := []struct {
		s1       string
		s2       string
		expected float64
	}{
		{"abc", "abc", 1.0},
		{"", "", 1.0},
		{"abc", "xyz", 0.0},
		{"abc", "ab", 0.666},
	}

	for _, tt := range tests {
		result := LevenshteinSimilarity(tt.s1, tt.s2)
		if !floatEqual(result, tt.expected, 0.01) {
			t.Errorf("LevenshteinSimilarity(%q, %q) = %f, expected %f", tt.s1, tt.s2, result, tt.expected)
		}
	}
}

func TestHamming(t *testing.T) {
	tests := []struct {
		s1       string
		s2       string
		expected int
	}{
		{"", "", 0},
		{"abc", "abc", 0},
		{"abc", "abd", 1},
		{"abc", "xyz", 3},
		{"1011101", "1001001", 2},
		{"abc", "ab", -1}, // Different lengths
	}

	for _, tt := range tests {
		result := Hamming(tt.s1, tt.s2)
		if result != tt.expected {
			t.Errorf("Hamming(%q, %q) = %d, expected %d", tt.s1, tt.s2, result, tt.expected)
		}
	}
}

func TestHammingSimilarity(t *testing.T) {
	tests := []struct {
		s1       string
		s2       string
		expected float64
	}{
		{"abc", "abc", 1.0},
		{"abc", "xyz", 0.0},
		{"abc", "ab", -1.0}, // Different lengths
	}

	for _, tt := range tests {
		result := HammingSimilarity(tt.s1, tt.s2)
		if !floatEqual(result, tt.expected, 0.01) {
			t.Errorf("HammingSimilarity(%q, %q) = %f, expected %f", tt.s1, tt.s2, result, tt.expected)
		}
	}
}

func TestJaro(t *testing.T) {
	tests := []struct {
		s1       string
		s2       string
		expected float64
	}{
		{"", "", 1.0},
		{"abc", "abc", 1.0},
		{"martha", "marhta", 0.944},
		{"dixon", "dicksonx", 0.767},
		{"", "abc", 0.0},
	}

	for _, tt := range tests {
		result := Jaro(tt.s1, tt.s2)
		if !floatEqual(result, tt.expected, 0.01) {
			t.Errorf("Jaro(%q, %q) = %f, expected %f", tt.s1, tt.s2, result, tt.expected)
		}
	}
}

func TestJaroWinkler(t *testing.T) {
	tests := []struct {
		s1       string
		s2       string
		expected float64
	}{
		{"", "", 1.0},
		{"abc", "abc", 1.0},
		{"martha", "marhta", 0.961},
		{"dwayne", "duane", 0.84},
		{"dixon", "dicksonx", 0.813},
	}

	for _, tt := range tests {
		result := JaroWinkler(tt.s1, tt.s2)
		if !floatEqual(result, tt.expected, 0.01) {
			t.Errorf("JaroWinkler(%q, %q) = %f, expected %f", tt.s1, tt.s2, result, tt.expected)
		}
	}
}

func TestCosine(t *testing.T) {
	tests := []struct {
		s1       string
		s2       string
		expected float64
	}{
		{"", "", 1.0},
		{"abc", "abc", 1.0},
		{"abc", "xyz", 0.0},
		{"hello", "hallo", 0.5},
	}

	for _, tt := range tests {
		result := Cosine(tt.s1, tt.s2)
		if !floatEqual(result, tt.expected, 0.15) {
			t.Errorf("Cosine(%q, %q) = %f, expected ~%f", tt.s1, tt.s2, result, tt.expected)
		}
	}
}

func TestJaccard(t *testing.T) {
	tests := []struct {
		s1       string
		s2       string
		expected float64
	}{
		{"", "", 1.0},
		{"abc", "abc", 1.0},
		{"abc", "xyz", 0.0},
		{"hello", "hallo", 0.4},
	}

	for _, tt := range tests {
		result := Jaccard(tt.s1, tt.s2)
		if !floatEqual(result, tt.expected, 0.15) {
			t.Errorf("Jaccard(%q, %q) = %f, expected ~%f", tt.s1, tt.s2, result, tt.expected)
		}
	}
}

func TestFuzzyMatch(t *testing.T) {
	tests := []struct {
		query    string
		target   string
		expected bool
	}{
		{"", "", true},
		{"", "abc", true},
		{"abc", "abc", true},
		{"abc", "aabbcc", true},
		{"abc", "axbxcx", true},
		{"abc", "xyz", false},
		{"abc", "acb", false}, // Order matters
		{"fb", "FooBar", true}, // Case insensitive
		{"fb", "foobar", true},
	}

	for _, tt := range tests {
		result := FuzzyMatch(tt.query, tt.target)
		if result != tt.expected {
			t.Errorf("FuzzyMatch(%q, %q) = %v, expected %v", tt.query, tt.target, result, tt.expected)
		}
	}
}

func TestFuzzyScore(t *testing.T) {
	tests := []struct {
		query    string
		target   string
		minScore float64
	}{
		{"abc", "abc", 0.9},       // Consecutive matches = high score
		{"abc", "aabbcc", 0.5},    // Non-consecutive but close
		{"abc", "axbxcx", 0.3},    // Spread out
		{"abc", "xyz", 0.0},       // No match
		{"fb", "FooBar", 0.5},     // Case insensitive
		{"", "anything", 1.0},     // Empty query
	}

	for _, tt := range tests {
		result := FuzzyScore(tt.query, tt.target)
		if result < tt.minScore-0.1 {
			t.Errorf("FuzzyScore(%q, %q) = %f, expected >= %f", tt.query, tt.target, result, tt.minScore)
		}
	}
}

func TestGetBigrams(t *testing.T) {
	tests := []struct {
		input    string
		expected map[string]int
	}{
		{"", map[string]int{}},
		{"a", map[string]int{}},
		{"ab", map[string]int{"ab": 1}},
		{"abc", map[string]int{"ab": 1, "bc": 1}},
		{"aaa", map[string]int{"aa": 2}},
	}

	for _, tt := range tests {
		result := getBigrams(tt.input)
		if len(result) != len(tt.expected) {
			t.Errorf("getBigrams(%q) returned %d bigrams, expected %d", tt.input, len(result), len(tt.expected))
			continue
		}

		for bigram, count := range tt.expected {
			if result[bigram] != count {
				t.Errorf("getBigrams(%q)[%q] = %d, expected %d", tt.input, bigram, result[bigram], count)
			}
		}
	}
}

func BenchmarkLevenshtein(b *testing.B) {
	s1 := "kitten"
	s2 := "sitting"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Levenshtein(s1, s2)
	}
}

func BenchmarkJaroWinkler(b *testing.B) {
	s1 := "martha"
	s2 := "marhta"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		JaroWinkler(s1, s2)
	}
}

func BenchmarkCosine(b *testing.B) {
	s1 := "hello world"
	s2 := "hallo world"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Cosine(s1, s2)
	}
}

func BenchmarkFuzzyMatch(b *testing.B) {
	query := "fb"
	target := "FooBar"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		FuzzyMatch(query, target)
	}
}
