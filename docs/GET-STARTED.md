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

1. Go 개발 환경 설정
```bash
# Go 설치 (1.16 이상)
# https://golang.org/dl/

# 프로젝트 클론
git clone https://github.com/your-username/rulesctl.git
cd rulesctl

# 의존성 설치
go mod download
```

2. 개발 도구 설치
```bash
# Cobra CLI 설치
go install github.com/spf13/cobra-cli@latest

# 테스트 도구 설치
go install github.com/stretchr/testify@latest
```

## 구현 가이드

rulesctl의 구현 방법과 NPM 배포 방법은 다음 문서를 참조하세요:
- [rulesctl 구현 가이드](rulesctl/2-HOW.md)
- [NPM 배포 가이드](npm/2-HOW.md)

## 테스트

1. 단위 테스트 실행
```bash
go test ./...
```

2. 통합 테스트 실행
```bash
go test -tags=integration ./...
```

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