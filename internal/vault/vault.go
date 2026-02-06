package vault

import (
	"encoding/json"
	"time"

	"github.com/zulfikawr/gotp/internal/crypto"
)

// Vault represents the top-level structure of the encrypted vault.
type Vault struct {
	Version    string              `json:"version"`
	CreatedAt  time.Time           `json:"created_at"`
	ModifiedAt time.Time           `json:"modified_at"`
	KDFParams  crypto.Argon2Params `json:"kdf_params"`
	Salt       []byte              `json:"salt"`
	Accounts   []Account           `json:"accounts"`
}

// NewVault creates a new, empty vault with default parameters.
func NewVault(salt []byte) *Vault {
	now := time.Now()
	return &Vault{
		Version:    "1.0",
		CreatedAt:  now,
		ModifiedAt: now,
		KDFParams:  crypto.DefaultArgon2Params(),
		Salt:       salt,
		Accounts:   []Account{},
	}
}

// Marshal serializes the vault into an encrypted JSON blob using a password.
func (v *Vault) Marshal(password []byte) ([]byte, error) {
	key := crypto.DeriveKey(password, v.Salt, v.KDFParams)
	defer crypto.ZeroBytes(key)
	return v.MarshalWithKey(key)
}

// MarshalWithKey serializes the vault using a pre-derived key.
func (v *Vault) MarshalWithKey(key []byte) ([]byte, error) {
	v.ModifiedAt = time.Now()
	plaintext, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}
	return crypto.Encrypt(plaintext, key)
}

// UnmarshalVault decrypts and deserializes a vault from an encrypted blob.
func UnmarshalVault(data []byte, password []byte, salt []byte, params crypto.Argon2Params) (*Vault, error) {
	key := crypto.DeriveKey(password, salt, params)
	defer crypto.ZeroBytes(key)

	plaintext, err := crypto.Decrypt(data, key)
	if err != nil {
		return nil, err
	}

	var v Vault
	if err := json.Unmarshal(plaintext, &v); err != nil {
		return nil, err
	}

	return &v, nil
}
