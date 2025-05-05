# rulesctl Development Roadmap

## Task Criteria
This document breaks down each task phase to be completable within **1 hour**. This serves as a standard for measuring and managing task progress.

---

# Milestone 1: Basic Feature Implementation

## 0. Planning and Implementation Method Review (2 hours)

### 0-1. Review and Organize Planning Documents (1 hour)
- [x] Review GIST structured storage method - Manage each rule set as separate GIST
- [x] Review meta.json schema - Include only basic file structure and MD5 hash
- [x] Review file conflict prevention strategy - Simplify with force/skip options

### 0-2. Document Implementation Methods (1 hour)
- [x] Write implementation guide - Define CLI structure and core commands
- [x] Write API design document - Define GitHub Gist API integration method
- [x] Write data structure document - Define meta.json schema and file structure
- [x] Document error handling strategy - Define force/skip option handling method

## 1. Initial Project Setup (2 hours)

### 1-1. Development Environment Setup (1 hour)
- [x] Initialize Go project (go.mod, go.sum)
- [x] Install Cobra package and set up basic CLI structure
- [x] Create directory structure (cmd, internal, pkg, etc.)
- [x] Configure .gitignore

### 1-2. Basic Command Structure Setup (1 hour)
- [x] Write root.go (basic command setup)
- [x] Set up basic flags and global options
- [x] Implement version info command
- [x] Run tests and verify basic functionality
- [x] Improve init command - Separate basic directory creation and sample file creation (add --sample flag)

## 2. Configuration and Authentication (3 hours)

### 2-1. Configuration File Management (1 hour)
- [x] Implement configuration directory creation logic (~/.rulesctl)
- [x] Implement configuration file save and load functionality
- [x] Write tests

### 2-2. GitHub API Client Implementation (1 hour)
- [x] Implement GitHub API client basic structure
- [x] Implement token authentication mechanism
- [x] Write API request basic wrapper functions

### 2-3. Implement auth Command (1 hour)
- [x] Implement auth command (cmd/auth.go)
- [x] Token input and validation logic
- [x] Token storage functionality
- [x] Test and integration verification
- [x] Add token input masking feature
- [x] Add test token management feature through .env.local file

## 3. Gist Core Features (5 hours)

### 3-1. Gist List View Feature (1 hour)
- [x] Implement list command (cmd/list.go)
- [x] Integrate Gist API and list view
- [x] Implement output format (table, JSON, etc.)
- [x] Write tests
- [x] Configure test environment through .env.local file
- [x] Improve output format (fixed width, ellipsis handling)
- [x] Add Gist ID display
- [x] Modify to show only Gists within the last month
- [x] Filter only Gists with .rulesctl.meta.json file

### 3-2. File System Utilities (1 hour)
- [x] Implement .cursor/rules directory management functions
- [x] File exploration and listing functionality
- [x] Path handling utility functions
- [x] Write tests

### 3-3. Rule Upload Feature (1 hour)
- [x] Implement upload command (cmd/upload.go)
- [x] Local rule file collection logic
- [x] Integrate Gist creation/update API
- [x] Write tests
- [x] Add metadata preview feature (--preview)
- [x] Display appropriate guidance message when no files exist

### 3-4. Rule Download Feature (1 hour)
- [x] Implement download command (cmd/download.go)
- [x] File download logic from Gist
- [x] Maintain local file structure functionality
- [x] Conflict handling logic
- [x] Write tests
- [x] Add download by title feature
- [x] Add download by Gist ID feature (support public Gists)

### 3-5. Rule Delete Feature (1 hour)
- [x] Implement delete command (cmd/delete.go)
- [x] Integrate Gist delete API
- [x] Implement confirmation prompt
- [x] Implement delete by title search feature
- [x] Implement delete without confirmation using --force option
- [x] Write tests

## 4. User Experience Improvements (3 hours)

### 4-1. Error Handling and Logging (1 hour)
- [x] Implement consistent error handling mechanism
- [x] Implement logging system
- [x] User-friendly error messages
- [x] Implement debug mode
- [x] Translate all error messages to English for global distribution

### 4-2. Progress Display and Prompts (1 hour)
- [ ] Progress display (progress bar, etc.)
- [ ] Implement interactive prompts
- [ ] Support color output
- [ ] Improve table format output

### 4-3. Usage and Help Improvements (1 hour)
- [ ] Improve help for each command
- [ ] Add examples and usage
- [ ] Generate auto-completion scripts
- [ ] Support man pages

## 5. NPM Packaging and Deployment (4 hours)

