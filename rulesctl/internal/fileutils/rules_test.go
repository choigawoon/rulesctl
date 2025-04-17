package fileutils

import (
	"os"
	"path/filepath"
	"testing"
)

func TestFileUtils(t *testing.T) {
	// 테스트용 임시 디렉토리 생성
	tempDir, err := os.MkdirTemp("", "rulesctl-test-*")
	if err != nil {
		t.Fatalf("임시 디렉토리 생성 실패: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// 작업 디렉토리를 임시 디렉토리로 변경
	oldDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("현재 작업 디렉토리 확인 실패: %v", err)
	}
	err = os.Chdir(tempDir)
	if err != nil {
		t.Fatalf("작업 디렉토리 변경 실패: %v", err)
	}
	defer os.Chdir(oldDir)

	t.Run("GetRulesDirPath", func(t *testing.T) {
		path, err := GetRulesDirPath()
		if err != nil {
			t.Errorf("GetRulesDirPath 실패: %v", err)
		}

		// 실제 경로 얻기 (심볼릭 링크 해결)
		realTempDir, err := filepath.EvalSymlinks(tempDir)
		if err != nil {
			t.Fatalf("심볼릭 링크 해결 실패: %v", err)
		}

		// 두 경로를 정규화하여 비교
		expected := filepath.Clean(filepath.Join(realTempDir, RulesDirName))
		actual := filepath.Clean(path)
		if actual != expected {
			t.Errorf("예상 경로: %s, 실제 경로: %s", expected, actual)
		}
	})

	t.Run("EnsureRulesDir", func(t *testing.T) {
		err := EnsureRulesDir()
		if err != nil {
			t.Errorf("EnsureRulesDir 실패: %v", err)
		}

		// 디렉토리가 실제로 생성되었는지 확인
		path, _ := GetRulesDirPath()
		if _, err := os.Stat(path); os.IsNotExist(err) {
			t.Errorf("규칙 디렉토리가 생성되지 않음: %s", path)
		}
	})

	t.Run("SaveRuleFile", func(t *testing.T) {
		// 테스트용 규칙 파일 생성
		content := []byte("test rule content")
		err := SaveRuleFile("test/test.mdc", content)
		if err != nil {
			t.Errorf("SaveRuleFile 실패: %v", err)
		}

		// 파일이 실제로 생성되었는지 확인
		path, _ := GetRulesDirPath()
		fullPath := filepath.Join(path, "test/test.mdc")
		if _, err := os.Stat(fullPath); os.IsNotExist(err) {
			t.Errorf("규칙 파일이 생성되지 않음: %s", fullPath)
		}
	})

	t.Run("ListLocalRules", func(t *testing.T) {
		// 여러 테스트 파일 생성
		files := map[string][]byte{
			"test1.mdc": []byte("test1"),
			"test2.mdc": []byte("test2"),
			"sub/test3.mdc": []byte("test3"),
		}

		for path, content := range files {
			err := SaveRuleFile(path, content)
			if err != nil {
				t.Errorf("테스트 파일 생성 실패 %s: %v", path, err)
			}
		}

		// 파일 목록 조회
		rules, err := ListLocalRules()
		if err != nil {
			t.Errorf("ListLocalRules 실패: %v", err)
		}

		// 모든 파일이 목록에 포함되어 있는지 확인
		for path := range files {
			if _, exists := rules[path]; !exists {
				t.Errorf("파일이 목록에 없음: %s", path)
			}
		}
	})

	t.Run("DeleteRuleFile", func(t *testing.T) {
		// 테스트 파일 생성
		path := "test-delete.mdc"
		err := SaveRuleFile(path, []byte("test"))
		if err != nil {
			t.Errorf("테스트 파일 생성 실패: %v", err)
		}

		// 파일 삭제
		err = DeleteRuleFile(path)
		if err != nil {
			t.Errorf("DeleteRuleFile 실패: %v", err)
		}

		// 파일이 실제로 삭제되었는지 확인
		rulesDir, _ := GetRulesDirPath()
		fullPath := filepath.Join(rulesDir, path)
		if _, err := os.Stat(fullPath); !os.IsNotExist(err) {
			t.Errorf("파일이 삭제되지 않음: %s", fullPath)
		}
	})
} 