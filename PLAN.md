# Comprehensive LLM Prompt for Building `gotp` - A Terminal-Based TOTP Authenticator

---

## 📋 Project Specification Document

### Executive Summary

You are tasked with building `gotp`, a secure, cross-platform, terminal-based TOTP (Time-based One-Time Password) authenticator written in Go. This application will allow users to generate, store, and manage two-factor authentication codes entirely from their terminal environment.

---

## 1. Project Metadata

```yaml
Project Name: gotp
Language: Go (1.21+)
License: MIT
Target Platforms:
  - Linux (x86_64, arm64)
  - macOS (x86_64, arm64)
  - Windows (x86_64)
Architecture: Monolithic CLI application with optional TUI mode
```

---

## 2. Detailed Requirements Breakdown

### 2.1 Core TOTP Implementation (RFC 6238)

#### 2.1.1 Algorithm Requirements

```markdown
**MUST implement from scratch (zero external dependencies):**

1. **HMAC-SHA1 Implementation**
   - Input: Key (byte array), Message (byte array)
   - Output: 20-byte hash
   - Follow RFC 2104 specification
   - Support key lengths > 64 bytes (hash the key first)
   - Support key lengths < 64 bytes (pad with zeros)

2. **HOTP Algorithm (RFC 4226)**
   - Counter-based OTP generation
   - Dynamic truncation of HMAC result
   - Configurable digit length (default: 6, support 6-8)
   - Formula: HOTP(K, C) = Truncate(HMAC-SHA1(K, C))

3. **TOTP Algorithm (RFC 6238)**
   - Time-based counter derivation
   - Default time step: 30 seconds
   - Support configurable time steps (30, 60 seconds)
   - Unix epoch as T0 (default)
   - Formula: TOTP = HOTP(K, T) where T = (Current Unix Time - T0) / X

4. **Dynamic Truncation**
   - Extract 4-bit offset from last byte of HMAC
   - Extract 31-bit integer from offset position
   - Modulo 10^d where d = digit count
```

#### 2.1.2 Supported Algorithms

```go
// Must support these hash algorithms
type HashAlgorithm string

const (
    SHA1   HashAlgorithm = "SHA1"   // Default, most common
    SHA256 HashAlgorithm = "SHA256" // Enhanced security
    SHA512 HashAlgorithm = "SHA512" // Maximum security
)
```

#### 2.1.3 TOTP Generation Code Structure

```go
// Core interface to implement
type TOTPGenerator interface {
    // Generate current TOTP code
    Generate(secret []byte, timestamp time.Time) (string, error)
    
    // Generate with specific parameters
    GenerateWithParams(params TOTPParams) (string, error)
    
    // Validate a provided code (with time window tolerance)
    Validate(secret []byte, code string, timestamp time.Time, window int) (bool, error)
    
    // Get remaining seconds until next code
    RemainingSeconds(timestamp time.Time) int
}

type TOTPParams struct {
    Secret    []byte
    Timestamp time.Time
    Period    int           // 30 or 60 seconds
    Digits    int           // 6, 7, or 8
    Algorithm HashAlgorithm // SHA1, SHA256, SHA512
}
```

---

### 2.2 Encrypted Storage System

#### 2.2.1 Encryption Requirements

```markdown
**Security Specifications:**

1. **Master Password System**
   - Minimum 8 characters
   - Argon2id for key derivation (recommended parameters):
     - Memory: 64MB
     - Iterations: 3
     - Parallelism: 4
     - Salt: 16 bytes (random)
     - Key length: 32 bytes

2. **Data Encryption**
   - Algorithm: AES-256-GCM
   - Nonce: 12 bytes (random, unique per encryption)
   - Store nonce with ciphertext
   - Authenticate additional data (AAD) for metadata

3. **Storage Format**
   - JSON structure encrypted as single blob
   - Version field for future migrations
   - Integrity verification on load
```

#### 2.2.2 Storage Schema

```json
{
  "version": "1.0",
  "created_at": "2024-01-15T10:30:00Z",
  "modified_at": "2024-01-15T14:22:00Z",
  "kdf_params": {
    "algorithm": "argon2id",
    "memory": 65536,
    "iterations": 3,
    "parallelism": 4,
    "salt": "base64_encoded_salt"
  },
  "accounts": [
    {
      "id": "uuid-v4",
      "name": "GitHub",
      "issuer": "GitHub",
      "username": "user@example.com",
      "secret": "encrypted_base32_secret",
      "algorithm": "SHA1",
      "digits": 6,
      "period": 30,
      "created_at": "2024-01-15T10:30:00Z",
      "last_used": "2024-01-15T14:22:00Z",
      "tags": ["work", "development"],
      "icon": "github",
      "sort_order": 0
    }
  ]
}
```

