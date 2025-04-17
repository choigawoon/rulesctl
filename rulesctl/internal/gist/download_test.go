package gist

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

func TestFetchGist(t *testing.T) {
	// 테스트 서버 설정
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 인증 헤더 확인
		token := r.Header.Get("Authorization")
		if token != "token test-token" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		// 테스트 Gist 데이터
		gist := Gist{
			ID:          "test-gist",
			Description: "Test Gist",
			Files: map[string]struct {
				Filename string `json:"filename"`
				Type     string `json:"type"`
				Language string `json:"language"`
				RawURL   string `json:"raw_url"`
				Size     int    `json:"size"`
			}{
				MetaFileName: {
					Filename: MetaFileName,
					Type:     "text/plain",
					RawURL:   "http://example.com/raw",
					Size:     100,
				},
			},
		}

		json.NewEncoder(w).Encode(gist)
	}))
	defer server.Close()

	// baseURL 임시 변경
	originalBaseURL := baseURL
	baseURL = server.URL
	defer func() { baseURL = originalBaseURL }()

	// 테스트 실행
	gist, err := FetchGist("test-token", "test-gist")
	if err != nil {
		t.Errorf("FetchGist 실패: %v", err)
	}

	if gist.ID != "test-gist" {
		t.Errorf("잘못된 Gist ID: got %s, want test-gist", gist.ID)
	}

	if _, exists := gist.Files[MetaFileName]; !exists {
		t.Error("메타데이터 파일이 없습니다")
	}
}

func TestParseMetadataFromGist(t *testing.T) {
	// 테스트 메타데이터
	testMeta := &Metadata{
		SchemaVersion: "1.0.0",
		CLIVersion:    "0.1.0",
		Files: []FileMetadata{
			{
				Path:     "test/file.mdc",
				GistName: "test_file_mdc",
				Size:     100,
				MD5:      "test-hash",
			},
		},
	}

	// JSON 직렬화
	content, err := json.Marshal(testMeta)
	if err != nil {
		t.Fatalf("메타데이터 직렬화 실패: %v", err)
	}

	// 테스트 실행
	meta, err := ParseMetadataFromGist(string(content))
	if err != nil {
		t.Errorf("ParseMetadataFromGist 실패: %v", err)
	}

	if meta.SchemaVersion != testMeta.SchemaVersion {
		t.Errorf("잘못된 스키마 버전: got %s, want %s", meta.SchemaVersion, testMeta.SchemaVersion)
	}

	if len(meta.Files) != 1 {
		t.Errorf("잘못된 파일 수: got %d, want 1", len(meta.Files))
	}

	if meta.Files[0].Path != testMeta.Files[0].Path {
		t.Errorf("잘못된 파일 경로: got %s, want %s", meta.Files[0].Path, testMeta.Files[0].Path)
	}
}

func TestCheckConflicts(t *testing.T) {
	// 임시 디렉토리 생성
	tmpDir, err := os.MkdirTemp("", "rulesctl-test-*")
	if err != nil {
		t.Fatalf("임시 디렉토리 생성 실패: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// 현재 디렉토리 저장
	originalWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("현재 디렉토리 확인 실패: %v", err)
	}
	defer os.Chdir(originalWd)

	// 임시 디렉토리로 이동
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("임시 디렉토리로 이동 실패: %v", err)
	}

	// 테스트 파일 생성
	testPath := filepath.Join(".cursor", "rules", "test", "file.mdc")
	if err := os.MkdirAll(filepath.Dir(testPath), 0755); err != nil {
		t.Fatalf("테스트 디렉토리 생성 실패: %v", err)
	}
	if err := os.WriteFile(testPath, []byte("test"), 0644); err != nil {
		t.Fatalf("테스트 파일 생성 실패: %v", err)
	}

	// 테스트 메타데이터
	meta := &Metadata{
		Files: []FileMetadata{
			{Path: "test/file.mdc"},
			{Path: "test/nonexistent.mdc"},
		},
	}

	// 테스트 실행
	conflicts, err := CheckConflicts(meta)
	if err != nil {
		t.Errorf("CheckConflicts 실패: %v", err)
	}

	if len(conflicts) != 1 {
		t.Errorf("잘못된 충돌 수: got %d, want 1", len(conflicts))
	}

	if conflicts[0] != "test/file.mdc" {
		t.Errorf("잘못된 충돌 파일: got %s, want test/file.mdc", conflicts[0])
	}
}

func TestDownloadFiles(t *testing.T) {
	// 테스트 서버 설정
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("test content"))
	}))
	defer server.Close()

	// 임시 디렉토리 생성
	tmpDir, err := os.MkdirTemp("", "rulesctl-test-*")
	if err != nil {
		t.Fatalf("임시 디렉토리 생성 실패: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// 현재 디렉토리 저장
	originalWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("현재 디렉토리 확인 실패: %v", err)
	}
	defer os.Chdir(originalWd)

	// 임시 디렉토리로 이동
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("임시 디렉토리로 이동 실패: %v", err)
	}

	// 테스트 메타데이터
	meta := &Metadata{
		Files: []FileMetadata{
			{
				Path:     "test/file.mdc",
				GistName: "test_file_mdc",
			},
		},
	}

	// 테스트 Gist
	testGist := &Gist{
		Files: map[string]struct {
			Filename string `json:"filename"`
			Type     string `json:"type"`
			Language string `json:"language"`
			RawURL   string `json:"raw_url"`
			Size     int    `json:"size"`
		}{
			"test_file_mdc": {
				Filename: "test_file_mdc",
				RawURL:   server.URL,
			},
		},
	}

	// baseURL 임시 변경
	originalBaseURL := baseURL
	baseURL = server.URL
	defer func() { baseURL = originalBaseURL }()

	// 테스트 실행
	if err := DownloadFiles("test-token", "test-gist", meta, true); err != nil {
		t.Errorf("DownloadFiles 실패: %v", err)
	}

	// 파일 확인
	content, err := os.ReadFile(filepath.Join(".cursor", "rules", "test", "file.mdc"))
	if err != nil {
		t.Errorf("다운로드된 파일 읽기 실패: %v", err)
	}

	if string(content) != "test content" {
		t.Errorf("잘못된 파일 내용: got %s, want test content", string(content))
	}
} 