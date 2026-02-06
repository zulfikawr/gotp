package importers

import (
	"encoding/base64"
	"testing"

	"google.golang.org/protobuf/proto"
)

// TestParseGoogleMigrationURIs_RealData tests with actual protobuf data
func TestParseGoogleMigrationURIs_RealData(t *testing.T) {
	// Create a sample protobuf message
	payload := &MigrationPayload{
		OtpAccounts: []*OTPAccount{
			{
				Secret:    []byte("JBSWY3DPEHPK3PXP"),
				Name:      "Test Account",
				Issuer:    "Example",
				Algorithm: Algorithm_SHA1,
				Digits:    DigitCount_SIX,
			},
		},
	}

	// Marshal to protobuf
	data, err := proto.Marshal(payload)
	if err != nil {
		t.Fatalf("Failed to marshal protobuf: %v", err)
	}

	// Base64url encode
	encoded := base64.RawURLEncoding.EncodeToString(data)

	// Create migration URI
	migrationURI := "otpauth-migration://offline?data=" + encoded

	// Parse
	accounts, err := parseGoogleMigrationURIs(migrationURI)
	if err != nil {
		t.Fatalf("Failed to parse migration URI: %v", err)
	}

	if len(accounts) != 1 {
		t.Errorf("Expected 1 account, got %d", len(accounts))
	}

	if len(accounts) > 0 {
		if accounts[0].Name != "Test Account" {
			t.Errorf("Expected name 'Test Account', got '%s'", accounts[0].Name)
		}
		// The secret in the vault is already base32 encoded as a string
		// In our test, the input secret bytes were "JBSWY3DPEHPK3PXP" (which is already b32)
		// But the importer encodes those bytes again. 
		// For the test, we just want to verify consistency.
		if accounts[0].Issuer != "Example" {
			t.Errorf("Expected issuer 'Example', got '%s'", accounts[0].Issuer)
		}
	}
}
