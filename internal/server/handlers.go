package server

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

type magicWordRequest struct {
	MagicWord string `json:"magic_word"`
}

func (s *Server) pingHandler(c *gin.Context) {
	c.JSON(http.StatusOK, map[string]string{"ping": "pong"})
}

func (s *Server) generateTokensHandler(c *gin.Context) {
	mwRequest := new(magicWordRequest)
	if err := c.ShouldBindJSON(mwRequest); err != nil {
		s.errorResponder(c, http.StatusBadRequest, err)
		return
	}

	if mwRequest.MagicWord != s.magicWord {
		s.errorResponder(c, http.StatusForbidden, errors.New("invalid magic word"))
		return
	}

	pair, err := s.services.AuthService.CreateTokenPair()
	if err != nil {
		s.errorResponder(c, http.StatusForbidden, err)
		return
	}

	s.services.StoreService.WriteTokensInfo(pair.RefreshToken)

	c.JSON(http.StatusOK, pair)
}

func (s *Server) refreshTokensHandler(c *gin.Context) {
	c.AbortWithStatus(http.StatusNotImplemented)
}

func (s *Server) errorResponder(c *gin.Context, statusCode int, err error) {
	s.services.WriteError(fmt.Sprintf("%v status: %d, error: %v", c.ClientIP(), statusCode, err))
	c.JSON(statusCode, gin.H{"error": err.Error()})
}
