package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// Config represents rulesctl configuration
type Config struct {
	Token    string `json:"token"`
	LastUsed string `json:"last_used"`
}

var (
	configDir  string
	configFile string
)

func init() {
	// Config directory can be overridden by environment variable
	if envDir := os.Getenv("RULESCTL_CONFIG_DIR"); envDir != "" {
		configDir = envDir
	} else {
		// Use os.UserHomeDir() to find home directory
		homeDir, err := os.UserHomeDir()
		if err != nil {
			fmt.Printf("failed to find home directory: %v\n", err)
			os.Exit(1)
		}
		configDir = filepath.Join(homeDir, ".rulesctl")
	}
	
	configFile = filepath.Join(configDir, "config.json")
	
	if err := os.MkdirAll(configDir, 0700); err != nil {
		fmt.Printf("failed to create config directory: %v\n", err)
		os.Exit(1)
	}
}

// GetConfigDir returns the path of the config directory
func GetConfigDir() (string, error) {
	return configDir, nil
}

// getConfigPath returns the path of the config file
func getConfigPath() (string, error) {
	return configFile, nil
}

// EnsureConfigDir checks if the config directory exists and creates it if not
func EnsureConfigDir() error {
	return os.MkdirAll(configDir, 0700)
}

// LoadConfig loads configuration
func LoadConfig() (*Config, error) {
	// 1. Check token in environment variable
	if token := os.Getenv("GITHUB_TOKEN"); token != "" {
		return &Config{Token: token}, nil
	}

	// 2. Check token in config file
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

// SaveConfig saves configuration to file
func SaveConfig(config *Config) error {
	if err := EnsureConfigDir(); err != nil {
		return err
	}

	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to convert config to JSON: %w", err)
	}

	if err := os.WriteFile(configFile, data, 0600); err != nil {
		return fmt.Errorf("failed to save config file: %w", err)
	}

	// Debug log
	fmt.Printf("Config file saved: %s\n", configFile)
	return nil
}

// SaveToken saves GitHub token
func SaveToken(token string) error {
	if token == "" {
		return fmt.Errorf("token is empty")
	}

	config, err := LoadConfig()
	if err != nil {
		return err
	}

	config.Token = token
	return SaveConfig(config)
} 