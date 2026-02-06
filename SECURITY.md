# Security Policy

## Overview

gotp is designed with security as a primary concern. This document outlines the security features, threat model, and best practices for using gotp.

## Security Features

### Encryption

#### Algorithm: AES-256-GCM
- **Mode**: Galois/Counter Mode (GCM) - authenticated encryption
- **Key Size**: 256 bits (32 bytes)
- **Nonce**: 12-byte random nonce per encryption
- **Authentication**: Built-in message authentication (MAC)

#### Key Derivation: Argon2id
- **Algorithm**: Argon2id (hybrid of Argon2i and Argon2d)
- **Memory**: 64 MB (65536 KB)
- **Iterations**: 3
- **Parallelism**: 4
- **Salt Length**: 16 bytes (random)
- **Output Key Length**: 32 bytes

### Memory Safety

#### Zeroing Sensitive Data
- Passwords are zeroed from memory immediately after use
- Encryption keys are zeroed after vault operations
- Secrets are cleared from memory when accounts are removed

#### Secure Byte Handling
- Uses `crypto/rand` for all random number generation
- No predictable random number generators
- Constant-time comparison where applicable

### Session Management

#### Session Tokens
- Session tokens are generated using cryptographically secure random
- Tokens are stored in memory only (not persisted to disk)
- Automatic expiration after configurable timeout (default: 5 minutes)

#### Session Locking
- Sessions are locked when the terminal is inactive
- Requires re-authentication to unlock

## Threat Model

### Protected Against

1. **Vault File Theft**
   - Encrypted at rest with strong encryption
   - Cannot be decrypted without the master password

2. **Memory Dump Attacks**
   - Sensitive data is zeroed from memory
   - Keys are not stored in plaintext

3. **Brute Force Attacks**
   - Argon2id is memory-hard and slow to compute
   - 64MB memory requirement makes parallel attacks difficult

4. **Man-in-the-Middle**
   - No network communication (local-only tool)
   - Clipboard operations are local

5. **Shoulder Surfing**
   - Password input is hidden (no echo)
   - Codes can be copied directly to clipboard

### Not Protected Against

1. **Keyloggers**
   - If your system is compromised with a keylogger, the master password can be captured

2. **Malware with Memory Access**
   - Advanced malware could potentially read memory before data is zeroed

3. **Physical Access**
   - Someone with physical access to your machine could potentially access the vault file

4. **Clipboard Hijacking**
   - Malware could intercept clipboard contents

## Password Requirements

### Minimum Recommendations
- **Length**: 12+ characters (20+ recommended)
- **Complexity**: Mix of uppercase, lowercase, numbers, and symbols
- **Uniqueness**: Never reuse passwords from other services
- **Memorability**: Use a passphrase (e.g., "correct-horse-battery-staple")

### Password Strength
gotp does not enforce password complexity, but you should:
- Use a password manager to generate and store your master password
- Consider using a diceware passphrase
- Avoid common passwords or patterns

## Best Practices

### Initial Setup
1. **Create a Strong Master Password**
   - Use 20+ characters
   - Consider a passphrase
   - Store it in a password manager

2. **Backup Your Vault**
   - Create regular backups using `gotp export`
   - Store backups in encrypted containers (e.g., VeraCrypt)
   - Keep backups in multiple secure locations

3. **Test Recovery**
   - Verify you can restore from backup
   - Test password recovery process

### Daily Usage
1. **Lock Your Session**
   - Sessions auto-lock after inactivity
   - Manually lock with Ctrl+C if needed

2. **Clear Clipboard**
   - Clipboard auto-clears after 30 seconds
   - Manually clear sensitive data when done

3. **Update Regularly**
   - Keep gotp updated for security patches
   - Update dependencies regularly

4. **Monitor Access**
   - Check vault file permissions (should be 600 on Unix)
   - Review system logs for unauthorized access

### Backup Strategy
1. **Frequency**: Weekly or after adding new accounts
2. **Location**: Encrypted cloud storage + offline storage
3. **Format**: Use encrypted export format
4. **Testing**: Verify backups periodically

