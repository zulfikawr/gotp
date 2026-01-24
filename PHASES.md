# gotp Development Phases - Detailed Step-by-Step

---

## Phase 1: Core TOTP Engine Implementation

### 1.1 Setup & Project Initialization
- [x] Initialize Go module (`go mod init github.com/zulfikawr/gotp`)
- [x] Create project directory structure
- [x] Set up `cmd/gotp/main.go` (entry point)
- [x] Create `internal/` package structure
- [x] Set up `.gitignore` and git repository
- [x] Create `go.mod` and `go.sum`

### 1.2 HMAC-SHA1/SHA256/SHA512 Implementation
- [x] Implement `internal/totp/hmac.go` with zero external dependencies
- [x] Support configurable hash algorithms (SHA1, SHA256, SHA512)
- [x] Handle key padding (< 64 bytes pad with zeros, > 64 bytes hash first)
- [x] Implement RFC 2104 HMAC specification
- [x] Create unit tests for HMAC with test vectors
- [x] Add test coverage for all three hash algorithms
- [x] Achieve 100% code coverage for HMAC implementation

### 1.3 Base32 Encoding/Decoding
- [x] Implement `pkg/base32/base32.go` (custom, no external deps)
- [x] Support standard Base32 alphabet (RFC 4648)
- [x] Implement encoder for secret generation
- [x] Implement decoder for secret validation
- [x] Handle padding and edge cases
- [x] Add comprehensive unit tests
- [x] Validate against RFC test vectors

### 1.4 HOTP Algorithm Implementation
- [x] Implement `internal/totp/hotp.go` (RFC 4226)
- [x] Create function for counter-based OTP generation
- [x] Implement dynamic truncation of HMAC result
- [x] Support configurable digit length (6, 7, 8 digits)
- [x] Extract 4-bit offset from last HMAC byte
- [x] Extract 31-bit integer and apply modulo 10^d
- [x] Add unit tests with RFC test vectors
- [x] Test all supported digit configurations

### 1.5 TOTP Algorithm Implementation
- [x] Implement `internal/totp/totp.go` (RFC 6238)
- [x] Implement time-based counter derivation from Unix timestamp
- [x] Support configurable time steps (30, 60 seconds)
- [x] Implement TOTP generation function using HOTP
- [x] Add time validation and remaining seconds calculation
- [x] Create test cases with RFC 6238 test vectors
- [x] Validate against reference implementations

### 1.6 TOTP Validation & Utilities
- [x] Implement code validation with time window tolerance
- [x] Create `RemainingSeconds()` function
- [x] Implement `ValidateWithWindow()` function
- [x] Add fuzzy matching for time drift
- [x] Create helper functions for timestamp conversion
- [x] Add comprehensive tests for validation

### 1.7 Comprehensive Testing
- [x] Run all RFC 6238 test vectors (5 standard test cases)
- [x] Test with SHA1, SHA256, SHA512 algorithms
- [x] Achieve 100% test coverage on crypto code
- [x] Implement security tests (constant-time comparison)
- [x] Add edge case tests (boundary conditions, padding)
- [x] Document test vectors and expected results

### 1.8 Documentation & Validation
- [x] Document HMAC implementation with RFC reference
- [x] Document HOTP algorithm with examples
- [x] Document TOTP algorithm with time handling
- [x] Add code comments explaining cryptographic operations
- [x] Create example usage documentation
- [x] Validate against online TOTP generators (e.g., Google Authenticator)

---

## Phase 2: Encryption & Storage System

### 2.1 Secure Memory Operations
- [ ] Implement `internal/crypto/secure.go`
- [ ] Create `ZeroBytes()` function for memory wiping
- [ ] Implement secure byte slice creation
- [ ] Add helpers to clear sensitive data after use
- [ ] Create tests verifying memory is actually zeroed
- [ ] Document memory safety practices

### 2.2 Argon2id Key Derivation
- [ ] Implement `internal/crypto/argon2.go`
- [ ] Use `golang.org/x/crypto/argon2` package
- [ ] Configure recommended parameters:
  - [ ] Memory: 65536 KB (64MB)
  - [ ] Iterations: 3
  - [ ] Parallelism: 4
  - [ ] Salt length: 16 bytes (random)
  - [ ] Output key length: 32 bytes
- [ ] Implement `DeriveKey(password, salt)` function
- [ ] Add secure random salt generation
- [ ] Create tests with known test vectors
- [ ] Implement parameter validation

