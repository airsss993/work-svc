package v1

import (
	"net/http"

	"github.com/airsss993/work-svc/internal/handler/dto"
	"github.com/gin-gonic/gin"
)

func (h *Handler) searchStudents(c *gin.Context) {
	var req dto.StudentSearchReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request body",
		})
		return
	}

	teachers, err := h.services.StudentService.SearchTeachers(
		c.Request.Context(),
		req.Query,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"students": teachers,
		"total":    len(teachers),
	})
}
