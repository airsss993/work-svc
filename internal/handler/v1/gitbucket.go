package v1

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

const (
	DefaultPerPage = 30
	MaxPerPage     = 100
)

func (h *Handler) getRepoContent(c *gin.Context) {
	owner := c.Param("owner")
	if owner == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "owner is required",
		})
		return
	}

	repo := c.Param("repo")
	if repo == "" {
		repo = "Work"
	}

	path := c.Query("path")

	content, err := h.services.GitBucketService.GetRepositoryContent(c.Request.Context(), owner, repo, path)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, content)
}

func (h *Handler) getListCommits(c *gin.Context) {
	owner := c.Param("owner")
	if owner == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "owner is required"})
		return
	}

	repo := c.Param("repo")
	if repo == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "repo is required"})
		return
	}

	perPage := DefaultPerPage
	page := 1

	perPageStr := c.Query("per_page")
	pageStr := c.Query("page")

	if perPageStr != "" {
		if v, err := strconv.Atoi(perPageStr); err == nil && v > 0 {
			if v > MaxPerPage {
				perPage = MaxPerPage
			} else {
				perPage = v
			}
		}
	}

	if pageStr != "" {
		if v, err := strconv.Atoi(pageStr); err == nil && v > 0 {
			page = v
		}
	}

	commits, err := h.services.GitBucketService.GetCommitsList(
		c.Request.Context(),
		owner,
		repo,
		perPage,
		page,
	)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, commits)
}

func (h *Handler) getRepoContentWithDates(c *gin.Context) {
	owner := c.Param("owner")
	if owner == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "owner is required"})
		return
	}

	repo := c.Param("repo")
	if repo == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "repo is required"})
		return
	}

	path := c.Query("path")

	content, err := h.services.GitBucketService.GetRepositoryContentWithDates(
		c.Request.Context(),
		owner,
		repo,
		path,
	)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, content)
}

func (h *Handler) getUserRepositories(c *gin.Context) {
	owner := c.Param("owner")
	if owner == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "owner is required",
		})
		return
	}

	repositories, err := h.services.GitBucketService.GetUserRepositories(c.Request.Context(), owner)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, repositories)
}
