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
	preview     bool
)

var uploadCmd = &cobra.Command{
	Use:   "upload [name]",
	Short: "로컬 규칙을 GIST에 업로드",
	Long: `로컬 .cursor/rules 디렉토리의 규칙 파일들을 GIST에 업로드합니다.
규칙 세트 이름은 따옴표로 감싸서 지정합니다.

--preview 플래그를 사용하면 실제 업로드 없이 메타데이터를 미리 확인할 수 있습니다.`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		// 0. 규칙 디렉토리 확인 및 생성
		if err := fileutils.EnsureRulesDir(); err != nil {
			return fmt.Errorf("규칙 디렉토리 생성 실패: %v", err)
		}

		// 1. 메타데이터 미리보기
		meta, err := gist.PreviewMetadataFromWorkingDir()
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

		// 파일이 없는 경우 처리
		if len(meta.Files) == 0 {
			rulesDir, _ := fileutils.GetRulesDirPath()
			fmt.Printf("현재 규칙 디렉토리: %s\n", rulesDir)
			fmt.Printf("다음 명령어로 규칙 파일을 추가할 수 있습니다:\n")
			fmt.Printf("  mkdir -p %s/python\n", rulesDir)
			fmt.Printf("  echo \"Python 린팅 규칙\" > %s/python/linting.mdc\n", rulesDir)
			return fmt.Errorf("업로드할 규칙 파일이 없습니다")
		}

		// 미리보기 모드인 경우 메타데이터만 출력하고 종료
		if preview {
			jsonData, err := meta.WriteMetadataPreview()
			if err != nil {
				return fmt.Errorf("메타데이터 JSON 생성 실패: %v", err)
			}
			fmt.Println(string(jsonData))
			return nil
		}

		// 2. 파일 내용 읽기 및 Gist 파일 맵 생성
		files := make(map[string]gist.File)
		rulesDir, err := fileutils.GetRulesDirPath()
		if err != nil {
			return fmt.Errorf("규칙 디렉토리 경로 확인 실패: %v", err)
		}

		for _, fileInfo := range meta.Files {
			fullPath := filepath.Join(rulesDir, fileInfo.Path)
			content, err := os.ReadFile(fullPath)
			if err != nil {
				return fmt.Errorf("파일 읽기 실패 %s: %v", fileInfo.Path, err)
			}

			files[fileInfo.GistName] = gist.File{
				Content: string(content),
			}
		}

		// 3. 메타데이터 파일 추가
		metaContent, err := meta.ToJSON()
		if err != nil {
			return fmt.Errorf("메타데이터 JSON 생성 실패: %v", err)
		}
		files[gist.MetaFileName] = gist.File{
			Content: string(metaContent),
		}

		// 4. Gist 클라이언트 초기화
		client, err := gist.NewClient()
		if err != nil {
			return fmt.Errorf("Gist 클라이언트 초기화 실패: %v", err)
		}

		// 5. Gist 생성 또는 업데이트
		name := args[0]
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
	uploadCmd.Flags().BoolVarP(&preview, "preview", "p", false, "업로드 전 메타데이터 미리보기")
} 