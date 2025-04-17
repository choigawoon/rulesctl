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

var configDir = filepath.Join(os.Getenv("HOME"), ".rulesctl")
var configFile = filepath.Join(configDir, "config.json")

func init() {
	if err := os.MkdirAll(configDir, 0700); err != nil {
		fmt.Printf("설정 디렉토리를 생성할 수 없습니다: %v\n", err)
		os.Exit(1)
	}
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
	data, err := os.ReadFile(configFile)
	if err != nil {
		if os.IsNotExist(err) {
			return &Config{}, nil
		}
		return nil, fmt.Errorf("설정 파일을 읽을 수 없습니다: %w", err)
	}

	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("설정 파일을 파싱할 수 없습니다: %w", err)
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