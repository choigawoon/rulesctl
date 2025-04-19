# rulesctl

GitHub Gist를 이용한 Cursor Rules 관리 및 공유 도구

## 문제 → 해결책

Cursor Rules를 효과적으로 관리하고 다른 개발자들과 손쉽게 공유하기 위한 도구가 필요했습니다. rulesctl은 GitHub Gist를 활용하여 규칙을 체계적으로 저장하고, public/private 옵션을 통해 선택적으로 공유할 수 있게 해줍니다.

## 사용법

### 설치

NPM을 통해 설치할 수 있습니다:

```bash
npm install -g rulesctl
```

### 인증 설정

GitHub 토큰은 다음 기능을 사용할 때 필요합니다:
- 규칙 업로드 (public/private)
- 내 Gist 목록 조회
- 제목으로 규칙 검색
- private Gist 다운로드

> **Note**: Public Gist를 ID로 다운로드할 때는 토큰이 필요하지 않습니다!

토큰 설정은 다음 두 가지 방법으로 할 수 있습니다:

1. 환경 변수 사용 (권장)
```bash
export GITHUB_TOKEN="your_github_token"
# 환경 변수 사용 시 'rulesctl auth' 명령어 실행이 필요하지 않습니다
```

2. auth 명령어 사용
```bash
rulesctl auth  # 프롬프트에서 토큰 입력
```

auth 명령어를 사용할 경우 토큰은 홈 디렉토리의 `~/.rulesctl/config` 파일에 JSON 형식으로 저장됩니다.

