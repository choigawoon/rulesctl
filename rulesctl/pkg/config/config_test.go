package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestConfigDir(t *testing.T) {
	dir, err := GetConfigDir()
	if err != nil {
		t.Fatalf("GetConfigDir 실패: %v", err)
	}

	if dir == "" {
		t.Error("설정 디렉토리 경로가 비어있습니다")
	}

	if !filepath.IsAbs(dir) {
		t.Error("설정 디렉토리 경로가 절대 경로가 아닙니다")
	}
}

func TestConfigOperations(t *testing.T) {
	// 테스트를 위한 임시 홈 디렉토리 설정
	tmpHome := t.TempDir()
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tmpHome)
	defer os.Setenv("HOME", originalHome)

	// 설정 디렉토리 생성 테스트
	if err := EnsureConfigDir(); err != nil {
		t.Fatalf("EnsureConfigDir 실패: %v", err)
	}

	dir, _ := GetConfigDir()
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		t.Error("설정 디렉토리가 생성되지 않았습니다")
	}

	// 토큰 저장 테스트
	testToken := "test-token"
	if err := SaveToken(testToken); err != nil {
		t.Fatalf("SaveToken 실패: %v", err)
	}

	// 설정 로드 테스트
	config, err := LoadConfig()
	if err != nil {
		t.Fatalf("LoadConfig 실패: %v", err)
	}

	if config.Token != testToken {
		t.Errorf("토큰이 일치하지 않습니다. 기대값: %s, 실제값: %s", testToken, config.Token)
	}

	if config.LastUsed == "" {
		t.Error("LastUsed가 설정되지 않았습니다")
	}
} 