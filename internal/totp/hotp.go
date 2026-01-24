package totp

import (
	"encoding/binary"
	"fmt"
	"math"
)

// GenerateHOTP generates a HMAC-based One-Time Password (HOTP) as defined in RFC 4226.
// It takes a secret key, a counter value, the number of digits for the OTP, and
// the hash algorithm to use for HMAC.
//
// The function follows the dynamic truncation algorithm to extract a 31-bit
// integer from the HMAC result and then applies modulo 10^digits to get the OTP.
func GenerateHOTP(secret []byte, counter uint64, digits int, algo HashAlgorithm) (string, error) {
	if digits < 6 || digits > 8 {
		return "", fmt.Errorf("digits must be between 6 and 8")
	}

	// 1. Generate HMAC-SHA-1 (or SHA-256/512) result HS.
	// The counter is converted to an 8-byte big-endian integer.
	counterBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(counterBytes, counter)
	hs := HMAC(secret, counterBytes, algo)

	// 2. Dynamic Truncation (DT)
	// Extract the low-order 4 bits of the last byte of HS to use as an offset.
	offset := hs[len(hs)-1] & 0x0f
	
	// Extract a 4-byte sequence starting at HS[offset].
	p := hs[offset : offset+4]
	
	// Convert p to a 31-bit integer by ignoring the most significant bit (0x7fffffff).
	binaryValue := uint32(binary.BigEndian.Uint32(p) & 0x7fffffff)

	// 3. Compute the OTP value.
	// Apply modulo 10^digits to the 31-bit integer.
	otp := binaryValue % uint32(math.Pow10(digits))

	// Format the resulting OTP as a string with leading zeros if necessary.
	format := fmt.Sprintf("%%0%dd", digits)
	return fmt.Sprintf(format, otp), nil
}