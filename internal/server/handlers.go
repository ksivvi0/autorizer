package server

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (s *Server) pingHandler(c *gin.Context) {
	s.services.LoggerService.WriteNotice(getRequestInfo(c))
	c.JSON(http.StatusOK, map[string]string{"ping": "pong"})
}

func (s *Server) generateTokensHandler(c *gin.Context) {
	s.services.LoggerService.WriteNotice(getRequestInfo(c))
	c.AbortWithStatus(http.StatusNotImplemented)
}

func (s *Server) refreshTokensHandler(c *gin.Context) {
	s.services.LoggerService.WriteNotice(getRequestInfo(c))
	c.AbortWithStatus(http.StatusNotImplemented)
}

func (s *Server) errorResponder(c *gin.Context, statusCode int, err error) {
	s.services.LoggerService.WriteNotice(getRequestInfo(c))
	c.AbortWithStatusJSON(statusCode, err.Error())
}

//func (s *Server) getResponseObject(key string, value interface{}) (*bytes.Buffer, error) {
//	tmp := new(bytes.Buffer)
//	ro := make(map[string]interface{})
//	ro[key] = value
//
//	if err := json.NewEncoder(tmp).Encode(ro); err != nil {
//		return nil, err
//	}
//	return tmp, nil
//}

func getRequestInfo(c *gin.Context) string {
	return fmt.Sprintf("%v: %v",
		c.Request.Header.Get("Host"),
		c.Request.Header.Get("Method"),
	)
}