#### 2.2.3 File Locations

```markdown
**Platform-Specific Paths:**

- **Linux**: `~/.config/gotp/vault.enc`
- **macOS**: `~/Library/Application Support/gotp/vault.enc`
- **Windows**: `%APPDATA%\gotp\vault.enc`

**Backup Location:**
- Same directory with `.bak` extension
- Keep last 3 backups with timestamps
```

---

### 2.3 CLI Interface Specification

#### 2.3.1 Command Structure

```bash
gotp [global flags] <command> [command flags] [arguments]
```

#### 2.3.2 Global Flags

```markdown
| Flag | Short | Description | Default |
|------|-------|-------------|---------|
| --vault | -v | Path to vault file | Platform default |
| --config | -c | Path to config file | Platform default |
| --no-color | | Disable colored output | false |
| --json | -j | Output in JSON format | false |
| --quiet | -q | Minimal output | false |
| --verbose | | Verbose output | false |
| --version | -V | Show version | - |
| --help | -h | Show help | - |
```

#### 2.3.3 Commands Specification

```markdown
### `gotp init`
Initialize a new vault with master password.

**Flags:**
- `--force, -f`: Overwrite existing vault

**Behavior:**
1. Check if vault exists (error if exists without --force)
2. Prompt for master password (twice for confirmation)
3. Validate password strength
4. Generate salt and derive key
5. Create empty encrypted vault
6. Display success message

**Example:**
```bash
$ gotp init
Creating new vault...
Enter master password: ********
Confirm master password: ********
✓ Vault created successfully at ~/.config/gotp/vault.enc
```

---

### `gotp add <name>`
Add a new TOTP account.

**Flags:**
- `--secret, -s <secret>`: Base32 secret (prompted if not provided)
- `--issuer, -i <issuer>`: Service issuer name
- `--username, -u <username>`: Account username/email
- `--algorithm, -a <algo>`: Hash algorithm (SHA1|SHA256|SHA512)
- `--digits, -d <num>`: Code digits (6|7|8)
- `--period, -p <seconds>`: Time period (30|60)
- `--tags, -t <tags>`: Comma-separated tags
- `--uri`: Add from otpauth:// URI
- `--qr <file>`: Add from QR code image file
- `--scan`: Scan QR from screen (if supported)

**Behavior:**
1. Prompt for master password
2. Decrypt and load vault
3. Validate secret (base32 format)
4. Parse otpauth:// URI if provided
5. Generate UUID for account
6. Add account to vault
7. Re-encrypt and save vault
8. Display current code as confirmation

**Example:**
```bash
$ gotp add GitHub --secret JBSWY3DPEHPK3PXP --username user@example.com
Enter master password: ********
✓ Added account: GitHub (user@example.com)
Current code: 123456 (expires in 18s)
```

---

### `gotp get <name|pattern>`
Get TOTP code for an account.

**Flags:**
- `--copy, -c`: Copy to clipboard
- `--continuous, -w`: Watch mode, update every second
- `--qr`: Display as QR code in terminal

**Behavior:**
1. Prompt for master password (or use cached session)
2. Search for account by name (fuzzy match if multiple)
3. Generate current TOTP
4. Display code with remaining time
5. Optionally copy to clipboard

**Example:**
```bash
$ gotp get github
Enter master password: ********
GitHub (user@example.com)
┌────────────┐
│   847293   │ ████████░░ 18s
└────────────┘