## Encryption Details

### Vault Structure
```json
{
  "version": "1.0",
  "created_at": "2024-01-01T00:00:00Z",
  "modified_at": "2024-01-01T00:00:00Z",
  "kdf_params": {
    "memory": 65536,
    "iterations": 3,
    "parallelism": 4,
    "salt_length": 16,
    "key_length": 32
  },
  "salt": "<random 16 bytes>",
  "accounts": [
    {
      "id": "uuid",
      "name": "Account Name",
      "issuer": "Service",
      "username": "user@example.com",
      "secret": "<encrypted base32 secret>",
      "algorithm": "SHA1",
      "digits": 6,
      "period": 30
    }
  ]
}
```

### Encryption Process
1. User enters master password
2. Argon2id derives 32-byte key from password + salt
3. Vault JSON is marshaled
4. AES-256-GCM encrypts the JSON with random nonce
5. Ciphertext + nonce + salt are stored

### Decryption Process
1. User enters master password
2. Argon2id derives key from password + salt (from file)
3. AES-256-GCM decrypts ciphertext using nonce (from file)
4. JSON is unmarshaled into vault structure

## Vulnerability Reporting

### Reporting a Vulnerability
If you discover a security vulnerability in gotp:

1. **Do not** open a public issue
2. Email security concerns to: [security@gotp.dev](mailto:security@gotp.dev)
3. Include:
   - Description of the vulnerability
   - Steps to reproduce
   - Potential impact
   - Suggested fix (if any)

### Response Timeline
- **Initial Response**: Within 48 hours
- **Fix Development**: Within 7 days for critical issues
- **Public Disclosure**: Coordinated release after fix is deployed

## Compliance

### Data Protection
- gotp does not transmit data over the network
- All data is stored locally and encrypted at rest
- No telemetry or analytics are collected

### Privacy
- No account registration required
- No cloud synchronization
- No data collection

### Import Formats Security

#### Aegis (Android)
- **Format**: JSON (unencrypted backup)
- **Security**: No encryption in export format
- **Recommendation**: Delete export file after import

#### Authy
- **Format**: JSON (plaintext export)
- **Security**: No encryption in export format
- **Recommendation**: Delete export file after import

#### Google Authenticator
- **JSON Format**: Standard JSON export
- **Migration Format**: `otpauth-migration://` with protobuf-encoded data
  - **Security**: Protobuf data is base64url-encoded (not encrypted)
  - **Recommendation**: Delete migration QR code/image after import
  - **Note**: Migration format contains sensitive secrets in plaintext

#### otpauth:// URIs
- **Format**: Plain text URIs
- **Security**: Secrets are visible in plaintext
- **Recommendation**: Store in encrypted file, delete after import

## Audit

### Code Review
- All cryptographic operations are in `internal/crypto/`
- Key derivation in `internal/crypto/argon2.go`
- Encryption/decryption in `internal/crypto/aes.go`
- Memory zeroing in `internal/crypto/secure.go`

### Testing
- Unit tests for all cryptographic functions
- RFC test vectors for TOTP/HOTP
- Integration tests for vault operations

## Limitations

### Platform Security
- **Linux**: Uses system keyring if available (future enhancement)
- **macOS**: Keychain integration (future enhancement)
- **Windows**: Credential Manager (future enhancement)

### Clipboard Security
- Clipboard contents may be visible to other applications
- Some platforms clear clipboard on application exit
- Clipboard managers may persist data

### Terminal Security
- Terminal output may be captured by screen recording
- Command history may contain sensitive data
- Use `HISTCONTROL=ignorespace` for password commands

## Future Enhancements

### Planned Security Features
- Hardware security module (HSM) support
- YubiKey integration for 2FA
- System keyring integration
- Secure session tokens with hardware-backed storage
- Audit logging for vault access

## License

Security-related code is covered by the same MIT License as the project.

## Contact

For security questions: [security@gotp.dev](mailto:security@gotp.dev)
