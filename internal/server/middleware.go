package server

import (
	"fmt"
	"github.com/gin-gonic/gin"
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
