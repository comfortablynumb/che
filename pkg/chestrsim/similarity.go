package chestrsim

import (
	"strings"
	"unicode/utf8"
)

// Levenshtein calculates the Levenshtein distance between two strings.
// The Levenshtein distance is the minimum number of single-character edits
// (insertions, deletions, or substitutions) required to change one string into another.
func Levenshtein(s1, s2 string) int {
	if s1 == s2 {
		return 0
	}

	// Convert to rune slices to handle Unicode properly
	r1 := []rune(s1)
	r2 := []rune(s2)

	len1 := len(r1)
	len2 := len(r2)

	if len1 == 0 {
		return len2
	}
	if len2 == 0 {
		return len1
	}

	// Create matrix
	matrix := make([][]int, len1+1)
	for i := range matrix {
		matrix[i] = make([]int, len2+1)
	}

	// Initialize first column and row
	for i := 0; i <= len1; i++ {
		matrix[i][0] = i
	}
	for j := 0; j <= len2; j++ {
		matrix[0][j] = j
	}

	// Fill matrix
	for i := 1; i <= len1; i++ {
		for j := 1; j <= len2; j++ {
			cost := 1
			if r1[i-1] == r2[j-1] {
				cost = 0
			}

			matrix[i][j] = min(
				min(matrix[i-1][j]+1, matrix[i][j-1]+1),
				matrix[i-1][j-1]+cost,
			)
		}
	}

	return matrix[len1][len2]
}

// LevenshteinSimilarity calculates a similarity score between 0 and 1 based on Levenshtein distance.
// Returns 1.0 for identical strings and approaches 0 as strings become more different.
func LevenshteinSimilarity(s1, s2 string) float64 {
	distance := Levenshtein(s1, s2)
	maxLen := max(utf8.RuneCountInString(s1), utf8.RuneCountInString(s2))

	if maxLen == 0 {
		return 1.0
	}

	return 1.0 - (float64(distance) / float64(maxLen))
}

// Hamming calculates the Hamming distance between two strings of equal length.
// The Hamming distance is the number of positions at which the corresponding characters differ.
// Returns -1 if strings have different lengths.
func Hamming(s1, s2 string) int {
	r1 := []rune(s1)
	r2 := []rune(s2)

	if len(r1) != len(r2) {
		return -1
	}

	distance := 0
	for i := range r1 {
		if r1[i] != r2[i] {
			distance++
		}
	}

	return distance
}

// HammingSimilarity calculates a similarity score between 0 and 1 based on Hamming distance.
// Returns -1 if strings have different lengths.
func HammingSimilarity(s1, s2 string) float64 {
	distance := Hamming(s1, s2)
	if distance == -1 {
		return -1
	}

	length := utf8.RuneCountInString(s1)
	if length == 0 {
		return 1.0
	}

	return 1.0 - (float64(distance) / float64(length))
}

// JaroWinkler calculates the Jaro-Winkler similarity between two strings.
// Returns a value between 0 (no similarity) and 1 (exact match).
// Particularly good for short strings like names.
func JaroWinkler(s1, s2 string) float64 {
	jaro := Jaro(s1, s2)

	// If the Jaro similarity is below the threshold, don't apply the Winkler modification
	if jaro < 0.7 {
		return jaro
	}

	// Find length of common prefix (up to 4 characters)
	r1 := []rune(s1)
	r2 := []rune(s2)
	prefixLen := 0
	maxPrefix := min(min(len(r1), len(r2)), 4)

	for i := 0; i < maxPrefix; i++ {
		if r1[i] == r2[i] {
			prefixLen++
		} else {
			break
		}
	}

	// Apply Winkler modification
	return jaro + float64(prefixLen)*0.1*(1.0-jaro)
}

// Jaro calculates the Jaro similarity between two strings.
// Returns a value between 0 (no similarity) and 1 (exact match).
func Jaro(s1, s2 string) float64 {
	if s1 == s2 {
		return 1.0
	}

	r1 := []rune(s1)
	r2 := []rune(s2)
	len1 := len(r1)
	len2 := len(r2)

	if len1 == 0 && len2 == 0 {
		return 1.0
	}
	if len1 == 0 || len2 == 0 {
		return 0.0
	}

	// Calculate match window
	matchWindow := max(len1, len2)/2 - 1
	if matchWindow < 0 {
		matchWindow = 0
	}

	s1Matches := make([]bool, len1)
	s2Matches := make([]bool, len2)

	matches := 0
	transpositions := 0

	// Find matches
	for i := 0; i < len1; i++ {
		start := max(0, i-matchWindow)
		end := min(i+matchWindow+1, len2)

		for j := start; j < end; j++ {
			if s2Matches[j] || r1[i] != r2[j] {
				continue
			}
			s1Matches[i] = true
			s2Matches[j] = true
			matches++
			break
		}
	}

	if matches == 0 {
		return 0.0
	}

	// Find transpositions
	k := 0
	for i := 0; i < len1; i++ {
		if !s1Matches[i] {
			continue
		}
		for !s2Matches[k] {
			k++
		}
		if r1[i] != r2[k] {
			transpositions++
		}
		k++
	}

	return (float64(matches)/float64(len1) +
		float64(matches)/float64(len2) +
		float64(matches-transpositions/2)/float64(matches)) / 3.0
}

