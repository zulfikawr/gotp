package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestConfig_Basics(t *testing.T) {
	tmpDir, _ := os.MkdirTemp("", "gotp-config-test-*")
	defer os.RemoveAll(tmpDir)
	configPath := filepath.Join(tmpDir, "config.yaml")

	cfg := DefaultConfig()
	err := cfg.SaveConfig(configPath)
	if err != nil {
		t.Fatalf("SaveConfig failed: %v", err)
	}

	cfg2, err := LoadConfig(configPath)
	if err != nil {
		t.Fatalf("LoadConfig failed: %v", err)
	}

	if cfg2.TUI.Theme != cfg.TUI.Theme {
		t.Error("Loaded config mismatch")
	}
}

func TestPaths(t *testing.T) {
	SetVaultPathOverride("/tmp/test.enc")
	if GetVaultPath() != "/tmp/test.enc" {
		t.Error("Vault path override failed")
	}
	SetVaultPathOverride("")

	dir := GetDefaultConfigDir()
	if dir == "" {
		t.Error("Config dir should not be empty")
	}
}
