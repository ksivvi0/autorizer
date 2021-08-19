package server

import (
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func (s *Server) loggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		rStart := time.Now()
		c.Next()

		latency := time.Since(rStart)
		responseStatus := c.Writer.Status()
		s.services.LoggerService.WriteNotice(fmt.Sprintf("%s: status: %v, latency: %v", c.ClientIP(), responseStatus, latency))
	}
}

func (s *Server) authMiddleware(ctx context.Context) gin.HandlerFunc {
	return func(c *gin.Context) {
		bearerHeader := c.Request.Header.Get("Authorization")
		if len(bearerHeader) == 0 {
			s.errorResponder(c, http.StatusUnauthorized, errors.New("empty authorization header"))
			c.Abort()
			return
		}

		uid, err := s.services.AuthService.ValidateToken(bearerHeader, false)
		if err != nil {
			s.errorResponder(c, http.StatusForbidden, err)
			c.Abort()
			return
		}
		_, err = s.services.StoreService.GetTokensInfo(ctx, uid)
		if err != nil {
			s.errorResponder(c, http.StatusForbidden, err)
			c.Abort()
			return
		}

		//c.Set("tokenPair", pair)
		c.Next()
	}
}
