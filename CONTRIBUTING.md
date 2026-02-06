# Contributing to gotp

Thank you for your interest in contributing to gotp! This document provides guidelines and instructions for contributing to the project.

## Table of Contents

- [Getting Started](#getting-started)
- [Development Setup](#development-setup)
- [Code Style](#code-style)
- [Testing](#testing)
- [Pull Request Process](#pull-request-process)
- [Reporting Bugs](#reporting-bugs)
- [Requesting Features](#requesting-features)
- [Security Vulnerabilities](#security-vulnerabilities)
- [Code of Conduct](#code-of-conduct)

## Getting Started

### Prerequisites

- Go 1.25 or higher
- Git
- A GitHub account

### Fork and Clone

1. Fork the repository on GitHub
2. Clone your fork locally:
   ```bash
   git clone https://github.com/YOUR_USERNAME/gotp.git
   cd gotp
   ```

3. Add the upstream repository:
   ```bash
   git remote add upstream https://github.com/zulfikawr/gotp.git
   ```

## Development Setup

### Install Dependencies

```bash
go mod download
go mod tidy
```

### Build

```bash
go build -o gotp ./cmd/gotp
```

### Run Tests

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run specific package tests
go test ./internal/crypto/...
go test ./internal/totp/...
```

### Run Linter

```bash
# Install golangci-lint if not already installed
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Run linter
golangci-lint run
```

### Development Workflow

1. Create a feature branch:
   ```bash
   git checkout -b feature/your-feature-name
   ```

2. Make your changes

3. Run tests and linter:
   ```bash
   go test ./...
   golangci-lint run
   ```

4. Commit your changes:
   ```bash
   git add .
   git commit -m "feat: add your feature description"
   ```

5. Push to your fork:
   ```bash
   git push origin feature/your-feature-name
   ```

6. Create a Pull Request on GitHub

## Code Style

### Go Conventions

- Follow [Effective Go](https://golang.org/doc/effective_go)
- Use `gofmt` for formatting (automatic with most editors)
- Use `golangci-lint` for static analysis

### Naming Conventions

- **Packages**: Short, lowercase, single word (e.g., `crypto`, `vault`)
- **Functions**: CamelCase, descriptive (e.g., `GenerateTOTP`, `LoadVault`)
- **Variables**: camelCase, descriptive (e.g., `vaultPath`, `targetAccount`)
- **Constants**: UPPER_SNAKE_CASE (e.g., `SHA1`, `SHA256`)

### Error Handling

- Return errors instead of panicking
- Use `fmt.Errorf` with `%w` for wrapping errors
- Check errors immediately after function calls
- Provide context in error messages

```go
// Good
key, err := crypto.DeriveKey(password, salt, params)
if err != nil {
    return nil, fmt.Errorf("failed to derive key: %w", err)
}

// Bad
key, _ := crypto.DeriveKey(password, salt, params)
```

### Comments

- Document exported functions with Godoc comments
- Use inline comments for complex logic
- Keep comments concise and meaningful

```go
// GenerateTOTP generates a Time-based One-Time Password (TOTP) as defined in RFC 6238.
// It calculates the time step counter based on the provided timestamp and period,
// then calls GenerateHOTP to produce the code.
func GenerateTOTP(params TOTPParams) (string, error) {
    // ...
}
```

### File Organization

- One type per file (when possible)
- Test files: `*_test.go`
- Platform-specific files: `*_linux.go`, `*_darwin.go`, `*_windows.go`

## Testing

### Writing Tests

- Test files must be named `*_test.go`
- Test functions must start with `Test`
- Use table-driven tests for multiple cases
- Test both success and failure cases

```go
func TestGenerateTOTP(t *testing.T) {
    tests := []struct {
        name     string
        params   TOTPParams
        expected string
        wantErr  bool
    }{
        {
            name: "RFC 6238 test vector 1",
            params: TOTPParams{
                Secret:    []byte{0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37, 0x38, 0x39, 0x30, 0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37, 0x38, 0x39, 0x30},
                Timestamp: time.Unix(59, 0),
                Period:    30,
                Digits:    6,
                Algorithm: SHA1,
            },
            expected: "287082",
            wantErr:  false,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result, err := GenerateTOTP(tt.params)
            if (err != nil) != tt.wantErr {
                t.Errorf("GenerateTOTP() error = %v, wantErr %v", err, tt.wantErr)
                return
            }
            if !tt.wantErr && result != tt.expected {
                t.Errorf("GenerateTOTP() = %v, want %v", result, tt.expected)
            }
        })
    }
}
```

### Test Coverage

- Aim for 100% coverage on cryptographic code
- Aim for >80% overall coverage
- Use `go test -cover` to check coverage

### Benchmarks

Add benchmarks for performance-critical functions:

```go
func BenchmarkGenerateTOTP(b *testing.B) {
    params := TOTPParams{
        Secret:    []byte{0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37, 0x38, 0x39, 0x30, 0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37, 0x38, 0x39, 0x30},
        Timestamp: time.Now(),
        Period:    30,
        Digits:    6,
        Algorithm: SHA1,
    }

    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        GenerateTOTP(params)
    }
}
```

## Pull Request Process

### PR Checklist

Before submitting a PR, ensure:

- [ ] Code compiles without errors
- [ ] All tests pass
- [ ] Linter passes with no warnings
- [ ] Code follows project style guidelines
- [ ] New functionality has tests
- [ ] Documentation is updated
- [ ] Commit messages follow Conventional Commits

### Commit Message Format

Use [Conventional Commits](https://www.conventionalcommits.org/):

```
<type>[optional scope]: <description>

[optional body]

[optional footer(s)]
```

**Types:**
- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation changes
- `style`: Code style changes (formatting, semicolons)
- `refactor`: Code refactoring without behavior change
- `test`: Adding or updating tests
- `chore`: Maintenance tasks

**Examples:**
```
feat: add QR code generation command

feat(crypto): add AES-256-GCM encryption

fix: handle empty vault gracefully

docs: update README with import examples
```

### PR Description

Include:
1. **What**: What does this PR do?
2. **Why**: Why is this change necessary?
3. **How**: How does it work?
4. **Testing**: How was it tested?
5. **Breaking Changes**: Any breaking changes?

### Review Process

1. PR will be reviewed by maintainers
2. Address review comments
3. PR will be merged once approved
4. Squash and merge is preferred

## Reporting Bugs

### Bug Report Template

Create an issue with the following information:

```markdown
## Bug Description
[Clear description of the bug]

## Steps to Reproduce
1. Step 1
2. Step 2
3. Step 3

## Expected Behavior
[What you expected to happen]

## Actual Behavior
[What actually happened]

## Environment
- OS: [e.g., Ubuntu 22.04]
- Go Version: [e.g., 1.25.5]
- gotp Version: [e.g., 0.1.0]

## Additional Context
[Screenshots, logs, etc.]
```

### Bug Triage

- Critical bugs will be prioritized
- Security bugs should be reported via security email
- Bugs without reproduction steps may be closed

## Requesting Features

### Feature Request Template

```markdown
## Feature Description
[Clear description of the feature]

## Use Case
[Why is this feature needed?]

## Proposed Solution
[How should it work?]

## Alternatives Considered
[Other approaches you considered]

## Additional Context
[Any other relevant information]
```

### Feature Prioritization

Features are prioritized based on:
- User demand
- Security implications
- Implementation complexity
- Alignment with project goals

## Security Vulnerabilities

**⚠️ IMPORTANT: Do not open public issues for security vulnerabilities**

### Reporting Security Issues

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

## Code of Conduct

### Our Pledge

We are committed to providing a friendly, safe, and welcoming environment for all contributors.

### Our Standards

**Positive Behavior:**
- Use welcoming and inclusive language
- Be respectful of differing viewpoints
- Gracefully accept constructive criticism
- Focus on what is best for the community

**Unacceptable Behavior:**
- Harassment, insults, or derogatory comments
- Trolling or insulting/derogatory remarks
- Publishing others' private information
- Other conduct which could reasonably be considered inappropriate

### Enforcement

Project maintainers are responsible for clarifying and enforcing acceptable behavior. They will take appropriate and fair corrective action in response to any behavior that they deem inappropriate.

### Consequences

Violations may result in:
- Temporary ban from the community
- Permanent ban from the community
- Removal from project contributor list

## Recognition

Contributors will be recognized in:
- The project's contributors list
- Release notes
- README.md (for significant contributions)

## Questions?

Feel free to:
- Open a discussion on GitHub
- Join our community chat (if available)
- Email the maintainers

## Thank You!

Thank you for your interest in contributing to gotp. Your contributions help make this project better for everyone!
