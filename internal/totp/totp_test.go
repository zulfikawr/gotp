package totp

import (
	"testing"
	"time"
)

func TestTOTP_RFC6238(t *testing.T) {
	seed1 := []byte("12345678901234567890")
	seed256 := []byte("12345678901234567890123456789012")
	seed512 := []byte("1234567890123456789012345678901234567890123456789012345678901234")

	testCases := []struct {
		time     int64
		sha1     string
		sha256   string
		sha512   string
	}{
		{59, "94287082", "46119246", "90693936"},
		{1111111109, "07081804", "68084774", "25091201"},
		{1111111111, "14050471", "67062674", "99943326"},
		{1234567890, "89005924", "91819424", "93441116"},
		{2000000000, "69279037", "90698825", "38618901"},
	}

	for _, tc := range testCases {
		ts := time.Unix(tc.time, 0)

		// Test SHA1
		sha1, _ := GenerateTOTP(TOTPParams{
			Secret:    seed1,
			Timestamp: ts,
			Period:    30,
			Digits:    8,
			Algorithm: SHA1,
		})
		if sha1 != tc.sha1 {
			t.Errorf("At time %d, SHA1 expected %s, got %s", tc.time, tc.sha1, sha1)
		}

		// Test SHA256
		sha256, _ := GenerateTOTP(TOTPParams{
			Secret:    seed256,
			Timestamp: ts,
			Period:    30,
			Digits:    8,
			Algorithm: SHA256,
		})
		if sha256 != tc.sha256 {
			t.Errorf("At time %d, SHA256 expected %s, got %s", tc.time, tc.sha256, sha256)
		}

		// Test SHA512
		sha512, _ := GenerateTOTP(TOTPParams{
			Secret:    seed512,
			Timestamp: ts,
			Period:    30,
			Digits:    8,
			Algorithm: SHA512,
		})
		if sha512 != tc.sha512 {
			t.Errorf("At time %d, SHA512 expected %s, got %s", tc.time, tc.sha512, sha512)
		}
	}
}

func TestTOTP_Validation(t *testing.T) {
	seed := []byte("JBSWY3DPEHPK3PXP")
	ts := time.Unix(1234567890, 0)
	params := TOTPParams{
		Secret:    seed,
		Timestamp: ts,
		Period:    30,
		Digits:    6,
		Algorithm: SHA1,
	}

	code, _ := GenerateTOTP(params)

	// Exact match
	valid, _ := ValidateTOTP(code, params, 1)
	if !valid {
		t.Error("Expected code to be valid (exact match)")
	}

	// Within window (one period back)
	paramsPast := params
	paramsPast.Timestamp = ts.Add(30 * time.Second)
	valid, _ = ValidateTOTP(code, paramsPast, 1)
	if !valid {
		t.Error("Expected code to be valid (within window -1)")
	}

	// Outside window
	valid, _ = ValidateTOTP(code, paramsPast, 0)
	if valid {
		t.Error("Expected code to be invalid (outside window)")
	}
}

func TestRemainingSeconds(t *testing.T) {
	ts := time.Unix(30, 0) // Exactly at period start
	rem := RemainingSeconds(ts, 30)
	if rem != 30 {
		t.Errorf("Expected 30s remaining at period start, got %d", rem)
	}

	ts = time.Unix(45, 0) // Middle of period
	rem = RemainingSeconds(ts, 30)
	if rem != 15 {
		t.Errorf("Expected 15s remaining at middle of period, got %d", rem)
	}

	ts = time.Unix(59, 0) // End of period
	rem = RemainingSeconds(ts, 30)
	if rem != 1 {
		t.Errorf("Expected 1s remaining at end of period, got %d", rem)
	}
}
