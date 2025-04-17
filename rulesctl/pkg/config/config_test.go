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
	// 테스트용 임시 디렉토리 생성
	tempDir, err := os.MkdirTemp("", "rulesctl-test")
	if err != nil {
		t.Fatalf("임시 디렉토리 생성 실패: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// 테스트용 설정 파일 경로 설정
	oldConfigDir := configDir
	oldConfigFile := configFile
	configDir = tempDir
	configFile = filepath.Join(tempDir, "config.json")
	defer func() {
		configDir = oldConfigDir
		configFile = oldConfigFile
	}()

	// 테스트 케이스
	tests := []struct {
		name    string
		token   string
		wantErr bool
	}{
		{
			name:    "유효한 토큰 저장",
			token:   "test-token",
			wantErr: false,
		},
		{
			name:    "빈 토큰 저장",
			token:   "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 토큰 저장 테스트
			err := SaveToken(tt.token)
			if (err != nil) != tt.wantErr {
				t.Errorf("SaveToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				return
			}

			// 설정 로드 테스트
			config, err := LoadConfig()
			if err != nil {
				t.Errorf("LoadConfig() error = %v", err)
				return
			}

			if config.Token != tt.token {
				t.Errorf("LoadConfig() token = %v, want %v", config.Token, tt.token)
			}
		})
	}
} 