$ gotp get github -c
✓ Code copied to clipboard (expires in 18s)
```

---

### `gotp list`
List all stored accounts.

**Flags:**
- `--with-codes`: Show current codes
- `--filter, -f <tag>`: Filter by tag
- `--sort <field>`: Sort by name|issuer|last_used|created

**Behavior:**
1. Prompt for master password
2. List accounts in table format
3. Show codes if requested

**Example:**
```bash
$ gotp list
Enter master password: ********
┌──────────────┬─────────────────────┬──────────┬────────────┐
│ Name         │ Username            │ Issuer   │ Tags       │
├──────────────┼─────────────────────┼──────────┼────────────┤
│ GitHub       │ user@example.com    │ GitHub   │ work, dev  │
│ AWS Console  │ admin@company.com   │ Amazon   │ work       │
│ Personal Gmail│ me@gmail.com       │ Google   │ personal   │
└──────────────┴─────────────────────┴──────────┴────────────┘
Total: 3 accounts
```

---

### `gotp remove <name>`
Remove an account from the vault.

**Flags:**
- `--force, -f`: Skip confirmation

**Behavior:**
1. Prompt for master password
2. Find account
3. Confirm deletion (unless --force)
4. Remove from vault
5. Save encrypted vault

**Example:**
```bash
$ gotp remove GitHub
Enter master password: ********
⚠ Are you sure you want to remove "GitHub"? This cannot be undone. [y/N]: y
✓ Removed account: GitHub
```

---

### `gotp edit <name>`
Edit an existing account.

**Flags:**
- `--name <new_name>`: Change account name
- `--username <username>`: Change username
- `--issuer <issuer>`: Change issuer
- `--secret <secret>`: Change secret (requires confirmation)
- `--tags <tags>`: Replace tags
- `--add-tag <tag>`: Add a tag
- `--remove-tag <tag>`: Remove a tag

**Example:**
```bash
$ gotp edit GitHub --username new@example.com --add-tag cloud
Enter master password: ********
✓ Updated account: GitHub
```

---

### `gotp export`
Export accounts for backup.

**Flags:**
- `--format <format>`: Export format (json|encrypted|uri)
- `--output, -o <file>`: Output file (stdout if not specified)
- `--accounts <names>`: Specific accounts (comma-separated)

**Behavior:**
1. Prompt for master password
2. Decrypt vault
3. For encrypted export: prompt for export password
4. Export in specified format
5. Warn about security implications

**Example:**
```bash
$ gotp export --format json -o backup.json
Enter master password: ********
⚠ WARNING: This will export secrets in PLAIN TEXT.
Continue? [y/N]: y
✓ Exported 3 accounts to backup.json
```

---

### `gotp import <file>`
Import accounts from backup or other authenticators.

**Flags:**
- `--format <format>`: Import format (json|encrypted|uri|aegis|authy|google)
- `--merge`: Merge with existing (skip duplicates)
- `--overwrite`: Overwrite duplicates

**Example:**
```bash
$ gotp import backup.json --merge
Enter master password: ********
Found 5 accounts to import.
  ✓ GitHub (new)
  ✓ AWS (new)
  ⏭ Google (skipped - duplicate)
  ✓ Dropbox (new)
  ✓ Twitter (new)
Imported: 4, Skipped: 1
```

---

### `gotp passwd`
Change master password.

**Behavior:**
1. Prompt for current password
2. Decrypt and validate vault
3. Prompt for new password (twice)
4. Re-derive key with new password
5. Re-encrypt vault

**Example:**
```bash
$ gotp passwd
Enter current password: ********
Enter new password: ********
Confirm new password: ********
✓ Master password changed successfully
```

---

### `gotp tui`
Launch interactive Terminal User Interface.

**Flags:**
- `--theme <theme>`: Color theme (dark|light|nord|dracula)

---

### `gotp completion <shell>`
Generate shell completion scripts.

**Supported shells:** bash, zsh, fish, powershell
```

---

### 2.4 TUI Interface Specification

#### 2.4.1 Layout Design

```
┌─────────────────────────────────────────────────────────────────────────┐
│  gotp v1.0.0                                        Press ? for help    │
├────────────────────────┬────────────────────────────────────────────────┤
│  Accounts              │  GitHub                                        │
│  ─────────────────────-│  ────────────────────────────────────────────  │
│  > GitHub         ★    │                                                │
│    AWS Console         │     ┌──────────────────────────────────┐      │
│    Personal Gmail      │     │                                  │      │
│    Dropbox             │     │            482 918               │      │
│    Twitter             │     │                                  │      │
│                        │     └──────────────────────────────────┘      │
│                        │                                                │
│                        │     ████████████████░░░░░░  18s remaining     │
│                        │                                                │
│                        │  Username: user@example.com                   │
│                        │  Issuer:   GitHub                             │
│                        │  Tags:     work, development                  │
│                        │                                                │
│                        │  [c] Copy  [e] Edit  [d] Delete  [Enter] Copy │
├────────────────────────┴────────────────────────────────────────────────┤
│  Filter: _                            3/5 accounts  │ ↑↓ navigate │ q quit│
└─────────────────────────────────────────────────────────────────────────┘
```

