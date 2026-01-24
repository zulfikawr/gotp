package totp

import (
	"time"
)

// TOTPParams holds the configuration parameters for generating or
// validating a Time-based One-Time Password (TOTP).
type TOTPParams struct {
	// Secret is the Base32-decoded secret key.
	Secret    []byte
	// Timestamp is the time for which the OTP is generated (usually time.Now()).
	Timestamp time.Time
	// Period is the time step in seconds (usually 30 or 60).
	Period    int
	// Digits is the number of digits in the generated code (6, 7, or 8).
	Digits    int
	// Algorithm is the hash function to use (SHA1, SHA256, or SHA512).
	Algorithm HashAlgorithm
}

// GenerateTOTP generates a Time-based One-Time Password (TOTP) as defined in RFC 6238.
// It calculates the time step counter based on the provided timestamp and period,
// then calls GenerateHOTP to produce the code.
func GenerateTOTP(params TOTPParams) (string, error) {
	// Apply default values if not specified.
	if params.Period <= 0 {
		params.Period = 30
	}
	if params.Digits == 0 {
		params.Digits = 6
	}
	if params.Algorithm == "" {
		params.Algorithm = SHA1
	}

	// T = (Current Unix Time - T0) / X
	// T0 is usually 0 (Unix epoch). X is the time step (Period).
	counter := uint64(params.Timestamp.Unix() / int64(params.Period))
	
	return GenerateHOTP(params.Secret, counter, params.Digits, params.Algorithm)
}

// ValidateTOTP validates a TOTP code with a given time window tolerance.
// A window of 1 allows for the current, previous, and next time steps to be valid,
// helping to account for clock drift between the client and server.
func ValidateTOTP(code string, params TOTPParams, window int) (bool, error) {
	if params.Period <= 0 {
		params.Period = 30
	}
	
	// Check the codes within the specified window of time steps.
	for i := -window; i <= window; i++ {
		t := params.Timestamp.Add(time.Duration(i*params.Period) * time.Second)
		p := params
		p.Timestamp = t
		
		generated, err := GenerateTOTP(p)
		if err != nil {
			return false, err
		}
		
		// Note: In a production environment with high security requirements,
		// a constant-time comparison should be used here to prevent timing attacks.
		if generated == code {
			return true, nil
		}
	}
	
	return false, nil
}

// RemainingSeconds returns the number of seconds remaining until the
// TOTP code generated for the given timestamp expires.
func RemainingSeconds(timestamp time.Time, period int) int {
	if period <= 0 {
		period = 30
	}
	return period - int(timestamp.Unix()%int64(period))
}

// NextExpiration returns the absolute time when the current TOTP period will end.
func NextExpiration(timestamp time.Time, period int) time.Time {
	if period <= 0 {
		period = 30
	}
	remaining := RemainingSeconds(timestamp, period)
	return timestamp.Truncate(time.Second).Add(time.Duration(remaining) * time.Second)
}
