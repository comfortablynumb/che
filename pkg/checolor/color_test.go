package checolor

import (
	"bytes"
	"strings"
	"testing"

	"github.com/comfortablynumb/che/pkg/chetest"
)

func TestColorize(t *testing.T) {
	// Save and restore NoColor
	oldNoColor := NoColor
	defer func() { NoColor = oldNoColor }()

	NoColor = false
	result := Colorize("hello", Red)
	expected := "\033[31mhello\033[0m"
	chetest.RequireEqual(t, result, expected)
}

func TestColorize_NoColor(t *testing.T) {
	oldNoColor := NoColor
	defer func() { NoColor = oldNoColor }()

	NoColor = true
	result := Colorize("hello", Red)
	chetest.RequireEqual(t, result, "hello")
}

func TestColorize_NoCodes(t *testing.T) {
	oldNoColor := NoColor
	defer func() { NoColor = oldNoColor }()

	NoColor = false
	result := Colorize("hello")
	chetest.RequireEqual(t, result, "hello")
}

func TestColor_Sprint(t *testing.T) {
	oldNoColor := NoColor
	defer func() { NoColor = oldNoColor }()

	NoColor = false
	c := New(Red, Bold)
	result := c.Sprint("hello", " ", "world")
	expected := "\033[31m\033[1mhello world\033[0m"
	chetest.RequireEqual(t, result, expected)
}

func TestColor_Sprintf(t *testing.T) {
	oldNoColor := NoColor
	defer func() { NoColor = oldNoColor }()

	NoColor = false
	c := New(Green)
	result := c.Sprintf("hello %s", "world")
	expected := "\033[32mhello world\033[0m"
	chetest.RequireEqual(t, result, expected)
}

func TestColor_Sprintln(t *testing.T) {
	oldNoColor := NoColor
	defer func() { NoColor = oldNoColor }()

	NoColor = false
	c := New(Blue)
	result := c.Sprintln("hello")
	expected := "\033[34mhello\n\033[0m"
	chetest.RequireEqual(t, result, expected)
}

func TestColor_Print(t *testing.T) {
	oldNoColor := NoColor
	defer func() { NoColor = oldNoColor }()

	NoColor = false
	buf := &bytes.Buffer{}
	c := New(Red).SetOutput(buf)

	n, err := c.Print("hello")
	chetest.RequireEqual(t, err, nil)
	chetest.RequireEqual(t, n > 0, true)
	chetest.RequireEqual(t, buf.String(), "\033[31mhello\033[0m")
}

func TestColor_Printf(t *testing.T) {
	oldNoColor := NoColor
	defer func() { NoColor = oldNoColor }()

	NoColor = false
	buf := &bytes.Buffer{}
	c := New(Green).SetOutput(buf)

	n, err := c.Printf("hello %s", "world")
	chetest.RequireEqual(t, err, nil)
	chetest.RequireEqual(t, n > 0, true)
	chetest.RequireEqual(t, buf.String(), "\033[32mhello world\033[0m")
}

func TestColor_Println(t *testing.T) {
	oldNoColor := NoColor
	defer func() { NoColor = oldNoColor }()

	NoColor = false
	buf := &bytes.Buffer{}
	c := New(Blue).SetOutput(buf)

	n, err := c.Println("hello")
	chetest.RequireEqual(t, err, nil)
	chetest.RequireEqual(t, n > 0, true)

	output := buf.String()
	chetest.RequireEqual(t, strings.HasPrefix(output, "\033[34mhello"), true)
	chetest.RequireEqual(t, strings.HasSuffix(output, "\033[0m\n"), true)
}

func TestColor_Add(t *testing.T) {
	oldNoColor := NoColor
	defer func() { NoColor = oldNoColor }()

	NoColor = false
	c := New(Red)
	c.Add(Bold, Underline)

	result := c.Sprint("hello")
	expected := "\033[31m\033[1m\033[4mhello\033[0m"
	chetest.RequireEqual(t, result, expected)
}

func TestColor_NoColor(t *testing.T) {
	oldNoColor := NoColor
	defer func() { NoColor = oldNoColor }()

	NoColor = true
	c := New(Red, Bold)
	result := c.Sprint("hello")
	chetest.RequireEqual(t, result, "hello")
}

func TestSprintFunctions(t *testing.T) {
	oldNoColor := NoColor
	defer func() { NoColor = oldNoColor }()

	NoColor = false

	result := Sprint(Red, "hello")
	chetest.RequireEqual(t, result, "\033[31mhello\033[0m")

	result = Sprintf(Green, "hello %s", "world")
	chetest.RequireEqual(t, result, "\033[32mhello world\033[0m")

	result = Sprintln(Blue, "hello")
	chetest.RequireEqual(t, strings.HasPrefix(result, "\033[34mhello"), true)
}

func TestConvenienceFunctions(t *testing.T) {
	oldNoColor := NoColor
	defer func() { NoColor = oldNoColor }()

	NoColor = false

	tests := []struct {
		fn       func(string) string
		input    string
		expected string
	}{
		{RedString, "hello", "\033[31mhello\033[0m"},
		{GreenString, "hello", "\033[32mhello\033[0m"},
		{YellowString, "hello", "\033[33mhello\033[0m"},
		{BlueString, "hello", "\033[34mhello\033[0m"},
		{MagentaString, "hello", "\033[35mhello\033[0m"},
		{CyanString, "hello", "\033[36mhello\033[0m"},
		{WhiteString, "hello", "\033[37mhello\033[0m"},
		{BoldString, "hello", "\033[1mhello\033[0m"},
		{UnderlineString, "hello", "\033[4mhello\033[0m"},
	}

	for _, tt := range tests {
		result := tt.fn(tt.input)
		chetest.RequireEqual(t, result, tt.expected)
	}
}

func TestSemanticFunctions(t *testing.T) {
	oldNoColor := NoColor
	defer func() { NoColor = oldNoColor }()

	NoColor = false

	tests := []struct {
		fn       func(string) string
		input    string
		expected string
	}{
		{Success, "done", "\033[32mdone\033[0m"},
		{Error, "failed", "\033[31mfailed\033[0m"},
		{Warning, "careful", "\033[33mcareful\033[0m"},
		{Info, "note", "\033[34mnote\033[0m"},
	}

	for _, tt := range tests {
		result := tt.fn(tt.input)
		chetest.RequireEqual(t, result, tt.expected)
	}
}

func TestMultipleCodes(t *testing.T) {
	oldNoColor := NoColor
	defer func() { NoColor = oldNoColor }()

	NoColor = false
	result := Colorize("hello", Red, Bold, Underline)
	expected := "\033[31m\033[1m\033[4mhello\033[0m"
	chetest.RequireEqual(t, result, expected)
}
