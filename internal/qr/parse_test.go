package qr

import (
	"testing"
)

func TestValidateOTPAuthURI(t *testing.T) {
	tests := []struct {
		name    string
		uri     string
		wantErr bool
	}{
		{
			name:    "valid otpauth URI",
			uri:     "otpauth://totp/Example:alice@google.com?secret=JBSWY3DPEHPK3PXP&issuer=Example",
			wantErr: false,
		},
		{
			name:    "invalid URI - too short",
			uri:     "otpauth:/",
			wantErr: true,
		},
		{
			name:    "invalid URI - wrong prefix",
			uri:     "http://example.com",
			wantErr: true,
		},
		{
			name:    "empty URI",
			uri:     "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateOTPAuthURI(tt.uri)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateOTPAuthURI() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestParseTerminalQR(t *testing.T) {
	_, err := ParseTerminalQR("some terminal QR code")
	if err == nil {
		t.Error("ParseTerminalQR() should return an error")
	}
}
