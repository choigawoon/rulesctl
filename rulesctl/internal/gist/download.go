package gist

import (
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

	// .cursor/rules 디렉토리 생성
	rulesDir := filepath.Join(workDir, ".cursor", "rules")
	if err := os.MkdirAll(rulesDir, 0755); err != nil {
		return fmt.Errorf(".cursor/rules 디렉토리 생성 실패: %w", err)
	}

	// Gist 가져오기
	gist, err := FetchGist(token, gistID)
	if err != nil {
		return err
	}

	// 각 파일 다운로드
	for _, file := range meta.Files {
		// Gist에서 파일 찾기
		gistFile, exists := gist.Files[file.GistName]
		if !exists {
			return fmt.Errorf("Gist에서 파일을 찾을 수 없습니다: %s", file.GistName)
		}

		// 로컬 파일 경로
		localPath := filepath.Join(rulesDir, file.Path)

		// 디렉토리 생성
		if err := os.MkdirAll(filepath.Dir(localPath), 0755); err != nil {
			return fmt.Errorf("디렉토리 생성 실패: %w", err)
		}

		// 파일이 이미 존재하고 force가 false인 경우 건너뛰기
		if !force {
			if _, err := os.Stat(localPath); err == nil {
				continue
			}
		}

		// 파일 다운로드
		if err := downloadFile(gistFile.RawURL, localPath); err != nil {
			return fmt.Errorf("파일 다운로드 실패 (%s): %w", file.Path, err)
		}
	}

	return nil
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