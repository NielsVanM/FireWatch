package tools

import (
	"strings"
)

// PasswordVerification applies the password verification rules
// it allows the software to check for valid passwords
func PasswordVerification(password string) bool {

	// Check length
	if len(password) < 8 {
		return false
	}

	// Check numeric value
	if !strings.ContainsAny(password, "0123456789") {
		return false
	}

	// Reached the end, return true
	return true
}
