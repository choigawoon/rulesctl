package gist

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
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
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	// 파일 크기 확인
	info, err := file.Stat()
	if err != nil {
		return err
	}

	// MD5 해시 계산
	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		return err
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
	for _, file := range m.Files {
		if file.Path == path {
			return file.GistName
		}
	}
	return convertToGistName(path)
} 