package service

import (
	"context"

	"github.com/airsss993/work-svc/internal/config"
	"github.com/airsss993/work-svc/internal/domain"
)

type Services struct {
	StudentService StudentService
	GroupService   GroupService
}

type RepositoryService interface {
	GetRepositoryContent(ctx context.Context, userID string, path string) (domain.RepoContent, error)
}

type Repositories struct {
}

type Deps struct {
	Repos  *Repositories
	Config *config.Config
}

func NewServices(deps Deps) *Services {
	return &Services{
		StudentService: NewStudentService(deps.Config, &deps.Config.App),
		GroupService:   NewGroupService(deps.Config, &deps.Config.App),
	}
}
