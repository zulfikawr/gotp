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
		time   int64
		sha1   string
		sha256 string
		sha512 string
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

func TestNextExpiration(t *testing.T) {
	ts := time.Unix(30, 0)
	next := NextExpiration(ts, 30)
	if next.Unix() != 60 {
		t.Errorf("Expected next expiration at 60, got %d", next.Unix())
	}
}

func TestHOTP_ErrorCases(t *testing.T) {
	_, err := GenerateHOTP([]byte("secret"), 1, 5, SHA1)
	if err == nil {
		t.Error("Expected error for digits < 6")
	}
	_, err = GenerateHOTP([]byte("secret"), 1, 9, SHA1)
	if err == nil {
		t.Error("Expected error for digits > 8")
	}
}

func TestHMAC_KeyPadding(t *testing.T) {
	// Test key longer than block size (64 for SHA1)
	longKey := make([]byte, 100)
	for i := range longKey {
		longKey[i] = byte(i)
	}
	res1 := HMAC(longKey, []byte("msg"), SHA1)
	if len(res1) != 20 {
		t.Error("HMAC-SHA1 should return 20 bytes")
	}

	// Test key shorter than block size
	shortKey := []byte("short")
	res2 := HMAC(shortKey, []byte("msg"), SHA1)
	if len(res2) != 20 {
		t.Error("HMAC-SHA1 should return 20 bytes")
	}

	// Test SHA512 (block size 128)
	res3 := HMAC(longKey, []byte("msg"), SHA512)
	if len(res3) != 64 {
		t.Error("HMAC-SHA512 should return 64 bytes")
	}
}

func TestTOTP_DefaultValues(t *testing.T) {
	// Test that default values are applied
	code, err := GenerateTOTP(TOTPParams{
		Secret:    []byte("JBSWY3DPEHPK3PXP"),
		Timestamp: time.Unix(1234567890, 0),
	})
	if err != nil {
		t.Fatalf("GenerateTOTP failed: %v", err)
	}
	if len(code) != 6 {
		t.Errorf("Expected 6 digits, got %d", len(code))
	}
}

func TestTOTP_ValidationWindow(t *testing.T) {
	params := TOTPParams{
		Secret:    []byte("JBSWY3DPEHPK3PXP"),
		Timestamp: time.Unix(1234567890, 0),
	}
	code, _ := GenerateTOTP(params)

	// Test with window = 0
	valid, _ := ValidateTOTP(code, params, 0)
	if !valid {
		t.Error("Expected code to be valid with window 0")
	}
}
