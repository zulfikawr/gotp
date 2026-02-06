package qr

import (
	"strings"
)

// TerminalQR represents a QR code rendered for terminal display.
type TerminalQR struct {
	// QR code data (encoded as a string of 0s and 1s)
	data []byte
	// Width of the QR code in modules
	width int
	// Height of the QR code in modules
	height int
}

// NewTerminalQR creates a new terminal QR code from raw QR data.
func NewTerminalQR(data []byte, width, height int) *TerminalQR {
	return &TerminalQR{
		data:   data,
		width:  width,
		height: height,
	}
}

// Render renders the QR code as ASCII/Unicode art for terminal display.
// Uses block characters for better visual representation.
func (t *TerminalQR) Render() string {
	var builder strings.Builder

	// Use Unicode block characters for better density
	// █ (U+2588) for black, space for white
	for y := 0; y < t.height; y++ {
		for x := 0; x < t.width; x++ {
			idx := y*t.width + x
			if idx < len(t.data) && t.data[idx] == 1 {
				builder.WriteString("██")
			} else {
				builder.WriteString("  ")
			}
		}
		builder.WriteString("\n")
	}

	return builder.String()
}

// RenderWithBorder renders the QR code with a border for better visibility.
func (t *TerminalQR) RenderWithBorder(borderSize int) string {
	var builder strings.Builder

	// Top border
	borderLine := strings.Repeat("██", t.width+borderSize*2) + "\n"
	for i := 0; i < borderSize; i++ {
		builder.WriteString(borderLine)
	}

	// QR code with side borders
	for y := 0; y < t.height; y++ {
		// Left border
		builder.WriteString(strings.Repeat("██", borderSize))

		// QR code line
		for x := 0; x < t.width; x++ {
			idx := y*t.width + x
			if idx < len(t.data) && t.data[idx] == 1 {
				builder.WriteString("██")
			} else {
				builder.WriteString("  ")
			}
		}

		// Right border
		builder.WriteString(strings.Repeat("██", borderSize))
		builder.WriteString("\n")
	}

	// Bottom border
	for i := 0; i < borderSize; i++ {
		builder.WriteString(borderLine)
	}

	return builder.String()
}

// RenderCompact renders a compact version of the QR code using half-block characters.
// This allows for higher density display.
func (t *TerminalQR) RenderCompact() string {
	var builder strings.Builder

	// Use half-block characters for higher density
	// Upper half: ▀ (U+2580), Lower half: ▄ (U+2584), Full: █ (U+2588), Empty: space
	for y := 0; y < t.height; y += 2 {
		for x := 0; x < t.width; x++ {
			topIdx := y*t.width + x
			bottomIdx := (y+1)*t.width + x

			topBit := topIdx < len(t.data) && t.data[topIdx] == 1
			bottomBit := bottomIdx < len(t.data) && t.data[bottomIdx] == 1

			switch {
			case topBit && bottomBit:
				builder.WriteString("█")
			case topBit && !bottomBit:
				builder.WriteString("▀")
			case !topBit && bottomBit:
				builder.WriteString("▄")
			default:
				builder.WriteString(" ")
			}
		}
		builder.WriteString("\n")
	}

	return builder.String()
}

// RenderAsText renders the QR code as simple text (0 for white, 1 for black).
func (t *TerminalQR) RenderAsText() string {
	var builder strings.Builder

	for y := 0; y < t.height; y++ {
		for x := 0; x < t.width; x++ {
			idx := y*t.width + x
			if idx < len(t.data) && t.data[idx] == 1 {
				builder.WriteString("1")
			} else {
				builder.WriteString("0")
			}
		}
		builder.WriteString("\n")
	}

	return builder.String()
}

// GetDimensions returns the width and height of the QR code.
func (t *TerminalQR) GetDimensions() (int, int) {
	return t.width, t.height
}
