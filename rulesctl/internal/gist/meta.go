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
		CLIVersion:    "0.1.3", // TODO: 버전 관리
		UpdatedAt:     time.Now(),
		Structure:     make(DirectoryStructure),
		Files:         make([]FileMetadata, 0),
	}
}

// convertToGistName은 원본 파일 경로를 Gist 파일 이름으로 변환합니다.
// 예: python/linting.mdc -> python_linting_mdc
func convertToGistName(path string) string {
	// 디렉토리와 파일명 분리
	dir := filepath.Dir(path)
	base := filepath.Base(path)

	// 파일명에서 확장자 분리
	ext := filepath.Ext(base)
	name := strings.TrimSuffix(base, ext)
	
	// 확장자에서 점(.) 제거
	ext = strings.TrimPrefix(ext, ".")

	if dir == "." {
		// 디렉토리가 없는 경우
		if ext != "" {
			return name + "_" + ext
		}
		return name
	}

	// 디렉토리 구분자를 언더스코어로 변환
	dirParts := strings.Split(dir, string(filepath.Separator))
	if ext != "" {
		return strings.Join(append(dirParts, name, ext), "_")
	}
	return strings.Join(append(dirParts, name), "_")
}

func (m *Metadata) AddFile(path string) error {
	// Check if path is already absolute
	var fullPath string
	var relativePath string
	
	if filepath.IsAbs(path) {
		fullPath = path
		// Extract path after .cursor/rules/
		if idx := strings.Index(path, ".cursor/rules/"); idx != -1 {
			relativePath = path[idx+len(".cursor/rules/"):]
		} else {
			return fmt.Errorf("path is not within .cursor/rules/ directory: %s", path)
		}
	} else {
		// For relative paths
		dir, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("failed to get working directory: %w", err)
		}
		fullPath = filepath.Join(dir, ".cursor/rules", path)
		relativePath = path
	}

	// Open file
	file, err := os.Open(fullPath)
	if err != nil {
		return fmt.Errorf("failed to open file %s: %w", path, err)
	}
	defer file.Close()

	// Get file info
	info, err := file.Stat()
	if err != nil {
		return fmt.Errorf("failed to get file info %s: %w", path, err)
	}

	// Calculate MD5 hash
	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		return fmt.Errorf("failed to calculate hash %s: %w", path, err)
	}

	// Generate Gist file name (using relative path)
	gistName := convertToGistName(relativePath)

	metadata := FileMetadata{
		Path:     relativePath,
		GistName: gistName,
		Size:     info.Size(),
		MD5:      hex.EncodeToString(hash.Sum(nil)),
	}
	m.Files = append(m.Files, metadata)

	// Update directory structure (using relative path)
	m.updateStructure(relativePath)

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
	// 입력된 경로를 상대 경로로 변환
	var relativePath string
	if filepath.IsAbs(path) {
		if idx := strings.Index(path, ".cursor/rules/"); idx != -1 {
			relativePath = path[idx+len(".cursor/rules/"):]
		} else {
			// .cursor/rules/ 경로가 없는 경우 기본 변환 사용
			return convertToGistName(filepath.Base(path))
		}
	} else {
		relativePath = path
	}

	// 상대 경로로 파일 메타데이터 찾기
	for _, file := range m.Files {
		if file.Path == relativePath {
			return file.GistName
		}
	}

	// 메타데이터에서 찾지 못한 경우 기본 변환 사용
	return convertToGistName(relativePath)
}

// PreviewMetadata generates metadata for the given file paths.
// This allows previewing metadata without actually uploading files.
func PreviewMetadata(paths []string) (*Metadata, error) {
	meta := NewMetadata()

	for _, path := range paths {
		if err := meta.AddFile(path); err != nil {
			return nil, fmt.Errorf("failed to generate metadata (%s): %w", path, err)
		}
	}

	return meta, nil
}

// PreviewMetadataFromDir finds .mdc files in the specified directory and generates metadata.
func PreviewMetadataFromDir(dir string) (*Metadata, error) {
	var paths []string

	// Find .mdc files
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(info.Name(), ".mdc") {
			paths = append(paths, path)
		}
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to scan directory: %w", err)
	}

	return PreviewMetadata(paths)
}

// PreviewMetadataFromWorkingDir generates metadata from the current working directory's .cursor/rules/ path
// and returns it in meta.json format.
func PreviewMetadataFromWorkingDir() (*Metadata, error) {
	// Check current working directory
	workDir, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("failed to get working directory: %w", err)
	}

	// Create .cursor/rules path
	rulesDir := filepath.Join(workDir, ".cursor", "rules")

	// Check if directory exists
	if _, err := os.Stat(rulesDir); os.IsNotExist(err) {
		return nil, fmt.Errorf(".cursor/rules directory not found: %s", rulesDir)
	}

	return PreviewMetadataFromDir(rulesDir)
}

// WriteMetadataPreview는 메타데이터를 meta.json 형식으로 반환합니다.
func (m *Metadata) WriteMetadataPreview() ([]byte, error) {
	return json.MarshalIndent(struct {
		SchemaVersion string             `json:"schema_version"`
		CLIVersion    string            `json:"cli_version"`
		UpdatedAt     time.Time         `json:"updated_at"`
		Structure     DirectoryStructure `json:"structure"`
		Files         []FileMetadata    `json:"files"`
	}{
		SchemaVersion: m.SchemaVersion,
		CLIVersion:    m.CLIVersion,
		UpdatedAt:     m.UpdatedAt,
		Structure:     m.Structure,
		Files:         m.Files,
	}, "", "  ")
}

// String은 메타데이터를 보기 좋은 형식으로 출력합니다.
func (m *Metadata) String() string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("Schema Version: %s\n", m.SchemaVersion))
	sb.WriteString(fmt.Sprintf("CLI Version: %s\n", m.CLIVersion))
	sb.WriteString(fmt.Sprintf("Updated At: %s\n\n", m.UpdatedAt.Format(time.RFC3339)))
	
	sb.WriteString("Files:\n")
	for _, file := range m.Files {
		sb.WriteString(fmt.Sprintf("- %s\n", file.Path))
		sb.WriteString(fmt.Sprintf("  → Gist Name: %s\n", file.GistName))
		sb.WriteString(fmt.Sprintf("  → Size: %d bytes\n", file.Size))
		sb.WriteString(fmt.Sprintf("  → MD5: %s\n", file.MD5))
	}

	return sb.String()
} 