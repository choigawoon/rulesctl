package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/choigawoon/rulesctl/internal/fileutils"
	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "규칙 디렉토리 초기화",
	Long: `.cursor/rules 디렉토리를 생성하고 샘플 규칙 파일을 추가합니다.
이 명령어는 rulesctl을 처음 사용할 때 유용합니다.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// 1. 규칙 디렉토리 생성
		if err := fileutils.EnsureRulesDir(); err != nil {
			return fmt.Errorf("규칙 디렉토리 생성 실패: %v", err)
		}

		// 2. 샘플 파일 경로 생성
		rulesDir, err := fileutils.GetRulesDirPath()
		if err != nil {
			return fmt.Errorf("규칙 디렉토리 경로 확인 실패: %v", err)
		}

		// 3. 샘플 디렉토리 및 파일 생성
		sampleFiles := map[string]string{
			"example/hello.mdc": "# 예제 규칙 파일\n\n이것은 예제 규칙 파일입니다.\n사용자 정의 규칙을 작성하여 Cursor AI 동작을 맞춤화할 수 있습니다.",
			"python/linting.mdc": "# Python 린팅 규칙\n\n항상 PEP8 스타일 가이드를 따르세요.\n들여쓰기는 공백 4칸을 사용하고, 한 줄은 79자를 넘지 않도록 합니다.",
			"javascript/coding.mdc": "# JavaScript 코딩 규칙\n\n- 세미콜론을 항상 사용하세요.\n- 변수 선언에는 const와 let을 사용하세요. var는 사용하지 마세요.\n- 함수는 화살표 함수를 선호합니다.",
		}

		for path, content := range sampleFiles {
			fullPath := filepath.Join(rulesDir, path)
			dirPath := filepath.Dir(fullPath)

			// 디렉토리 생성
			if err := os.MkdirAll(dirPath, 0755); err != nil {
				return fmt.Errorf("디렉토리 생성 실패 %s: %v", dirPath, err)
			}

			// 파일이 이미 존재하는지 확인
			if _, err := os.Stat(fullPath); !os.IsNotExist(err) {
				fmt.Printf("파일이 이미 존재합니다: %s\n", path)
				continue
			}

			// 파일 생성
			if err := os.WriteFile(fullPath, []byte(content), 0644); err != nil {
				return fmt.Errorf("파일 생성 실패 %s: %v", path, err)
			}
			fmt.Printf("샘플 파일 생성: %s\n", path)
		}

		fmt.Printf("\n초기화가 완료되었습니다.\n")
		fmt.Printf("생성된 디렉토리: %s\n", rulesDir)
		fmt.Printf("\n다음 명령어로 규칙을 업로드할 수 있습니다:\n")
		fmt.Printf("  rulesctl upload \"my-rules\"\n")

		return nil
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
} 