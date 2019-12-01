package server

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/graphql-go/graphql"
	"github.com/mlhamel/survilleray/models"
	"github.com/mlhamel/survilleray/pkg/config"
)

type GraphQLController struct {
	cfg *config.Config
}

func NewGraphQLController(cfg *config.Config) *GraphQLController {
	return &GraphQLController{cfg}
}

func (controller *GraphQLController) Query(c *gin.Context) {
	s, _ := controller.buildSchema()
	result := controller.executeQuery(c.Query("query"), s)
	c.JSON(http.StatusOK, result)
}

func (controller *GraphQLController) pointType() *graphql.Object {
	return graphql.NewObject(graphql.ObjectConfig{
		Name: "Point",
		Fields: graphql.Fields{
			"icao24": &graphql.Field{
				Type: graphql.String,
			},
		},
	})
}

func (controller *GraphQLController) queryType() *graphql.Object {
	return graphql.NewObject(
		graphql.ObjectConfig{
			Name: "Query",
			Fields: graphql.Fields{
				"points": &graphql.Field{
					Type:        graphql.NewList(controller.pointType()),
					Description: "Get point list",
					Resolve: func(params graphql.ResolveParams) (interface{}, error) {
						var points []models.Point
						controller.cfg.Database().Find(&points)

						return points, nil
					},
				},
			},
		},
	)
}

func (controller *GraphQLController) mutationType() *graphql.Object {
	return graphql.NewObject(
		graphql.ObjectConfig{
			Name:   "Mutation",
			Fields: graphql.Fields{},
		},
	)
}

func (controller *GraphQLController) buildSchema() (graphql.Schema, error) {
	return graphql.NewSchema(
		graphql.SchemaConfig{
			Query:    controller.queryType(),
			Mutation: controller.mutationType(),
		},
	)
}

func (controller *GraphQLController) executeQuery(query string, schema graphql.Schema) *graphql.Result {
	s, _ := controller.buildSchema()
	result := graphql.Do(graphql.Params{
		Schema:        s,
		RequestString: query,
	})
	if len(result.Errors) > 0 {
		fmt.Printf("errors: %v", result.Errors)
	}
	return result
}
