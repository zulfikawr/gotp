package config

import (
	"os"

	"github.com/zulfikawr/gotp/internal/totp"
	"gopkg.in/yaml.v3"
)

// Config represents the application configuration.
type Config struct {
	General  GeneralConfig  `yaml:"general"`
	CLI      CLIConfig      `yaml:"cli"`
	TUI      TUIConfig      `yaml:"tui"`
	Security SecurityConfig `yaml:"security"`
}

type GeneralConfig struct {
	DefaultDigits       int                `yaml:"default_digits"`
	DefaultPeriod       int                `yaml:"default_period"`
	DefaultAlgorithm    totp.HashAlgorithm `yaml:"default_algorithm"`
	AutoCopy            bool               `yaml:"auto_copy"`
	ClearClipboardAfter int                `yaml:"clear_clipboard_after"`
	SessionTimeout      int                `yaml:"session_timeout"`
}

type CLIConfig struct {
	Color      bool   `yaml:"color"`
	JSONOutput bool   `yaml:"json_output"`
	DateFormat string `yaml:"date_format"`
}

type TUIConfig struct {
	Theme           string `yaml:"theme"`
	ShowCodesInList bool   `yaml:"show_codes_in_list"`
	ConfirmDelete   bool   `yaml:"confirm_delete"`
	AnimateProgress bool   `yaml:"animate_progress"`
	RefreshRate     int    `yaml:"refresh_rate"`
}

type SecurityConfig struct {
	Argon2Memory      uint32 `yaml:"argon2_memory"`
	Argon2Iterations  uint32 `yaml:"argon2_iterations"`
	Argon2Parallelism uint8  `yaml:"argon2_parallelism"`
	BackupCount       int    `yaml:"backup_count"`
	AutoLock          bool   `yaml:"auto_lock"`
	AutoLockTimeout   int    `yaml:"auto_lock_timeout"`
}

// DefaultConfig returns the default configuration.
func DefaultConfig() *Config {
	return &Config{
		General: GeneralConfig{
			DefaultDigits:       6,
			DefaultPeriod:       30,
			DefaultAlgorithm:    totp.SHA1,
			AutoCopy:            false,
			ClearClipboardAfter: 30,
			SessionTimeout:      300,
		},
		CLI: CLIConfig{
			Color:      true,
			JSONOutput: false,
			DateFormat: "2006-01-02 15:04:05",
		},
		TUI: TUIConfig{
			Theme:           "dark",
			ShowCodesInList: true,
			ConfirmDelete:   true,
			AnimateProgress: true,
			RefreshRate:     100,
		},
		Security: SecurityConfig{
			Argon2Memory:      65536,
			Argon2Iterations:  3,
			Argon2Parallelism: 4,
			BackupCount:       3,
			AutoLock:          true,
			AutoLockTimeout:   300,
		},
	}
}

// LoadConfig loads the configuration from a YAML file.
func LoadConfig(path string) (*Config, error) {
	cfg := DefaultConfig()
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return cfg, nil
		}
		return nil, err
	}

	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}

// SaveConfig saves the configuration to a YAML file.
func (c *Config) SaveConfig(path string) error {
	data, err := yaml.Marshal(c)
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}
