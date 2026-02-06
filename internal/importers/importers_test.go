package importers

import (
	"testing"
)

func TestDetectFormat(t *testing.T) {
	tests := []struct {
		name     string
		data     []byte
		expected ImportFormat
	}{
		{
			name: "Aegis backup",
			data: []byte(`{
				"version": 1,
				"entries": [{"id": "test", "name": "Test", "secret": "JBSWY3DPEHPK3PXP", "type": "totp"}],
				"header": {"slots": []}
			}`),
			expected: FormatAegis,
		},
		{
			name: "Authy export",
			data: []byte(`{
				"accounts": [{"id": 12345, "name": "Test", "secret": "JBSWY3DPEHPK3PXP"}]
			}`),
			expected: FormatAuthy,
		},
		{
			name: "Google export",
			data: []byte(`{
				"accounts": [{"name": "Test", "secret": "JBSWY3DPEHPK3PXP"}]
			}`),
			expected: FormatGoogle,
		},
		{
			name:     "Google URIs",
			data:     []byte(`otpauth://totp/Example:alice@google.com?secret=JBSWY3DPEHPK3PXP&issuer=Example`),
			expected: FormatGoogle,
		},
		{
			name:     "JSON fallback",
			data:     []byte(`{"name": "Test"}`),
			expected: FormatJSON,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := DetectFormat(tt.data)
			if result != tt.expected {
				t.Errorf("DetectFormat() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestImportData(t *testing.T) {
	aegisData := []byte(`{
		"version": 1,
		"entries": [{"id": "test", "name": "Test", "secret": "JBSWY3DPEHPK3PXP", "type": "totp"}],
		"header": {"slots": []}
	}`)

	accounts, err := ImportData(aegisData, FormatAegis)
	if err != nil {
		t.Fatalf("ImportData() error = %v", err)
	}

	if len(accounts) != 1 {
		t.Errorf("ImportData() returned %d accounts, want 1", len(accounts))
	}
}

func TestImportData_UnsupportedFormat(t *testing.T) {
	_, err := ImportData([]byte("{}"), "unsupported")
	if err == nil {
		t.Error("ImportData() should return error for unsupported format")
	}
}
