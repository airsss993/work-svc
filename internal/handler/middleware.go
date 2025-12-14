package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) corsMiddleware(c *gin.Context) {
	allowedOrigins := map[string]bool{
		h.cfg.App.WebURL:                      true,
		"http://localhost:5172":               true,
		"https://work.students.it-college.ru": true,
		"http://work.students.it-college.ru":  true,
	}

	origin := c.Request.Header.Get("Origin")

	if allowedOrigins[origin] {
		c.Header("Access-Control-Allow-Origin", origin)
	}

	c.Header("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
	c.Header("Access-Control-Allow-Credentials", "true")

	if c.Request.Method != "OPTIONS" {
		c.Next()
	} else {
		c.AbortWithStatus(http.StatusOK)
	}
}
