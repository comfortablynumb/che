# checolor

Terminal color and text formatting utilities for Go.

## Features

- ANSI color codes for terminal output
- Text styling (bold, underline, italic, etc.)
- Background colors
- Semantic color functions (Success, Error, Warning, Info)
- NO_COLOR environment variable support
- Custom output writers
- Zero dependencies

## Installation

```bash
go get github.com/comfortablynumb/che/pkg/checolor
```

## Usage

### Basic Colors

```go
package main

import (
    "fmt"
    "github.com/comfortablynumb/che/pkg/checolor"
)

func main() {
    // Simple colored strings
    fmt.Println(checolor.RedString("Error message"))
    fmt.Println(checolor.GreenString("Success message"))
    fmt.Println(checolor.YellowString("Warning message"))
    fmt.Println(checolor.BlueString("Info message"))
}
```

### Semantic Functions

```go
// Semantic color functions for common use cases
fmt.Println(checolor.Success("Operation completed"))
fmt.Println(checolor.Error("Operation failed"))
fmt.Println(checolor.Warning("Be careful"))
fmt.Println(checolor.Info("For your information"))
```

### Custom Colors

```go
// Create custom color combinations
c := checolor.New(checolor.Red, checolor.Bold, checolor.Underline)
fmt.Println(c.Sprint("Important message"))

// Add more styles
c.Add(checolor.BgWhite)
fmt.Println(c.Sprint("Even more important"))
```

### Text Formatting

```go
// Bold text
fmt.Println(checolor.BoldString("Bold text"))

// Underlined text
fmt.Println(checolor.UnderlineString("Underlined text"))

// Combined formatting
combined := checolor.Colorize("Bold Red", checolor.Red, checolor.Bold)
fmt.Println(combined)
```

### Print Functions

```go
// Print directly to stdout
c := checolor.New(checolor.Green)
c.Print("Success: ")
c.Println("Operation completed")

// Printf with formatting
c.Printf("Processed %d items\n", 42)
```

### Custom Output

```go
import "os"

// Write to stderr
c := checolor.New(checolor.Red)
c.SetOutput(os.Stderr)
c.Println("Error message to stderr")

// Write to a buffer
var buf bytes.Buffer
c.SetOutput(&buf)
c.Println("Colored output in buffer")
```

### Disable Colors

```go
// Disable colors globally
checolor.NoColor = true
fmt.Println(checolor.RedString("This won't be colored"))

// Or use NO_COLOR environment variable
// $ NO_COLOR=1 go run main.go
```

## Available Colors

### Foreground Colors
- `Black`, `Red`, `Green`, `Yellow`, `Blue`, `Magenta`, `Cyan`, `White`
- `BoldBlack`, `BoldRed`, `BoldGreen`, `BoldYellow`, `BoldBlue`, `BoldMagenta`, `BoldCyan`, `BoldWhite`

### Background Colors
- `BgBlack`, `BgRed`, `BgGreen`, `BgYellow`, `BgBlue`, `BgMagenta`, `BgCyan`, `BgWhite`

### Text Styles
- `Bold` - Bold text
- `Dim` - Dimmed text
- `Italic` - Italic text (not widely supported)
- `Underline` - Underlined text
- `Blink` - Blinking text (not widely supported)
- `Reverse` - Reversed foreground/background
- `Hidden` - Hidden text

## Examples

### CLI Application

```go
package main

import (
    "fmt"
    "github.com/comfortablynumb/che/pkg/checolor"
)

func main() {
    // Header
    header := checolor.New(checolor.Bold, checolor.Cyan)
    header.Println("=== My Application ===")

    // Success output
    fmt.Println(checolor.Success("✓ Database connected"))
    fmt.Println(checolor.Success("✓ Server started on :8080"))

    // Warning
    fmt.Println(checolor.Warning("⚠ Using development mode"))

    // Error
    fmt.Println(checolor.Error("✗ Cache connection failed"))

    // Info
    fmt.Println(checolor.Info("ℹ Press Ctrl+C to exit"))
}
```

### Log Levels

