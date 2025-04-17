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

// FetchGist는 지정된 ID의 Gist를 가져옵니다.
func FetchGist(token, gistID string) (*Gist, error) {
	client := &http.Client{}
	url := fmt.Sprintf("%s/gists/%s", baseURL, gistID)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("요청 생성 실패: %w", err)
	}

	req.Header.Set("Authorization", "token "+token)
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("API 요청 실패: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("Gist를 찾을 수 없습니다: %s", gistID)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API 요청 실패: %s", resp.Status)
	}

	var gist Gist
	if err := json.NewDecoder(resp.Body).Decode(&gist); err != nil {
		return nil, fmt.Errorf("응답 파싱 실패: %w", err)
	}

	return &gist, nil
}

// ParseMetadataFromGist는 Gist의 메타데이터 파일 내용을 파싱합니다.
func ParseMetadataFromGist(content string) (*Metadata, error) {
	var meta Metadata
	if err := json.Unmarshal([]byte(content), &meta); err != nil {
		return nil, fmt.Errorf("메타데이터 파싱 실패: %w", err)
	}
	return &meta, nil
}

// CheckConflicts는 다운로드할 파일과 로컬 파일 간의 충돌을 검사합니다.
func CheckConflicts(meta *Metadata) ([]string, error) {
	var conflicts []string

	// 현재 작업 디렉토리 확인
	workDir, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("작업 디렉토리 확인 실패: %w", err)
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

// DownloadFiles는 Gist의 파일들을 로컬에 다운로드합니다.
func DownloadFiles(token, gistID string, meta *Metadata, force bool) error {
	// 현재 작업 디렉토리 확인
	workDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("작업 디렉토리 확인 실패: %w", err)
	}

	// 임시 디렉토리 경로 설정
	tmpDir := filepath.Join(workDir, ".rulesctl", "tmp", gistID)
	rulesDir := filepath.Join(workDir, ".cursor", "rules")

	// 임시 디렉토리가 이미 존재하면 제거
	if err := os.RemoveAll(tmpDir); err != nil {
		return fmt.Errorf("기존 임시 디렉토리 제거 실패: %w", err)
	}

	// 임시 디렉토리 생성
	if err := os.MkdirAll(tmpDir, 0755); err != nil {
		return fmt.Errorf("임시 디렉토리 생성 실패: %w", err)
	}
	defer os.RemoveAll(tmpDir) // 작업 완료 후 임시 디렉토리 제거

	// Gist 가져오기
	gist, err := FetchGist(token, gistID)
	if err != nil {
		return err
	}

	// 각 파일 다운로드 및 검증
	for _, file := range meta.Files {
		// Gist에서 파일 찾기
		gistFile, exists := gist.Files[file.GistName]
		if !exists {
			return fmt.Errorf("Gist에서 파일을 찾을 수 없습니다: %s", file.GistName)
		}

		// 임시 디렉토리에 파일 다운로드
		tmpPath := filepath.Join(tmpDir, file.Path)
		if err := os.MkdirAll(filepath.Dir(tmpPath), 0755); err != nil {
			return fmt.Errorf("임시 디렉토리 생성 실패: %w", err)
		}

		if err := downloadFile(gistFile.RawURL, tmpPath); err != nil {
			return fmt.Errorf("파일 다운로드 실패 (%s): %w", file.Path, err)
		}

		// MD5 해시 검증
		hash, err := calculateMD5(tmpPath)
		if err != nil {
			return fmt.Errorf("MD5 해시 계산 실패 (%s): %w", file.Path, err)
		}

		if hash != file.MD5 {
			return fmt.Errorf("MD5 해시 불일치 (%s): expected %s, got %s", file.Path, file.MD5, hash)
		}
	}

	// 파일 충돌 검사
	if !force {
		conflicts, err := CheckConflicts(meta)
		if err != nil {
			return fmt.Errorf("충돌 검사 실패: %w", err)
		}
		if len(conflicts) > 0 {
			return fmt.Errorf("파일 충돌이 발생했습니다. --force 옵션을 사용하여 덮어쓸 수 있습니다")
		}
	}

	// .cursor/rules 디렉토리 생성
	if err := os.MkdirAll(rulesDir, 0755); err != nil {
		return fmt.Errorf(".cursor/rules 디렉토리 생성 실패: %w", err)
	}

	// 검증이 완료된 파일들을 최종 위치로 이동
	for _, file := range meta.Files {
		tmpPath := filepath.Join(tmpDir, file.Path)
		finalPath := filepath.Join(rulesDir, file.Path)

		// 대상 디렉토리 생성
		if err := os.MkdirAll(filepath.Dir(finalPath), 0755); err != nil {
			return fmt.Errorf("디렉토리 생성 실패: %w", err)
		}

		// 기존 파일이 있으면 제거 (force 옵션이 true일 때)
		if force {
			if err := os.RemoveAll(finalPath); err != nil {
				return fmt.Errorf("기존 파일 제거 실패 (%s): %w", file.Path, err)
			}
		}

		// 파일 이동
		if err := os.Rename(tmpPath, finalPath); err != nil {
			return fmt.Errorf("파일 이동 실패 (%s): %w", file.Path, err)
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
		return fmt.Errorf("파일 다운로드 실패: %s", resp.Status)
	}

	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
} 