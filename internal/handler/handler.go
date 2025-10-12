package handlers

import (
	"github.com/airsss993/work-svc/internal/config"
	v1 "github.com/airsss993/work-svc/internal/handler/v1"
	"github.com/airsss993/work-svc/internal/service"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	services *service.Services
	cfg      *config.Config
}

func NewHandler(services *service.Services, cfg *config.Config) *Handler {
	return &Handler{
		services: services,
		cfg:      cfg,
	}
}

func (h *Handler) Init() *gin.Engine {
	router := gin.New()

	router.Use(
		gin.Recovery(),
		gin.Logger(),
		h.corsMiddleware,
	)

	h.initAPI(router)

	return router
}

func (h *Handler) initAPI(router *gin.Engine) {
	handlerV1 := v1.NewHandler(h.services, h.cfg)
	api := router.Group("/api")
	handlerV1.Init(api)
}
