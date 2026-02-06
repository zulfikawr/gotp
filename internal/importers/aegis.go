package importers

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/zulfikawr/gotp/internal/totp"
	"github.com/zulfikawr/gotp/internal/vault"
)

// AegisBackup represents the structure of an Aegis backup JSON file.
type AegisBackup struct {
	Version int          `json:"version"`
	Entries []AegisEntry `json:"entries"`
	Header  AegisHeader  `json:"header"`
}

// AegisHeader represents the header of an Aegis backup.
type AegisHeader struct {
	Slots []AegisSlot `json:"slots"`
}

// AegisSlot represents a slot in the Aegis header.
type AegisSlot struct {
	ID       int    `json:"id"`
	Key      string `json:"key"`
	Salt     string `json:"salt"`
	Recovery string `json:"recovery"`
}

// AegisEntry represents a single entry in an Aegis backup.
type AegisEntry struct {
	ID        string            `json:"id"`
	Name      string            `json:"name"`
	Issuer    string            `json:"issuer"`
	Username  string            `json:"username"`
	Secret    string            `json:"secret"`
	Type      string            `json:"type"`
	Algorithm string            `json:"algorithm"`
	Digits    int               `json:"digits"`
	Period    int               `json:"period"`
	Note      string            `json:"note"`
	Favorite  bool              `json:"favorite"`
	Icon      string            `json:"icon"`
	IconColor string            `json:"icon_color"`
	Tags      []string          `json:"tags"`
	Info      map[string]string `json:"info"`
}

// ParseAegisBackup parses an Aegis backup JSON file and returns a list of accounts.
func ParseAegisBackup(data []byte) ([]vault.Account, error) {
	var backup AegisBackup
	if err := json.Unmarshal(data, &backup); err != nil {
		return nil, fmt.Errorf("failed to parse Aegis backup: %w", err)
	}

	if backup.Version != 1 {
		return nil, fmt.Errorf("unsupported Aegis backup version: %d (expected 1)", backup.Version)
	}

	var accounts []vault.Account
	for _, entry := range backup.Entries {
		if entry.Type != "totp" {
			continue // Skip non-TOTP entries
		}

		acc := vault.NewAccount(entry.Name, []byte(entry.Secret))
		acc.ID = entry.ID
		if acc.ID == "" {
			acc.ID = uuid.New().String()
		}
		acc.Issuer = entry.Issuer
		acc.Username = entry.Username

		// Parse algorithm
		switch strings.ToUpper(entry.Algorithm) {
		case "SHA1":
			acc.Algorithm = totp.SHA1
		case "SHA256":
			acc.Algorithm = totp.SHA256
		case "SHA512":
			acc.Algorithm = totp.SHA512
		default:
			acc.Algorithm = totp.SHA1
		}

		// Parse digits
		if entry.Digits > 0 {
			acc.Digits = entry.Digits
		}

		// Parse period
		if entry.Period > 0 {
			acc.Period = entry.Period
		}

		// Parse tags
		if len(entry.Tags) > 0 {
			acc.Tags = entry.Tags
		}

		// Parse icon
		if entry.Icon != "" {
			acc.Icon = entry.Icon
		}

		// Note field
		if entry.Note != "" {
			if acc.Tags == nil {
				acc.Tags = []string{}
			}
			acc.Tags = append(acc.Tags, "note:"+entry.Note)
		}

		accounts = append(accounts, *acc)
	}

	return accounts, nil
}

// IsAegisBackup checks if the data appears to be an Aegis backup.
func IsAegisBackup(data []byte) bool {
	var backup AegisBackup
	if err := json.Unmarshal(data, &backup); err != nil {
		return false
	}
	return backup.Version == 1 && len(backup.Entries) > 0
}
