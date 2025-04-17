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

// GetRulesDirPath는 현재 작업 디렉토리에서 .cursor/rules 디렉토리의 경로를 반환합니다.
func GetRulesDirPath() (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("작업 디렉토리 확인 실패: %w", err)
	}

	// 실제 경로 얻기 (심볼릭 링크 해결)
	realCwd, err := filepath.EvalSymlinks(cwd)
	if err != nil {
		return "", fmt.Errorf("심볼릭 링크 해결 실패: %w", err)
	}

	return filepath.Join(realCwd, RulesDirName), nil
}

// EnsureRulesDir는 .cursor/rules 디렉토리가 존재하는지 확인하고, 없으면 생성합니다.
func EnsureRulesDir() error {
	dirPath, err := GetRulesDirPath()
	if err != nil {
		return err
	}

	if err := os.MkdirAll(dirPath, 0755); err != nil {
		return fmt.Errorf("규칙 디렉토리 생성 실패: %w", err)
	}
	return nil
}

// ListLocalRules는 로컬 .cursor/rules 디렉토리의 모든 .mdc 파일을 탐색하여 경로와 MD5 해시를 반환합니다.
func ListLocalRules() (map[string]string, error) {
	rulesDir, err := GetRulesDirPath()
	if err != nil {
		return nil, err
	}

	// .cursor/rules 디렉토리가 존재하는지 확인
	if _, err := os.Stat(rulesDir); os.IsNotExist(err) {
		return nil, fmt.Errorf(".cursor/rules 디렉토리가 존재하지 않습니다. 'rulesctl init' 명령어로 초기화하거나 직접 디렉토리를 생성해주세요")
	}

	files := make(map[string]string)
	err = filepath.Walk(rulesDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && strings.HasSuffix(path, ".mdc") {
			hash, err := calculateMD5(path)
			if err != nil {
				return fmt.Errorf("파일 해시 계산 실패 %s: %w", path, err)
			}

			// 상대 경로로 변환
			relPath, err := filepath.Rel(rulesDir, path)
			if err != nil {
				return fmt.Errorf("상대 경로 변환 실패 %s: %w", path, err)
			}

			files[relPath] = hash
		}
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("규칙 파일 탐색 실패: %w", err)
	}

	// 파일이 없는 경우
	if len(files) == 0 {
		return nil, fmt.Errorf(".cursor/rules 디렉토리에 .mdc 파일이 없습니다. 규칙 파일을 추가해주세요")
	}

	return files, nil
}

// calculateMD5는 파일의 MD5 해시를 계산합니다.
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

// SaveRuleFile은 규칙 파일을 지정된 경로에 저장합니다.
func SaveRuleFile(relativePath string, content []byte) error {
	rulesDir, err := GetRulesDirPath()
	if err != nil {
		return err
	}

	fullPath := filepath.Join(rulesDir, relativePath)
	dirPath := filepath.Dir(fullPath)

	// 디렉토리가 없으면 생성
	if err := os.MkdirAll(dirPath, 0755); err != nil {
		return fmt.Errorf("디렉토리 생성 실패 %s: %w", dirPath, err)
	}

	// 파일 저장
	if err := os.WriteFile(fullPath, content, 0644); err != nil {
		return fmt.Errorf("파일 저장 실패 %s: %w", fullPath, err)
	}

	return nil
}

// DeleteRuleFile은 지정된 규칙 파일을 삭제합니다.
func DeleteRuleFile(relativePath string) error {
	rulesDir, err := GetRulesDirPath()
	if err != nil {
		return err
	}

	fullPath := filepath.Join(rulesDir, relativePath)
	if err := os.Remove(fullPath); err != nil {
		return fmt.Errorf("파일 삭제 실패 %s: %w", fullPath, err)
	}

	return nil
} 