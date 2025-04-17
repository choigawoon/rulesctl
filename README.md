# rulesctl

Cursor Rules management tool using GitHub Gist

## Problem → Solution

We needed a tool to effectively manage and share Cursor Rules. rulesctl allows you to systematically store and manage rules using GitHub Gist.

## Usage

### Installation

You can install via NPM:

```bash
npm install -g rulesctl
```

### Authentication Setup (Required)

You **must** set up GitHub authentication before using rulesctl. There are two ways to set up authentication:

1. Using environment variable (recommended)
```bash
export GITHUB_TOKEN="your_github_token"
# No need to run 'rulesctl auth' when using environment variable
```

2. Using the auth command
```bash
rulesctl auth  # Enter token at prompt
```

The token will be stored in JSON format in the `~/.rulesctl/config` file in your home directory when using the auth command.

> **Important**: Your Personal Access Token needs the following permissions:
> - Gist (read/write) permissions
> - repo permissions (for accessing file lists at https://github.com/PatrickJS/awesome-cursorrules/tree/main/rules-new)
![1](docs/images/how-to-get-token-1.png)
![2](docs/images/how-to-get-token-2.png)

For information on how to create a token, refer to the [GitHub official documentation](https://docs.github.com/en/authentication/keeping-your-account-and-data-secure/creating-a-personal-access-token).

### Basic Commands

```bash
# View help
rulesctl --help

# Create rules directory
rulesctl init
rulesctl init --sample  # Also create example rule files

# View rule list (only shows those from the last month)
rulesctl list                # Show basic information
rulesctl list --detail      # Show detailed information including revision

# Upload rules (private by default)
rulesctl upload "RuleSetName"
rulesctl upload "RuleSetName" --public  # Upload as public Gist

# Download rules
rulesctl download "RuleSetName"         # Search by title in my Gist
rulesctl download --gistid abc123       # Download by public Gist ID

# Delete rules
rulesctl delete "RuleSetName"           # Search and delete by title
rulesctl delete "RuleSetName" --force   # Delete immediately without confirmation
```

### Usage Examples

First, set up authentication:
```bash
# Authenticate with GitHub token
rulesctl auth
# Enter Personal Access Token at prompt
```

Creating a rules directory:
```bash
# Create .cursor/rules directory
rulesctl init

# Start with examples
rulesctl init --sample
```

Uploading a rule set:
```bash
# Upload rules from current directory with a specific name (private by default)
rulesctl upload "my-python-ruleset"

# Upload with a name and description
rulesctl upload "my-python-ruleset" --desc "Rule collection for Python projects"

# Public upload for sharing with others
rulesctl upload "my-python-ruleset" --public

# Force upload with duplicate name (no confirmation prompt)
rulesctl upload "my-python-ruleset" --force
```

Deleting a rule set:
```bash
# Search and delete by title
rulesctl delete "my-python-ruleset"

# Delete immediately without confirmation
rulesctl delete "my-python-ruleset" --force
```

Checking the rules list:
```bash
# View list of rules uploaded in the last month
rulesctl list
```

Downloading a rule set:
```bash
# Download by searching for title in my Gist
rulesctl download "my-python-ruleset"

# Download directly by ID from a public Gist
rulesctl download --gistid abc123

# Force download even with conflicts
rulesctl download "my-python-ruleset" --force
rulesctl download --gistid abc123 --force
```

Deleting a rule set:
```bash
# Search and delete by title
rulesctl delete "my-python-ruleset"

# Delete immediately without confirmation
rulesctl delete "my-python-ruleset" --force
```

> **Important**:
> - Download supports two methods:
>   1. Download by title: Finds rules with exactly matching title from your Gist list
>   2. Download by Gist ID: Directly download by specifying the ID of a public Gist
> - If the `.cursor/rules` directory doesn't exist in the current path during download, it's created automatically
> - The original directory structure and files are restored exactly as they were uploaded
> - Files are ready to use immediately after download

## Supported Platforms

rulesctl supports the following platforms:
- macOS (darwin)
- Linux
- Windows

## Developer Guide

For development and testing methods, refer to the [Getting Started Guide](docs/en/GET-STARTED.md).

## Contributing

For contribution guidelines, refer to the [Contribution Guide](docs/en/GET-STARTED.md#contributing). 