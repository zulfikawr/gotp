package importers

import (
	"testing"
)

func TestParseAuthyExport(t *testing.T) {
	authyJSON := `{
		"accounts": [
			{
				"id": 12345,
				"name": "Test Account",
				"issuer": "TestIssuer",
				"username": "test@example.com",
				"secret": "JBSWY3DPEHPK3PXP",
				"type": "totp",
				"algorithm": "SHA1",
				"digits": 6,
				"period": 30
			}
		]
	}`

	accounts, err := ParseAuthyExport([]byte(authyJSON))
	if err != nil {
		t.Fatalf("ParseAuthyExport() error = %v", err)
	}

	if len(accounts) != 1 {
		t.Errorf("ParseAuthyExport() returned %d accounts, want 1", len(accounts))
	}

	if accounts[0].Name != "Test Account" {
		t.Errorf("ParseAuthyExport() account name = %s, want Test Account", accounts[0].Name)
	}
}

func TestParseAuthyExport_WithBase32Secret(t *testing.T) {
	authyJSON := `{
		"accounts": [
			{
				"id": 12345,
				"name": "Test Account",
				"issuer": "TestIssuer",
				"secret": "plaintext",
				"secret_base32": "JBSWY3DPEHPK3PXP",
				"type": "totp"
			}
		]
	}`

	accounts, err := ParseAuthyExport([]byte(authyJSON))
	if err != nil {
		t.Fatalf("ParseAuthyExport() error = %v", err)
	}

	if len(accounts) != 1 {
		t.Errorf("ParseAuthyExport() returned %d accounts, want 1", len(accounts))
	}

	// Should use base32 secret
	if string(accounts[0].Secret) != "JBSWY3DPEHPK3PXP" {
		t.Errorf("ParseAuthyExport() secret = %s, want JBSWY3DPEHPK3PXP", accounts[0].Secret)
	}
}

func TestIsAuthyExport(t *testing.T) {
	validAuthy := `{
		"accounts": [
			{
				"id": 12345,
				"name": "Test",
				"secret": "JBSWY3DPEHPK3PXP"
			}
		]
	}`

	if !IsAuthyExport([]byte(validAuthy)) {
		t.Error("IsAuthyExport() should return true for valid Authy export")
	}

	invalidJSON := `{invalid}`
	if IsAuthyExport([]byte(invalidJSON)) {
		t.Error("IsAuthyExport() should return false for invalid JSON")
	}
}

func TestParseAuthyEncrypted(t *testing.T) {
	_, err := ParseAuthyEncrypted([]byte{}, "password")
	if err == nil {
		t.Error("ParseAuthyEncrypted() should return error (not implemented)")
	}
}
