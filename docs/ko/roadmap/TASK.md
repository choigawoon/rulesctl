# rulesctl 개발 로드맵

## 히스토리
- 2025.04.17 - [TASK_v1.md](TASK_v1.md): 초기 개발 단계의 작업 히스토리

## 작업 기준
본 문서는 각 작업 단계를 **1시간 이내**에 완료할 수 있도록 세분화했습니다. 이는 작업의 진행 상황을 측정하고 관리하기 위한 기준입니다.

---

# 1차 배포 준비 (현재 기능)

## 버전 관리 개선 (2시간)
1. 버전 정보 중앙화
   - [x] internal/version/version.go 패키지 생성
   - [x] goreleaser ldflags 설정 업데이트
   - [ ] npm 패키지 버전 동기화 문제
     - 현재 상황: npm install -g ./rulesctl-0.1.2.tgz로 로컬 설치 시 버전이 0.1.0으로 표시됨
     - 원인: 
       1. goreleaser가 빌드한 바이너리와 npm 패키지의 바이너리가 다름
       2. npm/rulesctl.js의 VERSION과 실제 바이너리 버전이 불일치
       3. GitHub Release API 캐시 문제
          - GitHub release 배포 후 일정 시간이 지나야 새 버전이 정상적으로 다운로드됨
          - API 응답이 캐시되어 있어 즉시 반영되지 않는 것으로 추정
     - 해결 방안 검토:
       1. npm pack 전에 goreleaser 빌드 바이너리를 npm/bin에 복사
       2. 로컬 설치 시에도 GitHub releases에서 바이너리 다운로드
       3. goreleaser post hooks 수정 검토
       4. GitHub Release 배포 후 일정 시간(약 5-10분) 대기 후 npm 배포 진행

## Windows 설치 문제 해결 (2시간)
- [ ] Windows npm 설치 프로세스 개선
  - [x] PowerShell 의존성 제거
  - [x] adm-zip 패키지를 사용한 ZIP 파일 처리 구현
  - [ ] 설치 프로세스 테스트
    - [ ] Windows 10 테스트
    - [ ] Windows 11 테스트
    - [ ] 다양한 Node.js 버전 테스트
  - [ ] 에러 처리 및 사용자 피드백 개선

## 배포 준비 (1시간)
1. 최종 테스트
   - [ ] macOS 설치 테스트
   - [ ] Linux 설치 테스트
   - [ ] Windows 설치 테스트
   - [ ] 기본 기능 동작 확인

2. 배포
   - [ ] GitHub Release 생성
   - [ ] NPM 패키지 배포
   - [ ] 설치 가이드 최종 검토

---

# 2차 마일스톤 (사용자 피드백 후)

## 사용자 피드백 수집 및 분석
- [ ] GitHub Issues를 통한 피드백 수집
- [ ] Reddit 등 커뮤니티 피드백 수집
- [ ] 사용자 경험 개선점 도출

## 개선 계획 (피드백 기반)
- [ ] 에러 처리 및 로깅 개선
- [ ] 진행 상황 표시 개선
- [ ] 사용법 및 도움말 개선
- [ ] GitHub Actions 워크플로우 설정

## 추가 기능 검토
- [ ] 규칙 검색 기능
- [ ] 규칙 공유 및 팀 협업 기능
- [ ] 규칙 버전 관리 기능
- [ ] 웹 인터페이스 개발
- [ ] 규칙 템플릿 및 스캐폴딩 기능
- [ ] Gist revision 지정 다운로드 기능

## 추가 배포 방법
- [ ] node-windows 패키지 통합
- [ ] Chocolatey 패키지 제공

## 마케팅 및 홍보
- [ ] Reddit 홍보
- [ ] 기술 블로그 포스팅
- [ ] 사용자 가이드 영상 제작

---

이전 버전의 작업 계획은 [TASK_2024_0319.md](TASK_2024_0319.md)를 참조하세요.
