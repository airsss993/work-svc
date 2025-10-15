package service

import (
	"context"

	"github.com/airsss993/work-svc/internal/client"
	"github.com/airsss993/work-svc/internal/config"
	"github.com/airsss993/work-svc/internal/domain"
)

type Services struct {
	StudentService   StudentService
	GroupService     GroupService
	GitBucketService RepositoryService
	ProxyService     *ProxyService
}

type RepositoryService interface {
	GetRepositoryContent(ctx context.Context, userID string, path string) (domain.RepoContent, error)
	GetCommitsList(ctx context.Context, owner, repo string, perPage, page int) (domain.CommitResp, error)
	GetRepositoryContentWithDates(ctx context.Context, owner, repo, path string) (domain.RepoContent, error)
}

type Repositories struct {
}

type Deps struct {
	Repos     *Repositories
	GitClient *client.GitBucketClient
	Config    *config.Config
}

func NewServices(deps Deps) *Services {
	studentService := NewStudentService(deps.Config, &deps.Config.App)
	groupService := NewGroupService(deps.Config, &deps.Config.App)
	gitBucketService := NewGitBucketService(deps.GitClient, deps.Config)
	proxyService := NewProxyService(deps.GitClient)

	return &Services{
		StudentService:   studentService,
		GroupService:     groupService,
		GitBucketService: gitBucketService,
		ProxyService:     proxyService,
	}
}