### 2.3 AES-256-GCM Encryption
- [ ] Implement `internal/crypto/aes.go`
- [ ] Use `crypto/aes` and `crypto/cipher` from stdlib
- [ ] Implement 256-bit key encryption
- [ ] Generate random 12-byte nonces
- [ ] Implement authenticated encryption with GCM mode
- [ ] Add functions for encryption and decryption
- [ ] Handle nonce storage with ciphertext
- [ ] Create comprehensive encryption tests
- [ ] Test with various payload sizes
- [ ] Verify authentication failure handling

### 2.4 Account Data Structure
- [ ] Implement `internal/vault/account.go`
- [ ] Define Account struct with fields:
  - [ ] ID (UUID v4)
  - [ ] Name
  - [ ] Issuer
  - [ ] Username/Email
  - [ ] Secret (encrypted)
  - [ ] Algorithm (SHA1, SHA256, SHA512)
  - [ ] Digits (6, 7, 8)
  - [ ] Period (30, 60 seconds)
  - [ ] Tags (array of strings)
  - [ ] Icon name
  - [ ] Sort order
  - [ ] Created timestamp
  - [ ] Last used timestamp
- [ ] Implement JSON marshaling/unmarshaling
- [ ] Add validation methods
- [ ] Create helper functions for account operations

### 2.5 Vault Management
- [ ] Implement `internal/vault/vault.go`
- [ ] Create Vault struct containing:
  - [ ] Version string
  - [ ] Encrypted accounts array
  - [ ] KDF parameters
  - [ ] Created/modified timestamps
- [ ] Implement vault creation
- [ ] Implement vault unlock/decrypt
- [ ] Implement account CRUD operations
- [ ] Add in-memory vault locking
- [ ] Create tests for all vault operations

### 2.6 Storage & File I/O
- [ ] Implement `internal/vault/storage.go`
- [ ] Create platform-specific path resolution:
  - [ ] Linux: `~/.config/gotp/vault.enc`
  - [ ] macOS: `~/Library/Application Support/gotp/vault.enc`
  - [ ] Windows: `%APPDATA%\gotp\vault.enc`
- [ ] Implement file save with atomic writes
- [ ] Implement file load with integrity verification
- [ ] Add file permission handling (600 on Unix)
- [ ] Create tests for file operations

### 2.7 Backup System
- [ ] Implement `internal/vault/backup.go`
- [ ] Create backup on vault save
- [ ] Keep last 3 backups with timestamps
- [ ] Implement backup restoration
- [ ] Add backup path management
- [ ] Create backup tests
- [ ] Document backup directory structure

### 2.8 Configuration System
- [ ] Implement `internal/config/config.go`
- [ ] Create config struct with all settings
- [ ] Implement YAML parsing for config.yaml
- [ ] Set default values for all settings
- [ ] Implement per-platform config paths
- [ ] Add config validation
- [ ] Create tests for config loading

### 2.9 Platform Paths Resolution
- [ ] Implement `internal/config/paths.go`
- [ ] Create platform detection
- [ ] Implement vault path resolution
- [ ] Implement config path resolution
- [ ] Implement backup path resolution
- [ ] Add directory creation helpers
- [ ] Test all platform paths

### 2.10 Encryption Testing & Validation
- [ ] Test encryption roundtrip (encrypt then decrypt)
- [ ] Test with large payloads (1000+ accounts)
- [ ] Verify authentication failure detection
- [ ] Test nonce uniqueness
- [ ] Test key derivation determinism
- [ ] Benchmark encryption/decryption speed
- [ ] Verify no plaintext leakage in storage

---

## Phase 3: CLI Interface Implementation

### 3.1 CLI Framework Setup
- [ ] Implement `internal/cli/cli.go`
- [ ] Choose CLI framework (flag package or external)
- [ ] Set up command routing
- [ ] Implement global flags:
  - [ ] `--vault` / `-v`
  - [ ] `--config` / `-c`
  - [ ] `--no-color`
  - [ ] `--json` / `-j`
  - [ ] `--quiet` / `-q`
  - [ ] `--verbose`
  - [ ] `--version` / `-V`
  - [ ] `--help` / `-h`
- [ ] Create help text system
- [ ] Implement version display

