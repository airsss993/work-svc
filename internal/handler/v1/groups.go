package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) getITGroups(c *gin.Context) {
	groups, err := h.services.GroupService.GetITGroups(c.Request.Context())

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"groups": groups,
		"total":  len(groups),
	})
}
