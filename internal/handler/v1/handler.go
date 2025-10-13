package v1

import (
	"net/http"

	"github.com/airsss993/work-svc/internal/config"
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

func (h *Handler) Init(api *gin.RouterGroup) {
	api.GET("/ping", h.ping)
	v1 := api.Group("/v1")
	{
		search := v1.Group("/search")
		{
			search.POST("/students", h.searchStudents)
		}

		groups := v1.Group("/groups")
		{
			groups.GET("/it", h.getITGroups)
			groups.GET("/:groupName/students", h.getGroupStudents)
		}

		repository := v1.Group("/repos")
		{
			repository.GET("/:user_id/contents", h.getRepoContent)
		}
	}

	repos := api.Group("/repos/:owner/:repo/branches/:ref")
	{
		repos.GET("/html/*filepath", h.getHTMLProxy)
		repos.GET("/raw/*filepath", h.getRawProxy)
	}
}

func (h *Handler) ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
		"version": "v1",
	})
}
