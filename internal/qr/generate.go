package qr

import (
	"fmt"
	"os"

	"github.com/skip2/go-qrcode"
)

// GenerateQRCode generates a QR code from a string (typically an otpauth:// URI).
func GenerateQRCode(uri string, size int) ([]byte, error) {
	if size < 128 {
		size = 256 // Default size
	}

	qr, err := qrcode.New(uri, qrcode.Medium)
	if err != nil {
		return nil, fmt.Errorf("failed to create QR code: %w", err)
	}

	// Generate PNG bytes
	png, err := qr.PNG(size)
	if err != nil {
		return nil, fmt.Errorf("failed to generate PNG: %w", err)
	}

	return png, nil
}

// GenerateQRCodeToFile generates a QR code and saves it to a file.
func GenerateQRCodeToFile(uri string, filePath string, size int) error {
	png, err := GenerateQRCode(uri, size)
	if err != nil {
		return err
	}

	err = os.WriteFile(filePath, png, 0644)
	if err != nil {
		return fmt.Errorf("failed to write QR code to file: %w", err)
	}

	return nil
}

// GenerateQRCodeToTerminal generates a QR code and prints it to the terminal.
// This uses ASCII/Unicode characters to render the QR code.
func GenerateQRCodeToTerminal(uri string) error {
	qr, err := qrcode.New(uri, qrcode.Medium)
	if err != nil {
		return fmt.Errorf("failed to create QR code: %w", err)
	}

	// Use ASCII art for terminal display
	fmt.Println(qr.ToString(false))
	return nil
}
