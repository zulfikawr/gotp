package vault

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/zulfikawr/gotp/internal/totp"
)

// Secret is a custom type for TOTP secrets that handles Base32 string
// conversion during JSON marshaling/unmarshaling while keeping the
// data as a mutable byte slice in memory.
type Secret []byte

// MarshalJSON converts the secret to a Base32 string for JSON storage.
func (s Secret) MarshalJSON() ([]byte, error) {
	return json.Marshal(string(s))
}

// UnmarshalJSON converts a Base32 string from JSON into a byte slice.
func (s *Secret) UnmarshalJSON(data []byte) error {
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return err
	}
	*s = []byte(str)
	return nil
}

// Account represents a single TOTP account entry in the vault.
type Account struct {
	ID         string             `json:"id"`
	Name       string             `json:"name"`
	Issuer     string             `json:"issuer"`
	Username   string             `json:"username"`
	Secret     Secret             `json:"secret"` // Encrypted Base32 secret (stored as bytes for memory safety)
	Algorithm  totp.HashAlgorithm `json:"algorithm"`
	Digits     int                `json:"digits"`
	Period     int                `json:"period"`
	Tags       []string           `json:"tags"`
	Icon       string             `json:"icon"`
	SortOrder  int                `json:"sort_order"`
	CreatedAt  time.Time          `json:"created_at"`
	LastUsedAt time.Time          `json:"last_used_at"`
}

// NewAccount creates a new account with default values.
func NewAccount(name string, secret []byte) *Account {
	now := time.Now()
	return &Account{
		Name:       name,
		Secret:     Secret(secret),
		Algorithm:  totp.SHA1,
		Digits:     6,
		Period:     30,
		Tags:       []string{},
		CreatedAt:  now,
		LastUsedAt: now,
	}
}

// ToURI returns the otpauth:// URI representation of the account.
func (a *Account) ToURI() string {
	return fmt.Sprintf("otpauth://totp/%s:%s?secret=%s&issuer=%s&algorithm=%s&digits=%d&period=%d",
		a.Issuer, a.Username, string(a.Secret), a.Issuer, a.Algorithm, a.Digits, a.Period)
}

// FromURI parses an otpauth:// URI into an Account.
func FromURI(uriStr string) (*Account, error) {
	if !strings.HasPrefix(uriStr, "otpauth://totp/") {
		return nil, fmt.Errorf("invalid URI format: must start with otpauth://totp/")
	}

	u, err := url.Parse(uriStr)
	if err != nil {
		return nil, err
	}

	label := strings.TrimPrefix(u.Path, "/")
	issuer := u.Query().Get("issuer")
	username := label
	if strings.Contains(label, ":") {
		parts := strings.SplitN(label, ":", 2)
		if issuer == "" {
			issuer = parts[0]
		}
		username = strings.TrimSpace(parts[1])
	}

	secret := u.Query().Get("secret")
	if secret == "" {
		return nil, fmt.Errorf("missing secret in URI")
	}

	algo := totp.HashAlgorithm(strings.ToUpper(u.Query().Get("algorithm")))
	if algo == "" {
		algo = totp.SHA1
	}

	digits := 6
	if d := u.Query().Get("digits"); d != "" {
		if _, err := fmt.Sscanf(d, "%d", &digits); err != nil {
			return nil, fmt.Errorf("invalid digits: %v", err)
		}
	}

	period := 30
	if p := u.Query().Get("period"); p != "" {
		if _, err := fmt.Sscanf(p, "%d", &period); err != nil {
			return nil, fmt.Errorf("invalid period: %v", err)
		}
	}

	acc := NewAccount(username, []byte(secret))
	acc.Issuer = issuer
	acc.Username = username
	acc.Algorithm = algo
	acc.Digits = digits
	acc.Period = period

	return acc, nil
}
