package chestring

import (
	"testing"

	"github.com/comfortablynumb/che/pkg/chetest"
)

func TestToCamelCase(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"hello_world", "helloWorld"},
		{"hello-world", "helloWorld"},
		{"hello world", "helloWorld"},
		{"HelloWorld", "helloWorld"},
		{"helloWorld", "helloWorld"},
		{"HELLO_WORLD", "helloWorld"},
		{"", ""},
		{"a", "a"},
		{"HTTP_SERVER", "httpServer"},
	}

	for _, tt := range tests {
		result := ToCamelCase(tt.input)
		chetest.RequireEqual(t, result, tt.expected)
	}
}

func TestToPascalCase(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"hello_world", "HelloWorld"},
		{"hello-world", "HelloWorld"},
		{"hello world", "HelloWorld"},
		{"helloWorld", "HelloWorld"},
		{"HELLO_WORLD", "HelloWorld"},
		{"", ""},
		{"a", "A"},
	}

	for _, tt := range tests {
		result := ToPascalCase(tt.input)
		chetest.RequireEqual(t, result, tt.expected)
	}
}

func TestToSnakeCase(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"helloWorld", "hello_world"},
		{"HelloWorld", "hello_world"},
		{"hello-world", "hello_world"},
		{"hello world", "hello_world"},
		{"hello_world", "hello_world"},
		{"HTTPServer", "http_server"},
		{"", ""},
		{"A", "a"},
	}

	for _, tt := range tests {
		result := ToSnakeCase(tt.input)
		chetest.RequireEqual(t, result, tt.expected)
	}
}

func TestToKebabCase(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"helloWorld", "hello-world"},
		{"HelloWorld", "hello-world"},
		{"hello_world", "hello-world"},
		{"hello world", "hello-world"},
		{"hello-world", "hello-world"},
		{"", ""},
	}

	for _, tt := range tests {
		result := ToKebabCase(tt.input)
		chetest.RequireEqual(t, result, tt.expected)
	}
}

func TestToScreamingSnakeCase(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"helloWorld", "HELLO_WORLD"},
		{"HelloWorld", "HELLO_WORLD"},
		{"hello-world", "HELLO_WORLD"},
		{"hello_world", "HELLO_WORLD"},
		{"", ""},
	}

	for _, tt := range tests {
		result := ToScreamingSnakeCase(tt.input)
		chetest.RequireEqual(t, result, tt.expected)
	}
}

func TestCapitalize(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"hello", "Hello"},
		{"Hello", "Hello"},
		{"HELLO", "HELLO"},
		{"h", "H"},
		{"", ""},
	}

	for _, tt := range tests {
		result := Capitalize(tt.input)
		chetest.RequireEqual(t, result, tt.expected)
	}
}

func TestUncapitalize(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"Hello", "hello"},
		{"hello", "hello"},
		{"HELLO", "hELLO"},
		{"H", "h"},
		{"", ""},
	}

	for _, tt := range tests {
		result := Uncapitalize(tt.input)
		chetest.RequireEqual(t, result, tt.expected)
	}
}

func TestReverse(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"hello", "olleh"},
		{"Hello World", "dlroW olleH"},
		{"12345", "54321"},
		{"a", "a"},
		{"", ""},
	}

	for _, tt := range tests {
		result := Reverse(tt.input)
		chetest.RequireEqual(t, result, tt.expected)
	}
}

func TestIsEmpty(t *testing.T) {
	chetest.RequireEqual(t, IsEmpty(""), true)
	chetest.RequireEqual(t, IsEmpty(" "), false)
	chetest.RequireEqual(t, IsEmpty("hello"), false)
}

func TestIsBlank(t *testing.T) {
	chetest.RequireEqual(t, IsBlank(""), true)
	chetest.RequireEqual(t, IsBlank(" "), true)
	chetest.RequireEqual(t, IsBlank("  \t\n  "), true)
	chetest.RequireEqual(t, IsBlank("hello"), false)
	chetest.RequireEqual(t, IsBlank(" hello "), false)
}

func TestIsNotEmpty(t *testing.T) {
	chetest.RequireEqual(t, IsNotEmpty(""), false)
	chetest.RequireEqual(t, IsNotEmpty(" "), true)
	chetest.RequireEqual(t, IsNotEmpty("hello"), true)
}

func TestIsNotBlank(t *testing.T) {
	chetest.RequireEqual(t, IsNotBlank(""), false)
	chetest.RequireEqual(t, IsNotBlank(" "), false)
	chetest.RequireEqual(t, IsNotBlank("  \t\n  "), false)
	chetest.RequireEqual(t, IsNotBlank("hello"), true)
}

