package gist

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
	
	"github.com/choigawoon/rulesctl/pkg/config"
)

var baseURL = "https://api.github.com"

const MetaFileName = ".rulesctl.meta.json"

// Gist는 GitHub Gist의 정보를 나타냅니다
type Gist struct {
	ID          string    `json:"id"`
	Description string    `json:"description"`
	Public      bool      `json:"public"`
	UpdatedAt   time.Time `json:"updated_at"`
	Files       map[string]struct {
		Filename string `json:"filename"`
		Type     string `json:"type"`
		Language string `json:"language"`
		RawURL   string `json:"raw_url"`
		Size     int    `json:"size"`
		Content  string `json:"content"`
	} `json:"files"`
	History []struct {
		Version   string    `json:"version"`
		CommitID  string    `json:"commit_id"`
		UpdatedAt time.Time `json:"updated_at"`
	} `json:"history"`
	RevisionNumber int // GitHub 웹 UI와 동일한 순번 (최신이 1)
}

// FetchUserGists는 사용자의 Gist를 가져옵니다
// since가 지정된 경우 해당 시간 이후의 Gist만 가져옵니다
func FetchUserGists(since *time.Time) ([]Gist, error) {
	// 설정에서 토큰 로드
	cfg, err := config.LoadConfig()
	if err != nil {
		return nil, fmt.Errorf("설정을 로드할 수 없습니다: %w", err)
	}
	
	if cfg.Token == "" {
		return nil, fmt.Errorf("GitHub 토큰이 설정되지 않았습니다")
	}
	
	client := &http.Client{}
	
	url := fmt.Sprintf("%s/gists", baseURL)
	if since != nil {
		url = fmt.Sprintf("%s?since=%s", url, since.Format(time.RFC3339))
	}
	
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("요청 생성 실패: %w", err)
	}

	req.Header.Set("Authorization", "token "+cfg.Token)
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("API 요청 실패: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API 요청 실패: %s", resp.Status)
	}

	var gists []Gist
	if err := json.NewDecoder(resp.Body).Decode(&gists); err != nil {
		return nil, fmt.Errorf("응답 파싱 실패: %w", err)
	}

	// .rulesctl.meta.json 파일이 있는 Gist만 필터링
	var rulesctlGists []Gist
	for _, g := range gists {
		if _, hasRulesctlMeta := g.Files[MetaFileName]; hasRulesctlMeta {
			rulesctlGists = append(rulesctlGists, g)
		}
	}

	return rulesctlGists, nil
}

// FetchGistWithHistory는 특정 Gist의 상세 정보와 히스토리를 가져옵니다
func FetchGistWithHistory(token, gistID string) (*Gist, error) {
	url := fmt.Sprintf("%s/gists/%s", baseURL, gistID)
	
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	
	req.Header.Set("Accept", "application/vnd.github.v3+json")
	req.Header.Set("Authorization", fmt.Sprintf("token %s", token))
	
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API 요청 실패: %s", resp.Status)
	}
	
	var gist Gist
	if err := json.NewDecoder(resp.Body).Decode(&gist); err != nil {
		return nil, err
	}

	// History 길이를 기준으로 RevisionNumber 설정
	// 최신이 Rev 1이 되도록 함
	gist.RevisionNumber = len(gist.History)
	
	return &gist, nil
}

// DeleteGist는 지정된 ID의 Gist를 삭제합니다.
func DeleteGist(gistID string) error {
	// 설정에서 토큰 로드
	cfg, err := config.LoadConfig()
	if err != nil {
		return fmt.Errorf("설정을 로드할 수 없습니다: %w", err)
	}
	
	if cfg.Token == "" {
		return fmt.Errorf("GitHub 토큰이 설정되지 않았습니다")
	}
	
	client := &http.Client{}
	url := fmt.Sprintf("%s/gists/%s", baseURL, gistID)
	
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return fmt.Errorf("요청 생성 실패: %w", err)
	}
	
	req.Header.Set("Authorization", "token "+cfg.Token)
	req.Header.Set("Accept", "application/vnd.github.v3+json")
	
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("API 요청 실패: %w", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("Gist 삭제 실패: %s", resp.Status)
	}
	
	return nil
} 