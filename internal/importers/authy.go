package importers

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/zulfikawr/gotp/internal/totp"
	"github.com/zulfikawr/gotp/internal/vault"
)

// AuthyExport represents the structure of an Authy export.
// Authy exports can be in different formats, this handles the common JSON format.
type AuthyExport struct {
	Accounts []AuthyAccount `json:"accounts"`
}

// AuthyAccount represents a single account in an Authy export.
type AuthyAccount struct {
	ID        int64  `json:"id"`
	Name      string `json:"name"`
	Issuer    string `json:"issuer"`
	Username  string `json:"username"`
	Secret    string `json:"secret"`
	Algorithm string `json:"algorithm"`
	Digits    int    `json:"digits"`
	Period    int    `json:"period"`
	Type      string `json:"type"`
	// Authy-specific fields
	OriginalName string `json:"original_name"`
	// Sometimes the secret is base32 encoded
	SecretBase32 string `json:"secret_base32"`
}

// ParseAuthyExport parses an Authy export and returns a list of accounts.
func ParseAuthyExport(data []byte) ([]vault.Account, error) {
	var export AuthyExport
	if err := json.Unmarshal(data, &export); err != nil {
		return nil, fmt.Errorf("failed to parse Authy export: %w", err)
	}

	var accounts []vault.Account
	for _, acc := range export.Accounts {
		if acc.Type != "totp" && acc.Type != "" {
			continue // Skip non-TOTP entries
		}

		vaultAcc := vault.NewAccount(acc.Name, []byte(acc.Secret))
		vaultAcc.ID = uuid.New().String()
		vaultAcc.Issuer = acc.Issuer
		vaultAcc.Username = acc.Username

		// Handle secret - Authy might use base32 encoded secrets
		if acc.SecretBase32 != "" {
			vaultAcc.Secret = vault.Secret(acc.SecretBase32)
		} else if acc.Secret != "" {
			vaultAcc.Secret = vault.Secret(acc.Secret)
		}

		// Parse algorithm
		switch strings.ToUpper(acc.Algorithm) {
		case "SHA1":
			vaultAcc.Algorithm = totp.SHA1
		case "SHA256":
			vaultAcc.Algorithm = totp.SHA256
		case "SHA512":
			vaultAcc.Algorithm = totp.SHA512
		default:
			vaultAcc.Algorithm = totp.SHA1
		}

		// Parse digits
		if acc.Digits > 0 {
			vaultAcc.Digits = acc.Digits
		}

		// Parse period
		if acc.Period > 0 {
			vaultAcc.Period = acc.Period
		}

		// Add authy-specific tag
		if acc.OriginalName != "" {
			vaultAcc.Tags = []string{"authy", "original:" + acc.OriginalName}
		} else {
			vaultAcc.Tags = []string{"authy"}
		}

		accounts = append(accounts, *vaultAcc)
	}

	return accounts, nil
}

// IsAuthyExport checks if the data appears to be an Authy export.
func IsAuthyExport(data []byte) bool {
	var export AuthyExport
	if err := json.Unmarshal(data, &export); err != nil {
		return false
	}
	// Authy exports have accounts with an 'id' field
	if len(export.Accounts) == 0 {
		return false
	}
	// Check if at least one account has an ID field
	for _, acc := range export.Accounts {
		if acc.ID != 0 {
			return true
		}
	}
	return false
}

// ParseAuthyEncrypted parses an encrypted Authy backup.
// Note: Authy's encryption is proprietary and requires the Authy app to decrypt.
// This function provides a placeholder for future implementation.
func ParseAuthyEncrypted(data []byte, password string) ([]vault.Account, error) {
	// Authy uses a proprietary encryption format
	// In practice, users need to export from Authy app directly
	return nil, fmt.Errorf("encrypted Authy backups are not supported. Please export from Authy app in plaintext format")
}
