# pwm - Password Manager

A simple, secure, and lightweight command-line password manager written in Go. Store, retrieve, and manage your
passwords safely with AES encryption.

## Table of Contents

- [Features](#features)
- [How It Works](#how-it-works)
- [Installation](#installation)
    - [Prerequisites](#prerequisites)
    - [Building from Source](#building-from-source)
- [Usage](#usage)
    - [Basic Commands](#basic-commands)
    - [Command Reference](#command-reference)
- [Configuration](#configuration)
- [Security](#security)
- [Project Structure](#project-structure)
- [License](#license)

## Features

- 🔐 **AES Encryption**: Passwords are encrypted using AES encryption with PKCS7 padding
- 💾 **Local Storage**: All secrets are stored locally on your machine in `~/.pwm`
- 🖥️ **Terminal-based**: Simple and intuitive command-line interface powered by Cobra
- 📋 **Clipboard Integration**: Quickly copy passwords to your clipboard
- 🔄 **Full CRUD Operations**: Create, read, update, and delete secrets
- ⚡ **Fast & Lightweight**: Written in Go with minimal dependencies
- 🎯 **Rich Secret Information**: Store URLs, usernames, passwords, and descriptions

## How It Works

### Secret Storage

Secrets are stored as JSON files in your local `~/.pwm` directory (or a custom location of your choice). Each secret
contains:

- **Name**: The identifier for the secret
- **URL**: The website or service URL
- **Username**: Your username for the service
- **Password**: Your password (encrypted)
- **Description**: Additional notes or details

### Encryption

When you create a secret:

1. The tool prompts you to set a master password (used as the encryption key)
2. Your password and other sensitive data are encrypted using AES encryption
3. The encrypted data is stored as hex-encoded strings in JSON format
4. Your master password is never stored—you must remember it

When you retrieve a secret:

1. You provide your master password
2. The tool decrypts the stored data using your password as the key
3. The plaintext password is displayed or copied to your clipboard

## Installation

### Prerequisites

- **Go 1.25.0 or later** (if building from source)
- **macOS, Linux, or Windows** with a Unix-like shell

### Building from Source

1. Clone the repository:

```bash
git clone https://github.com/SpyrosMoux/pwm.git
cd pwm
```

2. Build and install using the provided Makefile:

```bash
make install
```

This will:

- Run `go mod tidy` to download dependencies
- Build the binary to `bin/pwm`
- Copy the binary to `/usr/local/bin/pwm` (requires sudo)
- Make it executable

Alternatively, build without installing to `/usr/local/bin`:

```bash
make build
# Binary will be at ./bin/pwm
```

Or use Go directly:

```bash
go mod tidy
go build -o bin/pwm main.go
```

## Usage

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

### Command Reference

#### `pwm create <secret_name>`

Creates a new secret. You'll be prompted to enter:

- URL (optional)
- Username
- Password
- Description (optional)

```bash
pwm create github
```

#### `pwm <secret_name>` or `pwm get <secret_name>`

Retrieves and displays a secret. Prompts you for your master password.

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

Lists all secrets stored in your default storage location.

```bash
pwm ls

# Output:
# /Users/username/.pwm
# github
# gmail
# twitter
```

#### `pwm rm <secret_name>`

Deletes a secret permanently.

```bash
pwm rm old_service
# Removed secret: old_service
```

#### `pwm --help`

Shows help information and available commands.

#### Custom Storage Location

By default, secrets are stored in `~/.pwm`. You can specify a different location:

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

This will persist for that specific command. To use a custom location permanently, you can create an alias:

```bash
alias pwm_custom='pwm --location ~/.my_passwords'
pwm_custom ls
```

### Master Password

Your master password is used to encrypt and decrypt secrets. It is **never stored** anywhere. You must remember it or
create a new master password for each use. If you forget your master password, you won't be able to decrypt your
secrets.

## Security

### Important Security Notes

⚠️ **Current Limitations**:

1. **Hardcoded Cipher Key**: Currently uses a hardcoded cipher key for encryption, which is **not production-ready**.
   All instances of pwm on a system share the same encryption key.
2. **Future Enhancement**: A user-set master password feature is planned to improve security.

### Best Practices

1. **Protect Your Storage Directory**: The `~/.pwm` directory contains encrypted data. Keep it secure and don't share
   it.
2. **Keep Go Updated**: Make sure you're using an up-to-date version of Go to benefit from the latest security patches.
3. **System Security**: This tool is only as secure as your system. If your computer is compromised, your secrets may be
   at risk.

### Encryption Details

- **Algorithm**: AES (Advanced Encryption Standard)
- **Padding**: PKCS7
- **Current Key**: Uses a hardcoded cipher key (not user-configurable)
- **Future Improvements**: Planned features include user-set master passwords and key derivation functions (KDF) like
  PBKDF2 or Argon2.

## Project Structure

```
pwm/
├── main.go                 # Entry point
├── go.mod                  # Go module definition
├── go.sum                  # Dependency checksums
├── Makefile               # Build and install scripts
├── README.md              # This file
├── LICENSE                # MIT License
├── bin/
│   └── pwm                # Compiled binary
├── cmd/
│   ├── root.go            # Root command and main logic
│   ├── create.go          # Create command
│   ├── ls.go              # List command
│   ├── cp.go              # Copy command
│   ├── rm.go              # Remove command
│   └── pwm.go             # Core command implementations
└── internal/
    ├── crypto/
    │   └── crypto.go      # AES encryption/decryption logic
    ├── helpers/
    │   └── inputs.go      # User input handling
    └── models/
        └── secret.go      # Secret data model and serialization
```

## Development

### Dependencies

- **github.com/spf13/cobra**: CLI framework
- **github.com/SpyrosMoux/passwdgen**: Password generation
- **golang.design/x/clipboard**: Clipboard integration
- **golang.org/x/term**: Terminal handling

Install dependencies:

```bash
go mod tidy
```

### Building

```bash
go build -o bin/pwm main.go
```

### Running Tests

```bash
go test ./...
```

## License

This project is licensed under the GNU License - see the [LICENSE](LICENSE) file for details.

## Author

Created by Spyros Mouchlianitis

## Contributing

Contributions are welcome! Feel free to open an issue or submit a pull request.
