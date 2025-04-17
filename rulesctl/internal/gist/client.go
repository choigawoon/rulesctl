package gist

import (
	"context"
	"fmt"
	"time"

	"github.com/google/go-github/v58/github"
)

type File struct {
	Content string
}

type Client struct {
	client *github.Client
	ctx    context.Context
}

func NewClient() (*Client, error) {
	token, err := getToken()
	if err != nil {
		return nil, err
	}

	ctx := context.Background()
	client := github.NewTokenClient(ctx, token)

	return &Client{
		client: client,
		ctx:    ctx,
	}, nil
}

func (c *Client) CreateOrUpdateGist(name string, files map[string]File, force bool) (string, error) {
	// 기존 Gist 검색
	gists, _, err := c.client.Gists.List(c.ctx, "", &github.GistListOptions{})
	if err != nil {
		return "", fmt.Errorf("Gist 목록 조회 실패: %v", err)
	}

	var existingGist *github.Gist
	for _, gist := range gists {
		if gist.Description != nil && *gist.Description == name {
			existingGist = gist
			break
		}
	}

	// Gist 파일 생성
	gistFiles := make(map[github.GistFilename]github.GistFile)
	for path, file := range files {
		filename := github.GistFilename(path)
		content := file.Content
		gistFiles[filename] = github.GistFile{
			Content: &content,
		}
	}

	if existingGist != nil {
		if !force {
			return "", fmt.Errorf("이미 존재하는 Gist입니다. --force 옵션을 사용하여 강제 업데이트하세요")
		}
		// Gist 업데이트
		existingGist.Files = gistFiles
		updatedGist, _, err := c.client.Gists.Edit(c.ctx, *existingGist.ID, existingGist)
		if err != nil {
			return "", fmt.Errorf("Gist 업데이트 실패: %v", err)
		}
		return *updatedGist.ID, nil
	}

	// 새 Gist 생성
	description := name
	public := false
	newGist := &github.Gist{
		Description: &description,
		Public:      &public,
		Files:       gistFiles,
	}

	createdGist, _, err := c.client.Gists.Create(c.ctx, newGist)
	if err != nil {
		return "", fmt.Errorf("Gist 생성 실패: %v", err)
	}

	return *createdGist.ID, nil
}

// FetchUserGists는 사용자의 모든 Gist를 가져옵니다
func (c *Client) FetchUserGists() ([]struct {
	ID          string
	Description string
	UpdatedAt   time.Time
}, error) {
	gists, _, err := c.client.Gists.List(c.ctx, "", &github.GistListOptions{})
	if err != nil {
		return nil, fmt.Errorf("Gist 목록 조회 실패: %v", err)
	}

	result := make([]struct {
		ID          string
		Description string
		UpdatedAt   time.Time
	}, 0, len(gists))

	for _, gist := range gists {
		if gist.ID == nil || gist.Description == nil {
			continue
		}
		result = append(result, struct {
			ID          string
			Description string
			UpdatedAt   time.Time
		}{
			ID:          *gist.ID,
			Description: *gist.Description,
			UpdatedAt:   gist.UpdatedAt.Time,
		})
	}

	return result, nil
} 