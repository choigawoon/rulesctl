# rulesctl 개발자 가이드

## 프로젝트 구조

```
rulesctl/
├── cmd/                    # Cobra 명령어 정의
│   ├── root.go
│   ├── auth.go
│   ├── upload.go
│   └── ...
├── internal/               # 내부 로직 모듈
│   ├── gist/
│   └── fileutils/
├── pkg/                    # 설정 등 외부 노출 가능 패키지
│   └── config/
├── go.mod
└── main.go
```

## 개발 환경 설정

1. Go 설치
   - Go 1.21 이상 버전 설치
   - [공식 다운로드 페이지](https://golang.org/dl/)에서 다운로드

2. 프로젝트 클론
```bash
git clone https://github.com/choigawoon/rulesctl.git
cd rulesctl
```

3. 의존성 설치
```bash
go mod download
```

## 구현 가이드

rulesctl의 구현 방법과 NPM 배포 방법은 다음 문서를 참조하세요:
- [rulesctl 구현 가이드](rulesctl/2-HOW.md)
- [NPM 배포 가이드](npm/2-HOW.md)

## 테스트

### 테스트 환경 설정

1. `.env.local` 파일 생성
```bash
echo "GITHUB_PERSONAL_ACCESS_TOKEN=your_github_token" > .env.local
```

2. GitHub 토큰 설정
   - GitHub 설정 → Developer settings → Personal access tokens → Tokens (classic)
   - 필요한 권한: `gist` (read/write)

### 테스트 실행 방법

1. 특정 패키지 테스트
```bash
# gist 패키지 테스트
go test ./internal/gist

# config 패키지 테스트
go test ./pkg/config

# cmd 패키지 테스트
go test ./cmd
```

2. 모든 패키지 테스트
```bash
go test ./...
```

3. 상세 테스트 결과 확인
```bash
go test -v ./...
```

4. 테스트 커버리지 확인
```bash
go test -cover ./...
```

### 테스트 주의사항

1. `.env.local` 파일은 Git에 커밋되지 않도록 `.gitignore`에 추가되어 있습니다.
2. 테스트 토큰은 실제 GitHub 계정의 Gist에 접근할 수 있는 권한이 필요합니다.
3. 테스트 실행 전에 `.env.local` 파일에 유효한 토큰이 설정되어 있어야 합니다.

## 기여하기

rulesctl 프로젝트에 기여하고 싶다면 다음 단계를 따르세요:

1. 이 저장소를 포크합니다.
2. 새로운 브랜치를 생성합니다 (`git checkout -b feature/amazing-feature`).
3. 변경사항을 커밋합니다 (`git commit -m 'Add some amazing feature'`).
4. 브랜치를 푸시합니다 (`git push origin feature/amazing-feature`).
5. Pull Request를 생성합니다.

### PR 가이드라인

- 모든 PR은 테스트를 포함해야 합니다.
- 문서 업데이트가 필요한 경우 함께 포함해야 합니다.
- 코드 스타일은 `gofmt`를 따릅니다.
- 커밋 메시지는 [Conventional Commits](https://www.conventionalcommits.org/)를 따릅니다.

## 라이선스

이 프로젝트는 MIT 라이선스 하에 배포됩니다. 자세한 내용은 [LICENSE](LICENSE) 파일을 참조하세요.

## 빌드 및 배포

### goreleaser를 사용한 릴리즈

1. goreleaser 설정
```bash
# goreleaser 설치
go install github.com/goreleaser/goreleaser@latest

# .goreleaser.yaml 파일 생성 및 설정
# - 크로스 플랫폼 빌드 설정
# - 체인지로그 자동화
# - 릴리즈 노트 템플릿
```

2. 로컬 빌드 테스트
```bash
# 스냅샷 빌드 (태그 없이 테스트)
~/go/bin/goreleaser build --snapshot --clean

# 결과물 확인
ls -l bin/           # 최신 빌드 결과물
ls -l dist/          # 모든 아키텍처별 빌드 결과물
```

3. GitHub 릴리즈 생성
```bash
# 1. GitHub 토큰 설정
export GITHUB_TOKEN="your_token"

# 2. 버전 태그 생성
git tag -a v0.1.0 -m "First release

주요 기능:
- GitHub Gist를 이용한 Cursor Rules 관리
- 규칙 업로드/다운로드
- 크로스 플랫폼 지원 (macOS, Linux, Windows)

변경사항:
- 초기 버전 릴리즈
- goreleaser를 통한 크로스 플랫폼 빌드 지원
- NPM 패키지 배포 준비"

# 3. 태그 푸시
git push origin v0.1.0

# 4. 릴리즈 생성
~/go/bin/goreleaser release --clean
```

4. 릴리즈 결과물
- GitHub 릴리즈 페이지에 드래프트 생성
- 크로스 플랫폼 바이너리 자동 업로드:
  - macOS (arm64/amd64)
  - Linux (arm64/amd64)
  - Windows (amd64)
- 체크섬 파일 생성
- 변경 로그 자동 생성

5. 릴리즈 공개
- GitHub 릴리즈 페이지에서 내용 확인
- "Publish release" 버튼 클릭하여 공개

이렇게 생성된 릴리즈의 바이너리 URL들은 NPM 패키지의 설치 스크립트에서 사용됩니다.

### NPM 패키지 배포

1. NPM 패키지 구조
```bash
npm/
├── package.json     # 패키지 메타데이터 및 설치 스크립트
├── rulesctl.js      # 실행 및 설치 스크립트
└── bin/            # 바이너리가 설치될 디렉토리 (비어있어도 됨)
```

2. `package.json` 설정
```json
{
  "name": "rulesctl",
  "version": "0.1.0",
  "bin": {
    "rulesctl": "./rulesctl.js"
  },
  "files": [
    "rulesctl.js",
    "bin/"
  ],
  "scripts": {
    "postinstall": "node rulesctl.js install"
  }
}
```

3. 설치 스크립트 작동 방식
- `npm install -g rulesctl` 실행 시:
  1. NPM이 패키지 설치
  2. postinstall 스크립트 실행
  3. GitHub 릴리즈에서 해당 플랫폼 바이너리 자동 다운로드
  4. 실행 권한 설정

4. 로컬 테스트
```bash
# 패키지 디렉토리로 이동
cd npm

# bin 디렉토리 생성
mkdir -p bin

# 패키지 생성
npm pack

# 로컬 설치 테스트
npm install -g ./rulesctl-0.1.0.tgz

# 설치 확인
rulesctl --version
```

5. NPM 배포
```bash
# NPM 로그인
npm login

# 패키지 배포
npm publish

# 특정 태그로 배포 (옵션)
npm publish --tag beta
```

6. 설치 및 사용
```bash
# 글로벌 설치
npm install -g rulesctl

# 사용
rulesctl --help
```

주의사항:
- GitHub 릴리즈가 공개되어 있어야 함
- 릴리즈의 바이너리 URL이 올바르게 설정되어 있어야 함
- `package.json`의 버전과 GitHub 릴리즈 태그가 일치해야 함