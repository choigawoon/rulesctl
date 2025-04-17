# rulesctl

GitHub Gist를 이용한 Cursor Rules 관리 도구

## 문제 → 해결책

Cursor Rules를 효과적으로 관리하고 공유하기 위한 도구가 필요했습니다. rulesctl은 GitHub Gist를 활용하여 규칙을 체계적으로 저장하고 관리할 수 있게 해줍니다.

## 사용법

### 설치

NPM을 통해 설치할 수 있습니다:

```bash
npm install -g rulesctl
```

### 인증 설정 (필수)

rulesctl을 사용하기 위해서는 GitHub 인증 설정이 **반드시** 필요합니다. 인증 설정은 다음 두 가지 방법으로 할 수 있습니다:

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

### 기본 명령어

```bash
# 도움말 보기
rulesctl --help

# 규칙 디렉토리 생성
rulesctl init
rulesctl init --sample  # 예제 규칙 파일도 함께 생성

# 규칙 목록 보기 (최근 1달 이내만 표시)
rulesctl list

# 규칙 업로드하기 (기본적으로 비공개)
rulesctl upload "규칙세트이름"
rulesctl upload "규칙세트이름" --public  # 공개 Gist로 업로드

# 규칙 다운로드하기
rulesctl download "규칙세트이름"         # 내 Gist에서 제목으로 검색
rulesctl download --gistid abc123       # 공개 Gist ID로 다운로드

# 규칙 삭제하기
rulesctl delete "규칙세트이름"           # 제목으로 검색하여 삭제
rulesctl delete "규칙세트이름" --force    # 확인 없이 바로 삭제
```

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
# 현재 디렉토리의 규칙을 특정 이름으로 업로드 (기본적으로 비공개)
rulesctl upload "my-python-ruleset"

# 특정 이름과 설명으로 업로드
rulesctl upload "my-python-ruleset" --desc "Python 프로젝트를 위한 규칙 모음"

# 다른 사람과 공유하기 위해 공개로 업로드
rulesctl upload "my-python-ruleset" --public

# 중복된 이름으로 강제 업로드 (확인 프롬프트 없음)
rulesctl upload "my-python-ruleset" --force
```

규칙 세트 삭제하기:
```bash
# 제목으로 검색하여 삭제
rulesctl delete "my-python-ruleset"

# 확인 없이 바로 삭제
rulesctl delete "my-python-ruleset" --force
```

규칙 목록 확인하기:
```bash
# 최근 1달 이내에 업로드된 규칙 목록 보기
rulesctl list
```

규칙 세트 다운로드하기:
```bash
# 내 Gist에서 제목으로 검색하여 다운로드
rulesctl download "my-python-ruleset"

# 공개된 Gist를 ID로 직접 다운로드
rulesctl download --gistid abc123

# 충돌이 있어도 강제로 다운로드
rulesctl download "my-python-ruleset" --force
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

## 지원 플랫폼

rulesctl은 다음 플랫폼을 지원합니다:
- macOS (darwin)
- Linux
- Windows

## 개발자 가이드

개발 및 테스트 방법은 [개발 시작 가이드](docs/ko/GET-STARTED.md)를 참조하세요.

## 기여하기

기여 방법은 [기여 가이드](docs/ko/GET-STARTED.md#기여-가이드)를 참조하세요. 