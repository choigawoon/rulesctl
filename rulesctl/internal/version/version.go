package version

// Version은 rulesctl의 현재 버전입니다.
// 이 값은 goreleaser에 의해 빌드 시점에 업데이트됩니다.
var Version = "0.1.3"

// BuildTime은 빌드 시점의 시간입니다.
var BuildTime = "unknown"

// GitCommit은 빌드된 소스코드의 Git 커밋 해시입니다.
var GitCommit = "unknown" 