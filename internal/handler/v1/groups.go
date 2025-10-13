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

func (h *Handler) getGroupStudents(c *gin.Context) {
	groupName := c.Param("groupName")

	if groupName == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "group name is required",
		})
		return
	}

	students, err := h.services.GroupService.GetGroupStudents(c.Request.Context(), groupName)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"students": students,
		"total":    len(students),
	})
}
