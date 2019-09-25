package server

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/graphql-go/graphql"
	"github.com/mlhamel/survilleray/pkg/config"
	"github.com/mlhamel/survilleray/pkg/models"
)

type GraphQLController struct {
	Config *config.Config
}

var vectorType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Vector",
		Fields: graphql.Fields{
			"icao24": &graphql.Field{
				Type: graphql.String,
			},
		},
	},
)

func queryType(c *config.Config) *graphql.Object {
	return graphql.NewObject(
		graphql.ObjectConfig{
			Name: "Query",
			Fields: graphql.Fields{
				"vectors": &graphql.Field{
					Type:        graphql.NewList(vectorType),
					Description: "Get vector list",
					Resolve: func(params graphql.ResolveParams) (interface{}, error) {
						var vectors []models.Vector
						c.DB().Find(&vectors)

						return vectors, nil
					},
				},
			},
		},
	)
}

func mutationType(c *config.Config) *graphql.Object {
	return graphql.NewObject(
		graphql.ObjectConfig{
			Name:   "Mutation",
			Fields: graphql.Fields{},
		},
	)
}

func buildSchema(c *config.Config) (graphql.Schema, error) {
	return graphql.NewSchema(
		graphql.SchemaConfig{
			Query:    queryType(c),
			Mutation: mutationType(c),
		},
	)
}

func executeQuery(query string, c *config.Config, schema graphql.Schema) *graphql.Result {
	s, _ := buildSchema(c)
	result := graphql.Do(graphql.Params{
		Schema:        s,
		RequestString: query,
	})
	if len(result.Errors) > 0 {
		fmt.Printf("errors: %v", result.Errors)
	}
	return result
}

func (g GraphQLController) Query(c *gin.Context) {
	s, _ := buildSchema(g.Config)
	result := executeQuery(c.Query("query"), g.Config, s)
	c.JSON(http.StatusOK, result)
}