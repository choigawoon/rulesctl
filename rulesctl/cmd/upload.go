package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/choigawoon/rulesctl/internal/fileutils"
	"github.com/choigawoon/rulesctl/internal/gist"
	"github.com/spf13/cobra"
)

var (
	forceUpload bool
)

var uploadCmd = &cobra.Command{
	Use:   "upload [name]",
	Short: "로컬 규칙을 GIST에 업로드",
	Long: `로컬 .cursor/rules 디렉토리의 규칙 파일들을 GIST에 업로드합니다.
규칙 세트 이름은 따옴표로 감싸서 지정합니다.`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]
		
		// 1. 로컬 규칙 파일 수집
		rules, err := fileutils.ListLocalRules()
		if err != nil {
			return fmt.Errorf("로컬 규칙 파일 수집 실패: %v", err)
		}

		if len(rules) == 0 {
			return fmt.Errorf("업로드할 규칙 파일이 없습니다. .cursor/rules 디렉토리에 규칙 파일을 추가해주세요")
		}

		// 2. Gist 클라이언트 초기화
		client, err := gist.NewClient()
		if err != nil {
			return fmt.Errorf("Gist 클라이언트 초기화 실패: %v", err)
		}

		// 3. 파일 내용 읽기 및 Gist 파일 생성
		files := make(map[string]gist.File)
		for path, _ := range rules {
			content, err := os.ReadFile(filepath.Join(fileutils.RulesDirName, path))
			if err != nil {
				return fmt.Errorf("파일 읽기 실패 %s: %v", path, err)
			}
			files[path] = gist.File{
				Content: string(content),
			}
		}

		// 4. Gist 생성 또는 업데이트
		gistID, err := client.CreateOrUpdateGist(name, files, forceUpload)
		if err != nil {
			return fmt.Errorf("Gist 업로드 실패: %v", err)
		}

		fmt.Printf("규칙이 성공적으로 업로드되었습니다. Gist ID: %s\n", gistID)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(uploadCmd)
	uploadCmd.Flags().BoolVarP(&forceUpload, "force", "f", false, "충돌이 있는 경우 강제로 업로드")
} 