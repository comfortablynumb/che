# chestring - String Utilities

Comprehensive string manipulation utilities for Go, including case conversions, transformations, and common string operations.

## Features

- **Case Conversions**: camelCase, PascalCase, snake_case, kebab-case, SCREAMING_SNAKE_CASE
- **Transformations**: Capitalize, Uncapitalize, Reverse
- **Validation**: IsEmpty, IsBlank, IsNotEmpty, IsNotBlank
- **Truncation**: Truncate by length or words
- **Search**: ContainsAny, ContainsAll
- **Utilities**: Repeat, RemoveWhitespace, DefaultIfEmpty, SplitAndTrim
- **Zero Dependencies**: Only uses Go standard library

## Installation

```bash
go get github.com/comfortablynumb/che/pkg/chestring
```

## Quick Start

```go
package main

import (
    "fmt"
    "github.com/comfortablynumb/che/pkg/chestring"
)

func main() {
    // Case conversions
    fmt.Println(chestring.ToCamelCase("hello_world"))        // helloWorld
    fmt.Println(chestring.ToPascalCase("hello_world"))       // HelloWorld
    fmt.Println(chestring.ToSnakeCase("HelloWorld"))         // hello_world
    fmt.Println(chestring.ToKebabCase("HelloWorld"))         // hello-world
    fmt.Println(chestring.ToScreamingSnakeCase("helloWorld")) // HELLO_WORLD

    // Validation
    fmt.Println(chestring.IsBlank("   "))                    // true
    fmt.Println(chestring.IsNotEmpty("hello"))               // true

    // Truncation
    fmt.Println(chestring.Truncate("hello world", 8))        // hello...
    fmt.Println(chestring.TruncateWords("hello world foo", 2)) // hello world...
}
```

## Usage

### Case Conversions

#### ToCamelCase

Converts strings to camelCase:

```go
chestring.ToCamelCase("hello_world")    // "helloWorld"
chestring.ToCamelCase("HelloWorld")     // "helloWorld"
chestring.ToCamelCase("hello-world")    // "helloWorld"
chestring.ToCamelCase("hello world")    // "helloWorld"
```

#### ToPascalCase

Converts strings to PascalCase:

```go
chestring.ToPascalCase("hello_world")   // "HelloWorld"
chestring.ToPascalCase("helloWorld")    // "HelloWorld"
chestring.ToPascalCase("hello-world")   // "HelloWorld"
```

#### ToSnakeCase

Converts strings to snake_case:

```go
chestring.ToSnakeCase("HelloWorld")     // "hello_world"
chestring.ToSnakeCase("helloWorld")     // "hello_world"
chestring.ToSnakeCase("HTTPServer")     // "http_server"
```

#### ToKebabCase

Converts strings to kebab-case:

```go
chestring.ToKebabCase("HelloWorld")     // "hello-world"
chestring.ToKebabCase("hello_world")    // "hello-world"
```

#### ToScreamingSnakeCase

Converts strings to SCREAMING_SNAKE_CASE:

```go
chestring.ToScreamingSnakeCase("helloWorld")  // "HELLO_WORLD"
chestring.ToScreamingSnakeCase("hello-world") // "HELLO_WORLD"
```

### String Transformations

#### Capitalize / Uncapitalize

```go
chestring.Capitalize("hello")     // "Hello"
chestring.Uncapitalize("Hello")   // "hello"
```

#### Reverse

```go
chestring.Reverse("hello")        // "olleh"
chestring.Reverse("Hello World")  // "dlroW olleH"
```

### Validation

```go
chestring.IsEmpty("")             // true
chestring.IsEmpty(" ")            // false

chestring.IsBlank("")             // true
chestring.IsBlank("   ")          // true
chestring.IsBlank("\t\n")         // true

chestring.IsNotEmpty("hello")     // true
chestring.IsNotBlank(" hello ")   // true
```

### Truncation

#### Truncate by Length

```go
chestring.Truncate("hello world", 20)  // "hello world"
chestring.Truncate("hello world", 8)   // "hello..."
chestring.Truncate("hello world", 5)   // "he..."
```

#### Truncate by Words

```go
chestring.TruncateWords("hello world foo bar", 5)  // "hello world foo bar"
chestring.TruncateWords("hello world foo bar", 3)  // "hello world foo..."
chestring.TruncateWords("hello world foo bar", 1)  // "hello..."
```

### Search Operations

#### ContainsAny

Checks if string contains any of the substrings:

```go
chestring.ContainsAny("hello world", "hello", "foo")  // true
chestring.ContainsAny("hello world", "foo", "bar")    // false
```

#### ContainsAll

Checks if string contains all of the substrings:

```go
chestring.ContainsAll("hello world", "hello", "world")  // true
chestring.ContainsAll("hello world", "hello", "foo")    // false
```

### Utility Functions

#### Repeat

