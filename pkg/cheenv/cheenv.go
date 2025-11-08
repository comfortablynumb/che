package cheenv

import (
	"os"
	"strconv"
	"strings"
	"time"
)

// Get retrieves an environment variable as a string
// Returns the default value if the variable is not set
func Get(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

// MustGet retrieves an environment variable as a string
// Panics if the variable is not set
func MustGet(key string) string {
	value := os.Getenv(key)
	if value == "" {
		panic("cheenv: required environment variable not set: " + key)
	}
	return value
}

// GetInt retrieves an environment variable as an int
// Returns the default value if the variable is not set or cannot be parsed
func GetInt(key string, defaultValue int) int {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}

	intValue, err := strconv.Atoi(value)
	if err != nil {
		return defaultValue
	}

	return intValue
}

// MustGetInt retrieves an environment variable as an int
// Panics if the variable is not set or cannot be parsed
func MustGetInt(key string) int {
	value := os.Getenv(key)
	if value == "" {
		panic("cheenv: required environment variable not set: " + key)
	}

	intValue, err := strconv.Atoi(value)
	if err != nil {
		panic("cheenv: failed to parse int from environment variable " + key + ": " + err.Error())
	}

	return intValue
}

// GetInt64 retrieves an environment variable as an int64
// Returns the default value if the variable is not set or cannot be parsed
func GetInt64(key string, defaultValue int64) int64 {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}

	int64Value, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return defaultValue
	}

	return int64Value
}

// GetFloat retrieves an environment variable as a float64
// Returns the default value if the variable is not set or cannot be parsed
func GetFloat(key string, defaultValue float64) float64 {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}

	floatValue, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return defaultValue
	}

	return floatValue
}

// GetBool retrieves an environment variable as a bool
// Accepts: true, false, 1, 0, yes, no, on, off (case-insensitive)
// Returns the default value if the variable is not set or cannot be parsed
func GetBool(key string, defaultValue bool) bool {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}

	boolValue, err := parseBool(value)
	if err != nil {
		return defaultValue
	}

	return boolValue
}

// MustGetBool retrieves an environment variable as a bool
// Panics if the variable is not set or cannot be parsed
func MustGetBool(key string) bool {
	value := os.Getenv(key)
	if value == "" {
		panic("cheenv: required environment variable not set: " + key)
	}

	boolValue, err := parseBool(value)
	if err != nil {
		panic("cheenv: failed to parse bool from environment variable " + key + ": " + err.Error())
	}

	return boolValue
}

// GetDuration retrieves an environment variable as a time.Duration
// Accepts values like "1s", "5m", "2h", etc.
// Returns the default value if the variable is not set or cannot be parsed
func GetDuration(key string, defaultValue time.Duration) time.Duration {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}

	duration, err := time.ParseDuration(value)
	if err != nil {
		return defaultValue
	}

	return duration
}

// GetStringList retrieves an environment variable as a slice of strings
// Splits by the specified separator and trims whitespace
// Returns the default value if the variable is not set
func GetStringList(key, separator string, defaultValue []string) []string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}

	parts := strings.Split(value, separator)
	result := make([]string, 0, len(parts))
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}

	if len(result) == 0 {
		return defaultValue
	}

	return result
}

// GetIntList retrieves an environment variable as a slice of ints
// Splits by the specified separator and parses each value
// Returns the default value if the variable is not set or any value cannot be parsed
func GetIntList(key, separator string, defaultValue []int) []int {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}

	parts := strings.Split(value, separator)
	result := make([]int, 0, len(parts))
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed == "" {
			continue
		}

		intValue, err := strconv.Atoi(trimmed)
		if err != nil {
			return defaultValue
		}
		result = append(result, intValue)
	}

	if len(result) == 0 {
		return defaultValue
	}

	return result
}

// Set sets an environment variable
func Set(key, value string) error {
	return os.Setenv(key, value)
}

// Unset unsets an environment variable
func Unset(key string) error {
	return os.Unsetenv(key)
}

// Has checks if an environment variable is set (even if empty)
func Has(key string) bool {
	_, exists := os.LookupEnv(key)
	return exists
}

// GetAll returns all environment variables as a map
func GetAll() map[string]string {
	environ := os.Environ()
	result := make(map[string]string, len(environ))

	for _, env := range environ {
		parts := strings.SplitN(env, "=", 2)
		if len(parts) == 2 {
			result[parts[0]] = parts[1]
		}
	}

	return result
}

// GetWithPrefix returns all environment variables with a given prefix
// The prefix is removed from the keys in the result
func GetWithPrefix(prefix string) map[string]string {
	environ := os.Environ()
	result := make(map[string]string)

	for _, env := range environ {
		parts := strings.SplitN(env, "=", 2)
		if len(parts) == 2 && strings.HasPrefix(parts[0], prefix) {
			key := strings.TrimPrefix(parts[0], prefix)
			result[key] = parts[1]
		}
	}

	return result
}

// parseBool parses a bool from a string with more flexibility than strconv.ParseBool
func parseBool(value string) (bool, error) {
	value = strings.ToLower(strings.TrimSpace(value))

	switch value {
	case "true", "1", "yes", "on", "y", "t":
		return true, nil
	case "false", "0", "no", "off", "n", "f":
		return false, nil
	default:
		return false, strconv.ErrSyntax
	}
}
