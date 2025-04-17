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
		CLIVersion:    "0.1.0", // TODO: 버전 관리
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
	// 먼저 path가 이미 절대 경로인지 확인
	var fullPath string
	var relativePath string
	
	if filepath.IsAbs(path) {
		fullPath = path
		// 절대 경로에서 .cursor/rules/ 이후의 경로만 추출
		if idx := strings.Index(path, ".cursor/rules/"); idx != -1 {
			relativePath = path[idx+len(".cursor/rules/"):]
		} else {
			return fmt.Errorf("경로가 .cursor/rules/ 디렉토리 내에 있지 않습니다: %s", path)
		}
	} else {
		// 상대 경로인 경우
		dir, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("작업 디렉토리 확인 실패: %w", err)
		}
		fullPath = filepath.Join(dir, ".cursor/rules", path)
		relativePath = path
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

	// Gist 파일 이름 생성 (상대 경로 사용)
	gistName := convertToGistName(relativePath)

	metadata := FileMetadata{
		Path:     relativePath,
		GistName: gistName,
		Size:     info.Size(),
		MD5:      hex.EncodeToString(hash.Sum(nil)),
	}
	m.Files = append(m.Files, metadata)

	// 디렉토리 구조 업데이트 (상대 경로 사용)
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

// PreviewMetadata는 주어진 경로의 파일들에 대한 메타데이터를 생성합니다.
// 실제로 파일을 업로드하지 않고 메타데이터만 미리 확인할 수 있습니다.
func PreviewMetadata(paths []string) (*Metadata, error) {
	meta := NewMetadata()

	for _, path := range paths {
		if err := meta.AddFile(path); err != nil {
			return nil, fmt.Errorf("메타데이터 생성 실패 (%s): %w", path, err)
		}
	}

	return meta, nil
}

// PreviewMetadataFromDir는 지정된 디렉토리에서 .mdc 파일들을 찾아 메타데이터를 생성합니다.
func PreviewMetadataFromDir(dir string) (*Metadata, error) {
	var paths []string

	// .mdc 파일 찾기
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
		return nil, fmt.Errorf("디렉토리 탐색 실패: %w", err)
	}

	return PreviewMetadata(paths)
}

// PreviewMetadataFromWorkingDir는 현재 작업 디렉토리의 .cursor/rules/ 경로에서
// 메타데이터를 생성하고 meta.json 형식으로 반환합니다.
func PreviewMetadataFromWorkingDir() (*Metadata, error) {
	// 현재 작업 디렉토리 확인
	workDir, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("작업 디렉토리 확인 실패: %w", err)
	}

	// .cursor/rules 경로 생성
	rulesDir := filepath.Join(workDir, ".cursor", "rules")

	// 디렉토리 존재 여부 확인
	if _, err := os.Stat(rulesDir); os.IsNotExist(err) {
		return nil, fmt.Errorf(".cursor/rules 디렉토리를 찾을 수 없습니다: %s", rulesDir)
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