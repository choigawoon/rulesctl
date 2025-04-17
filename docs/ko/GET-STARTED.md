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

### goreleaser를 사용한 빌드

1. goreleaser 설치
```bash
go install github.com/goreleaser/goreleaser@latest
```

2. 로컬 빌드 테스트
```bash
# 스냅샷 빌드 (태그 없이 테스트)
~/go/bin/goreleaser build --snapshot --clean

# 결과물 확인
ls -l bin/           # 최신 빌드 결과물
ls -l dist/          # 모든 아키텍처별 빌드 결과물
```

3. 지원되는 빌드 타겟
- macOS (Darwin)
  - arm64 (Apple Silicon M1/M2)
  - amd64 (Intel)
- Linux
  - arm64
  - amd64
- Windows
  - amd64

### 릴리즈 배포

1. 버전 태그 생성
```bash
git tag -a v0.1.0 -m "First release"
git push origin v0.1.0
```

2. 릴리즈 빌드
```bash
~/go/bin/goreleaser release
```

이 명령은 다음 작업을 수행합니다:
- 모든 플랫폼용 바이너리 빌드
- 체크섬 생성
- 변경 로그 생성
- GitHub 릴리즈 페이지에 드래프트 생성

### 빌드 설정 커스터마이징

빌드 설정은 `.goreleaser.yaml` 파일에서 관리됩니다:
- 빌드 대상 플랫폼/아키텍처 추가/제거
- 빌드 후크 스크립트 수정
- 아카이브 포맷 변경
- 릴리즈 설정 변경

자세한 설정 옵션은 [goreleaser 공식 문서](https://goreleaser.com/customization/)를 참조하세요.

### NPM 패키지 배포

1. NPM 패키지 구조 설정
```bash
# 패키지 디렉토리 생성
mkdir -p npm
cd npm

# package.json 초기화
npm init -y

# 필요한 의존성 설치
npm install --save-dev @octokit/rest
```

2. `package.json` 설정
```json
{
  "name": "rulesctl",
  "version": "0.1.0",
  "description": "GitHub Gist를 이용한 Cursor Rules 관리 도구",
  "bin": {
    "rulesctl": "./bin/rulesctl"
  },
  "scripts": {
    "postinstall": "node scripts/install.js"
  },
  "files": [
    "bin",
    "scripts"
  ],
  "keywords": ["cursor", "rules", "gist"],
  "author": "choigawoon",
  "license": "MIT"
}
```

3. 바이너리 다운로드 스크립트 작성 (`scripts/install.js`)
```javascript
const { Octokit } = require('@octokit/rest');
const fs = require('fs');
const path = require('path');
const https = require('https');
const os = require('os');

const OWNER = 'choigawoon';
const REPO = 'rulesctl';

async function getLatestRelease() {
  const octokit = new Octokit();
  const { data } = await octokit.repos.getLatestRelease({
    owner: OWNER,
    repo: REPO
  });
  return data;
}

function getPlatformAsset(assets) {
  const platform = os.platform();
  const arch = os.arch();
  
  const platformMap = {
    'darwin': 'Darwin',
    'linux': 'Linux',
    'win32': 'Windows'
  };
  
  const archMap = {
    'x64': 'x86_64',
    'arm64': 'arm64'
  };
  
  const assetPattern = `${REPO}_${platformMap[platform]}_${archMap[arch]}`;
  return assets.find(asset => asset.name.includes(assetPattern));
}

async function downloadBinary(url, dest) {
  return new Promise((resolve, reject) => {
    https.get(url, response => {
      if (response.statusCode === 302) {
        https.get(response.headers.location, response => {
          const file = fs.createWriteStream(dest);
          response.pipe(file);
          file.on('finish', () => {
            file.close();
            fs.chmodSync(dest, '755');
            resolve();
          });
        }).on('error', reject);
      }
    }).on('error', reject);
  });
}

async function install() {
  try {
    const release = await getLatestRelease();
    const asset = getPlatformAsset(release.assets);
    
    if (!asset) {
      throw new Error('No matching binary found for your platform');
    }
    
    const binPath = path.join(__dirname, '..', 'bin');
    if (!fs.existsSync(binPath)) {
      fs.mkdirSync(binPath, { recursive: true });
    }
    
    const dest = path.join(binPath, os.platform() === 'win32' ? 'rulesctl.exe' : 'rulesctl');
    await downloadBinary(asset.browser_download_url, dest);
    
    console.log('rulesctl binary installed successfully!');
  } catch (error) {
    console.error('Failed to install rulesctl:', error);
    process.exit(1);
  }
}

install();
```

4. NPM 패키지 배포
```bash
# NPM 로그인
npm login

# 패키지 배포
npm publish

# 특정 태그로 배포
npm publish --tag beta  # 베타 버전으로 배포
```

5. 배포 후 설치 테스트
```bash
# 글로벌 설치 테스트
npm install -g rulesctl

# 실행 테스트
rulesctl --version
```

자세한 설정 옵션은 [goreleaser 공식 문서](https://goreleaser.com/customization/)를 참조하세요.