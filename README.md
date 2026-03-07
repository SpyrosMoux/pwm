# pwm - Password Manager

`pwm` is a simple command-line password manager written in Go.

It stores secrets as local files and encrypts secret fields using AES with PKCS7 padding.

## Table of Contents

- [Features](#features)
- [How It Works](#how-it-works)
- [Installation](#installation)
  - [Prerequisites](#prerequisites)
  - [From GitHub Releases](#from-github-releases)
  - [Building from Source](#building-from-source)
- [Usage](#usage)
  - [Setup](#setup)
  - [Basic Commands](#basic-commands)
  - [Index Selection](#index-selection)
  - [Command Reference](#command-reference)
- [Configuration](#configuration)
- [Security](#security)
- [Project Structure](#project-structure)
- [Development](#development)
- [License](#license)

## Features

- AES encryption for secret fields (`url`, `username`, `password`, `description`)
- Cipher key loaded from `PWM_CIPHER_KEY` (must be 32 characters)
- Local file storage in `~/.pwm` by default
- Numbered secret listing via `pwm ls`
- Name-or-index targeting for read/copy/remove/update commands
- Create, read, copy, update, and delete workflows
- Password auto-generation during create and update
- Atomic secret updates with original modification time preserved
- Custom storage location via `--location`
- Cross-platform builds for Linux, macOS, and Windows (amd64 + arm64)

## How It Works

### Secret Storage

Each secret is stored as a JSON file in `~/.pwm` (or a custom `--location` directory).

A secret contains:

- `Name`
- `Url` (encrypted)
- `Username` (encrypted)
- `Password` (encrypted)
- `Description` (encrypted)

### Encryption Flow

When creating or updating a secret:

1. `pwm` reads `PWM_CIPHER_KEY` from the environment.
2. Sensitive fields are encrypted.
3. Encrypted values are stored as hex strings in JSON.

When reading/copying a secret:

1. `pwm` reads `PWM_CIPHER_KEY`.
2. Stored values are decrypted.
3. The secret is printed or copied to the clipboard.

## Installation

### Prerequisites

- Go `1.25.0` or later (for source builds)
- macOS, Linux, or Windows
- `PWM_CIPHER_KEY` configured before usage

### From GitHub Releases

Prebuilt binaries are available at:

- <https://github.com/SpyrosMoux/pwm/releases>

Supported release targets:

- Linux: `amd64`, `arm64`
- macOS: `amd64`, `arm64`
- Windows: `amd64`, `arm64`

### Building from Source

```bash
git clone https://github.com/SpyrosMoux/pwm.git
cd pwm
make build
```

This creates binaries under `bin/`:

- `bin/pwm-linux-amd64`
- `bin/pwm-linux-arm64`
- `bin/pwm-darwin-amd64`
- `bin/pwm-darwin-arm64`
- `bin/pwm-windows-amd64.exe`
- `bin/pwm-windows-arm64.exe`

Build a single target:

```bash
make build-linux
make build-linux-arm
make build-mac-intel
make build-mac-arm
make build-windows
make build-windows-arm
```

Or build directly with Go:

```bash
go mod tidy
go build -o bin/pwm .
```

## Usage

### Setup

Set `PWM_CIPHER_KEY` to an exact 32-character value.

Linux/macOS:

```bash
export PWM_CIPHER_KEY="your-32-character-secret-key!!!"
```

Windows (PowerShell):

```powershell
$env:PWM_CIPHER_KEY = "your-32-character-secret-key!!!"
```

If the key is missing or not 32 characters long, `pwm` exits with an error.

### Basic Commands

```bash
# Create a secret (name cannot be numeric-only)
pwm create github

# List secrets with numbered indices
pwm ls

# Show secret by name
pwm github

# Show secret by index
pwm 2

# Copy password by name or index
pwm cp github
pwm cp 2

# Update by name or index
pwm update github
pwm update 2

# Remove by name or index
pwm rm github
pwm rm 2
```

### Index Selection

`pwm ls` prints a tree view with 1-based numeric indices.

```bash
pwm ls
# /home/username/.pwm
# ├── {1} github
# ├── {2} gmail
# └── {3} twitter
```

These indices can be used anywhere a secret identifier is accepted (`pwm`, `cp`, `rm`, `update`).

### Command Reference

#### Argument Resolution (`name|index`)

When a command accepts `<name|index>`, `pwm` resolves the argument using these rules:

- Digits-only input is treated as an index.
- Index parsing uses `strconv.ParseInt` (base 10, 64-bit).
- Valid index range is `1..N` where `N` is the number of stored secrets.
- Out-of-range values return `index out of range (1..N)`.
- Overflow values return `invalid index: number too large`.
- Non-numeric input is treated as a secret name.

#### `pwm <name|index>`

Reads, decrypts, and prints a secret.

#### `pwm create <name>`

Creates a new secret by prompting for:

- URL
- Username
- Password (type `a` to auto-generate)
- Description

Notes:

- Secret names cannot be numeric-only.
- Use plain names such as `github`, `gmail`, `work_vpn`.

#### `pwm cp <name|index>`

Decrypts the secret and copies its password to the clipboard.

#### `pwm ls`

Lists secrets from the storage directory as a numbered tree.

#### `pwm rm <name|index>`

Permanently deletes a secret.

#### `pwm update <name|index>`

Updates an existing secret in place.

- Prompts each field with current value
- Press Enter to keep a field unchanged
- Password prompt supports `a` for auto-generation
- Writes through a temporary file and atomically replaces the original file
- Restores the original modification time so list index order stays stable

#### `pwm --location <dir> ...`

Overrides the default storage directory for a command.

Examples:

```bash
pwm --location /path/to/vault ls
pwm --location /path/to/vault create github
pwm --location /path/to/vault cp github
```

## Configuration

### `PWM_CIPHER_KEY`

- Required environment variable
- Must be exactly 32 characters
- Used for all encrypt/decrypt operations

### Storage Location

- Default: `~/.pwm`
- Override per command with `--location`

Example alias:

```bash
alias pwm-work='pwm --location ~/.pwm-work'
pwm-work ls
```

## Security

- Anyone with both your `PWM_CIPHER_KEY` and storage directory can decrypt your secrets.
- Keep your environment, shell history, and local storage secure.
- Do not commit `PWM_CIPHER_KEY` to source control.
- If you lose the key, stored secrets cannot be decrypted.

## Project Structure

```text
pwm/
├── main.go
├── go.mod
├── go.sum
├── Makefile
├── README.md
├── LICENSE
├── cmd/
│   ├── root.go        # root command and global flags
│   ├── pwm.go         # name/index argument resolution
│   ├── create.go
│   ├── ls.go
│   ├── cp.go
│   ├── rm.go
│   └── update.go
├── internal/
│   ├── crypto/        # encryption/decryption helpers
│   ├── models/        # Secret model
│   ├── helpers/       # input/output helpers
│   ├── store/         # filesystem storage abstraction
│   └── secrets/       # create/get/list/copy/remove/update service methods
└── .github/workflows/
    ├── ci.yml
    └── release.yml
```

## Development

### Dependencies

- `github.com/spf13/cobra`
- `github.com/SpyrosMoux/passwdgen`
- `golang.design/x/clipboard`
- `golang.org/x/term`

Install/update dependencies:

```bash
go mod tidy
```

Run tests:

```bash
go test ./...
```

Build all release binaries:

```bash
make build
```

### CI/CD

GitHub Actions workflows:

- `ci.yml`: builds all supported targets on push to `main` and pull requests
- `release.yml`: on `v*` tags, builds and packages artifacts, generates checksums, and creates a GitHub Release

## License

This project is licensed under the GNU General Public License v3.0.

See `LICENSE` for details.

## Author

Created by Spyros Mouchlianitis.

## Contributing

Contributions are welcome via issues and pull requests.
