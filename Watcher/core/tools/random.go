package tools

import (
	"crypto/rand"
	"fmt"
)

// RandomToken creates a random token of length bytes, the function returns
// the randomly generated token in string format.
func RandomToken(length int) string {
	b := make([]byte, length)
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}
