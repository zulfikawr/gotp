# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [0.1.1] - 2026-02-05

### Changed
- **Secure Session Storage**: Replaced insecure plaintext session caching in `/tmp` with an encrypted, cross-platform implementation.
- **AES-256-GCM Encryption**: Sessions are now encrypted using a machine-specific key (derived from hostname and UID).
- **Private Config Storage**: Moved session files from `/tmp` to the user's private configuration directory (`~/.config/gotp/session.bin`) with `0600` permissions.
- **Memory Safety**: Replaced immutable `string` types with mutable `[]byte` for sensitive TOTP secrets to prevent lingering data in memory.
- **Timing Attack Protection**: Implemented `subtle.ConstantTimeCompare` for TOTP code validation to prevent timing-based side-channel attacks.
- **Hardened Password Input**: Strictly require a terminal (TTY) for master password entry to prevent accidental leakage in non-interactive environments or logs.
- **Improved Watch Mode**: Refactored watch mode cleanup logic for better terminal compatibility and cursor restoration.
- **Robust Importers**: Refactored Google Authenticator migration URI parsing using standard protobuf and base64 handling for improved reliability and security.
- **Atomic Vault Saving**: Implemented atomic file saving using `os.CreateTemp` and `os.Rename` to prevent data corruption and race conditions.

### Fixed
- Vulnerability where any local user could read the master key from `/tmp/gotp.session`.

## [0.1.0] - 2026-01-25

### Added
- **Core TOTP Engine**: Full implementation of RFC 6238 (TOTP) and RFC 4226 (HOTP).
- **Multi-Algorithm Support**: Support for SHA1, SHA256, and SHA512 hash algorithms.
- **Secure Storage**: AES-256-GCM authenticated encryption for vault storage.
- **Key Derivation**: Argon2id for strong master password key derivation.
- **Interactive CLI**: Full-featured CLI with commands for `init`, `add`, `get`, `list`, `remove`, `edit`, `passwd`, and `completion`.
- **Import/Export**: Support for importing from Aegis, Authy, and Google Authenticator (including migration QR codes).
- **QR Code Support**: Built-in QR code generation and parsing (from image files or terminal).
- **Clipboard Integration**: Secure copy-to-clipboard with automatic timeout clearing.
- **Session Caching**: Optional session caching to avoid repeated password prompts.
- **Theming**: Integrated Gruvbox-inspired color theming for the terminal UI.
- **Cross-Platform**: Support for Linux, macOS, and Windows.
