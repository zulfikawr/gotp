package qr

import (
	"os"
	"testing"
)

func TestGenerateQRCode(t *testing.T) {
	uri := "otpauth://totp/Example:alice@google.com?secret=JBSWY3DPEHPK3PXP&issuer=Example"

	png, err := GenerateQRCode(uri, 256)
	if err != nil {
		t.Fatalf("GenerateQRCode() error = %v", err)
	}

	if len(png) == 0 {
		t.Error("GenerateQRCode() returned empty PNG data")
	}

	// Check PNG header
	if len(png) < 8 {
		t.Error("PNG data too short")
	}
}

func TestGenerateQRCodeToFile(t *testing.T) {
	uri := "otpauth://totp/Example:alice@google.com?secret=JBSWY3DPEHPK3PXP&issuer=Example"
	tmpFile := "/tmp/test_qr.png"
	defer os.Remove(tmpFile)

	err := GenerateQRCodeToFile(uri, tmpFile, 256)
	if err != nil {
		t.Fatalf("GenerateQRCodeToFile() error = %v", err)
	}

	// Check file exists
	if _, err := os.Stat(tmpFile); os.IsNotExist(err) {
		t.Error("QR code file was not created")
	}
}

func TestGenerateQRCodeToTerminal(t *testing.T) {
	uri := "otpauth://totp/Example:alice@google.com?secret=JBSWY3DPEHPK3PXP&issuer=Example"

	// This will print to stdout, just check it doesn't error
	err := GenerateQRCodeToTerminal(uri)
	if err != nil {
		t.Fatalf("GenerateQRCodeToTerminal() error = %v", err)
	}
}
