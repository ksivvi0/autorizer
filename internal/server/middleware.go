package server

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
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
			s.errorResponder(c, http.StatusUnauthorized, errors.New("empty authorization header"))
			c.Abort()
		}

		bearerHeaderArr := strings.Split(bearerHeader, " ")
		if len(bearerHeaderArr) != 2 || bearerHeaderArr[0] != "Bearer" {
			s.errorResponder(c, http.StatusUnauthorized, errors.New("invalid authorization header"))
			c.Abort()
		}

		token := bearerHeaderArr[1]
		if token == "" {
			s.errorResponder(c, http.StatusUnauthorized, errors.New("empty authorization token"))
			c.Abort()
		}

		_uuid, err := s.services.AuthService.ParseToken(token)
		if err != nil {
			s.errorResponder(c, http.StatusForbidden, err)
			c.Abort()
		}
		c.Set("uuid", _uuid)
		c.Next()
	}
}
