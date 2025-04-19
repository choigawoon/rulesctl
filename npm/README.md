# rulesctl

> Cursor Rules management tool using GitHub Gist

[![npm version](https://badge.fury.io/js/rulesctl.svg)](https://badge.fury.io/js/rulesctl)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

`rulesctl` is a command-line tool that helps you manage and share Cursor Rules using GitHub Gist. It provides an easy way to upload, download, and manage your Cursor Rules across different machines.

## Features

- 📤 Upload rules to GitHub Gist (private by default)
- 📥 Download rules by name or Gist ID
- 📋 List your uploaded rules
- 🔄 Easy synchronization across machines
- 🔒 Secure authentication using GitHub token

## Installation

```bash
npm install -g rulesctl
```

## Authentication

Before using rulesctl, you need to set up GitHub authentication. There are two ways:

1. Using environment variable (recommended):
```bash
export GITHUB_TOKEN="your_github_token"
```

2. Using the auth command:
```bash
rulesctl auth  # Enter token at prompt
```

> **Note**: Your GitHub token needs Gist (read/write) and repo permissions.

## Usage

### Basic Commands

```bash
# Create rules directory
rulesctl init
rulesctl init --sample  # Create with example rules

# List your rules
rulesctl list
rulesctl list --detail  # Show detailed information

# Upload rules
rulesctl upload "MyRules"
rulesctl upload "MyRules" --public  # Upload as public Gist

# Download rules
rulesctl download "MyRules"         # Search by title
rulesctl download --gistid abc123   # Download by Gist ID

# Delete rules
rulesctl delete "MyRules"
rulesctl delete "MyRules" --force   # Skip confirmation
```

### Examples

```bash
# Initialize with example rules
rulesctl init --sample

# Upload rules with description
rulesctl upload "python-rules" --desc "My Python project rules"

# Share rules publicly
rulesctl upload "shared-rules" --public

# Download specific rules
rulesctl download "python-rules"

# Download from public Gist
rulesctl download --gistid abc123def456
```

## Supported Platforms

- macOS (darwin)
- Linux
- Windows

## Requirements

- Node.js >= 14
- GitHub account with personal access token

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request. For major changes, please open an issue first to discuss what you would like to change.

## License

MIT © [choigawoon](https://github.com/choigawoon)

---

For more detailed information and documentation, visit our [GitHub repository](https://github.com/choigawoon/rulesctl). 