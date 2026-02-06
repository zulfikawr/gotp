package importers

import (
	"fmt"
	"os"

	"github.com/zulfikawr/gotp/internal/vault"
)

// ImportFormat represents the type of import format.
type ImportFormat string

const (
	FormatAegis     ImportFormat = "aegis"
	FormatAuthy     ImportFormat = "authy"
	FormatGoogle    ImportFormat = "google"
	FormatJSON      ImportFormat = "json"
	FormatURI       ImportFormat = "uri"
	FormatEncrypted ImportFormat = "encrypted"
)

// ImportResult represents the result of an import operation.
type ImportResult struct {
	Accounts []vault.Account
	Format   ImportFormat
	Count    int
	Skipped  int
	Errors   []string
}

// DetectFormat detects the format of the import data.
func DetectFormat(data []byte) ImportFormat {
	// Check for Aegis first (has version field)
	if IsAegisBackup(data) {
		return FormatAegis
	}
	// Check for Authy (has accounts with numeric id field) before Google
	if IsAuthyExport(data) {
		return FormatAuthy
	}
	// Check for Google (has accounts with name/secret, or URIs)
	if IsGoogleExport(data) {
		return FormatGoogle
	}
	// Default to JSON if it parses as JSON
	return FormatJSON
}

// ImportData imports accounts from the specified format.
func ImportData(data []byte, format ImportFormat) ([]vault.Account, error) {
	switch format {
	case FormatAegis:
		return ParseAegisBackup(data)
	case FormatAuthy:
		return ParseAuthyExport(data)
	case FormatGoogle:
		return ParseGoogleExport(data)
	default:
		return nil, fmt.Errorf("unsupported import format: %s", format)
	}
}

// ImportFromFile imports accounts from a file with automatic format detection.
func ImportFromFile(filePath string) ([]vault.Account, ImportFormat, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, "", fmt.Errorf("failed to read file: %w", err)
	}

	format := DetectFormat(data)
	accounts, err := ImportData(data, format)
	if err != nil {
		return nil, "", err
	}

	return accounts, format, nil
}
