# NPM Deployment Guide

## Deployment Structure

```
rulesctl/
├── npm/                    # NPM packaging files
│   ├── bin/               # Platform-specific binaries
│   │   ├── rulesctl-darwin
│   │   ├── rulesctl-linux
│   │   └── rulesctl-win.exe
│   ├── rulesctl.js        # Launcher script
│   └── package.json       # NPM metadata
```

## Package Configuration

### 1. package.json
```json
{
  "name": "rulesctl",
  "version": "1.0.0",
  "description": "Cursor Rules Management CLI Tool",
  "bin": {
    "rulesctl": "./rulesctl.js"
  },
  "files": [
    "bin/",
    "rulesctl.js"
  ],
  "os": [
    "darwin",
    "linux",
    "win32"
  ],
  "cpu": [
    "x64",
    "arm64"
  ]
}
```

### 2. rulesctl.js (Launcher Script)
```javascript
#!/usr/bin/env node

const os = require("os");
const path = require("path");
const { spawn } = require("child_process");

const platform = os.platform();
let bin = "rulesctl-";

if (platform === "darwin") bin += "darwin";
else if (platform === "linux") bin += "linux";
else if (platform === "win32") bin += "win.exe";
else {
  console.error("Unsupported platform:", platform);
  process.exit(1);
}

const binPath = path.join(__dirname, "bin", bin);
const args = process.argv.slice(2);

spawn(binPath, args, { stdio: "inherit" })
  .on("exit", code => process.exit(code));
```

## Build and Deployment Process

1. **Build Go Binaries**
```bash
# macOS
GOOS=darwin GOARCH=amd64 go build -o npm/bin/rulesctl-darwin
GOOS=darwin GOARCH=arm64 go build -o npm/bin/rulesctl-darwin-arm64

# Linux
GOOS=linux GOARCH=amd64 go build -o npm/bin/rulesctl-linux
GOOS=linux GOARCH=arm64 go build -o npm/bin/rulesctl-linux-arm64

# Windows
GOOS=windows GOARCH=amd64 go build -o npm/bin/rulesctl-win.exe
```

2. **Package Deployment**
```bash
# Deploy NPM package
cd npm
npm publish
```

## Version Management

1. **Follow Semantic Versioning**
   - MAJOR: Incompatible API changes
   - MINOR: Backward-compatible feature additions
   - PATCH: Backward-compatible bug fixes

2. **Version Update Process**
```bash
# Update package.json version
npm version [major|minor|patch]

# Commit changes and create tag
git commit -am "v1.0.0"
git tag v1.0.0
git push origin main --tags
```

## Testing Process

1. **Local Testing**
```bash
# Link package locally
cd npm
npm link

# Test in another terminal
rulesctl --version
```

2. **CI/CD Pipeline**
```yaml
# GitHub Actions example
name: Publish
on:
  push:
    tags:
      - 'v*'

jobs:
  publish:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v2
      - uses: actions/setup-node@v2
      - run: npm ci
      - run: npm publish
```

## Security Considerations

1. **Binary Signing**
   - Apply digital signatures to all binaries
   - Verify signatures on user systems

2. **Dependency Checks**
   - Regular security vulnerability checks
   - Automated dependency updates

3. **Token Management**
   - Secure NPM token storage
   - Safe token usage in CI/CD pipeline 