### 3.2 Password Prompting
- [ ] Implement `internal/cli/ui/prompt.go`
- [ ] Create secure password input (no echo)
- [ ] Use `golang.org/x/term` for terminal handling
- [ ] Implement password confirmation prompts
- [ ] Add password strength validation
- [ ] Create password retry logic
- [ ] Test on all platforms

### 3.3 Table & Formatted Output
- [ ] Implement `internal/cli/ui/table.go`
- [ ] Create table formatting functions
- [ ] Implement column alignment
- [ ] Add color support (optional)
- [ ] Create success/error message helpers
- [ ] Implement JSON output formatting

### 3.4 Progress & Countdown Display
- [ ] Implement `internal/cli/ui/progress.go`
- [ ] Create progress bar display
- [ ] Implement countdown timer
- [ ] Add color transitions (green → yellow → red)
- [ ] Create animated progress updates
- [ ] Test update frequency

### 3.5 `gotp init` Command
- [ ] Implement `internal/cli/commands/init.go`
- [ ] Check for existing vault
- [ ] Prompt for master password (twice)
- [ ] Validate password strength
- [ ] Generate random salt
- [ ] Derive encryption key
- [ ] Create empty encrypted vault
- [ ] Display success message
- [ ] Test all scenarios (new, force overwrite, etc.)

### 3.6 `gotp add` Command
- [ ] Implement `internal/cli/commands/add.go`
- [ ] Parse command flags:
  - [ ] `--secret` / `-s`
  - [ ] `--issuer` / `-i`
  - [ ] `--username` / `-u`
  - [ ] `--algorithm` / `-a`
  - [ ] `--digits` / `-d`
  - [ ] `--period` / `-p`
  - [ ] `--tags` / `-t`
  - [ ] `--uri`
  - [ ] `--qr`
  - [ ] `--scan`
- [ ] Prompt for master password
- [ ] Decrypt vault
- [ ] Validate secret (base32 format)
- [ ] Parse otpauth:// URI if provided
- [ ] Generate UUID for account
- [ ] Add account to vault
- [ ] Re-encrypt and save vault
- [ ] Display current code confirmation
- [ ] Test all input paths

### 3.7 `gotp get` Command
- [ ] Implement `internal/cli/commands/get.go`
- [ ] Parse command flags:
  - [ ] `--copy` / `-c`
  - [ ] `--continuous` / `-w`
  - [ ] `--qr`
- [ ] Prompt for master password
- [ ] Implement session caching
- [ ] Search for account (fuzzy match if multiple)
- [ ] Generate current TOTP code
- [ ] Display code with remaining time
- [ ] Implement copy to clipboard
- [ ] Implement watch mode (continuous update)
- [ ] Test all modes

### 3.8 `gotp list` Command
- [ ] Implement `internal/cli/commands/list.go`
- [ ] Parse command flags:
  - [ ] `--with-codes`
  - [ ] `--filter` / `-f`
  - [ ] `--sort`
- [ ] Prompt for master password
- [ ] List accounts in table format
- [ ] Implement filtering by tag
- [ ] Implement sorting options
- [ ] Show current codes if requested
- [ ] Display total count
- [ ] Test all filter/sort combinations

### 3.9 `gotp remove` Command
- [ ] Implement `internal/cli/commands/remove.go`
- [ ] Parse command flags:
  - [ ] `--force` / `-f`
- [ ] Prompt for master password
- [ ] Find account by name
- [ ] Show confirmation prompt (unless --force)
- [ ] Remove from vault
- [ ] Re-encrypt and save vault
- [ ] Display success message
- [ ] Test confirmation logic

### 3.10 `gotp edit` Command
- [ ] Implement `internal/cli/commands/edit.go`
- [ ] Parse command flags for all editable fields
- [ ] Prompt for master password
- [ ] Load account from vault
- [ ] Apply field changes
- [ ] For secret change: require confirmation
- [ ] Re-encrypt and save vault
- [ ] Display updated account info
- [ ] Test all field updates

### 3.11 `gotp export` Command
- [ ] Implement `internal/cli/commands/export.go`
- [ ] Parse command flags:
  - [ ] `--format`
  - [ ] `--output` / `-o`
  - [ ] `--accounts`
- [ ] Prompt for master password
- [ ] Support export formats:
  - [ ] JSON (plaintext)
  - [ ] Encrypted (password-protected)
  - [ ] otpauth:// URIs