// Cosine calculates the cosine similarity between two strings based on character bigrams.
// Returns a value between 0 (no similarity) and 1 (identical).
func Cosine(s1, s2 string) float64 {
	if s1 == s2 {
		return 1.0
	}

	bigrams1 := getBigrams(s1)
	bigrams2 := getBigrams(s2)

	if len(bigrams1) == 0 && len(bigrams2) == 0 {
		return 1.0
	}
	if len(bigrams1) == 0 || len(bigrams2) == 0 {
		return 0.0
	}

	// Calculate dot product
	dotProduct := 0.0
	for bigram := range bigrams1 {
		if count2, ok := bigrams2[bigram]; ok {
			dotProduct += float64(bigrams1[bigram] * count2)
		}
	}

	// Calculate magnitudes
	mag1 := 0.0
	for _, count := range bigrams1 {
		mag1 += float64(count * count)
	}

	mag2 := 0.0
	for _, count := range bigrams2 {
		mag2 += float64(count * count)
	}

	if mag1 == 0 || mag2 == 0 {
		return 0.0
	}

	return dotProduct / (sqrt(mag1) * sqrt(mag2))
}

// Jaccard calculates the Jaccard similarity coefficient between two strings.
// Returns a value between 0 (no similarity) and 1 (identical).
// Based on character bigrams.
func Jaccard(s1, s2 string) float64 {
	if s1 == s2 {
		return 1.0
	}

	bigrams1 := getBigrams(s1)
	bigrams2 := getBigrams(s2)

	if len(bigrams1) == 0 && len(bigrams2) == 0 {
		return 1.0
	}

	// Calculate intersection and union
	intersection := 0
	union := make(map[string]bool)

	for bigram := range bigrams1 {
		union[bigram] = true
		if _, ok := bigrams2[bigram]; ok {
			intersection++
		}
	}

	for bigram := range bigrams2 {
		union[bigram] = true
	}

	if len(union) == 0 {
		return 0.0
	}

	return float64(intersection) / float64(len(union))
}

// FuzzyMatch checks if a query fuzzy matches a target string.
// Returns true if all characters in query appear in order in target (case-insensitive).
func FuzzyMatch(query, target string) bool {
	query = strings.ToLower(query)
	target = strings.ToLower(target)

	qRunes := []rune(query)
	tRunes := []rune(target)

	qi := 0
	ti := 0

	for qi < len(qRunes) && ti < len(tRunes) {
		if qRunes[qi] == tRunes[ti] {
			qi++
		}
		ti++
	}

	return qi == len(qRunes)
}

// FuzzyScore calculates a fuzzy match score between 0 and 1.
// Higher scores indicate better matches (considers character proximity).
func FuzzyScore(query, target string) float64 {
	if !FuzzyMatch(query, target) {
		return 0.0
	}

	query = strings.ToLower(query)
	target = strings.ToLower(target)

	qRunes := []rune(query)
	tRunes := []rune(target)

	if len(qRunes) == 0 {
		return 1.0
	}

	score := 0.0
	qi := 0
	ti := 0
	lastMatchPos := -1

	for qi < len(qRunes) && ti < len(tRunes) {
		if qRunes[qi] == tRunes[ti] {
			// Award points based on proximity to last match
			if lastMatchPos == -1 {
				score += 1.0
			} else {
				gap := ti - lastMatchPos - 1
				if gap == 0 {
					score += 1.0 // Consecutive match
				} else {
					score += 1.0 / float64(gap+1)
				}
			}
			lastMatchPos = ti
			qi++
		}
		ti++
	}

	// Normalize by query length
	return score / float64(len(qRunes))
}

// Helper functions

func getBigrams(s string) map[string]int {
	runes := []rune(s)
	bigrams := make(map[string]int)

	if len(runes) < 2 {
		return bigrams
	}

	for i := 0; i < len(runes)-1; i++ {
		bigram := string(runes[i : i+2])
		bigrams[bigram]++
	}

	return bigrams
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func sqrt(x float64) float64 {
	if x == 0 {
		return 0
	}

	// Newton's method for square root
	z := x
	for i := 0; i < 10; i++ {
		z = z - (z*z-x)/(2*z)
	}
	return z
}
