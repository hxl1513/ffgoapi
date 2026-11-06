package common

import "strings"

// IsEmpty returns true if the string is empty or contains only whitespace.
func IsEmpty(s string) bool {
	return strings.TrimSpace(s) == ""
}

// IsNotEmpty returns true if the string is non-empty and contains non-whitespace characters.
func IsNotEmpty(s string) bool {
	return !IsEmpty(s)
}