- [ ] Show security warning for plaintext
- [ ] Write to file or stdout
- [ ] Test all export formats

### 3.12 `gotp import` Command
- [ ] Implement `internal/cli/commands/import.go`
- [ ] Parse command flags:
  - [ ] `--format`
  - [ ] `--merge`
  - [ ] `--overwrite`
- [ ] Prompt for master password
- [ ] Support import formats:
  - [ ] JSON
  - [ ] Encrypted
  - [ ] otpauth:// URIs
  - [ ] Aegis backup
  - [ ] Authy backup
  - [ ] Google Authenticator
- [ ] Handle duplicate accounts
- [ ] Display import summary
- [ ] Test all format imports

### 3.13 `gotp passwd` Command
- [ ] Implement `internal/cli/commands/passwd.go`
- [ ] Prompt for current password
- [ ] Verify vault decryption
- [ ] Prompt for new password (twice)
- [ ] Re-derive key with new password
- [ ] Re-encrypt vault
- [ ] Save updated vault
- [ ] Display success message
- [ ] Test password change flow

### 3.14 `gotp completion` Command
- [ ] Implement `internal/cli/commands/completion.go`
- [ ] Support shell types:
  - [ ] bash
  - [ ] zsh
  - [ ] fish
  - [ ] powershell
- [ ] Generate completion script
- [ ] Output to stdout
- [ ] Test completion generation

### 3.15 Main Application Entry
- [ ] Implement `cmd/gotp/main.go`
- [ ] Set version string
- [ ] Initialize CLI framework
- [ ] Register all commands
- [ ] Handle global flags
- [ ] Implement error handling
- [ ] Exit with appropriate status codes
- [ ] Test application startup

### 3.16 Clipboard Integration
- [ ] Implement `internal/clipboard/clipboard.go`
- [ ] Create platform-specific implementations:
  - [ ] `internal/clipboard/clipboard_linux.go` (xclip/wl-copy)
  - [ ] `internal/clipboard/clipboard_darwin.go` (pbcopy)
  - [ ] `internal/clipboard/clipboard_windows.go` (PowerShell)
- [ ] Implement code copy functionality
- [ ] Implement clipboard clear after timeout
- [ ] Test on all platforms

### 3.17 Session Management
- [ ] Implement session caching
- [ ] Create session timeout logic
- [ ] Store session token securely
- [ ] Implement session locking
- [ ] Add session expiration
- [ ] Test session workflow

### 3.18 CLI Testing & Integration
- [ ] Create integration tests for all commands
- [ ] Test error conditions
- [ ] Test input validation
- [ ] Test with various account configurations
- [ ] Test color/non-color modes
- [ ] Test JSON output format
- [ ] Achieve >80% test coverage

---

## Phase 4: TUI Interface Implementation

### 4.1 TUI Framework Setup
- [ ] Add `github.com/charmbracelet/bubbletea` dependency
- [ ] Add `github.com/charmbracelet/lipgloss` for styling
- [ ] Implement `internal/tui/app.go` (main application)
- [ ] Create application state model
- [ ] Set up event handling loop
- [ ] Implement graceful shutdown

### 4.2 Application State Management
- [ ] Implement `internal/tui/model.go`
- [ ] Create app state struct
- [ ] Implement view state tracking
- [ ] Add session state management
- [ ] Create account management state
- [ ] Add UI state (selection, filter, etc.)
- [ ] Implement state serialization for testing

### 4.3 Theme System
- [ ] Implement `internal/tui/styles/themes.go`
- [ ] Create color themes:
  - [ ] Dark (default)
  - [ ] Light
  - [ ] Nord
  - [ ] Dracula
- [ ] Define color values per theme
- [ ] Create style helpers
- [ ] Test theme rendering

### 4.4 Account List Component
- [ ] Implement `internal/tui/components/list.go`
- [ ] Create scrollable list with selection
- [ ] Implement keyboard navigation (↑↓ or j/k)
- [ ] Add filtering/search capability
- [ ] Display account name, issuer, current code
- [ ] Show favorite/pin indicator
- [ ] Implement list wrapping
- [ ] Test navigation and selection

### 4.5 Code Display Component
- [ ] Implement `internal/tui/components/code.go`
- [ ] Create large code display
- [ ] Implement animated progress bar
- [ ] Add color transitions based on time:
  - [ ] Green (>20s remaining)
  - [ ] Yellow (10-20s remaining)
  - [ ] Red (<10s remaining)
