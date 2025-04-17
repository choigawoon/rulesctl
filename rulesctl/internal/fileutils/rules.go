package fileutils

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

const (
	RulesDirName = ".cursor/rules"
)

// GetRulesDirPath returns the path of .cursor/rules directory from the current working directory.
func GetRulesDirPath() (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("failed to get working directory: %w", err)
	}

	// Get real path (resolve symlinks)
	realCwd, err := filepath.EvalSymlinks(cwd)
	if err != nil {
		return "", fmt.Errorf("failed to resolve symlinks: %w", err)
	}

	return filepath.Join(realCwd, RulesDirName), nil
}

// EnsureRulesDir checks if .cursor/rules directory exists and creates it if not.
func EnsureRulesDir() error {
	dirPath, err := GetRulesDirPath()
	if err != nil {
		return err
	}

	if err := os.MkdirAll(dirPath, 0755); err != nil {
		return fmt.Errorf("failed to create rules directory: %w", err)
	}
	return nil
}

// ListLocalRules searches all .mdc files in local .cursor/rules directory and returns their paths and MD5 hashes.
func ListLocalRules() (map[string]string, error) {
	rulesDir, err := GetRulesDirPath()
	if err != nil {
		return nil, err
	}

	// Check if .cursor/rules directory exists
	if _, err := os.Stat(rulesDir); os.IsNotExist(err) {
		return nil, fmt.Errorf(".cursor/rules directory does not exist. Initialize with 'rulesctl init' command or create it manually")
	}

	files := make(map[string]string)
	err = filepath.Walk(rulesDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && strings.HasSuffix(path, ".mdc") {
			hash, err := calculateMD5(path)
			if err != nil {
				return fmt.Errorf("failed to calculate file hash %s: %w", path, err)
			}

			// Convert to relative path
			relPath, err := filepath.Rel(rulesDir, path)
			if err != nil {
				return fmt.Errorf("failed to convert to relative path %s: %w", path, err)
			}

			files[relPath] = hash
		}
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to scan rule files: %w", err)
	}

	// No files found
	if len(files) == 0 {
		return nil, fmt.Errorf("no .mdc files found in .cursor/rules directory. Please add rule files")
	}

	return files, nil
}

// calculateMD5 calculates the MD5 hash of a file.
func calculateMD5(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}

	return hex.EncodeToString(hash.Sum(nil)), nil
}

// SaveRuleFile saves a rule file at the specified path.
func SaveRuleFile(relativePath string, content []byte) error {
	rulesDir, err := GetRulesDirPath()
	if err != nil {
		return err
	}

	fullPath := filepath.Join(rulesDir, relativePath)
	dirPath := filepath.Dir(fullPath)

	// Create directory if it doesn't exist
	if err := os.MkdirAll(dirPath, 0755); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", dirPath, err)
	}

	// Save file
	if err := os.WriteFile(fullPath, content, 0644); err != nil {
		return fmt.Errorf("failed to save file %s: %w", fullPath, err)
	}

	return nil
}

// DeleteRuleFile deletes the specified rule file.
func DeleteRuleFile(relativePath string) error {
	rulesDir, err := GetRulesDirPath()
	if err != nil {
		return err
	}

	fullPath := filepath.Join(rulesDir, relativePath)
	if err := os.Remove(fullPath); err != nil {
		return fmt.Errorf("failed to delete file %s: %w", fullPath, err)
	}

	return nil
} 