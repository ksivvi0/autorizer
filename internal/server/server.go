package server

import (
	"authorized/internal/services"
	"errors"
	"github.com/gin-gonic/gin"
)

type Server struct {
	services  *services.Services
	engine    *gin.Engine
	address   string
	magicWord string
}

func NewServer() *Server {
	return &Server{
		engine:    gin.Default(),
		magicWord: "sudo", // :)
	}
}

func (s *Server) initRoutes() {
	s.engine.GET("/ping", s.pingHandler)
	s.engine.POST("/tokens", s.generateTokensHandler)
}

func (s *Server) Run(addr string) error {
	s.initRoutes()
	if addr == "" {
		return errors.New("invalid address string")
	}
	return s.engine.Run(addr)
}