- [ ] Update animation at 100ms intervals
- [ ] Display time remaining
- [ ] Test animation updates

### 4.6 Account Form Component
- [ ] Implement `internal/tui/components/form.go`
- [ ] Create input fields for:
  - [ ] Name
  - [ ] Username
  - [ ] Issuer
  - [ ] Secret (hidden input)
  - [ ] Algorithm (dropdown)
  - [ ] Digits (dropdown)
  - [ ] Period (dropdown)
  - [ ] Tags (comma-separated)
- [ ] Implement field validation feedback
- [ ] Add secret show/hide toggle
- [ ] Create focus management
- [ ] Test form interactions

### 4.7 Search Component
- [ ] Implement `internal/tui/components/search.go`
- [ ] Create search bar with fuzzy matching
- [ ] Real-time filtering as user types
- [ ] Implement match highlighting
- [ ] Add search result count
- [ ] Test fuzzy matching algorithm

### 4.8 Dialog Component
- [ ] Implement `internal/tui/components/dialog.go`
- [ ] Create confirmation dialogs
- [ ] Create error message dialogs
- [ ] Create password prompt dialogs
- [ ] Implement focus management
- [ ] Add keyboard shortcuts (y/n for confirmation)
- [ ] Test dialog flow

### 4.9 Status Bar Component
- [ ] Implement status bar display
- [ ] Show current filter
- [ ] Display account count (e.g., "3/5 accounts")
- [ ] Show keyboard hints
- [ ] Display notifications
- [ ] Update dynamically

### 4.10 Key Bindings
- [ ] Implement `internal/tui/keys.go`
- [ ] Define key bindings:
  - [ ] ↑↓/j/k: Navigate
  - [ ] Tab: Switch panels
  - [ ] /: Activate search
  - [ ] Enter: Copy code
  - [ ] e: Edit account
  - [ ] a: Add account
  - [ ] d: Delete account
  - [ ] c: Copy code
  - [ ] ?: Show help
  - [ ] q: Quit
  - [ ] Esc: Cancel/Back
- [ ] Document all shortcuts
- [ ] Make shortcuts customizable

### 4.11 Update Logic
- [ ] Implement `internal/tui/update.go`
- [ ] Handle keyboard input
- [ ] Handle mouse input
- [ ] Update codes on timer
- [ ] Handle command results
- [ ] Update UI state
- [ ] Test state transitions

### 4.12 View/Rendering
- [ ] Implement `internal/tui/view.go`
- [ ] Create main view layout (split pane)
- [ ] Implement account list view
- [ ] Implement account detail view
- [ ] Implement add/edit form view
- [ ] Implement settings view
- [ ] Implement help view
- [ ] Test rendering on various terminal sizes

### 4.13 Main TUI Views

#### 4.13.1 Main View
- [ ] Implement main application view
- [ ] Left pane: Account list
- [ ] Right pane: Account details + code
- [ ] Bottom: Status bar with keyboard hints
- [ ] Live code updates
- [ ] Test layout responsiveness

#### 4.13.2 Add Account View
- [ ] Implement add account form view
- [ ] All account fields
- [ ] Input validation
- [ ] Submit/Cancel buttons
- [ ] Test form submission

#### 4.13.3 Edit Account View
- [ ] Implement edit account form view
- [ ] Pre-populate existing values
- [ ] Field validation
- [ ] Submit/Cancel buttons
- [ ] Test field updates

#### 4.13.4 Settings View
- [ ] Implement settings configuration
- [ ] Editable config options
- [ ] Theme selection
- [ ] Save/Cancel buttons
- [ ] Test settings persistence

#### 4.13.5 Help View
- [ ] Implement help/shortcuts view
- [ ] Display all key bindings
- [ ] Organize by category
- [ ] Make searchable
- [ ] Test help display

### 4.14 Live Updates
- [ ] Implement timer for code updates
- [ ] Update every 1 second
- [ ] Refresh codes at period boundary
- [ ] Update progress bar animation
- [ ] Smooth time remaining display
- [ ] Test timer accuracy

### 4.15 Clipboard Integration in TUI
- [ ] Implement copy on Enter
- [ ] Show copy notification
- [ ] Auto-clear clipboard after timeout
- [ ] Test clipboard operations

