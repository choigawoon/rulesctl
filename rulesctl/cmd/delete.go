package cmd

import (
	"fmt"
	"strings"

	"github.com/choigawoon/rulesctl/internal/gist"
	"github.com/spf13/cobra"
)

var deleteCmd = &cobra.Command{
	Use:           "delete [name]",
	Short:         "규칙 세트 삭제",
	Long: `GitHub Gist에 저장된 규칙 세트를 삭제합니다.
규칙 세트는 제목으로 검색하여 삭제합니다.
최근 1달 이내에 업로드된 규칙 세트만 삭제할 수 있습니다.

사용 예시:
  rulesctl delete "my-python-ruleset"    # 제목으로 검색하여 삭제
  rulesctl delete "my-ruleset" --force   # 확인 없이 바로 삭제`,
	Args:          cobra.ExactArgs(1),
	SilenceErrors: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		force, _ := cmd.Flags().GetBool("force")
		title := args[0]

		// 전체 Gist 목록 조회
		gists, err := gist.FetchUserGists(nil)
		if err != nil {
			return fmt.Errorf("Gist 목록 조회 실패: %w", err)
		}

		// 제목으로 Gist 찾기
		var targetGist gist.Gist
		found := false
		for _, g := range gists {
			if g.Description == title {
				targetGist = g
				found = true
				break
			}
		}

		if !found {
			cmd.SilenceUsage = true
			return fmt.Errorf("규칙 세트를 찾을 수 없습니다: %s", title)
		}

		// 삭제 전 확인
		if !force {
			fmt.Printf("'%s' 규칙 세트를 삭제하시겠습니까? (y/N): ", title)
			var response string
			fmt.Scanln(&response)
			if !strings.EqualFold(response, "y") {
				fmt.Println("삭제가 취소되었습니다.")
				return nil
			}
		}

		if err := gist.DeleteGist(targetGist.ID); err != nil {
			cmd.SilenceUsage = true
			return fmt.Errorf("Gist 삭제 실패: %w", err)
		}

		fmt.Printf("'%s' 규칙 세트가 삭제되었습니다.\n", title)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)
	deleteCmd.Flags().Bool("force", false, "확인 없이 바로 삭제")
} 