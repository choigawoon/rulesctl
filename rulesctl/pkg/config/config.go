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
}

// GetConfigDir는 설정 디렉토리의 경로를 반환합니다
func GetConfigDir() (string, error) {
	return configDir, nil
}

// getConfigPath는 설정 파일의 경로를 반환합니다
func getConfigPath() (string, error) {
	return configFile, nil
}

// EnsureConfigDir는 설정 디렉토리가 존재하는지 확인하고, 없으면 생성합니다
func EnsureConfigDir() error {
	return os.MkdirAll(configDir, 0700)
}

// LoadConfig는 설정을 로드합니다
func LoadConfig() (*Config, error) {
	// 1. 환경 변수에서 토큰 확인
	if token := os.Getenv("GITHUB_TOKEN"); token != "" {
		return &Config{Token: token}, nil
	}

	// 2. 설정 파일에서 토큰 확인
	configPath, err := getConfigPath()
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			return &Config{}, nil
		}
		return nil, err
	}

	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, err
	}

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