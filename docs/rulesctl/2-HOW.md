# rulesctl 구현 가이드

## 기술 스택
- **Go**: 메인 개발 언어
- **Cobra**: CLI 프레임워크
- **GitHub API v3**: GIST 관리
- **JSON**: 설정 및 메타데이터 저장
- **MD5**: 파일 무결성 검증

## 아키텍처 설계
```
project-root/
├── cmd/
│   ├── root.go       # 루트 커맨드 설정
│   ├── auth.go       # GitHub 인증 처리
│   ├── list.go       # GIST 목록 출력
│   ├── upload.go     # 규칙 업로드
│   └── download.go   # 규칙 다운로드
├── internal/
│   ├── gist/         # GIST API 래퍼
│   └── fileutils/    # 파일 시스템 유틸리티
└── pkg/
    └── config/       # 설정 파일 관리
```

## 핵심 명령어 구현

### 1. 인증 처리 (`auth`)
```go
// cmd/auth.go
var authCmd = &cobra.Command{
    Use:   "auth",
    Short: "GitHub Personal Access Token 설정",
    RunE: func(cmd *cobra.Command, args []string) error {
        token, _ := cmd.Flags().GetString("token")
        return config.SaveToken(token)
    },
}
```

인증 정보는 `~/.rulesctl/config` 파일에 저장됩니다:
```json
{
  "token": "ghp_YourPersonalAccessTokenHere",
  "last_used": "2023-08-15T12:34:56Z"
}
```

> **중요**: Personal Access Token에는 다음 권한이 필요합니다:
> - Gist (read/write) 권한
> - repo 권한 (https://github.com/PatrickJS/awesome-cursorrules/tree/main/rules-new 에서 파일 목록 접근용)

### 2. 규칙 목록 조회 (`list`)
```go
// cmd/list.go
var listCmd = &cobra.Command{
    Use:   "list",
    Short: "GIST에 저장된 규칙 목록 출력",
    RunE: func(cmd *cobra.Command, args []string) error {
        gists, err := gist.FetchUserGists()
        // [최종수정시각] 제목 형식으로 정렬 출력
    },
}
```

### 3. 규칙 업로드 (`upload`)
```go
// cmd/upload.go
var uploadCmd = &cobra.Command{
    Use:   "upload [name]",
    Short: "로컬 규칙을 GIST에 업로드",
    RunE: func(cmd *cobra.Command, args []string) error {
        return fileutils.WalkRulesDir(".cursor/rules", func(path string) {
            gist.AddFile(path, content)
        })
    },
}
```

> **중요**: 
> - rulesctl은 현재 실행 경로에 `.cursor/rules/**/*.mdc` 구조가 있어야만 사용할 수 있습니다.
> - 규칙 세트 이름은 따옴표로 감싸서 지정합니다.

### 4. 규칙 다운로드 (`download`)
```go
// cmd/download.go
var downloadCmd = &cobra.Command{
    Use:   "download [name]",
    Short: "GIST에서 규칙 다운로드",
    RunE: func(cmd *cobra.Command, args []string) error {
        if !force && checkConflicts() {
            return errors.New("충돌 파일 존재. --force 옵션 사용")
        }
        return gist.DownloadFiles(args[0])
    },
}
```

> **중요**:
> - 다운로드 시 현재 실행 경로에 `.cursor/rules` 디렉토리가 없으면 자동으로 생성됩니다.
> - 원래 업로드된 디렉토리 구조와 파일들이 그대로 복원됩니다.

## 파일 충돌 검사 로직
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

## 경로 구조 예시
```
.cursor/
└── rules/
    ├── python/
    │   ├── linting.mdc
    │   └── testing.mdc
    └── database/
        └── postgres.mdc
```

## GIST 구조화 저장 방식
```
gist/
├── python_linting.mdc
├── python_testing.mdc
├── database_postgres.mdc
└── meta.json  # 디렉토리 구조 및 파일 메타데이터
```

`meta.json` 파일 구조:
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

> **중요**: 
> - rulesctl은 현재 실행 경로에 `.cursor/rules/**/*.mdc` 구조가 있어야만 사용할 수 있습니다.
> - 규칙 세트 이름은 따옴표로 감싸서 지정합니다.
> - Gist에 업로드될 때 파일 이름은 디렉토리 구조를 반영하여 변환됩니다. (예: `python/linting.mdc` → `python_linting.mdc`)

## NPM 배포

NPM을 통한 배포 방법은 [NPM 배포 가이드](../npm/2-HOW.md)를 참조하세요.

이 구현체는 Cobra의 서브커맨드 체계와 GitHub API 클라이언트를 결합하여, 사용자가 CLI를 통해 cursorrules를 체계적으로 관리할 수 있도록 설계되었습니다. 특히 디렉토리 구조 유지와 충돌 검사 기능을 통해 팀 개발 환경에서의 협업 효율성을 높였습니다.

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
Perplexity로부터의 답변: pplx.ai/share