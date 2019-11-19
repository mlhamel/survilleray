package server

import (
	"github.com/gin-gonic/gin"
	"github.com/mlhamel/survilleray/pkg/config"
	"github.com/mlhamel/survilleray/pkg/runtime"
)

func NewRouter() *gin.Engine {
	cfg := config.NewConfig()
	context := runtime.NewContext(cfg, nil)
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	health := HealthController{}
	graphql := GraphQLController{context}

	router.GET("/health", health.Status)
	router.POST("/graphql", graphql.Query)

	return router
}
