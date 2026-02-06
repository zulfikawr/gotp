# Architecture Documentation

## Project Structure

```
gotp/
├── cmd/gotp/              # Application entry point
│   └── main.go           # Main function, CLI initialization
├── internal/
│   ├── cli/              # CLI interface and commands
│   │   ├── commands/     # Individual CLI commands
│   │   │   ├── add.go
│   │   │   ├── edit.go
│   │   │   ├── export.go
│   │   │   ├── get.go
│   │   │   ├── import.go
│   │   │   ├── init.go
│   │   │   ├── list.go
│   │   │   ├── passwd.go
│   │   │   ├── qr.go
│   │   │   └── remove.go
│   │   ├── ui/           # User interface components
│   │   │   ├── progress.go
│   │   │   ├── prompt.go
│   │   │   ├── table.go
│   │   │   └── ui_test.go
│   │   └── cli.go        # Root command setup
│   ├── config/           # Configuration management
│   │   ├── config.go
│   │   ├── config_test.go
│   │   └── paths.go
│   ├── crypto/           # Cryptographic operations
│   │   ├── aes.go        # AES-256-GCM encryption
│   │   ├── argon2.go     # Argon2id key derivation
│   │   ├── crypto_test.go
│   │   └── secure.go     # Memory safety utilities
│   ├── importers/        # Import from other authenticators
│   │   ├── aegis.go      # Aegis backup format
│   │   ├── authy.go      # Authy export format
│   │   ├── google.go     # Google Authenticator (JSON + migration)
│   │   ├── importers.go  # Format detection & coordination
│   │   ├── migration.pb.go  # Protobuf for migration format
│   │   └── *_test.go     # Import tests
│   ├── qr/               # QR code operations
│   │   ├── generate.go   # QR code generation
│   │   ├── parse.go      # QR code parsing
│   │   ├── terminal.go   # Terminal QR display
│   │   └── *_test.go
│   ├── totp/             # TOTP/HOTP implementation
│   │   ├── hmac.go       # HMAC implementation
│   │   ├── hotp.go       # HOTP (RFC 4226)
│   │   ├── totp.go       # TOTP (RFC 6238)
│   │   └── *_test.go
│   └── vault/            # Vault management
│       ├── account.go    # Account data structure
│       ├── backup.go     # Backup system
│       ├── session.go    # Session management
│       ├── storage.go    # File I/O operations
│       ├── vault.go      # Vault structure and operations
│       └── *_test.go
├── pkg/
│   └── base32/           # Base32 encoding/decoding
│       ├── base32.go
│       └── base32_test.go
├── go.mod
├── go.sum
├── README.md
├── SECURITY.md
├── ARCHITECTURE.md
├── CONTRIBUTING.md
└── PHASES.md
```

## Module Responsibilities

### `cmd/gotp/main.go`
**Purpose**: Application entry point
**Responsibilities**:
- Initialize CLI framework
- Register all commands
- Set version string
- Handle global error handling
- Exit with appropriate status codes

### `internal/cli/`
**Purpose**: CLI interface and command routing
**Responsibilities**:
- Command registration and routing
- Global flags handling
- Help text generation
- UI component orchestration

### `internal/cli/commands/`
**Purpose**: Individual CLI commands
**Responsibilities**:
- Command-specific flag parsing
- Business logic for each operation
- User interaction (prompts, confirmations)
- Error handling and reporting

### `internal/cli/ui/`
**Purpose**: User interface components
**Responsibilities**:
- Terminal formatting (colors, styles)
- Table rendering
- Progress indicators
- Password prompts (secure input)
- Message output (success, error, warning)

### `internal/config/`
**Purpose**: Configuration management
**Responsibilities**:
- Load/save configuration files
- Platform-specific path resolution
- Default configuration values
- Configuration validation

### `internal/crypto/`
**Purpose**: Cryptographic operations
**Responsibilities**:
- Key derivation (Argon2id)
- Encryption/decryption (AES-256-GCM)
- Memory safety (zeroing sensitive data)
- Random number generation

