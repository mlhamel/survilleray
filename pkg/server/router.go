package server

import (
	ginzerolog "github.com/dn365/gin-zerolog"
	"github.com/gin-gonic/gin"
	"github.com/mlhamel/survilleray/pkg/config"
)

func NewRouter(cfg *config.Config) *gin.Engine {
	router := gin.New()
	router.Use(ginzerolog.Logger("gin"))
	router.Use(gin.Recovery())

	health := HealthController{}
	graphql := GraphQLController{cfg}

	router.GET("/health", health.Status)
	router.POST("/graphql", graphql.Query)

	return router
}
