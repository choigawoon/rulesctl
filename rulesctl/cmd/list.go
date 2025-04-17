package cmd

import (
	"fmt"
	"sort"

	"github.com/choigawoon/rulesctl/internal/gist"
	"github.com/choigawoon/rulesctl/pkg/config"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "GIST에 저장된 규칙 목록 출력",
	Long: `GIST에 저장된 모든 규칙 목록을 출력합니다.
각 규칙은 [최종수정시각] 제목 형식으로 정렬되어 표시됩니다.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		config, err := config.LoadConfig()
		if err != nil {
			return fmt.Errorf("설정을 로드할 수 없습니다: %w", err)
		}

		if config.Token == "" {
			return fmt.Errorf("GitHub 토큰이 설정되지 않았습니다. 'rulesctl auth' 명령어로 토큰을 설정해주세요.")
		}

		gists, err := gist.FetchUserGists(config.Token)
		if err != nil {
			return fmt.Errorf("GIST 목록을 가져올 수 없습니다: %w", err)
		}

		// GIST를 최종 수정 시간 기준으로 정렬
		sort.Slice(gists, func(i, j int) bool {
			return gists[i].UpdatedAt.After(gists[j].UpdatedAt)
		})

		// 테이블 헤더 출력
		fmt.Printf("%-30s %-20s %s\n", "제목", "최종 수정", "설명")
		fmt.Println("----------------------------------------------------------------")

		// 각 GIST 정보 출력
		for _, g := range gists {
			updatedAt := g.UpdatedAt.Format("2006-01-02 15:04:05")
			fmt.Printf("%-30s %-20s %s\n", g.Description, updatedAt, g.Description)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
} 