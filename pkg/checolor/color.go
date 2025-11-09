// Package checolor provides terminal color and text formatting utilities.
package checolor

import (
	"fmt"
	"io"
	"os"
	"strings"
)

// Color codes
const (
	Reset = "\033[0m"

	// Regular colors
	Black   = "\033[30m"
	Red     = "\033[31m"
	Green   = "\033[32m"
	Yellow  = "\033[33m"
	Blue    = "\033[34m"
	Magenta = "\033[35m"
	Cyan    = "\033[36m"
	White   = "\033[37m"

	// Bold colors
	BoldBlack   = "\033[1;30m"
	BoldRed     = "\033[1;31m"
	BoldGreen   = "\033[1;32m"
	BoldYellow  = "\033[1;33m"
	BoldBlue    = "\033[1;34m"
	BoldMagenta = "\033[1;35m"
	BoldCyan    = "\033[1;36m"
	BoldWhite   = "\033[1;37m"

	// Background colors
	BgBlack   = "\033[40m"
	BgRed     = "\033[41m"
	BgGreen   = "\033[42m"
	BgYellow  = "\033[43m"
	BgBlue    = "\033[44m"
	BgMagenta = "\033[45m"
	BgCyan    = "\033[46m"
	BgWhite   = "\033[47m"

	// Text styles
	Bold      = "\033[1m"
	Dim       = "\033[2m"
	Italic    = "\033[3m"
	Underline = "\033[4m"
	Blink     = "\033[5m"
	Reverse   = "\033[7m"
	Hidden    = "\033[8m"
)

var (
	// NoColor disables all color output globally
	NoColor = false

	// Output is the default output writer (defaults to os.Stdout)
	Output io.Writer = os.Stdout
)

func init() {
	// Check NO_COLOR environment variable
	if os.Getenv("NO_COLOR") != "" {
		NoColor = true
	}
}

// Color represents a color/style configuration.
type Color struct {
	codes  []string
	writer io.Writer
}

// New creates a new Color with the given codes.
func New(codes ...string) *Color {
	return &Color{
		codes:  codes,
		writer: Output,
	}
}

// SetOutput sets the output writer for this color.
func (c *Color) SetOutput(w io.Writer) *Color {
	c.writer = w
	return c
}

// Sprint returns a colored string.
func (c *Color) Sprint(a ...interface{}) string {
	if NoColor {
		return fmt.Sprint(a...)
	}
	return c.wrap(fmt.Sprint(a...))
}

// Sprintf returns a colored formatted string.
func (c *Color) Sprintf(format string, a ...interface{}) string {
	if NoColor {
		return fmt.Sprintf(format, a...)
	}
	return c.wrap(fmt.Sprintf(format, a...))
}

// Sprintln returns a colored string with newline.
func (c *Color) Sprintln(a ...interface{}) string {
	if NoColor {
		return fmt.Sprintln(a...)
	}
	return c.wrap(fmt.Sprintln(a...))
}

// Print prints colored text to the configured output.
func (c *Color) Print(a ...interface{}) (n int, err error) {
	return fmt.Fprint(c.writer, c.Sprint(a...))
}

// Printf prints colored formatted text to the configured output.
func (c *Color) Printf(format string, a ...interface{}) (n int, err error) {
	return fmt.Fprint(c.writer, c.Sprintf(format, a...))
}

// Println prints colored text with newline to the configured output.
func (c *Color) Println(a ...interface{}) (n int, err error) {
	return fmt.Fprintln(c.writer, c.Sprint(a...))
}

// Add adds more codes to this color.
func (c *Color) Add(codes ...string) *Color {
	c.codes = append(c.codes, codes...)
	return c
}

func (c *Color) wrap(s string) string {
	if len(c.codes) == 0 {
		return s
	}
	return strings.Join(c.codes, "") + s + Reset
}

// Colorize wraps the string with the given codes.
func Colorize(s string, codes ...string) string {
	if NoColor || len(codes) == 0 {
		return s
	}
	return strings.Join(codes, "") + s + Reset
}

// Predefined color functions
func Sprint(code string, a ...interface{}) string {
	return New(code).Sprint(a...)
}

func Sprintf(code, format string, a ...interface{}) string {
	return New(code).Sprintf(format, a...)
}

func Sprintln(code string, a ...interface{}) string {
	return New(code).Sprintln(a...)
}

// Convenience color functions
func RedString(s string) string       { return Colorize(s, Red) }
func GreenString(s string) string     { return Colorize(s, Green) }
func YellowString(s string) string    { return Colorize(s, Yellow) }
func BlueString(s string) string      { return Colorize(s, Blue) }
func MagentaString(s string) string   { return Colorize(s, Magenta) }
func CyanString(s string) string      { return Colorize(s, Cyan) }
func WhiteString(s string) string     { return Colorize(s, White) }
func BoldString(s string) string      { return Colorize(s, Bold) }
func UnderlineString(s string) string { return Colorize(s, Underline) }

// Success prints green text
func Success(s string) string { return GreenString(s) }

// Error prints red text
func Error(s string) string { return RedString(s) }

// Warning prints yellow text
func Warning(s string) string { return YellowString(s) }

// Info prints blue text
func Info(s string) string { return BlueString(s) }
