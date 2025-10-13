package service

import (
	"context"

	"github.com/airsss993/work-svc/internal/client"
	"github.com/airsss993/work-svc/internal/config"
	"github.com/airsss993/work-svc/internal/domain"
)

type Services struct {
	StudentService StudentService
	GroupService   GroupService
	StudentService   StudentService
	GitBucketService RepositoryService
}
type StudentService interface {
	SearchStudents(ctx context.Context, query string) ([]domain.StudentInfo, error)
}

type RepositoryService interface {
	GetRepositoryContent(ctx context.Context, userID string, path string) (domain.RepoContent, error)
}

type Repositories struct {
}

type Deps struct {
	Repos     *Repositories
	GitClient *client.GitBucketClient
	Config    *config.Config
}

func NewServices(deps Deps) *Services {
	return &Services{
		StudentService: NewStudentService(deps.Config, &deps.Config.App),
		GroupService:   NewGroupService(deps.Config, &deps.Config.App),
	studentService := NewStudentService(deps.Config, &deps.Config.App)
	gitBucketService := NewGitBucketService(deps.GitClient, deps.Config)

	return &Services{
		StudentService:   studentService,
		GitBucketService: gitBucketService,
	}
}
