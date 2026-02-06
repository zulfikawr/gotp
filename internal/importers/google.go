package importers

import (
	"bufio"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/url"
	"strings"

	"github.com/google/uuid"
	"github.com/zulfikawr/gotp/internal/totp"
	"github.com/zulfikawr/gotp/internal/vault"
	"github.com/zulfikawr/gotp/pkg/base32"
	"google.golang.org/protobuf/proto"
)

// GoogleExport represents the structure of a Google Authenticator export.
// Google Authenticator exports can be in different formats.
type GoogleExport struct {
	Accounts []GoogleAccount `json:"accounts"`
}

// GoogleAccount represents a single account in a Google Authenticator export.
type GoogleAccount struct {
	Name     string `json:"name"`
	Issuer   string `json:"issuer"`
	Username string `json:"username"`
	Secret   string `json:"secret"`
	// Google sometimes uses different field names
	SecretBase32 string `json:"secret_base32"`
	// For QR code migration
	URI string `json:"uri"`
	// Algorithm
	Algorithm string `json:"algorithm"`
	// Digits
	Digits int `json:"digits"`
	// Period
	Period int `json:"period"`
}

// ParseGoogleExport parses a Google Authenticator export and returns a list of accounts.
func ParseGoogleExport(data []byte) ([]vault.Account, error) {
	// Try to parse as JSON first
	var export GoogleExport
	if err := json.Unmarshal(data, &export); err == nil && len(export.Accounts) > 0 {
		return parseGoogleJSONExport(export)
	}

	// Try to parse as otpauth-migration:// (Google Authenticator migration format)
	if strings.HasPrefix(string(data), "otpauth-migration://") {
		return parseGoogleMigrationURIs(string(data))
	}

	// Try to parse as otpauth:// URIs (one per line)
	return parseGoogleURIs(data)
}

// parseGoogleJSONExport parses JSON format Google Authenticator export.
func parseGoogleJSONExport(export GoogleExport) ([]vault.Account, error) {
	var accounts []vault.Account

	for _, acc := range export.Accounts {
		// Try to get secret from different fields
		secret := acc.Secret
		if secret == "" && acc.SecretBase32 != "" {
			secret = acc.SecretBase32
		}

		if secret == "" {
			continue // Skip accounts without secrets
		}

		vaultAcc := vault.NewAccount(acc.Name, []byte(secret))
		vaultAcc.ID = uuid.New().String()

		// Use issuer if available, otherwise use name
		if acc.Issuer != "" {
			vaultAcc.Issuer = acc.Issuer
		} else {
			vaultAcc.Issuer = acc.Name
		}

		// Username
		if acc.Username != "" {
			vaultAcc.Username = acc.Username
		}

		// Parse algorithm
		if acc.Algorithm != "" {
			switch strings.ToUpper(acc.Algorithm) {
			case "SHA1":
				vaultAcc.Algorithm = totp.SHA1
			case "SHA256":
				vaultAcc.Algorithm = totp.SHA256
			case "SHA512":
				vaultAcc.Algorithm = totp.SHA512
			}
		}

		// Parse digits
		if acc.Digits > 0 {
			vaultAcc.Digits = acc.Digits
		}

		// Parse period
		if acc.Period > 0 {
			vaultAcc.Period = acc.Period
		}

		// Add google-specific tag
		vaultAcc.Tags = []string{"google"}

		accounts = append(accounts, *vaultAcc)
	}

	return accounts, nil
}

// parseGoogleURIs parses otpauth:// URIs from Google Authenticator.
func parseGoogleURIs(data []byte) ([]vault.Account, error) {
	var accounts []vault.Account
	scanner := bufio.NewScanner(strings.NewReader(string(data)))

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		// Check if it's an otpauth:// URI
		if strings.HasPrefix(line, "otpauth://") {
			acc, err := vault.FromURI(line)
			if err != nil {
				continue // Skip invalid URIs
			}
			acc.Tags = []string{"google"}
			accounts = append(accounts, *acc)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("failed to scan URIs: %w", err)
	}

	return accounts, nil
}

