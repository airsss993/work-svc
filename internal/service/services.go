package service

import (
	"github.com/airsss993/work-svc/internal/config"
)

type Services struct {
	StudentService StudentService
}

type Repositories struct {
}

type Deps struct {
	Repos  *Repositories
	Config *config.Config
}

func NewServices(deps Deps) *Services {
	return &Services{}
}
