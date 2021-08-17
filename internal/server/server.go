package server

import (
	"authorizer/internal/services"
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

func NewServerInstance(addr string, s *services.Services, debug bool) (*Server, error) {
	if len(addr) == 0 {
		return nil, errors.New("invalid address")
	}
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
	}, nil
}

func (s *Server) initRoutes() {
	router := gin.Default()
	router.Use(s.loggerMiddleware())

	api := router.Group("/api/")
	api.Use(s.authMiddleware())
	{
		api.GET("/ping", s.pingHandler)
	}
	router.POST("/tokens", s.generateTokensHandler)
	s.engine.Handler = router
}

func (s *Server) Run() error {
	s.initRoutes()

	return s.engine.ListenAndServe()
}
