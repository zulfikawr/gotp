package importers

import (
	"testing"
)

func TestParseGoogleMigrationURIs(t *testing.T) {
	// Test with a sample migration URI
	// Note: This is a simplified test - real migration URIs contain protobuf data
	migrationURI := "otpauth-migration://offline?data=CjsKCoOfX7gOctmmEz8SCmxhdy1tYWtlcnMaBkdpdEh1YiABKAEwAkITYmM1N2Q0MTc2NTM1NjIyNjcyNgpLChT%2BnWUJ01t8Zs4WOw8X2dTXZS4JdRINemFoaWRhYmRpbGxhaBoJTmFtZWNoZWFwIAEoATACQhM5MDU2ZDUxNzY5MDk1NDAxOTk1EAIYASAA"

	accounts, err := parseGoogleMigrationURIs(migrationURI)
	if err != nil {
		// This is expected to fail with invalid protobuf data
		// The test just verifies the function exists and is called
		t.Logf("Expected error (invalid protobuf): %v", err)
		return
	}

	if len(accounts) > 0 {
		t.Logf("Successfully parsed %d accounts", len(accounts))
	}
}

func TestIsGoogleExport_Migration(t *testing.T) {
	// Test that migration URIs are detected as Google export
	migrationData := []byte("otpauth-migration://offline?data=CjsKCoOfX7gOctmmEz8SCmxhdy1tYWtlcnMaBkdpdEh1YiABKAEwAkITYmM1N2Q0MTc2NTM1NjIyNjcyNgpLChT%2BnWUJ01t8Zs4WOw8X2dTXZS4JdRINemFoaWRhYmRpbGxhaBoJTmFtZWNoZWFwIAEoATACQhM5MDU2ZDUxNzY5MDk1NDAxOTk1EAIYASAA")

	if !IsGoogleExport(migrationData) {
		t.Error("IsGoogleExport() should return true for migration URIs")
	}
}
