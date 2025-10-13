package v1

import (
	"fmt"
	"net/http"

	"github.com/airsss993/work-svc/internal/utils"
	"github.com/gin-gonic/gin"
)

func (h *Handler) getHTMLProxy(c *gin.Context) {
	owner := c.Param("owner")
	repo := c.Param("repo")
	ref := c.Param("ref")
	filePath := c.Param("filepath")

	if owner == "" || repo == "" || ref == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "owner, repo, and ref are required",
		})
		return
	}

	if len(filePath) > 0 && filePath[0] == '/' {
		filePath = filePath[1:]
	}

	if filePath == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "filepath is required",
		})
		return
	}

	baseURL := fmt.Sprintf("/api/repos/%s/%s/branches/%s/raw", owner, repo, ref)

	content, err := h.services.ProxyService.GetHTMLWithBase(
		c.Request.Context(),
		owner,
		repo,
		ref,
		filePath,
		baseURL,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("failed to get HTML: %v", err),
		})
		return
	}

	c.Data(http.StatusOK, "text/html; charset=utf-8", content)
}

func (h *Handler) getRawProxy(c *gin.Context) {
	owner := c.Param("owner")
	repo := c.Param("repo")
	ref := c.Param("ref")
	filePath := c.Param("filepath")

	if owner == "" || repo == "" || ref == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "owner, repo, and ref are required",
		})
		return
	}

	if len(filePath) > 0 && filePath[0] == '/' {
		filePath = filePath[1:]
	}

	if filePath == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "filepath is required",
		})
		return
	}

	content, err := h.services.ProxyService.GetRawFile(
		c.Request.Context(),
		owner,
		repo,
		ref,
		filePath,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("failed to get file: %v", err),
		})
		return
	}

	contentType := utils.GetContentType(filePath)

	c.Data(http.StatusOK, contentType, content)
}
