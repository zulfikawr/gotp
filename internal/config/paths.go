package config

import (
	"os"
	"path/filepath"
	"runtime"
)

var vaultPathOverride string

// SetVaultPathOverride allows overriding the vault path for testing purposes.
func SetVaultPathOverride(path string) {
	vaultPathOverride = path
}

// GetDefaultConfigDir returns the platform-specific default configuration directory.
func GetDefaultConfigDir() string {
	var path string
	home, _ := os.UserHomeDir()

	switch runtime.GOOS {
	case "windows":
		path = filepath.Join(os.Getenv("APPDATA"), "gotp")
	case "darwin":
		path = filepath.Join(home, "Library", "Application Support", "gotp")
	default: // Linux and others
		path = filepath.Join(home, ".config", "gotp")
	}

	return path
}

// GetVaultPath returns the full path to the default vault file or the override if set.
func GetVaultPath() string {
	if vaultPathOverride != "" {
		return vaultPathOverride
	}
	return filepath.Join(GetDefaultConfigDir(), "vault.enc")
}

// GetConfigPath returns the full path to the default configuration file.
func GetConfigPath() string {
	return filepath.Join(GetDefaultConfigDir(), "config.yaml")
}
