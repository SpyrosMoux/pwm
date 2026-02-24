# pwm - Password Manager

A simple, secure, and lightweight command-line password manager written in Go. Store, retrieve, and manage your
passwords safely with AES encryption.

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
    - [Command Reference](#command-reference)
- [Configuration](#configuration)
- [Security](#security)
- [Project Structure](#project-structure)
- [Development](#development)
- [License](#license)

## Features

- 🔐 **AES Encryption**: All sensitive fields (URL, username, password, description) are encrypted using AES with PKCS7
  padding
- 🔑 **User-Configurable Cipher Key**: Encryption key is set via the `PWM_CIPHER_KEY` environment variable
- 💾 **Local Storage**: All secrets are stored locally on your machine in `~/.pwm`
- 🖥️ **Terminal-based**: Simple and intuitive command-line interface powered by Cobra
- 📋 **Clipboard Integration**: Quickly copy passwords to your clipboard
- 🔄 **Full CRUD Operations**: Create, read, and delete secrets
- ⚡ **Fast & Lightweight**: Written in Go with minimal dependencies
- 🎯 **Rich Secret Information**: Store URLs, usernames, passwords, and descriptions
- 🔧 **Password Auto-Generation**: Type `a` when prompted for a password to auto-generate a secure one
- 🌐 **Cross-Platform**: Pre-built binaries for Linux, macOS, and Windows (amd64 and arm64)

## How It Works

### Secret Storage

Secrets are stored as JSON files in your local `~/.pwm` directory (or a custom location of your choice). Each secret
contains:

- **Name**: The identifier for the secret
- **URL**: The website or service URL (encrypted)
- **Username**: Your username for the service (encrypted)
- **Password**: Your password (encrypted)
- **Description**: Additional notes or details (encrypted)

### Encryption

When you create a secret:

1. The tool reads your 32-character cipher key from the `PWM_CIPHER_KEY` environment variable
2. All sensitive fields (URL, username, password, description) are encrypted using AES encryption
3. The encrypted data is stored as hex-encoded strings in JSON format
4. Your cipher key is never stored on disk—it only exists in your environment

When you retrieve a secret:

1. The tool reads the cipher key from `PWM_CIPHER_KEY`
2. The stored data is decrypted using your key
3. The plaintext secret is displayed or the password is copied to your clipboard

## Installation

### Prerequisites

- **Go 1.25.0 or later** (if building from source)
- **macOS, Linux, or Windows**
- The `PWM_CIPHER_KEY` environment variable must be set (see [Setup](#setup))

### From GitHub Releases

Pre-built binaries are available on the [Releases](https://github.com/SpyrosMoux/pwm/releases) page for:

| Platform       | Architecture                         |
|----------------|--------------------------------------|
| Linux          | amd64, arm64                         |
| macOS (Darwin) | Intel (amd64), Apple Silicon (arm64) |
| Windows        | amd64, arm64                         |

Download the appropriate archive for your platform, extract it, and place the binary in your `PATH`.

### Building from Source

1. Clone the repository:

```bash
git clone https://github.com/SpyrosMoux/pwm.git
cd pwm
```

2. Build for all platforms using the provided Makefile:

```bash
make build
```

This will produce cross-compiled binaries in the `bin/` directory:

- `bin/pwm-linux-amd64`
- `bin/pwm-linux-arm64`
- `bin/pwm-darwin-amd64`
- `bin/pwm-darwin-arm64`
- `bin/pwm-windows-amd64.exe`
- `bin/pwm-windows-arm64.exe`

To build for a specific platform only:

```bash
make build-linux          # Linux amd64
make build-linux-arm      # Linux arm64
make build-mac-intel      # macOS Intel
make build-mac-arm        # macOS Apple Silicon
make build-windows        # Windows amd64
make build-windows-arm    # Windows arm64
```

Or use Go directly:

```bash
go mod tidy
go build -o bin/pwm .
```

## Usage

### Setup

Before using `pwm`, you must set the `PWM_CIPHER_KEY` environment variable to a **32-character** string. This key is
used to encrypt and decrypt your secrets.

**Linux / macOS:**

```bash
export PWM_CIPHER_KEY="your-32-character-secret-key!!!"
```

To persist across sessions, add it to your shell profile (`~/.bashrc`, `~/.zshrc`, etc.).

**Windows (PowerShell):**

```powershell
$env:PWM_CIPHER_KEY = "your-32-character-secret-key!!!"
```

To persist, set it as a system or user environment variable.

> ⚠️ **Important**: The key must be exactly 32 characters long. Keep it secret and don't lose it—without it, you cannot
> decrypt your secrets.

### Basic Commands

```bash
# Create a new secret
pwm create my_secret

# List all stored secrets
pwm ls

# Retrieve and display a secret
pwm my_secret

# Copy a secret's password to clipboard
pwm cp my_secret

# Delete a secret
pwm rm my_secret

# Show help
pwm --help
```

### Index Selection

You can also select secrets by numeric index instead of typing the full name. Run `pwm ls` to see a numbered tree (indices are 1-based and oldest secrets have lower numbers):

```bash
./pwm ls
# /home/username/.pwm
# ├── {1} old_secret
# ├── {2} work/github
# └── {3} personal/email
```

Use the index in any command that accepts a secret name. For example:

```bash
# Show secret with index 2 (same as `pwm "work/github"`)
pwm 2

# Copy password using index
pwm cp 2

# Remove using index
pwm rm 2
```


### Command Reference

#### `pwm create <secret_name>`

Creates a new secret. You'll be prompted to enter:

- URL (optional)
- Username
- Password (enter `a` to auto-generate a secure password)
- Description (optional)

```bash
pwm create github
```

#### `pwm <secret_name>`

Retrieves and displays a secret.

```bash
pwm github
```

Output:

```
Name: github
Url: https://github.com
Username: your_username
Password: your_password
Description: My GitHub account
```

#### `pwm cp <secret_name>`

Copies the password of a secret to your clipboard.

```bash
pwm cp github
# Password copied to clipboard
```

#### `pwm ls`

Lists all secrets stored in your storage location as a tree structure.

```bash
pwm ls

# Output:
# /home/username/.pwm
# ├── github
# ├── gmail
# └── twitter
```

#### `pwm rm <secret_name>`

Deletes a secret permanently.

```bash
pwm rm old_service
# Removed secret: old_service
```

#### `pwm update <secret_name|index>`

Updates an existing secret's fields (URL, username, password, description). You'll be prompted to enter new values for each field. Press Enter to skip a field and keep its current value.

```bash
# Update by name
pwm update github

# Update by numeric index
pwm update 2

# Example session:
# Updating secret: github
# (Press Enter to skip a field)
#
# URL [https://github.com]: https://new-github-url.com
# Username [old_user]: 
# Password [***]: new_secure_password
# Description [My GitHub account]: Updated GitHub account
# Secret updated successfully: github
```

**Atomic Updates**: The update operation writes to a temporary file and atomically replaces the original file to prevent corruption if the process is interrupted.

#### `pwm --help`

Shows help information and available commands.

#### Custom Storage Location

By default, secrets are stored in `~/.pwm`. You can specify a different location using the `--location` flag:

```bash
pwm --location /path/to/custom/location create my_secret
pwm --location /path/to/custom/location ls
pwm --location /path/to/custom/location my_secret
```

## Configuration

### Storage Location

The default storage location is `~/.pwm`. This can be customized using the `--location` flag:

```bash
pwm --location ~/.my_passwords ls
```

This applies to that specific command. To use a custom location permanently, you can create an alias:

```bash
alias pwm_custom='pwm --location ~/.my_passwords'
pwm_custom ls
```

### Cipher Key

Your cipher key is read from the `PWM_CIPHER_KEY` environment variable. It must be exactly **32 characters** long and
is used as the AES encryption key for all encrypt/decrypt operations.

The cipher key is **never stored** by `pwm`. If you lose your cipher key, you will not be able to decrypt your secrets.

## Security

### Important Security Notes

- The `PWM_CIPHER_KEY` environment variable must be set and kept secure. Anyone with access to this key and your
  `~/.pwm` directory can decrypt your secrets.
- All sensitive fields (URL, username, password, and description) are encrypted—not just the password.

### Best Practices

1. **Protect Your Cipher Key**: Store your `PWM_CIPHER_KEY` securely. Avoid committing it to version control or sharing
   it.
2. **Protect Your Storage Directory**: The `~/.pwm` directory contains encrypted data. Keep it secure and don't share
   it.
3. **Keep Go Updated**: Make sure you're using an up-to-date version of Go to benefit from the latest security patches.
4. **System Security**: This tool is only as secure as your system. If your computer is compromised, your secrets may be
   at risk.

### Encryption Details

- **Algorithm**: AES (Advanced Encryption Standard)
- **Padding**: PKCS7
- **Key**: User-configurable via the `PWM_CIPHER_KEY` environment variable (32 characters / 256-bit)
- **Encrypted Fields**: URL, Username, Password, Description
- **Storage Format**: Hex-encoded ciphertext in JSON files

## Project Structure

```
pwm/
├── main.go                 # Entry point
├── go.mod                  # Go module definition
├── go.sum                  # Dependency checksums
├── Makefile                # Cross-platform build scripts
├── README.md               # This file
├── LICENSE                 # GNU GPLv3 License
├── .github/
│   └── workflows/
│       ├── ci.yml          # CI pipeline (build on push/PR)
│       └── release.yml     # Release pipeline (build & publish on tag)
├── cmd/
│   ├── root.go             # Root command, flags, and cipher key loading
│   ├── create.go           # Create command
│   ├── ls.go               # List command
│   ├── cp.go               # Copy command
│   ├── rm.go               # Remove command
│   └── pwm.go              # Core command implementations
└── internal/
    ├── crypto/
    │   └── crypto.go       # AES encryption/decryption and PKCS7 padding
    ├── helpers/
    │   ├── inputs.go       # User input handling (text & secret input)
    │   └── outputs.go      # Output helpers (info, warn, error)
    └── models/
        └── secret.go       # Secret data model, encrypt/decrypt methods
```

## Development

### Dependencies

- **[github.com/spf13/cobra](https://github.com/spf13/cobra)**: CLI framework
- **[github.com/SpyrosMoux/passwdgen](https://github.com/SpyrosMoux/passwdgen)**: Password auto-generation
- **[golang.design/x/clipboard](https://pkg.go.dev/golang.design/x/clipboard)**: Clipboard integration
- **[golang.org/x/term](https://pkg.go.dev/golang.org/x/term)**: Terminal handling for secret input

Install dependencies:

```bash
go mod tidy
```

### Building

```bash
# Build for all platforms
make build

# Or build with Go directly
go build -o bin/pwm .
```

### Running Tests

```bash
go test ./...
```

### CI/CD

The project uses GitHub Actions for continuous integration and releases:

- **CI** (`.github/workflows/ci.yml`): Builds on every push to `main` and on pull requests across all supported
  platforms.
- **Release** (`.github/workflows/release.yml`): On tagged pushes (`v*`), builds for all platforms, packages archives,
  generates SHA-256 checksums, and creates a GitHub Release with all assets.

## License

This project is licensed under the GNU General Public License v3.0 - see the [LICENSE](LICENSE) file for details.

## Author

Created by Spyros Mouchlianitis

## Contributing

Contributions are welcome! Feel free to open an issue or submit a pull request.