> **중요**: Personal Access Token에는 다음 권한이 필요합니다:
> - Gist (읽기/쓰기) 권한
> - repo 권한 (https://github.com/PatrickJS/awesome-cursorrules/tree/main/rules-new 의 파일 목록 접근용)

토큰 생성 방법은 [GitHub 공식 문서](https://docs.github.com/ko/authentication/keeping-your-account-and-data-secure/creating-a-personal-access-token)를 참조하세요.

### 시작하기

시작하는 방법에는 두 가지가 있습니다:

1. 예제 규칙으로 시작하기 (권장)
```bash
# 공유된 예제 규칙 다운로드 (토큰 불필요)
rulesctl download --gistid 74abf627d19e4114ac51bf0b6fbec99d

# 또는 직접 예제 생성
rulesctl init --sample
```

2. 새로 시작하기
```bash
# 빈 규칙 디렉토리 생성
rulesctl init
```

### 기본 명령어

```bash
# 도움말 보기
rulesctl --help

# 규칙 디렉토리 생성
rulesctl init
rulesctl init --sample  # 예제 규칙 파일도 함께 생성

# 예제 규칙 다운로드 (토큰 불필요)
rulesctl download --gistid 74abf627d19e4114ac51bf0b6fbec99d

# 규칙 목록 보기 (최근 1달 이내만 표시)
rulesctl list                # Public/Private 여부 및 기본 정보 표시
rulesctl list --detail      # revision 정보 포함하여 상세 표시

# 규칙 업로드하기
rulesctl upload "규칙세트이름"        # private으로 업로드 (기본값)
rulesctl upload "규칙세트이름" --public  # public으로 업로드 (다른 사용자와 공유 가능)

# 규칙 다운로드하기
rulesctl download "규칙세트이름"         # 내 Gist에서 제목으로 검색
rulesctl download --gistid abc123       # 공개된 Gist ID로 다운로드 (토큰 불필요)
```

### 규칙 공유하기 📢

rulesctl을 사용하면 다른 개발자들과 손쉽게 규칙을 공유할 수 있습니다:

1. 규칙 공유하기 (업로드)
```bash
# 규칙을 public으로 업로드
rulesctl upload "python-best-practices" --public

# 업로드 후 list 명령어로 Gist ID 확인
rulesctl list
# Type     Title                    Last Modified         Gist ID
# --------------------------------------------------------------
# Public   python-best-practices    2024-03-20 15:04:05  abc123...
```

2. 규칙 받아오기 (다운로드)
```bash
# 다른 사용자의 public 규칙을 Gist ID로 다운로드
rulesctl download --gistid abc123  # GitHub 토큰 없이도 가능!

# 충돌이 있는 경우 강제 다운로드
rulesctl download --gistid abc123 --force
```

> **Tip**: Public으로 업로드된 규칙은 GitHub 토큰 없이도 다운로드할 수 있어, 팀원들과 쉽게 공유할 수 있습니다!

### 사용 예시

먼저 인증 설정을 합니다:
```bash
# GitHub 토큰으로 인증
rulesctl auth
# 프롬프트에 Personal Access Token 입력
```

규칙 디렉토리 생성하기:
```bash
# .cursor/rules 디렉토리 생성
rulesctl init

# 예제와 함께 시작하기
rulesctl init --sample
```

규칙 세트 업로드하기:
```bash
# 현재 디렉토리의 규칙을 private으로 업로드 (기본값)
rulesctl upload "my-python-ruleset"

# 다른 사람과 공유하기 위해 public으로 업로드
rulesctl upload "my-python-ruleset" --public

# 업로드된 규칙의 public/private 상태 확인
rulesctl list

# 중복된 이름으로 강제 업로드
rulesctl upload "my-python-ruleset" --force
```

규칙 세트 다운로드하기:
```bash
# 내 Gist에서 제목으로 검색하여 다운로드
rulesctl download "my-python-ruleset"

# 다른 사람의 public Gist를 ID로 다운로드 (토큰 불필요)
rulesctl download --gistid abc123

# 충돌이 있어도 강제로 다운로드
rulesctl download --gistid abc123 --force
```

규칙 세트 삭제하기:
```bash
# 제목으로 검색하여 삭제
rulesctl delete "my-python-ruleset"

# 확인 없이 바로 삭제
rulesctl delete "my-python-ruleset" --force
```

> **중요**:
> - 다운로드는 두 가지 방식을 지원합니다:
>   1. 제목으로 다운로드: 내 Gist 목록에서 제목이 정확히 일치하는 규칙을 찾아 다운로드
>   2. Gist ID로 다운로드: 공개된 Gist의 ID를 직접 지정하여 다운로드
> - 다운로드 시 현재 실행 경로에 `.cursor/rules` 디렉토리가 없으면 자동으로 생성됩니다.
> - 원래 업로드된 디렉토리 구조와 파일들이 그대로 복원됩니다.
> - 다운로드 후 바로 사용할 수 있는 상태로 준비됩니다.
![1](docs/images/how-to-get-token-1.png)
![2](docs/images/how-to-get-token-2.png)

## 지원 플랫폼

rulesctl은 다음 플랫폼을 지원합니다:
- macOS (darwin)
- Linux
- Windows

## 개발자 가이드

개발 및 테스트 방법은 [개발 시작 가이드](docs/ko/GET-STARTED.md)를 참조하세요.

## 기여하기

기여 방법은 [기여 가이드](docs/ko/GET-STARTED.md#기여-가이드)를 참조하세요. 

## 로드맵 🚀

앞으로 제공될 예정인 기능들입니다:

### 설치 개선
- [ ] Windows 사용자를 위한 설치 프로세스 개선
  - Chocolatey 패키지 제공
  - 자동 PATH 환경변수 설정
  - 설치 과정 간소화

### 템플릿 규칙 세트
- [ ] 주요 기술 스택별 템플릿 규칙 세트 제공
  - 프론트엔드
    - React 개발 환경
    - Vue.js 개발 환경
    - Next.js/Nuxt.js 개발 환경
  - 백엔드
    - FastAPI 개발 환경
    - NestJS 개발 환경
    - Spring Boot 개발 환경
  - DevOps
    - Kubernetes/kubectl 작업 환경
    - Terraform 인프라 관리
    - Docker 컨테이너 관리

### 사용자 경험 개선
- [ ] 규칙 검색 및 필터링 기능
- [ ] 규칙 세트 버전 관리 기능
- [ ] 팀 협업을 위한 규칙 공유 기능
- [ ] 웹 인터페이스 제공

### 커뮤니티 참여 🤝
- [ ] 사용자 피드백 기반 개선
  - GitHub Discussions를 통한 의견 수렴
  - Reddit 커뮤니티 피드백 수집
  - 기술 블로그를 통한 사용 경험 공유
- [ ] 템플릿 규칙 세트 기여 가이드 제공
  - 커뮤니티 기반 템플릿 제작 및 공유
  - 템플릿 품질 관리 기준 수립
- [ ] 다국어 문서화 지원 확대

진행 상황과 새로운 기능 요청은 [GitHub Issues](https://github.com/choigawoon/rulesctl/issues)를 통해 확인하고 제안하실 수 있습니다.

더 나은 도구를 만들기 위해 여러분의 의견이 필요합니다:
- 💡 아이디어 제안: [GitHub Discussions](https://github.com/choigawoon/rulesctl/discussions/categories/ideas)
- 🐛 버그 리포트: [GitHub Issues](https://github.com/choigawoon/rulesctl/issues)
- 💬 질문하기: [GitHub Discussions Q&A](https://github.com/choigawoon/rulesctl/discussions/categories/q-a)
- 📝 사용 경험 공유: [GitHub Discussions Show and tell](https://github.com/choigawoon/rulesctl/discussions/categories/show-and-tell) 