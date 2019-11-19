package server

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/graphql-go/graphql"
	"github.com/mlhamel/survilleray/pkg/models"
	"github.com/mlhamel/survilleray/pkg/runtime"
)

type GraphQLController struct {
	context *runtime.Context
}

var pointType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Point",
		Fields: graphql.Fields{
			"icao24": &graphql.Field{
				Type: graphql.String,
			},
		},
	},
)

func queryType(context *runtime.Context) *graphql.Object {
	return graphql.NewObject(
		graphql.ObjectConfig{
			Name: "Query",
			Fields: graphql.Fields{
				"points": &graphql.Field{
					Type:        graphql.NewList(pointType),
					Description: "Get point list",
					Resolve: func(params graphql.ResolveParams) (interface{}, error) {
						var points []models.Point
						context.Database().Find(&points)

						return points, nil
					},
				},
			},
		},
	)
}

func mutationType(context *runtime.Context) *graphql.Object {
	return graphql.NewObject(
		graphql.ObjectConfig{
			Name:   "Mutation",
			Fields: graphql.Fields{},
		},
	)
}

func buildSchema(context *runtime.Context) (graphql.Schema, error) {
	return graphql.NewSchema(
		graphql.SchemaConfig{
			Query:    queryType(context),
			Mutation: mutationType(context),
		},
	)
}

func executeQuery(query string, context *runtime.Context, schema graphql.Schema) *graphql.Result {
	s, _ := buildSchema(context)
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
	s, _ := buildSchema(g.context)
	result := executeQuery(c.Query("query"), g.context, s)
	c.JSON(http.StatusOK, result)
}