#### 2.4.2 TUI Features

```markdown
**Navigation:**
- Arrow keys / vim keys (j/k) for account navigation
- Tab to switch between panels
- / to activate filter/search
- Enter to copy current code
- ESC to cancel/go back

**Views:**
1. **Main View** - Account list with live codes
2. **Add Account View** - Form for new account
3. **Edit Account View** - Modify existing account
4. **Settings View** - App configuration
5. **Help View** - Keyboard shortcuts reference

**Features:**
- Live-updating codes with countdown timer
- Fuzzy search/filter
- Favorites/pinned accounts (shown first)
- Color-coded countdown (green > yellow > red)
- Clipboard integration with notification
- QR code display for account sharing
```

#### 2.4.3 TUI Components to Build

```markdown
1. **AccountList Component**
   - Scrollable list with selection highlight
   - Shows: name, issuer, current code (optional), favorite indicator
   - Supports filtering

2. **CodeDisplay Component**
   - Large code display
   - Animated progress bar
   - Color transitions based on time remaining

3. **AccountForm Component**
   - Input fields for all account properties
   - Validation feedback
   - Secret input (hidden/revealed toggle)

4. **SearchBar Component**
   - Fuzzy matching
   - Real-time filtering
   - Match highlighting

5. **StatusBar Component**
   - Current filter
   - Account count
   - Keyboard hints
   - Notifications

6. **Dialog Component**
   - Confirmation dialogs
   - Error messages
   - Password prompts
```

---

### 2.5 Configuration System

#### 2.5.1 Configuration File

```yaml
# ~/.config/gotp/config.yaml

# General settings
general:
  default_digits: 6
  default_period: 30
  default_algorithm: SHA1
  auto_copy: false
  clear_clipboard_after: 30  # seconds, 0 to disable
  session_timeout: 300       # seconds before requiring password again

# CLI settings  
cli:
  color: true
  json_output: false
  date_format: "2006-01-02 15:04:05"

# TUI settings
tui:
  theme: dark
  show_codes_in_list: true
  confirm_delete: true
  animate_progress: true
  refresh_rate: 100  # milliseconds

# Security settings
security:
  argon2_memory: 65536    # KB
  argon2_iterations: 3
  argon2_parallelism: 4
  backup_count: 3
  auto_lock: true
  auto_lock_timeout: 300  # seconds

# Paths (override defaults)
paths:
  vault: ""   # empty = platform default
  backup: ""  # empty = same as vault directory
```

---

### 2.6 Security Requirements

#### 2.6.1 Security Checklist

```markdown
**Memory Security:**
- [ ] Zero secrets from memory after use
- [ ] Use secure memory allocation where available
- [ ] Avoid string operations on secrets (use byte slices)
- [ ] Clear clipboard after configurable timeout

**Storage Security:**
- [ ] Never store plaintext secrets
- [ ] Use authenticated encryption (AES-GCM)
- [ ] Derive keys with Argon2id
- [ ] Random salt per vault
- [ ] Random nonce per encryption operation

**Input Validation:**
- [ ] Validate base32 secret format
- [ ] Sanitize account names
- [ ] Validate otpauth:// URI format
- [ ] Limit input lengths

**Session Management:**
- [ ] Optional session caching with timeout
- [ ] Secure session token storage
- [ ] Clear session on explicit lock

**Audit & Logging:**
- [ ] No secrets in logs
- [ ] Log authentication attempts
- [ ] Optional audit trail for account changes
```

#### 2.6.2 Threat Model

```markdown
**Protected Against:**
- Unauthorized vault access (encryption)
- Brute force attacks (Argon2id)
- Memory dumps (zeroing)
- Clipboard sniffing (auto-clear)
- Shoulder surfing (optional code hiding)

**Not Protected Against:**
- Keyloggers (master password capture)
- Compromised system (full access)
- Physical access while unlocked
- Screen capture while codes visible
```

---

## 3. Project Structure

