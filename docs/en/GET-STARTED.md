# rulesctl Developer Guide

## Project Structure

```
rulesctl/
├── cmd/                    # Cobra command definitions
│   ├── root.go
│   ├── auth.go
│   ├── upload.go
│   └── ...
├── internal/               # Internal logic modules
│   ├── gist/
│   └── fileutils/
├── pkg/                    # Externally exposable packages (e.g., config)
│   └── config/
├── go.mod
└── main.go
```

## Development Environment Setup

1. Install Go
   - Install Go version 1.21 or higher
   - Download from [official download page](https://golang.org/dl/)

2. Clone the Project
```bash
git clone https://github.com/choigawoon/rulesctl.git
cd rulesctl
```

3. Install Dependencies
```bash
go mod download
```

## Implementation Guide

For rulesctl implementation and NPM deployment methods, refer to the following documents:
- [rulesctl Implementation Guide](rulesctl/2-HOW.md)
- [NPM Deployment Guide](npm/2-HOW.md)

## Testing

### Test Environment Setup

1. Create `.env.local` file
```bash
echo "GITHUB_PERSONAL_ACCESS_TOKEN=your_github_token" > .env.local
```

2. GitHub Token Setup
   - GitHub Settings → Developer settings → Personal access tokens → Tokens (classic)
   - Required permissions: `gist` (read/write)

### How to Run Tests

1. Test Specific Package
```bash
# Test gist package
go test ./internal/gist

# Test config package
go test ./pkg/config

# Test cmd package
go test ./cmd
```

2. Test All Packages
```bash
go test ./...
```

3. View Detailed Test Results
```bash
go test -v ./...
```

4. Check Test Coverage
```bash
go test -cover ./...
```

### Testing Precautions

1. The `.env.local` file is added to `.gitignore` to prevent committing to Git.
2. The test token needs permissions to access Gists on the actual GitHub account.
3. A valid token must be set in the `.env.local` file before running tests.

## Contributing

To contribute to the rulesctl project, follow these steps:

1. Fork this repository
2. Create a new branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Create a Pull Request

### PR Guidelines

- All PRs must include tests
- Documentation updates should be included when necessary
- Code style should follow `gofmt`
- Commit messages should follow [Conventional Commits](https://www.conventionalcommits.org/)

## License

This project is distributed under the MIT License. See [LICENSE](LICENSE) file for more information. 