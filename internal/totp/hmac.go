// Package totp provides a core engine for generating and validating
// Time-based One-Time Passwords (TOTP) as specified in RFC 6238 and
// HMAC-based One-Time Passwords (HOTP) as specified in RFC 4226.
package totp

import (
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"hash"
)

// HashAlgorithm represents the supported cryptographic hash algorithms
// for HMAC-based OTP generation.
type HashAlgorithm string

const (
	// SHA1 is the default hash algorithm used by most TOTP implementations (RFC 4226).
	SHA1   HashAlgorithm = "SHA1"
	// SHA256 provides enhanced security over SHA1 (RFC 6238).
	SHA256 HashAlgorithm = "SHA256"
	// SHA512 provides maximum security for TOTP generation (RFC 6238).
	SHA512 HashAlgorithm = "SHA512"
)

// HMAC implements the HMAC keyed-hash message authentication code algorithm
// as defined in RFC 2104. This implementation is dependency-free and
// specifically designed for OTP generation.
//
// HMAC(K, m) = H((K ^ opad) || H((K ^ ipad) || m))
func HMAC(key []byte, message []byte, algo HashAlgorithm) []byte {
	var h func() hash.Hash
	var blockSize int

	switch algo {
	case SHA256:
		h = sha256.New
		blockSize = 64
	case SHA512:
		h = sha512.New
		blockSize = 128
	default: // Default to SHA1
		h = sha1.New
		blockSize = 64
	}

	// If key is longer than blockSize, hash it first according to RFC 2104.
	if len(key) > blockSize {
		hasher := h()
		hasher.Write(key)
		key = hasher.Sum(nil)
	}

	// If key is shorter than blockSize, pad it with zeros.
	if len(key) < blockSize {
		paddedKey := make([]byte, blockSize)
		copy(paddedKey, key)
		key = paddedKey
	}

	// ipad and opad are the inner and outer padding constants.
	ipad := make([]byte, blockSize)
	opad := make([]byte, blockSize)
	for i := range key {
		ipad[i] = key[i] ^ 0x36
		opad[i] = key[i] ^ 0x5c
	}

	// Inner hash: H((K ^ ipad) || m)
	innerHasher := h()
	innerHasher.Write(ipad)
	innerHasher.Write(message)
	innerHash := innerHasher.Sum(nil)

	// Outer hash: H((K ^ opad) || innerHash)
	outerHasher := h()
	outerHasher.Write(opad)
	outerHasher.Write(innerHash)
	return outerHasher.Sum(nil)
}