```
gotp/
├── cmd/
│   └── gotp/
│       └── main.go                 # Entry point
├── internal/
│   ├── totp/
│   │   ├── hotp.go                 # HOTP implementation
│   │   ├── totp.go                 # TOTP implementation
│   │   ├── hmac.go                 # HMAC-SHA1/256/512
│   │   └── totp_test.go            # RFC test vectors
│   ├── crypto/
│   │   ├── argon2.go               # Key derivation
│   │   ├── aes.go                  # AES-GCM encryption
│   │   ├── secure.go               # Secure memory operations
│   │   └── crypto_test.go
│   ├── vault/
│   │   ├── vault.go                # Vault management
│   │   ├── account.go              # Account struct & methods
│   │   ├── storage.go              # File I/O
│   │   ├── backup.go               # Backup management
│   │   └── vault_test.go
│   ├── cli/
│   │   ├── cli.go                  # CLI setup
│   │   ├── commands/
│   │   │   ├── init.go
│   │   │   ├── add.go
│   │   │   ├── get.go
│   │   │   ├── list.go
│   │   │   ├── remove.go
│   │   │   ├── edit.go
│   │   │   ├── export.go
│   │   │   ├── import.go
│   │   │   ├── passwd.go
│   │   │   └── completion.go
│   │   └── ui/
│   │       ├── prompt.go           # Password prompts
│   │       ├── table.go            # Table output
│   │       └── progress.go         # Progress display
│   ├── tui/
│   │   ├── app.go                  # TUI application
│   │   ├── model.go                # Application state
│   │   ├── update.go               # State updates
│   │   ├── view.go                 # Rendering
│   │   ├── components/
│   │   │   ├── list.go
│   │   │   ├── code.go
│   │   │   ├── form.go
│   │   │   ├── search.go
│   │   │   └── dialog.go
│   │   ├── styles/
│   │   │   └── themes.go
│   │   └── keys.go                 # Key bindings
│   ├── config/
│   │   ├── config.go               # Configuration loading
│   │   └── paths.go                # Platform paths
│   ├── qr/
│   │   ├── parse.go                # QR code parsing
│   │   ├── generate.go             # QR code generation
│   │   └── terminal.go             # Terminal QR display
│   └── clipboard/
│       ├── clipboard.go            # Clipboard interface
│       ├── clipboard_linux.go
│       ├── clipboard_darwin.go
│       └── clipboard_windows.go
├── pkg/
│   └── base32/
│       └── base32.go               # Custom base32 for secrets
├── testdata/
│   ├── test_vault.enc
│   └── test_qr.png
├── docs/
│   ├── SECURITY.md
│   ├── ARCHITECTURE.md
│   └── CONTRIBUTING.md
├── scripts/
│   ├── build.sh
│   └── release.sh
├── .github/
│   └── workflows/
│       ├── ci.yml
│       └── release.yml
├── go.mod
├── go.sum
├── Makefile
├── README.md
└── LICENSE
```

---

## 4. Implementation Phases

### Phase 1: Core TOTP Engine (Week 1)

```markdown
**Deliverables:**
1. HMAC-SHA1/256/512 implementation (no external deps)
2. HOTP algorithm with test vectors
3. TOTP algorithm with test vectors
4. Base32 encoding/decoding
5. 100% test coverage for crypto primitives

**Test Vectors (RFC 6238):**
| Time (seconds) | SHA1     | SHA256   | SHA512   |
|----------------|----------|----------|----------|
| 59             | 94287082 | 46119246 | 90693936 |
| 1111111109     | 07081804 | 68084774 | 25091201 |
| 1111111111     | 14050471 | 67062674 | 99943326 |
| 1234567890     | 89005924 | 91819424 | 93441116 |
| 2000000000     | 69279037 | 90698825 | 38618901 |
```

### Phase 2: Encryption & Storage (Week 2)

```markdown
**Deliverables:**
1. Argon2id key derivation
2. AES-256-GCM encryption/decryption
3. Vault struct and serialization
4. Secure file I/O
5. Backup system
6. Memory zeroing utilities
```

### Phase 3: CLI Implementation (Week 3)

```markdown
**Deliverables:**
1. Command parser setup
2. All commands implemented
3. Password prompting (secure, no echo)
4. Table/formatted output
5. JSON output mode
6. Shell completions
```

### Phase 4: TUI Implementation (Week 4)

```markdown
**Deliverables:**
1. Main application loop
2. Account list view
3. Code display with live updates
4. Add/edit forms
5. Search/filter functionality
6. Themes and styling
```