```go
chestring.Repeat("ab", 3)    // "ababab"
chestring.Repeat("x", 5)     // "xxxxx"
```

#### RemoveWhitespace

```go
chestring.RemoveWhitespace("hello world")         // "helloworld"
chestring.RemoveWhitespace("  hello  world  ")    // "helloworld"
chestring.RemoveWhitespace("hello\tworld\n")      // "helloworld"
```

#### DefaultIfEmpty / DefaultIfBlank

```go
chestring.DefaultIfEmpty("", "default")        // "default"
chestring.DefaultIfEmpty("hello", "default")   // "hello"

chestring.DefaultIfBlank("", "default")        // "default"
chestring.DefaultIfBlank("  ", "default")      // "default"
chestring.DefaultIfBlank("hello", "default")   // "hello"
```

#### SplitAndTrim

Splits string and trims whitespace from each part:

```go
parts := chestring.SplitAndTrim("a, b, c", ",")
// []string{"a", "b", "c"}

parts = chestring.SplitAndTrim("a,  ,c", ",")
// []string{"a", "c"} (empty parts removed)
```

## Examples

### API Field Name Conversion

```go
package main

import (
    "fmt"
    "github.com/comfortablynumb/che/pkg/chestring"
)

type APIResponse struct {
    UserID   int
    UserName string
}

func convertFieldNames(structName string) map[string]string {
    fields := map[string]string{
        "UserID":   chestring.ToSnakeCase("UserID"),     // user_id
        "UserName": chestring.ToCamelCase("user_name"),  // userName
    }
    return fields
}
```

### Configuration Keys

```go
package main

import (
    "github.com/comfortablynumb/che/pkg/chestring"
)

func generateConfigKey(name string) string {
    // Convert to screaming snake case for environment variables
    return chestring.ToScreamingSnakeCase(name)
}

func main() {
    key := generateConfigKey("databaseUrl")  // "DATABASE_URL"
    key = generateConfigKey("maxConnections") // "MAX_CONNECTIONS"
}
```

### String Sanitization

```go
package main

import (
    "github.com/comfortablynumb/che/pkg/chestring"
)

func sanitizeInput(input string) string {
    // Return default if blank
    if chestring.IsBlank(input) {
        return "N/A"
    }

    // Truncate long strings
    return chestring.Truncate(input, 100)
}

func formatDescription(desc string) string {
    // Truncate to 10 words for preview
    return chestring.TruncateWords(desc, 10)
}
```

### URL Slugs

```go
package main

import (
    "github.com/comfortablynumb/che/pkg/chestring"
)

func createSlug(title string) string {
    // Convert to kebab-case for URL-friendly slugs
    return chestring.ToKebabCase(title)
}

func main() {
    slug := createSlug("Hello World")        // "hello-world"
    slug = createSlug("My Awesome Post")     // "my-awesome-post"
}
```

### CSV Parsing

```go
package main

import (
    "github.com/comfortablynumb/che/pkg/chestring"
)

func parseCSV(line string) []string {
    return chestring.SplitAndTrim(line, ",")
}

func main() {
    values := parseCSV("apple, banana, cherry")
    // []string{"apple", "banana", "cherry"}
}
```

## API Reference

### Case Conversions
- `ToCamelCase(s string) string` - Convert to camelCase
- `ToPascalCase(s string) string` - Convert to PascalCase
- `ToSnakeCase(s string) string` - Convert to snake_case
- `ToKebabCase(s string) string` - Convert to kebab-case
- `ToScreamingSnakeCase(s string) string` - Convert to SCREAMING_SNAKE_CASE

### Transformations
- `Capitalize(s string) string` - Capitalize first letter
- `Uncapitalize(s string) string` - Lowercase first letter
- `Reverse(s string) string` - Reverse string

### Validation
- `IsEmpty(s string) bool` - Check if empty
- `IsBlank(s string) bool` - Check if empty or whitespace only
- `IsNotEmpty(s string) bool` - Check if not empty
- `IsNotBlank(s string) bool` - Check if not blank

### Truncation
- `Truncate(s string, maxLen int) string` - Truncate to max length
- `TruncateWords(s string, maxWords int) string` - Truncate to max words

### Search
- `ContainsAny(s string, substrs ...string) bool` - Contains any substring
- `ContainsAll(s string, substrs ...string) bool` - Contains all substrings

### Utilities
- `Repeat(s string, n int) string` - Repeat string n times
- `RemoveWhitespace(s string) string` - Remove all whitespace
- `DefaultIfEmpty(s, defaultValue string) string` - Return default if empty
- `DefaultIfBlank(s, defaultValue string) string` - Return default if blank
- `SplitAndTrim(s, sep string) []string` - Split and trim each part

## Related Packages

- **[cheenv](../cheenv)** - Environment variable utilities
- **[cheslice](../cheslice)** - Slice utility functions
- **[chemap](../chemap)** - Map utility functions

## License

This package is part of the Che library and shares the same license.
