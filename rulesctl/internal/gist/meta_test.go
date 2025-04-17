package gist

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestGetGistName(t *testing.T) {
	meta := &Metadata{
		Files: []FileMetadata{
			{
				Path:     "test/example.txt",
				GistName: "example_txt_12345",
			},
		},
	}

	tests := []struct {
		name     string
		path     string
		expected string
	}{
		{
			name:     "절대 경로 (.cursor/rules/ 포함)",
			path:     "/Users/test/.cursor/rules/test/example.txt",
			expected: "example_txt_12345",
		},
		{
			name:     "절대 경로 (.cursor/rules/ 미포함)",
			path:     "/Users/test/example.txt",
			expected: "example_txt",
		},
		{
			name:     "상대 경로 (메타데이터에 존재)",
			path:     "test/example.txt",
			expected: "example_txt_12345",
		},
		{
			name:     "상대 경로 (메타데이터에 없음)",
			path:     "new/file.txt",
			expected: "new_file_txt",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := meta.GetGistName(tt.path)
			if result != tt.expected {
				t.Errorf("GetGistName(%s) = %s; want %s", tt.path, result, tt.expected)
			}
		})
	}
}

func TestPreviewMetadata(t *testing.T) {
	// 임시 디렉토리 생성
	tempDir := t.TempDir()
	rulesDir := filepath.Join(tempDir, ".cursor", "rules")

	// 테스트 파일 생성
	testFiles := []struct {
		path    string
		content string
	}{
		{
			path:    "python/linting.mdc",
			content: "파이썬 린팅 규칙",
		},
		{
			path:    "database/postgres.mdc",
			content: "PostgreSQL 규칙",
		},
	}

	for _, tf := range testFiles {
		fullPath := filepath.Join(rulesDir, tf.path)
		err := os.MkdirAll(filepath.Dir(fullPath), 0755)
		if err != nil {
			t.Fatalf("디렉토리 생성 실패: %v", err)
		}
		err = os.WriteFile(fullPath, []byte(tf.content), 0644)
		if err != nil {
			t.Fatalf("파일 생성 실패: %v", err)
		}
	}

	// 개별 파일 메타데이터 테스트
	t.Run("개별 파일 메타데이터", func(t *testing.T) {
		paths := []string{
			filepath.Join(rulesDir, "python/linting.mdc"),
			filepath.Join(rulesDir, "database/postgres.mdc"),
		}
		
		meta, err := PreviewMetadata(paths)
		if err != nil {
			t.Fatalf("메타데이터 생성 실패: %v", err)
		}

		if len(meta.Files) != 2 {
			t.Errorf("예상 파일 수: 2, 실제: %d", len(meta.Files))
		}

		// Gist 이름 확인
		expectedNames := map[string]string{
			"python/linting.mdc":    "python_linting_mdc",
			"database/postgres.mdc": "database_postgres_mdc",
		}

		for _, file := range meta.Files {
			expected, ok := expectedNames[file.Path]
			if !ok {
				t.Errorf("예상치 못한 파일 경로: %s", file.Path)
				continue
			}
			if file.GistName != expected {
				t.Errorf("파일 %s의 Gist 이름이 잘못됨. 예상: %s, 실제: %s", 
					file.Path, expected, file.GistName)
			}
		}
	})

	// 디렉토리 메타데이터 테스트
	t.Run("디렉토리 메타데이터", func(t *testing.T) {
		meta, err := PreviewMetadataFromDir(rulesDir)
		if err != nil {
			t.Fatalf("디렉토리 메타데이터 생성 실패: %v", err)
		}

		if len(meta.Files) != 2 {
			t.Errorf("예상 파일 수: 2, 실제: %d", len(meta.Files))
		}

		// 출력 형식 확인
		output := meta.String()
		expectedSubstrings := []string{
			"Schema Version:",
			"CLI Version:",
			"Updated At:",
			"Files:",
			"python/linting.mdc",
			"database/postgres.mdc",
			"→ Gist Name:",
			"→ Size:",
			"→ MD5:",
		}

		for _, substr := range expectedSubstrings {
			if !strings.Contains(output, substr) {
				t.Errorf("출력에서 '%s'를 찾을 수 없음", substr)
			}
		}
	})
}

func TestPreviewMetadataFromWorkingDir(t *testing.T) {
	// 현재 작업 디렉토리 저장
	originalWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("현재 작업 디렉토리 확인 실패: %v", err)
	}
	defer os.Chdir(originalWd)

	// 임시 디렉토리 생성 및 이동
	tempDir := t.TempDir()
	if err := os.Chdir(tempDir); err != nil {
		t.Fatalf("임시 디렉토리로 이동 실패: %v", err)
	}

	// .cursor/rules 디렉토리가 없는 경우 테스트
	t.Run("디렉토리 없음", func(t *testing.T) {
		_, err := PreviewMetadataFromWorkingDir()
		if err == nil {
			t.Error("에러가 발생해야 하지만 발생하지 않음")
		}
		if !strings.Contains(err.Error(), ".cursor/rules 디렉토리를 찾을 수 없습니다") {
			t.Errorf("예상치 못한 에러: %v", err)
		}
	})

	// .cursor/rules 디렉토리 및 테스트 파일 생성
	rulesDir := filepath.Join(tempDir, ".cursor", "rules")
	testFiles := []struct {
		path    string
		content string
	}{
		{
			path:    "python/linting.mdc",
			content: "파이썬 린팅 규칙",
		},
		{
			path:    "database/postgres.mdc",
			content: "PostgreSQL 규칙",
		},
	}

	for _, tf := range testFiles {
		fullPath := filepath.Join(rulesDir, tf.path)
		err := os.MkdirAll(filepath.Dir(fullPath), 0755)
		if err != nil {
			t.Fatalf("디렉토리 생성 실패: %v", err)
		}
		err = os.WriteFile(fullPath, []byte(tf.content), 0644)
		if err != nil {
			t.Fatalf("파일 생성 실패: %v", err)
		}
	}

	// 메타데이터 생성 테스트
	t.Run("메타데이터 생성", func(t *testing.T) {
		meta, err := PreviewMetadataFromWorkingDir()
		if err != nil {
			t.Fatalf("메타데이터 생성 실패: %v", err)
		}

		// 파일 수 확인
		if len(meta.Files) != 2 {
			t.Errorf("예상 파일 수: 2, 실제: %d", len(meta.Files))
		}

		// meta.json 형식으로 변환
		jsonData, err := meta.WriteMetadataPreview()
		if err != nil {
			t.Fatalf("JSON 변환 실패: %v", err)
		}

		// JSON 형식 검증
		var preview struct {
			SchemaVersion string    `json:"schema_version"`
			CLIVersion    string    `json:"cli_version"`
			UpdatedAt     time.Time `json:"updated_at"`
			Files         []FileMetadata
		}
		if err := json.Unmarshal(jsonData, &preview); err != nil {
			t.Fatalf("JSON 파싱 실패: %v", err)
		}

		// 필수 필드 확인
		if preview.SchemaVersion == "" {
			t.Error("SchemaVersion이 비어있음")
		}
		if preview.CLIVersion == "" {
			t.Error("CLIVersion이 비어있음")
		}
		if preview.UpdatedAt.IsZero() {
			t.Error("UpdatedAt이 설정되지 않음")
		}
		if len(preview.Files) != 2 {
			t.Errorf("예상 파일 수: 2, 실제: %d", len(preview.Files))
		}
	})
} 