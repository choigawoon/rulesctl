# NPM 배포 가이드

## 배포 구조

```
rulesctl/
├── npm/                    # NPM 패키징 관련 파일
│   ├── bin/               # 플랫폼별 바이너리
│   │   ├── rulesctl-darwin
│   │   ├── rulesctl-linux
│   │   └── rulesctl-win.exe
│   ├── rulesctl.js        # 런처 스크립트
│   └── package.json       # NPM 메타데이터
```

## 패키지 구성

### 1. package.json
```json
{
  "name": "rulesctl",
  "version": "1.0.0",
  "description": "Cursor Rules 관리 CLI 도구",
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

### 2. rulesctl.js (런처 스크립트)
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

## 빌드 및 배포 프로세스

1. **Go 바이너리 빌드**
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

2. **패키지 배포**
```bash
# NPM 패키지 배포
cd npm
npm publish
```

## 버전 관리

1. **시맨틱 버저닝 준수**
   - MAJOR: 호환되지 않는 API 변경
   - MINOR: 이전 버전과 호환되는 기능 추가
   - PATCH: 이전 버전과 호환되는 버그 수정

2. **버전 업데이트 프로세스**
```bash
# package.json 버전 업데이트
npm version [major|minor|patch]

# 변경사항 커밋 및 태그 생성
git commit -am "v1.0.0"
git tag v1.0.0
git push origin main --tags
```

## 테스트 프로세스

1. **로컬 테스트**
```bash
# 로컬에 패키지 링크
cd npm
npm link

# 다른 터미널에서 테스트
rulesctl --version
```

2. **CI/CD 파이프라인**
```yaml
# GitHub Actions 예시
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

## 보안 고려사항

1. **바이너리 서명**
   - 모든 바이너리에 디지털 서명 적용
   - 사용자 시스템에서 서명 검증

2. **의존성 검사**
   - 정기적인 보안 취약점 검사
   - 의존성 업데이트 자동화

3. **토큰 관리**
   - NPM 토큰 보안 저장
   - CI/CD 파이프라인에서 안전한 토큰 사용