func TestTruncate(t *testing.T) {
	chetest.RequireEqual(t, Truncate("hello world", 20), "hello world")
	chetest.RequireEqual(t, Truncate("hello world", 8), "hello...")
	chetest.RequireEqual(t, Truncate("hello world", 5), "he...")
	chetest.RequireEqual(t, Truncate("hello", 3), "hel")
	chetest.RequireEqual(t, Truncate("", 5), "")
}

func TestTruncateWords(t *testing.T) {
	chetest.RequireEqual(t, TruncateWords("hello world foo bar", 5), "hello world foo bar")
	chetest.RequireEqual(t, TruncateWords("hello world foo bar", 3), "hello world foo...")
	chetest.RequireEqual(t, TruncateWords("hello world foo bar", 1), "hello...")
	chetest.RequireEqual(t, TruncateWords("", 5), "")
}

func TestContainsAny(t *testing.T) {
	chetest.RequireEqual(t, ContainsAny("hello world", "hello", "foo"), true)
	chetest.RequireEqual(t, ContainsAny("hello world", "world", "foo"), true)
	chetest.RequireEqual(t, ContainsAny("hello world", "foo", "bar"), false)
	chetest.RequireEqual(t, ContainsAny("hello world"), false)
}

func TestContainsAll(t *testing.T) {
	chetest.RequireEqual(t, ContainsAll("hello world", "hello", "world"), true)
	chetest.RequireEqual(t, ContainsAll("hello world", "hello", "foo"), false)
	chetest.RequireEqual(t, ContainsAll("hello world", "foo", "bar"), false)
	chetest.RequireEqual(t, ContainsAll("hello world"), true)
}

func TestRepeat(t *testing.T) {
	chetest.RequireEqual(t, Repeat("ab", 3), "ababab")
	chetest.RequireEqual(t, Repeat("x", 5), "xxxxx")
	chetest.RequireEqual(t, Repeat("hello", 0), "")
	chetest.RequireEqual(t, Repeat("", 5), "")
}

func TestRemoveWhitespace(t *testing.T) {
	chetest.RequireEqual(t, RemoveWhitespace("hello world"), "helloworld")
	chetest.RequireEqual(t, RemoveWhitespace("  hello  world  "), "helloworld")
	chetest.RequireEqual(t, RemoveWhitespace("hello\tworld\n"), "helloworld")
	chetest.RequireEqual(t, RemoveWhitespace(""), "")
}

func TestDefaultIfEmpty(t *testing.T) {
	chetest.RequireEqual(t, DefaultIfEmpty("", "default"), "default")
	chetest.RequireEqual(t, DefaultIfEmpty(" ", "default"), " ")
	chetest.RequireEqual(t, DefaultIfEmpty("hello", "default"), "hello")
}

func TestDefaultIfBlank(t *testing.T) {
	chetest.RequireEqual(t, DefaultIfBlank("", "default"), "default")
	chetest.RequireEqual(t, DefaultIfBlank(" ", "default"), "default")
	chetest.RequireEqual(t, DefaultIfBlank("  \t\n  ", "default"), "default")
	chetest.RequireEqual(t, DefaultIfBlank("hello", "default"), "hello")
}

func TestSplitAndTrim(t *testing.T) {
	result := SplitAndTrim("a, b, c", ",")
	chetest.RequireEqual(t, len(result), 3)
	chetest.RequireEqual(t, result[0], "a")
	chetest.RequireEqual(t, result[1], "b")
	chetest.RequireEqual(t, result[2], "c")

	result = SplitAndTrim("a,  ,c", ",")
	chetest.RequireEqual(t, len(result), 2)
	chetest.RequireEqual(t, result[0], "a")
	chetest.RequireEqual(t, result[1], "c")

	result = SplitAndTrim("", ",")
	chetest.RequireEqual(t, len(result), 0)
}

func TestSplitWords(t *testing.T) {
	tests := []struct {
		input    string
		expected []string
	}{
		{"helloWorld", []string{"hello", "World"}},
		{"HelloWorld", []string{"Hello", "World"}},
		{"hello_world", []string{"hello", "world"}},
		{"hello-world", []string{"hello", "world"}},
		{"HTTPServer", []string{"HTTP", "Server"}},
		{"getHTTPResponseCode", []string{"get", "HTTP", "Response", "Code"}},
	}

	for _, tt := range tests {
		result := splitWords(tt.input)
		chetest.RequireEqual(t, len(result), len(tt.expected))
		for i := range result {
			chetest.RequireEqual(t, result[i], tt.expected[i])
		}
	}
}
