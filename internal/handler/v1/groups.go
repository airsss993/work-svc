package v1

import (
	"net/http"

	"github.com/airsss993/work-svc/internal/domain"
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

	ctx := c.Request.Context()
	subgroup := c.Query("subgroup")

	var students []domain.Student
	var err error

	if subgroup != "" {
		students, err = h.services.GroupService.GetGroupStudentsFiltered(ctx, groupName, subgroup)
	} else {
		students, err = h.services.GroupService.GetGroupStudents(ctx, groupName)
	}

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

func (h *Handler) getGroupSubgroups(c *gin.Context) {
	groupName := c.Param("groupName")

	if groupName == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "group name is required",
		})
		return
	}

	subgroups, err := h.services.GroupService.GetGroupSubgroups(c.Request.Context(), groupName)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, subgroups)
}
