package app

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/mlhamel/survilleray/pkg/config"
	"github.com/mlhamel/survilleray/pkg/server"
)

type ServerApp struct {
	cfg    *config.Config
	router *gin.Engine
}

func NewServerApp(cfg *config.Config) *ServerApp {
	return &ServerApp{
		cfg:    cfg,
		router: server.NewRouter(),
	}
}

func (s *ServerApp) Run() error {
	return s.router.Run(s.connexionString())
}

func (s *ServerApp) connexionString() string {
	return fmt.Sprintf(":%s", s.cfg.HttpPort())
}