### `internal/importers/`
**Purpose**: Import from other authenticators
**Responsibilities**:
- Parse Aegis backup format
- Parse Authy export format
- Parse Google Authenticator format
- Format detection
- Data transformation to gotp format

### `internal/qr/`
**Purpose**: QR code operations
**Responsibilities**:
- QR code generation from URI
- QR code parsing from images
- Terminal QR display
- URI validation

### `internal/totp/`
**Purpose**: TOTP/HOTP implementation
**Responsibilities**:
- HMAC implementation (SHA1, SHA256, SHA512)
- HOTP generation (RFC 4226)
- TOTP generation (RFC 6238)
- Code validation with time windows
- Time remaining calculation

### `internal/vault/`
**Purpose**: Vault management
**Responsibilities**:
- Account CRUD operations
- Vault encryption/decryption
- File I/O operations
- Backup management
- Session management
- Data validation

### `pkg/base32/`
**Purpose**: Base32 encoding/decoding
**Responsibilities**:
- RFC 4648 compliant Base32
- Padding handling
- Error detection

## Data Flow

### Vault Creation Flow
```
User → CLI (init) → Vault.NewVault() → Crypto.DeriveKey() →
Vault.Marshal() → Crypto.Encrypt() → Storage.SaveVault()
```

### Vault Unlock Flow
```
User → CLI (any command) → Storage.LoadVault() →
Crypto.DeriveKey() → Vault.UnmarshalVault() →
Crypto.Decrypt() → JSON.Unmarshal() → Vault object
```

### TOTP Generation Flow
```
User → CLI (get) → Vault.LoadVault() → Account.ToURI() →
TOTP.GenerateTOTP() → HOTP.GenerateHOTP() →
HMAC.Compute() → Dynamic Truncation → Code Output
```

### Import Flow
```
User → CLI (import) → Importers.DetectFormat() →
Importers.ImportData() → Parse specific format →
Transform to Account → Vault.AddAccount() →
Vault.SaveVault()
```

### QR Code Flow
```
User → CLI (qr) → QR.GenerateQRCode() →
QR.GenerateQRCodeToFile() → Save PNG OR
QR.GenerateQRCodeToTerminal() → Display ASCII
```

## Design Patterns

### Command Pattern
Each CLI command is implemented as a separate function that returns a `cobra.Command`. This allows for:
- Clean separation of concerns
- Easy testing
- Modular command registration

### Strategy Pattern
Used in:
- **Import formats**: Different parsers for different formats
- **Clipboard operations**: Platform-specific implementations
- **Storage paths**: Platform-specific path resolution

### Factory Pattern
- **Account creation**: `NewAccount()` factory function
- **Vault creation**: `NewVault()` factory function
- **Command creation**: `New*Cmd()` factory functions

### Repository Pattern
The `vault` package acts as a repository:
- Abstracts storage operations
- Handles encryption/decryption
- Manages account lifecycle

## Security Architecture

### Encryption Layer
```
┌─────────────────────────────────────┐
│         User Interface              │
└─────────────────────────────────────┘
                  ↓
┌─────────────────────────────────────┐
│      Password Prompt (Secure)       │
└─────────────────────────────────────┘
                  ↓
┌─────────────────────────────────────┐
│    Argon2id Key Derivation          │
│    (64MB, 3 iterations, 4 threads)  │
└─────────────────────────────────────┘
                  ↓
┌─────────────────────────────────────┐
│    AES-256-GCM Encryption           │
│    (Authenticated Encryption)       │
└─────────────────────────────────────┘
                  ↓
┌─────────────────────────────────────┐
│    File Storage (Encrypted)         │
└─────────────────────────────────────┘
```

### Memory Safety Layer
```
┌─────────────────────────────────────┐
│    Sensitive Data Allocation        │
│    (Passwords, Keys, Secrets)       │
└─────────────────────────────────────┘
                  ↓
┌─────────────────────────────────────┐
│    Use for Operation                │
└─────────────────────────────────────┘
                  ↓
┌─────────────────────────────────────┐
│    Zero Memory (defer crypto.ZeroBytes)
└─────────────────────────────────────┘
```

## Concurrency Model

