package server

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (s *Server) pingHandler(c *gin.Context) {
	c.JSON(http.StatusOK, map[string]string{"ping": "pong"})
}

func (s *Server) generateTokensHandler(c *gin.Context) {
	c.AbortWithStatus(http.StatusNotImplemented)
}

func (s *Server) refreshTokensHandler(c *gin.Context) {
	c.AbortWithStatus(http.StatusNotImplemented)
}

func (s *Server) errorResponder(c *gin.Context, statusCode int, err error) {
	c.AbortWithStatusJSON(statusCode, err.Error())
}

func (s *Server) getResponseObject(key string, value interface{}) (*bytes.Buffer, error) {
	tmp := new(bytes.Buffer)
	ro := make(map[string]interface{})
	ro[key] = value

	if err := json.NewEncoder(tmp).Encode(ro); err != nil {
		return nil, err
	}
	return tmp, nil
}
