package vault

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/zulfikawr/gotp/internal/crypto"
)

// VaultMetadata stores the unencrypted part of the vault required to decrypt it.
type VaultMetadata struct {
	Salt       []byte              `json:"salt"`
	KDFParams  crypto.Argon2Params `json:"kdf_params"`
	Ciphertext []byte              `json:"ciphertext"`
}

// SaveVault writes the vault to its encrypted file atomically using a password.
func SaveVault(path string, vault *Vault, password []byte) error {
	key := crypto.DeriveKey(password, vault.Salt, vault.KDFParams)
	defer crypto.ZeroBytes(key)
	return SaveVaultWithKey(path, vault, key)
}

// SaveVaultWithKey writes the vault to its encrypted file atomically using a pre-derived key.
func SaveVaultWithKey(path string, vault *Vault, key []byte) error {
	ciphertext, err := vault.MarshalWithKey(key)
	if err != nil {
		return err
	}

	metadata := VaultMetadata{
		Salt:       vault.Salt,
		KDFParams:  vault.KDFParams,
		Ciphertext: ciphertext,
	}

	data, err := json.Marshal(metadata)
	if err != nil {
		return err
	}

	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0700); err != nil {
		return err
	}

	// Use CreateTemp to ensure a unique temporary file and avoid race conditions
	f, err := os.CreateTemp(dir, "vault-*.tmp")
	if err != nil {
		return fmt.Errorf("failed to create temp file: %w", err)
	}
	tmpName := f.Name()
	
	// Ensure cleanup if we fail before the rename
	defer func() {
		if _, err := os.Stat(tmpName); err == nil {
			_ = os.Remove(tmpName)
		}
	}()

	if _, err := f.Write(data); err != nil {
		f.Close()
		return fmt.Errorf("failed to write to temp file: %w", err)
	}
	
	// Sync to ensure data is actually on disk before we rename
	if err := f.Sync(); err != nil {
		f.Close()
		return fmt.Errorf("failed to sync temp file: %w", err)
	}
	f.Close()

	if err := os.Chmod(tmpName, 0600); err != nil {
		return fmt.Errorf("failed to set permissions on temp file: %w", err)
	}

	return os.Rename(tmpName, path)
}

// LoadVault reads and decrypts the vault from a file using a password.
func LoadVault(path string, password []byte) (*Vault, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var metadata VaultMetadata
	if err := json.Unmarshal(data, &metadata); err != nil {
		return nil, err
	}

	return UnmarshalVault(metadata.Ciphertext, password, metadata.Salt, metadata.KDFParams)
}

// LoadVaultWithKey reads and decrypts the vault using a pre-derived key.
func LoadVaultWithKey(path string, key []byte) (*Vault, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var metadata VaultMetadata
	if err := json.Unmarshal(data, &metadata); err != nil {
		return nil, err
	}

	plaintext, err := crypto.Decrypt(metadata.Ciphertext, key)
	if err != nil {
		return nil, err
	}

	var v Vault
	if err := json.Unmarshal(plaintext, &v); err != nil {
		return nil, err
	}

	return &v, nil
}

// LoadVaultInteractive attempts to load the vault using a session key, or prompts for a password if needed.
func LoadVaultInteractive(path string, promptFunc func(string) ([]byte, error)) (*Vault, []byte, error) {
	key, _ := GetSession()
	if key != nil {
		v, err := LoadVaultWithKey(path, key)
		if err == nil {
			return v, key, nil
		}
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, nil, fmt.Errorf("could not read vault file: %v", err)
	}

	var metadata VaultMetadata
	if err := json.Unmarshal(data, &metadata); err != nil {
		return nil, nil, fmt.Errorf("failed to parse vault metadata: %v", err)
	}

	// If unmarshaling failed or file is older format, Salt might be empty.
	// We'll provide a clearer error message.
	if len(metadata.Salt) == 0 {
		return nil, nil, fmt.Errorf("vault file is corrupted or in an incompatible format (missing salt)")
	}

	password, err := promptFunc("Enter master password: ")
	if err != nil {
		return nil, nil, err
	}

	key = crypto.DeriveKey(password, metadata.Salt, metadata.KDFParams)
	v, err := LoadVaultWithKey(path, key)
	if err != nil {
		return nil, nil, fmt.Errorf("invalid master password")
	}

	_ = SaveSession(key, 5*time.Minute)

	return v, key, nil
}
