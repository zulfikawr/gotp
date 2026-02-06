package vault

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/zulfikawr/gotp/internal/crypto"
)

func TestVaultOperations(t *testing.T) {
	password := []byte("password123")
	salt, _ := crypto.GenerateSalt(16)
	v := NewVault(salt)

	acc := NewAccount("Test", []byte("SECRET"))
	v.Accounts = append(v.Accounts, *acc)

	// Test Roundtrip in memory
	data, err := v.Marshal(password)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	v2, err := UnmarshalVault(data, password, salt, v.KDFParams)
	if err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	if len(v2.Accounts) != 1 || v2.Accounts[0].Name != "Test" {
		t.Errorf("Recovered vault data mismatch")
	}

	// Test Storage
	tmpDir, _ := os.MkdirTemp("", "gotp-test-*")
	defer os.RemoveAll(tmpDir)
	vaultPath := filepath.Join(tmpDir, "vault.enc")

	err = SaveVault(vaultPath, v, password)
	if err != nil {
		t.Fatalf("SaveVault failed: %v", err)
	}

	v3, err := LoadVault(vaultPath, password)
	if err != nil {
		t.Fatalf("LoadVault failed: %v", err)
	}

	if v3.Accounts[0].Name != "Test" {
		t.Error("Loaded vault name mismatch")
	}
}

func TestBackupSystem(t *testing.T) {
	tmpDir, _ := os.MkdirTemp("", "gotp-backup-test-*")
	defer os.RemoveAll(tmpDir)
	vaultPath := filepath.Join(tmpDir, "vault.enc")

	if err := os.WriteFile(vaultPath, []byte("data"), 0600); err != nil {
		t.Fatalf("Failed to create test vault file: %v", err)
	}

	err := CreateBackup(vaultPath, 3)
	if err != nil {
		t.Fatalf("CreateBackup failed: %v", err)
	}

	matches, _ := filepath.Glob(vaultPath + ".*.bak")
	if len(matches) != 1 {
		t.Errorf("Expected 1 backup file, got %d", len(matches))
	}
}

func TestAccountURI(t *testing.T) {
	acc := NewAccount("user@example.com", []byte("JBSWY3DPEHPK3PXP"))
	acc.Issuer = "GitHub"
	acc.Username = "user@example.com"

	uri := acc.ToURI()
	acc2, err := FromURI(uri)
	if err != nil {
		t.Fatalf("FromURI failed: %v", err)
	}

	if !bytes.Equal(acc2.Secret, acc.Secret) || acc2.Username != acc.Username || acc2.Issuer != acc.Issuer {
		t.Errorf("Roundtrip URI conversion failed. Got: %+v, Expected: %+v", acc2, acc)
	}
}

func TestSessionManagement(t *testing.T) {
	key := []byte("secret-key-32-bytes-long-exactly!!")
	err := SaveSession(key, 1*time.Second)
	if err != nil {
		t.Fatalf("SaveSession failed: %v", err)
	}

	cached, err := GetSession()
	if err != nil {
		t.Fatalf("GetSession failed: %v", err)
	}

	if !bytes.Equal(cached, key) {
		t.Error("Cached key mismatch")
	}

	_ = ClearSession()
	cached, _ = GetSession()
	if cached != nil {
		t.Error("Session should be cleared")
	}
}

func TestLoadVaultInteractive(t *testing.T) {
	tmpDir, _ := os.MkdirTemp("", "gotp-interactive-test-*")
	defer os.RemoveAll(tmpDir)
	vaultPath := filepath.Join(tmpDir, "vault.enc")

	password := []byte("password")
	salt, _ := crypto.GenerateSalt(16)
	v := NewVault(salt)
	if err := SaveVault(vaultPath, v, password); err != nil {
		t.Fatalf("SaveVault failed: %v", err)
	}

	// Test loading with password prompt
	promptFunc := func(p string) ([]byte, error) {
		return password, nil
	}

	v2, key, err := LoadVaultInteractive(vaultPath, promptFunc)
	if err != nil {
		t.Fatalf("LoadVaultInteractive failed: %v", err)
	}

	if v2 == nil || key == nil {
		t.Fatal("Expected vault and key")
	}

	// Test loading from session (should not call prompt)
	v3, key3, err := LoadVaultInteractive(vaultPath, func(p string) ([]byte, error) {
		t.Fatal("Prompt should not be called when session exists")
		return nil, nil
	})

	if err != nil || v3 == nil || key3 == nil {
		t.Fatal("Failed to load from session")
	}
}
