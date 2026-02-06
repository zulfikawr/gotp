package importers

import (
	"testing"
)

func TestParseAegisBackup(t *testing.T) {
	// Minimal Aegis backup JSON
	aegisJSON := `{
		"version": 1,
		"entries": [
			{
				"id": "test-id",
				"name": "Test Account",
				"issuer": "TestIssuer",
				"username": "test@example.com",
				"secret": "JBSWY3DPEHPK3PXP",
				"type": "totp",
				"algorithm": "SHA1",
				"digits": 6,
				"period": 30,
				"tags": ["test"]
			}
		],
		"header": {
			"slots": []
		}
	}`

	accounts, err := ParseAegisBackup([]byte(aegisJSON))
	if err != nil {
		t.Fatalf("ParseAegisBackup() error = %v", err)
	}

	if len(accounts) != 1 {
		t.Errorf("ParseAegisBackup() returned %d accounts, want 1", len(accounts))
	}

	if accounts[0].Name != "Test Account" {
		t.Errorf("ParseAegisBackup() account name = %s, want Test Account", accounts[0].Name)
	}
}

func TestParseAegisBackup_InvalidJSON(t *testing.T) {
	invalidJSON := `{invalid json}`

	_, err := ParseAegisBackup([]byte(invalidJSON))
	if err == nil {
		t.Error("ParseAegisBackup() should return error for invalid JSON")
	}
}

func TestParseAegisBackup_WrongVersion(t *testing.T) {
	wrongVersion := `{
		"version": 2,
		"entries": []
	}`

	_, err := ParseAegisBackup([]byte(wrongVersion))
	if err == nil {
		t.Error("ParseAegisBackup() should return error for wrong version")
	}
}

func TestIsAegisBackup(t *testing.T) {
	validAegis := `{
		"version": 1,
		"entries": [
			{
				"id": "test",
				"name": "Test",
				"secret": "JBSWY3DPEHPK3PXP",
				"type": "totp"
			}
		],
		"header": {"slots": []}
	}`

	if !IsAegisBackup([]byte(validAegis)) {
		t.Error("IsAegisBackup() should return true for valid Aegis backup")
	}

	invalidJSON := `{invalid}`
	if IsAegisBackup([]byte(invalidJSON)) {
		t.Error("IsAegisBackup() should return false for invalid JSON")
	}
}