### Single-threaded by Design
gotp is designed to be single-threaded for simplicity and security:
- No goroutines for cryptographic operations
- Sequential command execution
- No race conditions possible

### Session Management
- Session tokens are stored in memory only
- No concurrent access to session data
- Automatic cleanup on exit

## Error Handling

### Error Types
- **User Errors**: Invalid input, missing data (displayed to user)
- **System Errors**: File I/O, encryption (logged with details)
- **Security Errors**: Authentication failures (generic messages)

### Error Propagation
```
CLI Command → Error → User Message (colored) → Exit Code
```

## Testing Strategy

### Unit Tests
- **Crypto**: Test vectors from RFCs
- **TOTP**: RFC 6238 test cases
- **Base32**: RFC 4648 test vectors
- **Importers**: Sample data from each format

### Integration Tests
- End-to-end vault operations
- CLI command testing
- File I/O operations

### Coverage Goals
- **Crypto**: 100% coverage
- **TOTP**: 100% coverage
- **Overall**: >80% coverage

## Performance Characteristics

### Time Complexity
- **TOTP Generation**: O(1) - constant time
- **Vault Encryption**: O(n) where n = JSON size
- **Vault Decryption**: O(n) where n = ciphertext size
- **Account Search**: O(n) where n = number of accounts

### Space Complexity
- **Vault in Memory**: O(n) where n = number of accounts
- **Session Cache**: O(1) - single token
- **Encryption Buffer**: O(1) - fixed size

### Benchmarks
- **Startup**: <100ms
- **TOTP Generation**: <1ms
- **Vault Load (1000 accounts)**: <500ms
- **Vault Save (1000 accounts)**: <500ms

## Platform Support

### Linux
- **Paths**: `~/.config/gotp/`
- **Clipboard**: xclip/wl-copy
- **Terminal**: Full color support

### macOS
- **Paths**: `~/Library/Application Support/gotp/`
- **Clipboard**: pbcopy/pbpaste
- **Terminal**: Full color support

### Windows
- **Paths**: `%APPDATA%\gotp\`
- **Clipboard**: PowerShell
- **Terminal**: ANSI color support

## Future Extensions

### Planned
- Hardware security module (HSM) support
- YubiKey integration
- System keyring integration
- WebAuthn support
- Cloud sync (optional, encrypted)

### Architecture Considerations
- Plugin system for storage backends
- Extension points for new authenticator formats
- API for programmatic access
- Web interface (optional)

## Dependencies

### Core Dependencies
- `github.com/spf13/cobra`: CLI framework
- `github.com/google/uuid`: UUID generation
- `golang.org/x/crypto`: Cryptographic functions
- `golang.org/x/term`: Terminal handling
- `gopkg.in/yaml.v3`: Configuration format

### QR Code Dependencies
- `github.com/makiuchi-d/gozxing`: QR code parsing
- `github.com/skip2/go-qrcode`: QR code generation

### No External Dependencies For
- TOTP/HOTP implementation (RFC-compliant)
- Base32 encoding/decoding
- Memory safety utilities
- Session management

## Code Quality

### Standards
- Follows Go best practices
- Idiomatic Go code
- Comprehensive error handling
- Extensive documentation
- 100% test coverage for crypto

### Linting
- `golangci-lint` for static analysis
- No warnings or errors
- Consistent code style

## Deployment

### Build Process
```bash
go build -o gotp ./cmd/gotp
```

### Distribution
- Single binary (no dependencies)
- Cross-platform builds
- Static linking for portability

### Installation
```bash
# From source
go install github.com/zulfikawr/gotp@latest

# Manual
chmod +x gotp
sudo mv gotp /usr/local/bin/
```

## Maintenance

### Versioning
- Semantic versioning (MAJOR.MINOR.PATCH)
- Breaking changes in MAJOR versions
- New features in MINOR versions
- Bug fixes in PATCH versions

### Release Process
1. Run full test suite
2. Update documentation
3. Create release tag
4. Build binaries for all platforms
5. Publish to GitHub Releases
6. Update package managers

### Support
- Security issues: security@gotp.dev
- Bug reports: GitHub Issues
- Feature requests: GitHub Discussions
