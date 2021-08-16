package server

import (
	"authorized/internal/services"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type Server struct {
	services  *services.Services
	address   string
	magicWord string
	engine    *http.Server
}

func NewServer(addr string, s *services.Services, debug bool) *Server {
	if !debug {
		gin.SetMode(gin.ReleaseMode)
	}
	return &Server{
		magicWord: "sudo", // :)
		services:  s,
		engine: &http.Server{
			Addr:         addr,
			ReadTimeout:  time.Second * 10,
			WriteTimeout: time.Second * 10,
		},
	}
}

func (s *Server) initRoutes() {
	router := gin.Default()
	router.Use(s.loggerMiddleware())
	router.GET("/ping", s.pingHandler)
	router.POST("/tokens", s.generateTokensHandler)
	s.engine.Handler = router
}

func (s *Server) Run(addr string) error {
	s.initRoutes()
	if addr == "" {
		return errors.New("invalid address string")
	}
	return s.engine.ListenAndServe()
}