// IsGoogleExport checks if the data appears to be a Google Authenticator export.
func IsGoogleExport(data []byte) bool {
	// Check for JSON format
	var export GoogleExport
	if err := json.Unmarshal(data, &export); err == nil && len(export.Accounts) > 0 {
		// Google exports typically don't have numeric 'id' fields like Authy
		// Check if accounts have name/secret but no numeric id
		for _, acc := range export.Accounts {
			// Google accounts have name and secret, but not a numeric id
			// This helps distinguish from Authy which has numeric ids
			if acc.Name != "" && acc.Secret != "" {
				return true
			}
		}
	}

	// Check for otpauth:// URIs or otpauth-migration:// URIs
	scanner := bufio.NewScanner(strings.NewReader(string(data)))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(line, "otpauth://") || strings.HasPrefix(line, "otpauth-migration://") {
			return true
		}
	}

	return false
}

// ParseGoogleMigration parses Google Authenticator migration QR code data.
// Google Authenticator uses a special format for migration QR codes.
func ParseGoogleMigration(migrationData string) ([]vault.Account, error) {
	// Google migration QR codes contain otpauth:// URIs
	// The format is typically: otpauth://totp/...?secret=...
	return parseGoogleURIs([]byte(migrationData))
}

// parseGoogleMigrationURIs parses otpauth-migration:// URIs from Google Authenticator.
// This format contains base64-encoded protobuf data in the "data" query parameter.
func parseGoogleMigrationURIs(migrationURI string) ([]vault.Account, error) {
	u, err := url.Parse(strings.TrimSpace(migrationURI))
	if err != nil {
		return nil, fmt.Errorf("failed to parse migration URI: %w", err)
	}

	dataParam := u.Query().Get("data")
	if dataParam == "" {
		return nil, fmt.Errorf("migration URI missing data parameter")
	}

	// Google migration data is base64 encoded.
	// It's often "standard" base64 but can be "URL-safe" (base64url).
	// url.Query().Get() already handles the percent-decoding (including +).
	
	// Try standard base64 first
	decoded, err := base64.StdEncoding.DecodeString(dataParam)
	if err != nil {
		// Fallback to URL-safe base64 if standard fails
		decoded, err = base64.URLEncoding.DecodeString(dataParam)
		if err != nil {
			// Some versions might skip padding
			decoded, err = base64.RawStdEncoding.DecodeString(dataParam)
			if err != nil {
				decoded, err = base64.RawURLEncoding.DecodeString(dataParam)
				if err != nil {
					return nil, fmt.Errorf("failed to decode base64 data: %w", err)
				}
			}
		}
	}

	// Parse the protobuf message
	payload := &MigrationPayload{}
	if err := proto.Unmarshal(decoded, payload); err != nil {
		return nil, fmt.Errorf("failed to unmarshal migration payload: %w", err)
	}

	// Convert protobuf accounts to vault accounts
	var accounts []vault.Account
	for _, otpAcc := range payload.OtpAccounts {
		if len(otpAcc.Secret) == 0 {
			continue
		}

		// Google secrets are raw bytes, we encode them to base32 for internal storage
		secretBase32 := base32.Encode(otpAcc.Secret)
		
		acc := vault.NewAccount(otpAcc.Name, []byte(secretBase32))
		acc.ID = uuid.New().String()

		// Use issuer if available
		if otpAcc.Issuer != "" {
			acc.Issuer = otpAcc.Issuer
		} else {
			acc.Issuer = otpAcc.Name
		}

		// Username (parsed from name if it contains a colon)
		if strings.Contains(otpAcc.Name, ":") {
			parts := strings.SplitN(otpAcc.Name, ":", 2)
			acc.Username = strings.TrimSpace(parts[1])
		}

		// Algorithm (0=Unspecified, 1=SHA1, 2=SHA256, 3=SHA512)
		switch otpAcc.Algorithm {
		case Algorithm_SHA1:
			acc.Algorithm = totp.SHA1
		case Algorithm_SHA256:
			acc.Algorithm = totp.SHA256
		case Algorithm_SHA512:
			acc.Algorithm = totp.SHA512
		default:
			acc.Algorithm = totp.SHA1
		}

		// Digits (0=Unspecified, 1=6, 2=8)
		switch otpAcc.Digits {
		case DigitCount_SIX:
			acc.Digits = 6
		case DigitCount_EIGHT:
			acc.Digits = 8
		default:
			acc.Digits = 6
		}

		// Period (0=Unspecified)
		acc.Period = 30 // Default to 30 as Google migration usually doesn't provide it reliably

		// Add google-specific tags
		acc.Tags = []string{"google", "migration"}

		accounts = append(accounts, *acc)
	}

	return accounts, nil
}
