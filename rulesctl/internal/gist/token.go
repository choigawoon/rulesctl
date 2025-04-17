package gist

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

const (
	configDir  = ".rulesctl"
	configFile = "config.json"
)

type Config struct {
	Token    string `json:"token"`
	LastUsed string `json:"last_used"`
}

func getToken() (string, error) {
	// 1. 환경 변수에서 토큰 확인
	if token := os.Getenv("GITHUB_TOKEN"); token != "" {
		return token, nil
	}

	// 2. 설정 파일에서 토큰 확인
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("홈 디렉토리 확인 실패: %v", err)
	}

	configPath := filepath.Join(homeDir, configDir, configFile)
	configData, err := os.ReadFile(configPath)
	if err != nil {
		return "", fmt.Errorf("설정 파일 읽기 실패: %v", err)
	}

	var config Config
	if err := json.Unmarshal(configData, &config); err != nil {
		return "", fmt.Errorf("설정 파일 파싱 실패: %v", err)
	}

	if config.Token == "" {
		return "", fmt.Errorf("GitHub 토큰이 설정되지 않았습니다. 'rulesctl auth' 명령어를 실행하여 토큰을 설정해주세요")
	}

	return config.Token, nil
}

func SaveToken(token string) error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("홈 디렉토리 확인 실패: %v", err)
	}

	configDirPath := filepath.Join(homeDir, configDir)
	if err := os.MkdirAll(configDirPath, 0755); err != nil {
		return fmt.Errorf("설정 디렉토리 생성 실패: %v", err)
	}

	config := Config{
		Token:    token,
		LastUsed: "now", // TODO: 실제 시간으로 변경
	}

	configData, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("설정 데이터 마샬링 실패: %v", err)
	}

	configPath := filepath.Join(configDirPath, configFile)
	if err := os.WriteFile(configPath, configData, 0600); err != nil {
		return fmt.Errorf("설정 파일 쓰기 실패: %v", err)
	}

	return nil
} 