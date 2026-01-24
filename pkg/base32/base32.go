// Package base32 provides a dependency-free implementation of the Base32
// encoding and decoding as specified in RFC 4648.
package base32

import (
	"errors"
	"strings"
)

// alphabet is the standard Base32 character set (A-Z, 2-7).
const alphabet = "ABCDEFGHIJKLMNOPQRSTUVWXYZ234567"

// Encode encodes a byte slice into a Base32 string with padding.
// It follows the algorithm described in RFC 4648.
func Encode(data []byte) string {
	if len(data) == 0 {
		return ""
	}

	var result strings.Builder
	var buffer uint64
	var bitsLeft uint

	for _, b := range data {
		// Shift the buffer and add 8 new bits.
		buffer = (buffer << 8) | uint64(b)
		bitsLeft += 8
		
		// Extract 5-bit chunks from the buffer.
		for bitsLeft >= 5 {
			bitsLeft -= 5
			index := (buffer >> bitsLeft) & 0x1F
			result.WriteByte(alphabet[index])
		}
	}

	// Handle remaining bits.
	if bitsLeft > 0 {
		index := (buffer << (5 - bitsLeft)) & 0x1F
		result.WriteByte(alphabet[index])
	}

	// Add '=' padding to make the length a multiple of 8.
	for result.Len()%8 != 0 {
		result.WriteByte('=')
	}

	return result.String()
}

// Decode decodes a Base32 string into its original byte slice.
// It supports strings with or without padding and is case-insensitive.
// Returns an error if the input contains characters outside the Base32 alphabet.
func Decode(s string) ([]byte, error) {
	// Normalize the input: remove padding and convert to uppercase.
	s = strings.ToUpper(strings.TrimRight(s, "="))
	if s == "" {
		return []byte{}, nil
	}

	var result []byte
	var buffer uint64
	var bitsLeft uint

	for i := 0; i < len(s); i++ {
		char := s[i]
		var val int
		
		// Map the character to its 5-bit value.
		if char >= 'A' && char <= 'Z' {
			val = int(char - 'A')
		} else if char >= '2' && char <= '7' {
			val = int(char - '2' + 26)
		} else {
			return nil, errors.New("invalid base32 character")
		}

		// Shift the buffer and add 5 new bits.
		buffer = (buffer << 5) | uint64(val)
		bitsLeft += 5

		// Extract 8-bit bytes from the buffer.
		if bitsLeft >= 8 {
			bitsLeft -= 8
			result = append(result, byte(buffer>>bitsLeft))
		}
	}

	return result, nil
}
