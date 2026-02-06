// Package crypto provides cryptographic primitives for secure key derivation
// and authenticated encryption.
package crypto

import (
	"crypto/subtle"
)

// ZeroBytes overwrites the given byte slice with zeros to ensure sensitive
// data is removed from memory.
func ZeroBytes(data []byte) {
	for i := range data {
		data[i] = 0
	}
}

// SecureCompare performs a constant-time comparison of two byte slices
// to prevent timing attacks.
func SecureCompare(a, b []byte) bool {
	return subtle.ConstantTimeCompare(a, b) == 1
}
