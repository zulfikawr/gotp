# gotp - Terminal-based TOTP Authenticator

![gotp](https://img.shields.io/badge/gotp-TOTP%20Authenticator-blue)
![Go](https://img.shields.io/badge/Go-1.25+-00ADD8)
![License](https://img.shields.io/badge/License-MIT-green)

**gotp** is a secure, cross-platform, terminal-based TOTP (Time-based One-Time Password) authenticator that allows you to manage your two-factor authentication codes.

## Features

- 🔐 **Secure Storage**: AES-256-GCM encryption with Argon2id key derivation
- 📱 **Cross-Platform**: Works on Linux, macOS, and Windows
- 💾 **Session Caching**: Avoid repeated password prompts
- 📤 **Import Support**: Aegis, Authy, Google Authenticator, and more
- 📥 **Export Support**: JSON, encrypted, and otpauth:// URIs
- 📷 **QR Code Support**: Generate and parse QR codes
- 🔗 **URI Support**: Import from `otpauth://` URIs
- 📋 **Clipboard**: Auto-copy codes to clipboard
- ⏱️ **Live Countdown**: Visual timer for code expiration

## Installation

### From Source

```bash
go install github.com/zulfikawr/gotp@latest
```

### Build from Source

```bash
git clone https://github.com/zulfikawr/gotp.git
cd gotp
go build -o gotp ./cmd/gotp
```

## Quick Start

### 1. Initialize Your Vault

```bash
gotp init
```

You'll be prompted to create a master password. This password encrypts your vault.

### 2. Add an Account

```bash
gotp add "My Account" --secret JBSWY3DPEHPK3PXP --issuer "Example"
```

Or use an otpauth:// URI:

```bash
gotp add --uri "otpauth://totp/Example:alice@google.com?secret=JBSWY3DPEHPK3PXP&issuer=Example"
```

### 3. Generate a TOTP Code

```bash
gotp get "My Account"
```

Copy to clipboard:

```bash
gotp get "My Account" --copy
```

### 4. List All Accounts

```bash
gotp list
```

### 5. Import from Other Authenticators

```bash
# Auto-detect format
gotp import aegis_backup.json

# Specify format
gotp import authy_export.json --format authy
gotp import google_export.json --format google

# Import Google Authenticator migration QR code
gotp import migration.txt --format google
```

**Note:** Google Authenticator migration format (`otpauth-migration://`) is supported and automatically decoded from protobuf data.

### 6. Generate QR Codes

```bash
# Generate QR code image
gotp qr "My Account" --output qr.png

# Display in terminal
gotp qr "My Account" --terminal

# Parse QR code from image
gotp qr --parse qr.png
```

## CLI Commands

### `gotp init`
Initialize a new vault with a master password.

**Flags:**
- `--force`, `-f`: Overwrite existing vault

### `gotp add`
Add a new TOTP account.

**Flags:**
- `--secret`, `-s`: Base32-encoded secret (required if no URI)
- `--issuer`, `-i`: Issuer name (e.g., "Google", "GitHub")
- `--username`, `-u`: Username/email
- `--algorithm`, `-a`: Hash algorithm (SHA1, SHA256, SHA512)
- `--digits`, `-d`: Number of digits (6, 7, 8)
- `--period`, `-p`: Time period in seconds (30, 60)
- `--tags`, `-t`: Comma-separated tags
- `--uri`: otpauth:// URI (alternative to manual flags)

### `gotp get`
Generate and display a TOTP code.

**Flags:**
- `--copy`, `-c`: Copy code to clipboard
- `--continuous`, `-w`: Watch mode (auto-update)
- `--qr`: Display QR code

### `gotp list`
List all accounts.

**Flags:**
- `--with-codes`: Show current TOTP codes
- `--filter`, `-f`: Filter by tag or name
- `--sort`: Sort by name, issuer, or last used

### `gotp edit`
Edit an account's details.

**Flags:**
- `--name`: New name
- `--issuer`: New issuer
- `--username`: New username
- `--secret`: New secret (requires confirmation)
- `--algorithm`: New algorithm
- `--digits`: New digit count
- `--period`: New period
- `--tags`: New tags

### `gotp remove`
Remove an account.

**Flags:**
- `--force`, `-f`: Skip confirmation

### `gotp export`
Export accounts to a file.

**Flags:**
- `--format`: Export format (json, encrypted, uri)
- `--output`, `-o`: Output file path
- `--accounts`: Specific accounts to export (comma-separated)

### `gotp import`
Import accounts from a file.

**Flags:**
- `--format`: Import format (auto, json, uri, encrypted, aegis, authy, google)

### `gotp passwd`
Change the vault master password.

### `gotp qr`
Generate or parse QR codes.

**Flags:**
- `--output`, `-o`: Output file path (for generation)
- `--size`, `-s`: QR code size in pixels (default: 256)
- `--terminal`: Display QR code in terminal
- `--parse`: Parse a QR code image file

### `gotp completion`
Generate shell completion scripts.

**Supported shells:** bash, zsh, fish, powershell

## Security

### Encryption
- **Algorithm**: AES-256-GCM (authenticated encryption)
- **Key Derivation**: Argon2id with 64MB memory, 3 iterations, 4 parallelism
- **Salt**: 16-byte random salt per vault
- **Nonce**: 12-byte random nonce per encryption

### Memory Safety
- Sensitive data (passwords, secrets) is zeroed from memory after use
- Uses Go's secure memory handling

### Session Management
- Session tokens are stored securely
- Automatic session expiration
- Session locking for multi-user systems

### Best Practices
- Use a strong master password (12+ characters, mixed case, numbers, symbols)
- Never share your vault file
- Keep backups of your vault in secure locations
- Use `gotp passwd` periodically to change your master password

## Configuration

Configuration is stored in `~/.config/gotp/config.yaml` (Linux/macOS) or `%APPDATA%\gotp\config.yaml` (Windows).

```yaml
# Default configuration
vault_path: ~/.config/gotp/vault.enc
session_timeout: 300  # 5 minutes
clipboard_timeout: 30  # 30 seconds
color: true
```

## Platform-Specific Paths

- **Linux**: `~/.config/gotp/vault.enc`
- **macOS**: `~/Library/Application Support/gotp/vault.enc`
- **Windows**: `%APPDATA%\gotp\vault.enc`

## Import Formats

### Aegis (Android)
Export from Aegis → Backup → JSON (unencrypted)

### Authy
Export from Authy app (plaintext format)

### Google Authenticator
Export via QR code migration or JSON export

### otpauth:// URIs
Standard format used by most authenticators:
```
otpauth://totp/Issuer:username?secret=SECRET&issuer=Issuer&algorithm=SHA1&digits=6&period=30
```

## Export Formats

### JSON (Plaintext)
```json
[
  {
    "id": "uuid",
    "name": "Account Name",
    "issuer": "Service",
    "username": "user@example.com",
    "secret": "JBSWY3DPEHPK3PXP",
    "algorithm": "SHA1",
    "digits": 6,
    "period": 30
  }
]
```

### Encrypted
Password-protected export using the same encryption as the vault.

### otpauth:// URIs
One URI per line, suitable for importing into other authenticators.

## Troubleshooting

### "Vault not found"
Run `gotp init` to create a new vault.

### "Invalid password"
Ensure you're using the correct master password. Passwords are case-sensitive.

### "Account not found"
Check spelling with `gotp list`. Names are case-insensitive.

### QR code parsing fails
- Ensure the image is clear and well-lit
- Try using a higher resolution image
- Verify the QR code contains a valid otpauth:// URI

## Contributing

Contributions are welcome! Please see [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

- [gozxing](https://github.com/makiuchi-d/gozxing) - QR code parsing
- [go-qrcode](https://github.com/skip2/go-qrcode) - QR code generation
- [cobra](https://github.com/spf13/cobra) - CLI framework
- [argon2](https://github.com/golang/crypto/tree/master/argon2) - Key derivation
