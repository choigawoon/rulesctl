package gist

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

// FetchGist fetches a Gist with the specified ID.
func FetchGist(token, gistID string) (*Gist, error) {
	client := &http.Client{}
	url := fmt.Sprintf("%s/gists/%s", baseURL, gistID)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Add authorization header only if token is provided
	if token != "" {
		req.Header.Set("Authorization", "token "+token)
	}
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make API request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("Gist not found: %s", gistID)
	}

	if resp.StatusCode == http.StatusUnauthorized {
		if token == "" {
			return nil, fmt.Errorf("this Gist requires authentication. Please run 'rulesctl auth' to set your token")
		}
		return nil, fmt.Errorf("invalid or expired token")
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed: %s", resp.Status)
	}

	var gist Gist
	if err := json.NewDecoder(resp.Body).Decode(&gist); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &gist, nil
}

// ParseMetadataFromGist는 Gist의 메타데이터 파일 내용을 파싱합니다.
func ParseMetadataFromGist(content string) (*Metadata, error) {
	var meta Metadata
	if err := json.Unmarshal([]byte(content), &meta); err != nil {
		return nil, fmt.Errorf("failed to parse metadata: %w", err)
	}
	return &meta, nil
}

// CheckConflicts는 다운로드할 파일과 로컬 파일 간의 충돌을 검사합니다.
func CheckConflicts(meta *Metadata) ([]string, error) {
	var conflicts []string

	// 현재 작업 디렉토리 확인
	workDir, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("failed to get working directory: %w", err)
	}

	// 각 파일에 대해 충돌 검사
	for _, file := range meta.Files {
		localPath := filepath.Join(workDir, ".cursor", "rules", file.Path)
		if _, err := os.Stat(localPath); err == nil {
			conflicts = append(conflicts, file.Path)
		}
	}

	return conflicts, nil
}

// DownloadFiles downloads files from a Gist to local.
func DownloadFiles(token, gistID string, meta *Metadata, force bool) error {
	// Check current working directory
	workDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get working directory: %w", err)
	}

	// Set temporary directory path
	tmpDir := filepath.Join(workDir, ".rulesctl", "tmp", gistID)
	rulesDir := filepath.Join(workDir, ".cursor", "rules")

	// Remove existing temporary directory if exists
	if err := os.RemoveAll(tmpDir); err != nil {
		return fmt.Errorf("failed to remove existing temporary directory: %w", err)
	}

	// Create temporary directory
	if err := os.MkdirAll(tmpDir, 0755); err != nil {
		return fmt.Errorf("failed to create temporary directory: %w", err)
	}
	defer os.RemoveAll(tmpDir) // Remove temporary directory after completion

	// Fetch Gist
	gist, err := FetchGist(token, gistID)
	if err != nil {
		return err
	}

	// Download and verify each file
	for _, file := range meta.Files {
		// Find file in Gist
		gistFile, exists := gist.Files[file.GistName]
		if !exists {
			return fmt.Errorf("file not found in Gist: %s", file.GistName)
		}

		// Download file to temporary directory
		tmpPath := filepath.Join(tmpDir, file.Path)
		if err := os.MkdirAll(filepath.Dir(tmpPath), 0755); err != nil {
			return fmt.Errorf("failed to create temporary directory: %w", err)
		}

		if err := downloadFile(gistFile.RawURL, tmpPath); err != nil {
			return fmt.Errorf("failed to download file (%s): %w", file.Path, err)
		}

		// Verify MD5 hash
		hash, err := calculateMD5(tmpPath)
		if err != nil {
			return fmt.Errorf("failed to calculate MD5 hash (%s): %w", file.Path, err)
		}

		if hash != file.MD5 {
			return fmt.Errorf("MD5 hash mismatch (%s): expected %s, got %s", file.Path, file.MD5, hash)
		}
	}

	// Check for file conflicts
	if !force {
		conflicts, err := CheckConflicts(meta)
		if err != nil {
			return fmt.Errorf("failed to check conflicts: %w", err)
		}
		if len(conflicts) > 0 {
			return fmt.Errorf("file conflicts detected. Use --force option to overwrite")
		}
	}

	// Create .cursor/rules directory
	if err := os.MkdirAll(rulesDir, 0755); err != nil {
		return fmt.Errorf("failed to create .cursor/rules directory: %w", err)
	}

	// Move verified files to final location
	for _, file := range meta.Files {
		tmpPath := filepath.Join(tmpDir, file.Path)
		finalPath := filepath.Join(rulesDir, file.Path)

		// Create target directory
		if err := os.MkdirAll(filepath.Dir(finalPath), 0755); err != nil {
			return fmt.Errorf("failed to create directory: %w", err)
		}

		// Remove existing file if force option is true
		if force {
			if err := os.RemoveAll(finalPath); err != nil {
				return fmt.Errorf("failed to remove existing file (%s): %w", file.Path, err)
			}
		}

		// Move file
		if err := os.Rename(tmpPath, finalPath); err != nil {
			return fmt.Errorf("failed to move file (%s): %w", file.Path, err)
		}
	}

	return nil
}

// calculateMD5는 파일의 MD5 해시를 계산합니다.
func calculateMD5(filepath string) (string, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", hash.Sum(nil)), nil
}

// downloadFile은 URL에서 파일을 다운로드합니다.
func downloadFile(url, filepath string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to download file: %s", resp.Status)
	}

	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
} 