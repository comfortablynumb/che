package chestring

import (
	"strings"
	"unicode"
)

// ToCamelCase converts a string to camelCase
// Example: "hello_world" -> "helloWorld"
func ToCamelCase(s string) string {
	if s == "" {
		return s
	}

	words := splitWords(s)
	if len(words) == 0 {
		return ""
	}

	result := strings.ToLower(words[0])
	for i := 1; i < len(words); i++ {
		result += capitalize(words[i])
	}

	return result
}

// ToPascalCase converts a string to PascalCase
// Example: "hello_world" -> "HelloWorld"
func ToPascalCase(s string) string {
	if s == "" {
		return s
	}

	words := splitWords(s)
	var result string
	for _, word := range words {
		result += capitalize(word)
	}

	return result
}

// ToSnakeCase converts a string to snake_case
// Example: "HelloWorld" -> "hello_world"
func ToSnakeCase(s string) string {
	if s == "" {
		return s
	}

	words := splitWords(s)
	return strings.ToLower(strings.Join(words, "_"))
}

// ToKebabCase converts a string to kebab-case
// Example: "HelloWorld" -> "hello-world"
func ToKebabCase(s string) string {
	if s == "" {
		return s
	}

	words := splitWords(s)
	return strings.ToLower(strings.Join(words, "-"))
}

// ToScreamingSnakeCase converts a string to SCREAMING_SNAKE_CASE
// Example: "helloWorld" -> "HELLO_WORLD"
func ToScreamingSnakeCase(s string) string {
	if s == "" {
		return s
	}

	words := splitWords(s)
	return strings.ToUpper(strings.Join(words, "_"))
}

// Capitalize capitalizes the first letter of a string
// Example: "hello" -> "Hello"
func Capitalize(s string) string {
	if s == "" {
		return s
	}

	runes := []rune(s)
	runes[0] = unicode.ToUpper(runes[0])
	return string(runes)
}

// Uncapitalize lowercases the first letter of a string
// Example: "Hello" -> "hello"
func Uncapitalize(s string) string {
	if s == "" {
		return s
	}

	runes := []rune(s)
	runes[0] = unicode.ToLower(runes[0])
	return string(runes)
}

// Reverse reverses a string
// Example: "hello" -> "olleh"
func Reverse(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

// IsEmpty checks if a string is empty
func IsEmpty(s string) bool {
	return len(s) == 0
}

// IsBlank checks if a string is empty or contains only whitespace
func IsBlank(s string) bool {
	return len(strings.TrimSpace(s)) == 0
}

// IsNotEmpty checks if a string is not empty
func IsNotEmpty(s string) bool {
	return !IsEmpty(s)
}

// IsNotBlank checks if a string is not blank
func IsNotBlank(s string) bool {
	return !IsBlank(s)
}

// Truncate truncates a string to a maximum length
// If the string is longer, it adds an ellipsis
func Truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}

	if maxLen <= 3 {
		return s[:maxLen]
	}

	return s[:maxLen-3] + "..."
}

// TruncateWords truncates a string to a maximum number of words
func TruncateWords(s string, maxWords int) string {
	words := strings.Fields(s)
	if len(words) <= maxWords {
		return s
	}

	return strings.Join(words[:maxWords], " ") + "..."
}

// ContainsAny checks if a string contains any of the substrings
func ContainsAny(s string, substrs ...string) bool {
	for _, substr := range substrs {
		if strings.Contains(s, substr) {
			return true
		}
	}
	return false
}

// ContainsAll checks if a string contains all of the substrings
func ContainsAll(s string, substrs ...string) bool {
	for _, substr := range substrs {
		if !strings.Contains(s, substr) {
			return false
		}
	}
	return true
}

// Repeat repeats a string n times
func Repeat(s string, n int) string {
	return strings.Repeat(s, n)
}

// RemoveWhitespace removes all whitespace from a string
func RemoveWhitespace(s string) string {
	return strings.Map(func(r rune) rune {
		if unicode.IsSpace(r) {
			return -1
		}
		return r
	}, s)
}

// DefaultIfEmpty returns the default value if the string is empty
func DefaultIfEmpty(s, defaultValue string) string {
	if IsEmpty(s) {
		return defaultValue
	}
	return s
}

// DefaultIfBlank returns the default value if the string is blank
func DefaultIfBlank(s, defaultValue string) string {
	if IsBlank(s) {
		return defaultValue
	}
	return s
}

// SplitAndTrim splits a string by delimiter and trims each part
func SplitAndTrim(s, sep string) []string {
	parts := strings.Split(s, sep)
	result := make([]string, 0, len(parts))
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}

// splitWords splits a string into words considering various naming conventions
func splitWords(s string) []string {
	var words []string
	var currentWord strings.Builder

	runes := []rune(s)
	for i := 0; i < len(runes); i++ {
		r := runes[i]

		// Handle delimiters
		if r == '_' || r == '-' || r == ' ' || r == '.' {
			if currentWord.Len() > 0 {
				words = append(words, currentWord.String())
				currentWord.Reset()
			}
			continue
		}

		// Handle camelCase/PascalCase transitions
		if unicode.IsUpper(r) {
			// Check if this is a transition from lowercase to uppercase
			if currentWord.Len() > 0 && i > 0 && unicode.IsLower(runes[i-1]) {
				words = append(words, currentWord.String())
				currentWord.Reset()
			}

			// Check if this is the end of an acronym (e.g., "HTTPServer" -> "HTTP", "Server")
			if i+1 < len(runes) && unicode.IsLower(runes[i+1]) && currentWord.Len() > 0 {
				words = append(words, currentWord.String())
				currentWord.Reset()
			}
		}

		currentWord.WriteRune(r)
	}

	if currentWord.Len() > 0 {
		words = append(words, currentWord.String())
	}

	return words
}

// capitalize capitalizes the first letter and lowercases the rest
func capitalize(s string) string {
	if s == "" {
		return s
	}

	runes := []rune(strings.ToLower(s))
	runes[0] = unicode.ToUpper(runes[0])
	return string(runes)
}
