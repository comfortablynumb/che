package chevalid

import (
	"net"
	"net/url"
	"regexp"
	"strings"
)

// Email validation regex (simplified but covers most cases)
var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)

// UUID validation regex
var uuidRegex = regexp.MustCompile(`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`)

// IsEmail checks if a string is a valid email address.
func IsEmail(email string) bool {
	if len(email) < 3 || len(email) > 254 {
		return false
	}
	return emailRegex.MatchString(email)
}

// IsURL checks if a string is a valid URL.
func IsURL(str string) bool {
	u, err := url.Parse(str)
	if err != nil {
		return false
	}

	return u.Scheme != "" && u.Host != ""
}

// IsIP checks if a string is a valid IP address (IPv4 or IPv6).
func IsIP(str string) bool {
	return net.ParseIP(str) != nil
}

// IsIPv4 checks if a string is a valid IPv4 address.
func IsIPv4(str string) bool {
	ip := net.ParseIP(str)
	if ip == nil {
		return false
	}
	return ip.To4() != nil
}

// IsIPv6 checks if a string is a valid IPv6 address.
func IsIPv6(str string) bool {
	ip := net.ParseIP(str)
	if ip == nil {
		return false
	}
	return ip.To4() == nil
}

// IsUUID checks if a string is a valid UUID.
func IsUUID(str string) bool {
	return uuidRegex.MatchString(str)
}

// IsAlpha checks if a string contains only alphabetic characters.
func IsAlpha(str string) bool {
	for _, r := range str {
		if (r < 'a' || r > 'z') && (r < 'A' || r > 'Z') {
			return false
		}
	}
	return len(str) > 0
}

// IsAlphanumeric checks if a string contains only alphanumeric characters.
func IsAlphanumeric(str string) bool {
	for _, r := range str {
		if (r < 'a' || r > 'z') && (r < 'A' || r > 'Z') && (r < '0' || r > '9') {
			return false
		}
	}
	return len(str) > 0
}

// IsNumeric checks if a string contains only numeric characters.
func IsNumeric(str string) bool {
	for _, r := range str {
		if r < '0' || r > '9' {
			return false
		}
	}
	return len(str) > 0
}

// IsLuhn validates a string using the Luhn algorithm (used for credit cards).
func IsLuhn(str string) bool {
	str = strings.ReplaceAll(str, " ", "")
	str = strings.ReplaceAll(str, "-", "")

	if !IsNumeric(str) {
		return false
	}

	sum := 0
	double := false

	for i := len(str) - 1; i >= 0; i-- {
		digit := int(str[i] - '0')

		if double {
			digit *= 2
			if digit > 9 {
				digit -= 9
			}
		}

		sum += digit
		double = !double
	}

	return sum%10 == 0
}

// MinLength checks if a string has at least the specified length.
func MinLength(str string, min int) bool {
	return len(str) >= min
}

// MaxLength checks if a string has at most the specified length.
func MaxLength(str string, max int) bool {
	return len(str) <= max
}

// LengthBetween checks if a string length is between min and max (inclusive).
func LengthBetween(str string, min, max int) bool {
	length := len(str)
	return length >= min && length <= max
}

// Contains checks if a string contains a substring.
func Contains(str, substr string) bool {
	return strings.Contains(str, substr)
}

// HasPrefix checks if a string starts with a prefix.
func HasPrefix(str, prefix string) bool {
	return strings.HasPrefix(str, prefix)
}

// HasSuffix checks if a string ends with a suffix.
func HasSuffix(str, suffix string) bool {
	return strings.HasSuffix(str, suffix)
}
