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
├── docs/                   # 문서
│   ├── rulesctl/
│   └── npm/
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