# Deployment Guide

This document describes the process of deploying a new version of rulesctl.

## 1. Version Update

Version information is managed in the following files:

- `internal/version/version.go`: Central version management for Go code
- `npm/rulesctl.js`: Version information for NPM package

### How to Update Version

1. Update `Version` variable in `internal/version/version.go`:
```go
var Version = "x.y.z"  // Change to new version
```

2. Update `VERSION` constant in `npm/rulesctl.js`:
```javascript
const VERSION = "vx.y.z";  // Change to new version (v prefix required)
```

## 2. Create Release

1. Commit changes:
```bash
git add .
git commit -m "chore: bump version to x.y.z"
```

2. Create and push tag:
```bash
git tag -a vx.y.z -m "Release vx.y.z"
git push origin vx.y.z
```

3. Create release with goreleaser:
```bash
goreleaser release --clean
```

4. Check and Publish Release on GitHub
- Check the draft release at https://github.com/choigawoon/rulesctl/releases
- Click "Publish release" button to make it public

## 3. Deploy NPM Package

1. Create and test package:
```bash
cd npm
npm pack
npm install -g ./rulesctl-x.y.z.tgz  # Local test
rulesctl version  # Version check
```

2. NPM login and publish:
```bash
npm login  # Login with NPM account
npm publish  # Publish package
```

## 4. Post-deployment Verification

1. Test NPM package installation:
```bash
npm install -g rulesctl
rulesctl version  # Check new version
```

2. Test basic functionality:
```bash
rulesctl init --sample
rulesctl list
rulesctl upload "test"
rulesctl download "test"
rulesctl delete "test"
```

## Important Notes

1. Version Management
- Version must be managed consistently in both `internal/version/version.go` and `npm/rulesctl.js`
- goreleaser updates version information at build time through ldflags

2. GitHub Token
- goreleaser uses the `GITHUB_TOKEN` environment variable
- Ensure token has `repo` permissions

3. NPM Deployment
- Verify version in `package.json` is correct
- Always test locally before publishing 