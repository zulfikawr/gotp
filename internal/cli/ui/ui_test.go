package ui

import (
	"bytes"
	"strings"
	"testing"
)

func TestUI_Basics(t *testing.T) {
	// Test SetColor
	SetColor(false)
	if Reset != "" {
		t.Error("Reset should be empty when color is disabled")
	}
	SetColor(true)
	if Reset == "" {
		t.Error("Reset should not be empty when color is enabled")
	}

	// Test Dimmed
	d := Dimmed("test")
	if !strings.Contains(d, "test") {
		t.Error("Dimmed text missing original string")
	}

	// Test PromptString
	In = strings.NewReader("input\n")
	ResetScanner()
	Out = new(bytes.Buffer)
	res := PromptString("Prompt", "default")
	if res != "input" {
		t.Errorf("Expected 'input', got %q", res)
	}

	// Test PromptRequired
	In = strings.NewReader("\ninput\n")
	ResetScanner()
	res = PromptRequired("Prompt")
	if res != "input" {
		t.Errorf("Expected 'input', got %q", res)
	}
}

func TestUI_Confirm(t *testing.T) {
	In = strings.NewReader("y\n")
	ResetScanner()
	if !PromptConfirm("Confirm", false) {
		t.Error("Expected true for 'y'")
	}

	In = strings.NewReader("n\n")
	ResetScanner()
	if PromptConfirm("Confirm", true) {
		t.Error("Expected false for 'n'")
	}
}
