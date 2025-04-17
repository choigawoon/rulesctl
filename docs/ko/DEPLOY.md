# 배포 가이드

이 문서는 rulesctl의 새로운 버전을 배포하는 과정을 설명합니다.

## 1. 버전 업데이트

버전 정보는 다음 파일들에서 관리됩니다:

- `internal/version/version.go`: Go 코드의 중앙 버전 관리
- `npm/rulesctl.js`: NPM 패키지의 버전 정보

### 버전 업데이트 방법

1. `internal/version/version.go` 파일에서 `Version` 변수 업데이트:
```go
var Version = "x.y.z"  // 새로운 버전으로 변경
```

2. `npm/rulesctl.js` 파일에서 `VERSION` 상수 업데이트:
```javascript
const VERSION = "vx.y.z";  // 새로운 버전으로 변경 (v 접두사 필수)
```

## 2. 릴리즈 생성

1. 변경사항 커밋:
```bash
git add .
git commit -m "chore: bump version to x.y.z"
```

2. 태그 생성 및 푸시:
```bash
git tag -a vx.y.z -m "Release vx.y.z"
git push origin vx.y.z
```

3. goreleaser로 릴리즈 생성:
```bash
goreleaser release --clean
```

4. GitHub에서 릴리즈 확인 및 Publish
- https://github.com/choigawoon/rulesctl/releases 에서 draft 상태의 릴리즈를 확인
- "Publish release" 버튼을 클릭하여 공개

## 3. NPM 패키지 배포

1. 패키지 생성 및 테스트:
```bash
cd npm
npm pack
npm install -g ./rulesctl-x.y.z.tgz  # 로컬 테스트
rulesctl version  # 버전 확인
```

2. NPM 로그인 및 배포:
```bash
npm login  # NPM 계정으로 로그인
npm publish  # 패키지 배포
```

## 4. 배포 후 확인

1. NPM 패키지 설치 테스트:
```bash
npm install -g rulesctl
rulesctl version  # 새 버전 확인
```

2. 기본 기능 테스트:
```bash
rulesctl init --sample
rulesctl list
rulesctl upload "test"
rulesctl download "test"
rulesctl delete "test"
```

## 주의사항

1. 버전 관리
- 버전은 반드시 `internal/version/version.go`와 `npm/rulesctl.js`에서 동일하게 관리
- goreleaser는 빌드 시점에 ldflags를 통해 버전 정보를 업데이트

2. GitHub 토큰
- goreleaser는 `GITHUB_TOKEN` 환경 변수를 사용
- 토큰에 `repo` 권한이 있는지 확인

3. NPM 배포
- `package.json`의 버전이 올바른지 확인
- 배포 전 로컬에서 반드시 테스트 