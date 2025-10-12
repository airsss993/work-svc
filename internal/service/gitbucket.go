package service

import (
	"context"

	"github.com/airsss993/work-svc/internal/config"
	"github.com/airsss993/work-svc/internal/domain"
)

type GitBucketService struct {
	cfg *config.Config
}

func NewGitBucketService(cfg *config.Config) *GitBucketService {
	return &GitBucketService{
		cfg: cfg,
	}
}

func (g *GitBucketService) GetRepositoryContent(ctx context.Context, userID string, path string) (domain.RepoContent, error) {
	return domain.RepoContent{}, nil
}
