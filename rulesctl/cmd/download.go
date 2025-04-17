package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/choigawoon/rulesctl/internal/gist"
	"github.com/choigawoon/rulesctl/pkg/config"
)

var (
	gistID string
)

var downloadCmd = &cobra.Command{
	Use:   "download [제목]",
	Short: "GIST에서 규칙 다운로드",
	Long: `제목 또는 Gist ID로 규칙을 .cursor/rules 디렉토리로 다운로드합니다.
기존 파일이 있는 경우 --force 옵션으로 덮어쓸 수 있습니다.

사용 예시:
  # 제목으로 다운로드 (내 Gist에서 검색)
  rulesctl download "파이썬 린팅 규칙"
  rulesctl download "파이썬 린팅 규칙" --force

  # Gist ID로 다운로드 (공개 Gist)
  rulesctl download --gistid abc123
  rulesctl download --gistid abc123 --force`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		// 설정 로드
		cfg, err := config.LoadConfig()
		if err != nil {
			return fmt.Errorf("설정을 로드할 수 없습니다: %w", err)
		}

		if cfg.Token == "" {
			cmd.SilenceUsage = true
			return fmt.Errorf("GitHub 토큰이 설정되지 않았습니다. 'rulesctl auth' 명령어로 토큰을 설정해주세요")
		}

		var targetGistID string

		if gistID != "" {
			// Gist ID로 다운로드
			targetGistID = gistID
		} else {
			// 제목으로 다운로드
			if len(args) == 0 {
				cmd.SilenceUsage = true
				return fmt.Errorf("제목 또는 --gistid 옵션을 지정해주세요")
			}
			title := args[0]

			// 전체 Gist 목록 조회
			gists, err := gist.FetchUserGists(nil)
			if err != nil {
				cmd.SilenceUsage = true
				return fmt.Errorf("Gist 목록 조회 실패: %w", err)
			}

			// 제목으로 Gist 찾기
			found := false
			for _, g := range gists {
				if g.Description == title {
					targetGistID = g.ID
					found = true
					break
				}
			}

			if !found {
				cmd.SilenceUsage = true
				return fmt.Errorf("제목과 일치하는 Gist를 찾을 수 없습니다: %s", title)
			}
		}

		// Gist 가져오기
		g, err := gist.FetchGist(cfg.Token, targetGistID)
		if err != nil {
			cmd.SilenceUsage = true
			return fmt.Errorf("Gist를 가져올 수 없습니다: %w", err)
		}

		// .rulesctl.meta.json 파일 확인
		metaFile, exists := g.Files[gist.MetaFileName]
		if !exists {
			cmd.SilenceUsage = true
			return fmt.Errorf("이 Gist는 rulesctl로 관리되지 않습니다 (메타데이터 파일이 없음)")
		}

		// 메타데이터 파싱
		meta, err := gist.ParseMetadataFromGist(metaFile.Content)
		if err != nil {
			cmd.SilenceUsage = true
			return fmt.Errorf("메타데이터 파싱 실패: %w", err)
		}

		// 파일 충돌 검사
		if !force {
			conflicts, err := gist.CheckConflicts(meta)
			if err != nil {
				cmd.SilenceUsage = true
				return fmt.Errorf("충돌 검사 실패: %w", err)
			}
			if len(conflicts) > 0 {
				fmt.Println("다음 파일들이 이미 존재합니다:")
				for _, path := range conflicts {
					fmt.Printf("  - %s\n", path)
				}
				cmd.SilenceUsage = true
				return fmt.Errorf("파일 충돌이 발생했습니다. --force 옵션을 사용하여 덮어쓸 수 있습니다")
			}
		}

		// 파일 다운로드
		fmt.Printf("규칙 다운로드 중... (Gist ID: %s)\n", targetGistID)
		if err := gist.DownloadFiles(cfg.Token, targetGistID, meta, force); err != nil {
			return fmt.Errorf("다운로드 실패: %w", err)
		}

		fmt.Println("다운로드가 완료되었습니다.")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(downloadCmd)
	downloadCmd.Flags().StringVar(&gistID, "gistid", "", "다운로드할 Gist의 ID")
} 