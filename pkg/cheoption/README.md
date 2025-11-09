# cheoption

Functional programming patterns for Go: Optional and Result types.

## Features

- **Optional[T]**: Represents a value that may or may not be present
- **Result[T]**: Represents the result of an operation that can succeed or fail
- Type-safe with generics
- Functional operations: Map, FlatMap, Filter, AndThen
- Zero dependencies

## Installation

```bash
go get github.com/comfortablynumb/che/pkg/cheoption
```

## Usage

### Optional

```go
package main

import (
    "fmt"
    "github.com/comfortablynumb/che/pkg/cheoption"
)

func findUser(id int) cheoption.Optional[string] {
    if id == 1 {
        return cheoption.Some("Alice")
    }
    return cheoption.None[string]()
}

func main() {
    // Basic usage
    user := findUser(1)
    if user.IsPresent() {
        fmt.Println("Found:", user.Get())
    }

    // With default value
    name := findUser(999).GetOr("Guest")
    fmt.Println(name) // "Guest"

    // Map operation
    upper := findUser(1).Map(func(s string) string {
        return strings.ToUpper(s)
    })
    fmt.Println(upper.Get()) // "ALICE"

    // Filter
    admin := findUser(1).Filter(func(s string) bool {
        return s == "Admin"
    })
    fmt.Println(admin.IsEmpty()) // true

    // Chaining
    result := findUser(1).
        Map(func(s string) string { return "Hello, " + s }).
        GetOr("Hello, Guest")
    fmt.Println(result) // "Hello, Alice"
}
```

### Result

```go
package main

import (
    "errors"
    "fmt"
    "github.com/comfortablynumb/che/pkg/cheoption"
)

func divide(a, b float64) cheoption.Result[float64] {
    if b == 0 {
        return cheoption.Err[float64](errors.New("division by zero"))
    }
    return cheoption.Ok(a / b)
}

func main() {
    // Successful result
    result := divide(10, 2)
    if result.IsOk() {
        fmt.Println("Result:", result.Unwrap()) // "Result: 5"
    }

    // Error result
    result = divide(10, 0)
    if result.IsErr() {
        fmt.Println("Error:", result.Error()) // "Error: division by zero"
    }

    // With default value
    value := divide(10, 0).UnwrapOr(0)
    fmt.Println(value) // 0

    // Map operation
    squared := divide(10, 2).Map(func(x float64) float64 {
        return x * x
    })
    fmt.Println(squared.Unwrap()) // 25

    // Chaining operations
    result = divide(10, 2).
        FlatMap(func(x float64) cheoption.Result[float64] {
            return divide(x, 2)
        }).
        Map(func(x float64) float64 {
            return x * 10
        })
    fmt.Println(result.Unwrap()) // 25 (10/2/2*10)
}
```

## API

### Optional[T]

- `Some[T](value T) Optional[T]` - Create Optional with value
- `None[T]() Optional[T]` - Create empty Optional
- `IsPresent() bool` - Check if value exists
- `IsEmpty() bool` - Check if empty
- `Get() T` - Get value (panics if empty)
- `GetOr(defaultValue T) T` - Get value or default
- `GetOrElse(fn func() T) T` - Get value or call function
- `OrElse(other Optional[T]) Optional[T]` - Return this or other if empty
- `Map(fn func(T) T) Optional[T]` - Transform value
- `FlatMap(fn func(T) Optional[T]) Optional[T]` - Transform and flatten
- `Filter(predicate func(T) bool) Optional[T]` - Keep only if predicate matches
- `IfPresent(fn func(T))` - Execute function if present
- `IfPresentOrElse(presentFn func(T), emptyFn func())` - Execute based on presence
- `Unpack() (T, bool)` - Return value and presence flag

### Result[T]

- `Ok[T](value T) Result[T]` - Create successful result
- `Err[T](err error) Result[T]` - Create error result
- `IsOk() bool` - Check if successful
- `IsErr() bool` - Check if error
- `Unwrap() T` - Get value (panics if error)
- `UnwrapOr(defaultValue T) T` - Get value or default
- `UnwrapOrElse(fn func(error) T) T` - Get value or call function
- `Error() error` - Get error
- `Map(fn func(T) T) Result[T]` - Transform value if ok
- `MapErr(fn func(error) error) Result[T]` - Transform error
- `FlatMap(fn func(T) Result[T]) Result[T]` - Transform and flatten
- `AndThen(fn func(T) Result[T]) Result[T]` - Alias for FlatMap
- `Or(other Result[T]) Result[T]` - Return this or other if error
- `Unpack() (T, error)` - Return value and error

## License

MIT