### 4.16 Error Handling & Notifications
- [ ] Show error dialogs
- [ ] Display validation errors
- [ ] Show success notifications
- [ ] Handle decryption failures
- [ ] Test error scenarios

### 4.17 TUI Testing
- [ ] Create unit tests for components
- [ ] Test state transitions
- [ ] Test rendering (snapshot tests)
- [ ] Test keyboard input handling
- [ ] Test with different terminal sizes
- [ ] Performance testing (no lag)
- [ ] Achieve >70% test coverage

---

## Phase 5: QR Code Support & Polish

### 5.1 QR Code Parsing
- [ ] Implement `internal/qr/parse.go`
- [ ] Parse image files (PNG, JPEG)
- [ ] Extract QR code data
- [ ] Support otpauth:// URI extraction
- [ ] Validate extracted data
- [ ] Add error handling for invalid QR codes
- [ ] Test with sample QR codes

### 5.2 QR Code Generation
- [ ] Implement `internal/qr/generate.go`
- [ ] Generate otpauth:// URI from account
- [ ] Create QR code from URI
- [ ] Format QR code for file output
- [ ] Test QR generation and scanning

### 5.3 Terminal QR Display
- [ ] Implement `internal/qr/terminal.go`
- [ ] Render QR code in terminal
- [ ] Support ASCII/Unicode rendering
- [ ] Test display on various terminals

### 5.4 Import from Other Authenticators

#### 5.4.1 Aegis Import
- [ ] Parse Aegis backup format
- [ ] Extract account data
- [ ] Handle encrypted Aegis backups
- [ ] Map fields to gotp format
- [ ] Test Aegis import

#### 5.4.2 Authy Import
- [ ] Parse Authy export format
- [ ] Extract account data
- [ ] Handle Authy-specific fields
- [ ] Map fields to gotp format
- [ ] Test Authy import

#### 5.4.3 Google Authenticator Import
- [ ] Parse Google Authenticator export
- [ ] Handle QR code migration codes
- [ ] Extract account data
- [ ] Map fields to gotp format
- [ ] Test Google import

### 5.5 Documentation

#### 5.5.1 README.md
- [ ] Write project description
- [ ] Add feature list
- [ ] Document installation instructions
- [ ] Create quick start guide
- [ ] Document all CLI commands
- [ ] Document TUI mode
- [ ] Add security information
- [ ] Document configuration
- [ ] Add contributing guidelines
- [ ] Add license information

#### 5.5.2 Security Documentation
- [ ] Create SECURITY.md
- [ ] Explain encryption approach
- [ ] Document threat model
- [ ] List security features
- [ ] Explain memory safety
- [ ] Document password requirements
- [ ] Add security best practices

#### 5.5.3 Architecture Documentation
- [ ] Create ARCHITECTURE.md
- [ ] Document project structure
- [ ] Explain module responsibilities
- [ ] Document data flow
- [ ] Add component diagrams
- [ ] Explain design decisions

#### 5.5.4 Contributing Guidelines
- [ ] Create CONTRIBUTING.md
- [ ] Explain development setup
- [ ] Document coding standards
- [ ] Add testing requirements
- [ ] Explain PR process
- [ ] List code of conduct

### 5.6 Build & Release Setup

#### 5.6.1 Makefile
- [ ] Create comprehensive Makefile
- [ ] `make build` - Build binaries
- [ ] `make test` - Run tests
- [ ] `make test-coverage` - Generate coverage report
- [ ] `make lint` - Run linter
- [ ] `make clean` - Clean build artifacts
- [ ] `make install` - Install binary
- [ ] `make release` - Build release binaries

#### 5.6.2 Build Scripts
- [ ] Create build.sh for local builds
- [ ] Create release.sh for release automation
- [ ] Set version from git tags
- [ ] Build for all platforms
- [ ] Generate checksums
- [ ] Create release artifacts

### 5.7 CI/CD Pipeline

#### 5.7.1 GitHub Actions Workflow
- [ ] Create `.github/workflows/ci.yml`
- [ ] Run tests on push/PR
- [ ] Test on Linux, macOS, Windows
- [ ] Test on Go 1.21+
- [ ] Run linter checks
- [ ] Generate coverage reports
- [ ] Upload coverage to Codecov

#### 5.7.2 Release Workflow
- [ ] Create `.github/workflows/release.yml`
- [ ] Build on tag push
- [ ] Create release artifacts
- [ ] Generate checksums
- [ ] Create GPG signature
- [ ] Publish to GitHub Releases
- [ ] Publish to package managers

