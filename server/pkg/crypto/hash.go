package crypto

import (
	"crypto/sha256"
	"encoding/hex"
)

// HashSHA256 hashes a string using SHA-256 and returns the hex-encoded digest.
func HashSHA256(s string) string {
	h := sha256.Sum256([]byte(s))
	return hex.EncodeToString(h[:])
}
