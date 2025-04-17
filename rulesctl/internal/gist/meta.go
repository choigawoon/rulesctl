package gist

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const MetaFileName = "meta.json"

type FileMetadata struct {
	Path     string `json:"path"`     // 원본 파일 경로
	GistName string `json:"gist_name"` // Gist에서의 파일 이름
	Size     int64  `json:"size"`     // 파일 크기
	MD5      string `json:"md5"`      // MD5 해시
}

type DirectoryStructure map[string]interface{} // 중첩된 디렉토리 구조

type Metadata struct {
	SchemaVersion string             `json:"schema_version"`
	CLIVersion    string            `json:"cli_version"`
	UpdatedAt     time.Time         `json:"updated_at"`
	Structure     DirectoryStructure `json:"structure"`
	Files         []FileMetadata    `json:"files"`
}

func NewMetadata() *Metadata {
	return &Metadata{
		SchemaVersion: "1.0.0",
		CLIVersion:    "0.1.0", // TODO: 버전 관리
		UpdatedAt:     time.Now(),
		Structure:     make(DirectoryStructure),
		Files:         make([]FileMetadata, 0),
	}
}

// convertToGistName은 원본 파일 경로를 Gist 파일 이름으로 변환합니다.
// 예: python/linting.mdc -> python_linting.mdc
func convertToGistName(path string) string {
	dir := filepath.Dir(path)
	base := filepath.Base(path)
	ext := filepath.Ext(base)
	name := strings.TrimSuffix(base, ext)

	if dir == "." {
		return name + ext
	}

	// 디렉토리 구분자를 언더스코어로 변환
	dirParts := strings.Split(dir, string(filepath.Separator))
	return strings.Join(dirParts, "_") + "_" + name + ext
}

func (m *Metadata) AddFile(path string) error {
	// 이 함수에서는 path가 상대 경로이며, 실제 파일은 .cursor/rules 디렉토리 내에 있습니다.
	// 따라서 파일을 열기 전에 절대 경로를 구성해야 합니다.
	
	// 먼저 path가 이미 절대 경로인지 확인
	var fullPath string
	if filepath.IsAbs(path) {
		fullPath = path
	} else {
		// rulesDir 경로 가져오기
		dir, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("작업 디렉토리 확인 실패: %w", err)
		}
		fullPath = filepath.Join(dir, ".cursor/rules", path)
	}

	// 파일 열기
	file, err := os.Open(fullPath)
	if err != nil {
		return fmt.Errorf("파일 열기 실패 %s: %w", path, err)
	}
	defer file.Close()

	// 파일 크기 확인
	info, err := file.Stat()
	if err != nil {
		return fmt.Errorf("파일 정보 조회 실패 %s: %w", path, err)
	}

	// MD5 해시 계산
	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		return fmt.Errorf("해시 계산 실패 %s: %w", path, err)
	}

	// Gist 파일 이름 생성
	gistName := convertToGistName(path)

	metadata := FileMetadata{
		Path:     path,
		GistName: gistName,
		Size:     info.Size(),
		MD5:      hex.EncodeToString(hash.Sum(nil)),
	}
	m.Files = append(m.Files, metadata)

	// 디렉토리 구조 업데이트
	m.updateStructure(path)

	return nil
}

func (m *Metadata) updateStructure(path string) {
	parts := strings.Split(filepath.ToSlash(path), "/")
	current := m.Structure

	// 경로의 각 부분을 순회하며 구조 생성
	for i, part := range parts {
		if i == len(parts)-1 {
			// 마지막 부분은 파일
			current[part] = nil
		} else {
			// 디렉토리인 경우
			if _, exists := current[part]; !exists {
				current[part] = make(DirectoryStructure)
			}
			current = current[part].(DirectoryStructure)
		}
	}
}

func (m *Metadata) ToJSON() ([]byte, error) {
	return json.MarshalIndent(m, "", "  ")
}

// GetGistName은 원본 파일 경로에 대응하는 Gist 파일 이름을 반환합니다.
func (m *Metadata) GetGistName(path string) string {
	// 전체 경로와 상대 경로 모두 처리
	for _, file := range m.Files {
		// 전체 경로와 일치하는지 확인
		if file.Path == path {
			return file.GistName
		}
		
		// 상대 경로와 일치하는지 확인
		if filepath.Base(file.Path) == filepath.Base(path) {
			dir := filepath.Dir(path)
			if dir != "." && strings.HasSuffix(filepath.Dir(file.Path), dir) {
				return file.GistName
			}
		}
	}
	return convertToGistName(path)
} 