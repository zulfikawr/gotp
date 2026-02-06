package crypto

import (
	"bytes"
	"testing"
)

func TestZeroBytes(t *testing.T) {
	data := []byte{1, 2, 3, 4}
	ZeroBytes(data)
	for i, b := range data {
		if b != 0 {
			t.Errorf("Byte at index %d not zeroed", i)
		}
	}
}

func TestEncryptionRoundtrip(t *testing.T) {
	password := []byte("master-password")
	salt, _ := GenerateSalt(16)
	params := DefaultArgon2Params()

	key := DeriveKey(password, salt, params)
	plaintext := []byte("secret account data")

	ciphertext, err := Encrypt(plaintext, key)
	if err != nil {
		t.Fatalf("Encryption failed: %v", err)
	}

	decrypted, err := Decrypt(ciphertext, key)
	if err != nil {
		t.Fatalf("Decryption failed: %v", err)
	}

	if !bytes.Equal(plaintext, decrypted) {
		t.Errorf("Decrypted data doesn't match plaintext")
	}
}

func TestDecryptionFailure(t *testing.T) {
	key := make([]byte, 32)
	wrongKey := make([]byte, 32)
	wrongKey[0] = 1

	plaintext := []byte("data")
	ciphertext, _ := Encrypt(plaintext, key)

	_, err := Decrypt(ciphertext, wrongKey)
	if err == nil {
		t.Error("Expected error for decryption with wrong key, got nil")
	}
}

func TestSecureCompare(t *testing.T) {
	a := []byte("hello")
	b := []byte("hello")
	c := []byte("world")

	if !SecureCompare(a, b) {
		t.Error("SecureCompare(a, b) should be true")
	}
	if SecureCompare(a, c) {
		t.Error("SecureCompare(a, c) should be false")
	}
}

func TestGenerateSalt(t *testing.T) {
	salt, err := GenerateSalt(16)
	if err != nil {
		t.Fatalf("GenerateSalt failed: %v", err)
	}
	if len(salt) != 16 {
		t.Errorf("Expected salt length 16, got %d", len(salt))
	}
}

func TestEncryptionErrorCases(t *testing.T) {
	// Invalid key size
	_, err := Encrypt([]byte("data"), []byte("short"))
	if err == nil {
		t.Error("Expected error for invalid key size in Encrypt")
	}

	_, err = Decrypt([]byte("data"), []byte("short"))
	if err == nil {
		t.Error("Expected error for invalid key size in Decrypt")
	}

	// Ciphertext too short
	_, err = Decrypt([]byte("abc"), make([]byte, 32))
	if err == nil {
		t.Error("Expected error for short ciphertext in Decrypt")
	}
}
