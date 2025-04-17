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

	// 실제 환경 테스트
	t.Run("실제 환경 테스트", func(t *testing.T) {
		// 1. 실제 .cursor/rules 디렉토리 구조 테스트
		t.Run("디렉토리 구조", func(t *testing.T) {
			// 실제 디렉토리 구조 생성
			structure := map[string]string{
				"python/linting.mdc":    "Python 린팅 규칙",
				"python/testing.mdc":    "Python 테스트 규칙",
				"database/postgres.mdc": "PostgreSQL 규칙",
			}

			for path, content := range structure {
				err := SaveRuleFile(path, []byte(content))
				if err != nil {
					t.Errorf("디렉토리 구조 생성 실패 %s: %v", path, err)
				}
			}

			// 디렉토리 구조 확인
			rules, err := ListLocalRules()
			if err != nil {
				t.Errorf("디렉토리 구조 확인 실패: %v", err)
			}

			for path := range structure {
				if _, exists := rules[path]; !exists {
					t.Errorf("예상된 파일이 없음: %s", path)
				}
			}
		})

		// 2. 권한 테스트
		t.Run("파일 권한", func(t *testing.T) {
			// 읽기 전용 파일 생성
			path := "readonly.mdc"
			err := SaveRuleFile(path, []byte("readonly content"))
			if err != nil {
				t.Errorf("읽기 전용 파일 생성 실패: %v", err)
			}

			// 파일 권한 변경
			rulesDir, _ := GetRulesDirPath()
			fullPath := filepath.Join(rulesDir, path)
			err = os.Chmod(fullPath, 0444) // 읽기 전용
			if err != nil {
				t.Errorf("파일 권한 변경 실패: %v", err)
			}

			// 읽기 시도
			_, err = os.ReadFile(fullPath)
			if err != nil {
				t.Errorf("읽기 전용 파일 읽기 실패: %v", err)
			}

			// 쓰기 시도 (실패해야 함)
			err = os.WriteFile(fullPath, []byte("new content"), 0644)
			if err == nil {
				t.Errorf("읽기 전용 파일에 쓰기가 성공했음 (실패해야 함)")
			}
		})

		// 3. 경로 처리 테스트
		t.Run("경로 처리", func(t *testing.T) {
			// 상대 경로 테스트
			relPath := "path/test.mdc"
			err := SaveRuleFile(relPath, []byte("relative path test"))
			if err != nil {
				t.Errorf("상대 경로 파일 생성 실패: %v", err)
			}

			// 절대 경로 테스트
			rulesDir, _ := GetRulesDirPath()
			absPath := filepath.Join(rulesDir, "abs/test.mdc")
			err = os.MkdirAll(filepath.Dir(absPath), 0755)
			if err != nil {
				t.Errorf("절대 경로 디렉토리 생성 실패: %v", err)
			}
			err = os.WriteFile(absPath, []byte("absolute path test"), 0644)
			if err != nil {
				t.Errorf("절대 경로 파일 생성 실패: %v", err)
			}

			// 경로 정규화 테스트
			normalizedPath := filepath.Clean("path/../path/test.mdc")
			if normalizedPath != "path/test.mdc" {
				t.Errorf("경로 정규화 실패: %s", normalizedPath)
			}
		})

		// 4. 심볼릭 링크 테스트
		t.Run("심볼릭 링크", func(t *testing.T) {
			// 원본 파일 생성
			originalPath := "original.mdc"
			err := SaveRuleFile(originalPath, []byte("original content"))
			if err != nil {
				t.Errorf("원본 파일 생성 실패: %v", err)
			}

			// 심볼릭 링크 생성
			rulesDir, _ := GetRulesDirPath()
			originalFullPath := filepath.Join(rulesDir, originalPath)
			linkPath := filepath.Join(rulesDir, "link.mdc")
			err = os.Symlink(originalFullPath, linkPath)
			if err != nil {
				t.Errorf("심볼릭 링크 생성 실패: %v", err)
			}

			// 심볼릭 링크를 통한 파일 접근
			content, err := os.ReadFile(linkPath)
			if err != nil {
				t.Errorf("심볼릭 링크를 통한 파일 읽기 실패: %v", err)
			}
			if string(content) != "original content" {
				t.Errorf("심볼릭 링크 내용 불일치")
			}

			// 심볼릭 링크 경로 처리
			linkInfo, err := os.Lstat(linkPath)
			if err != nil {
				t.Errorf("심볼릭 링크 정보 조회 실패: %v", err)
			}
			if linkInfo.Mode()&os.ModeSymlink == 0 {
				t.Errorf("심볼릭 링크가 아님")
			}
		})
	})
} 