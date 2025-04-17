package gist

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/choigawoon/rulesctl/pkg/config"
)

const (
	baseURL = "https://api.github.com"
)

// Client는 GitHub Gist API 클라이언트를 나타냅니다
type Client struct {
	httpClient *http.Client
	token      string
}

// NewClient는 새로운 GitHub Gist API 클라이언트를 생성합니다
func NewClient() (*Client, error) {
	config, err := config.LoadConfig()
	if err != nil {
		return nil, fmt.Errorf("설정을 로드할 수 없습니다: %w", err)
	}

	if config.Token == "" {
		return nil, fmt.Errorf("GitHub 토큰이 설정되지 않았습니다. 'rulesctl auth' 명령어로 토큰을 설정하세요")
	}

	return &Client{
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
		token: config.Token,
	}, nil
}

// doRequest는 HTTP 요청을 보내고 응답을 처리합니다
func (c *Client) doRequest(ctx context.Context, method, path string, body interface{}) ([]byte, error) {
	var reqBody io.Reader
	if body != nil {
		jsonData, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("요청 본문을 JSON으로 변환할 수 없습니다: %w", err)
		}
		reqBody = bytes.NewBuffer(jsonData)
	}

	req, err := http.NewRequestWithContext(ctx, method, baseURL+path, reqBody)
	if err != nil {
		return nil, fmt.Errorf("HTTP 요청을 생성할 수 없습니다: %w", err)
	}

	req.Header.Set("Authorization", "token "+c.token)
	req.Header.Set("Accept", "application/vnd.github.v3+json")
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("HTTP 요청을 보낼 수 없습니다: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API 요청 실패 (상태 코드: %d): %s", resp.StatusCode, string(body))
	}

	return io.ReadAll(resp.Body)
}

// GistFile은 Gist의 파일을 나타냅니다
type GistFile struct {
	Content string `json:"content"`
}

// Gist는 GitHub Gist를 나타냅니다
type Gist struct {
	ID          string              `json:"id"`
	Description string              `json:"description"`
	Public      bool                `json:"public"`
	Files       map[string]GistFile `json:"files"`
}

// CreateGist는 새로운 Gist를 생성합니다
func (c *Client) CreateGist(ctx context.Context, description string, files map[string]string) (*Gist, error) {
	gist := &Gist{
		Description: description,
		Public:      false,
		Files:       make(map[string]GistFile),
	}

	for filename, content := range files {
		gist.Files[filename] = GistFile{Content: content}
	}

	data, err := c.doRequest(ctx, "POST", "/gists", gist)
	if err != nil {
		return nil, err
	}

	var result Gist
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("응답을 파싱할 수 없습니다: %w", err)
	}

	return &result, nil
}

// GetGist는 Gist를 조회합니다
func (c *Client) GetGist(ctx context.Context, id string) (*Gist, error) {
	data, err := c.doRequest(ctx, "GET", "/gists/"+id, nil)
	if err != nil {
		return nil, err
	}

	var result Gist
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("응답을 파싱할 수 없습니다: %w", err)
	}

	return &result, nil
}

// ListGists는 사용자의 Gist 목록을 조회합니다
func (c *Client) ListGists(ctx context.Context) ([]Gist, error) {
	data, err := c.doRequest(ctx, "GET", "/gists", nil)
	if err != nil {
		return nil, err
	}

	var result []Gist
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("응답을 파싱할 수 없습니다: %w", err)
	}

	return result, nil
}

// DeleteGist는 Gist를 삭제합니다
func (c *Client) DeleteGist(ctx context.Context, id string) error {
	_, err := c.doRequest(ctx, "DELETE", "/gists/"+id, nil)
	return err
} 