### 5.8 Installation Methods

#### 5.8.1 Homebrew (macOS)
- [ ] Create Homebrew formula
- [ ] Test installation
- [ ] Document installation

#### 5.8.2 Package Managers (Linux)
- [ ] Create AUR package (Arch)
- [ ] Support apt/deb (Debian/Ubuntu)
- [ ] Support rpm (Red Hat/Fedora)
- [ ] Test all package managers

#### 5.8.3 Scoop (Windows)
- [ ] Create Scoop manifest
- [ ] Test installation
- [ ] Document installation

### 5.9 Completion Scripts
- [ ] Generate bash completion
- [ ] Generate zsh completion
- [ ] Generate fish completion
- [ ] Generate PowerShell completion
- [ ] Test completions on each shell

### 5.10 Man Pages
- [ ] Create man page for gotp(1)
- [ ] Create man pages for each command:
  - [ ] gotp-init(1)
  - [ ] gotp-add(1)
  - [ ] gotp-get(1)
  - [ ] gotp-list(1)
  - [ ] gotp-remove(1)
  - [ ] gotp-edit(1)
  - [ ] gotp-export(1)
  - [ ] gotp-import(1)
  - [ ] gotp-passwd(1)
  - [ ] gotp-tui(1)
  - [ ] gotp-completion(1)

### 5.11 Performance Optimization
- [ ] Profile startup time (target: <100ms)
- [ ] Profile code generation (target: <1ms)
- [ ] Profile vault loading (target: <500ms for 1000 accounts)
- [ ] Optimize hot paths
- [ ] Memory usage optimization
- [ ] Benchmark suite

### 5.12 Final Testing & QA

#### 5.12.1 Functional Testing
- [ ] Test all CLI commands
- [ ] Test TUI functionality
- [ ] Test on all target platforms
- [ ] Test with various account counts
- [ ] Test import/export workflows
- [ ] Test password change
- [ ] Test backup/restore

#### 5.12.2 Security Testing
- [ ] Verify no plaintext secrets in logs
- [ ] Check memory for unzeroed data
- [ ] Test encryption strength
- [ ] Test key derivation
- [ ] Test authentication failure
- [ ] Verify session timeout
- [ ] Test clipboard clearing

#### 5.12.3 Cross-Platform Testing
- [ ] Test on Linux (x86_64, arm64)
- [ ] Test on macOS (x86_64, arm64)
- [ ] Test on Windows (x86_64)
- [ ] Test various terminal emulators
- [ ] Test clipboard on all platforms
- [ ] Test file paths on all platforms

### 5.13 Release Preparation
- [ ] Finalize version number
- [ ] Update CHANGELOG.md
- [ ] Create release tag
- [ ] Generate release notes
- [ ] Build all release artifacts
- [ ] Sign release artifacts
- [ ] Create GitHub Release
- [ ] Publish to package managers
- [ ] Announce release

---

## Success Criteria Tracking

### Functional Criteria
- [ ] Generate valid TOTP codes matching reference implementations
- [ ] All CLI commands work as specified
- [ ] TUI launches and is fully functional
- [ ] Import/export works with common formats
- [ ] Cross-platform builds successful

### Security Criteria
- [ ] Passes security review checklist
- [ ] No plaintext secrets in logs or errors
- [ ] Encryption implementation auditable
- [ ] Memory properly zeroed after use
- [ ] Session management working correctly

### Quality Criteria
- [ ] >90% test coverage on crypto code
- [ ] >80% test coverage overall
- [ ] Zero linter warnings
- [ ] Documentation complete
- [ ] Code follows Go best practices

### Performance Criteria
- [ ] Startup time <100ms
- [ ] Code generation <1ms
- [ ] Vault loading <500ms for 1000 accounts
- [ ] TUI renders smoothly (60+ FPS)
- [ ] No memory leaks

---

## Overall Status Summary

**Phase 1: Core TOTP Engine** - [x] Complete
**Phase 2: Encryption & Storage** - [ ] Complete
**Phase 3: CLI Interface** - [ ] Complete
**Phase 4: TUI Interface** - [ ] Complete
**Phase 5: QR & Polish** - [ ] Complete

**PROJECT COMPLETE** - [ ] All phases finished and tested