```go
type Logger struct {
    info  *checolor.Color
    warn  *checolor.Color
    err   *checolor.Color
    debug *checolor.Color
}

func NewLogger() *Logger {
    return &Logger{
        info:  checolor.New(checolor.Blue),
        warn:  checolor.New(checolor.Yellow),
        err:   checolor.New(checolor.Red, checolor.Bold),
        debug: checolor.New(checolor.Magenta),
    }
}

func (l *Logger) Info(msg string) {
    l.info.Printf("[INFO] %s\n", msg)
}

func (l *Logger) Warn(msg string) {
    l.warn.Printf("[WARN] %s\n", msg)
}

func (l *Logger) Error(msg string) {
    l.err.Printf("[ERROR] %s\n", msg)
}

func (l *Logger) Debug(msg string) {
    l.debug.Printf("[DEBUG] %s\n", msg)
}

func main() {
    log := NewLogger()
    log.Info("Application started")
    log.Warn("Low memory")
    log.Error("Connection failed")
    log.Debug("Request payload: {...}")
}
```

### Progress Indicator

```go
func showProgress(current, total int) {
    percent := float64(current) / float64(total) * 100

    var color *checolor.Color
    switch {
    case percent < 33:
        color = checolor.New(checolor.Red)
    case percent < 66:
        color = checolor.New(checolor.Yellow)
    default:
        color = checolor.New(checolor.Green)
    }

    color.Printf("Progress: %d/%d (%.1f%%)\n", current, total, percent)
}

func main() {
    total := 100
    for i := 0; i <= total; i += 25 {
        showProgress(i, total)
    }
}
```

### Table Output

```go
func printTable() {
    header := checolor.New(checolor.Bold, checolor.White, checolor.BgBlue)
    success := checolor.New(checolor.Green)
    failure := checolor.New(checolor.Red)

    // Header
    header.Printf("%-20s %-10s %-10s\n", "Service", "Status", "Uptime")

    // Rows
    fmt.Printf("%-20s ", "API Server")
    success.Printf("%-10s ", "Running")
    fmt.Printf("%-10s\n", "99.9%")

    fmt.Printf("%-20s ", "Database")
    success.Printf("%-10s ", "Running")
    fmt.Printf("%-10s\n", "99.5%")

    fmt.Printf("%-20s ", "Cache")
    failure.Printf("%-10s ", "Down")
    fmt.Printf("%-10s\n", "0%")
}
```

## API Reference

### Functions

- `Colorize(s string, codes ...string) string` - Wrap string with color codes
- `Sprint(code string, a ...interface{}) string` - Sprint with single color code
- `Sprintf(code, format string, a ...interface{}) string` - Sprintf with color code
- `Sprintln(code string, a ...interface{}) string` - Sprintln with color code

### Color Methods

- `Sprint(a ...interface{}) string` - Return colored string
- `Sprintf(format string, a ...interface{}) string` - Return colored formatted string
- `Sprintln(a ...interface{}) string` - Return colored string with newline
- `Print(a ...interface{}) (int, error)` - Print colored text
- `Printf(format string, a ...interface{}) (int, error)` - Print colored formatted text
- `Println(a ...interface{}) (int, error)` - Print colored text with newline
- `Add(codes ...string) *Color` - Add more style codes
- `SetOutput(w io.Writer) *Color` - Set output writer

### Convenience Functions

- `RedString(s string) string`
- `GreenString(s string) string`
- `YellowString(s string) string`
- `BlueString(s string) string`
- `MagentaString(s string) string`
- `CyanString(s string) string`
- `WhiteString(s string) string`
- `BoldString(s string) string`
- `UnderlineString(s string) string`

### Semantic Functions

- `Success(s string) string` - Green text for success messages
- `Error(s string) string` - Red text for error messages
- `Warning(s string) string` - Yellow text for warnings
- `Info(s string) string` - Blue text for informational messages

### Global Variables

- `NoColor bool` - Disable all colors globally (also respects NO_COLOR env var)
- `Output io.Writer` - Default output writer (defaults to os.Stdout)

## NO_COLOR Support

This package respects the [NO_COLOR](https://no-color.org/) standard. If the `NO_COLOR` environment variable is set (to any value), all color output will be disabled.

```bash
# Disable colors via environment variable
NO_COLOR=1 ./myapp

# Or in code
checolor.NoColor = true
```

## License

MIT
