package service

import (
	"context"
	"fmt"

	"github.com/airsss993/work-svc/internal/client"
	"github.com/airsss993/work-svc/internal/config"
	"github.com/airsss993/work-svc/internal/domain"
)

type GitBucketService struct {
	cfg             *config.Config
	gitbucketClient *client.GitBucketClient
}

func NewGitBucketService(gitbucketClient *client.GitBucketClient, cfg *config.Config) *GitBucketService {
	return &GitBucketService{
		gitbucketClient: gitbucketClient,
		cfg:             cfg,
	}
}

func (g *GitBucketService) GetRepositoryContent(ctx context.Context, owner, path string) (domain.RepoContent, error) {
	if owner == "" {
		return domain.RepoContent{}, fmt.Errorf("owner is required")
	}

	content, err := g.gitbucketClient.GetRepositoryContent(ctx, owner, path)
	if err != nil {
		return domain.RepoContent{}, fmt.Errorf("failed to get repository content: %w", err)
	}

	var repoItems []domain.RepoContentItem

	for _, item := range content {
		if item.Name == ".DS_Store" {
			continue
		}

		repoItem := domain.RepoContentItem{
			Type:        item.Type,
			Name:        item.Name,
			Path:        item.Path,
			DownloadURL: item.DownloadURL,
		}

		repoItems = append(repoItems, repoItem)
	}

	repoContent := domain.RepoContent{
		Items: repoItems,
	}

	return repoContent, nil
}
