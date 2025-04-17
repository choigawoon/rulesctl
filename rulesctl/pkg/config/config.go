package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// Config는 rulesctl의 설정을 나타냅니다
type Config struct {
	Token    string `json:"token"`
	LastUsed string `json:"last_used"`
}

var (
	configDir  string
	configFile string
)

func init() {
	// 환경 변수로 설정 디렉토리 오버라이드 가능
	if envDir := os.Getenv("RULESCTL_CONFIG_DIR"); envDir != "" {
		configDir = envDir
	} else {
		// os.UserHomeDir()를 사용하여 홈 디렉토리 찾기
		homeDir, err := os.UserHomeDir()
		if err != nil {
			fmt.Printf("홈 디렉토리를 찾을 수 없습니다: %v\n", err)
			os.Exit(1)
		}
		configDir = filepath.Join(homeDir, ".rulesctl")
	}
	
	configFile = filepath.Join(configDir, "config.json")
	
	if err := os.MkdirAll(configDir, 0700); err != nil {
		fmt.Printf("설정 디렉토리를 생성할 수 없습니다: %v\n", err)
		os.Exit(1)
	}
	
	// 디버그 로그
	fmt.Printf("설정 디렉토리: %s\n", configDir)
	fmt.Printf("설정 파일: %s\n", configFile)
}

// GetConfigDir는 설정 디렉토리의 경로를 반환합니다
func GetConfigDir() (string, error) {
	return configDir, nil
}

// EnsureConfigDir는 설정 디렉토리가 존재하는지 확인하고, 없으면 생성합니다
func EnsureConfigDir() error {
	return os.MkdirAll(configDir, 0700)
}

// LoadConfig는 설정 파일을 로드합니다
func LoadConfig() (*Config, error) {
	// 디버그 로그
	fmt.Printf("설정 파일 로드 시도: %s\n", configFile)
	
	data, err := os.ReadFile(configFile)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Printf("설정 파일이 없습니다. 새로 생성합니다.\n")
			return &Config{}, nil
		}
		return nil, fmt.Errorf("설정 파일을 읽을 수 없습니다: %w", err)
	}

	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("설정 파일을 파싱할 수 없습니다: %w", err)
	}

	// 디버그 로그
	fmt.Printf("설정 파일 로드 완료. 토큰 길이: %d\n", len(config.Token))
	return &config, nil
}

// SaveConfig는 설정을 파일에 저장합니다
func SaveConfig(config *Config) error {
	if err := EnsureConfigDir(); err != nil {
		return err
	}

	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("설정을 JSON으로 변환할 수 없습니다: %w", err)
	}

	if err := os.WriteFile(configFile, data, 0600); err != nil {
		return fmt.Errorf("설정 파일을 저장할 수 없습니다: %w", err)
	}

	// 디버그 로그
	fmt.Printf("설정 파일 저장 완료: %s\n", configFile)
	return nil
}

// SaveToken은 GitHub 토큰을 저장합니다
func SaveToken(token string) error {
	if token == "" {
		return fmt.Errorf("토큰이 비어있습니다")
	}

	config, err := LoadConfig()
	if err != nil {
		return err
	}

	config.Token = token
	return SaveConfig(config)
} 