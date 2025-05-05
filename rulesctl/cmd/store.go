package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/choigawoon/rulesctl/internal/fileutils"
	"github.com/choigawoon/rulesctl/internal/gist"
	"github.com/spf13/cobra"
)

// StoreItem은 public-store.json의 각 항목을 표현하는 구조체
type StoreItem struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	GistID      string `json:"gist_id"`
	Source      string `json:"source"`
	Category    string `json:"category"`
}

// storeCmd는 'rulesctl store' 명령어 그룹
var storeCmd = &cobra.Command{
	Use:   "store",
	Short: "Manage and use public rule store",
	Long: `Manage and use public rule store.
The store contains curated cursor rules for various technology stacks.`,
}

// storeListCmd는 'rulesctl store list' 명령어
var storeListCmd = &cobra.Command{
	Use:   "list",
	Short: "List available rules in the store",
	Long: `List all available rules in the public store.
Shows name, description, category, and full Gist ID for each rule.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// 1. GitHub에서 public-store.json 다운로드
		const remoteURL = "https://raw.githubusercontent.com/choigawoon/rulesctl/main/public-store.json"
		jsonPath := filepath.Join("public-store.json")

		remoteData, err := fileutils.DownloadFileFromURL(remoteURL)
		if err != nil {
			fmt.Printf("[경고] 원격 스토어 목록을 내려받지 못했습니다: %v\n", err)
		}

		// 2. 로컬 파일과 해시 비교
		localData, readErr := os.ReadFile(jsonPath)
		updateNeeded := false
		if err == nil && readErr == nil {
			remoteHash := fileutils.CalculateMD5FromBytes(remoteData)
			localHash := fileutils.CalculateMD5FromBytes(localData)
			if remoteHash != localHash {
				updateNeeded = true
			}
		} else if err == nil && readErr != nil {
			// 로컬 파일이 없으면 무조건 갱신
			updateNeeded = true
		}

		if err == nil && updateNeeded {
			err := os.WriteFile(jsonPath, remoteData, 0644)
			if err != nil {
				fmt.Printf("[경고] 최신 스토어 목록을 저장하지 못했습니다: %v\n", err)
			}
		}

		// 3. (최신) 로컬 파일 읽어서 출력
		file, err := os.Open(jsonPath)
		if err != nil {
			return fmt.Errorf("public-store.json 파일을 열 수 없습니다: %w", err)
		}
		defer file.Close()

		var storeItems []StoreItem
		decoder := json.NewDecoder(file)
		if err := decoder.Decode(&storeItems); err != nil {
			return fmt.Errorf("public-store.json 파싱 오류: %w", err)
		}

		if len(storeItems) == 0 {
			fmt.Println("등록된 스토어 항목이 없습니다.")
			return nil
		}

		// 4. 결과 출력 (nameWidth, descWidth는 유지, idWidth는 전체 표시)
		nameWidth := 25
		descWidth := 40
		categoryWidth := 10

		fmt.Printf("%-*s  %-*s  %-*s  %s\n", nameWidth, "Name", descWidth, "Description", categoryWidth, "Category", "Gist ID")
		fmt.Println(strings.Repeat("-", nameWidth+descWidth+categoryWidth+4+36)) // +36은 Gist ID 길이

		for _, item := range storeItems {
			name := item.Name
			if len(name) > nameWidth {
				name = name[:nameWidth-3] + "..."
			} else {
				name = fmt.Sprintf("%-*s", nameWidth, name)
			}

			desc := item.Description
			if len(desc) > descWidth {
				desc = desc[:descWidth-3] + "..."
			} else {
				desc = fmt.Sprintf("%-*s", descWidth, desc)
			}

			category := item.Category
			if len(category) > categoryWidth {
				category = category[:categoryWidth-3] + "..."
			} else {
				category = fmt.Sprintf("%-*s", categoryWidth, category)
			}

			// Gist ID는 전체 표시 (truncate 없이)
			fmt.Printf("%s  %s  %s  %s\n", name, desc, category, item.GistID)
		}

		return nil
	},
}

// storeDownloadCmd는 'rulesctl store download' 명령어
var storeDownloadCmd = &cobra.Command{
	Use:   "download [name]",
	Short: "Download rule by name from the store",
	Long: `Download rule by name from the store.
Finds the Gist ID from the store and downloads it.

Example:
  rulesctl store download fastapi-patrickjs`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]

		// 1. public-store.json 읽기
		jsonPath := filepath.Join("public-store.json")
		file, err := os.Open(jsonPath)
		if os.IsNotExist(err) {
			// 파일이 없으면 원격에서 다운로드 시도
			const remoteURL = "https://raw.githubusercontent.com/choigawoon/rulesctl/main/public-store.json"
			remoteData, dlErr := fileutils.DownloadFileFromURL(remoteURL)
			if dlErr != nil {
				return fmt.Errorf("스토어 목록을 가져올 수 없습니다: %w", dlErr)
			}
			
			if err := os.WriteFile(jsonPath, remoteData, 0644); err != nil {
				return fmt.Errorf("스토어 목록을 저장할 수 없습니다: %w", err)
			}
			
			file, err = os.Open(jsonPath)
			if err != nil {
				return fmt.Errorf("public-store.json 파일을 열 수 없습니다: %w", err)
			}
		} else if err != nil {
			return fmt.Errorf("public-store.json 파일을 열 수 없습니다: %w", err)
		}
		defer file.Close()

		// 2. 항목 검색
		var storeItems []StoreItem
		decoder := json.NewDecoder(file)
		if err := decoder.Decode(&storeItems); err != nil {
			return fmt.Errorf("public-store.json 파싱 오류: %w", err)
		}

		// 이름으로 검색
		var targetGistID string
		var targetName string
		for _, item := range storeItems {
			if item.Name == name {
				targetGistID = item.GistID
				targetName = item.Name
				break
			}
		}

		if targetGistID == "" {
			return fmt.Errorf("스토어에서 '%s' 항목을 찾을 수 없습니다", name)
		}

		// 3. 다운로드 실행 (gist ID로)
		fmt.Printf("'%s' 룰셋을 다운로드합니다. (Gist ID: %s)\n", targetName, targetGistID)

		// Fetch Gist
		g, err := gist.FetchGist("", targetGistID) // 공개 Gist는 토큰 필요 없음
		if err != nil {
			return fmt.Errorf("Gist를 가져오지 못했습니다: %w", err)
		}

		// Check .rulesctl.meta.json file
		metaFile, exists := g.Files[gist.MetaFileName]
		if !exists {
			return fmt.Errorf("이 Gist는 rulesctl로 관리되지 않습니다 (메타데이터 없음)")
		}

		// Parse metadata
		meta, err := gist.ParseMetadataFromGist(metaFile.Content)
		if err != nil {
			return fmt.Errorf("메타데이터 파싱 오류: %w", err)
		}

		// Check for file conflicts
		forceDownload, _ := cmd.Flags().GetBool("force")
		if !forceDownload {
			conflicts, err := gist.CheckConflicts(meta)
			if err != nil {
				return fmt.Errorf("충돌 확인 오류: %w", err)
			}
			if len(conflicts) > 0 {
				fmt.Println("다음 파일이 이미 존재합니다:")
				for _, path := range conflicts {
					fmt.Printf("  - %s\n", path)
				}
				return fmt.Errorf("파일 충돌이 감지되었습니다. --force 옵션을 사용하여 덮어쓰기 가능합니다")
			}
		}

		// Download files
		if err := gist.DownloadFiles("", targetGistID, meta, forceDownload); err != nil {
			return fmt.Errorf("다운로드 실패: %w", err)
		}

		fmt.Println("다운로드가 성공적으로 완료되었습니다.")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(storeCmd)
	storeCmd.AddCommand(storeListCmd)
	storeCmd.AddCommand(storeDownloadCmd)
	
	// 다운로드 시 force 옵션 추가
	storeDownloadCmd.Flags().Bool("force", false, "Force overwrite if files already exist")
} 