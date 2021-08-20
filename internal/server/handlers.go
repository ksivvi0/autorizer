package server

import (
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type magicWordRequest struct {
	MagicWord string `json:"magic_word"`
}
type refreshTokenRequest struct {
	RefreshToken string `json:"refresh_token"`
	AccessToken  string `json:"access_token"`
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

	pair, err := s.services.AuthService.CreateTokenPair("")
	if err != nil {
		s.errorResponder(c, http.StatusForbidden, err)
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	id, err := s.services.StoreService.WriteTokensInfo(ctx, *pair)
	defer cancel()
	if err != nil {
		s.errorResponder(c, http.StatusInternalServerError, err)
		return
	}
	s.services.LoggerService.WriteNotice(fmt.Sprintf("created auth with %v", id))

	c.JSON(http.StatusOK, pair)
}

func (s *Server) refreshTokensHandler(c *gin.Context) {
	rTokenRequest := new(refreshTokenRequest)
	if err := c.ShouldBindJSON(rTokenRequest); err != nil {
		s.errorResponder(c, http.StatusBadRequest, err)
		return
	}

	ctx := context.Background()
	rUid, err := s.services.AuthService.GetDataFromToken(rTokenRequest.RefreshToken, true)
	if err != nil {
		s.errorResponder(c, http.StatusForbidden, err)
		return
	}

	tPairs, err := s.services.StoreService.GetTokensInfo(ctx, "refresh_token_uid", rUid)
	if err != nil {
		s.errorResponder(c, http.StatusForbidden, err)
		return
	}

	refreshTokenExists := false
	for _, v := range tPairs {
		if v.AccessUID != "" {
			dropId, err := s.services.StoreService.DropTokensInfo(ctx, "access_token_uid", v.AccessUID)
			if err != nil {
				s.errorResponder(c, http.StatusForbidden, err)
				return
			}
			s.services.LoggerService.WriteNotice(fmt.Sprintf("Drop %d access token with uid %s", dropId, v.AccessUID))
		}
		if v.RefreshToken != "" {
			refreshTokenExists = true
		}
	}
	if !refreshTokenExists {
		s.errorResponder(c, http.StatusForbidden, errors.New("refresh token not found"))
		return
	}

	pair, err := s.services.AuthService.CreateTokenPair(rUid)
	if err != nil {
		s.errorResponder(c, http.StatusForbidden, err)
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	_, err = s.services.StoreService.WriteTokensInfo(ctx, *pair)
	defer cancel()
	if err != nil {
		s.errorResponder(c, http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, pair)
}

func (s *Server) errorResponder(c *gin.Context, statusCode int, err error) {
	s.services.WriteError(fmt.Sprintf("%v status: %d, error: %v", c.ClientIP(), statusCode, err))
	c.JSON(statusCode, gin.H{"error": err.Error()})
}
