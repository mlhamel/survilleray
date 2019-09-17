package server

import (
	"github.com/gin-gonic/gin"
	"github.com/mlhamel/survilleray/pkg/config"
)

func NewRouter() *gin.Engine {
	c := config.NewConfig()
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	health := HealthController{}
	graphql := GraphQLController{Config: c}

	router.GET("/health", health.Status)
	router.POST("/graphql", graphql.Query)

	return router
}
