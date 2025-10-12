package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) getRepoContent(c *gin.Context) {
	owner := c.Param("owner")
	if owner == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "owner is required",
		})
		return
	}

	path := c.Query("path")

	content, err := h.services.GitBucketService.GetRepositoryContent(c.Request.Context(), owner, path)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, content)
}
