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
		
		// 0. 규칙 디렉토리 확인 및 생성
		if err := fileutils.EnsureRulesDir(); err != nil {
			return fmt.Errorf("규칙 디렉토리 생성 실패: %v", err)
		}
		
		// 1. 로컬 규칙 파일 수집
		rules, err := fileutils.ListLocalRules()
		if err != nil {
			// 파일이 없는 경우 안내 메시지 출력
			rulesDir, _ := fileutils.GetRulesDirPath()
			fmt.Printf("안내: %v\n", err)
			fmt.Printf("현재 규칙 디렉토리: %s\n", rulesDir)
			fmt.Printf("다음 명령어로 규칙 파일을 추가할 수 있습니다:\n")
			fmt.Printf("  mkdir -p %s/python\n", rulesDir)
			fmt.Printf("  echo \"Python 린팅 규칙\" > %s/python/linting.mdc\n", rulesDir)
			return fmt.Errorf("업로드할 규칙 파일이 없습니다")
		}

		// 2. 메타데이터 생성
		metadata := gist.NewMetadata()
		files := make(map[string]gist.File)

		// 3. 파일 내용 읽기 및 메타데이터 수집
		rulesDir, err := fileutils.GetRulesDirPath()
		if err != nil {
			return fmt.Errorf("규칙 디렉토리 경로 확인 실패: %v", err)
		}

		for path := range rules {
			fullPath := filepath.Join(rulesDir, path)
			content, err := os.ReadFile(fullPath)
			if err != nil {
				return fmt.Errorf("파일 읽기 실패 %s: %v", path, err)
			}

			// 메타데이터에 파일 추가
			if err := metadata.AddFile(fullPath); err != nil {
				return fmt.Errorf("메타데이터 추가 실패 %s: %v", path, err)
			}

			// Gist 파일 맵에 변환된 이름으로 추가
			gistName := metadata.GetGistName(path)
			files[gistName] = gist.File{
				Content: string(content),
			}
		}

		// 4. 메타데이터 파일 추가
		metaContent, err := metadata.ToJSON()
		if err != nil {
			return fmt.Errorf("메타데이터 JSON 생성 실패: %v", err)
		}
		files[gist.MetaFileName] = gist.File{
			Content: string(metaContent),
		}

		// 5. Gist 클라이언트 초기화
		client, err := gist.NewClient()
		if err != nil {
			return fmt.Errorf("Gist 클라이언트 초기화 실패: %v", err)
		}

		// 6. Gist 생성 또는 업데이트
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