### Phase 5: Polish & Release (Week 5)

```markdown
**Deliverables:**
1. QR code support (parse and display)
2. Import from other authenticators
3. Documentation
4. CI/CD pipeline
5. Release binaries
6. Installation scripts
```

---

## 5. Testing Requirements

### 5.1 Unit Tests

```go
// Example test structure for TOTP
func TestTOTP_RFC6238_SHA1(t *testing.T) {
    testCases := []struct {
        time     int64
        expected string
    }{
        {59, "94287082"},
        {1111111109, "07081804"},
        {1111111111, "14050471"},
        {1234567890, "89005924"},
        {2000000000, "69279037"},
    }
    
    secret := []byte("12345678901234567890")
    
    for _, tc := range testCases {
        t.Run(fmt.Sprintf("time_%d", tc.time), func(t *testing.T) {
            result := GenerateTOTP(secret, time.Unix(tc.time, 0), 8, SHA1)
            if result != tc.expected {
                t.Errorf("expected %s, got %s", tc.expected, result)
            }
        })
    }
}
```

### 5.2 Integration Tests

```markdown
- Vault creation and loading
- Add/remove/edit account workflows
- Export/import roundtrip
- Password change
- Session management
```

### 5.3 Security Tests

```markdown
- Memory scanning for secrets after operations
- Encryption strength validation
- Invalid input handling
- Timing attack resistance (constant-time comparison)
```

---

## 6. External Dependencies Policy

### 6.1 Allowed Dependencies

```markdown
**Core (Zero External):**
- TOTP/HOTP algorithms: Must be implemented from scratch
- HMAC: Must be implemented from scratch
- Base32: Must be implemented from scratch

**Allowed External:**
- `golang.org/x/crypto/argon2` - Argon2id (too complex to implement safely)
- `golang.org/x/term` - Terminal handling
- `github.com/charmbracelet/bubbletea` - TUI framework (optional)
- `github.com/charmbracelet/lipgloss` - TUI styling (optional)

**Explicitly Forbidden:**
- Any TOTP/OTP library
- Any HMAC library (use crypto/hmac from stdlib only for reference)
- Any base32 library
```

---

## 7. Build & Release

### 7.1 Makefile

```makefile
.PHONY: build test clean release

VERSION := $(shell git describe --tags --always --dirty)
LDFLAGS := -ldflags "-X main.version=$(VERSION)"

build:
	go build $(LDFLAGS) -o bin/gotp ./cmd/gotp

test:
	go test -v -race -cover ./...

test-coverage:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out

lint:
	golangci-lint run

clean:
	rm -rf bin/ coverage.out

release:
	goreleaser release --clean

install:
	go install $(LDFLAGS) ./cmd/gotp
```

### 7.2 Release Artifacts

```markdown
- gotp_linux_amd64.tar.gz
- gotp_linux_arm64.tar.gz
- gotp_darwin_amd64.tar.gz
- gotp_darwin_arm64.tar.gz
- gotp_windows_amd64.zip
- SHA256SUMS
- SHA256SUMS.sig (GPG signed)
```

---

## 8. Documentation Requirements

### 8.1 README.md Structure

```markdown
# gotp

Badges: Build, Coverage, Release, License

## Features
## Installation
## Quick Start
## Usage
### CLI Commands
### TUI Mode
## Security
## Configuration
## Contributing
## License
```

### 8.2 Man Pages

```markdown
Generate man pages for:
- gotp(1) - Main command
- gotp-add(1) - Add command
- gotp-get(1) - Get command
- etc.
```

---

## 9. Success Criteria

```markdown
**Functional:**
- [ ] Generate valid TOTP codes matching reference implementations
- [ ] All CLI commands work as specified
- [ ] TUI launches and is fully functional
- [ ] Import/export works with common formats
- [ ] Cross-platform builds successful

**Security:**
- [ ] Passes security review checklist
- [ ] No plaintext secrets in logs or errors
- [ ] Encryption implementation auditable

**Quality:**
- [ ] >90% test coverage on crypto code
- [ ] >80% test coverage overall
- [ ] Zero linter warnings
- [ ] Documentation complete

**Performance:**
- [ ] Startup time <100ms
- [ ] Code generation <1ms
- [ ] Vault loading <500ms for 1000 accounts
```