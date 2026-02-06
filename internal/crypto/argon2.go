package crypto

import (
	"crypto/rand"
	"io"

	"golang.org/x/crypto/argon2"
)

// Argon2Params defines the parameters for the Argon2id key derivation function.
type Argon2Params struct {
	Memory      uint32 `json:"memory"`
	Iterations  uint32 `json:"iterations"`
	Parallelism uint8  `json:"parallelism"`
	SaltLength  uint32 `json:"salt_length"`
	KeyLength   uint32 `json:"key_length"`
}

// DefaultArgon2Params returns the recommended parameters for Argon2id.
func DefaultArgon2Params() Argon2Params {
	return Argon2Params{
		Memory:      65536, // 64MB
		Iterations:  3,
		Parallelism: 4,
		SaltLength:  16,
		KeyLength:   32,
	}
}

// DeriveKey derives a cryptographic key from a password and salt using Argon2id.
func DeriveKey(password []byte, salt []byte, params Argon2Params) []byte {
	return argon2.IDKey(
		password,
		salt,
		params.Iterations,
		params.Memory,
		params.Parallelism,
		params.KeyLength,
	)
}

// GenerateSalt generates a random salt of the specified length.
func GenerateSalt(length uint32) ([]byte, error) {
	salt := make([]byte, length)
	if _, err := io.ReadFull(rand.Reader, salt); err != nil {
		return nil, err
	}
	return salt, nil
}
