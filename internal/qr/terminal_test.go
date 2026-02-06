package qr

import (
	"strings"
	"testing"
)

func TestNewTerminalQR(t *testing.T) {
	data := []byte{1, 0, 1, 0}
	qr := NewTerminalQR(data, 2, 2)

	if qr == nil {
		t.Error("NewTerminalQR() returned nil")
		return
	}

	if qr.width != 2 || qr.height != 2 {
		t.Errorf("NewTerminalQR() dimensions incorrect: got %dx%d, want 2x2", qr.width, qr.height)
	}
}

func TestTerminalQR_Render(t *testing.T) {
	data := []byte{1, 0, 1, 0}
	qr := NewTerminalQR(data, 2, 2)

	rendered := qr.Render()
	if rendered == "" {
		t.Error("Render() returned empty string")
	}

	// Check that it contains the expected characters
	if !strings.Contains(rendered, "██") {
		t.Error("Render() should contain block characters")
	}
}

func TestTerminalQR_RenderWithBorder(t *testing.T) {
	data := []byte{1, 0, 1, 0}
	qr := NewTerminalQR(data, 2, 2)

	rendered := qr.RenderWithBorder(1)
	if rendered == "" {
		t.Error("RenderWithBorder() returned empty string")
	}

	// Check that it contains border
	if !strings.Contains(rendered, "██") {
		t.Error("RenderWithBorder() should contain block characters")
	}
}

func TestTerminalQR_RenderCompact(t *testing.T) {
	data := []byte{1, 0, 1, 0, 1, 0, 1, 0}
	qr := NewTerminalQR(data, 2, 4)

	rendered := qr.RenderCompact()
	if rendered == "" {
		t.Error("RenderCompact() returned empty string")
	}

	// Check that it contains compact characters
	if !strings.Contains(rendered, "█") && !strings.Contains(rendered, "▀") && !strings.Contains(rendered, "▄") {
		t.Error("RenderCompact() should contain compact block characters")
	}
}

func TestTerminalQR_RenderAsText(t *testing.T) {
	data := []byte{1, 0, 1, 0}
	qr := NewTerminalQR(data, 2, 2)

	rendered := qr.RenderAsText()
	if rendered == "" {
		t.Error("RenderAsText() returned empty string")
	}

	// Check that it contains 0s and 1s
	if !strings.Contains(rendered, "0") || !strings.Contains(rendered, "1") {
		t.Error("RenderAsText() should contain 0s and 1s")
	}
}

func TestTerminalQR_GetDimensions(t *testing.T) {
	data := []byte{1, 0, 1, 0}
	qr := NewTerminalQR(data, 3, 4)

	width, height := qr.GetDimensions()
	if width != 3 || height != 4 {
		t.Errorf("GetDimensions() = %d, %d, want 3, 4", width, height)
	}
}
