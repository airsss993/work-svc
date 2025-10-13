package app

import (
	"fmt"

	"github.com/airsss993/work-svc/internal/client"
	"github.com/airsss993/work-svc/internal/config"
	handlers "github.com/airsss993/work-svc/internal/handler"
	"github.com/airsss993/work-svc/internal/server"
	"github.com/airsss993/work-svc/internal/service"
	"github.com/airsss993/work-svc/pkg/logger"
)

func Run() {
	cfg, err := config.Init()
	if err != nil {
		logger.Fatal(err)
	}

	gitClient := client.NewGitBucketClient(cfg)
	services := service.NewServices(service.Deps{
		Repos:     nil,
		GitClient: gitClient,
		Config:    cfg,
	})

	handler := handlers.NewHandler(services, cfg)

	router := handler.Init()

	srv := server.NewServer(cfg, router)

	logger.Info(fmt.Sprintf("College Work Service started - PORT: %s", cfg.Server.Port))

	if err := srv.Run(); err != nil {
		logger.Fatal(err)
	}
}
