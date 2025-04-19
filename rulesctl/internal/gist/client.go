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

func (c *Client) CreateOrUpdateGist(name string, files map[string]File, force bool, public bool) (string, error) {
	// Search for existing Gist
	gists, _, err := c.client.Gists.List(c.ctx, "", &github.GistListOptions{})
	if err != nil {
		return "", fmt.Errorf("failed to list gists: %v", err)
	}

	var existingGist *github.Gist
	for _, gist := range gists {
		if gist.Description != nil && *gist.Description == name {
			existingGist = gist
			break
		}
	}

	// Create Gist files
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
			return "", fmt.Errorf("Gist already exists. Use --force option to force update")
		}
		// Update Gist
		existingGist.Files = gistFiles
		updatedGist, _, err := c.client.Gists.Edit(c.ctx, *existingGist.ID, existingGist)
		if err != nil {
			return "", fmt.Errorf("failed to update Gist: %v", err)
		}
		return *updatedGist.ID, nil
	}

	// Create new Gist
	description := name
	newGist := &github.Gist{
		Description: &description,
		Public:      &public,
		Files:       gistFiles,
	}

	createdGist, _, err := c.client.Gists.Create(c.ctx, newGist)
	if err != nil {
		return "", fmt.Errorf("failed to create Gist: %v", err)
	}

	return *createdGist.ID, nil
}

// FetchUserGists fetches all Gists of the user
func (c *Client) FetchUserGists() ([]struct {
	ID          string
	Description string
	UpdatedAt   time.Time
}, error) {
	gists, _, err := c.client.Gists.List(c.ctx, "", &github.GistListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to fetch Gist list: %v", err)
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