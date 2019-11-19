package app

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/mlhamel/survilleray/pkg/runtime"
	"github.com/mlhamel/survilleray/pkg/server"
)

type ServerApp struct {
	context *runtime.Context
	router  *gin.Engine
}

func NewServerApp(context *runtime.Context) *ServerApp {
	return &ServerApp{
		context: context,
		router:  server.NewRouter(),
	}
}

func (s *ServerApp) Run() error {
	return s.router.Run(s.connexionString())
}

func (s *ServerApp) connexionString() string {
	return fmt.Sprintf(":%s", s.context.Config().HttpPort())
}
