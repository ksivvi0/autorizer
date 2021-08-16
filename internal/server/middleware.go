package server

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func (s *Server) loggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		t := time.Now()

		c.Next()

		latency := time.Since(t)
		s.services.LoggerService.WriteNotice(fmt.Sprintf("%s: latency: %v", c.ClientIP(), latency))

		status := c.Writer.Status()
		s.services.LoggerService.WriteNotice(fmt.Sprintf("%s: status: %v", c.ClientIP(), status))
	}
}

func (s *Server) authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		bearerHeader := c.Request.Header.Get("Authorization")
		if len(bearerHeader) == 0 {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "empty authorization header"})
			c.Abort()
		}
		//todo: add parsing auth token
		c.Next()
	}
}
