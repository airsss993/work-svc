package service

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/airsss993/work-svc/internal/client"
	"github.com/airsss993/work-svc/internal/config"
	"github.com/airsss993/work-svc/internal/domain"
	"github.com/airsss993/work-svc/pkg/logger"
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

func (g *GitBucketService) GetRepositoryContentWithDates(ctx context.Context, owner, repo, path string) (domain.RepoContent, error) {
	repoContent, err := g.gitbucketClient.GetRepositoryContent(ctx, owner, path)
	if err != nil {
		logger.Error(fmt.Errorf("GetRepositoryContent failed for %s/%s path=%s: %w", owner, repo, path, err))
		return domain.RepoContent{}, fmt.Errorf("failed to get repo content: %w", err)
	}

	repoCommits, err := g.gitbucketClient.GetCommitsList(ctx, owner, repo, 500, 1)
	if err != nil {
		logger.Error(fmt.Errorf("GetCommitsList failed for %s/%s: %w", owner, repo, err))
		return domain.RepoContent{}, fmt.Errorf("failed to get repo commits: %w", err)
	}

	lastModified := map[string]time.Time{}

	for _, commit := range repoCommits {
		for _, file := range commit.Files {
			filename := file.Filename
			date := commit.Commit.Author.Date

			if prev, ok := lastModified[filename]; !ok || date.After(prev) {
				lastModified[filename] = *date
			}
		}
	}

	var repoContentResp domain.RepoContent
	var repoContentItems []domain.RepoContentItem

	for _, c := range repoContent {
		if c.Name == ".DS_Store" {
			continue
		}

		var last time.Time

		if c.Type == "dir" {
			for fileName, fileTime := range lastModified {
				if strings.HasPrefix(fileName, c.Path+"/") {
					if fileTime.After(last) {
						last = fileTime
					}
				}
			}
		} else {
			if t, ok := lastModified[c.Path]; ok {
				last = t
			}
		}

		file := domain.RepoContentItem{
			Type:         c.Type,
			Name:         c.Name,
			Path:         c.Path,
			DownloadURL:  c.DownloadURL,
			LastModified: last,
		}

		repoContentItems = append(repoContentItems, file)
	}

	repoContentResp.Items = repoContentItems

	return repoContentResp, nil
}

func (g *GitBucketService) GetCommitsList(ctx context.Context, owner, repo string, perPage, page int) (domain.CommitResp, error) {
	if owner == "" {
		return domain.CommitResp{}, fmt.Errorf("owner is required")
	}

	commits, err := g.gitbucketClient.GetCommitsList(ctx, owner, repo, perPage, page)
	if err != nil {
		return domain.CommitResp{}, fmt.Errorf("failed to get commit list: %w", err)
	}

	var result domain.CommitResp

	result.Count = len(commits)
	result.Page = page
	result.PerPage = perPage
	result.HasNext = len(commits) == perPage

	result.Commits = make([]domain.Commit, 0, len(commits))

	for _, c := range commits {
		files := make([]domain.CommitFile, 0, len(c.Files))
		for _, f := range c.Files {
			files = append(files, domain.CommitFile{
				Filename: f.Filename,
			})
		}

		commitShort := domain.Commit{
			Message: c.Commit.Message,
			Date:    c.Commit.Author.Date,
			Files:   files,
		}

		result.Commits = append(result.Commits, commitShort)
	}

	return result, nil
}