### 5-1. Go Binary Build Scripts (1 hour)
- [ ] Write cross-platform build scripts
- [ ] Support various architectures (amd64, arm64)
- [ ] Build automation scripts

### 5-2. NPM Package Structure Setup (1 hour)
- [ ] Create npm directory structure
- [ ] Write package.json
- [ ] Write launcher script (rulesctl.js)
- [ ] Configure binary location

### 5-3. Package Testing and Deployment Preparation (1 hour)
- [ ] Local package testing (npm link)
- [ ] Test on various platforms
- [ ] Create deployment checklist

### 5-4. CI/CD and Automatic Deployment Setup (1 hour)
- [ ] Write GitHub Actions workflow
- [ ] Automatic tag creation and release
- [ ] Configure NPM automatic deployment
- [ ] Version management automation

## 6. Documentation and Finalization (5 hours)

### 6-1. User Documentation (1 hour)
- [x] Complete README.md
- [x] Document installation and usage methods
- [x] Add examples and screenshots

### 6-2. Developer Documentation (1 hour)
- [x] Complete developer guide (GET-STARTED.md)
- [x] Write API documentation
- [x] Write contribution guidelines
- [x] Add test environment setup documentation

### 6-3. Document Internationalization (2 hours)
- [x] Restructure documentation directory
- [x] Set English version as default
- [x] Create Korean version documents (_ko.md)
- [x] Update cross-reference links between documents
- [x] Write multilingual document management guidelines

### 6-4. Testing and Final Check (1 hour)
- [ ] Complete feature integration testing
- [ ] Usability testing
- [ ] Final release preparation check

---

# Milestone 2 Plan (Future Work)

## Features Under Consideration
- Rule search functionality
- Rule sharing and team collaboration features
- Rule version management features
- Web interface development
- Rule templates and scaffolding features
- Gist revision specific download feature
  - Feature to download rule sets of specific revisions
  - Specify specific version through `--revision` flag
  - Add revision history view feature

## Current Issues

### Version Management Improvement
1. Version Information Centralization
   - [x] Created internal/version/version.go package
   - [x] Updated goreleaser ldflags configuration
   - [ ] NPM package version synchronization issue
     - Current situation: Version shows as 0.1.0 when installing locally with npm install -g ./rulesctl-0.1.2.tgz
     - Causes: 
       1. Binary built by goreleaser differs from npm package binary
       2. VERSION constant in npm/rulesctl.js doesn't match actual binary version
       3. GitHub Release API caching issue
          - New version downloads correctly only after some time has passed since GitHub release deployment
          - API responses are likely cached and not immediately reflected
     - Potential solutions:
       1. Copy goreleaser build binaries to npm/bin before npm pack
       2. Download binaries from GitHub releases even for local installation
       3. Review goreleaser post hooks modification
       4. Wait for a period (about 5-10 minutes) after GitHub Release deployment before proceeding with npm deployment

### Internationalization and Localization
1. Message Translation
   - [x] Translate all CLI messages to English
   - [x] Update client.go messages
   - [x] Update list.go messages
   - [x] Update download.go messages
   - [ ] Review and update remaining files for any Korean messages
   - [ ] Implement i18n system for future multi-language support

### Next Steps
1. Resolve Version Management Issue
   - [ ] Decide npm package installation method (local binary vs GitHub releases)
   - [ ] Modify rulesctl.js
   - [ ] Test installation process
   - [ ] Update DEPLOY.md documentation

2. Public Rules Store Implementation
   - [x] Implement `store` subcommand structure
   - [x] Create `public-store.json` with metadata (name, description, gistid, source, category)
   - [x] Implement `store list` command to fetch and display templates with full Gist ID
   - [x] Implement `store download` command to download rules by name
   - [x] Update documentation for new store commands
   - [x] Add redirects from deprecated `list --store` to new commands
   - [x] Set up GitHub Actions workflow for community contribution (auto-merge on thumbs up)
   - [ ] Add filtering by category
   - [ ] Add search functionality

3. Remaining Tasks
   - [ ] Improve error handling and logging
   - [ ] Enhance progress display
   - [ ] Improve usage and help documentation
   - [ ] Set up GitHub Actions workflow

## Future Features (2nd Milestone)
- Rule search functionality
- Rule sharing and team collaboration features
- Rule version management
- Web interface development
- Rule templates and scaffolding features
- Gist revision specific download feature
  - Ability to download specific revision of rule sets
  - Use `--revision` flag to specify version
  - Add revision history viewing functionality 