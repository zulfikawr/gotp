package vault

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
)

// Session represents a temporary authenticated session.
type Session struct {
	Key       []byte    `json:"key"`
	ExpiresAt time.Time `json:"expires_at"`
}

// getMachineKey generates a machine-specific key for session encryption.
func getMachineKey() ([]byte, error) {
	hostname, err := os.Hostname()
	if err != nil {
		hostname = "unknown-host"
	}
	uid := os.Getuid()
	
	// Create a unique seed for this machine/user
	seed := fmt.Sprintf("%s-%d-gotp-session-secret", hostname, uid)
	hash := sha256.Sum256([]byte(seed))
	return hash[:], nil
}

// GetSessionPath returns the path to the temporary session file in the user's config directory.
func GetSessionPath() (string, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	sessionDir := filepath.Join(configDir, "gotp")
	if err := os.MkdirAll(sessionDir, 0700); err != nil {
		return "", err
	}
	return filepath.Join(sessionDir, "session.bin"), nil
}

// SaveSession encrypts and saves the derived key to a private config file.
func SaveSession(key []byte, duration time.Duration) error {
	session := Session{
		Key:       key,
		ExpiresAt: time.Now().Add(duration),
	}

	plaintext, err := json.Marshal(session)
	if err != nil {
		return err
	}

	machineKey, err := getMachineKey()
	if err != nil {
		return err
	}

	block, err := aes.NewCipher(machineKey)
	if err != nil {
		return err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return err
	}

	ciphertext := gcm.Seal(nonce, nonce, plaintext, nil)

	path, err := GetSessionPath()
	if err != nil {
		return err
	}

	return os.WriteFile(path, ciphertext, 0600)
}

// GetSession retrieves and decrypts the cached key if it hasn't expired.
func GetSession() ([]byte, error) {
	path, err := GetSessionPath()
	if err != nil {
		return nil, err
	}

	ciphertext, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}

	machineKey, err := getMachineKey()
	if err != nil {
		return nil, err
	}

	block, err := aes.NewCipher(machineKey)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return nil, fmt.Errorf("ciphertext too short")
	}

	nonce, encryptedData := ciphertext[:nonceSize], ciphertext[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, encryptedData, nil)
	if err != nil {
		// If decryption fails, the machine key might have changed (e.g. hostname change)
		// Better to clear the session than to crash.
		_ = os.Remove(path)
		return nil, nil
	}

	var session Session
	if err := json.Unmarshal(plaintext, &session); err != nil {
		return nil, err
	}

	if time.Now().After(session.ExpiresAt) {
		_ = os.Remove(path)
		return nil, nil
	}

	return session.Key, nil
}

// ClearSession removes the temporary session file.
func ClearSession() error {
	path, err := GetSessionPath()
	if err != nil {
		return err
	}
	return os.Remove(path)
}
