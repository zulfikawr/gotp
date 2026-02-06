package importers

import (
	"testing"
)

func TestParseGoogleExport(t *testing.T) {
	googleJSON := `{
		"accounts": [
			{
				"name": "Test Account",
				"issuer": "TestIssuer",
				"username": "test@example.com",
				"secret": "JBSWY3DPEHPK3PXP",
				"algorithm": "SHA1",
				"digits": 6,
				"period": 30
			}
		]
	}`

	accounts, err := ParseGoogleExport([]byte(googleJSON))
	if err != nil {
		t.Fatalf("ParseGoogleExport() error = %v", err)
	}

	if len(accounts) != 1 {
		t.Errorf("ParseGoogleExport() returned %d accounts, want 1", len(accounts))
	}

	if accounts[0].Name != "Test Account" {
		t.Errorf("ParseGoogleExport() account name = %s, want Test Account", accounts[0].Name)
	}
}

func TestParseGoogleExport_URIs(t *testing.T) {
	uris := `otpauth://totp/Example:alice@google.com?secret=JBSWY3DPEHPK3PXP&issuer=Example
otpauth://totp/Example:bob@google.com?secret=JBSWY3DPEHPK3PXP&issuer=Example`

	accounts, err := ParseGoogleExport([]byte(uris))
	if err != nil {
		t.Fatalf("ParseGoogleExport() error = %v", err)
	}

	if len(accounts) != 2 {
		t.Errorf("ParseGoogleExport() returned %d accounts, want 2", len(accounts))
	}
}

func TestParseGoogleExport_WithBase32Secret(t *testing.T) {
	googleJSON := `{
		"accounts": [
			{
				"name": "Test Account",
				"secret_base32": "JBSWY3DPEHPK3PXP",
				"type": "totp"
			}
		]
	}`

	accounts, err := ParseGoogleExport([]byte(googleJSON))
	if err != nil {
		t.Fatalf("ParseGoogleExport() error = %v", err)
	}

	if len(accounts) != 1 {
		t.Errorf("ParseGoogleExport() returned %d accounts, want 1", len(accounts))
	}

	if string(accounts[0].Secret) != "JBSWY3DPEHPK3PXP" {
		t.Errorf("ParseGoogleExport() secret = %s, want JBSWY3DPEHPK3PXP", string(accounts[0].Secret))
	}
}

func TestIsGoogleExport(t *testing.T) {
	validGoogle := `{
		"accounts": [
			{
				"name": "Test",
				"secret": "JBSWY3DPEHPK3PXP"
			}
		]
	}`

	if !IsGoogleExport([]byte(validGoogle)) {
		t.Error("IsGoogleExport() should return true for valid Google export")
	}

	uris := `otpauth://totp/Example:alice@google.com?secret=JBSWY3DPEHPK3PXP&issuer=Example`
	if !IsGoogleExport([]byte(uris)) {
		t.Error("IsGoogleExport() should return true for URI format")
	}

	invalidJSON := `{invalid}`
	if IsGoogleExport([]byte(invalidJSON)) {
		t.Error("IsGoogleExport() should return false for invalid JSON")
	}
}

func TestParseGoogleMigration(t *testing.T) {
	migrationData := `otpauth://totp/Example:alice@google.com?secret=JBSWY3DPEHPK3PXP&issuer=Example`

	accounts, err := ParseGoogleMigration(migrationData)
	if err != nil {
		t.Fatalf("ParseGoogleMigration() error = %v", err)
	}

	if len(accounts) != 1 {
		t.Errorf("ParseGoogleMigration() returned %d accounts, want 1", len(accounts))
	}
}
