package base32

import (
	"bytes"
	"testing"
)

func TestBase32(t *testing.T) {
	testCases := []struct {
		input    string
		expected string
	}{
		{"", ""},
		{"f", "MY======"},
		{"fo", "MZXQ===="},
		{"foo", "MZXW6==="},
		{"foob", "MZXW6YQ="},
		{"fooba", "MZXW6YTB"},
		{"foobar", "MZXW6YTBOI======"},
	}

	for _, tc := range testCases {
		encoded := Encode([]byte(tc.input))
		if encoded != tc.expected {
			t.Errorf("Encode(%q) = %q, expected %q", tc.input, encoded, tc.expected)
		}

		decoded, err := Decode(tc.expected)
		if err != nil {
			t.Errorf("Decode(%q) error: %v", tc.expected, err)
			continue
		}
		if !bytes.Equal(decoded, []byte(tc.input)) {
			t.Errorf("Decode(%q) = %q, expected %q", tc.expected, string(decoded), tc.input)
		}
	}
}

func TestBase32DecodeInvalid(t *testing.T) {
	_, err := Decode("123") // '1' is invalid
	if err == nil {
		t.Error("Expected error for invalid base32 string, got nil")
	}
}
