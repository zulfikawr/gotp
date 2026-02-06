package qr

import (
	"bytes"
	"fmt"
	"image"
	_ "image/jpeg" // Register JPEG decoder
	_ "image/png"  // Register PNG decoder
	"os"
	"strings"

	"github.com/makiuchi-d/gozxing"
	"github.com/makiuchi-d/gozxing/qrcode"
)

// ParseImageFile parses a QR code from an image file (PNG, JPEG).
func ParseImageFile(filePath string) (string, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to read image file: %w", err)
	}

	return ParseImageBytes(data)
}

// ParseImageBytes parses a QR code from image bytes.
func ParseImageBytes(data []byte) (string, error) {
	img, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		return "", fmt.Errorf("failed to decode image: %w", err)
	}

	return ParseImage(img)
}

// ParseImage parses a QR code from an image.Image.
func ParseImage(img image.Image) (string, error) {
	// Create a binary bitmap from the image
	binaryBitmap, err := gozxing.NewBinaryBitmapFromImage(img)
	if err != nil {
		return "", fmt.Errorf("failed to create binary bitmap: %w", err)
	}

	// Create a QR code reader
	reader := qrcode.NewQRCodeReader()

	// Decode the QR code
	result, err := reader.Decode(binaryBitmap, nil)
	if err != nil {
		return "", fmt.Errorf("failed to decode QR code: %w", err)
	}

	return result.GetText(), nil
}

// ParseTerminalQR parses a terminal-rendered QR code (ASCII/Unicode art).
// This is useful for importing QR codes that were displayed in terminal.
func ParseTerminalQR(terminalQR string) (string, error) {
	// This is a simplified parser that expects a specific format
	// In practice, users would scan the QR code with a phone camera
	// and extract the URI, or use ParseImageFile for image files
	return "", fmt.Errorf("terminal QR parsing is not supported. Please use image files instead")
}

// ValidateOTPAuthURI validates that the parsed data is a valid otpauth:// URI.
func ValidateOTPAuthURI(uri string) error {
	if len(uri) < 10 {
		return fmt.Errorf("URI too short")
	}
	// Check for both otpauth:// and otpauth-migration:// formats
	if uri[:10] != "otpauth://" && !strings.HasPrefix(uri, "otpauth-migration://") {
		return fmt.Errorf("invalid URI format: must start with otpauth:// or otpauth-migration://")
	}
	return nil
}
