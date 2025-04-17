package cmd

import (
	"fmt"
	"sort"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/choigawoon/rulesctl/internal/gist"
	"github.com/choigawoon/rulesctl/pkg/config"
	"github.com/spf13/cobra"
)

const (
	titleWidth = 25  // 제목 최대 너비
	dateWidth  = 19  // 날짜 너비
	idWidth    = 32  // Gist ID 너비
	separator  = "..."
)

// truncateString은 문자열을 지정된 너비로 자르고 필요한 경우 말줄임표를 추가합니다.
func truncateString(s string, width int) string {
	if utf8.RuneCountInString(s) <= width {
		return s + strings.Repeat(" ", width-utf8.RuneCountInString(s))
	}
	return string([]rune(s)[:width-len(separator)]) + separator
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "GIST에 저장된 규칙 목록 출력",
	Long: `GIST에 저장된 모든 규칙 목록을 출력합니다.
각 규칙은 [제목] [최종수정시각] [Gist ID] 형식으로 정렬되어 표시됩니다.
기본적으로 최근 1달 이내의 규칙만 표시됩니다.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		config, err := config.LoadConfig()
		if err != nil {
			return fmt.Errorf("설정을 로드할 수 없습니다: %w", err)
		}

		if config.Token == "" {
			return fmt.Errorf("GitHub 토큰이 설정되지 않았습니다. 'rulesctl auth' 명령어로 토큰을 설정해주세요.")
		}

		// 1달 전 날짜 계산
		oneMonthAgo := time.Now().AddDate(0, -1, 0)

		gists, err := gist.FetchUserGists(config.Token, oneMonthAgo)
		if err != nil {
			return fmt.Errorf("GIST 목록을 가져올 수 없습니다: %w", err)
		}

		// GIST를 최종 수정 시간 기준으로 정렬
		sort.Slice(gists, func(i, j int) bool {
			return gists[i].UpdatedAt.After(gists[j].UpdatedAt)
		})

		// 테이블 헤더 출력
		titleHeader := truncateString("제목", titleWidth)
		dateHeader := truncateString("최종 수정", dateWidth)
		idHeader := truncateString("Gist ID", idWidth)
		
		fmt.Printf("%s  %s  %s\n", titleHeader, dateHeader, idHeader)
		fmt.Println(strings.Repeat("-", titleWidth+dateWidth+idWidth+4))

		// 각 GIST 정보 출력
		for _, g := range gists {
			description := g.Description
			if description == "" {
				description = "(제목 없음)"
			}
			
			title := truncateString(description, titleWidth)
			date := truncateString(g.UpdatedAt.Format("2006-01-02 15:04:05"), dateWidth)
			id := truncateString(g.ID, idWidth)
			
			fmt.Printf("%s  %s  %s\n", title, date, id)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
} 