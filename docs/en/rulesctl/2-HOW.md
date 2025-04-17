# rulesctl Implementation Guide

## Tech Stack
- **Go**: Main development language
- **Cobra**: CLI framework
- **GitHub API v3**: GIST management
- **JSON**: Configuration and metadata storage
- **MD5**: File integrity verification

## Architecture Design
```
project-root/
├── cmd/
│   ├── root.go       # Root command setup
│   ├── auth.go       # GitHub authentication handling
│   ├── list.go       # GIST list output
│   ├── upload.go     # Rule upload
│   └── download.go   # Rule download
├── internal/
│   ├── gist/         # GIST API wrapper
│   └── fileutils/    # File system utilities
└── pkg/
    └── config/       # Configuration file management
```

## Core Command Implementation

### 1. Authentication Handling (`auth`)
```go
// cmd/auth.go
var authCmd = &cobra.Command{
    Use:   "auth",
    Short: "Set GitHub Personal Access Token",
    RunE: func(cmd *cobra.Command, args []string) error {
        token, _ := cmd.Flags().GetString("token")
        return config.SaveToken(token)
    },
}
```

Authentication information is stored in `~/.rulesctl/config` file:
```json
{
  "token": "ghp_YourPersonalAccessTokenHere",
  "last_used": "2023-08-15T12:34:56Z"
}
```

> **Important**: Personal Access Token requires the following permissions:
> - Gist (read/write) permission
> - repo permission (for accessing file list at https://github.com/PatrickJS/awesome-cursorrules/tree/main/rules-new)

### 2. List Rules (`list`)
```go
// cmd/list.go
var listCmd = &cobra.Command{
    Use:   "list",
    Short: "Display list of rules stored in GIST",
    RunE: func(cmd *cobra.Command, args []string) error {
        gists, err := gist.FetchUserGists()
        // Sort and display in [last_modified] title format
    },
}
```

### 3. Upload Rules (`upload`)
```go
// cmd/upload.go
var uploadCmd = &cobra.Command{
    Use:   "upload [name]",
    Short: "Upload local rules to GIST",
    RunE: func(cmd *cobra.Command, args []string) error {
        return fileutils.WalkRulesDir(".cursor/rules", func(path string) {
            gist.AddFile(path, content)
        })
    },
}
```

> **Important**: 
> - rulesctl requires `.cursor/rules/**/*.mdc` structure in the current execution path.
> - Rule set names should be enclosed in quotes.

### 4. Download Rules (`download`)
```go
// cmd/download.go
var downloadCmd = &cobra.Command{
    Use:   "download [name]",
    Short: "Download rules from GIST",
    RunE: func(cmd *cobra.Command, args []string) error {
        if !force && checkConflicts() {
            return errors.New("Conflict files exist. Use --force option")
        }
        return gist.DownloadFiles(args[0])
    },
}
```

> **Important**:
> - If `.cursor/rules` directory doesn't exist during download, it will be created automatically.
> - The original uploaded directory structure and files are restored as is.

## File Conflict Check Logic
```go
func checkConflicts(gistID string) bool {
    localFiles := fileutils.ListLocalRules()
    remoteFiles := gist.GetFileList(gistID)
    
    for f := range remoteFiles {
        if _, exists := localFiles[f]; exists {
            return true
        }
    }
    return false
}
```

## Path Structure Example
```
.cursor/
└── rules/
    ├── python/
    │   ├── linting.mdc
    │   └── testing.mdc
    └── database/
        └── postgres.mdc
```

## GIST Structured Storage Method
```
gist/
├── python_linting.mdc
├── python_testing.mdc
├── database_postgres.mdc
└── meta.json  # Directory structure and file metadata
```

`meta.json` file structure:
```json
{
  "schema_version": "1.0.0",
  "cli_version": "0.1.0",
  "updated_at": "2024-03-17T12:34:56Z",
  "structure": {
    "python": {
      "linting.mdc": {
        "path": "python/linting.mdc",
        "gist_name": "python_linting.mdc",
        "size": 1234,
        "md5": "a1b2c3d4e5f6g7h8i9j0"
      },
      "testing.mdc": {
        "path": "python/testing.mdc",
        "gist_name": "python_testing.mdc",
        "size": 2345,
        "md5": "b2c3d4e5f6g7h8i9j0a1"
      }
    },
    "database": {
      "postgres.mdc": {
        "path": "database/postgres.mdc",
        "gist_name": "database_postgres.mdc",
        "size": 3456,
        "md5": "c3d4e5f6g7h8i9j0a1b2"
      }
    }
  }
}
```

> **Important**: 
> - rulesctl requires `.cursor/rules/**/*.mdc` structure in the current execution path.
> - Rule set names should be enclosed in quotes.
> - When uploading to Gist, file names are converted to reflect the directory structure. (e.g., `python/linting.mdc` → `python_linting.mdc`)

## NPM Deployment

For NPM deployment methods, refer to the [NPM Deployment Guide](../npm/2-HOW.md).

This implementation combines Cobra's subcommand system with the GitHub API client to allow users to systematically manage cursorrules through CLI. In particular, it enhances collaboration efficiency in team development environments through directory structure maintenance and conflict detection features.

Citations:
[1] https://docs.cursor.com/context/rules
[2] https://github.com/spf13/cobra
[3] https://apidog.com/blog/awesome-cursor-rules/
[4] https://www.bytesizego.com/blog/structure-go-cli-app
[5] https://github.com/spf13/cobra-cli/blob/main/README.md
[6] https://github.com/Qwertic/cursorrules
[7] https://www.digitalocean.com/community/tutorials/how-to-use-the-cobra-package-in-go
[8] https://dev.to/kgoedert/create-a-command-line-tool-with-go-and-cobra-eel
[9] https://www.jetbrains.com/guide/go/tutorials/cli-apps-go-cobra/creating_cli/
[10] https://www.sktenterprise.com/bizInsight/blogDetail/dev/2755
[11] https://www.reddit.com/r/golang/comments/16igrq6/a_better_cobra_project_structure/
[12] https://apidog.com/kr/blog/awesome-cursor-rules-kr/
[13] https://stackoverflow.com/questions/25161774/what-are-conventions-for-filenames-in-go
[14] https://nangman14.tistory.com/97
[15] https://github.com/sanjeed5/awesome-cursor-rules-mdc
[16] https://dev.to/tuhinbar/my-first-cli-with-go-4eig
[17] https://mcauto.github.io/back-end/2018/10/30/golang-cobra/
[18] https://forum.cursor.com/t/how-to-force-your-cursor-ai-agent-to-always-follow-your-rules-using-auto-rule-generation-techniques/80199
[19] https://github.com/golang-standards/project-layout
[20] https://forum.cursor.com/t/prompting-the-perfect-coding-partner-through-cursorrules/39907
[21] https://google.github.io/styleguide/go/best-practices.html

---
Response from Perplexity: pplx.ai/share 